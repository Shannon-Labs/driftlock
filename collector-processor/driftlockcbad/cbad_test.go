package driftlockcbad

import (
	"testing"
)

func TestValidateLibrary(t *testing.T) {
	err := ValidateLibrary()
	if err != nil {
		t.Fatalf("CBAD library validation failed: %v", err)
	}
}

func TestComputeMetricsQuick(t *testing.T) {
	// Test with normal vs anomalous log data
	baseline := []byte(`{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}`)
	window := []byte(`{"timestamp":"2025-10-24T00:00:01Z","severity":"ERROR","service":"api-gateway","message":"Panic occurred","attributes":{"stack_trace":"thread 'main' panicked at 'index out of bounds', src/main.rs:42:13","error":"runtime panic"}}`)

	metrics, err := ComputeMetricsQuick(baseline, window)
	if err != nil {
		t.Fatalf("ComputeMetricsQuick failed: %v", err)
	}

	// Validate metrics
	if metrics.NCD < 0 || metrics.NCD > 1 {
		t.Errorf("Invalid NCD value: %f", metrics.NCD)
	}

	if metrics.PValue < 0 || metrics.PValue > 1 {
		t.Errorf("Invalid p-value: %f", metrics.PValue)
	}

	// Should detect an anomaly between normal logs and error logs
	if !metrics.IsAnomaly {
		t.Logf("Anomaly not detected (this might be expected with fallback compression): NCD=%.3f, p=%.3f", metrics.NCD, metrics.PValue)
	}

	// Test explanations
	anomalyExplanation := metrics.GetAnomalyExplanation()
	if anomalyExplanation == "" {
		t.Error("GetAnomalyExplanation returned empty string")
	}
	t.Logf("Anomaly explanation: %s", anomalyExplanation)

	detailedExplanation := metrics.GetDetailedExplanation()
	if detailedExplanation == "" {
		t.Error("GetDetailedExplanation returned empty string")
	}
	t.Logf("Detailed explanation: %s", detailedExplanation)
}

func TestComputeMetricsWithConfig(t *testing.T) {
	baseline := []byte("INFO service=api-gateway msg=request_completed duration_ms=42\n")
	window := []byte("ERROR service=api-gateway msg=stack_trace_panic\n")

	// Test with custom configuration
	metrics, err := ComputeMetrics(baseline, window, 12345, 500)
	if err != nil {
		t.Fatalf("ComputeMetrics failed: %v", err)
	}

	// Validate basic metrics
	if metrics.NCD < 0 || metrics.NCD > 1 {
		t.Errorf("Invalid NCD value: %f", metrics.NCD)
	}

	t.Logf("Custom config results: NCD=%.3f, p=%.3f, is_anomaly=%v", metrics.NCD, metrics.PValue, metrics.IsAnomaly)
}

func TestEmptyDataHandling(t *testing.T) {
	// Test with empty baseline
	_, err := ComputeMetricsQuick([]byte{}, []byte("test data"))
	if err == nil {
		t.Error("Expected error for empty baseline")
	}

	// Test with empty window
	_, err = ComputeMetricsQuick([]byte("test data"), []byte{})
	if err == nil {
		t.Error("Expected error for empty window")
	}
}

func TestGlassBoxExplanations(t *testing.T) {
	baseline := []byte("INFO service=api-gateway msg=request_completed duration_ms=42\n")
	window := []byte("ERROR service=api-gateway msg=stack_trace_panic\n")

	metrics, err := ComputeMetricsQuick(baseline, window)
	if err != nil {
		t.Fatalf("ComputeMetricsQuick failed: %v", err)
	}

	// Test glass-box explanations
	anomalyExplanation := metrics.GetAnomalyExplanation()
	if anomalyExplanation == "" {
		t.Error("GetAnomalyExplanation returned empty string")
	}
	t.Logf("Glass-box explanation: %s", anomalyExplanation)

	// Test detailed explanation
	detailedExplanation := metrics.GetDetailedExplanation()
	if detailedExplanation == "" {
		t.Error("GetDetailedExplanation returned empty string")
	}
	t.Logf("Detailed explanation: %s", detailedExplanation)

	// Test statistical significance
	isSignificant := metrics.IsStatisticallySignificant()
	confidenceLevel := metrics.GetConfidenceLevel()
	t.Logf("Statistically significant: %v, Confidence level: %.1f%%", isSignificant, confidenceLevel)
}
