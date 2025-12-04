-- +goose Up
-- Waitlist table for pre-launch email capture

CREATE TABLE IF NOT EXISTS waitlist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    source TEXT NOT NULL DEFAULT 'website',
    ip_address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for duplicate checking and analytics
CREATE INDEX IF NOT EXISTS idx_waitlist_email ON waitlist(email);
CREATE INDEX IF NOT EXISTS idx_waitlist_created_at ON waitlist(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_waitlist_created_at;
DROP INDEX IF EXISTS idx_waitlist_email;
DROP TABLE IF EXISTS waitlist;
