-- Evidence-Native Compliance: Add fields for full audit reproducibility
-- This enables DORA/SOX/HIPAA compliance with provable detection decisions

-- Add evidence-native fields to anomalies table
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS detector_config_hash TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS cbad_core_version TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS api_version TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS baseline_anchor_id UUID REFERENCES stream_anchors(id);
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS baseline_events_hash TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS window_events_hash TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS ncd_ci_lower FLOAT8;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS ncd_ci_upper FLOAT8;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS composite_score FLOAT8;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS conditional_novelty FLOAT8;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS processing_seed BIGINT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS permutation_count INTEGER;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS compression_algorithm TEXT;
ALTER TABLE anomalies ADD COLUMN IF NOT EXISTS tokenizer_config JSONB;

-- Index for evidence queries
CREATE INDEX IF NOT EXISTS idx_anomalies_baseline_anchor ON anomalies(baseline_anchor_id);
CREATE INDEX IF NOT EXISTS idx_anomalies_config_hash ON anomalies(detector_config_hash);
CREATE INDEX IF NOT EXISTS idx_anomalies_cbad_version ON anomalies(cbad_core_version);

-- Add comment for documentation
COMMENT ON COLUMN anomalies.detector_config_hash IS 'SHA-256 hash of detector configuration at detection time';
COMMENT ON COLUMN anomalies.cbad_core_version IS 'cbad-core crate version used for detection';
COMMENT ON COLUMN anomalies.api_version IS 'API server version';
COMMENT ON COLUMN anomalies.baseline_anchor_id IS 'Reference to stream_anchors.id - baseline snapshot used';
COMMENT ON COLUMN anomalies.baseline_events_hash IS 'SHA-256 hash of baseline data for reproducibility';
COMMENT ON COLUMN anomalies.window_events_hash IS 'SHA-256 hash of window data for reproducibility';
COMMENT ON COLUMN anomalies.ncd_ci_lower IS '95% confidence interval lower bound for NCD';
COMMENT ON COLUMN anomalies.ncd_ci_upper IS '95% confidence interval upper bound for NCD';
COMMENT ON COLUMN anomalies.composite_score IS 'Weighted fusion score from multiple detectors';
COMMENT ON COLUMN anomalies.conditional_novelty IS 'Conditional novelty score: (C(b+w)-C(b))/C(w)';
COMMENT ON COLUMN anomalies.processing_seed IS 'Random seed used for deterministic replay';
COMMENT ON COLUMN anomalies.permutation_count IS 'Number of permutations used for p-value computation';
COMMENT ON COLUMN anomalies.compression_algorithm IS 'Compression algorithm used (zstd, zlab, etc.)';
COMMENT ON COLUMN anomalies.tokenizer_config IS 'Tokenizer configuration snapshot as JSON';
