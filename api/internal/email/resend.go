package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const resendAPIURL = "https://api.resend.com/emails"

// ResendConfig configures outbound visitor replies.
type ResendConfig struct {
	APIKey      string // re_...
	FromAddress string // "Carlo <replies@mathsanalysis.com>"
	Domain      string // "mathsanalysis.com" — used to build per-ticket Reply-To
}

// ResendSender delivers messages via the Resend HTTPS API.
type ResendSender struct {
	cfg    ResendConfig
	log    *slog.Logger
	client *http.Client
}

// NewResendSender returns a no-op sender when APIKey is empty.
func NewResendSender(cfg ResendConfig, log *slog.Logger) (Sender, bool) {
	if cfg.APIKey == "" || cfg.FromAddress == "" || cfg.Domain == "" {
		log.Info("resend sender disabled",
			"reason", "missing RESEND_API_KEY, RESEND_FROM or RESEND_DOMAIN",
		)
		return NoopSender(), false
	}
	return &ResendSender{
		cfg:    cfg,
		log:    log,
		client: &http.Client{Timeout: 10 * time.Second},
	}, true
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html,omitempty"`
	ReplyTo string   `json:"reply_to,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type resendResponse struct {
	ID      string `json:"id"`
	Message string `json:"message,omitempty"`
}

func (s *ResendSender) SendReply(ctx context.Context, t Ticket, bodyText, fromName string) (string, error) {
	subject := "Re: " + t.Subject
	replyTo := fmt.Sprintf("tickets+%s@%s", t.PublicCode, s.cfg.Domain)

	from := s.cfg.FromAddress
	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, addressOnly(s.cfg.FromAddress))
	}

	payload := resendPayload{
		From:    from,
		To:      []string{t.Email},
		Subject: subject,
		Text:    bodyText,
		ReplyTo: replyTo,
		Headers: map[string]string{
			"X-Portfolio-Ticket": t.PublicCode,
		},
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("resend call: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 16*1024))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("resend status %d: %s", resp.StatusCode, string(body))
	}

	var out resendResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("resend decode: %w", err)
	}
	return out.ID, nil
}
