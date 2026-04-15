// Package session implements HMAC-signed session cookies.
//
// Cookie value format:  <email>|<exp-unix>|<base64url(hmac)>
//
// We deliberately don't use a random session ID + server store — there is
// exactly one admin, login state is simple, and a stateless signed cookie
// avoids any shared DB state for auth. Rotating SESSION_SECRET invalidates
// all sessions.
package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const CookieName = "portfolio_session"

type Manager struct {
	secret []byte
	ttl    time.Duration
	secure bool
}

func New(secret string, ttl time.Duration, secure bool) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
		secure: secure,
	}
}

func (m *Manager) Enabled() bool { return len(m.secret) >= 16 }

// Sign returns the cookie value and expiry for the given email.
func (m *Manager) Sign(email string) (string, time.Time) {
	exp := time.Now().Add(m.ttl)
	expStr := strconv.FormatInt(exp.Unix(), 10)
	payload := email + "|" + expStr
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "|" + sig, exp
}

// Verify checks the signature + expiry and returns the email.
func (m *Manager) Verify(cookieValue string) (string, error) {
	parts := strings.SplitN(cookieValue, "|", 3)
	if len(parts) != 3 {
		return "", errors.New("malformed cookie")
	}
	email, expStr, sigStr := parts[0], parts[1], parts[2]
	if email == "" {
		return "", errors.New("empty email")
	}
	expUnix, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return "", errors.New("bad exp")
	}
	if time.Now().Unix() > expUnix {
		return "", errors.New("expired")
	}
	gotSig, err := base64.RawURLEncoding.DecodeString(sigStr)
	if err != nil {
		return "", errors.New("bad signature encoding")
	}
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(email + "|" + expStr))
	wantSig := mac.Sum(nil)
	if subtle.ConstantTimeCompare(gotSig, wantSig) != 1 {
		return "", errors.New("signature mismatch")
	}
	return email, nil
}

// SetCookie writes the session cookie on the response.
func (m *Manager) SetCookie(w http.ResponseWriter, value string, exp time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    value,
		Path:     "/admin",
		Expires:  exp,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearCookie writes an expired cookie to force logout.
func (m *Manager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// FromRequest reads and verifies the session cookie from the incoming request.
// Returns email or an error.
func (m *Manager) FromRequest(r *http.Request) (string, error) {
	c, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return m.Verify(c.Value)
}
