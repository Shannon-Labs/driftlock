package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
)

const defaultAnchorThreshold = 0.35

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildSnapshot serializes a subset of events into a canonical, newline-delimited payload.
// If tail is true, the last "limit" events are used; otherwise the first "limit" events are taken.
func buildSnapshot(events []json.RawMessage, tokenizer *driftlockcbad.Tokenizer, limit int, tail bool) ([]byte, int, error) {
	if len(events) == 0 {
		return nil, 0, nil
	}
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}

	start, end := 0, len(events)
	if limit < len(events) {
		if tail {
			start = len(events) - limit
		} else {
			end = limit
		}
	}

	var buf bytes.Buffer
	count := 0
	for i := start; i < end; i++ {
		compact, err := compactEvent(events[i], tokenizer)
		if err != nil {
			return nil, count, err
		}
		buf.Write(compact)
		if i < end-1 {
			buf.WriteByte('\n')
		}
		count++
	}
	return buf.Bytes(), count, nil
}

func compactEvent(ev json.RawMessage, tokenizer *driftlockcbad.Tokenizer) ([]byte, error) {
	data := bytes.TrimSpace(ev)
	if tokenizer != nil {
		data = tokenizer.Tokenize(data)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("event is empty")
	}
	if !json.Valid(data) {
		return nil, fmt.Errorf("event is not valid json")
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return nil, fmt.Errorf("compact json: %w", err)
	}
	return buf.Bytes(), nil
}

// decodeAnchorData returns the anchor payload in its raw form regardless of compression.
func decodeAnchorData(data []byte, compressor string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("anchor data is empty")
	}
	switch strings.ToLower(compressor) {
	case "", "raw", "identity", "json", "none":
		return data, nil
	case "zstd":
		decoder, err := zstd.NewReader(nil)
		if err != nil {
			return nil, fmt.Errorf("init zstd decoder: %w", err)
		}
		defer decoder.Close()
		decoded, err := decoder.DecodeAll(data, nil)
		if err != nil {
			return nil, fmt.Errorf("decode zstd anchor: %w", err)
		}
		return decoded, nil
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("init gzip reader: %w", err)
		}
		defer reader.Close()
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("decode gzip anchor: %w", err)
		}
		return decoded, nil
	default:
		return nil, fmt.Errorf("unsupported anchor compressor %q", compressor)
	}
}

// newAnchorRecord builds a StreamAnchor populated with baseline metrics when available.
func newAnchorRecord(streamID uuid.UUID, compressor string, baseline []byte, eventCount int, completedAt time.Time, driftThreshold float64, seed uint64, permutations int) *StreamAnchor {
	if len(baseline) == 0 || streamID == uuid.Nil {
		return nil
	}
	if driftThreshold <= 0 {
		driftThreshold = defaultAnchorThreshold
	}
	snapshot := append([]byte(nil), baseline...)
	anchor := &StreamAnchor{
		StreamID:               streamID,
		AnchorData:             snapshot,
		Compressor:             compressor,
		EventCount:             eventCount,
		CalibrationCompletedAt: completedAt,
		IsActive:               true,
		DriftNCDThreshold:      driftThreshold,
	}

	metrics, err := driftlockcbad.ComputeMetrics(snapshot, snapshot, seed, permutations)
	if err != nil {
		log.Printf("anchor: compute baseline metrics failed (non-fatal): %v", err)
		return anchor
	}
	if metrics != nil {
		entropy := metrics.BaselineEntropy
		cr := metrics.BaselineCompressionRatio
		ncd := metrics.NCD
		anchor.BaselineEntropy = &entropy
		anchor.BaselineCompressionRatio = &cr
		anchor.BaselineNCDSelf = &ncd
	}
	return anchor
}
