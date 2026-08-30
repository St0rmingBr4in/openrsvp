package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed magic_link.html rsvp_confirmation.html event_reminder.html retention_warning.html organizer_rsvp_notification.html feedback_confirmation.html rsvp_lookup.html waitlist_promotion.html cohost_invitation.html
var templateFS embed.FS

var (
	magicLinkTmpl               *template.Template
	rsvpConfirmationTmpl        *template.Template
	eventReminderTmpl           *template.Template
	retentionWarningTmpl        *template.Template
	organizerRSVPNotifyTmpl     *template.Template
	feedbackConfirmationTmpl    *template.Template
	rsvpLookupTmpl              *template.Template
	waitlistPromotionTmpl       *template.Template
	cohostInvitationTmpl        *template.Template
)

func init() {
	// "t" exposes the localized message catalog to templates as {{(t).Field}}.
	fm := template.FuncMap{"t": Locale}
	parse := func(name string) *template.Template {
		return template.Must(template.New(name).Funcs(fm).ParseFS(templateFS, name))
	}
	magicLinkTmpl = parse("magic_link.html")
	rsvpConfirmationTmpl = parse("rsvp_confirmation.html")
	eventReminderTmpl = parse("event_reminder.html")
	retentionWarningTmpl = parse("retention_warning.html")
	organizerRSVPNotifyTmpl = parse("organizer_rsvp_notification.html")
	feedbackConfirmationTmpl = parse("feedback_confirmation.html")
	rsvpLookupTmpl = parse("rsvp_lookup.html")
	waitlistPromotionTmpl = parse("waitlist_promotion.html")
	cohostInvitationTmpl = parse("cohost_invitation.html")
}

// magicLinkData holds the template data for a magic link email.
type magicLinkData struct {
	URL           string
	ExpiryMinutes int
	Colors        EmailColors
}

// rsvpConfirmationData holds the template data for an RSVP confirmation email.
type rsvpConfirmationData struct {
	EventTitle string
	EventDate  string
	Location   string
	RSVPStatus string
	ModifyURL  string
	Colors     EmailColors
}

// eventReminderData holds the template data for an event reminder email.
type eventReminderData struct {
	EventTitle string
	EventDate  string
	Location   string
	Message    string
	InviteURL  string
	Colors     EmailColors
}

// retentionWarningData holds the template data for a retention warning email.
type retentionWarningData struct {
	EventTitle   string
	ExpiresAt    string
	DashboardURL string
	Colors       EmailColors
}

// organizerRSVPNotificationData holds the template data for notifying an
// organizer about a new or updated RSVP.
type organizerRSVPNotificationData struct {
	EventTitle   string
	GuestName    string
	RSVPStatus   string
	GuestEmail   string
	GuestPhone   string
	PlusOnes     int
	DashboardURL string
	Colors       EmailColors
}

// displayStatus returns a human-friendly, localized label for an RSVP status.
func displayStatus(status string) string {
	m := Locale()
	switch status {
	case "attending":
		return m.StatusAttending
	case "maybe":
		return m.StatusMaybe
	case "declined":
		return m.StatusDeclined
	case "pending":
		return m.StatusPending
	case "waitlisted":
		return m.StatusWaitlisted
	default:
		return status
	}
}

// RenderMagicLink renders the magic link email template and returns the HTML
// body and a plain text fallback.
func RenderMagicLink(baseURL, token string, expiryMinutes int) (html, plain string, err error) {
	url := fmt.Sprintf("%s/auth/verify?token=%s", baseURL, token)

	data := magicLinkData{
		URL:           url,
		ExpiryMinutes: expiryMinutes,
		Colors:        DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := magicLinkTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render magic link template: %w", err)
	}

	m := Locale()
	plainText := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s\n\n%s",
		m.MagicPlainTitle, m.MagicPlainBody, url, fmt.Sprintf(m.MagicPlainExpiry, expiryMinutes), m.MagicPlainIgnore,
	)

	return buf.String(), plainText, nil
}

