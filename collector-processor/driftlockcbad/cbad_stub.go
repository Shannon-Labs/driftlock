//go:build !cgo || !driftlock_cbad_cgo

package driftlockcbad

import "fmt"

// ComputeMetrics is a stub implementation used when the CBAD Rust library is unavailable.
func ComputeMetrics(baseline []byte, window []byte, seed uint64, permutations int) (*Metrics, error) {
	return nil, fmt.Errorf("cbad: Rust core not available (build with CGO enabled and tag driftlock_cbad_cgo)")
}

// ComputeMetricsQuick is a stub implementation used when the CBAD Rust library is unavailable.
func ComputeMetricsQuick(baseline []byte, window []byte) (*Metrics, error) {
	return ComputeMetrics(baseline, window, defaultSeed, defaultPermutations)
}

// ValidateLibrary verifies that the Rust library is linked; in stub builds it reports an error.
func ValidateLibrary() error {
	return fmt.Errorf("cbad: ValidateLibrary requires CGO-enabled build with driftlock_cbad_cgo tag")
}
