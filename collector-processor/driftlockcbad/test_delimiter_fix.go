package driftlockcbad

import (
	"testing"
)

func TestDelimiterFixWithStatisticalSignificance(t *testing.T) {
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
		RequireStatisticalSignificance: true, // Enable statistical significance
	})
	if err != nil {
		t.Fatalf("create detector: %v", err)
	}
	defer detector.Close()

	normal := []byte(`{"type":"normal","value":100}`)
	outlier := []byte(`{"type":"OUTLIER","value":999999}`)

	// Add baseline events
	for i := 0; i < detector.config.BaselineSize; i++ {
		if _, err := detector.AddData(normal); err != nil {
			t.Fatalf("add baseline %d: %v", i, err)
		}
	}

	// Add window events with one extreme outlier
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

	t.Logf("delimiter fix test metrics:")
	t.Logf("  NCD: %.4f", metrics.NCD)
	t.Logf("  P-Value: %.4f", metrics.PValue)
	t.Logf("  Compression Ratio Change: %.4f", metrics.CompressionRatioChange)
	t.Logf("  Entropy Change: %.4f", metrics.EntropyChange)
	t.Logf("  Confidence Level: %.3f", metrics.ConfidenceLevel)
	t.Logf("  Baseline Compression Ratio: %.3f", metrics.BaselineCompressionRatio)
	t.Logf("  Window Compression Ratio: %.3f", metrics.WindowCompressionRatio)

	// With the delimiter fix, permutation testing should work correctly
	// and we should get a low p-value (< 0.05) for the obvious outlier
	if metrics.PValue >= 0.05 {
		t.Fatalf("expected p-value < 0.05 for obvious outlier; got %.4f", metrics.PValue)
	}

	// The anomaly should be detected
	if !detected {
		t.Fatalf("expected anomaly when one event is a 10^6 spike; got none")
	}

	t.Logf("✓ Delimiter fix working correctly - p-value: %.4f, anomaly detected: %v", metrics.PValue, detected)
}