package middleware

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ctxKey int

const (
	ctxRequestID ctxKey = iota
	ctxAdminEmail
)

// Recover converts panics into 500s and logs the stack.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic",
						"err", rec,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Logger writes one structured log line per request.
func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(sw, r)
			log.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.status,
				"bytes", sw.bytes,
				"dur_ms", time.Since(start).Milliseconds(),
				"ip", ClientIP(r),
				"ua", r.UserAgent(),
			)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// ClientIP extracts the real client IP, preferring CF-Connecting-IP
// (set by Cloudflare Tunnel at the edge), falling back to X-Forwarded-For,
// then the remote addr.
func ClientIP(r *http.Request) string {
	if v := r.Header.Get("CF-Connecting-IP"); v != "" {
		return v
	}
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		// first entry
		for i := 0; i < len(v); i++ {
			if v[i] == ',' {
				return v[:i]
			}
		}
		return v
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// TokenVerifier is satisfied by *auth.CFAccess — injected so middleware
// doesn't depend on the auth package directly.
type TokenVerifier interface {
	Verify(ctx context.Context, token string) (email string, err error)
}

// BasicAuthConfig holds credentials for HTTP Basic Auth. The hash is bcrypt
// (cost ≥ 10) of the admin password.
type BasicAuthConfig struct {
	Username string // e.g. carlo4340@outlook.it
	Hash     string // $2[aby]$... bcrypt
	Realm    string // "mathsanalysis admin"
}

func (b BasicAuthConfig) Enabled() bool {
	return b.Username != "" && strings.HasPrefix(b.Hash, "$2")
}

// AdminAuth checks, in order:
//  1. HTTP Basic Auth (if configured) — primary production path
//  2. mockEmail — dev-only bypass
//  3. CF Access signed JWT (if configured)
//  4. Otherwise reject with 401.
//
// Returns 401 with WWW-Authenticate header when basic auth is configured,
// so browsers show the native login prompt on first visit.
func AdminAuth(
	basic BasicAuthConfig,
	mockEmail string,
	verifier TokenVerifier,
	log *slog.Logger,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var email string

			// 1. HTTP Basic Auth (primary)
			if basic.Enabled() {
				user, pass, ok := r.BasicAuth()
				if ok && subtle.ConstantTimeCompare([]byte(user), []byte(basic.Username)) == 1 {
					// bcrypt is slow-by-design (~100ms) — we do it only on a
					// well-formed credential. Rate-limit above protects us from
					// mass brute-force at the public edge.
					if err := bcrypt.CompareHashAndPassword([]byte(basic.Hash), []byte(pass)); err == nil {
						email = basic.Username
					} else {
						log.Warn("admin basic auth: wrong password",
							"user", user,
							"ip", ClientIP(r),
						)
					}
				}
				if email == "" {
					realm := basic.Realm
					if realm == "" {
						realm = "mathsanalysis admin"
					}
					w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`", charset="UTF-8"`)
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// 2. dev bypass (only if basic auth wasn't configured)
			if email == "" && mockEmail != "" {
				email = mockEmail
			}

			// 3. CF Access JWT
			if email == "" && verifier != nil {
				token := r.Header.Get("Cf-Access-Jwt-Assertion")
				if token == "" {
					if c, err := r.Cookie("CF_Authorization"); err == nil {
						token = c.Value
					}
				}
				if token != "" {
					e, err := verifier.Verify(r.Context(), token)
					if err != nil {
						log.Warn("admin auth: jwt invalid",
							"err", err,
							"path", r.URL.Path,
							"ip", ClientIP(r),
						)
					} else {
						email = e
					}
				}
			}

			if email == "" {
				log.Warn("admin auth: unauthenticated",
					"path", r.URL.Path,
					"ip", ClientIP(r),
					"basic", basic.Enabled(),
					"mock", mockEmail != "",
					"verifier", verifier != nil,
				)
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxAdminEmail, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminEmail(ctx context.Context) string {
	v, _ := ctx.Value(ctxAdminEmail).(string)
	return v
}
