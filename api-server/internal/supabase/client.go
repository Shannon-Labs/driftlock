package supabase

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
)

// Client is a minimal Supabase HTTP client for REST and Edge Functions.
type Client struct {
    BaseURL        string
    AnonKey        string
    ServiceRoleKey string
    HTTPClient     *http.Client

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

// NewClient creates a new minimal Supabase client
func NewClient(cfg Config) (*Client, error) {
    var redisClient *redis.Client
    if cfg.RedisAddr != "" {
        redisClient = redis.NewClient(&redis.Options{
            Addr:     cfg.RedisAddr,
            Password: cfg.RedisPassword,
            DB:       cfg.RedisDB,
        })
    }

    return &Client{
        BaseURL:        cfg.BaseURL,
        AnonKey:        cfg.AnonKey,
        ServiceRoleKey: cfg.ServiceRoleKey,
        HTTPClient:     &http.Client{Timeout: 10 * time.Second},
        redis:          redisClient,
        projectID:      cfg.ProjectID,
    }, nil
}

// HealthCheck verifies Supabase edge function health if available.
func (c *Client) HealthCheck(ctx context.Context) error {
    if c.BaseURL == "" {
        return fmt.Errorf("supabase base URL not configured")
    }
    u := fmt.Sprintf("%s/functions/v1/health", c.BaseURL)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("supabase health returned %d", resp.StatusCode)
    }
    return nil
}

// CreateAnomaly inserts an anomaly into Supabase via REST
func (c *Client) CreateAnomaly(ctx context.Context, anomaly map[string]interface{}) error {
    if c.BaseURL == "" || c.ServiceRoleKey == "" {
        return fmt.Errorf("supabase base URL or service role key not configured")
    }
    endpoint := fmt.Sprintf("%s/rest/v1/anomalies", c.BaseURL)
    body, err := json.Marshal(anomaly)
    if err != nil {
        return fmt.Errorf("marshal anomaly: %w", err)
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Prefer", "return=representation")
    req.Header.Set("apikey", c.ServiceRoleKey)
    req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("supabase error: %d", resp.StatusCode)
    }

    if c.redis != nil {
        c.redis.Del(ctx, fmt.Sprintf("anomalies:list:%s", c.projectID))
    }
    return nil
}

// GetAnomalies queries anomalies using REST filters
func (c *Client) GetAnomalies(ctx context.Context, filter map[string]interface{}) ([]map[string]interface{}, error) {
    if c.BaseURL == "" || c.AnonKey == "" {
        return nil, fmt.Errorf("supabase base URL or anon key not configured")
    }
    u, _ := url.Parse(fmt.Sprintf("%s/rest/v1/anomalies", c.BaseURL))
    q := u.Query()
    q.Set("select", "*")
    for k, v := range filter {
        q.Set(k, fmt.Sprintf("eq.%v", v))
    }
    u.RawQuery = q.Encode()
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
    req.Header.Set("apikey", c.AnonKey)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("supabase error: %d", resp.StatusCode)
    }
    var out []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode: %w", err)
    }
    return out, nil
}

// UpdateAnomalyStatus patches anomaly status via REST
func (c *Client) UpdateAnomalyStatus(ctx context.Context, id string, status string) error {
    if c.BaseURL == "" || c.ServiceRoleKey == "" {
        return fmt.Errorf("supabase base URL or service role key not configured")
    }
    endpoint := fmt.Sprintf("%s/rest/v1/anomalies?id=eq.%s", c.BaseURL, url.QueryEscape(id))
    body, _ := json.Marshal(map[string]interface{}{"status": status})
    req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Prefer", "return=representation")
    req.Header.Set("apikey", c.ServiceRoleKey)
    req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("supabase error: %d", resp.StatusCode)
    }
    if c.redis != nil {
        c.redis.Del(ctx, fmt.Sprintf("anomalies:list:%s", c.projectID))
    }
    return nil
}

// CreateUsageRecord inserts a usage record via REST
func (c *Client) CreateUsageRecord(ctx context.Context, usage map[string]interface{}) error {
    if c.BaseURL == "" || c.ServiceRoleKey == "" {
        return fmt.Errorf("supabase base URL or service role key not configured")
    }
    endpoint := fmt.Sprintf("%s/rest/v1/usage_records", c.BaseURL)
    body, _ := json.Marshal(usage)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Prefer", "return=representation")
    req.Header.Set("apikey", c.ServiceRoleKey)
    req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("supabase error: %d", resp.StatusCode)
    }
    return nil
}

// GetOrganization retrieves organization details via REST
func (c *Client) GetOrganization(ctx context.Context, orgID string) (map[string]interface{}, error) {
    if c.BaseURL == "" || c.AnonKey == "" {
        return nil, fmt.Errorf("supabase base URL or anon key not configured")
    }
    u := fmt.Sprintf("%s/rest/v1/organizations?id=eq.%s&select=*", c.BaseURL, url.QueryEscape(orgID))
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
    req.Header.Set("apikey", c.AnonKey)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("supabase error: %d", resp.StatusCode)
    }
    var results []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return nil, fmt.Errorf("decode: %w", err)
    }
    if len(results) == 0 {
        return nil, fmt.Errorf("organization not found: %s", orgID)
    }
    return results[0], nil
}

// NotifyWebhook sends a webhook notification to Supabase Edge Function
func (c *Client) NotifyWebhook(ctx context.Context, eventType string, payload interface{}) error {
    webhookURL := os.Getenv("SUPABASE_WEBHOOK_URL")
    if webhookURL == "" {
        return fmt.Errorf("SUPABASE_WEBHOOK_URL not configured")
    }
    webhookPayload := map[string]interface{}{
        "event_type": eventType,
        "timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
        "data":       payload,
    }
    payloadBytes, err := json.Marshal(webhookPayload)
    if err != nil {
        return fmt.Errorf("failed to marshal webhook payload: %w", err)
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payloadBytes))
    if err != nil {
        return fmt.Errorf("failed to create webhook request: %w", err)
    }
    // Authorize with service role when available
    if c.ServiceRoleKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
    }
    return nil
}

// MeterUsage calls the Supabase Edge Function to meter usage for billing
func (c *Client) MeterUsage(ctx context.Context, organizationID string, hasAnomaly bool, count int) error {
    if c.BaseURL == "" || c.ServiceRoleKey == "" {
        return fmt.Errorf("supabase base URL or service role key not configured")
    }
    endpoint := fmt.Sprintf("%s/functions/v1/meter-usage", c.BaseURL)
    payload := map[string]interface{}{
        "organization_id": organizationID,
        "anomaly":        hasAnomaly,
        "count":          count,
    }
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("meter-usage request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return fmt.Errorf("meter-usage returned status: %d", resp.StatusCode)
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
