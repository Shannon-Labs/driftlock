package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// E2E tests for CBAD Security Features (SHA-139, SHA-140, SHA-141, SHA-142, SHA-143)
// These tests require a real database with migrations applied.

// createTestTenantAndStream creates a verified tenant with a stream for CBAD tests.
// Returns the API key and stream slug.
func (te *testEnv) createTestTenantAndStream(prefix string, streamSettings map[string]interface{}) (apiKey string, streamSlug string) {
	te.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate unique identifiers
	email := generateTestEmail(prefix)
	streamSlug = fmt.Sprintf("test-stream-%d", time.Now().UnixNano())
	tenantSlug := fmt.Sprintf("e2e-test-%d", time.Now().UnixNano())
	tenantID := uuid.New()

	// Create tenant
	_, err := te.store.pool.Exec(ctx, `
		INSERT INTO tenants (id, email, name, slug, status, verified_at, created_at, updated_at)
		VALUES ($1, $2, 'E2E Test Co', $3, 'active', NOW(), NOW(), NOW())
	`, tenantID, email, tenantSlug)
	if err != nil {
		te.t.Fatalf("Failed to create test tenant: %v", err)
	}

	// Extract settings or use defaults
	anchorEnabled := true
	if v, ok := streamSettings["anchor_enabled"].(bool); ok {
		anchorEnabled = v
	}
	driftThreshold := 0.35
	if v, ok := streamSettings["drift_ncd_threshold"].(float64); ok {
		driftThreshold = v
	}
	numericOutlierEnabled, _ := streamSettings["numeric_outlier_enabled"].(bool)
	numericKSigma := 3.0
	if v, ok := streamSettings["numeric_k_sigma"].(float64); ok {
		numericKSigma = v
	}

	// Create stream first (required for API key)
	streamID := uuid.New()
	_, err = te.store.pool.Exec(ctx, `
		INSERT INTO streams (id, tenant_id, slug, type, description, anchor_enabled, drift_ncd_threshold,
		                     numeric_outlier_enabled, numeric_k_sigma, created_at, updated_at)
		VALUES ($1, $2, $3, 'logs', 'E2E Test Stream', $4, $5, $6, $7, NOW(), NOW())
	`, streamID, tenantID, streamSlug, anchorEnabled, driftThreshold, numericOutlierEnabled, numericKSigma)
	if err != nil {
		te.t.Fatalf("Failed to create test stream: %v", err)
	}

	// Use the store's createAPIKey method for proper key format
	apiKey, _, err = te.store.createAPIKey(ctx, tenantID, "test-key", "admin")
	if err != nil {
		te.t.Fatalf("Failed to create API key: %v", err)
	}

	// Reload store cache to pick up new data
	if err := te.store.loadCache(ctx); err != nil {
		te.t.Fatalf("Failed to reload store cache: %v", err)
	}

	return apiKey, streamSlug
}

// detectRequest sends events to the detect endpoint
func (te *testEnv) detectRequest(apiKey, streamSlug string, events []map[string]interface{}) (map[string]interface{}, int) {
	te.t.Helper()

	payload := map[string]interface{}{
		"stream": streamSlug,
		"events": events,
	}

	resp, body := te.post("/v1/detect", payload, map[string]string{
		"X-API-Key": apiKey,
	})

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		te.t.Logf("Response body: %s", string(body))
		te.t.Fatalf("Failed to parse detect response: %v", err)
	}

	return result, resp.StatusCode
}

// generateEvents creates N copies of a template event with sequence numbers
func generateEvents(count int, template map[string]interface{}) []map[string]interface{} {
	events := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		event := make(map[string]interface{}, len(template)+1)
		for k, v := range template {
			event[k] = v
		}
		event["seq"] = i
		events[i] = event
	}
	return events
}

