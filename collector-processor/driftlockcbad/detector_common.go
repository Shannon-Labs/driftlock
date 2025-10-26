package driftlockcbad

// DetectorConfig represents configuration for the anomaly detector.
type DetectorConfig struct {
	BaselineSize                   int
	WindowSize                     int
	HopSize                        int
	MaxCapacity                    int
	PValueThreshold                float64
	NCDThreshold                   float64
	PermutationCount               int
	Seed                           uint64
	RequireStatisticalSignificance bool
	CompressionAlgorithm           string
}

// EnhancedMetrics represents comprehensive anomaly detection results.
type EnhancedMetrics struct {
	NCD                        float64
	PValue                     float64
	BaselineCompressionRatio   float64
	WindowCompressionRatio     float64
	BaselineEntropy            float64
	WindowEntropy              float64
	IsAnomaly                  bool
	ConfidenceLevel            float64
	IsStatisticallySignificant bool
	CompressionRatioChange     float64
	EntropyChange              float64
	Explanation                string
}
