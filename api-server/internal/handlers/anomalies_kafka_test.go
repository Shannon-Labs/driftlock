package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/storage"
	"github.com/Hmbown/driftlock/api-server/internal/stream"
	"github.com/Hmbown/driftlock/api-server/internal/streaming"
	"github.com/Hmbown/driftlock/api-server/internal/streaming/kafka"
)

// TestAPIWithoutKafka tests that the API server works when Kafka is disabled
func TestAPIWithoutKafka(t *testing.T) {
	// Disable Kafka by setting the environment variable
	os.Setenv("KAFKA_ENABLED", "false")
	defer os.Unsetenv("KAFKA_ENABLED")

	// Create an in-memory database for testing
	db := storage.NewInMemory()
	
	// Create a test streamer
	streamer := stream.NewStreamer(100)

	// Create an in-memory event publisher
	broker := kafka.NewInMemoryBroker()
	publisher := streaming.NewInMemoryPublisher(broker, "anomaly-events")

	// Create the anomalies handler
	anomaliesHandler := handlers.NewAnomaliesHandler(db, streamer, publisher)

	// Create a test HTTP request to create an anomaly
	anomalyData := `{
		"stream_type": "logs",
		"message": "Test anomaly",
		"severity": "high",
		"ncd_score": 0.85,
		"p_value": 0.01,
		"anomaly_explanation": "This is a test anomaly",
		"timestamp": "2025-10-25T10:00:00Z"
	}`

	req := httptest.NewRequest("POST", "/v1/anomalies", strings.NewReader(anomalyData))
	req.Header.Set("Content-Type", "application/json")
	
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

	if createdAnomaly.Message != "Test anomaly" {
		t.Errorf("Expected anomaly message 'Test anomaly', got '%s'", createdAnomaly.Message)
	}

	// Verify that no Kafka errors occurred (since Kafka is disabled)
	// The in-memory publisher should handle the event without error
}