package driftlockcbad

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func ensureRustCore(t *testing.T) {
	t.Helper()
	if err := ValidateLibrary(); err != nil {
		t.Skipf("cbad-core not available: %v", err)
	}
}

func TestComputeMetricsDeterministic(t *testing.T) {
	ensureRustCore(t)

	baseline := []byte("baseline:payments:v1::" + strings.Repeat("N", 64))
	window := []byte("window:payments:v1::" + strings.Repeat("A", 16))

	first, err := ComputeMetrics(baseline, window, 1337, 256)
	require.NoError(t, err)
	second, err := ComputeMetrics(baseline, window, 1337, 256)
	require.NoError(t, err)

	require.Equal(t, first, second, "metrics should be deterministic for identical inputs")
}

func TestDetectorStreamingDeterminism(t *testing.T) {
	ensureRustCore(t)

	events := []string{
		`{"amount": 100, "endpoint": "/v1/charges", "latency_ms": 91}`,
		`{"amount": 105, "endpoint": "/v1/charges", "latency_ms": 89}`,
		`{"amount": 110, "endpoint": "/v1/charges", "latency_ms": 92}`,
		`{"amount": 100, "endpoint": "/v1/charges", "latency_ms": 90}`,
		`{"amount": 4200, "endpoint": "/v1/charges", "latency_ms": 810}`,
		`{"amount": 95, "endpoint": "/v1/refunds", "latency_ms": 95}`,
		`{"amount": 97, "endpoint": "/v1/refunds", "latency_ms": 96}`,
	}

	first := playbackDetector(t, events)
	second := playbackDetector(t, events)

	require.Equal(t, first, second, "detector output must remain deterministic with identical inputs")
}

type detectionSnapshot struct {
	Detected bool
	Metrics  EnhancedMetrics
}

func playbackDetector(t *testing.T, events []string) []detectionSnapshot {
	t.Helper()

	det, err := NewDetector(DetectorConfig{
		BaselineSize:         4,
		WindowSize:           1,
		HopSize:              1,
		MaxCapacity:          256,
		CompressionAlgorithm: "zstd",
		PermutationCount:     128,
		Seed:                 99,
	})
	require.NoError(t, err)
	t.Cleanup(det.Close)

	var snapshots []detectionSnapshot

	for _, raw := range events {
		added, err := det.AddData([]byte(raw))
		require.NoError(t, err)
		if !added {
			continue
		}
		ready, err := det.IsReady()
		require.NoError(t, err)
		if !ready {
			continue
		}
		detected, metrics, err := det.DetectAnomaly()
		require.NoError(t, err)
		snap := detectionSnapshot{Detected: detected}
		if metrics != nil {
			snap.Metrics = *metrics
		}
		snapshots = append(snapshots, snap)
	}

	return snapshots
}
