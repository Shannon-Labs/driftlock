// Package driftlockcbad provides compression-based anomaly detection.
// This file implements SHA-142: Numeric Value Histogramming for value outlier detection.

package driftlockcbad

import (
	"encoding/json"
	"math"
	"sync"
)

// NumericStats tracks running statistics using Welford's online algorithm.
// This allows single-pass computation of mean and variance without storing all values.
// Memory efficient: O(1) per field path regardless of event count.
type NumericStats struct {
	N       int64   // Count of values seen
	Mean    float64 // Running mean
	M2      float64 // Sum of squared differences from mean (for variance)
	Min     float64 // Minimum value seen
	Max     float64 // Maximum value seen
	mu      sync.Mutex
}

// NewNumericStats creates a new stats tracker
func NewNumericStats() *NumericStats {
	return &NumericStats{
		Min: math.MaxFloat64,
		Max: -math.MaxFloat64,
	}
}

// Update incorporates a new value using Welford's algorithm.
// This is numerically stable for computing variance in a single pass.
func (s *NumericStats) Update(x float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.N++
	delta := x - s.Mean
	s.Mean += delta / float64(s.N)
	s.M2 += delta * (x - s.Mean) // Note: uses updated mean

	if x < s.Min {
		s.Min = x
	}
	if x > s.Max {
		s.Max = x
	}
}

// Variance returns the sample variance (N-1 denominator).
// Returns 0 if fewer than 2 samples.
func (s *NumericStats) Variance() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.N < 2 {
		return 0
	}
	return s.M2 / float64(s.N-1)
}

// StdDev returns the sample standard deviation.
func (s *NumericStats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// IsOutlier checks if a value is more than k standard deviations from the mean.
// Common values for k: 2.0 (95%), 3.0 (99.7%), 4.0 (99.99%)
// Returns false if insufficient data (N < 30) for statistical significance.
func (s *NumericStats) IsOutlier(x float64, k float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Need enough samples for statistical significance
	if s.N < 30 {
		return false
	}

	stddev := math.Sqrt(s.M2 / float64(s.N-1))
	if stddev == 0 {
		// All values identical - any different value is an outlier
		return x != s.Mean
	}

	return math.Abs(x-s.Mean) > k*stddev
}

// ZScore calculates how many standard deviations x is from the mean.
// Returns 0 if insufficient data.
func (s *NumericStats) ZScore(x float64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.N < 2 {
		return 0
	}

	stddev := math.Sqrt(s.M2 / float64(s.N-1))
	if stddev == 0 {
		return 0
	}

	return (x - s.Mean) / stddev
}

// Snapshot returns a copy of the current statistics
func (s *NumericStats) Snapshot() NumericStatsSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	stddev := 0.0
	if s.N >= 2 {
		stddev = math.Sqrt(s.M2 / float64(s.N-1))
	}

	return NumericStatsSnapshot{
		N:      s.N,
		Mean:   s.Mean,
		StdDev: stddev,
		Min:    s.Min,
		Max:    s.Max,
	}
}

// NumericStatsSnapshot is a point-in-time copy of statistics
type NumericStatsSnapshot struct {
	N      int64   `json:"n"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

// NumericStatsRegistry manages stats for multiple streams and field paths.
// Thread-safe for concurrent access.
type NumericStatsRegistry struct {
	// Map structure: streamID -> fieldPath -> stats
	stats map[string]map[string]*NumericStats
	mu    sync.RWMutex
}

// NewNumericStatsRegistry creates a new registry
func NewNumericStatsRegistry() *NumericStatsRegistry {
	return &NumericStatsRegistry{
		stats: make(map[string]map[string]*NumericStats),
	}
}

// GetOrCreate returns stats for a stream/field, creating if needed
func (r *NumericStatsRegistry) GetOrCreate(streamID, fieldPath string) *NumericStats {
	r.mu.RLock()
	if streamStats, ok := r.stats[streamID]; ok {
		if stats, ok := streamStats[fieldPath]; ok {
			r.mu.RUnlock()
			return stats
		}
	}
	r.mu.RUnlock()

	// Upgrade to write lock
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if _, ok := r.stats[streamID]; !ok {
		r.stats[streamID] = make(map[string]*NumericStats)
	}
	if _, ok := r.stats[streamID][fieldPath]; !ok {
		r.stats[streamID][fieldPath] = NewNumericStats()
	}

	return r.stats[streamID][fieldPath]
}

// GetStreamStats returns all stats for a stream
func (r *NumericStatsRegistry) GetStreamStats(streamID string) map[string]NumericStatsSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]NumericStatsSnapshot)
	if streamStats, ok := r.stats[streamID]; ok {
		for path, stats := range streamStats {
			result[path] = stats.Snapshot()
		}
	}
	return result
}

// Clear removes all stats for a stream
func (r *NumericStatsRegistry) Clear(streamID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.stats, streamID)
}

// ValueOutlier represents a detected numeric value anomaly
type ValueOutlier struct {
	FieldPath    string  `json:"field_path"`
	Value        float64 `json:"value"`
	ExpectedMean float64 `json:"expected_mean"`
	StdDev       float64 `json:"std_dev"`
	ZScore       float64 `json:"z_score"`
	SampleCount  int64   `json:"sample_count"`
}

// NumericFieldExtractor extracts numeric values from JSON with their paths
type NumericFieldExtractor struct {
	// IncludePatterns: if non-empty, only fields matching these patterns are tracked
	IncludePatterns []string
	// ExcludePatterns: fields matching these patterns are ignored
	ExcludePatterns []string
}

// ExtractNumericFields recursively extracts all numeric fields from JSON.
// Returns map of field paths (dot-notation) to values.
// Example: {"user": {"age": 25}} -> {"user.age": 25}
func ExtractNumericFields(data []byte) (map[string]float64, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	extractNumeric("", obj, result)
	return result, nil
}

// extractNumeric recursively walks JSON and extracts numeric values
func extractNumeric(prefix string, v interface{}, result map[string]float64) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v := range val {
			path := k
			if prefix != "" {
				path = prefix + "." + k
			}
			extractNumeric(path, v, result)
		}
	case []interface{}:
		// For arrays, we could track by index or aggregate
		// For now, skip arrays to avoid path explosion
		// TODO: Consider tracking array element stats separately
	case float64:
		result[prefix] = val
	case int:
		result[prefix] = float64(val)
	case int64:
		result[prefix] = float64(val)
	}
}

// CheckForOutliers analyzes numeric fields in an event against historical stats.
// Returns any fields that exceed the k-sigma threshold.
func CheckForOutliers(
	registry *NumericStatsRegistry,
	streamID string,
	event []byte,
	kSigma float64,
	updateStats bool,
) ([]ValueOutlier, error) {
	fields, err := ExtractNumericFields(event)
	if err != nil {
		return nil, err
	}

	var outliers []ValueOutlier

	for path, value := range fields {
		stats := registry.GetOrCreate(streamID, path)

		// Check if outlier BEFORE updating (to compare against prior distribution)
		if stats.IsOutlier(value, kSigma) {
			snapshot := stats.Snapshot()
			outliers = append(outliers, ValueOutlier{
				FieldPath:    path,
				Value:        value,
				ExpectedMean: snapshot.Mean,
				StdDev:       snapshot.StdDev,
				ZScore:       stats.ZScore(value),
				SampleCount:  snapshot.N,
			})
		}

		// Update stats with new value
		if updateStats {
			stats.Update(value)
		}
	}

	return outliers, nil
}
