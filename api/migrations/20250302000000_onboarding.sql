-- +goose Up
-- Add onboarding fields to tenants table

ALTER TABLE tenants 
ADD COLUMN IF NOT EXISTS email TEXT,
ADD COLUMN IF NOT EXISTS signup_ip INET,
ADD COLUMN IF NOT EXISTS verification_token TEXT,
ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS trial_ends_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT,
ADD COLUMN IF NOT EXISTS plan_started_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS current_period_end TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS signup_source TEXT;

-- Create index for email lookups
CREATE INDEX IF NOT EXISTS tenants_email_idx ON tenants (email);

-- Create index for trial tracking
CREATE INDEX IF NOT EXISTS tenants_trial_ends_idx ON tenants (trial_ends_at) 
WHERE trial_ends_at IS NOT NULL;

-- Add usage tracking table
CREATE TABLE IF NOT EXISTS usage_metrics (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    event_count BIGINT NOT NULL DEFAULT 0,
    api_request_count BIGINT NOT NULL DEFAULT 0,
    anomaly_count BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (tenant_id, stream_id, date)
);

CREATE INDEX IF NOT EXISTS usage_metrics_date_idx ON usage_metrics (date DESC);
CREATE INDEX IF NOT EXISTS usage_metrics_tenant_idx ON usage_metrics (tenant_id, date DESC);

-- Add stripe customer tracking
CREATE TABLE IF NOT EXISTS stripe_customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
    stripe_customer_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
ALTER TABLE tenants 
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS signup_ip,
DROP COLUMN IF EXISTS verification_token,
DROP COLUMN IF EXISTS verified_at,
DROP COLUMN IF EXISTS trial_ends_at,
DROP COLUMN IF EXISTS stripe_customer_id,
DROP COLUMN IF EXISTS plan_started_at,
DROP COLUMN IF EXISTS current_period_end,
DROP COLUMN IF EXISTS signup_source;

DROP INDEX IF EXISTS tenants_email_idx;
DROP INDEX IF EXISTS tenants_trial_ends_idx;

DROP TABLE IF EXISTS usage_metrics;
DROP TABLE IF EXISTS stripe_customers;