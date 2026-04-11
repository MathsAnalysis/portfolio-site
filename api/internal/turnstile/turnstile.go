// Package turnstile verifies Cloudflare Turnstile tokens server-side.
//
// Turnstile is CF's privacy-friendly CAPTCHA alternative. The widget renders
// on the form page, generates a single-use token, and the browser submits it.
// The backend MUST re-verify the token with CF before trusting the submission.
//
// Doc: https://developers.cloudflare.com/turnstile/get-started/server-side-validation/
package turnstile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const verifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type Verifier struct {
	SecretKey string
	client    *http.Client
}

// New returns a Verifier. If secretKey is empty, Verify() is a no-op that
// always accepts (for Phase 1/local dev, where the widget isn't configured yet).
func New(secretKey string) *Verifier {
	return &Verifier{
		SecretKey: secretKey,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

// Enabled reports whether the verifier has a secret configured.
func (v *Verifier) Enabled() bool {
	return v != nil && v.SecretKey != ""
}

type verifyResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Action      string   `json:"action"`
}

// Verify returns nil on a good token. Returns a descriptive error otherwise.
// When Enabled() is false, it returns nil (bypass).
func (v *Verifier) Verify(ctx context.Context, token, remoteIP string) error {
	if !v.Enabled() {
		return nil
	}
	if token == "" {
		return errors.New("missing turnstile token")
	}

	form := url.Values{}
	form.Set("secret", v.SecretKey)
	form.Set("response", token)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL,
		strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("turnstile call: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("turnstile status %d", resp.StatusCode)
	}

	var out verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("turnstile decode: %w", err)
	}
	if !out.Success {
		return fmt.Errorf("turnstile rejected: %v", out.ErrorCodes)
	}
	return nil
}
