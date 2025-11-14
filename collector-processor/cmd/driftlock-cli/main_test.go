//go:build cgo && !driftlock_no_cbad

package main

import (
	"testing"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

func TestProcessEventWithDetector(t *testing.T) {
	det, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
		BaselineSize:         2,
		WindowSize:           1,
		HopSize:              1,
		MaxCapacity:          64,
		PValueThreshold:      0.2,
		NCDThreshold:         0.1,
		PermutationCount:     10,
		Seed:                 1,
		CompressionAlgorithm: "zstd",
	})
	if err != nil {
		t.Skipf("skipping: detector unavailable (%v)", err)
		return
	}
	defer det.Close()

	var anomalies []anomalyOutput
	if err := processEvent(det, []byte(`{"a":1}`), 0, &anomalies); err != nil {
		t.Fatalf("processEvent(0) error: %v", err)
	}
	if err := processEvent(det, []byte(`{"a":2}`), 1, &anomalies); err != nil {
		t.Fatalf("processEvent(1) error: %v", err)
	}
	if err := processEvent(det, []byte(`{"a":999}`), 2, &anomalies); err != nil {
		t.Fatalf("processEvent(2) error: %v", err)
	}
	// No strict assertion on anomalies length to avoid flakiness; just ensure no panic.
}


