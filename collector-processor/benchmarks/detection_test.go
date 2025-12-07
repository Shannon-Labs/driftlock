package benchmarks

import (
	"testing"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

func BenchmarkDetectionSmall(b *testing.B) {
	benchDetect(b, 100, 100, 20)
}

func BenchmarkDetectionMedium(b *testing.B) {
	benchDetect(b, 1000, 200, 40)
}

func BenchmarkDetectionLarge(b *testing.B) {
	benchDetect(b, 10000, 400, 80)
}

func BenchmarkMemoryUsage(b *testing.B) {
	// The detector uses bounded windows; track allocations during repeated runs.
	b.ReportAllocs()
	benchDetect(b, 1000, 200, 40)
}

// BenchmarkAlgorithmComparison compares all compression algorithms
func BenchmarkAlgorithmZlab(b *testing.B) {
	benchDetectWithAlgo(b, 500, 100, 20, "zlab")
}

func BenchmarkAlgorithmZstd(b *testing.B) {
	benchDetectWithAlgo(b, 500, 100, 20, "zstd")
}

func BenchmarkAlgorithmLz4(b *testing.B) {
	benchDetectWithAlgo(b, 500, 100, 20, "lz4")
}

func BenchmarkAlgorithmGzip(b *testing.B) {
	benchDetectWithAlgo(b, 500, 100, 20, "gzip")
}

func benchDetectWithAlgo(b *testing.B, totalEvents, baselineSize, windowSize int, algo string) {
	if !driftlockcbad.IsAvailable() {
		b.Skipf("CBAD core unavailable: %v", driftlockcbad.AvailabilityError())
	}

	detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
		BaselineSize:         baselineSize,
		WindowSize:           windowSize,
		HopSize:              windowSize / 2,
		MaxCapacity:          baselineSize + 4*windowSize,
		PValueThreshold:      0.05,
		NCDThreshold:         0.3,
		PermutationCount:     100,
		Seed:                 42,
		CompressionAlgorithm: algo,
	})
	if err != nil {
		b.Fatalf("failed to create detector: %v", err)
	}
	defer detector.Close()

	normal := []byte(`{"ts":"2025-12-06T20:00:00Z","level":"info","msg":"User login successful","user":"alice"}`)
	anomaly := []byte(`{"ts":"2025-12-06T20:00:30Z","level":"error","msg":"SQL INJECTION: SELECT * FROM users; DROP TABLE sessions;--","severity":"critical"}`)

	events := make([][]byte, 0, totalEvents)
	for i := 0; i < totalEvents; i++ {
		events = append(events, normal)
	}
	// Inject anomaly
	if totalEvents > 0 {
		events[len(events)/2] = anomaly
	}

	// Prime baseline
	for i := 0; i < baselineSize && i < len(events); i++ {
		_, _ = detector.AddData(events[i])
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ev := range events {
			_, _ = detector.AddData(ev)
		}
		ready, _ := detector.IsReady()
		if ready {
			start := time.Now()
			_, _, _ = detector.DetectAnomaly()
			b.ReportMetric(float64(time.Since(start).Microseconds()), "detect_us")
		}
	}
}

func benchDetect(b *testing.B, totalEvents, baselineSize, windowSize int) {
	if !driftlockcbad.IsAvailable() {
		b.Skipf("CBAD core unavailable: %v", driftlockcbad.AvailabilityError())
	}

	detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
		BaselineSize:         baselineSize,
		WindowSize:           windowSize,
		HopSize:              windowSize / 2,
		MaxCapacity:          baselineSize + 4*windowSize,
		PValueThreshold:      0.05,
		NCDThreshold:         0.3,
		PermutationCount:     100,
		Seed:                 42,
		CompressionAlgorithm: "zstd",
	})
	if err != nil {
		b.Fatalf("failed to create detector: %v", err)
	}
	defer detector.Close()

	normal := []byte("INFO service=api msg=ok\n")
	anomaly := []byte("ERROR service=db msg=connection_failed stacktrace=AAAA\n")

	events := make([][]byte, 0, totalEvents)
	for i := 0; i < totalEvents; i++ {
		events = append(events, normal)
	}
	// Inject one anomaly near the end to keep detection path exercised
	if totalEvents > 0 {
		events[len(events)/2] = anomaly
	}

	// Prime baseline
	for i := 0; i < baselineSize && i < len(events); i++ {
		_, _ = detector.AddData(events[i])
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reuse the detector; feed a window's worth of events each iteration
		for _, ev := range events {
			_, _ = detector.AddData(ev)
		}
		ready, _ := detector.IsReady()
		if ready {
			start := time.Now()
			_, _, _ = detector.DetectAnomaly()
			b.ReportMetric(float64(time.Since(start).Microseconds()), "detect_us")
		}
	}
}