// RenderRSVPConfirmation renders the RSVP confirmation email template and
// returns the HTML body and a plain text fallback.
func RenderRSVPConfirmation(eventTitle, eventDate, location, rsvpStatus, modifyURL string) (html, plain string, err error) {
	label := displayStatus(rsvpStatus)
	data := rsvpConfirmationData{
		EventTitle: eventTitle,
		EventDate:  eventDate,
		Location:   location,
		RSVPStatus: label,
		ModifyURL:  modifyURL,
		Colors:     DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := rsvpConfirmationTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render rsvp confirmation template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.RSVPConfPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEvent, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelDate, eventDate))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelLocation, location))
	sb.WriteString(fmt.Sprintf("%s: %s\n\n", m.LabelYourRSVP, label))
	sb.WriteString(fmt.Sprintf("%s\n%s\n", m.ModifyButton, modifyURL))

	return buf.String(), sb.String(), nil
}

// RenderEventReminder renders the event reminder email template and returns
// the HTML body and a plain text fallback.
func RenderEventReminder(eventTitle, eventDate, location, message, inviteURL string) (html, plain string, err error) {
	data := eventReminderData{
		EventTitle: eventTitle,
		EventDate:  eventDate,
		Location:   location,
		Message:    message,
		InviteURL:  inviteURL,
		Colors:     DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := eventReminderTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render event reminder template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.ReminderPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEvent, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelDate, eventDate))
	sb.WriteString(fmt.Sprintf("%s: %s\n\n", m.LabelLocation, location))
	if message != "" {
		sb.WriteString(fmt.Sprintf("%s\n%s\n\n", m.ReminderPlainOrgMsg, message))
	}
	sb.WriteString(fmt.Sprintf("%s\n%s\n", m.ReminderPlainInvite, inviteURL))

	return buf.String(), sb.String(), nil
}

// RenderRetentionWarning renders the retention warning email template and
// returns the HTML body and a plain text fallback.
func RenderRetentionWarning(eventTitle, expiresAt, dashboardURL string) (html, plain string, err error) {
	data := retentionWarningData{
		EventTitle:   eventTitle,
		ExpiresAt:    expiresAt,
		DashboardURL: dashboardURL,
		Colors:       DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := retentionWarningTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render retention warning template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.RetentionPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf(m.RetentionPlainBody+"\n\n", eventTitle, expiresAt))
	sb.WriteString(m.RetentionPlainNote + "\n\n")
	if dashboardURL != "" {
		sb.WriteString(fmt.Sprintf("%s\n%s\n", m.RetentionPlainExtend, dashboardURL))
	}

	return buf.String(), sb.String(), nil
}

// RenderOrganizerRSVPNotification renders the organizer RSVP notification email
// and returns the HTML body and a plain text fallback.
func RenderOrganizerRSVPNotification(eventTitle, guestName, rsvpStatus, guestEmail, guestPhone string, plusOnes int, dashboardURL string) (html, plain string, err error) {
	label := displayStatus(rsvpStatus)
	data := organizerRSVPNotificationData{
		EventTitle:   eventTitle,
		GuestName:    guestName,
		RSVPStatus:   label,
		GuestEmail:   guestEmail,
		GuestPhone:   guestPhone,
		PlusOnes:     plusOnes,
		DashboardURL: dashboardURL,
		Colors:       DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := organizerRSVPNotifyTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render organizer rsvp notification template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.OrgNotifyPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEvent, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelGuest, guestName))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelResponse, label))
	if guestEmail != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEmail, guestEmail))
	}
	if guestPhone != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelPhone, guestPhone))
	}
	if plusOnes > 0 {
		sb.WriteString(fmt.Sprintf("%s: +%d\n", m.LabelAdditional, plusOnes))
	}
	sb.WriteString(fmt.Sprintf("\n%s\n%s\n", m.OrgNotifyPlainView, dashboardURL))

	return buf.String(), sb.String(), nil
}

// rsvpLookupData holds the template data for an RSVP lookup email.
type rsvpLookupData struct {
	EventTitle string
	ModifyURL  string
	Colors     EmailColors
}

