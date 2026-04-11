// Package email provides two building blocks:
//
//   Notifier — sends internal notifications to Carlo (new ticket arrived).
//              Backed by SMTP (e.g. the axiom-mailserver container).
//
//   Sender   — sends outbound replies to visitors.
//              Backed by Resend (for deliverability, DKIM/SPF on mathsanalysis.com).
//
// Both interfaces have no-op default implementations. If their corresponding
// env vars are unset (dev / Phase 1), the app simply logs and continues — email
// operations never crash the request path.
package email

import (
	"context"
	"time"
)

// Ticket is the minimal subset of a ticket the email layer needs. Defined here
// to avoid a cyclic import with the db package.
type Ticket struct {
	ID         string
	PublicCode string
	FirstName  string
	LastName   string
	Email      string
	Category   string
	Subject    string
	Message    string
	CreatedAt  time.Time
}

// Notifier is called when a new ticket is created. Implementations should be
// non-blocking from the request's perspective (errors logged, not returned to
// the visitor).
type Notifier interface {
	NotifyNewTicket(ctx context.Context, t Ticket) error
}

// Sender delivers a reply to the visitor. Implementations set the From address
// to a per-ticket subaddress (tickets+<public_code>@<domain>) so inbound
// replies from the visitor land back on the same ticket thread.
type Sender interface {
	SendReply(ctx context.Context, t Ticket, bodyText, fromName string) (providerID string, err error)
}

// ---------- no-op ----------

type noop struct{}

func (noop) NotifyNewTicket(_ context.Context, _ Ticket) error { return nil }
func (noop) SendReply(_ context.Context, _ Ticket, _, _ string) (string, error) {
	return "noop", nil
}

// NoopNotifier is a drop-in Notifier that silently succeeds.
func NoopNotifier() Notifier { return noop{} }

// NoopSender is a drop-in Sender that silently succeeds.
func NoopSender() Sender { return noop{} }
