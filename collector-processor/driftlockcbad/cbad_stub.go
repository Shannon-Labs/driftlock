//go:build !cgo || driftlock_no_cbad

package driftlockcbad

import "fmt"

// ComputeMetrics is a stub implementation used when the CBAD Rust library is unavailable.
func ComputeMetrics(baseline []byte, window []byte, seed uint64, permutations int) (*Metrics, error) {
	return nil, fmt.Errorf("cbad: Rust core not available (enable CGO and omit the driftlock_no_cbad build tag)")
}

// ComputeMetricsQuick is a stub implementation used when the CBAD Rust library is unavailable.
func ComputeMetricsQuick(baseline []byte, window []byte) (*Metrics, error) {
	return ComputeMetrics(baseline, window, defaultSeed, defaultPermutations)
}

// ValidateLibrary verifies that the Rust library is linked; in stub builds it reports an error.
func ValidateLibrary() error {
	return fmt.Errorf("cbad: ValidateLibrary requires CGO and the CBAD library (remove driftlock_no_cbad if set)")
}

// HasOpenZL always reports false when the Rust core is unavailable.
func HasOpenZL() bool {
	return false
}
