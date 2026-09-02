package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/yannkr/openrsvp/internal/notification"
)

// implicitTLSPort is the SMTPS convention: the server expects a TLS
// handshake as the very first bytes on the wire, rather than a plaintext
// greeting followed by an opportunistic STARTTLS upgrade. Several providers
// (OVH's MX Plan among them) only offer this on port 465, with no STARTTLS
// alternative. A var, not a const, so tests can point it at an unprivileged
// ephemeral port instead of binding the real (privileged) 465.
var implicitTLSPort = "465"

// SMTPProvider sends emails via SMTP.
type SMTPProvider struct {
	host     string
	port     string
	username string
	password string
	from     string
	// tlsConfig overrides the default {ServerName: host} TLS config when
	// non-nil. Production leaves this nil (system root CA pool); tests set
	// it to trust an in-memory self-signed certificate.
	tlsConfig *tls.Config
}

// NewSMTPProvider creates a new SMTPProvider with the given SMTP configuration.
func NewSMTPProvider(host, port, username, password, from string) *SMTPProvider {
	return &SMTPProvider{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Name returns the provider identifier.
func (p *SMTPProvider) Name() string {
	return "smtp"
}

// Channel returns which channel this provider serves.
func (p *SMTPProvider) Channel() notification.Channel {
	return notification.ChannelEmail
}

// Send composes a proper MIME email and delivers it via SMTP.
// When attachments are present, the structure is:
//
//	multipart/mixed
//	  multipart/alternative
//	    text/plain
//	    text/html
//	  attachment(s)
//
// Without attachments, the structure is just multipart/alternative.
func (p *SMTPProvider) Send(ctx context.Context, msg *notification.Message) (*notification.SendResult, error) {
	var buf bytes.Buffer
	// Defensive: strip CR/LF from header values to defeat header injection
	// even if upstream validation is bypassed. mime.QEncoding handles the
	// Subject separately by encoding non-printable bytes.
	buf.WriteString(fmt.Sprintf("From: %s\r\n", stripCRLF(p.from)))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", stripCRLF(msg.To)))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("utf-8", stripCRLF(msg.Subject))))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z)))

	altBoundary := fmt.Sprintf("==OpenRSVP-alt==%d==", time.Now().UnixNano())

	if len(msg.Attachments) > 0 {
		// Wrap in multipart/mixed when attachments are present.
		mixedBoundary := fmt.Sprintf("==OpenRSVP-mix==%d==", time.Now().UnixNano())
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", mixedBoundary))
		buf.WriteString("\r\n")

		// Start the alternative part inside mixed.
		buf.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", altBoundary))
		buf.WriteString("\r\n")

		writeAlternativeParts(&buf, msg, altBoundary)

		buf.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))

		// Write each attachment.
		for _, att := range msg.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", att.ContentType, att.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", att.Filename))
			buf.WriteString("\r\n")
			encoded := base64.StdEncoding.EncodeToString(att.Data)
			// Wrap base64 at 76 characters per line (RFC 2045).
			for i := 0; i < len(encoded); i += 76 {
				end := i + 76
				if end > len(encoded) {
					end = len(encoded)
				}
				buf.WriteString(encoded[i:end])
				buf.WriteString("\r\n")
			}
		}

		buf.WriteString(fmt.Sprintf("--%s--\r\n", mixedBoundary))
	} else {
		// No attachments — just multipart/alternative.
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", altBoundary))
		buf.WriteString("\r\n")

		writeAlternativeParts(&buf, msg, altBoundary)

		buf.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))
	}

	if err := p.sendMail(ctx, []string{msg.To}, buf.Bytes()); err != nil {
		return nil, fmt.Errorf("smtp send: %w", err)
	}

	return &notification.SendResult{}, nil
}

// tlsConfigOrDefault returns the injected TLS config for tests, or the
// production default of trusting the system root CA pool for p.host.
func (p *SMTPProvider) tlsConfigOrDefault() *tls.Config {
	if p.tlsConfig != nil {
		return p.tlsConfig
	}
	return &tls.Config{ServerName: p.host}
}

// dial connects to the SMTP server and completes the EHLO/HELO exchange,
// choosing implicit TLS or plaintext-with-opportunistic-STARTTLS based on
// the configured port. The returned client is ready for MAIL/RCPT/DATA (or
// QUIT for a health check); the caller must Close it.
func (p *SMTPProvider) dial(ctx context.Context) (*smtp.Client, error) {
	addr := net.JoinHostPort(p.host, p.port)
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	var conn net.Conn
	var err error
	if p.port == implicitTLSPort {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: p.tlsConfigOrDefault()}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("handshake: %w", err)
	}

	if err := client.Hello("localhost"); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("hello: %w", err)
	}

	// Already encrypted from the first byte on an implicit-TLS port; only
	// negotiate STARTTLS when the connection started out in plaintext.
	if p.port != implicitTLSPort {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(p.tlsConfigOrDefault()); err != nil {
				_ = client.Close()
				return nil, fmt.Errorf("starttls: %w", err)
			}
		}
	}

	return client, nil
}

// sendMail delivers one message over a freshly dialed connection.
func (p *SMTPProvider) sendMail(ctx context.Context, to []string, body []byte) error {
	client, err := p.dial(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	if p.username != "" && p.password != "" {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(smtp.PlainAuth("", p.username, p.password, p.host)); err != nil {
				return fmt.Errorf("auth: %w", err)
			}
		}
	}
	if err := client.Mail(p.from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("rcpt to: %w", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}
	return client.Quit()
}

// writeAlternativeParts writes the text/plain and text/html MIME parts inside
// a multipart/alternative boundary.
func writeAlternativeParts(buf *bytes.Buffer, msg *notification.Message, boundary string) {
	plain := msg.Plain
	if plain == "" {
		plain = msg.Body
	}

	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	buf.WriteString("\r\n")
	qpw := quotedprintable.NewWriter(buf)
	_, _ = qpw.Write([]byte(plain))
	_ = qpw.Close()
	buf.WriteString("\r\n")

	if msg.Body != "" && msg.Body != plain {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		buf.WriteString("\r\n")
		qpw = quotedprintable.NewWriter(buf)
		_, _ = qpw.Write([]byte(msg.Body))
		_ = qpw.Close()
		buf.WriteString("\r\n")
	}
}

// SendBatch delivers multiple notifications by iterating and sending each
// one individually.
func (p *SMTPProvider) SendBatch(ctx context.Context, msgs []*notification.Message) ([]*notification.SendResult, []error) {
	results := make([]*notification.SendResult, len(msgs))
	errs := make([]error, len(msgs))
	for i, msg := range msgs {
		results[i], errs[i] = p.Send(ctx, msg)
	}
	return results, errs
}

// stripCRLF removes carriage returns and line feeds from a header value to
// prevent SMTP header injection. We replace with empty string rather than
// space because legitimate header values do not contain raw CR/LF.
func stripCRLF(s string) string {
	if !strings.ContainsAny(s, "\r\n") {
		return s
	}
	r := strings.NewReplacer("\r", "", "\n", "")
	return r.Replace(s)
}

// HealthCheck dials the SMTP server to verify connectivity.
func (p *SMTPProvider) HealthCheck(ctx context.Context) error {
	client, err := p.dial(ctx)
	if err != nil {
		return fmt.Errorf("smtp health check: %w", err)
	}
	defer func() { _ = client.Close() }()

	return client.Quit()
}
