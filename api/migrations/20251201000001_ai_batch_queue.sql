-- +goose Up
-- AI Batch Queue Table
CREATE TABLE IF NOT EXISTS ai_batch_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    model_type TEXT NOT NULL,
    event_payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    batch_id TEXT, -- ID returned by Anthropic Batch API
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

-- Create indexes for efficient polling
CREATE INDEX idx_ai_batch_queue_status_created ON ai_batch_queue(status, created_at);
CREATE INDEX idx_ai_batch_queue_tenant_id ON ai_batch_queue(tenant_id);
CREATE INDEX idx_ai_batch_queue_batch_id ON ai_batch_queue(batch_id);

-- Add comment for documentation
COMMENT ON TABLE ai_batch_queue is 'Buffer for AI requests to be processed in batch (Enterprise only)';

-- +goose Down
DROP TABLE IF EXISTS ai_batch_queue;
