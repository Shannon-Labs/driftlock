-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active',
    plan TEXT NOT NULL DEFAULT 'pilot',
    retention_days INTEGER NOT NULL DEFAULT 30,
    default_compressor TEXT NOT NULL DEFAULT 'zstd',
    rate_limit_rps INTEGER NOT NULL DEFAULT 60,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('logs','metrics','traces','llm')),
    description TEXT,
    seed BIGINT NOT NULL DEFAULT 42,
    compressor TEXT NOT NULL DEFAULT 'zstd',
    queue_mode TEXT NOT NULL DEFAULT 'memory',
    retention_days INTEGER NOT NULL DEFAULT 14,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, slug)
);

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin','stream')),
    key_hash TEXT NOT NULL,
    stream_id UUID REFERENCES streams(id) ON DELETE SET NULL,
    rate_limit_rps INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ
);
CREATE INDEX api_keys_tenant_id_idx ON api_keys (tenant_id);

CREATE TABLE stream_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    config JSONB NOT NULL,
    created_by TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (stream_id, version)
);

CREATE TABLE ingest_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    batch_hash TEXT NOT NULL,
    queued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','processing','completed','failed')),
    worker TEXT NOT NULL DEFAULT 'memory',
    error TEXT,
    UNIQUE (tenant_id, stream_id, batch_hash)
);
CREATE INDEX ingest_batches_stream_idx ON ingest_batches (tenant_id, stream_id, queued_at DESC);

CREATE TABLE anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    ingest_batch_id UUID REFERENCES ingest_batches(id) ON DELETE SET NULL,
    ncd NUMERIC(6,4) NOT NULL,
    compression_ratio NUMERIC(8,4) NOT NULL,
    entropy_change NUMERIC(8,4) NOT NULL,
    p_value NUMERIC(6,4) NOT NULL,
    confidence NUMERIC(6,4) NOT NULL,
    explanation TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new','acknowledged','exported')),
    detected_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    details JSONB,
    baseline_snapshot JSONB,
    window_snapshot JSONB
);
CREATE INDEX anomalies_lookup_idx ON anomalies (tenant_id, stream_id, detected_at DESC, id);

CREATE TABLE anomaly_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anomaly_id UUID NOT NULL REFERENCES anomalies(id) ON DELETE CASCADE,
    format TEXT NOT NULL CHECK (format IN ('markdown','html','pdf')),
    uri TEXT NOT NULL,
    checksum TEXT,
    size_bytes BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX anomaly_evidence_anomaly_idx ON anomaly_evidence (anomaly_id);

CREATE TABLE export_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    format TEXT NOT NULL CHECK (format IN ('json','markdown','html','pdf')),
    filters JSONB NOT NULL,
    delivery JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','running','completed','failed')),
    result_uri TEXT,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX export_jobs_tenant_idx ON export_jobs (tenant_id, created_at DESC);

