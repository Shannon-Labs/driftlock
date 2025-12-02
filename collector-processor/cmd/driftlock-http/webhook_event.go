package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WebhookEventStatus represents the processing status of a webhook event
type WebhookEventStatus string

const (
	StatusPending    WebhookEventStatus = "pending"
	StatusProcessing WebhookEventStatus = "processing"
	StatusCompleted  WebhookEventStatus = "completed"
	StatusFailed     WebhookEventStatus = "failed"
	StatusDeadLetter WebhookEventStatus = "dead_letter"
)

// WebhookEvent represents a stored Stripe webhook event
type WebhookEvent struct {
	ID            uuid.UUID
	StripeEventID string
	EventType     string
	EventData     json.RawMessage
	Status        WebhookEventStatus
	RetryCount    int
	MaxRetries    int
	NextRetryAt   *time.Time
	LastError     *string
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}

// WebhookRetryConfig holds retry behavior configuration
type WebhookRetryConfig struct {
	MaxRetries      int           // Default: 5
	InitialBackoff  time.Duration // Default: 1 minute
	MaxBackoff      time.Duration // Default: 1 hour
	BackoffMultiply float64       // Default: 2.0
}

// DefaultRetryConfig returns sensible defaults for retry behavior
func DefaultRetryConfig() WebhookRetryConfig {
	return WebhookRetryConfig{
		MaxRetries:      5,
		InitialBackoff:  1 * time.Minute,
		MaxBackoff:      1 * time.Hour,
		BackoffMultiply: 2.0,
	}
}

// WebhookEventStore handles persistence of webhook events
type WebhookEventStore struct {
	pool   *pgxpool.Pool
	config WebhookRetryConfig
}

// NewWebhookEventStore creates a new webhook event store
func NewWebhookEventStore(pool *pgxpool.Pool, config WebhookRetryConfig) *WebhookEventStore {
	return &WebhookEventStore{
		pool:   pool,
		config: config,
	}
}

// StoreEvent stores a webhook event for processing.
// Returns the event ID and whether it was newly created (vs duplicate).
// This is idempotent - duplicate stripe_event_id returns existing record.
func (s *WebhookEventStore) StoreEvent(ctx context.Context, stripeEventID, eventType string, eventData json.RawMessage) (uuid.UUID, bool, error) {
	var id uuid.UUID

	err := s.pool.QueryRow(ctx, `
		INSERT INTO stripe_webhook_events (stripe_event_id, event_type, event_data, max_retries)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (stripe_event_id) DO NOTHING
		RETURNING id
	`, stripeEventID, eventType, eventData, s.config.MaxRetries).Scan(&id)

	if err == pgx.ErrNoRows {
		// Event already exists, fetch existing ID
		err = s.pool.QueryRow(ctx, `
			SELECT id FROM stripe_webhook_events WHERE stripe_event_id = $1
		`, stripeEventID).Scan(&id)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("failed to fetch existing event: %w", err)
		}
		return id, false, nil
	}
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("failed to store event: %w", err)
	}

	return id, true, nil
}

