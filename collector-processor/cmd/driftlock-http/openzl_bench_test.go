package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

// BenchmarkOpenZLPreference runs only when OpenZL symbols are present.
// It compares openzl vs zstd over a deterministic payload to ensure
// we don't regress below the generic compressor performance envelope.
func BenchmarkOpenZLPreference(b *testing.B) {
	if !driftlockcbad.HasOpenZL() {
		b.Skip("OpenZL not compiled in")
	}

	events := buildBenchmarkEvents(256)

	bench := func(name, algo string) {
		b.Run(name, func(sb *testing.B) {
			cfg := driftlockcbad.DetectorConfig{
				BaselineSize:         64,
				WindowSize:           16,
				HopSize:              4,
				MaxCapacity:          64 + 4*16 + 1024,
				PValueThreshold:      0.05,
				NCDThreshold:         0.2,
				PermutationCount:     64,
				Seed:                 42,
				CompressionAlgorithm: algo,
			}
			detector, err := driftlockcbad.NewDetector(cfg)
			if err != nil {
				b.Fatalf("new detector: %v", err)
			}
			defer detector.Close()

			b.ResetTimer()
			for i := 0; i < sb.N; i++ {
				for _, ev := range events {
					added, err := detector.AddData(ev)
					if err != nil {
						b.Fatalf("add data: %v", err)
					}
					if !added {
						continue
					}
					ready, err := detector.IsReady()
					if err != nil {
						b.Fatalf("ready: %v", err)
					}
					if !ready {
						continue
					}
					_, _, err = detector.DetectAnomaly()
					if err != nil {
						b.Fatalf("detect: %v", err)
					}
				}
			}
		})
	}

	bench("openzl", "openzl")
	bench("zstd", "zstd")
}

func buildBenchmarkEvents(n int) []json.RawMessage {
	out := make([]json.RawMessage, 0, n)
	for i := 0; i < n; i++ {
		payload := map[string]any{
			"ts":    time.Unix(0, int64(i*1_000_000)).UTC().Format(time.RFC3339Nano),
			"value": i % 7,
			"msg":   "bench-payload",
		}
		buf, _ := json.Marshal(payload)
		out = append(out, buf)
	}
	return out
}