// RenderRSVPLookup renders the RSVP lookup magic link email template and
// returns the HTML body and a plain text fallback.
func RenderRSVPLookup(eventTitle, modifyURL string) (html, plain string, err error) {
	data := rsvpLookupData{
		EventTitle: eventTitle,
		ModifyURL:  modifyURL,
		Colors:     DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := rsvpLookupTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render rsvp lookup template: %w", err)
	}

	m := Locale()
	plainText := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s",
		m.LookupPlainTitle, fmt.Sprintf(m.LookupPlainBody, eventTitle), modifyURL, m.LookupPlainNote,
	)

	return buf.String(), plainText, nil
}

// feedbackConfirmationData holds the template data for a feedback confirmation email.
type feedbackConfirmationData struct {
	FeedbackType  string
	AllowFollowUp bool
	Colors        EmailColors
}

// RenderFeedbackConfirmation renders the feedback confirmation email template
// and returns the HTML body and a plain text fallback.
func RenderFeedbackConfirmation(feedbackType string, allowFollowUp bool) (htmlBody, plain string, err error) {
	data := feedbackConfirmationData{
		FeedbackType:  feedbackType,
		AllowFollowUp: allowFollowUp,
		Colors:        DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := feedbackConfirmationTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render feedback confirmation template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.FeedbackPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf(m.FeedbackPlainBody+"\n\n", feedbackType))
	if allowFollowUp {
		sb.WriteString(m.FeedbackFollowUp + "\n\n")
	}
	sb.WriteString(m.FeedbackPlainOutro + "\n")

	return buf.String(), sb.String(), nil
}

// cohostInvitationData holds the template data for a co-host invitation email.
type cohostInvitationData struct {
	EventTitle   string
	EventDate    string
	Location     string
	AddedByName  string
	DashboardURL string
	Colors       EmailColors
}

// RenderCoHostInvitation renders the co-host invitation email template and
// returns the HTML body and a plain text fallback.
func RenderCoHostInvitation(eventTitle, eventDate, location, addedByName, dashboardURL string) (html, plain string, err error) {
	data := cohostInvitationData{
		EventTitle:   eventTitle,
		EventDate:    eventDate,
		Location:     location,
		AddedByName:  addedByName,
		DashboardURL: dashboardURL,
		Colors:       DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := cohostInvitationTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render cohost invitation template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.CohostPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf(m.CohostPlainBody+"\n\n", addedByName, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEvent, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelDate, eventDate))
	sb.WriteString(fmt.Sprintf("%s: %s\n\n", m.LabelLocation, location))
	sb.WriteString(fmt.Sprintf("%s\n%s\n", m.CohostPlainView, dashboardURL))

	return buf.String(), sb.String(), nil
}

// waitlistPromotionData holds the template data for a waitlist promotion email.
type waitlistPromotionData struct {
	EventTitle string
	EventDate  string
	Location   string
	ModifyURL  string
	Colors     EmailColors
}

// RenderWaitlistPromotion renders the waitlist promotion email template and
// returns the HTML body and a plain text fallback.
func RenderWaitlistPromotion(eventTitle, eventDate, location, modifyURL string) (html, plain string, err error) {
	data := waitlistPromotionData{
		EventTitle: eventTitle,
		EventDate:  eventDate,
		Location:   location,
		ModifyURL:  modifyURL,
		Colors:     DefaultEmailColors(),
	}

	var buf bytes.Buffer
	if err := waitlistPromotionTmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("render waitlist promotion template: %w", err)
	}

	m := Locale()
	var sb strings.Builder
	sb.WriteString(m.WaitlistPlainTitle + "\n\n")
	sb.WriteString(fmt.Sprintf(m.WaitlistPlainBody+"\n\n", eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelEvent, eventTitle))
	sb.WriteString(fmt.Sprintf("%s: %s\n", m.LabelDate, eventDate))
	sb.WriteString(fmt.Sprintf("%s: %s\n\n", m.LabelLocation, location))
	sb.WriteString(fmt.Sprintf("%s\n%s\n", m.WaitlistPlainView, modifyURL))

	return buf.String(), sb.String(), nil
}
