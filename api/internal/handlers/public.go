package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/MathsAnalysis/portfolio-api/internal/db"
	"github.com/MathsAnalysis/portfolio-api/internal/email"
	"github.com/MathsAnalysis/portfolio-api/internal/middleware"
	"github.com/MathsAnalysis/portfolio-api/internal/turnstile"
)

type Public struct {
	Store     *db.TicketStore
	Log       *slog.Logger
	Turnstile *turnstile.Verifier
	Notifier  email.Notifier
}

type createTicketReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Category  string `json:"category"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	// Turnstile token — ignored in phase 1, validated in phase 3
	Turnstile string `json:"cf_turnstile"` //nolint:tagliatelle
}

type createTicketResp struct {
	OK         bool   `json:"ok"`
	PublicCode string `json:"public_code,omitempty"`
	Error      string `json:"error,omitempty"`
}

var validCategories = map[string]struct{}{
	"general":       {},
	"job":           {},
	"security":      {},
	"collaboration": {},
	"other":         {},
}

func (h *Public) CreateTicket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32*1024))
	dec.DisallowUnknownFields()

	var req createTicketReq
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, createTicketResp{Error: "invalid json: " + err.Error()})
		return
	}

	if err := validateTicket(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, createTicketResp{Error: err.Error()})
		return
	}

	// Turnstile (no-op when disabled)
	if err := h.Turnstile.Verify(r.Context(), req.Turnstile, middleware.ClientIP(r)); err != nil {
		h.Log.Warn("turnstile reject",
			"err", err,
			"ip", middleware.ClientIP(r),
		)
		writeJSON(w, http.StatusForbidden, createTicketResp{Error: "challenge failed"})
		return
	}

	ticket, err := h.Store.Create(r.Context(), db.CreateTicketInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     strings.ToLower(req.Email),
		Category:  req.Category,
		Subject:   req.Subject,
		Message:   req.Message,
		SourceIP:  middleware.ClientIP(r),
		UserAgent: r.UserAgent(),
		CFCountry: r.Header.Get("CF-IPCountry"),
	})
	if err != nil {
		h.Log.Error("create ticket", "err", err)
		writeJSON(w, http.StatusInternalServerError, createTicketResp{Error: "could not save ticket"})
		return
	}

	h.Log.Info("ticket created",
		"id", ticket.ID,
		"code", ticket.PublicCode,
		"category", ticket.Category,
		"email", ticket.Email,
	)

	// Fire-and-forget notification to Carlo.
	// We detach from the request context (which is canceled when the handler
	// returns) and give the async job a bounded 20s deadline.
	go func() {
		if h.Notifier == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := h.Notifier.NotifyNewTicket(ctx, email.Ticket{
			ID:         ticket.ID.String(),
			PublicCode: ticket.PublicCode,
			FirstName:  ticket.FirstName,
			LastName:   ticket.LastName,
			Email:      ticket.Email,
			Category:   ticket.Category,
			Subject:    ticket.Subject,
			Message:    req.Message,
			CreatedAt:  ticket.CreatedAt,
		}); err != nil {
			h.Log.Error("notify new ticket", "err", err, "code", ticket.PublicCode)
		}
	}()

	writeJSON(w, http.StatusCreated, createTicketResp{
		OK:         true,
		PublicCode: ticket.PublicCode,
	})
}

func validateTicket(r *createTicketReq) error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(r.Email)
	r.Category = strings.TrimSpace(r.Category)
	r.Subject = strings.TrimSpace(r.Subject)
	r.Message = strings.TrimSpace(r.Message)

	if len(r.FirstName) < 1 || len(r.FirstName) > 80 {
		return errors.New("first_name length must be 1-80")
	}
	if len(r.LastName) < 1 || len(r.LastName) > 80 {
		return errors.New("last_name length must be 1-80")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email")
	}
	if len(r.Email) > 254 {
		return errors.New("email too long")
	}
	if _, ok := validCategories[r.Category]; !ok {
		return errors.New("invalid category")
	}
	if len(r.Subject) < 3 || len(r.Subject) > 200 {
		return errors.New("subject length must be 3-200")
	}
	if len(r.Message) < 10 || len(r.Message) > 5000 {
		return errors.New("message length must be 10-5000")
	}
	return nil
}

func (h *Public) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
