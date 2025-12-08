package main

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

// StreamCharacteristics holds computed data characteristics for adaptive sizing
type StreamCharacteristics struct {
	AvgEventsPerHour   float64
	AvgEventSizeBytes  int
	AvgBaselineEntropy float64
	PatternDiversity   float64 // 0-1 scale, higher = more varied
}

// AdaptiveWindowConfig defines bounds and scaling factors for window sizing
type AdaptiveWindowConfig struct {
	// Baseline sizing
	MinBaseline          int
	MaxBaseline          int
	BaselinePerEventRate float64 // Events/hour -> baseline multiplier

	// Window sizing
	MinWindow             int
	MaxWindow             int
	WindowToBaselineRatio float64 // Default window = baseline * ratio

	// Memory-based limits
	MaxMemoryPerStreamMB int
	AvgEventSizeForCalc  int // Default event size assumption
}

var defaultAdaptiveConfig = AdaptiveWindowConfig{
	MinBaseline:           100,
	MaxBaseline:           2000,
	BaselinePerEventRate:  0.1, // 1000 events/hour -> baseline of 100
	MinWindow:             10,
	MaxWindow:             200,
	WindowToBaselineRatio: 0.125, // Window = baseline / 8
	MaxMemoryPerStreamMB:  50,
	AvgEventSizeForCalc:   1024,
}

// computeAdaptiveWindowSizes calculates optimal baseline and window sizes
// based on stream characteristics
func computeAdaptiveWindowSizes(chars StreamCharacteristics, cfg AdaptiveWindowConfig) (baseline, window int) {
	// Factor 1: Event frequency
	// More frequent events -> larger baseline for stability
	freqBaseline := int(chars.AvgEventsPerHour * cfg.BaselinePerEventRate)
	freqBaseline = clampInt(freqBaseline, cfg.MinBaseline, cfg.MaxBaseline)

	// Factor 2: Event size (memory constraint)
	// Larger events -> smaller windows to stay within memory budget
	eventSize := chars.AvgEventSizeBytes
	if eventSize <= 0 {
		eventSize = cfg.AvgEventSizeForCalc
	}
	maxEventsByMemory := (cfg.MaxMemoryPerStreamMB * 1024 * 1024) / eventSize
	memoryBaseline := maxEventsByMemory / 4 // Reserve room for 3 windows
	memoryBaseline = clampInt(memoryBaseline, cfg.MinBaseline, cfg.MaxBaseline)

	// Factor 3: Pattern complexity
	// More diverse patterns -> larger baseline to capture full range
	complexityMultiplier := 1.0 + (chars.PatternDiversity * 0.5) // 1.0 to 1.5x
	complexityBaseline := int(float64(cfg.MinBaseline*2) * complexityMultiplier)
	complexityBaseline = clampInt(complexityBaseline, cfg.MinBaseline, cfg.MaxBaseline)

	// Factor 4: Entropy (data randomness)
	// Higher entropy data needs more samples for stable baseline
	entropyMultiplier := 1.0
	if chars.AvgBaselineEntropy > 6.0 { // High entropy threshold (bits/byte, max ~8)
		entropyMultiplier = 1.0 + (chars.AvgBaselineEntropy-6.0)*0.1
	}
	entropyBaseline := int(float64(cfg.MinBaseline*2) * entropyMultiplier)
	entropyBaseline = clampInt(entropyBaseline, cfg.MinBaseline, cfg.MaxBaseline)

	// Combine factors: take minimum of memory constraint, weighted average of others
	combinedBaseline := (freqBaseline + complexityBaseline + entropyBaseline) / 3
	baseline = min(combinedBaseline, memoryBaseline)
	baseline = clampInt(baseline, cfg.MinBaseline, cfg.MaxBaseline)

	// Window size: ratio of baseline, clamped
	window = int(float64(baseline) * cfg.WindowToBaselineRatio)
	window = clampInt(window, cfg.MinWindow, cfg.MaxWindow)

	// Ensure window doesn't exceed baseline
	if window > baseline {
		window = baseline / 2
	}

	return baseline, window
}

// deriveStreamCharacteristics computes lightweight statistics from a request payload.
// It intentionally samples a small subset to avoid heavy CPU work on large batches.
func deriveStreamCharacteristics(events []json.RawMessage) StreamCharacteristics {
	if len(events) == 0 {
		return StreamCharacteristics{}
	}

	limit := min(len(events), 200)
	samples := make([][]byte, 0, limit)
	var totalSize float64

	combined := bytes.Buffer{}
	for i := 0; i < limit; i++ {
		ev := bytes.TrimSpace([]byte(events[i]))
		samples = append(samples, ev)
		size := float64(len(ev))
		totalSize += size
		if combined.Len() < 512*1024 { // cap entropy sample to 512KB
			combined.Write(ev)
		}
	}

	count := float64(len(samples))
	avgSize := 0.0
	if count > 0 {
		avgSize = totalSize / count
	}

	patternDiversity := computePatternDiversity(samples)
	entropy := computeApproximateEntropy(combined.Bytes())

	// Approximate events/hour using batch size; conservative to avoid oversizing
	approxEventsPerHour := float64(len(events)) * 6 // assume ~10 minute batches

	return StreamCharacteristics{
		AvgEventsPerHour:   approxEventsPerHour,
		AvgEventSizeBytes:  int(avgSize),
		AvgBaselineEntropy: entropy,
		PatternDiversity:   patternDiversity,
	}
}

