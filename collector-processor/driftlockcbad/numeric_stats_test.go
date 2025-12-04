package driftlockcbad

import (
	"math"
	"sync"
	"testing"
)

func TestNumericStats_WelfordAlgorithm(t *testing.T) {
	stats := NewNumericStats()

	// Known values: 2, 4, 4, 4, 5, 5, 7, 9
	// Mean = 5, Variance = 4, StdDev = 2
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	for _, v := range values {
		stats.Update(v)
	}

	if stats.N != 8 {
		t.Errorf("N: got %d, want 8", stats.N)
	}

	if math.Abs(stats.Mean-5.0) > 0.0001 {
		t.Errorf("Mean: got %f, want 5.0", stats.Mean)
	}

	variance := stats.Variance()
	if math.Abs(variance-4.0) > 0.0001 {
		t.Errorf("Variance: got %f, want 4.0", variance)
	}

	stddev := stats.StdDev()
	if math.Abs(stddev-2.0) > 0.0001 {
		t.Errorf("StdDev: got %f, want 2.0", stddev)
	}

	if stats.Min != 2 {
		t.Errorf("Min: got %f, want 2", stats.Min)
	}
	if stats.Max != 9 {
		t.Errorf("Max: got %f, want 9", stats.Max)
	}
}

func TestNumericStats_IsOutlier(t *testing.T) {
	stats := NewNumericStats()

	// Add 100 values around mean=50, stddev=10
	for i := 0; i < 100; i++ {
		// Values: 40, 41, 42, ..., 59, 40, 41, ... (cycling through 40-59)
		v := 40 + float64(i%20)
		stats.Update(v)
	}

	// Mean should be ~49.5, stddev ~5.77
	snapshot := stats.Snapshot()
	t.Logf("Mean: %f, StdDev: %f", snapshot.Mean, snapshot.StdDev)

	// Test outlier detection with 3-sigma
	tests := []struct {
		value    float64
		kSigma   float64
		expected bool
	}{
		{50, 3.0, false},  // Normal value
		{45, 3.0, false},  // Within range
		{55, 3.0, false},  // Within range
		{100, 3.0, true},  // Way outside
		{0, 3.0, true},    // Way outside
		{30, 3.0, true},   // > 3 stddev below mean
		{70, 3.0, true},   // > 3 stddev above mean
		{50, 2.0, false},  // Tighter threshold, still normal
	}

	for _, tt := range tests {
		got := stats.IsOutlier(tt.value, tt.kSigma)
		if got != tt.expected {
			t.Errorf("IsOutlier(%f, %f): got %v, want %v", tt.value, tt.kSigma, got, tt.expected)
		}
	}
}

func TestNumericStats_InsufficientData(t *testing.T) {
	stats := NewNumericStats()

	// With < 30 samples, should never detect outliers
	for i := 0; i < 29; i++ {
		stats.Update(float64(i))
	}

	// Even an extreme value shouldn't be flagged with insufficient data
	if stats.IsOutlier(1000, 3.0) {
		t.Error("Should not detect outlier with < 30 samples")
	}

	// Add one more to reach threshold
	stats.Update(29)
	// Now outliers can be detected
	if !stats.IsOutlier(1000, 3.0) {
		t.Error("Should detect outlier with >= 30 samples")
	}
}

func TestNumericStats_Concurrent(t *testing.T) {
	stats := NewNumericStats()
	var wg sync.WaitGroup

	// Concurrent updates
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stats.Update(float64(base + j))
			}
		}(i * 100)
	}

	wg.Wait()

	if stats.N != 1000 {
		t.Errorf("N: got %d, want 1000", stats.N)
	}
}

