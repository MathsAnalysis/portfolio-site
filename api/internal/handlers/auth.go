package handlers

import (
	"crypto/subtle"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/MathsAnalysis/portfolio-api/internal/middleware"
	"github.com/MathsAnalysis/portfolio-api/internal/session"
	"golang.org/x/crypto/bcrypt"
)

// Auth serves the admin login/logout flow.
type Auth struct {
	Templates TemplateSet
	Sessions  *session.Manager
	Log       *slog.Logger

	AdminUser string // must match posted email
	AdminHash string // bcrypt hash of the password
	Realm     string // shown in the login page heading
}

type loginData struct {
	Error string
	Next  string
	Email string
	Realm string
	Now   time.Time
}

// GET /admin/login
func (h *Auth) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, bounce to /admin
	if h.Sessions != nil && h.Sessions.Enabled() {
		if _, err := h.Sessions.FromRequest(r); err == nil {
			http.Redirect(w, r, safeNext(r.URL.Query().Get("next")), http.StatusFound)
			return
		}
	}

	data := loginData{
		Next:  safeNext(r.URL.Query().Get("next")),
		Realm: defaultRealm(h.Realm),
		Now:   time.Now(),
	}
	h.renderLogin(w, data, http.StatusOK)
}

// POST /admin/login
func (h *Auth) LoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	pass := r.FormValue("password")
	next := safeNext(r.FormValue("next"))

	data := loginData{
		Next:  next,
		Email: email,
		Realm: defaultRealm(h.Realm),
		Now:   time.Now(),
	}

	if h.Sessions == nil || !h.Sessions.Enabled() {
		h.Log.Error("login: session manager disabled")
		data.Error = "Login is currently disabled. Contact the administrator."
		h.renderLogin(w, data, http.StatusInternalServerError)
		return
	}
	if h.AdminUser == "" || !strings.HasPrefix(h.AdminHash, "$2") {
		h.Log.Error("login: admin credentials not configured")
		data.Error = "Admin account not configured."
		h.renderLogin(w, data, http.StatusInternalServerError)
		return
	}

	// Constant-time username compare, then bcrypt.
	userOK := subtle.ConstantTimeCompare([]byte(email), []byte(h.AdminUser)) == 1
	passErr := bcrypt.CompareHashAndPassword([]byte(h.AdminHash), []byte(pass))
	if !userOK || passErr != nil {
		h.Log.Warn("login: bad credentials",
			"email", email,
			"ip", middleware.ClientIP(r),
			"ua", r.UserAgent(),
		)
		// Always generic error — don't leak which of user/pass was wrong.
		data.Error = "Invalid email or password."
		h.renderLogin(w, data, http.StatusUnauthorized)
		return
	}

	// Success — set signed cookie
	value, exp := h.Sessions.Sign(email)
	h.Sessions.SetCookie(w, value, exp)
	h.Log.Info("login: success",
		"email", email,
		"ip", middleware.ClientIP(r),
		"expires", exp.Format(time.RFC3339),
	)
	http.Redirect(w, r, next, http.StatusFound)
}

// POST /admin/logout  (also accepts GET for convenience)
func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	if h.Sessions != nil {
		h.Sessions.ClearCookie(w)
	}
	h.Log.Info("logout", "ip", middleware.ClientIP(r))
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

// ---------------- helpers ----------------

func (h *Auth) renderLogin(w http.ResponseWriter, data loginData, code int) {
	t, ok := h.Templates["login.html"]
	if !ok {
		http.Error(w, "login template missing", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Frame-Options", "DENY")
	w.WriteHeader(code)
	if err := t.ExecuteTemplate(w, "login.html", data); err != nil {
		h.Log.Error("login template exec", "err", err)
	}
}

// safeNext constrains the post-login redirect to paths inside /admin so we
// can't be tricked into redirecting off-site.
func safeNext(raw string) string {
	if raw == "" ||
		!strings.HasPrefix(raw, "/admin") ||
		strings.HasPrefix(raw, "/admin/login") ||
		strings.HasPrefix(raw, "//") ||
		strings.Contains(raw, "\r") || strings.Contains(raw, "\n") {
		return "/admin"
	}
	return raw
}

func defaultRealm(r string) string {
	if r == "" {
		return "mathsanalysis admin"
	}
	return r
}

// compile-time: ensure html/template is used (for code that moves here later)
var _ = template.FuncMap{}
