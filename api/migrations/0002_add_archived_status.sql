-- Add 'archived' to ticket_status enum.
-- Safe re-run: the IF NOT EXISTS clause is available since PG 9.6 for
-- ALTER TYPE ... ADD VALUE.
DO $$ BEGIN
  ALTER TYPE ticket_status ADD VALUE IF NOT EXISTS 'archived';
EXCEPTION WHEN duplicate_object THEN
  -- already present
  NULL;
END $$;
