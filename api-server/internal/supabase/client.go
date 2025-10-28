package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bytes"

	"github.com/go-redis/redis/v8"
	"github.com/supabase-community/supabase-go"
)

// Client wraps the Supabase client with additional functionality
type Client struct {
	supabase  *supabase.Client
	redis     *redis.Client
	projectID string
}

// Config holds Supabase configuration
type Config struct {
	ProjectID      string
	AnonKey        string
	ServiceRoleKey string
	BaseURL        string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
}

// NewClient creates a new Supabase client
func NewClient(cfg Config) (*Client, error) {
	// Initialize Supabase client
	supabaseClient, err := supabase.NewClient(cfg.BaseURL, cfg.AnonKey, cfg.ServiceRoleKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	// Initialize Redis client for caching
	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
	}

	return &Client{
		supabase:  supabaseClient,
		redis:     redisClient,
		projectID: cfg.ProjectID,
	}, nil
}

// CreateAnomaly creates an anomaly record in Supabase
func (c *Client) CreateAnomaly(ctx context.Context, anomaly map[string]interface{}) error {
	_, err := c.supabase.DB.From("anomalies").Insert(anomaly).Execute()
	if err != nil {
		return fmt.Errorf("failed to create anomaly: %w", err)
	}

	// Invalidate cache if Redis is available
	if c.redis != nil {
		cacheKey := fmt.Sprintf("anomalies:list:%s", c.projectID)
		c.redis.Del(ctx, cacheKey)
	}

	return nil
}

// GetAnomalies retrieves anomalies from Supabase with optional filtering
func (c *Client) GetAnomalies(ctx context.Context, filter map[string]interface{}) ([]map[string]interface{}, error) {
	query := c.supabase.DB.From("anomalies").Select("*")

	// Apply filters if provided
	if len(filter) > 0 {
		for key, value := range filter {
			query = query.Eq(key, value)
		}
	}

	var results []map[string]interface{}
	err := query.Execute(&results)
	if err != nil {
		return nil, fmt.Errorf("failed to get anomalies: %w", err)
	}

	return results, nil
}

// UpdateAnomalyStatus updates an anomaly's status in Supabase
func (c *Client) UpdateAnomalyStatus(ctx context.Context, id string, status string) error {
	_, err := c.supabase.DB.From("anomalies").Update(map[string]interface{}{
		"status": status,
	}).Eq("id", id).Execute()
	if err != nil {
		return fmt.Errorf("failed to update anomaly status: %w", err)
	}

	// Invalidate cache if Redis is available
	if c.redis != nil {
		cacheKey := fmt.Sprintf("anomalies:list:%s", c.projectID)
		c.redis.Del(ctx, cacheKey)
	}

	return nil
}

// CreateUsageRecord creates a usage record for billing
func (c *Client) CreateUsageRecord(ctx context.Context, usage map[string]interface{}) error {
	_, err := c.supabase.DB.From("usage_records").Insert(usage).Execute()
	if err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	return nil
}

// GetOrganization retrieves organization details
func (c *Client) GetOrganization(ctx context.Context, orgID string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.supabase.DB.From("organizations").Select("*").Eq("id", orgID).Single().Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return result, nil
}

// NotifyWebhook sends a webhook notification to Supabase Edge Function
func (c *Client) NotifyWebhook(ctx context.Context, eventType string, payload interface{}) error {
	// Get webhook URL from environment
	webhookURL := os.Getenv("SUPABASE_WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("SUPABASE_WEBHOOK_URL not configured")
	}

	// Prepare webhook payload
	webhookPayload := map[string]interface{}{
		"event_type": eventType,
		"timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
		"data":       payload,
	}

	payloadBytes, err := json.Marshal(webhookPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Send webhook
	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.supabase.ServiceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}

// Close closes the Supabase client connections
func (c *Client) Close() error {
	if c.redis != nil {
		return c.redis.Close()
	}
	return nil
}
