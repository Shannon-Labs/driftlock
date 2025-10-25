-- Driftlock Initial Database Schema
-- This migration creates the core tables for the anomaly detection system

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Anomalies table - stores detected anomalies with CBAD metrics
CREATE TABLE IF NOT EXISTS anomalies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMPTZ NOT NULL,

    -- Stream identification
    stream_type VARCHAR(50) NOT NULL CHECK (stream_type IN ('logs', 'metrics', 'traces', 'llm')),

    -- Core CBAD metrics
    ncd_score FLOAT NOT NULL CHECK (ncd_score >= 0 AND ncd_score <= 1),
    p_value FLOAT NOT NULL CHECK (p_value >= 0 AND p_value <= 1),

    -- Status management
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'acknowledged', 'dismissed', 'investigating')),

    -- Explanations
    glass_box_explanation TEXT NOT NULL,
    detailed_explanation TEXT,

    -- Compression metrics
    compression_baseline FLOAT NOT NULL,
    compression_window FLOAT NOT NULL,
    compression_combined FLOAT NOT NULL,
    compression_ratio_change FLOAT,

    -- Entropy metrics (optional)
    baseline_entropy FLOAT,
    window_entropy FLOAT,
    entropy_change FLOAT,

    -- Statistical significance
    confidence_level FLOAT NOT NULL CHECK (confidence_level >= 0 AND confidence_level <= 1),
    is_statistically_significant BOOLEAN NOT NULL DEFAULT false,

    -- Data payloads (stored as JSONB for queryability)
    baseline_data JSONB,
    window_data JSONB,
    metadata JSONB,

    -- Tags for categorization
    tags TEXT[] DEFAULT '{}',

    -- User interaction tracking
    acknowledged_by VARCHAR(255),
    acknowledged_at TIMESTAMPTZ,
    dismissed_by VARCHAR(255),
    dismissed_at TIMESTAMPTZ,
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_anomalies_timestamp ON anomalies(timestamp DESC);
CREATE INDEX idx_anomalies_status ON anomalies(status);
CREATE INDEX idx_anomalies_stream_type ON anomalies(stream_type);
CREATE INDEX idx_anomalies_p_value ON anomalies(p_value);
CREATE INDEX idx_anomalies_ncd_score ON anomalies(ncd_score);
CREATE INDEX idx_anomalies_significant ON anomalies(is_statistically_significant) WHERE is_statistically_significant = true;
CREATE INDEX idx_anomalies_created_at ON anomalies(created_at DESC);
CREATE INDEX idx_anomalies_tags ON anomalies USING GIN(tags);
CREATE INDEX idx_anomalies_metadata ON anomalies USING GIN(metadata);

-- Detection configuration table
CREATE TABLE IF NOT EXISTS detection_config (
    id SERIAL PRIMARY KEY,

    -- Threshold settings
    ncd_threshold FLOAT NOT NULL DEFAULT 0.5 CHECK (ncd_threshold >= 0 AND ncd_threshold <= 1),
    p_value_threshold FLOAT NOT NULL DEFAULT 0.05 CHECK (p_value_threshold >= 0 AND p_value_threshold <= 1),

    -- Window settings
    baseline_size INTEGER NOT NULL DEFAULT 100 CHECK (baseline_size > 0),
    window_size INTEGER NOT NULL DEFAULT 10 CHECK (window_size > 0),
    hop_size INTEGER NOT NULL DEFAULT 1 CHECK (hop_size > 0),

    -- Per-stream overrides (JSONB for flexibility)
    stream_overrides JSONB DEFAULT '{}',

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT false
);

-- Ensure only one active configuration at a time
CREATE UNIQUE INDEX idx_detection_config_active ON detection_config(is_active) WHERE is_active = true;

-- Insert default configuration
INSERT INTO detection_config (ncd_threshold, p_value_threshold, baseline_size, window_size, hop_size, is_active, notes)
VALUES (0.5, 0.05, 100, 10, 1, true, 'Default CBAD configuration');

-- Performance metrics table for monitoring API performance
CREATE TABLE IF NOT EXISTS performance_metrics (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Metric identification
    metric_type VARCHAR(50) NOT NULL,
    endpoint VARCHAR(255),

    -- Performance data
    duration_ms INTEGER NOT NULL,
    success BOOLEAN NOT NULL DEFAULT true,
    error_message TEXT,

    -- Additional metadata
    metadata JSONB DEFAULT '{}'
);

-- Index for time-series queries
CREATE INDEX idx_performance_metrics_timestamp ON performance_metrics(timestamp DESC);
CREATE INDEX idx_performance_metrics_endpoint ON performance_metrics(endpoint);
CREATE INDEX idx_performance_metrics_metric_type ON performance_metrics(metric_type);

-- Audit log table for tracking changes
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Actor information
    username VARCHAR(255) NOT NULL,
    ip_address INET,

    -- Action details
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,

    -- Change tracking
    old_values JSONB,
    new_values JSONB,

    -- Metadata
    metadata JSONB DEFAULT '{}'
);

-- Index for audit queries
CREATE INDEX idx_audit_log_timestamp ON audit_log(timestamp DESC);
CREATE INDEX idx_audit_log_username ON audit_log(username);
CREATE INDEX idx_audit_log_resource ON audit_log(resource_type, resource_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at on anomalies table
CREATE TRIGGER update_anomalies_updated_at
    BEFORE UPDATE ON anomalies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to auto-update updated_at on detection_config table
CREATE TRIGGER update_detection_config_updated_at
    BEFORE UPDATE ON detection_config
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Views for common queries

-- View for recent significant anomalies
CREATE OR REPLACE VIEW recent_significant_anomalies AS
SELECT
    id, timestamp, stream_type, ncd_score, p_value, status,
    glass_box_explanation, compression_ratio_change,
    tags, created_at
FROM anomalies
WHERE is_statistically_significant = true
    AND p_value < 0.05
ORDER BY timestamp DESC
LIMIT 100;

-- View for anomaly statistics by stream type
CREATE OR REPLACE VIEW anomaly_stats_by_stream AS
SELECT
    stream_type,
    COUNT(*) as total_anomalies,
    COUNT(*) FILTER (WHERE is_statistically_significant = true) as significant_anomalies,
    AVG(ncd_score) as avg_ncd_score,
    AVG(p_value) as avg_p_value,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
    COUNT(*) FILTER (WHERE status = 'acknowledged') as acknowledged_count,
    COUNT(*) FILTER (WHERE status = 'dismissed') as dismissed_count
FROM anomalies
GROUP BY stream_type;

-- Grant permissions (adjust as needed for your environment)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO driftlock_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO driftlock_app;
