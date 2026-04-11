package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// SMTPConfig configures the internal notification path via the axiom mailserver.
type SMTPConfig struct {
	Host       string // e.g. "axiom-mailserver"
	Port       int    // 587 STARTTLS, 25 plain (LAN-only)
	Username   string // optional
	Password   string // optional
	From       string // "Portfolio <noreply@mathsanalysis.com>"
	To         string // where notifications go: "carlo4340@outlook.it"
	UseTLS     bool   // true = implicit TLS (465), false = STARTTLS or plain
	SkipVerify bool   // true = accept self-signed certs on LAN
}

// SMTPNotifier sends "new ticket" notifications to Carlo.
type SMTPNotifier struct {
	cfg SMTPConfig
	log *slog.Logger
}

// NewSMTPNotifier returns nil (and true for "disabled") when cfg.Host or cfg.To is empty.
func NewSMTPNotifier(cfg SMTPConfig, log *slog.Logger) (Notifier, bool) {
	if cfg.Host == "" || cfg.To == "" || cfg.From == "" {
		log.Info("smtp notifier disabled",
			"reason", "missing SMTP_HOST, SMTP_FROM or SMTP_TO",
		)
		return NoopNotifier(), false
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &SMTPNotifier{cfg: cfg, log: log}, true
}

func (s *SMTPNotifier) NotifyNewTicket(ctx context.Context, t Ticket) error {
	subject := fmt.Sprintf("[%s] %s — %s", t.PublicCode, categoryLabel(t.Category), t.Subject)

	var body bytes.Buffer
	fmt.Fprintf(&body, "A new ticket has arrived on mathsanalysis.com\n\n")
	fmt.Fprintf(&body, "Code:      %s\n", t.PublicCode)
	fmt.Fprintf(&body, "From:      %s %s <%s>\n", t.FirstName, t.LastName, t.Email)
	fmt.Fprintf(&body, "Category:  %s\n", categoryLabel(t.Category))
	fmt.Fprintf(&body, "Subject:   %s\n", t.Subject)
	fmt.Fprintf(&body, "Received:  %s\n\n", t.CreatedAt.UTC().Format(time.RFC3339))
	fmt.Fprintf(&body, "────────────────────────\n%s\n────────────────────────\n\n", t.Message)
	fmt.Fprintf(&body, "Open in admin: https://mathsanalysis.com/admin/tickets/%s\n", t.ID)

	msg := buildMIME(s.cfg.From, s.cfg.To, subject, body.String())
	return s.send(ctx, []string{s.cfg.To}, msg)
}

func (s *SMTPNotifier) send(ctx context.Context, to []string, msg []byte) error {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	var conn net.Conn
	var err error
	if s.cfg.UseTLS {
		tlsCfg := &tls.Config{
			ServerName:         s.cfg.Host,
			InsecureSkipVerify: s.cfg.SkipVerify, //nolint:gosec
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsCfg)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Quit() //nolint:errcheck

	if ok, _ := client.Extension("STARTTLS"); ok && !s.cfg.UseTLS {
		tlsCfg := &tls.Config{
			ServerName:         s.cfg.Host,
			InsecureSkipVerify: s.cfg.SkipVerify, //nolint:gosec
		}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if s.cfg.Username != "" && s.cfg.Password != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	from := addressOnly(s.cfg.From)
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("rcpt %s: %w", rcpt, err)
		}
	}
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := wc.Write(msg); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("data close: %w", err)
	}
	return nil
}

func buildMIME(from, to, subject, body string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&b, "Content-Type: text/plain; charset=utf-8\r\n")
	fmt.Fprintf(&b, "Content-Transfer-Encoding: 8bit\r\n")
	fmt.Fprintf(&b, "\r\n")
	b.WriteString(body)
	return b.Bytes()
}

func addressOnly(s string) string {
	if i := strings.Index(s, "<"); i >= 0 {
		if j := strings.Index(s, ">"); j > i {
			return s[i+1 : j]
		}
	}
	return s
}

func categoryLabel(c string) string {
	switch c {
	case "general":
		return "General"
	case "job":
		return "Job opportunity"
	case "security":
		return "Security / Pentest"
	case "collaboration":
		return "Collaboration / OSS"
	case "other":
		return "Other"
	}
	return c
}

var _ = errors.New // keep import
