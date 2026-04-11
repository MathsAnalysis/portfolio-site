package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/MathsAnalysis/portfolio-api/internal/db"
)

// Webhook contains the webhook handlers (inbound email from CF Worker, etc).
type Webhook struct {
	Store         *db.TicketStore
	Log           *slog.Logger
	InboundSecret string // HMAC-SHA256 shared secret with the CF Worker
}

// inboundPayload is the JSON body the CF Email Worker POSTs to us.
type inboundPayload struct {
	From      string            `json:"from"`
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	Text      string            `json:"text"`
	HTML      string            `json:"html"`
	RawMIME   string            `json:"raw_mime,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	FromName  string            `json:"from_name,omitempty"`
}

// InboundEmail processes a parsed inbound email from the CF Worker.
//
// Auth: body is signed with HMAC-SHA256 using InboundSecret.
// The signature is sent in the header X-Webhook-Signature: sha256=<hex>.
//
// The CF Worker must:
//   1. parse the incoming email
//   2. construct the JSON payload
//   3. compute HMAC-SHA256 of the raw JSON bytes with the shared secret
//   4. POST with header X-Webhook-Signature: sha256=<hex>
func (h *Webhook) InboundEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if h.InboundSecret == "" {
		h.Log.Warn("inbound webhook called but secret not configured")
		http.Error(w, "inbound webhook disabled", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024)) // 1 MB max
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	if err := verifyHMAC(r.Header.Get("X-Webhook-Signature"), body, h.InboundSecret); err != nil {
		h.Log.Warn("inbound webhook: bad signature", "err", err)
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var p inboundPayload
	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	code, err := extractTicketCode(p.To)
	if err != nil {
		// Fallback: try to extract from subject "Re: ... [TK-XXXXX]"
		if c := codeFromSubject(p.Subject); c != "" {
			code = c
		} else {
			h.Log.Warn("inbound webhook: could not extract ticket code",
				"to", p.To,
				"subject", p.Subject,
			)
			writeJSON(w, http.StatusOK, map[string]any{
				"ok":      false,
				"matched": false,
				"reason":  "no ticket code in envelope",
			})
			return
		}
	}

	ticket, err := h.Store.FindByPublicCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			h.Log.Warn("inbound webhook: unknown ticket code", "code", code)
			writeJSON(w, http.StatusOK, map[string]any{
				"ok":      false,
				"matched": false,
				"reason":  "ticket code not found",
			})
			return
		}
		h.Log.Error("inbound webhook: lookup", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// Use plain text if present, else strip tags from HTML
	text := strings.TrimSpace(p.Text)
	if text == "" {
		text = strings.TrimSpace(stripHTML(p.HTML))
	}
	if text == "" {
		text = "(empty message body)"
	}

	_, err = h.Store.AddInboundMessage(
		r.Context(),
		ticket.ID,
		sanitizeAddress(p.From),
		p.FromName,
		text,
		p.HTML,
		p.RawMIME,
	)
	if err != nil {
		h.Log.Error("inbound webhook: save message", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	h.Log.Info("inbound reply stored",
		"ticket_id", ticket.ID,
		"ticket_code", ticket.PublicCode,
		"from", p.From,
	)

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":          true,
		"matched":     true,
		"ticket_code": ticket.PublicCode,
	})
}

// verifyHMAC compares the signature header "sha256=<hex>" to the hex-encoded
// HMAC-SHA256 of `body` keyed with `secret`. Constant-time comparison.
func verifyHMAC(header string, body []byte, secret string) error {
	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) {
		return errors.New("missing or malformed X-Webhook-Signature")
	}
	got, err := hex.DecodeString(header[len(prefix):])
	if err != nil {
		return errors.New("signature not hex")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	want := mac.Sum(nil)
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return errors.New("signature mismatch")
	}
	return nil
}

// extractTicketCode parses an address like "tickets+TK-ABCD12@mathsanalysis.com"
// and returns "TK-ABCD12".
func extractTicketCode(to string) (string, error) {
	// Handle both "Name <addr>" and bare address
	if i := strings.Index(to, "<"); i >= 0 {
		if j := strings.Index(to, ">"); j > i {
			to = to[i+1 : j]
		}
	}
	at := strings.Index(to, "@")
	if at < 0 {
		return "", errors.New("no @")
	}
	local := to[:at]
	plus := strings.Index(local, "+")
	if plus < 0 {
		return "", errors.New("no subaddress")
	}
	code := local[plus+1:]
	if code == "" {
		return "", errors.New("empty subaddress")
	}
	return strings.ToUpper(code), nil
}

var codeRE = regexp.MustCompile(`\bTK-[A-Z0-9]{4,10}\b`)

func codeFromSubject(s string) string {
	return codeRE.FindString(strings.ToUpper(s))
}

var tagRE = regexp.MustCompile(`(?is)<[^>]+>`)

func stripHTML(s string) string {
	return tagRE.ReplaceAllString(s, " ")
}

func sanitizeAddress(s string) string {
	if i := strings.Index(s, "<"); i >= 0 {
		if j := strings.Index(s, ">"); j > i {
			return strings.ToLower(strings.TrimSpace(s[i+1 : j]))
		}
	}
	return strings.ToLower(strings.TrimSpace(s))
}
