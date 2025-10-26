package engine

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/Hmbown/driftlock/api-server/internal/cbad"
    "github.com/Hmbown/driftlock/api-server/internal/models"
    "go.opentelemetry.io/otel"
)

// Engine processes telemetry data and detects anomalies
type Engine struct {
    detector *cbad.Detector
}

// New creates a new engine with the provided CBAD detector
func New(detector *cbad.Detector) *Engine {
    return &Engine{
        detector: detector,
    }
}

// Process consumes telemetry data and performs anomaly detection using CBAD
func (e *Engine) Process(ctx context.Context, payload []byte) error {
    tracer := otel.Tracer("driftlock/engine")
    ctx, span := tracer.Start(ctx, "engine.Process")
    defer span.End()

    // Validate payload
    if len(payload) == 0 {
        log.Printf("engine: received empty payload")
        return nil
    }

    // If detector is nil (simple version), just return without processing
    if e.detector == nil {
        log.Printf("engine: detector not initialized (simple mode), skipping anomaly detection")
        return nil
    }

    // Attempt to parse the payload to determine its type
    var parsed interface{}
    if err := json.Unmarshal(payload, &parsed); err != nil {
        // If it's not valid JSON, we'll still process it with CBAD
        log.Printf("engine: payload is not JSON, processing as raw data (bytes: %d)", len(payload))
        return e.processWithCBAD(ctx, payload)
    }

    // If it's a structured payload, we'll process it based on its content
    return e.processWithCBAD(ctx, payload)
}

// processWithCBAD runs the CBAD analysis on the provided data
func (e *Engine) processWithCBAD(ctx context.Context, data []byte) error {
    // For now, we'll treat all incoming data as a potential anomaly to be compared
    // against a baseline of recent historical data.
    // In a real implementation, this would involve:
    // 1. Maintaining sliding windows of historical data
    // 2. Using appropriate baseline data for comparison
    // 3. Handling different types of telemetry (logs, metrics, traces)
    
    if e.detector == nil {
        log.Printf("engine: CBAD detector not available, skipping anomaly detection")
        return nil
    }

    // For now, we'll use a simple approach where we compare against the most recent data
    // In practice, you'd want to maintain proper baseline and window data
    baseline := getBaselineData(data) // This would be historical baseline data
    window := data                   // Current window of data to check

    // Process with CBAD detector
    streamType := detectStreamType(data)
    if err := e.detector.ProcessTelemetry(ctx, baseline, window, streamType); err != nil {
        return fmt.Errorf("CBAD processing failed: %w", err)
    }

    return nil
}

// getBaselineData returns baseline data for comparison (in a real implementation, this would come from historical data)
func getBaselineData(currentData []byte) []byte {
    // This is a simplified approach - in reality, you'd want to maintain 
    // proper sliding windows of historical data for comparison
    if len(currentData) < 100 {
        return []byte(`{"timestamp": "2025-01-01T00:00:00Z", "message": "normal operation", "value": 42}`)
    }
    
    // Return a portion of the current data as a simple baseline
    // In a real system, this would be historical baseline data
    end := len(currentData)
    if end > 100 {
        end = 100
    }
    return currentData[:end]
}

// detectStreamType determines the type of telemetry stream from the data
func detectStreamType(data []byte) models.StreamType {
    // Try to detect stream type based on the content
    if len(data) == 0 {
        return models.StreamTypeLogs // Default to logs
    }

    // Look for common patterns in the data
    stringData := string(data)
    if contains(stringData, "metric") || contains(stringData, "gauge") || contains(stringData, "counter") {
        return models.StreamTypeMetrics
    } else if contains(stringData, "trace") || contains(stringData, "span") {
        return models.StreamTypeTraces
    } else if contains(stringData, "prompt") || contains(stringData, "completion") || contains(stringData, "llm") {
        return models.StreamTypeLLM
    }

    // Default to logs
    return models.StreamTypeLogs
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
    // Simple case-insensitive search
    sLower := lower(s)
    substrLower := lower(substr)
    
    for i := 0; i <= len(sLower)-len(substrLower); i++ {
        if sLower[i:i+len(substrLower)] == substrLower {
            return true
        }
    }
    return false
}

// lower converts a string to lowercase
func lower(s string) string {
    var result []byte
    for i := 0; i < len(s); i++ {
        c := s[i]
        if c >= 'A' && c <= 'Z' {
            c = c + ('a' - 'A')
        }
        result = append(result, c)
    }
    return string(result)
}
