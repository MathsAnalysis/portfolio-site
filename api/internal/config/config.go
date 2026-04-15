package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env            string
	HTTPAddr       string
	DatabaseURL    string
	MigrationsDir  string
	TrustedProxies []string
	AdminMockEmail string // phase 1: bypass CF Access; set to "" in production
	PublicBaseURL  string

	// HTTP Basic Auth on /admin (fallback for curl / CI access)
	AdminBasicUser string
	AdminBasicHash string // bcrypt hash; must start with $2
	AdminRealm     string

	// Session cookie for the admin login page (HMAC secret)
	SessionSecret string
	SessionTTL    int // hours
	SessionSecure bool

	// CF Access JWT verification
	CFAccessTeamDomain string // e.g. "mathsanalysis.cloudflareaccess.com"
	CFAccessAudience   string // AUD tag from CF Access application

	// Cloudflare Turnstile
	TurnstileSecret string

	// Internal SMTP (axiom-mailserver) for new-ticket notifications
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
	SMTPTo         string
	SMTPUseTLS     bool
	SMTPSkipVerify bool

	// Resend for outbound visitor replies
	ResendAPIKey  string
	ResendFrom    string
	ResendDomain  string

	// HMAC secret for the inbound email webhook
	InboundWebhookSecret string
}

func FromEnv() (*Config, error) {
	c := &Config{
		Env:                  getEnv("APP_ENV", "development"),
		HTTPAddr:             getEnv("HTTP_ADDR", ":8182"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		MigrationsDir:        getEnv("MIGRATIONS_DIR", "/app/migrations"),
		AdminMockEmail:       os.Getenv("ADMIN_MOCK_EMAIL"),
		PublicBaseURL:        getEnv("PUBLIC_BASE_URL", "https://mathsanalysis.com"),
		AdminBasicUser:       os.Getenv("ADMIN_BASIC_USER"),
		AdminBasicHash:       os.Getenv("ADMIN_BASIC_HASH"),
		AdminRealm:           getEnv("ADMIN_REALM", "mathsanalysis admin"),
		SessionSecret:        os.Getenv("SESSION_SECRET"),
		SessionTTL:           getEnvInt("SESSION_TTL_HOURS", 12),
		SessionSecure:        getEnvBool("SESSION_SECURE", true),
		CFAccessTeamDomain:   os.Getenv("CF_ACCESS_TEAM_DOMAIN"),
		CFAccessAudience:     os.Getenv("CF_ACCESS_AUDIENCE"),
		TurnstileSecret:      os.Getenv("TURNSTILE_SECRET"),
		SMTPHost:             os.Getenv("SMTP_HOST"),
		SMTPPort:             getEnvInt("SMTP_PORT", 587),
		SMTPUser:             os.Getenv("SMTP_USER"),
		SMTPPassword:         os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:             os.Getenv("SMTP_FROM"),
		SMTPTo:               os.Getenv("SMTP_TO"),
		SMTPUseTLS:           getEnvBool("SMTP_USE_TLS", false),
		SMTPSkipVerify:       getEnvBool("SMTP_SKIP_VERIFY", false),
		ResendAPIKey:         os.Getenv("RESEND_API_KEY"),
		ResendFrom:           os.Getenv("RESEND_FROM"),
		ResendDomain:         os.Getenv("RESEND_DOMAIN"),
		InboundWebhookSecret: os.Getenv("INBOUND_WEBHOOK_SECRET"),
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return c, nil
}

func (c *Config) IsDev() bool { return c.Env == "development" }

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "TRUE", "yes", "on":
		return true
	case "0", "false", "FALSE", "no", "off":
		return false
	}
	return def
}
