-- Portfolio ticket system — initial schema
-- All timestamps in UTC, all IDs UUIDv7 (generated app-side for index locality).

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ───────────── enums ─────────────
DO $$ BEGIN
  CREATE TYPE ticket_status AS ENUM ('new','open','replied','closed','spam');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE ticket_category AS ENUM ('general','job','security','collaboration','other');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE msg_direction AS ENUM ('inbound','outbound','system');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- ───────────── tickets ─────────────
CREATE TABLE IF NOT EXISTS tickets (
  id            uuid PRIMARY KEY,
  public_code   text UNIQUE NOT NULL,
  first_name    text NOT NULL,
  last_name     text NOT NULL,
  email         citext NOT NULL,
  category      ticket_category NOT NULL DEFAULT 'general',
  subject       text NOT NULL,
  status        ticket_status NOT NULL DEFAULT 'new',
  priority      smallint NOT NULL DEFAULT 0,
  source_ip     inet,
  user_agent    text,
  cf_country    text,
  turnstile_ok  boolean NOT NULL DEFAULT false,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  closed_at     timestamptz,
  CONSTRAINT tickets_subject_len CHECK (char_length(subject) BETWEEN 3 AND 200),
  CONSTRAINT tickets_email_format CHECK (email ~ '^[^@\s]+@[^@\s]+\.[^@\s]+$')
);

CREATE INDEX IF NOT EXISTS tickets_status_created_idx
  ON tickets (status, created_at DESC);
CREATE INDEX IF NOT EXISTS tickets_email_idx
  ON tickets (email);
CREATE INDEX IF NOT EXISTS tickets_category_idx
  ON tickets (category);

-- ───────────── messages ─────────────
CREATE TABLE IF NOT EXISTS messages (
  id            uuid PRIMARY KEY,
  ticket_id     uuid NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
  direction     msg_direction NOT NULL,
  author_email  citext NOT NULL,
  author_name   text,
  body_text     text NOT NULL,
  body_html     text,
  raw_mime      text,
  provider_id   text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT messages_body_len CHECK (char_length(body_text) BETWEEN 1 AND 50000)
);

CREATE INDEX IF NOT EXISTS messages_ticket_created_idx
  ON messages (ticket_id, created_at ASC);

-- ───────────── admin_users ─────────────
CREATE TABLE IF NOT EXISTS admin_users (
  email         citext PRIMARY KEY,
  display_name  text,
  role          text NOT NULL DEFAULT 'admin',
  created_at    timestamptz NOT NULL DEFAULT now()
);

INSERT INTO admin_users (email, display_name, role)
VALUES ('carlo4340@outlook.it', 'Carlo Maria Cardi', 'owner')
ON CONFLICT (email) DO NOTHING;

-- ───────────── audit_log ─────────────
CREATE TABLE IF NOT EXISTS audit_log (
  id            bigserial PRIMARY KEY,
  actor_email   citext,
  action        text NOT NULL,
  target_type   text,
  target_id     text,
  metadata      jsonb NOT NULL DEFAULT '{}'::jsonb,
  ip            inet,
  created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS audit_log_created_idx
  ON audit_log (created_at DESC);
CREATE INDEX IF NOT EXISTS audit_log_target_idx
  ON audit_log (target_type, target_id);

-- ───────────── updated_at trigger ─────────────
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tickets_set_updated_at ON tickets;
CREATE TRIGGER tickets_set_updated_at
  BEFORE UPDATE ON tickets
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