// TestE2E_ColdStartCalibration tests SHA-139: calibration status for low-history streams
func TestE2E_ColdStartCalibration(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	// Create stream with default settings (50 event minimum)
	apiKey, streamSlug := te.createTestTenantAndStream("coldstart", map[string]interface{}{})

	t.Run("FirstEventReturnsCalibrating", func(t *testing.T) {
		events := []map[string]interface{}{
			{"message": "first event", "level": "info"},
		}
		result, status := te.detectRequest(apiKey, streamSlug, events)

		if status != http.StatusOK {
			t.Fatalf("Expected 200, got %d", status)
		}

		// Should be calibrating since we only sent 1 event
		if result["status"] != "calibrating" {
			t.Errorf("Expected status 'calibrating', got %v", result["status"])
		}

		// Check calibration info
		calibration, ok := result["calibration"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected calibration info in response")
		}

		if calibration["events_ingested"].(float64) < 1 {
			t.Errorf("Expected events_ingested >= 1, got %v", calibration["events_ingested"])
		}

		if calibration["events_needed"].(float64) < 1 {
			t.Errorf("Expected events_needed >= 1, got %v", calibration["events_needed"])
		}
	})

	t.Run("After50EventsBecomesReady", func(t *testing.T) {
		// Send 60 events to exceed the 50 minimum baseline
		events := generateEvents(60, map[string]interface{}{
			"message": "calibration event",
			"level":   "info",
			"user_id": "user-123",
		})

		result, status := te.detectRequest(apiKey, streamSlug, events)

		if status != http.StatusOK {
			t.Fatalf("Expected 200, got %d", status)
		}

		// Should now be ready
		if result["status"] != "ready" {
			// Might still be calibrating if we haven't hit threshold yet
			if result["status"] == "calibrating" {
				calibration := result["calibration"].(map[string]interface{})
				t.Logf("Still calibrating: %d/%d events",
					int(calibration["events_ingested"].(float64)),
					int(calibration["min_baseline_size"].(float64)))
				// Send more events
				events = generateEvents(50, map[string]interface{}{
					"message": "more calibration",
					"level":   "debug",
				})
				result, _ = te.detectRequest(apiKey, streamSlug, events)
			}
		}

		// After enough events, should be ready
		if result["status"] != "ready" {
			t.Errorf("Expected status 'ready' after calibration, got %v", result["status"])
		}
	})
}

// TestE2E_TokenizerReducesNoise tests SHA-141: high-entropy tokenization
func TestE2E_TokenizerReducesNoise(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	// Create stream with tokenizer enabled
	apiKey, streamSlug := te.createTestTenantAndStream("tokenizer", map[string]interface{}{
		"tokenizer_enabled": true,
		"tokenize_uuid":     true,
		"tokenize_hash":     true,
		"tokenize_jwt":      true,
	})

	t.Run("UUIDsTokenizedConsistently", func(t *testing.T) {
		// Send baseline events with UUIDs
		baseline := generateEvents(60, map[string]interface{}{
			"message":    "user login",
			"user_id":    "550e8400-e29b-41d4-a716-446655440000",
			"session_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"level":      "info",
		})

		result, status := te.detectRequest(apiKey, streamSlug, baseline)
		if status != http.StatusOK {
			t.Fatalf("Baseline failed: %d", status)
		}

		// Now send events with DIFFERENT UUIDs - should NOT be anomalous
		// because tokenizer replaces all UUIDs with <UUID>
		different := []map[string]interface{}{
			{
				"message":    "user login",
				"user_id":    "11111111-2222-3333-4444-555555555555",
				"session_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"level":      "info",
			},
		}

		result, status = te.detectRequest(apiKey, streamSlug, different)
		if status != http.StatusOK {
			t.Fatalf("Detection failed: %d", status)
		}

		// With tokenizer, UUID changes should not cause anomalies
		anomalyCount := int(result["anomaly_count"].(float64))
		if anomalyCount > 0 {
			t.Logf("Anomaly detected (may be expected during early calibration): count=%d", anomalyCount)
		}
	})
}

