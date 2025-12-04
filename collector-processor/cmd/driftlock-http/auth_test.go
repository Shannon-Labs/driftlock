package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestRevokedKeyRejected validates that revoked API keys are properly rejected
// This is an E2E test that requires DATABASE_URL to be set
func TestRevokedKeyRejected(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()

	// Clean up any existing test data
	te.cleanupTestTenants("%revoke-test%")

	// Step 1: Sign up a new tenant
	email := generateTestEmail("revoke-test")
	resp, _ := te.post("/v1/onboard/signup", map[string]interface{}{
		"email":        email,
		"company_name": "Revoke Test Corp",
		"plan":         "trial",
	}, nil)
	te.assertStatus(resp, http.StatusOK)

	// Step 2: Get verification token and verify
	token, err := te.getVerificationToken(email)
	if err != nil {
		t.Fatalf("Failed to get verification token: %v", err)
	}

	resp, body := te.get("/v1/onboard/verify?token="+token, nil)
	te.assertStatus(resp, http.StatusOK)
	te.assertJSONField(body, "success", true)

	// Step 3: Get the API key from the response
	var verifyResp struct {
		APIKey string `json:"api_key"`
	}
	if err := unmarshalJSON(body, &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}
	if verifyResp.APIKey == "" {
		t.Fatal("Expected API key in verify response, got empty string")
	}
	apiKey := verifyResp.APIKey

	// Step 4: Make authenticated request - should succeed
	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
	}
	resp, _ = te.post("/v1/detect", map[string]interface{}{
		"stream_id": "default",
		"events": []map[string]interface{}{
			{"value": 42, "timestamp": time.Now().Format(time.RFC3339)},
		},
	}, headers)
	te.assertStatus(resp, http.StatusOK)
	t.Log("Authenticated request succeeded with valid API key")

	// Step 5: Revoke the API key directly in database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = te.store.pool.Exec(ctx, `
		UPDATE api_keys SET revoked = true, revoked_at = NOW()
		WHERE key_hash = encode(sha256($1::bytea), 'hex')
	`, apiKey)
	if err != nil {
		t.Fatalf("Failed to revoke API key: %v", err)
	}

	// Force cache refresh to pick up revocation
	if err := te.store.loadCache(ctx); err != nil {
		t.Fatalf("Failed to refresh cache: %v", err)
	}

	// Step 6: Make authenticated request with revoked key - should fail
	resp, body = te.post("/v1/detect", map[string]interface{}{
		"stream_id": "default",
		"events": []map[string]interface{}{
			{"value": 43, "timestamp": time.Now().Format(time.RFC3339)},
		},
	}, headers)
	te.assertStatus(resp, http.StatusUnauthorized)
	t.Log("Revoked key correctly rejected with 401")

	// Cleanup
	te.cleanupTestTenants("%revoke-test%")
}

// TestExpiredKeyRejected validates that expired API keys are rejected
func TestExpiredKeyRejected(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()

	// Clean up any existing test data
	te.cleanupTestTenants("%expire-test%")

	// Step 1: Sign up and verify a new tenant
	email := generateTestEmail("expire-test")
	resp, _ := te.post("/v1/onboard/signup", map[string]interface{}{
		"email":        email,
		"company_name": "Expire Test Corp",
		"plan":         "trial",
	}, nil)
	te.assertStatus(resp, http.StatusOK)

	token, err := te.getVerificationToken(email)
	if err != nil {
		t.Fatalf("Failed to get verification token: %v", err)
	}

	resp, body := te.get("/v1/onboard/verify?token="+token, nil)
	te.assertStatus(resp, http.StatusOK)

	var verifyResp struct {
		APIKey string `json:"api_key"`
	}
	if err := unmarshalJSON(body, &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}
	apiKey := verifyResp.APIKey

	// Step 2: Set expiration to past date
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = te.store.pool.Exec(ctx, `
		UPDATE api_keys SET expires_at = NOW() - INTERVAL '1 hour'
		WHERE key_hash = encode(sha256($1::bytea), 'hex')
	`, apiKey)
	if err != nil {
		t.Fatalf("Failed to set key expiration: %v", err)
	}

	// Force cache refresh
	if err := te.store.loadCache(ctx); err != nil {
		t.Fatalf("Failed to refresh cache: %v", err)
	}

	// Step 3: Make authenticated request - should fail
	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
	}
	resp, _ = te.post("/v1/detect", map[string]interface{}{
		"stream_id": "default",
		"events": []map[string]interface{}{
			{"value": 42, "timestamp": time.Now().Format(time.RFC3339)},
		},
	}, headers)
	te.assertStatus(resp, http.StatusUnauthorized)
	t.Log("Expired key correctly rejected with 401")

	// Cleanup
	te.cleanupTestTenants("%expire-test%")
}

// unmarshalJSON is a helper for parsing JSON responses
func unmarshalJSON(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}
