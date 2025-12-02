package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BatchService handles buffering of AI requests for batch processing
type BatchService struct {
	pool *pgxpool.Pool
}

// NewBatchService creates a new batch service
func NewBatchService(pool *pgxpool.Pool) *BatchService {
	return &BatchService{
		pool: pool,
	}
}

// QueueRequest adds an AI request to the batch queue
func (s *BatchService) QueueRequest(ctx context.Context, tenantID string, model string, payload interface{}) error {
	// Serialize payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO ai_batch_queue (tenant_id, model_type, event_payload, status)
		VALUES ($1, $2, $3, 'pending')
	`

	_, err = s.pool.Exec(ctx, query, tenantID, model, jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to queue request: %w", err)
	}

	return nil
}

// GetPendingRequests retrieves pending requests for processing
func (s *BatchService) GetPendingRequests(ctx context.Context, limit int) ([]BatchRequest, error) {
	query := `
		SELECT id, tenant_id, model_type, event_payload, created_at
		FROM ai_batch_queue
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending requests: %w", err)
	}
	defer rows.Close()

	var requests []BatchRequest
	for rows.Next() {
		var req BatchRequest
		if err := rows.Scan(&req.ID, &req.TenantID, &req.ModelType, &req.Payload, &req.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// MarkAsProcessing updates status to processing
func (s *BatchService) MarkAsProcessing(ctx context.Context, ids []string, batchID string) error {
	query := `
		UPDATE ai_batch_queue
		SET status = 'processing', batch_id = $1, updated_at = NOW()
		WHERE id = ANY($2)
	`

	_, err := s.pool.Exec(ctx, query, batchID, ids)
	if err != nil {
		return fmt.Errorf("failed to mark requests as processing: %w", err)
	}

	return nil
}

// BatchRequest represents a queued request
type BatchRequest struct {
	ID        string          `json:"id"`
	TenantID  string          `json:"tenant_id"`
	ModelType string          `json:"model_type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