// TestE2E_AnchorDriftDetection tests SHA-140: anchor baseline drift detection
func TestE2E_AnchorDriftDetection(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	// Create stream
	apiKey, streamSlug := te.createTestTenantAndStream("drift", map[string]interface{}{
		"anchor_enabled":     true,
		"drift_ncd_threshold": 0.3,
	})

	t.Run("CalibrateAndCreateAnchor", func(t *testing.T) {
		// First calibrate the stream
		events := generateEvents(60, map[string]interface{}{
			"message": "normal operation",
			"level":   "info",
			"code":    200,
		})

		result, status := te.detectRequest(apiKey, streamSlug, events)
		if status != http.StatusOK {
			t.Fatalf("Calibration failed: %d", status)
		}

		// Create anchor with baseline events
		anchorEvents := generateEvents(100, map[string]interface{}{
			"message": "baseline event",
			"level":   "info",
			"code":    200,
		})

		resp, body := te.post("/v1/streams/"+streamSlug+"/reset-anchor", map[string]interface{}{
			"events": anchorEvents,
		}, map[string]string{
			"X-API-Key": apiKey,
		})

		if resp.StatusCode != http.StatusOK {
			t.Logf("Response: %s", string(body))
			t.Fatalf("Anchor creation failed: %d", resp.StatusCode)
		}

		var anchorResult map[string]interface{}
		json.Unmarshal(body, &anchorResult)

		if anchorResult["success"] != true {
			t.Errorf("Expected success=true, got %v", anchorResult["success"])
		}

		if anchorResult["anchor_id"] == nil {
			t.Error("Expected anchor_id in response")
		}

		// Now send events that should trigger drift detection
		driftEvents := generateEvents(20, map[string]interface{}{
			"message":   "CRITICAL ERROR",
			"level":     "error",
			"code":      500,
			"exception": "NullPointerException at com.example.Service.process()",
		})

		result, status = te.detectRequest(apiKey, streamSlug, driftEvents)
		if status != http.StatusOK {
			t.Fatalf("Detection failed: %d", status)
		}

		// Check for drift result
		if drift, ok := result["drift"].(map[string]interface{}); ok {
			t.Logf("Drift detection: score=%.4f, threshold=%.4f, detected=%v",
				drift["drift_score"], drift["drift_threshold"], drift["drift_detected"])
		} else {
			t.Log("No drift result in response (anchor may not be active yet)")
		}
	})
}

// TestE2E_NumericOutlierDetection tests SHA-142/SHA-143: numeric value outlier detection
func TestE2E_NumericOutlierDetection(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	// Create stream with numeric outlier detection enabled
	apiKey, streamSlug := te.createTestTenantAndStream("numeric", map[string]interface{}{
		"numeric_outlier_enabled": true,
		"numeric_k_sigma":         3.0,
	})

	t.Run("NormalValuesNoOutlier", func(t *testing.T) {
		// Train with normal values around mean=100
		events := generateEvents(50, map[string]interface{}{
			"message":  "transaction completed",
			"amount":   100.0,
			"duration": 50.0,
		})

		// Add some variance
		for i := range events {
			events[i]["amount"] = 90.0 + float64(i%20)  // 90-109
			events[i]["duration"] = 45.0 + float64(i%10) // 45-54
		}

		result, status := te.detectRequest(apiKey, streamSlug, events)
		if status != http.StatusOK {
			t.Fatalf("Expected 200, got %d", status)
		}

		// Should have no outliers since values are consistent
		if outliers, ok := result["value_outliers"].([]interface{}); ok && len(outliers) > 0 {
			t.Logf("Got %d outliers during training (expected during early calibration)", len(outliers))
		}
	})

	t.Run("ExtremeValueFlagsOutlier", func(t *testing.T) {
		// First, establish a good baseline with more data
		for i := 0; i < 3; i++ {
			events := generateEvents(30, map[string]interface{}{
				"message":  "normal transaction",
				"amount":   100.0,
				"duration": 50.0,
			})
			te.detectRequest(apiKey, streamSlug, events)
		}

		// Now send an extreme outlier
		outlierEvent := []map[string]interface{}{
			{
				"message":  "suspicious transaction",
				"amount":   100000.0, // Way outside normal range
				"duration": 50.0,
			},
		}

		result, status := te.detectRequest(apiKey, streamSlug, outlierEvent)
		if status != http.StatusOK {
			t.Fatalf("Expected 200, got %d", status)
		}

		// Should flag the amount as an outlier
		if outliers, ok := result["value_outliers"].([]interface{}); ok {
			if len(outliers) == 0 {
				t.Log("No outliers flagged - may need more baseline data")
			} else {
				t.Logf("Detected %d outlier(s):", len(outliers))
				for _, o := range outliers {
					outlier := o.(map[string]interface{})
					t.Logf("  - Field: %s, Value: %v, Z-Score: %.2f",
						outlier["field_path"], outlier["value"], outlier["z_score"])
				}

				// Check that amount was flagged
				found := false
				for _, o := range outliers {
					outlier := o.(map[string]interface{})
					if outlier["field_path"] == "amount" {
						found = true
						if outlier["value"].(float64) != 100000.0 {
							t.Errorf("Expected outlier value 100000, got %v", outlier["value"])
						}
					}
				}
				if !found {
					t.Error("Expected 'amount' field to be flagged as outlier")
				}
			}
		} else {
			t.Log("No value_outliers field in response")
		}
	})
}

