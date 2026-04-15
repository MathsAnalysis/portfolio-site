package db

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticket struct {
	ID          uuid.UUID
	PublicCode  string
	FirstName   string
	LastName    string
	Email       string
	Category    string
	Subject     string
	Status      string
	Priority    int16
	SourceIP    *string
	UserAgent   *string
	CFCountry   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    *time.Time
	MessageLast *time.Time // joined: last message timestamp, optional
}

type Message struct {
	ID          uuid.UUID
	TicketID    uuid.UUID
	Direction   string
	AuthorEmail string
	AuthorName  *string
	BodyText    string
	BodyHTML    *string
	CreatedAt   time.Time
}

type TicketStore struct {
	pool *pgxpool.Pool
}

func NewTicketStore(pool *pgxpool.Pool) *TicketStore {
	return &TicketStore{pool: pool}
}

type CreateTicketInput struct {
	FirstName string
	LastName  string
	Email     string
	Category  string
	Subject   string
	Message   string
	SourceIP  string
	UserAgent string
	CFCountry string
}

func (s *TicketStore) Create(ctx context.Context, in CreateTicketInput) (*Ticket, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	id := uuid.New()
	msgID := uuid.New()
	code, err := newPublicCode()
	if err != nil {
		return nil, fmt.Errorf("gen code: %w", err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO tickets
		  (id, public_code, first_name, last_name, email, category, subject, source_ip, user_agent, cf_country)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, NULLIF($8,'')::inet, NULLIF($9,''), NULLIF($10,''))
		RETURNING id, public_code, first_name, last_name, email, category, subject,
		          status, priority, created_at, updated_at
	`, id, code, in.FirstName, in.LastName, in.Email, in.Category, in.Subject,
		in.SourceIP, in.UserAgent, in.CFCountry)

	var t Ticket
	if err := row.Scan(&t.ID, &t.PublicCode, &t.FirstName, &t.LastName, &t.Email,
		&t.Category, &t.Subject, &t.Status, &t.Priority, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, fmt.Errorf("insert ticket: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO messages
		  (id, ticket_id, direction, author_email, author_name, body_text)
		VALUES
		  ($1, $2, 'inbound', $3, $4, $5)
	`, msgID, id, in.Email, in.FirstName+" "+in.LastName, in.Message)
	if err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &t, nil
}

// ListFilter lets admins filter the dashboard view.
type ListFilter struct {
	Status   string // "" = all
	Category string
	Search   string
	Limit    int
	Offset   int
}

func (s *TicketStore) List(ctx context.Context, f ListFilter) ([]Ticket, int, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}

	// Build dynamic WHERE
	args := []any{}
	where := "WHERE 1=1"
	add := func(cond string, v any) {
		args = append(args, v)
		where += fmt.Sprintf(" AND %s $%d", cond, len(args))
	}
	switch f.Status {
	case "":
		// Default view: show all active tickets, HIDE archived
		where += " AND status <> 'archived'"
	case "all":
		// show everything including archived — no filter
	default:
		add("status =", f.Status)
	}
	if f.Category != "" {
		add("category =", f.Category)
	}
	if f.Search != "" {
		args = append(args, "%"+f.Search+"%")
		where += fmt.Sprintf(" AND (subject ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
			len(args), len(args), len(args), len(args))
	}

	// total
	var total int
	if err := s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM tickets "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, f.Offset)
	query := fmt.Sprintf(`
		SELECT id, public_code, first_name, last_name, email, category, subject,
		       status, priority, created_at, updated_at, closed_at
		FROM tickets
		%s
		ORDER BY
		  CASE status WHEN 'new' THEN 0 WHEN 'open' THEN 1 WHEN 'replied' THEN 2
		              WHEN 'closed' THEN 3 WHEN 'spam' THEN 4 END,
		  created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)-1, len(args))

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []Ticket
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.PublicCode, &t.FirstName, &t.LastName, &t.Email,
			&t.Category, &t.Subject, &t.Status, &t.Priority,
			&t.CreatedAt, &t.UpdatedAt, &t.ClosedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, t)
	}
	return out, total, rows.Err()
}

func (s *TicketStore) Get(ctx context.Context, id uuid.UUID) (*Ticket, []Message, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, public_code, first_name, last_name, email, category, subject,
		       status, priority, created_at, updated_at, closed_at
		FROM tickets WHERE id = $1
	`, id)

	var t Ticket
	if err := row.Scan(&t.ID, &t.PublicCode, &t.FirstName, &t.LastName, &t.Email,
		&t.Category, &t.Subject, &t.Status, &t.Priority,
		&t.CreatedAt, &t.UpdatedAt, &t.ClosedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, ticket_id, direction, author_email, author_name, body_text, body_html, created_at
		FROM messages WHERE ticket_id = $1 ORDER BY created_at ASC
	`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.TicketID, &m.Direction, &m.AuthorEmail,
			&m.AuthorName, &m.BodyText, &m.BodyHTML, &m.CreatedAt); err != nil {
			return nil, nil, err
		}
		msgs = append(msgs, m)
	}
	return &t, msgs, rows.Err()
}

func (s *TicketStore) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	var closedAt any
	if status == "closed" || status == "archived" {
		closedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE tickets SET status = $1, closed_at = $2 WHERE id = $3`,
		status, closedAt, id,
	)
	return err
}

