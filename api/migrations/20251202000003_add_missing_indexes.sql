-- +goose Up
-- Add missing indexes for foreign key columns to improve query performance

-- Index on api_keys.stream_id for efficient stream-scoped key lookups
CREATE INDEX IF NOT EXISTS idx_api_keys_stream_id ON api_keys(stream_id);

-- Index on ai_usage.stream_id for efficient stream-level AI usage queries
CREATE INDEX IF NOT EXISTS idx_ai_usage_stream_id ON ai_usage(stream_id);

-- Note: ai_batch_queue already has idx_ai_batch_queue_status_created on (status, created_at)
-- which was created in 20251201000001_ai_batch_queue.sql

-- +goose Down
DROP INDEX IF EXISTS idx_api_keys_stream_id;
DROP INDEX IF EXISTS idx_ai_usage_stream_id;