func TestNumericStatsRegistry(t *testing.T) {
	registry := NewNumericStatsRegistry()

	// Get or create stats
	stats1 := registry.GetOrCreate("stream-1", "user.age")
	stats1.Update(25)
	stats1.Update(30)
	stats1.Update(35)

	stats2 := registry.GetOrCreate("stream-1", "order.amount")
	stats2.Update(100)
	stats2.Update(150)

	stats3 := registry.GetOrCreate("stream-2", "user.age")
	stats3.Update(40)

	// Verify isolation
	if stats1.N != 3 {
		t.Errorf("stream-1 user.age N: got %d, want 3", stats1.N)
	}
	if stats2.N != 2 {
		t.Errorf("stream-1 order.amount N: got %d, want 2", stats2.N)
	}
	if stats3.N != 1 {
		t.Errorf("stream-2 user.age N: got %d, want 1", stats3.N)
	}

	// Get stream stats
	stream1Stats := registry.GetStreamStats("stream-1")
	if len(stream1Stats) != 2 {
		t.Errorf("stream-1 stats count: got %d, want 2", len(stream1Stats))
	}

	// Clear stream
	registry.Clear("stream-1")
	stream1Stats = registry.GetStreamStats("stream-1")
	if len(stream1Stats) != 0 {
		t.Errorf("stream-1 stats after clear: got %d, want 0", len(stream1Stats))
	}

	// stream-2 should be unaffected
	stream2Stats := registry.GetStreamStats("stream-2")
	if len(stream2Stats) != 1 {
		t.Errorf("stream-2 stats: got %d, want 1", len(stream2Stats))
	}
}

func TestExtractNumericFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]float64
	}{
		{
			name:  "simple object",
			input: `{"age": 25, "score": 85.5}`,
			expected: map[string]float64{
				"age":   25,
				"score": 85.5,
			},
		},
		{
			name:  "nested object",
			input: `{"user": {"age": 30, "profile": {"score": 100}}}`,
			expected: map[string]float64{
				"user.age":           30,
				"user.profile.score": 100,
			},
		},
		{
			name:  "mixed types",
			input: `{"name": "test", "count": 42, "active": true}`,
			expected: map[string]float64{
				"count": 42,
			},
		},
		{
			name:     "arrays skipped",
			input:    `{"values": [1, 2, 3], "total": 6}`,
			expected: map[string]float64{"total": 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractNumericFields([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("field count: got %d, want %d", len(result), len(tt.expected))
			}

			for path, expectedVal := range tt.expected {
				if gotVal, ok := result[path]; !ok {
					t.Errorf("missing field: %s", path)
				} else if math.Abs(gotVal-expectedVal) > 0.0001 {
					t.Errorf("field %s: got %f, want %f", path, gotVal, expectedVal)
				}
			}
		})
	}
}

func TestCheckForOutliers(t *testing.T) {
	registry := NewNumericStatsRegistry()
	streamID := "test-stream"

	// Populate with normal values (need >= 30 for outlier detection)
	for i := 0; i < 50; i++ {
		event := []byte(`{"amount": 100, "count": 10}`)
		_, err := CheckForOutliers(registry, streamID, event, 3.0, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Now send an outlier
	outlierEvent := []byte(`{"amount": 10000, "count": 10}`)
	outliers, err := CheckForOutliers(registry, streamID, outlierEvent, 3.0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outliers) != 1 {
		t.Errorf("outlier count: got %d, want 1", len(outliers))
	} else {
		if outliers[0].FieldPath != "amount" {
			t.Errorf("outlier field: got %s, want amount", outliers[0].FieldPath)
		}
		if outliers[0].Value != 10000 {
			t.Errorf("outlier value: got %f, want 10000", outliers[0].Value)
		}
	}
}

func BenchmarkNumericStats_Update(b *testing.B) {
	stats := NewNumericStats()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats.Update(float64(i % 100))
	}
}

func BenchmarkExtractNumericFields(b *testing.B) {
	event := []byte(`{
		"timestamp": "2025-01-15T10:30:00Z",
		"user": {"id": "123", "age": 30, "score": 85.5},
		"transaction": {"amount": 150.00, "fee": 2.50, "count": 1}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractNumericFields(event)
	}
}

func BenchmarkCheckForOutliers(b *testing.B) {
	registry := NewNumericStatsRegistry()
	streamID := "bench-stream"
	event := []byte(`{"amount": 100, "count": 10, "score": 85.5}`)

	// Pre-populate with 100 events
	for i := 0; i < 100; i++ {
		CheckForOutliers(registry, streamID, event, 3.0, true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckForOutliers(registry, streamID, event, 3.0, true)
	}
}