// Delete permanently removes a ticket and its messages (FK cascades).
// Writes an audit_log entry before deleting for forensic trail.
func (s *TicketStore) Delete(ctx context.Context, id uuid.UUID, actorEmail, ip string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var publicCode, subject, requesterEmail string
	err = tx.QueryRow(ctx,
		`SELECT public_code, subject, email FROM tickets WHERE id = $1`, id,
	).Scan(&publicCode, &subject, &requesterEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (actor_email, action, target_type, target_id, metadata, ip)
		VALUES ($1, 'ticket.delete', 'ticket', $2,
		        jsonb_build_object('public_code', $3::text, 'subject', $4::text, 'requester_email', $5::text),
		        NULLIF($6,'')::inet)
	`, actorEmail, id.String(), publicCode, subject, requesterEmail, ip)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM tickets WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindByPublicCode returns the ticket whose public_code matches (case-insensitive).
func (s *TicketStore) FindByPublicCode(ctx context.Context, code string) (*Ticket, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, public_code, first_name, last_name, email, category, subject,
		       status, priority, created_at, updated_at, closed_at
		FROM tickets WHERE UPPER(public_code) = UPPER($1)
	`, code)
	var t Ticket
	if err := row.Scan(&t.ID, &t.PublicCode, &t.FirstName, &t.LastName, &t.Email,
		&t.Category, &t.Subject, &t.Status, &t.Priority,
		&t.CreatedAt, &t.UpdatedAt, &t.ClosedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// AddInboundMessage appends an inbound (visitor-sourced) message to an existing ticket
// and moves its status back to "open" if it was "replied" or "closed".
func (s *TicketStore) AddInboundMessage(
	ctx context.Context,
	ticketID uuid.UUID,
	authorEmail, authorName, bodyText, bodyHTML, rawMIME string,
) (*Message, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	msgID := uuid.New()
	row := tx.QueryRow(ctx, `
		INSERT INTO messages (id, ticket_id, direction, author_email, author_name, body_text, body_html, raw_mime)
		VALUES ($1, $2, 'inbound', $3, NULLIF($4,''), $5, NULLIF($6,''), NULLIF($7,''))
		RETURNING id, ticket_id, direction, author_email, author_name, body_text, created_at
	`, msgID, ticketID, authorEmail, authorName, bodyText, bodyHTML, rawMIME)

	var m Message
	if err := row.Scan(&m.ID, &m.TicketID, &m.Direction, &m.AuthorEmail,
		&m.AuthorName, &m.BodyText, &m.CreatedAt); err != nil {
		return nil, err
	}

	// Re-open ticket on visitor reply
	if _, err := tx.Exec(ctx,
		`UPDATE tickets SET status = 'open', closed_at = NULL WHERE id = $1 AND status IN ('replied','closed')`,
		ticketID,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *TicketStore) AddReply(ctx context.Context, ticketID uuid.UUID, authorEmail, authorName, body string) (*Message, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	msgID := uuid.New()
	row := tx.QueryRow(ctx, `
		INSERT INTO messages (id, ticket_id, direction, author_email, author_name, body_text)
		VALUES ($1, $2, 'outbound', $3, NULLIF($4,''), $5)
		RETURNING id, ticket_id, direction, author_email, author_name, body_text, created_at
	`, msgID, ticketID, authorEmail, authorName, body)

	var m Message
	if err := row.Scan(&m.ID, &m.TicketID, &m.Direction, &m.AuthorEmail,
		&m.AuthorName, &m.BodyText, &m.CreatedAt); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx,
		`UPDATE tickets SET status = 'replied' WHERE id = $1 AND status IN ('new','open')`,
		ticketID,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

// Stats returns counts per status — used by admin sidebar badges.
type Stats struct {
	New, Open, Replied, Closed, Spam, Archived, Total int
}

func (s *TicketStore) Stats(ctx context.Context) (Stats, error) {
	var st Stats
	err := s.pool.QueryRow(ctx, `
		SELECT
		  COUNT(*) FILTER (WHERE status = 'new'),
		  COUNT(*) FILTER (WHERE status = 'open'),
		  COUNT(*) FILTER (WHERE status = 'replied'),
		  COUNT(*) FILTER (WHERE status = 'closed'),
		  COUNT(*) FILTER (WHERE status = 'spam'),
		  COUNT(*) FILTER (WHERE status = 'archived'),
		  COUNT(*)
		FROM tickets
	`).Scan(&st.New, &st.Open, &st.Replied, &st.Closed, &st.Spam, &st.Archived, &st.Total)
	return st, err
}

// newPublicCode returns e.g. "TK-7F3A9B" — short, URL-safe, collision-resistant enough for our volume.
func newPublicCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // no 0/O/1/I
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, 6)
	for i, b := range buf {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return "TK-" + string(out), nil
}

// ErrNotFound is returned when a row doesn't exist.
var ErrNotFound = errNotFound{}

type errNotFound struct{}

func (errNotFound) Error() string { return "not found" }
