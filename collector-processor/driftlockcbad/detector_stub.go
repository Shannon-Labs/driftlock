//go:build !cgo || driftlock_no_cbad

package driftlockcbad

import "fmt"

// Detector is a stub implementation used when the Rust detector is unavailable.
type Detector struct{}

// NewDetector returns an error because the detector requires the Rust core.
func NewDetector(config DetectorConfig) (*Detector, error) {
	_ = config
	return nil, fmt.Errorf("cbad: detector requires CGO (remove the driftlock_no_cbad build tag)")
}

// Close is a no-op for the stub implementation.
func (d *Detector) Close() {
	_ = d
}

// AddData returns an error indicating the detector is unavailable.
func (d *Detector) AddData(data []byte) (bool, error) {
	_ = d
	return false, fmt.Errorf("cbad: detector requires CGO (remove the driftlock_no_cbad build tag); received %d bytes", len(data))
}

// IsReady always reports false in stub builds.
func (d *Detector) IsReady() (bool, error) {
	_ = d
	return false, fmt.Errorf("cbad: detector requires CGO (remove the driftlock_no_cbad build tag)")
}

// DetectAnomaly returns an error because detection is unavailable.
func (d *Detector) DetectAnomaly() (bool, *EnhancedMetrics, error) {
	_ = d
	return false, nil, fmt.Errorf("cbad: detector requires CGO (remove the driftlock_no_cbad build tag)")
}

// GetStats returns an error because no detector state is available.
func (d *Detector) GetStats() (uint64, int, bool, error) {
	_ = d
	return 0, 0, false, fmt.Errorf("cbad: detector requires CGO (remove the driftlock_no_cbad build tag)")
}
