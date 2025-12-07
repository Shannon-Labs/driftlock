package driftlockcbad

import "testing"

func TestAccuracyOnKnownAnomalies(t *testing.T) {
	if !IsAvailable() {
		t.Skipf("CBAD core unavailable: %v", AvailabilityError())
	}

	testCases := []struct {
		name         string
		normalEvents []string
		anomalyEvent string
		shouldDetect bool
	}{
		{
			name: "database_failure_in_logs",
			normalEvents: []string{
				"GET /api/users 200 12ms",
				"POST /api/login 200 45ms",
				"GET /api/users 200 10ms",
				"GET /api/users 200 11ms",
			},
			anomalyEvent: "ERROR: Database connection pool exhausted",
			shouldDetect: true,
		},
		{
			name: "steady_state_no_anomaly",
			normalEvents: []string{
				"User login successful",
				"User login successful",
				"User login successful",
				"User login successful",
			},
			anomalyEvent: "User login successful",
			shouldDetect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ncdThreshold := 0.5
			dropThreshold := 0.8
			composite := 0.6
			pThreshold := 0.2
			if tc.shouldDetect {
				ncdThreshold = 0.2
				dropThreshold = 0.1
				composite = 0.3
				pThreshold = 1.0
			}

			detector, err := NewDetector(DetectorConfig{
				BaselineSize:                   8,
				WindowSize:                     4,
				HopSize:                        2,
				MaxCapacity:                    64,
				PValueThreshold:                pThreshold,
				NCDThreshold:                   ncdThreshold,
				CompressionRatioDropThreshold:  dropThreshold,
				EntropyChangeThreshold:         0.0,
				CompositeThreshold:             composite,
				PermutationCount:               50,
				CompressionAlgorithm:           "zstd",
				RequireStatisticalSignificance: false,
			})
			if err != nil {
				t.Fatalf("create detector: %v", err)
			}
			defer detector.Close()

			// Fill baseline with normal events (repeat to reach BaselineSize)
			for i := 0; i < detector.config.BaselineSize; i++ {
				ev := tc.normalEvents[i%len(tc.normalEvents)]
				if _, err := detector.AddData([]byte(ev)); err != nil {
					t.Fatalf("add baseline: %v", err)
				}
			}

			// Fill window with anomaly/normal mixture
			for i := 0; i < detector.config.WindowSize; i++ {
				ev := tc.anomalyEvent
				if i%2 == 1 && !tc.shouldDetect {
					ev = tc.normalEvents[0]
				}
				if _, err := detector.AddData([]byte(ev)); err != nil {
					t.Fatalf("add window: %v", err)
				}
			}

			ready, err := detector.IsReady()
			if err != nil {
				t.Fatalf("isReady: %v", err)
			}
			if !ready {
				t.Fatalf("detector not ready after feeding events")
			}

			detected, metrics, err := detector.DetectAnomaly()
			if err != nil {
				t.Fatalf("detect anomaly: %v", err)
			}

			if detected != tc.shouldDetect {
				t.Fatalf("detected=%v want %v (ncd=%.3f p=%.3f compression_change=%.3f)", detected, tc.shouldDetect, metrics.NCD, metrics.PValue, metrics.CompressionRatioChange)
			}
		})
	}
}
