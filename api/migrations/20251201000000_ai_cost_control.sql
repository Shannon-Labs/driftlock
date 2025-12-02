-- +goose Up
-- AI Cost Control Configuration Table
CREATE TABLE IF NOT EXISTS ai_cost_control_configs (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT false,
    models JSONB NOT NULL DEFAULT '[]'::jsonb,
    max_calls_per_day INTEGER NOT NULL DEFAULT 0,
    max_calls_per_hour INTEGER NOT NULL DEFAULT 0,
    max_cost_per_month DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    analysis_threshold DECIMAL(3,2) NOT NULL DEFAULT 0.7 CHECK (analysis_threshold >= 0.0 AND analysis_threshold <= 1.0),
    batch_size INTEGER NOT NULL DEFAULT 50 CHECK (batch_size > 0),
    optimize_for TEXT NOT NULL DEFAULT 'cost' CHECK (optimize_for IN ('speed', 'cost', 'accuracy')),
    notify_threshold DECIMAL(3,2) NOT NULL DEFAULT 0.8 CHECK (notify_threshold >= 0.0 AND notify_threshold <= 1.0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for efficient lookups
CREATE INDEX idx_ai_cost_control_tenant_id ON ai_cost_control_configs(tenant_id);

-- AI Usage Tracking Table
CREATE TABLE IF NOT EXISTS ai_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    model_type TEXT NOT NULL CHECK (model_type IN ('claude-haiku-4-5-20251001', 'claude-sonnet-4-5-20250929', 'claude-opus-4-5-20251101')),
    input_tokens BIGINT NOT NULL DEFAULT 0,
    output_tokens BIGINT NOT NULL DEFAULT 0,
    cost_usd DECIMAL(10,6) NOT NULL DEFAULT 0.000000,
    margin_percent DECIMAL(5,2) NOT NULL DEFAULT 15.00,
    total_charge_usd DECIMAL(10,6) NOT NULL DEFAULT 0.000000,
    request_metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX idx_ai_usage_tenant_id ON ai_usage(tenant_id);
CREATE INDEX idx_ai_usage_tenant_date ON ai_usage(tenant_id, DATE(created_at AT TIME ZONE 'UTC'));
CREATE INDEX idx_ai_usage_tenant_month ON ai_usage(tenant_id, DATE_TRUNC('month', created_at AT TIME ZONE 'UTC'));
CREATE INDEX idx_ai_usage_model_type ON ai_usage(model_type);
CREATE INDEX idx_ai_usage_created_at ON ai_usage(created_at);

-- AI Usage Limits Tracking (for rate limiting)
CREATE TABLE IF NOT EXISTS ai_usage_limits (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    window_type TEXT NOT NULL CHECK (window_type IN ('hour', 'day', 'month')),
    window_start TIMESTAMPTZ NOT NULL,
    call_count INTEGER NOT NULL DEFAULT 0,
    total_cost DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, window_type, window_start)
);

-- Create index for efficient limit checks
CREATE INDEX idx_ai_usage_limits_window ON ai_usage_limits(window_type, window_start);

-- Extend usage_metrics table to include AI metrics
ALTER TABLE usage_metrics
ADD COLUMN IF NOT EXISTS ai_calls_count BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS ai_input_tokens BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS ai_output_tokens BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS ai_cost_usd DECIMAL(10,6) DEFAULT 0.000000,
ADD COLUMN IF NOT EXISTS ai_charge_usd DECIMAL(10,6) DEFAULT 0.000000;

-- Create trigger to update ai_usage_limits
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_ai_usage_limits()
RETURNS TRIGGER AS $$
BEGIN
    -- Update daily limit
    INSERT INTO ai_usage_limits (tenant_id, window_type, window_start, call_count, total_cost)
    VALUES (NEW.tenant_id, 'day', DATE_TRUNC('day', NEW.created_at), 1, NEW.total_charge_usd)
    ON CONFLICT (tenant_id, window_type, window_start) DO UPDATE SET
        call_count = ai_usage_limits.call_count + 1,
        total_cost = ai_usage_limits.total_cost + NEW.total_charge_usd,
        updated_at = NOW();

    -- Update hourly limit
    INSERT INTO ai_usage_limits (tenant_id, window_type, window_start, call_count, total_cost)
    VALUES (NEW.tenant_id, 'hour', DATE_TRUNC('hour', NEW.created_at), 1, NEW.total_charge_usd)
    ON CONFLICT (tenant_id, window_type, window_start) DO UPDATE SET
        call_count = ai_usage_limits.call_count + 1,
        total_cost = ai_usage_limits.total_cost + NEW.total_charge_usd,
        updated_at = NOW();

    -- Update monthly limit
    INSERT INTO ai_usage_limits (tenant_id, window_type, window_start, call_count, total_cost)
    VALUES (NEW.tenant_id, 'month', DATE_TRUNC('month', NEW.created_at), 1, NEW.total_charge_usd)
    ON CONFLICT (tenant_id, window_type, window_start) DO UPDATE SET
        call_count = ai_usage_limits.call_count + 1,
        total_cost = ai_usage_limits.total_cost + NEW.total_charge_usd,
        updated_at = NOW();

    -- Update daily usage_metrics
    INSERT INTO usage_metrics (tenant_id, stream_id, date, ai_calls_count, ai_input_tokens, ai_output_tokens, ai_cost_usd, ai_charge_usd)
    VALUES (NEW.tenant_id, NEW.stream_id, DATE(NEW.created_at), 1, NEW.input_tokens, NEW.output_tokens, NEW.cost_usd, NEW.total_charge_usd)
    ON CONFLICT (tenant_id, stream_id, date) DO UPDATE SET
        ai_calls_count = usage_metrics.ai_calls_count + 1,
        ai_input_tokens = usage_metrics.ai_input_tokens + NEW.input_tokens,
        ai_output_tokens = usage_metrics.ai_output_tokens + NEW.output_tokens,
        ai_cost_usd = usage_metrics.ai_cost_usd + NEW.cost_usd,
        ai_charge_usd = usage_metrics.ai_charge_usd + NEW.total_charge_usd;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create trigger
DROP TRIGGER IF EXISTS trigger_update_ai_usage_limits ON ai_usage;
CREATE TRIGGER trigger_update_ai_usage_limits
    AFTER INSERT ON ai_usage
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_usage_limits();

-- Add comment for documentation
COMMENT ON TABLE ai_cost_control_configs IS 'Configurable AI cost control settings per tenant';
COMMENT ON TABLE ai_usage IS 'Tracks all AI API usage for billing and analytics';
COMMENT ON TABLE ai_usage_limits IS 'Tracks usage within time windows for rate limiting';

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_ai_usage_limits ON ai_usage;
DROP FUNCTION IF EXISTS update_ai_usage_limits();
ALTER TABLE usage_metrics DROP COLUMN IF EXISTS ai_calls_count;
ALTER TABLE usage_metrics DROP COLUMN IF EXISTS ai_input_tokens;
ALTER TABLE usage_metrics DROP COLUMN IF EXISTS ai_output_tokens;
ALTER TABLE usage_metrics DROP COLUMN IF EXISTS ai_cost_usd;
ALTER TABLE usage_metrics DROP COLUMN IF EXISTS ai_charge_usd;
DROP TABLE IF EXISTS ai_usage_limits;
DROP TABLE IF EXISTS ai_usage;
DROP TABLE IF EXISTS ai_cost_control_configs;