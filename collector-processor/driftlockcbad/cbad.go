//go:build cgo && !driftlock_no_cbad
// +build cgo,!driftlock_no_cbad

// Driftlock - Compression-Based Anomaly Detection
// Patent-Pending Technology - Shannon Labs, LLC
// Licensed under Apache 2.0. Commercial licenses available.
// See PATENTS.md and LICENSE for details.

package driftlockcbad

/*
#cgo CFLAGS: -I../../cbad-core/src
#cgo LDFLAGS: -L../../cbad-core/target/release -lcbad_core
#include <stdint.h>

// C-compatible metrics structure from Rust FFI
typedef struct {
    double ncd;
    double p_value;
    double baseline_compression_ratio;
    double window_compression_ratio;
    double baseline_entropy;
    double window_entropy;
    int is_anomaly; // 0 = false, 1 = true
    double confidence_level;
} CBADMetrics;

// Rust FFI function declarations
extern void cbad_init_logging();

extern CBADMetrics cbad_compute_metrics(
    const uint8_t* baseline_ptr,
    size_t baseline_len,
    const uint8_t* window_ptr,
    size_t window_len,
    uint64_t seed,
    size_t permutations
);

extern double cbad_compute_metrics_len(const uint8_t* data, size_t len);
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

func init() {
	// Initialize Rust logging
	os.Setenv("RUST_LOG", "error")
	C.cbad_init_logging()
}

// ComputeMetrics calculates CBAD anomaly detection metrics using compression-based analysis
func ComputeMetrics(baseline []byte, window []byte, seed uint64, permutations int) (*Metrics, error) {
	if len(baseline) == 0 || len(window) == 0 {
		return nil, errors.New("baseline and window data must not be empty")
	}

	if permutations <= 0 {
		permutations = 1000 // Default permutation count
	}

	// Call the Rust FFI function
	cMetrics := C.cbad_compute_metrics(
		(*C.uint8_t)(unsafe.Pointer(&baseline[0])),
		C.size_t(len(baseline)),
		(*C.uint8_t)(unsafe.Pointer(&window[0])),
		C.size_t(len(window)),
		C.uint64_t(seed),
		C.size_t(permutations),
	)

	// Convert C struct to Go struct
	metrics := &Metrics{
		NCD:                      float64(cMetrics.ncd),
		PValue:                   float64(cMetrics.p_value),
		BaselineCompressionRatio: float64(cMetrics.baseline_compression_ratio),
		WindowCompressionRatio:   float64(cMetrics.window_compression_ratio),
		BaselineEntropy:          float64(cMetrics.baseline_entropy),
		WindowEntropy:            float64(cMetrics.window_entropy),
		IsAnomaly:                cMetrics.is_anomaly != 0,
		ConfidenceLevel:          float64(cMetrics.confidence_level),
	}

	return metrics, nil
}

// ComputeMetricsQuick calculates CBAD metrics with default configuration
func ComputeMetricsQuick(baseline []byte, window []byte) (*Metrics, error) {
	return ComputeMetrics(baseline, window, defaultSeed, defaultPermutations)
}

// ValidateLibrary checks if the CBAD library is properly loaded and functional
func ValidateLibrary() error {
	// Test with simple data
	baseline := []byte("INFO service=api-gateway msg=request_completed\n")
	window := []byte("ERROR service=api-gateway msg=stack_trace\n")

	metrics, err := ComputeMetricsQuick(baseline, window)
	if err != nil {
		return fmt.Errorf("CBAD library validation failed: %w", err)
	}

	// Basic sanity checks
	if metrics.NCD < 0 || metrics.NCD > 1 {
		return fmt.Errorf("invalid NCD value: %f", metrics.NCD)
	}

	if metrics.PValue < 0 || metrics.PValue > 1 {
		return fmt.Errorf("invalid p-value: %f", metrics.PValue)
	}

	return nil
}
