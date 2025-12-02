-- +goose Up
-- Webhook event storage for durability and retry support
CREATE TABLE IF NOT EXISTS stripe_webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_event_id TEXT UNIQUE NOT NULL,
    event_type TEXT NOT NULL,
    event_data JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'dead_letter')),
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 5,
    next_retry_at TIMESTAMPTZ,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for efficient polling of pending/failed events
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_retry
    ON stripe_webhook_events (status, next_retry_at)
    WHERE status IN ('pending', 'failed');

-- Index for deduplication checks
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_stripe_id
    ON stripe_webhook_events (stripe_event_id);

-- Index for cleanup of old completed events
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_cleanup
    ON stripe_webhook_events (created_at)
    WHERE status = 'completed';

-- Reconciliation tracking table
CREATE TABLE IF NOT EXISTS stripe_reconciliation_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'running'
        CHECK (status IN ('running', 'completed', 'failed')),
    tenants_checked INTEGER DEFAULT 0,
    discrepancies_found INTEGER DEFAULT 0,
    discrepancies_fixed INTEGER DEFAULT 0,
    error TEXT,
    details JSONB
);

CREATE INDEX IF NOT EXISTS idx_stripe_reconciliation_runs_recent
    ON stripe_reconciliation_runs (started_at DESC);

-- +goose Down
DROP TABLE IF EXISTS stripe_reconciliation_runs;
DROP TABLE IF EXISTS stripe_webhook_events;
