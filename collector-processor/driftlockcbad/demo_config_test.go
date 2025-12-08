package driftlockcbad

import (
	"math/rand"
	"testing"
	"time"
)

func TestDemoConfigDetectsSyntheticOutlier(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: false,
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	normalEvent := []byte(`{"type":"normal","value":100}`)
	for i := 0; i < detector.config.BaselineSize; i++ {
		if _, err := detector.AddData(normalEvent); err != nil {
			t.Fatalf("add baseline %d: %v", i, err)
		}
	}

	rand.Seed(time.Now().UnixNano())
	// Create a noisy window that should be obviously different from baseline
	for i := 0; i < detector.config.WindowSize; i++ {
		payload := make([]byte, 256)
		for j := range payload {
			payload[j] = byte(rand.Intn(26) + 97)
		}
		noisy := append([]byte(`{"type":"OUTLIER","blob":"`), payload...)
		noisy = append(noisy, []byte(`"}`)...)
		if _, err := detector.AddData(noisy); err != nil {
			t.Fatalf("add window %d: %v", i, err)
		}
	}

	ready, err := detector.IsReady()
	if err != nil {
		t.Fatalf("isReady: %v", err)
	}
	if !ready {
		t.Fatalf("detector not ready after %d events", detector.config.BaselineSize+detector.config.WindowSize)
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		t.Fatalf("detect anomaly: %v", err)
	}
	if metrics == nil {
		t.Fatalf("metrics nil")
	}

	if !detected {
		t.Fatalf("expected anomaly, got none (ncd=%.3f p=%.3f drop=%.3f entropy=%.3f)", metrics.NCD, metrics.PValue, metrics.CompressionRatioChange, metrics.EntropyChange)
	}
}

func TestDemoConfigDetectsSingleSpike(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: false,
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	normalEvent := []byte(`{"type":"normal","value":100}`)
	for i := 0; i < detector.config.BaselineSize; i++ {
		if _, err := detector.AddData(normalEvent); err != nil {
			t.Fatalf("add baseline %d: %v", i, err)
		}
	}

	// Window with mostly normal events plus a single large outlier
	for i := 0; i < detector.config.WindowSize; i++ {
		ev := normalEvent
		if i == detector.config.WindowSize-1 {
			ev = []byte(`{"type":"OUTLIER","value":9999999,"msg":"massive spike"}`)
		}
		if _, err := detector.AddData(ev); err != nil {
			t.Fatalf("add window %d: %v", i, err)
		}
	}

	ready, err := detector.IsReady()
	if err != nil {
		t.Fatalf("isReady: %v", err)
	}
	if !ready {
		t.Fatalf("detector not ready")
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		t.Fatalf("detect anomaly: %v", err)
	}
	if metrics == nil {
		t.Fatalf("metrics nil")
	}

	if !detected {
		t.Fatalf("expected anomaly for single spike, got none (ncd=%.3f p=%.3f drop=%.3f entropy=%.3f)", metrics.NCD, metrics.PValue, metrics.CompressionRatioChange, metrics.EntropyChange)
	}
}

func TestDemoConfigFlagsExtremeSingleOutlier(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: false,
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	normal := []byte(`{"type":"normal","value":100}`)
	outlier := []byte(`{"type":"OUTLIER","value":999999}`)

	for i := 0; i < detector.config.BaselineSize; i++ {
		if _, err := detector.AddData(normal); err != nil {
			t.Fatalf("add baseline %d: %v", i, err)
		}
	}

	for i := 0; i < detector.config.WindowSize; i++ {
		ev := normal
		if i == detector.config.WindowSize-1 {
			ev = outlier
		}
		if _, err := detector.AddData(ev); err != nil {
			t.Fatalf("add window %d: %v", i, err)
		}
	}

	ready, err := detector.IsReady()
	if err != nil {
		t.Fatalf("isReady: %v", err)
	}
	if !ready {
		t.Fatalf("detector not ready")
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		t.Fatalf("detect anomaly: %v", err)
	}
	if metrics == nil {
		t.Fatalf("metrics nil")
	}

	t.Logf("single spike metrics: ncd=%.4f p=%.4f drop=%.4f entropy=%.4f conf=%.3f base_ratio=%.3f window_ratio=%.3f", metrics.NCD, metrics.PValue, metrics.CompressionRatioChange, metrics.EntropyChange, metrics.ConfidenceLevel, metrics.BaselineCompressionRatio, metrics.WindowCompressionRatio)

	if !detected {
		t.Fatalf("expected anomaly when one event is a 10^6 spike; got none (ncd=%.3f p=%.3f drop=%.3f entropy=%.3f)", metrics.NCD, metrics.PValue, metrics.CompressionRatioChange, metrics.EntropyChange)
	}
}

