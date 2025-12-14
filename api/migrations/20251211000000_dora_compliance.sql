-- +goose Up
-- DORA Transaction Monitoring and Compliance (EU Regulation 2022/2554)
-- Implements incident management, transaction tracking, and Driftlog audit trail

-- Main incidents table (DORA Article 10 compliant)
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,

    -- Incident classification (DORA Article 10)
    incident_type TEXT NOT NULL CHECK (incident_type IN (
        'major_ict_incident',       -- Major ICT-related incident
        'significant_cyber_threat', -- Significant cyber threat
        'transaction_anomaly',      -- Transaction processing anomaly
        'data_breach',              -- Personal data breach
        'service_disruption',       -- Service availability issue
        'unauthorized_access',      -- Unauthorized system access
        'data_integrity',           -- Data integrity violation
        'compliance_violation'      -- Regulatory compliance issue
    )),
    severity TEXT NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')) DEFAULT 'medium',
    status TEXT NOT NULL CHECK (status IN (
        'detected',        -- Initial detection
        'investigating',   -- Under investigation
        'classified',      -- Classified per DORA criteria
        'reported',        -- Reported to authorities if required
        'mitigated',       -- Mitigation applied
        'resolved',        -- Incident resolved
        'closed'           -- Case closed
    )) DEFAULT 'detected',

    -- Transaction context
    transaction_id TEXT,
    transaction_type TEXT,  -- 'payment', 'transfer', 'settlement', etc.
    amount FLOAT8,
    currency TEXT,
    sender_account TEXT,
    receiver_account TEXT,

    -- Detection metadata
    risk_score FLOAT8 NOT NULL,
    confidence FLOAT8 NOT NULL,
    detection_method TEXT NOT NULL DEFAULT 'CBAD',
    explanation TEXT NOT NULL,
    recommended_action TEXT,

    -- DORA-specific fields (Article 19)
    regulatory_notification_required BOOLEAN NOT NULL DEFAULT FALSE,
    notification_deadline TIMESTAMPTZ,  -- 24h for initial, 72h for intermediate
    notification_sent_at TIMESTAMPTZ,
    notification_reference TEXT,

    -- Impact assessment (DORA Article 18)
    impact_assessment JSONB,
    affected_clients_count INTEGER,
    financial_impact_eur FLOAT8,

    -- Raw event data for audit
    raw_event JSONB NOT NULL,

    -- Timestamps (DORA timeline requirements)
    detected_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    classification_timestamp TIMESTAMPTZ,
    mitigation_timestamp TIMESTAMPTZ,
    resolution_timestamp TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Junction table: incidents to anomalies (many-to-many)
-- An incident may be correlated from multiple anomalies
CREATE TABLE IF NOT EXISTS incident_anomalies (
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    anomaly_id UUID NOT NULL REFERENCES anomalies(id) ON DELETE CASCADE,
    correlation_type TEXT NOT NULL CHECK (correlation_type IN (
        'primary',      -- Primary triggering anomaly
        'related',      -- Related anomaly (same transaction)
        'correlated',   -- Statistically correlated anomaly
        'subsequent'    -- Follow-on anomaly
    )) DEFAULT 'primary',
    correlation_score FLOAT8,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (incident_id, anomaly_id)
);

-- Driftlog: Complete audit trail of all detection decisions
-- Named after Driftlock + audit log = Driftlog
CREATE TABLE IF NOT EXISTS driftlog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,

    -- Event identification
    event_hash TEXT NOT NULL,  -- SHA-256 of raw_event for deduplication
    transaction_id TEXT,

    -- Detection decision
    decision TEXT NOT NULL CHECK (decision IN (
        'normal',       -- No anomaly detected
        'anomaly',      -- Anomaly detected
        'escalated',    -- Escalated to incident
        'suppressed',   -- Suppressed (below threshold)
        'skipped'       -- Skipped (calibrating/warmup)
    )),

    -- Detection metrics at decision time
    ncd FLOAT8,
    compression_ratio FLOAT8,
    entropy FLOAT8,
    p_value FLOAT8,
    confidence FLOAT8,

    -- Thresholds applied
    ncd_threshold_applied FLOAT8,
    profile_applied TEXT,

    -- References
    anomaly_id UUID REFERENCES anomalies(id) ON DELETE SET NULL,
    incident_id UUID REFERENCES incidents(id) ON DELETE SET NULL,

    -- Audit metadata
    processing_time_us INTEGER,  -- Microseconds
    api_key_id UUID,
    client_ip TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for incidents