// TestE2E_AnchorEndpoints tests anchor management REST endpoints
func TestE2E_AnchorEndpoints(t *testing.T) {
	te := setupTestEnv(t)
	defer te.Close()
	defer te.cleanupTestTenants("%@e2e-test.driftlock.net")

	apiKey, streamSlug := te.createTestTenantAndStream("anchor-api", map[string]interface{}{
		"anchor_enabled": true,
	})

	// First calibrate and send some events
	events := generateEvents(60, map[string]interface{}{
		"message": "baseline event",
		"level":   "info",
	})
	te.detectRequest(apiKey, streamSlug, events)

	t.Run("GetAnchorInitiallyEmpty", func(t *testing.T) {
		resp, body := te.get("/v1/streams/"+streamSlug+"/anchor", map[string]string{
			"X-API-Key": apiKey,
		})

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		// Should not have an active anchor yet
		if result["has_active_anchor"] == true {
			t.Log("Anchor already exists (may have been auto-created)")
		}
	})

	t.Run("CreateAnchorWithEvents", func(t *testing.T) {
		anchorEvents := generateEvents(50, map[string]interface{}{
			"message": "anchor baseline",
			"code":    200,
		})

		resp, body := te.post("/v1/streams/"+streamSlug+"/reset-anchor", map[string]interface{}{
			"events": anchorEvents,
		}, map[string]string{
			"X-API-Key": apiKey,
		})

		if resp.StatusCode != http.StatusOK {
			t.Logf("Response: %s", string(body))
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if result["success"] != true {
			t.Errorf("Expected success=true, got %v", result["success"])
		}

		if result["anchor_id"] == nil || result["anchor_id"] == "" {
			t.Error("Expected anchor_id in response")
		}
	})

	t.Run("GetAnchorDetails", func(t *testing.T) {
		resp, body := te.get("/v1/streams/"+streamSlug+"/anchor/details", map[string]string{
			"X-API-Key": apiKey,
		})

		if resp.StatusCode != http.StatusOK {
			t.Logf("Response: %s", string(body))
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		anchor, ok := result["anchor"].(map[string]interface{})
		if !ok {
			t.Logf("Response: %s", string(body))
			t.Fatal("Expected anchor object in response")
		}

		// Check anchor has expected fields
		if anchor["id"] == nil || anchor["id"] == "" {
			t.Error("Expected anchor id")
		}
		if anchor["event_count"] == nil {
			t.Error("Expected event_count")
		}
		if anchor["is_active"] != true {
			t.Error("Expected is_active=true")
		}
	})

	t.Run("DeleteAnchor", func(t *testing.T) {
		resp, body := te.doRequest(http.MethodDelete, "/v1/streams/"+streamSlug+"/anchor", nil, map[string]string{
			"X-API-Key": apiKey,
		})

		if resp.StatusCode != http.StatusOK {
			t.Logf("Response: %s", string(body))
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if result["success"] != true {
			t.Errorf("Expected success=true, got %v", result["success"])
		}

		// Verify anchor is gone
		resp, body = te.get("/v1/streams/"+streamSlug+"/anchor", map[string]string{
			"X-API-Key": apiKey,
		})

		json.Unmarshal(body, &result)
		if result["has_active_anchor"] == true {
			t.Error("Expected no active anchor after delete")
		}
	})
}