// Mirrors the demo endpoint flow: streaming add + detect per event with formatted JSON.
func TestDemoStreamDetectsWhitespaceOutlier(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: false,
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	normal := []byte(`{"type": "normal", "value": 100}`)
	outlier := []byte(`{"type": "OUTLIER", "value": 999999}`)

	events := make([][]byte, 0, detector.config.BaselineSize+detector.config.WindowSize)
	for i := 0; i < detector.config.BaselineSize+detector.config.WindowSize-1; i++ {
		events = append(events, normal)
	}
	events = append(events, outlier)

	detectedCount := 0
	var lastMetrics *EnhancedMetrics

	for _, ev := range events {
		added, err := detector.AddData(ev)
		if err != nil {
			t.Fatalf("add data: %v", err)
		}
		if !added {
			continue
		}

		ready, err := detector.IsReady()
		if err != nil {
			t.Fatalf("isReady: %v", err)
		}
		if !ready {
			continue
		}

		detected, metrics, err := detector.DetectAnomaly()
		if err != nil {
			t.Fatalf("detect: %v", err)
		}
		if metrics != nil {
			lastMetrics = metrics
		}
		if detected {
			detectedCount++
			break
		}
	}

	if detectedCount == 0 {
		if lastMetrics == nil {
			t.Fatalf("expected metrics after streaming detect, got nil")
		}
		t.Fatalf("expected anomaly during streaming detect; last metrics ncd=%.3f p=%.3f drop=%.3f entropy=%.3f", lastMetrics.NCD, lastMetrics.PValue, lastMetrics.CompressionRatioChange, lastMetrics.EntropyChange)
	}
}

// ============================================================================
// Tests for RequireStatisticalSignificance=true (strict mode)
// These verify the production configuration eliminates false positives
// ============================================================================

// TestStatisticalSignificanceReducesFalsePositives verifies that homogeneous data
// does NOT trigger anomalies when RequireStatisticalSignificance=true.
// This is the key behavioral test for the "no more 60% false positive spam" fix.
func TestStatisticalSignificanceReducesFalsePositives(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               100, // More permutations for stable p-value
		Seed:                           42,
		RequireStatisticalSignificance: true, // KEY: strict mode
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	// Feed 100% homogeneous data - all identical events
	normalEvent := []byte(`{"type":"normal","value":100,"msg":"routine log entry"}`)
	totalEvents := detector.config.BaselineSize + detector.config.WindowSize

	for i := 0; i < totalEvents; i++ {
		if _, err := detector.AddData(normalEvent); err != nil {
			t.Fatalf("add event %d: %v", i, err)
		}
	}

	ready, err := detector.IsReady()
	if err != nil {
		t.Fatalf("isReady: %v", err)
	}
	if !ready {
		t.Fatalf("detector not ready after %d events", totalEvents)
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		t.Fatalf("detect anomaly: %v", err)
	}
	if metrics == nil {
		t.Fatalf("metrics nil")
	}

	t.Logf("Homogeneous data metrics: ncd=%.4f p=%.4f conf=%.3f statistically_significant=%v",
		metrics.NCD, metrics.PValue, metrics.ConfidenceLevel, metrics.IsStatisticallySignificant)

	// With RequireStatisticalSignificance=true, homogeneous data should NOT trigger anomaly
	// because p-value will be high (data is indistinguishable from baseline)
	if detected {
		t.Fatalf("FAIL: Homogeneous data triggered false positive with strict mode! "+
			"ncd=%.3f p=%.3f (should not detect when p > 0.05)",
			metrics.NCD, metrics.PValue)
	}

	// P-value should be high for identical data (no statistical difference)
	if metrics.PValue < 0.05 {
		t.Errorf("WARNING: P-value unexpectedly low (%.4f) for homogeneous data", metrics.PValue)
	}
}

