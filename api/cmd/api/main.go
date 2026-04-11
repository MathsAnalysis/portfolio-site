package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/MathsAnalysis/portfolio-api/internal/auth"
	"github.com/MathsAnalysis/portfolio-api/internal/config"
	"github.com/MathsAnalysis/portfolio-api/internal/db"
	"github.com/MathsAnalysis/portfolio-api/internal/email"
	"github.com/MathsAnalysis/portfolio-api/internal/handlers"
	mw "github.com/MathsAnalysis/portfolio-api/internal/middleware"
	"github.com/MathsAnalysis/portfolio-api/internal/turnstile"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	cfg, err := config.FromEnv()
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(1)
	}
	log.Info("booting portfolio-api",
		"env", cfg.Env,
		"addr", cfg.HTTPAddr,
		"mock_admin", cfg.AdminMockEmail != "",
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool, cfg.MigrationsDir, log); err != nil {
		log.Error("migrations", "err", err)
		os.Exit(1)
	}

	store := db.NewTicketStore(pool)

	// Parse admin templates
	tmplDir := os.Getenv("TEMPLATES_DIR")
	if tmplDir == "" {
		tmplDir = "/app/web/templates"
	}
	tmplSet, err := loadTemplateSet(tmplDir)
	if err != nil {
		log.Error("templates", "err", err)
		os.Exit(1)
	}

	// ─── email: notifier (SMTP) + sender (Resend) ───
	notifier, smtpEnabled := email.NewSMTPNotifier(email.SMTPConfig{
		Host:       cfg.SMTPHost,
		Port:       cfg.SMTPPort,
		Username:   cfg.SMTPUser,
		Password:   cfg.SMTPPassword,
		From:       cfg.SMTPFrom,
		To:         cfg.SMTPTo,
		UseTLS:     cfg.SMTPUseTLS,
		SkipVerify: cfg.SMTPSkipVerify,
	}, log)

	sender, resendEnabled := email.NewResendSender(email.ResendConfig{
		APIKey:      cfg.ResendAPIKey,
		FromAddress: cfg.ResendFrom,
		Domain:      cfg.ResendDomain,
	}, log)

	// ─── turnstile verifier (noop when secret empty) ───
	tv := turnstile.New(cfg.TurnstileSecret)

	// ─── CF Access JWT verifier (nil when unconfigured) ───
	var jwtVerifier mw.TokenVerifier
	if cfg.CFAccessTeamDomain != "" && cfg.CFAccessAudience != "" {
		jwtVerifier = auth.New(auth.Config{
			TeamDomain: cfg.CFAccessTeamDomain,
			Audience:   cfg.CFAccessAudience,
		})
	}

	basicAuth := mw.BasicAuthConfig{
		Username: cfg.AdminBasicUser,
		Hash:     cfg.AdminBasicHash,
		Realm:    cfg.AdminRealm,
	}

	log.Info("subsystems",
		"smtp_notifier", smtpEnabled,
		"resend_sender", resendEnabled,
		"turnstile", tv.Enabled(),
		"cf_access_jwt", jwtVerifier != nil,
		"admin_basic_auth", basicAuth.Enabled(),
		"admin_mock", cfg.AdminMockEmail != "",
		"inbound_webhook", cfg.InboundWebhookSecret != "",
	)

	pub := &handlers.Public{
		Store:     store,
		Log:       log,
		Turnstile: tv,
		Notifier:  notifier,
	}
	adm := &handlers.Admin{
		Store:     store,
		Templates: tmplSet,
		Log:       log,
		Sender:    sender,
	}
	wh := &handlers.Webhook{
		Store:         store,
		Log:           log,
		InboundSecret: cfg.InboundWebhookSecret,
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mw.Recover(log))
	r.Use(mw.Logger(log))

	// -------- public API --------
	r.Get("/api/health", pub.Health)
	r.Group(func(g chi.Router) {
		g.Use(httprate.Limit(
			10,
			time.Minute,
			httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
				return mw.ClientIP(r), nil
			}),
		))
		g.Post("/api/tickets", pub.CreateTicket)
	})

	// -------- inbound email webhook (HMAC-signed by CF Worker) --------
	r.Post("/api/webhook/email", wh.InboundEmail)

	// -------- admin (auth required) --------
	r.Group(func(g chi.Router) {
		g.Use(mw.AdminAuth(basicAuth, cfg.AdminMockEmail, jwtVerifier, log))
		g.Get("/admin", adm.Index)
		g.Get("/admin/tickets", adm.TicketsList)
		g.Get("/admin/tickets/{id}", adm.TicketDetail)
		g.Post("/admin/tickets/{id}/reply", adm.TicketReply)
		g.Post("/admin/tickets/{id}/status", adm.TicketStatus)
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info("shutting down")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = srv.Shutdown(shutdownCtx)
		cancel()
	}()

	log.Info("listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("http", "err", err)
		os.Exit(1)
	}
}

// loadTemplateSet returns a map of page-name → *template.Template.
// Each page is parsed in its OWN namespace (layout + fragments + page) so that
// {{define "title"}} / {{define "content"}} blocks don't collide across pages.
// Fragment templates (those not listed as pages) are also exposed as standalone
// entries keyed by their file name, so handlers can render a fragment directly.
func loadTemplateSet(tmplDir string) (handlers.TemplateSet, error) {
	funcs := template.FuncMap{
		"humanTime": func(ts time.Time) string {
			return ts.UTC().Format("2006-01-02 15:04 UTC")
		},
		"statusColor": func(s string) string {
			switch s {
			case "new":
				return "bg-emerald-500/10 text-emerald-400 border-emerald-500/30"
			case "open":
				return "bg-cyan-500/10 text-cyan-400 border-cyan-500/30"
			case "replied":
				return "bg-violet-500/10 text-violet-400 border-violet-500/30"
			case "closed":
				return "bg-zinc-500/10 text-zinc-400 border-zinc-500/30"
			case "spam":
				return "bg-rose-500/10 text-rose-400 border-rose-500/30"
			}
			return "bg-zinc-500/10 text-zinc-400 border-zinc-500/30"
		},
		"categoryLabel": func(c string) string {
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
		},
	}

	// Collect all files
	var all []string
	err := filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		all = append(all, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	layoutPath := filepath.Join(tmplDir, "layout.html")

	var fragments []string
	var pages []string
	for _, p := range all {
		base := filepath.Base(p)
		switch {
		case base == "layout.html":
			// handled separately
		case strings.HasPrefix(base, "fragment_"):
			fragments = append(fragments, p)
		default:
			pages = append(pages, p)
		}
	}

	set := handlers.TemplateSet{}

	// Each page gets its own set: layout + ALL fragments + the page itself.
	for _, pagePath := range pages {
		name := filepath.Base(pagePath)
		files := append([]string{layoutPath}, fragments...)
		files = append(files, pagePath)
		t, err := template.New(name).Funcs(funcs).ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("parse page %q: %w", name, err)
		}
		set[name] = t
	}

	// Fragments rendered on their own (htmx requests): layout NOT included.
	// They may depend on other fragments, so parse all of them together.
	if len(fragments) > 0 {
		for _, fragPath := range fragments {
			name := filepath.Base(fragPath)
			t, err := template.New(name).Funcs(funcs).ParseFiles(fragments...)
			if err != nil {
				return nil, fmt.Errorf("parse fragment %q: %w", name, err)
			}
			set[name] = t
		}
	}

	return set, nil
}
