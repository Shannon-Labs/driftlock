package main

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
	"github.com/google/uuid"
)

func TestBuildSnapshotHeadTail(t *testing.T) {
	events := []json.RawMessage{
		[]byte(`{"a":1}`),
		[]byte(`{"b":2}`),
		[]byte(`{"c":3}`),
	}
	head, count, err := buildSnapshot(events, nil, 2, false)
	if err != nil {
		t.Fatalf("head snapshot failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 events, got %d", count)
	}
	if got := string(head); got != `{"a":1}
{"b":2}` {
		t.Fatalf("unexpected head snapshot: %q", got)
	}

	tail, tailCount, err := buildSnapshot(events, nil, 2, true)
	if err != nil {
		t.Fatalf("tail snapshot failed: %v", err)
	}
	if tailCount != 2 {
		t.Fatalf("expected 2 tail events, got %d", tailCount)
	}
	if got := string(tail); got != `{"b":2}
{"c":3}` {
		t.Fatalf("unexpected tail snapshot: %q", got)
	}
}

func TestDecodeAnchorDataRaw(t *testing.T) {
	raw := []byte("baseline-bytes")
	out, err := decodeAnchorData(raw, "raw")
	if err != nil {
		t.Fatalf("decodeAnchorData returned error: %v", err)
	}
	if string(out) != "baseline-bytes" {
		t.Fatalf("unexpected anchor payload: %q", string(out))
	}
}

func TestNewAnchorRecordDefaults(t *testing.T) {
	streamID := uuid.New()
	baseline := []byte(`{"foo":"bar"}`)
	anchor := newAnchorRecord(streamID, "raw", baseline, 3, time.Unix(0, 0), 0, 42, 10)
	if anchor == nil {
		t.Fatalf("expected anchor record")
	}
	if anchor.StreamID != streamID {
		t.Fatalf("stream id mismatch")
	}
	if anchor.DriftNCDThreshold != defaultAnchorThreshold {
		t.Fatalf("expected default threshold %v, got %v", defaultAnchorThreshold, anchor.DriftNCDThreshold)
	}
	if anchor.EventCount != 3 {
		t.Fatalf("expected event count 3, got %d", anchor.EventCount)
	}
	if &anchor.AnchorData[0] == &baseline[0] {
		t.Fatalf("baseline payload should be copied")
	}

	// Metrics are optional if the CBAD library is unavailable.
	if err := driftlockcbad.ValidateLibrary(); err != nil {
		return
	}
	if anchor.BaselineEntropy == nil || anchor.BaselineCompressionRatio == nil || anchor.BaselineNCDSelf == nil {
		t.Fatalf("expected baseline metrics when CBAD is available")
	}
	if anchor.BaselineNCDSelf != nil && math.IsNaN(*anchor.BaselineNCDSelf) {
		t.Fatalf("baseline NCD self should not be NaN")
	}
}
