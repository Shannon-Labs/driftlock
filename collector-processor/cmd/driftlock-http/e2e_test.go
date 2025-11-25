package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// testEnv holds the test environment
type testEnv struct {
	t       *testing.T
	server  *httptest.Server
	handler http.Handler
	store   *store
	cfg     config
	cleanup func()
}

// setupTestEnv creates a test environment with a real database
// Requires DATABASE_URL to be set or uses default test connection
func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	// Set dev mode for tests (bypasses license, allows Stripe signature bypass)
	os.Setenv("DRIFTLOCK_DEV_MODE", "true")

	// Load test config
	cfg := loadConfig()
	cfg.RateLimitRPS = 1000 // High rate limit for tests

	// Connect to test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL or DATABASE_URL not set, skipping E2E tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("Failed to ping test database: %v", err)
	}

	store := newStore(pool)
	if err := store.loadCache(ctx); err != nil {
		pool.Close()
		t.Fatalf("Failed to load cache: %v", err)
	}

	// Create minimal job queue for tests
	queue := newMemoryQueue(100)
	limiter := newTenantRateLimiter(cfg.DefaultRateLimit())

	// Use nil emailer - emails won't be sent, but we can get tokens from DB
	var emailer *emailService = nil
	tracker := newUsageTracker(store, emailer)

	handler := buildHTTPHandler(cfg, store, queue, limiter, emailer, tracker)
	server := httptest.NewServer(handler)

	return &testEnv{
		t:       t,
		server:  server,
		handler: handler,
		store:   store,
		cfg:     cfg,
		cleanup: func() {
			server.Close()
			pool.Close()
		},
	}
}

// Close cleans up the test environment
func (te *testEnv) Close() {
	if te.cleanup != nil {
		te.cleanup()
	}
}

// cleanupTestTenants removes test tenants by email pattern
func (te *testEnv) cleanupTestTenants(emailPattern string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete test data (cascades to related tables)
	_, err := te.store.pool.Exec(ctx, `
		DELETE FROM tenants WHERE email LIKE $1
	`, emailPattern)
	if err != nil {
		te.t.Logf("Warning: failed to cleanup test tenants: %v", err)
	}
}

// getVerificationToken retrieves the verification token from the database
func (te *testEnv) getVerificationToken(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var token string
	err := te.store.pool.QueryRow(ctx, `
		SELECT verification_token FROM tenants WHERE email = $1
	`, email).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("failed to get verification token: %w", err)
	}
	return token, nil
}

// getTenantStatus retrieves tenant status from database
func (te *testEnv) getTenantStatus(email string) (status string, verified bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var verifiedAt *time.Time
	err = te.store.pool.QueryRow(ctx, `
		SELECT status, email_verified_at FROM tenants WHERE email = $1
	`, email).Scan(&status, &verifiedAt)
	if err != nil {
		return "", false, err
	}
	return status, verifiedAt != nil, nil
}

// httpClient helpers

func (te *testEnv) doRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	te.t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			te.t.Fatalf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, te.server.URL+path, bodyReader)
	if err != nil {
		te.t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		te.t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		te.t.Fatalf("Failed to read response body: %v", err)
	}

	return resp, respBody
}

func (te *testEnv) post(path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	return te.doRequest(http.MethodPost, path, body, headers)
}

func (te *testEnv) get(path string, headers map[string]string) (*http.Response, []byte) {
	return te.doRequest(http.MethodGet, path, nil, headers)
}

// assertStatus checks response status code
func (te *testEnv) assertStatus(resp *http.Response, expected int) {
	te.t.Helper()
	if resp.StatusCode != expected {
		te.t.Errorf("Expected status %d, got %d", expected, resp.StatusCode)
	}
}

// assertJSONField checks a field in JSON response
func (te *testEnv) assertJSONField(body []byte, field string, expected interface{}) {
	te.t.Helper()
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		te.t.Fatalf("Failed to parse JSON: %v", err)
	}

	parts := strings.Split(field, ".")
	var current interface{} = data
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			te.t.Errorf("Field %s not found in response", field)
			return
		}
	}

	if current != expected {
		te.t.Errorf("Expected %s = %v, got %v", field, expected, current)
	}
}

// assertJSONContains checks if a field contains a substring
func (te *testEnv) assertJSONContains(body []byte, field, substring string) {
	te.t.Helper()
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		te.t.Fatalf("Failed to parse JSON: %v", err)
	}

	parts := strings.Split(field, ".")
	var current interface{} = data
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			te.t.Errorf("Field %s not found in response", field)
			return
		}
	}

	str, ok := current.(string)
	if !ok {
		te.t.Errorf("Field %s is not a string", field)
		return
	}

	if !strings.Contains(str, substring) {
		te.t.Errorf("Expected %s to contain %q, got %q", field, substring, str)
	}
}

// generateTestEmail creates a unique test email address
func generateTestEmail(prefix string) string {
	return fmt.Sprintf("test.%s.%d@e2e-test.driftlock.net", prefix, time.Now().UnixNano())
}
