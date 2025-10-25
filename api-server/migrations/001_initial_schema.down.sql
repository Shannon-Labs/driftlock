-- Rollback script for initial schema

-- Drop views
DROP VIEW IF EXISTS anomaly_stats_by_stream;
DROP VIEW IF EXISTS recent_significant_anomalies;

-- Drop triggers
DROP TRIGGER IF EXISTS update_detection_config_updated_at ON detection_config;
DROP TRIGGER IF EXISTS update_anomalies_updated_at ON anomalies;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS performance_metrics;
DROP TABLE IF EXISTS detection_config;
DROP TABLE IF EXISTS anomalies;

-- Drop extension (only if not used by other databases)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