// StreamStatistics represents computed statistics stored in database
type StreamStatistics struct {
	ID                      uuid.UUID
	StreamID                uuid.UUID
	AvgEventsPerHour        float64
	AvgEventSizeBytes       int
	AvgBaselineEntropy      float64
	AvgNCDSelf              float64
	PatternDiversityScore   float64
	RecommendedBaselineSize int
	RecommendedWindowSize   int
	SampleCount             int
	LastComputedAt          time.Time
}

// updateStreamStatistics computes and stores stream characteristics
func (s *store) updateStreamStatistics(ctx context.Context, streamID uuid.UUID, chars StreamCharacteristics) error {
	baseline, window := computeAdaptiveWindowSizes(chars, defaultAdaptiveConfig)

	_, err := s.pool.Exec(ctx, `
        INSERT INTO stream_statistics
        (stream_id, avg_events_per_hour, avg_event_size_bytes,
         avg_baseline_entropy, pattern_diversity_score,
         recommended_baseline_size, recommended_window_size,
         sample_count, last_computed_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, 1, NOW())
        ON CONFLICT (stream_id) DO UPDATE SET
            avg_events_per_hour = (stream_statistics.avg_events_per_hour + EXCLUDED.avg_events_per_hour) / 2,
            avg_event_size_bytes = (stream_statistics.avg_event_size_bytes + EXCLUDED.avg_event_size_bytes) / 2,
            avg_baseline_entropy = (stream_statistics.avg_baseline_entropy + EXCLUDED.avg_baseline_entropy) / 2,
            pattern_diversity_score = (stream_statistics.pattern_diversity_score + EXCLUDED.pattern_diversity_score) / 2,
            recommended_baseline_size = EXCLUDED.recommended_baseline_size,
            recommended_window_size = EXCLUDED.recommended_window_size,
            sample_count = stream_statistics.sample_count + 1,
            last_computed_at = NOW()`,
		streamID, chars.AvgEventsPerHour, chars.AvgEventSizeBytes,
		chars.AvgBaselineEntropy, chars.PatternDiversity, baseline, window)

	return err
}

// getStreamStatistics retrieves computed statistics for a stream
func (s *store) getStreamStatistics(ctx context.Context, streamID uuid.UUID) (*StreamStatistics, error) {
	stats := &StreamStatistics{}

	err := s.pool.QueryRow(ctx, `
        SELECT id, stream_id, avg_events_per_hour, avg_event_size_bytes,
               COALESCE(avg_baseline_entropy, 0),
               COALESCE(avg_ncd_self, 0), COALESCE(pattern_diversity_score, 0),
               recommended_baseline_size, recommended_window_size,
               sample_count, last_computed_at
        FROM stream_statistics
        WHERE stream_id = $1`,
		streamID).Scan(
		&stats.ID, &stats.StreamID, &stats.AvgEventsPerHour, &stats.AvgEventSizeBytes,
		&stats.AvgBaselineEntropy, &stats.AvgNCDSelf, &stats.PatternDiversityScore,
		&stats.RecommendedBaselineSize, &stats.RecommendedWindowSize,
		&stats.SampleCount, &stats.LastComputedAt)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// getAdaptiveWindowSizes retrieves computed window sizes for a stream
func (s *store) getAdaptiveWindowSizes(ctx context.Context, streamID uuid.UUID) (baseline, window int, err error) {
	err = s.pool.QueryRow(ctx, `
        SELECT recommended_baseline_size, recommended_window_size
        FROM stream_statistics WHERE stream_id = $1`,
		streamID).Scan(&baseline, &window)
	return
}

// AdaptiveWindowBounds holds per-stream bounds for adaptive sizing
type AdaptiveWindowBounds struct {
	BaselineMin int
	BaselineMax int
	WindowMin   int
	WindowMax   int
}

// getAdaptiveWindowBounds retrieves bounds for a stream from database
func (s *store) getAdaptiveWindowBounds(ctx context.Context, streamID uuid.UUID) (*AdaptiveWindowBounds, error) {
	bounds := &AdaptiveWindowBounds{}

	err := s.pool.QueryRow(ctx, `
        SELECT
            COALESCE(adaptive_baseline_min, 100),
            COALESCE(adaptive_baseline_max, 2000),
            COALESCE(adaptive_window_min, 10),
            COALESCE(adaptive_window_max, 200)
        FROM streams
        WHERE id = $1`,
		streamID).Scan(
		&bounds.BaselineMin, &bounds.BaselineMax,
		&bounds.WindowMin, &bounds.WindowMax)

	if err != nil {
		return nil, err
	}

	return bounds, nil
}

// computePatternDiversity estimates pattern diversity from a sample of events
// Returns a score from 0 (uniform) to 1 (highly diverse)
func computePatternDiversity(eventSamples [][]byte) float64 {
	if len(eventSamples) < 2 {
		return 0.0
	}

	// Simple heuristic: measure length variance as proxy for diversity
	var sumLen, sumSqLen float64
	for _, ev := range eventSamples {
		l := float64(len(ev))
		sumLen += l
		sumSqLen += l * l
	}

	n := float64(len(eventSamples))
	mean := sumLen / n
	variance := (sumSqLen / n) - (mean * mean)

	if mean == 0 {
		return 0.0
	}

	// Coefficient of variation as diversity score (capped at 1.0)
	cv := math.Sqrt(variance) / mean
	return math.Min(cv, 1.0)
}

// computeApproximateEntropy estimates entropy from event data
// Returns bits per byte (0-8 range)
func computeApproximateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	// Count byte frequencies
	var counts [256]int
	for _, b := range data {
		counts[b]++
	}

	// Calculate entropy
	total := float64(len(data))
	var entropy float64
	for _, count := range counts {
		if count > 0 {
			p := float64(count) / total
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}
