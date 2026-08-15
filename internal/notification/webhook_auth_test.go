package notification

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yannkr/openrsvp/internal/database"
	"github.com/yannkr/openrsvp/internal/testutil"
)

// testWebhookSecret is a 32-character secret, the shortest value config accepts.
const testWebhookSecret = "0123456789abcdef0123456789abcdef"

// sendGridComplaintBody is a spam complaint for one address. It reaches the
// global suppression path when the handler runs.
const sendGridComplaintBody = `[{"email":"victim@example.com","event":"spamreport"}]`

// deliveryStatus reads the delivery status of one notification log row.
func deliveryStatus(ctx context.Context, t *testing.T, db database.DB, logID string) string {
	t.Helper()
	var status string
	err := db.QueryRowContext(ctx, "SELECT delivery_status FROM notification_log WHERE id = ?", logID).Scan(&status)
	require.NoError(t, err)
	return status
}

// postWebhook sends a request through the notification router and returns the
// recorded response.
func postWebhook(h *Handler, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	return rec
}

func TestWebhookRoute_SecretUnset_NotFoundAndHandlerDoesNotRun(t *testing.T) {
	db := testutil.NewTestDB(t)
	fake := newFakeSuppression()
	h := newTestHandlerWithSecret(NewTrackingService(db, zerolog.Nop()), nil, fake, "")

	rec := postWebhook(h, "/webhooks/sendgrid", sendGridComplaintBody)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// A token cannot help while the secret is unset.
	rec = postWebhook(h, "/webhooks/sendgrid?token="+testWebhookSecret, sendGridComplaintBody)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// The handler never ran, so nothing was suppressed.
	assert.False(t, fake.suppressed["victim@example.com"])
	assert.Empty(t, fake.suppressed2)
}

func TestWebhookRoute_TokenMissing_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	fake := newFakeSuppression()
	h := newTestHandlerWithSecret(NewTrackingService(db, zerolog.Nop()), nil, fake, testWebhookSecret)

	rec := postWebhook(h, "/webhooks/sendgrid", sendGridComplaintBody)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	rec = postWebhook(h, "/webhooks/sendgrid?token=", sendGridComplaintBody)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	assert.False(t, fake.suppressed["victim@example.com"])
}

func TestWebhookRoute_TokenWrong_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	fake := newFakeSuppression()
	h := newTestHandlerWithSecret(NewTrackingService(db, zerolog.Nop()), nil, fake, testWebhookSecret)

	// Same length as the secret, different content.
	wrong := "0123456789abcdef0123456789abcdee"
	require.Len(t, wrong, len(testWebhookSecret))

	rec := postWebhook(h, "/webhooks/sendgrid?token="+wrong, sendGridComplaintBody)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// A prefix of the secret is not accepted either.
	rec = postWebhook(h, "/webhooks/ses?token="+testWebhookSecret[:16], `{"Type":"Notification","Message":"{}"}`)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	assert.False(t, fake.suppressed["victim@example.com"])
}

func TestWebhookRoute_TokenCorrect_ProcessesEvent(t *testing.T) {
	db := testutil.NewTestDB(t)
	ctx := context.Background()
	eventID := uuid.Must(uuid.NewV7()).String()
	attendeeID := uuid.Must(uuid.NewV7()).String()
	createParentRecordsForNotification(t, ctx, db, eventID, attendeeID)

	logID := uuid.Must(uuid.NewV7()).String()
	insertNotificationLog(t, ctx, db, logID, eventID, attendeeID, "sent", "delivered", "sg-gate-1")

	fake := newFakeSuppression()
	h := newTestHandlerWithSecret(NewTrackingService(db, zerolog.Nop()), nil, fake, testWebhookSecret)

	body := `[{"email":"bob@example.com","event":"bounce","type":"bounce","sg_message_id":"sg-gate-1","timestamp":1700000000}]`

	// The same event without the token must not reach the handler.
	rec := postWebhook(h, "/webhooks/sendgrid", body)
	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "delivered", deliveryStatus(ctx, t, db, logID))

	rec = postWebhook(h, "/webhooks/sendgrid?token="+testWebhookSecret, body)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "bounced", deliveryStatus(ctx, t, db, logID))
	assert.Equal(t, "bounce", fake.suppressed2["bob@example.com"])
}

func TestWebhookRoute_SESTokenCorrect_ProcessesEvent(t *testing.T) {
	db := testutil.NewTestDB(t)
	ctx := context.Background()
	eventID := uuid.Must(uuid.NewV7()).String()
	attendeeID := uuid.Must(uuid.NewV7()).String()
	createParentRecordsForNotification(t, ctx, db, eventID, attendeeID)

	logID := uuid.Must(uuid.NewV7()).String()
	insertNotificationLog(t, ctx, db, logID, eventID, attendeeID, "sent", "delivered", "ses-gate-1")

	fake := newFakeSuppression()
	h := newTestHandlerWithSecret(NewTrackingService(db, zerolog.Nop()), nil, fake, testWebhookSecret)

	body := `{"Type":"Notification","Message":"{\"notificationType\":\"Bounce\",\"mail\":{\"messageId\":\"ses-gate-1\"},\"bounce\":{\"bounceType\":\"Permanent\",\"bouncedRecipients\":[{\"emailAddress\":\"dave@example.com\"}]}}"}`

	// The same event without the token must not reach the handler.
	rec := postWebhook(h, "/webhooks/ses", body)
	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "delivered", deliveryStatus(ctx, t, db, logID))

	rec = postWebhook(h, "/webhooks/ses?token="+testWebhookSecret, body)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "bounced", deliveryStatus(ctx, t, db, logID))
	assert.Equal(t, "bounce", fake.suppressed2["dave@example.com"])
}

// TestWebhookRoute_TokenComparisonIsConstantTime guards the comparison itself.
// A behavioural test cannot tell a constant-time compare from a byte-by-byte
// one, so the test reads the source instead.
func TestWebhookRoute_TokenComparisonIsConstantTime(t *testing.T) {
	src, err := os.ReadFile("handler.go")
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(src), "subtle.ConstantTimeCompare"),
		"the webhook token must be compared with subtle.ConstantTimeCompare")
}
