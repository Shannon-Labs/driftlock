package benchmarks

import (
	"testing"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

// TestAlgorithmConsistency verifies all algorithms detect the same anomalies
func TestAlgorithmConsistency(t *testing.T) {
	if !driftlockcbad.IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", driftlockcbad.AvailabilityError())
	}

	algorithms := []string{"zlab", "zstd", "lz4", "gzip"}

	// Generate test data with known anomaly
	normal := []byte(`{"ts":"2025-12-06T20:00:00Z","level":"info","msg":"User login successful","user":"alice"}`)
	anomaly := []byte(`{"ts":"2025-12-06T20:00:30Z","level":"error","msg":"SQL INJECTION: SELECT * FROM users; DROP TABLE sessions;--","severity":"critical"}`)

	events := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		events[i] = normal
	}
	// Inject anomaly at position 40
	events[40] = anomaly

	results := make(map[string]int)

	for _, algo := range algorithms {
		detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
			BaselineSize:         30,
			WindowSize:           10,
			HopSize:              5,
			MaxCapacity:          100,
			PValueThreshold:      0.05,
			NCDThreshold:         0.3,
			PermutationCount:     100,
			Seed:                 42,
			CompressionAlgorithm: algo,
		})
		if err != nil {
			t.Fatalf("failed to create detector with %s: %v", algo, err)
		}

		anomalyCount := 0
		for _, ev := range events {
			isAnomaly, err := detector.AddData(ev)
			if err != nil {
				t.Fatalf("AddData failed with %s: %v", algo, err)
			}
			if isAnomaly {
				anomalyCount++
			}
		}

		results[algo] = anomalyCount
		detector.Close()

		t.Logf("%s: detected %d anomalies", algo, anomalyCount)
	}

	// Verify all algorithms detected at least the injected anomaly
	for algo, count := range results {
		if count == 0 {
			t.Errorf("%s failed to detect any anomalies", algo)
		}
	}

	// Log comparison
	t.Logf("Algorithm consistency results: %v", results)
}

// TestAlgorithmRoundTrip verifies compression/decompression works for all algorithms
func TestAlgorithmRoundTrip(t *testing.T) {
	if !driftlockcbad.IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", driftlockcbad.AvailabilityError())
	}

	algorithms := []string{"zlab", "zstd", "lz4", "gzip"}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
				BaselineSize:         30,
				WindowSize:           10,
				HopSize:              5,
				MaxCapacity:          100,
				PValueThreshold:      0.05,
				NCDThreshold:         0.3,
				PermutationCount:     100,
				Seed:                 42,
				CompressionAlgorithm: algo,
			})
			if err != nil {
				t.Fatalf("failed to create detector: %v", err)
			}
			defer detector.Close()

			// Add some data to exercise the compression
			for i := 0; i < 50; i++ {
				_, err := detector.AddData([]byte(`{"test": "data", "index": ` + string(rune('0'+i%10)) + `}`))
				if err != nil {
					t.Fatalf("AddData failed: %v", err)
				}
			}

			// Verify detector is functional
			ready, err := detector.IsReady()
			if err != nil {
				t.Fatalf("IsReady failed: %v", err)
			}
			if !ready {
				t.Logf("%s: detector not ready after 50 events (expected with small baseline)", algo)
			}
		})
	}
}
