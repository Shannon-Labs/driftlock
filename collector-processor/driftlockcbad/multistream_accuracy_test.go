package driftlockcbad

import "testing"

// Ensures baselines per stream do not interfere.
func TestMultipleStreamsIndependently(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	streams := []struct {
		name         string
		normalEvents []string
		anomaly      string
		expectAnom   bool
	}{
		{
			name: "auth",
			normalEvents: []string{
				"AUTH user=alice status=ok",
				"AUTH user=bob status=ok",
				"AUTH user=carol status=ok",
			},
			anomaly:    "AUTH user=mallory status=fail reason=locked_out",
			expectAnom: true,
		},
		{
			name: "payments",
			normalEvents: []string{
				"PAY charge=42 status=ok",
				"PAY charge=99 status=ok",
				"PAY charge=18 status=ok",
			},
			anomaly:    "PAY charge=5000 status=declined reason=fraud",
			expectAnom: true,
		},
		{
			name: "users",
			normalEvents: []string{
				"USER profile update ok",
				"USER profile update ok",
				"USER profile update ok",
			},
			anomaly:    "USER profile update ok",
			expectAnom: false,
		},
	}

	for _, s := range streams {
		t.Run(s.name, func(t *testing.T) {
			detector, err := NewDetector(DetectorConfig{
				BaselineSize:                   6,
				WindowSize:                     3,
				HopSize:                        2,
				MaxCapacity:                    64,
				PValueThreshold:                0.6,
				NCDThreshold:                   0.3,
				CompressionRatioDropThreshold:  0.2,
				EntropyChangeThreshold:         0.1,
				CompositeThreshold:             0.6,
				PermutationCount:               20,
				CompressionAlgorithm:           "zstd",
				RequireStatisticalSignificance: true,
			})
			if err != nil {
				t.Fatalf("create detector: %v", err)
			}
			defer detector.Close()

			for i := 0; i < detector.config.BaselineSize; i++ {
				_, _ = detector.AddData([]byte(s.normalEvents[i%len(s.normalEvents)]))
			}
			for i := 0; i < detector.config.WindowSize; i++ {
				ev := s.anomaly
				if i%2 == 0 && !s.expectAnom {
					ev = s.normalEvents[0]
				}
				_, _ = detector.AddData([]byte(ev))
			}

			ready, _ := detector.IsReady()
			if !ready {
				t.Fatalf("detector not ready")
			}
			detected, metrics, err := detector.DetectAnomaly()
			if err != nil {
				t.Fatalf("detect: %v", err)
			}
			if metrics != nil {
				t.Logf("[%s] detected=%v ncd=%.3f p=%.3f compression_change=%.3f entropy_change=%.3f",
					s.name, detected, metrics.NCD, metrics.PValue, metrics.CompressionRatioChange, metrics.EntropyChange)
			}
			if detected != s.expectAnom {
				t.Fatalf("detected=%v want %v", detected, s.expectAnom)
			}
		})
	}
}
