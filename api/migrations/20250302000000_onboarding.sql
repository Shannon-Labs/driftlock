-- +goose Up
-- Add onboarding and email fields to tenants table

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS email TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS signup_ip TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS signup_source TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS verification_token TEXT;

-- Create unique index on email (only for non-null emails)
CREATE UNIQUE INDEX IF NOT EXISTS tenants_email_unique_idx ON tenants (email) WHERE email IS NOT NULL;

-- Create usage_metrics table for tracking
CREATE TABLE IF NOT EXISTS usage_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID REFERENCES streams(id) ON DELETE SET NULL,
    metric_date DATE NOT NULL,
    event_count BIGINT NOT NULL DEFAULT 0,
    anomaly_count BIGINT NOT NULL DEFAULT 0,
    api_request_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, stream_id, metric_date)
);

CREATE INDEX IF NOT EXISTS usage_metrics_tenant_date_idx ON usage_metrics (tenant_id, metric_date DESC);

-- +goose Down
DROP INDEX IF EXISTS usage_metrics_tenant_date_idx;
DROP TABLE IF EXISTS usage_metrics;
DROP INDEX IF EXISTS tenants_email_unique_idx;
ALTER TABLE tenants DROP COLUMN IF EXISTS verification_token;
ALTER TABLE tenants DROP COLUMN IF EXISTS verified_at;
ALTER TABLE tenants DROP COLUMN IF EXISTS signup_source;
ALTER TABLE tenants DROP COLUMN IF EXISTS signup_ip;
ALTER TABLE tenants DROP COLUMN IF EXISTS email;