CREATE INDEX IF NOT EXISTS incidents_tenant_idx ON incidents (tenant_id, detected_at DESC);
CREATE INDEX IF NOT EXISTS incidents_stream_idx ON incidents (stream_id, detected_at DESC);
CREATE INDEX IF NOT EXISTS incidents_status_idx ON incidents (tenant_id, status, severity);
CREATE INDEX IF NOT EXISTS incidents_transaction_idx ON incidents (tenant_id, transaction_id) WHERE transaction_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS incidents_notification_idx ON incidents (notification_deadline)
    WHERE regulatory_notification_required = TRUE AND notification_sent_at IS NULL;

-- Indexes for driftlog
CREATE INDEX IF NOT EXISTS driftlog_tenant_idx ON driftlog (tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS driftlog_stream_idx ON driftlog (stream_id, created_at DESC);
CREATE INDEX IF NOT EXISTS driftlog_decision_idx ON driftlog (tenant_id, decision, created_at DESC);
CREATE INDEX IF NOT EXISTS driftlog_transaction_idx ON driftlog (tenant_id, transaction_id) WHERE transaction_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS driftlog_hash_idx ON driftlog (tenant_id, event_hash);

-- Indexes for incident_anomalies
CREATE INDEX IF NOT EXISTS incident_anomalies_anomaly_idx ON incident_anomalies (anomaly_id);

-- Add DORA compliance fields to tenants
ALTER TABLE tenants
ADD COLUMN IF NOT EXISTS dora_compliant BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS jurisdiction TEXT;  -- 'EU', 'UK', 'US', etc.

-- Add 'transactions' stream type (extend constraint)
-- First drop old constraint, then add new one with transactions type
ALTER TABLE streams DROP CONSTRAINT IF EXISTS streams_type_check;
ALTER TABLE streams ADD CONSTRAINT streams_type_check
    CHECK (type IN ('logs', 'metrics', 'traces', 'llm', 'transactions'));

-- Documentation
COMMENT ON TABLE incidents IS 'DORA Article 10 compliant ICT-related incident records';
COMMENT ON TABLE incident_anomalies IS 'Links incidents to underlying anomalies for audit trail';
COMMENT ON TABLE driftlog IS 'Complete audit trail of all detection decisions (Driftlog)';
COMMENT ON COLUMN incidents.regulatory_notification_required IS 'Per DORA Article 19: mandatory reporting threshold';
COMMENT ON COLUMN incidents.notification_deadline IS '24h for initial report, 72h for intermediate, 1 month for final';
COMMENT ON COLUMN driftlog.event_hash IS 'SHA-256 hash for deduplication and audit linkage';

-- +goose Down
-- Remove DORA compliance fields from tenants
ALTER TABLE tenants
DROP COLUMN IF EXISTS dora_compliant,
DROP COLUMN IF EXISTS jurisdiction;

-- Drop indexes
DROP INDEX IF EXISTS driftlog_hash_idx;
DROP INDEX IF EXISTS driftlog_transaction_idx;
DROP INDEX IF EXISTS driftlog_decision_idx;
DROP INDEX IF EXISTS driftlog_stream_idx;
DROP INDEX IF EXISTS driftlog_tenant_idx;

DROP INDEX IF EXISTS incident_anomalies_anomaly_idx;

DROP INDEX IF EXISTS incidents_notification_idx;
DROP INDEX IF EXISTS incidents_transaction_idx;
DROP INDEX IF EXISTS incidents_status_idx;
DROP INDEX IF EXISTS incidents_stream_idx;
DROP INDEX IF EXISTS incidents_tenant_idx;

-- Drop tables
DROP TABLE IF EXISTS driftlog;
DROP TABLE IF EXISTS incident_anomalies;
DROP TABLE IF EXISTS incidents;

-- Restore original stream types constraint
ALTER TABLE streams DROP CONSTRAINT IF EXISTS streams_type_check;
ALTER TABLE streams ADD CONSTRAINT streams_type_check
    CHECK (type IN ('logs', 'metrics', 'traces', 'llm'));
