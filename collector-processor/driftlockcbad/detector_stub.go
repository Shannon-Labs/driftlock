//go:build !cgo || !driftlock_cbad_cgo

package driftlockcbad

import "fmt"

// Detector is a stub implementation used when the Rust detector is unavailable.
type Detector struct{}

// NewDetector returns an error because the detector requires the Rust core.
func NewDetector(config DetectorConfig) (*Detector, error) {
	_ = config
	return nil, fmt.Errorf("cbad: detector requires CGO-enabled build with driftlock_cbad_cgo tag")
}

// Close is a no-op for the stub implementation.
func (d *Detector) Close() {
	_ = d
}

// AddData returns an error indicating the detector is unavailable.
func (d *Detector) AddData(data []byte) (bool, error) {
	_ = d
	return false, fmt.Errorf("cbad: detector requires CGO-enabled build with driftlock_cbad_cgo tag (received %d bytes)", len(data))
}

// IsReady always reports false in stub builds.
func (d *Detector) IsReady() (bool, error) {
	_ = d
	return false, fmt.Errorf("cbad: detector requires CGO-enabled build with driftlock_cbad_cgo tag")
}

// DetectAnomaly returns an error because detection is unavailable.
func (d *Detector) DetectAnomaly() (*EnhancedMetrics, error) {
	_ = d
	return nil, fmt.Errorf("cbad: detector requires CGO-enabled build with driftlock_cbad_cgo tag")
}

// GetStats returns an error because no detector state is available.
func (d *Detector) GetStats() (map[string]any, error) {
	_ = d
	return nil, fmt.Errorf("cbad: detector requires CGO-enabled build with driftlock_cbad_cgo tag")
}