// GetEvent retrieves a webhook event by ID
func (s *WebhookEventStore) GetEvent(ctx context.Context, id uuid.UUID) (*WebhookEvent, error) {
	var event WebhookEvent
	err := s.pool.QueryRow(ctx, `
		SELECT id, stripe_event_id, event_type, event_data, status, retry_count, max_retries,
		       next_retry_at, last_error, created_at, processed_at
		FROM stripe_webhook_events
		WHERE id = $1
	`, id).Scan(
		&event.ID, &event.StripeEventID, &event.EventType, &event.EventData,
		&event.Status, &event.RetryCount, &event.MaxRetries,
		&event.NextRetryAt, &event.LastError, &event.CreatedAt, &event.ProcessedAt,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// MarkProcessing marks an event as being processed (with optimistic locking)
// Uses FOR UPDATE SKIP LOCKED to allow concurrent workers
func (s *WebhookEventStore) MarkProcessing(ctx context.Context, id uuid.UUID) error {
	result, err := s.pool.Exec(ctx, `
		UPDATE stripe_webhook_events
		SET status = 'processing', updated_at = NOW()
		WHERE id = $1 AND status IN ('pending', 'failed')
	`, id)
	if err != nil {
		return fmt.Errorf("failed to mark processing: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("event not found or already processing")
	}
	return nil
}

// MarkCompleted marks an event as successfully processed
func (s *WebhookEventStore) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE stripe_webhook_events
		SET status = 'completed', processed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// MarkFailed marks an event as failed and schedules retry
// If max retries exceeded, moves to dead_letter status
func (s *WebhookEventStore) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	// Calculate next retry time with exponential backoff
	var retryCount, maxRetries int
	err := s.pool.QueryRow(ctx, `
		SELECT retry_count, max_retries FROM stripe_webhook_events WHERE id = $1
	`, id).Scan(&retryCount, &maxRetries)
	if err != nil {
		return fmt.Errorf("failed to get retry count: %w", err)
	}

	newRetryCount := retryCount + 1
	var newStatus WebhookEventStatus
	var nextRetry *time.Time

	if newRetryCount >= maxRetries {
		newStatus = StatusDeadLetter
		log.Printf("Webhook event %s moved to dead_letter after %d retries: %s", id, newRetryCount, errMsg)
	} else {
		newStatus = StatusFailed
		// Exponential backoff: 1min, 2min, 4min, 8min, 16min
		backoff := s.config.InitialBackoff * time.Duration(math.Pow(s.config.BackoffMultiply, float64(retryCount)))
		if backoff > s.config.MaxBackoff {
			backoff = s.config.MaxBackoff
		}
		t := time.Now().Add(backoff)
		nextRetry = &t
	}

	_, err = s.pool.Exec(ctx, `
		UPDATE stripe_webhook_events
		SET status = $2, retry_count = $3, next_retry_at = $4, last_error = $5, updated_at = NOW()
		WHERE id = $1
	`, id, newStatus, newRetryCount, nextRetry, errMsg)
	return err
}

// GetEventsForRetry fetches events due for retry processing
// Uses FOR UPDATE SKIP LOCKED to allow concurrent workers across Cloud Run instances
func (s *WebhookEventStore) GetEventsForRetry(ctx context.Context, batchSize int) ([]*WebhookEvent, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, stripe_event_id, event_type, event_data, status, retry_count, max_retries,
		       next_retry_at, last_error, created_at, processed_at
		FROM stripe_webhook_events
		WHERE status IN ('pending', 'failed')
		  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events for retry: %w", err)
	}
	defer rows.Close()

	var events []*WebhookEvent
	for rows.Next() {
		var event WebhookEvent
		err := rows.Scan(
			&event.ID, &event.StripeEventID, &event.EventType, &event.EventData,
			&event.Status, &event.RetryCount, &event.MaxRetries,
			&event.NextRetryAt, &event.LastError, &event.CreatedAt, &event.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

// GetDeadLetterEvents fetches events that have exhausted retries
func (s *WebhookEventStore) GetDeadLetterEvents(ctx context.Context, since time.Time, limit int) ([]*WebhookEvent, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, stripe_event_id, event_type, event_data, status, retry_count, max_retries,
		       next_retry_at, last_error, created_at, processed_at
		FROM stripe_webhook_events
		WHERE status = 'dead_letter' AND updated_at >= $1
		ORDER BY updated_at DESC
		LIMIT $2
	`, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*WebhookEvent
	for rows.Next() {
		var event WebhookEvent
		err := rows.Scan(
			&event.ID, &event.StripeEventID, &event.EventType, &event.EventData,
			&event.Status, &event.RetryCount, &event.MaxRetries,
			&event.NextRetryAt, &event.LastError, &event.CreatedAt, &event.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

// CleanupOldEvents removes completed events older than the specified duration
func (s *WebhookEventStore) CleanupOldEvents(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result, err := s.pool.Exec(ctx, `
		DELETE FROM stripe_webhook_events
		WHERE status = 'completed' AND created_at < $1
	`, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Stats returns queue statistics
type WebhookQueueStats struct {
	Pending    int64
	Processing int64
	Failed     int64
	DeadLetter int64
	Completed  int64
}

func (s *WebhookEventStore) Stats(ctx context.Context) (*WebhookQueueStats, error) {
	var stats WebhookQueueStats
	err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'processing') as processing,
			COUNT(*) FILTER (WHERE status = 'failed') as failed,
			COUNT(*) FILTER (WHERE status = 'dead_letter') as dead_letter,
			COUNT(*) FILTER (WHERE status = 'completed') as completed
		FROM stripe_webhook_events
	`).Scan(&stats.Pending, &stats.Processing, &stats.Failed, &stats.DeadLetter, &stats.Completed)
	return &stats, err
}
