-- Migration 001: Initial Driftlock schema for anomaly detection
-- Description: Creates core anomalies table with comprehensive glass-box explanation fields

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Anomalies table: Stores all detected anomalies with compression metrics
CREATE TABLE anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    stream_type TEXT NOT NULL, -- 'logs', 'metrics', 'traces', 'llm'

    -- Core CBAD metrics
    ncd_score DOUBLE PRECISION NOT NULL,
    p_value DOUBLE PRECISION NOT NULL,

    -- Status tracking
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'acknowledged', 'dismissed', 'investigating'

    -- Glass-box explanation
    glass_box_explanation TEXT NOT NULL,
    detailed_explanation TEXT,

    -- Compression metrics
    compression_baseline DOUBLE PRECISION NOT NULL,
    compression_window DOUBLE PRECISION NOT NULL,
    compression_combined DOUBLE PRECISION NOT NULL,
    compression_ratio_change DOUBLE PRECISION NOT NULL, -- Percentage change

    -- Entropy metrics
    baseline_entropy DOUBLE PRECISION,
    window_entropy DOUBLE PRECISION,
    entropy_change DOUBLE PRECISION,

    -- Statistical significance
    confidence_level DOUBLE PRECISION NOT NULL,
    is_statistically_significant BOOLEAN NOT NULL DEFAULT false,

    -- Data payloads (JSONB for efficient querying)
    baseline_data JSONB,
    window_data JSONB,
    metadata JSONB,

    -- Tags for categorization (array for efficient filtering)
    tags TEXT[],

    -- User interaction
    acknowledged_by TEXT,
    acknowledged_at TIMESTAMPTZ,
    dismissed_by TEXT,
    dismissed_at TIMESTAMPTZ,
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_anomalies_timestamp ON anomalies(timestamp DESC);
CREATE INDEX idx_anomalies_stream_type ON anomalies(stream_type);
CREATE INDEX idx_anomalies_status ON anomalies(status);
CREATE INDEX idx_anomalies_p_value ON anomalies(p_value);
CREATE INDEX idx_anomalies_ncd_score ON anomalies(ncd_score DESC);
CREATE INDEX idx_anomalies_created_at ON anomalies(created_at DESC);
CREATE INDEX idx_anomalies_tags ON anomalies USING GIN(tags);
CREATE INDEX idx_anomalies_metadata ON anomalies USING GIN(metadata);

-- Composite index for common query patterns
CREATE INDEX idx_anomalies_stream_status_time ON anomalies(stream_type, status, timestamp DESC);
CREATE INDEX idx_anomalies_significance ON anomalies(is_statistically_significant, p_value) WHERE is_statistically_significant = true;

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_anomalies_updated_at
    BEFORE UPDATE ON anomalies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Configuration table: Stores detection thresholds and settings
CREATE TABLE detection_config (
    id SERIAL PRIMARY KEY,
    ncd_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.3,
    p_value_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.05,
    baseline_size INT NOT NULL DEFAULT 100,
    window_size INT NOT NULL DEFAULT 50,
    hop_size INT NOT NULL DEFAULT 10,

    -- Stream-specific overrides (JSONB)
    stream_overrides JSONB,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by TEXT,
    notes TEXT,

    -- Only allow one active config
    is_active BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(is_active) WHERE is_active = true
);

CREATE TRIGGER update_detection_config_updated_at
    BEFORE UPDATE ON detection_config
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default configuration
INSERT INTO detection_config (ncd_threshold, p_value_threshold, baseline_size, window_size, hop_size, created_by, notes)
VALUES (0.3, 0.05, 100, 50, 10, 'system', 'Default CBAD configuration');

-- Performance metrics table: Tracks API and system performance
CREATE TABLE performance_metrics (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metric_type TEXT NOT NULL, -- 'api_request', 'cbad_computation', 'database_query'
    endpoint TEXT,
    duration_ms DOUBLE PRECISION NOT NULL,
    success BOOLEAN NOT NULL DEFAULT true,
    error_message TEXT,
    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_performance_metrics_timestamp ON performance_metrics(timestamp DESC);
CREATE INDEX idx_performance_metrics_type ON performance_metrics(metric_type);
CREATE INDEX idx_performance_metrics_endpoint ON performance_metrics(endpoint);

-- Retention policy: Keep performance metrics for 30 days
CREATE INDEX idx_performance_metrics_created_at ON performance_metrics(created_at);

-- API keys table for authentication
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_hash TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,

    -- Permissions
    role TEXT NOT NULL DEFAULT 'viewer', -- 'admin', 'analyst', 'viewer'
    scopes TEXT[] NOT NULL DEFAULT ARRAY['read:anomalies'],

    -- Lifecycle
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by TEXT,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_used_at TIMESTAMPTZ,

    -- Rate limiting
    rate_limit_per_minute INT DEFAULT 100
);

CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_active ON api_keys(is_active) WHERE is_active = true;

-- Comments for documentation
COMMENT ON TABLE anomalies IS 'Stores all detected anomalies from CBAD analysis with compression-based metrics';
COMMENT ON COLUMN anomalies.ncd_score IS 'Normalized Compression Distance (0-1): measures similarity between baseline and window';
COMMENT ON COLUMN anomalies.p_value IS 'Statistical significance (0-1): probability that difference is random';
COMMENT ON COLUMN anomalies.glass_box_explanation IS 'Human-readable explanation of why anomaly was detected';
COMMENT ON COLUMN anomalies.baseline_data IS 'JSON payload of baseline data used for comparison';
COMMENT ON COLUMN anomalies.window_data IS 'JSON payload of anomalous window data';

COMMENT ON TABLE detection_config IS 'Global and stream-specific CBAD detection configuration';
COMMENT ON TABLE performance_metrics IS 'API and system performance tracking for monitoring and optimization';
COMMENT ON TABLE api_keys IS 'API key authentication and authorization';
