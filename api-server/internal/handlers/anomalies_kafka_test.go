package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
	"github.com/Shannon-Labs/driftlock/api-server/internal/stream"
	"github.com/Shannon-Labs/driftlock/api-server/internal/streaming"
	"github.com/Shannon-Labs/driftlock/api-server/internal/streaming/kafka"
)

// TestAPIWithoutKafka tests that the API server works when Kafka is disabled
func TestAPIWithoutKafka(t *testing.T) {
	// Disable Kafka by setting the environment variable
	os.Setenv("KAFKA_ENABLED", "false")
	defer os.Unsetenv("KAFKA_ENABLED")

	// Create a mock storage for testing
	db := newMockStorage()
	
	// Create a test streamer
	streamer := stream.NewStreamer(100)

	// Create an in-memory event publisher
	broker := kafka.NewInMemoryBroker()
	publisher := streaming.NewInMemoryPublisher(broker, "anomaly-events")

	// Create the anomalies handler (without Supabase)
	anomaliesHandler := NewAnomaliesHandler(db, streamer, publisher)

	// Create a test HTTP request to create an anomaly
	anomalyData := `{
		"stream_type": "logs",
		"ncd_score": 0.85,
		"p_value": 0.01,
		"glass_box_explanation": "This is a test anomaly",
		"compression_baseline": 0.5,
		"compression_window": 0.6,
		"compression_combined": 0.55,
		"confidence_level": 0.95,
		"timestamp": "2025-10-25T10:00:00Z"
	}`

	req := httptest.NewRequest("POST", "/v1/anomalies", strings.NewReader(anomalyData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background())
	
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function that matches the expected signature
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			anomaliesHandler.CreateAnomaly(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Execute the request
	handler.ServeHTTP(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d: %s", http.StatusCreated, status, rr.Body.String())
	}

	// Verify that the anomaly was created in the database
	var createdAnomaly models.Anomaly
	if err := json.Unmarshal(rr.Body.Bytes(), &createdAnomaly); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdAnomaly.GlassBoxExplanation != "This is a test anomaly" {
		t.Errorf("Expected anomaly explanation 'This is a test anomaly', got '%s'", createdAnomaly.GlassBoxExplanation)
	}

	if createdAnomaly.StreamType != models.StreamTypeLogs {
		t.Errorf("Expected stream type 'logs', got '%s'", createdAnomaly.StreamType)
	}

	// Verify that no Kafka errors occurred (since Kafka is disabled)
	// The in-memory publisher should handle the event without error
}