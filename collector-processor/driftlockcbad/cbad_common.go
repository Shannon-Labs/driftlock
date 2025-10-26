package driftlockcbad

import "fmt"

// Metrics represents the anomaly detection results from CBAD.
type Metrics struct {
	NCD                      float64
	PValue                   float64
	BaselineCompressionRatio float64
	WindowCompressionRatio   float64
	BaselineEntropy          float64
	WindowEntropy            float64
	IsAnomaly                bool
	ConfidenceLevel          float64
}

const (
	defaultSeed         = 42
	defaultPermutations = 1000
)

// GetAnomalyExplanation generates a human-readable explanation of the anomaly detection result.
func (m *Metrics) GetAnomalyExplanation() string {
	if !m.IsAnomaly {
		return fmt.Sprintf(
			"No anomaly detected: NCD=%.3f, p=%.3f, compression ratios similar (baseline=%.2fx, window=%.2fx)",
			m.NCD,
			m.PValue,
			m.BaselineCompressionRatio,
			m.WindowCompressionRatio,
		)
	}

	return fmt.Sprintf(
		"Anomaly detected: NCD=%.3f, p=%.3f, compression ratio dropped from %.2fx to %.2fx due to data structure changes",
		m.NCD,
		m.PValue,
		m.BaselineCompressionRatio,
		m.WindowCompressionRatio,
	)
}

// GetDetailedExplanation provides a comprehensive explanation with statistical significance.
func (m *Metrics) GetDetailedExplanation() string {
	significance := "not statistically significant"
	if m.PValue < 0.05 {
		significance = "statistically significant"
	}

	return fmt.Sprintf(
		"CBAD Analysis: NCD=%.3f (%s), confidence=%.1f%%, baseline entropy=%.2f bits/byte, window entropy=%.2f bits/byte, compression change=%.1f%%",
		m.NCD,
		significance,
		m.ConfidenceLevel*100,
		m.BaselineEntropy,
		m.WindowEntropy,
		((m.WindowCompressionRatio-m.BaselineCompressionRatio)/m.BaselineCompressionRatio)*100,
	)
}

// IsStatisticallySignificant returns true if the anomaly is statistically significant (p < 0.05).
func (m *Metrics) IsStatisticallySignificant() bool {
	return m.PValue < 0.05
}

// GetConfidenceLevel returns the confidence level as a percentage (0-100).
func (m *Metrics) GetConfidenceLevel() float64 {
	return m.ConfidenceLevel * 100
}
