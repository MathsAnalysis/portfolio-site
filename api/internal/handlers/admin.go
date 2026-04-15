package handlers

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MathsAnalysis/portfolio-api/internal/db"
	"github.com/MathsAnalysis/portfolio-api/internal/email"
	"github.com/MathsAnalysis/portfolio-api/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TemplateSet is a map of template name → parsed template. Each entry is a
// fully-scoped template (its own set of defines) so pages don't collide.
type TemplateSet map[string]*template.Template

type Admin struct {
	Store     *db.TicketStore
	Templates TemplateSet
	Log       *slog.Logger
	Sender    email.Sender // outbound reply delivery (Resend)
}

// ---------- helpers ----------

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Admin) render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	data["AdminEmail"] = middleware.AdminEmail(r.Context())
	data["Now"] = time.Now()

	t, ok := h.Templates[name]
	if !ok {
		h.Log.Error("template not found", "name", name)
		http.Error(w, "template "+name+" not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	// For page templates (with layout), execute by the template's own name.
	// For fragments, executing the named template emits just the defined block.
	if err := t.ExecuteTemplate(w, name, data); err != nil {
		h.Log.Error("template", "name", name, "err", err)
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

// ---------- routes ----------

// GET /admin — redirect to /admin/tickets
func (h *Admin) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/tickets", http.StatusFound)
}

// GET /admin/tickets
func (h *Admin) TicketsList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := db.ListFilter{
		Status:   q.Get("status"),
		Category: q.Get("category"),
		Search:   q.Get("q"),
		Limit:    50,
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filter.Offset = n
		}
	}

	tickets, total, err := h.Store.List(r.Context(), filter)
	if err != nil {
		h.Log.Error("list tickets", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	stats, err := h.Store.Stats(r.Context())
	if err != nil {
		h.Log.Error("stats", "err", err)
	}

	data := map[string]any{
		"Tickets":  tickets,
		"Total":    total,
		"Filter":   filter,
		"Stats":    stats,
		"Category": filter.Category,
		"Status":   filter.Status,
	}

	if isHTMX(r) {
		h.render(w, r, "fragment_tickets_list.html", data)
		return
	}
	h.render(w, r, "tickets_list.html", data)
}

// GET /admin/tickets/{id}
func (h *Admin) TicketDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ticket, msgs, err := h.Store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.Log.Error("get ticket", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	stats, _ := h.Store.Stats(r.Context())

	data := map[string]any{
		"Ticket":   ticket,
		"Messages": msgs,
		"Stats":    stats,
	}
	if isHTMX(r) {
		h.render(w, r, "fragment_ticket_detail.html", data)
		return
	}
	h.render(w, r, "ticket_detail.html", data)
}

// POST /admin/tickets/{id}/reply
func (h *Admin) TicketReply(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	body := strings.TrimSpace(r.FormValue("body"))
	if len(body) < 1 || len(body) > 10000 {
		http.Error(w, "body length 1-10000", http.StatusBadRequest)
		return
	}

	adminEmail := middleware.AdminEmail(r.Context())
	if _, err := h.Store.AddReply(r.Context(), id, adminEmail, adminEmail, body); err != nil {
		h.Log.Error("add reply", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// Load the ticket for outbound email dispatch
	ticket, _, err := h.Store.Get(r.Context(), id)
	if err != nil {
		h.Log.Error("reload ticket", "err", err)
	} else if h.Sender != nil {
		// Fire-and-forget outbound via Resend — detached context with timeout.
		go func(t *db.Ticket, bodyText, fromName string) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			providerID, err := h.Sender.SendReply(ctx, email.Ticket{
				ID:         t.ID.String(),
				PublicCode: t.PublicCode,
				FirstName:  t.FirstName,
				LastName:   t.LastName,
				Email:      t.Email,
				Category:   t.Category,
				Subject:    t.Subject,
				CreatedAt:  t.CreatedAt,
			}, bodyText, fromName)
			if err != nil {
				h.Log.Error("send reply email",
					"err", err,
					"ticket", t.PublicCode,
				)
				return
			}
			h.Log.Info("reply email sent",
				"ticket", t.PublicCode,
				"provider_id", providerID,
			)
		}(ticket, body, adminEmail)
	}

	h.Log.Info("admin reply",
		"ticket", id,
		"admin", adminEmail,
	)

	// Return the updated detail fragment (htmx swap target)
	r.Method = http.MethodGet
	h.TicketDetail(w, r)
}

// POST /admin/tickets/{id}/status
func (h *Admin) TicketStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	status := r.FormValue("status")
	switch status {
	case "new", "open", "replied", "closed", "spam", "archived":
	default:
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}
	if err := h.Store.UpdateStatus(r.Context(), id, status); err != nil {
		h.Log.Error("update status", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	h.Log.Info("status changed",
		"ticket", id,
		"to", status,
		"admin", middleware.AdminEmail(r.Context()),
	)

	r.Method = http.MethodGet
	h.TicketDetail(w, r)
}

// POST /admin/tickets/{id}/delete
// Permanent deletion. Writes an audit_log entry with actor + requester
// metadata, then cascades to messages. Responds with an htmx redirect to
// /admin/tickets so the list refreshes after deletion.
func (h *Admin) TicketDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	actor := middleware.AdminEmail(r.Context())
	ip := middleware.ClientIP(r)

	if err := h.Store.Delete(r.Context(), id, actor, ip); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.Log.Error("delete ticket", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	h.Log.Warn("ticket deleted",
		"ticket", id,
		"admin", actor,
		"ip", ip,
	)

	// For htmx: send redirect header; for plain form: 303 to list
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/tickets")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/admin/tickets", http.StatusSeeOther)
}