// TestStatisticalSignificanceWithClearAnomaly verifies that a clear pattern break
// IS detected when RequireStatisticalSignificance=true, with p < 0.05.
func TestStatisticalSignificanceWithClearAnomaly(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	detector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               100,
		Seed:                           42,
		RequireStatisticalSignificance: true, // KEY: strict mode
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	// Baseline: normal structured JSON logs
	normalEvent := []byte(`{"type":"normal","level":"INFO","msg":"user login successful","status":200}`)
	for i := 0; i < detector.config.BaselineSize; i++ {
		if _, err := detector.AddData(normalEvent); err != nil {
			t.Fatalf("add baseline %d: %v", i, err)
		}
	}

	// Window: completely different structure - error stack traces
	errorEvent := []byte(`{"type":"PANIC","level":"FATAL","error":"OutOfMemoryError","stack":"java.lang.OutOfMemoryError: GC overhead limit exceeded\n\tat java.util.Arrays.copyOf(Arrays.java:3236)\n\tat java.util.ArrayList.grow(ArrayList.java:265)"}`)
	for i := 0; i < detector.config.WindowSize; i++ {
		if _, err := detector.AddData(errorEvent); err != nil {
			t.Fatalf("add window %d: %v", i, err)
		}
	}

	ready, err := detector.IsReady()
	if err != nil {
		t.Fatalf("isReady: %v", err)
	}
	if !ready {
		t.Fatalf("detector not ready")
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		t.Fatalf("detect anomaly: %v", err)
	}
	if metrics == nil {
		t.Fatalf("metrics nil")
	}

	t.Logf("Pattern break metrics: ncd=%.4f p=%.4f conf=%.3f statistically_significant=%v",
		metrics.NCD, metrics.PValue, metrics.ConfidenceLevel, metrics.IsStatisticallySignificant)

	// With RequireStatisticalSignificance=true, a clear pattern break SHOULD be detected
	if !detected {
		t.Fatalf("FAIL: Clear pattern break NOT detected in strict mode! "+
			"ncd=%.3f p=%.3f (expected detection with both ncd > 0.3 AND p < 0.05)",
			metrics.NCD, metrics.PValue)
	}

	// P-value should indicate statistical significance
	if metrics.PValue >= 0.05 {
		t.Errorf("WARNING: Pattern break detected but p-value not significant (%.4f >= 0.05)", metrics.PValue)
	}

	// Should be marked as statistically significant
	if !metrics.IsStatisticallySignificant {
		t.Errorf("WARNING: Detected anomaly not marked as statistically significant")
	}
}

// TestStrictModeVsRelaxedMode compares detection behavior between modes
// to demonstrate that strict mode is more conservative
func TestStrictModeVsRelaxedMode(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	// Create a dataset with a subtle pattern change (not dramatic)
	// This should trigger in relaxed mode but NOT in strict mode
	baseEvent := []byte(`{"type":"normal","value":100}`)
	subtleChange := []byte(`{"type":"normal","value":105}`) // Subtle 5% value change

	// Test relaxed mode (RequireStatisticalSignificance=false)
	relaxedDetector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: false, // Relaxed
	})
	if err != nil {
		t.Fatalf("create relaxed detector: %v", err)
	}
	defer relaxedDetector.Close()

	// Test strict mode (RequireStatisticalSignificance=true)
	strictDetector, err := NewDetector(DetectorConfig{
		BaselineSize:                   30,
		WindowSize:                     10,
		HopSize:                        5,
		MaxCapacity:                    30 + 4*10 + 1024,
		PValueThreshold:                0.05,
		NCDThreshold:                   0.3,
		CompressionAlgorithm:           "zstd",
		PermutationCount:               50,
		Seed:                           42,
		RequireStatisticalSignificance: true, // Strict
	})
	if err != nil {
		t.Fatalf("create strict detector: %v", err)
	}
	defer strictDetector.Close()

	// Feed identical data to both
	for i := 0; i < 30; i++ {
		relaxedDetector.AddData(baseEvent)
		strictDetector.AddData(baseEvent)
	}
	for i := 0; i < 10; i++ {
		relaxedDetector.AddData(subtleChange)
		strictDetector.AddData(subtleChange)
	}

	relaxedReady, _ := relaxedDetector.IsReady()
	strictReady, _ := strictDetector.IsReady()
	if !relaxedReady || !strictReady {
		t.Fatalf("detectors not ready")
	}

	relaxedDetected, relaxedMetrics, _ := relaxedDetector.DetectAnomaly()
	strictDetected, strictMetrics, _ := strictDetector.DetectAnomaly()

	t.Logf("Relaxed mode: detected=%v ncd=%.4f p=%.4f",
		relaxedDetected, relaxedMetrics.NCD, relaxedMetrics.PValue)
	t.Logf("Strict mode:  detected=%v ncd=%.4f p=%.4f",
		strictDetected, strictMetrics.NCD, strictMetrics.PValue)

	// The key insight: strict mode should be equal or MORE conservative
	// (fewer detections) than relaxed mode for the same data
	if strictDetected && !relaxedDetected {
		t.Errorf("Strict mode detected anomaly that relaxed mode missed - this is unexpected")
	}

	// Log comparison for visibility
	if relaxedDetected && !strictDetected {
		t.Logf("SUCCESS: Strict mode correctly filtered out a detection that relaxed mode flagged")
	}
}
