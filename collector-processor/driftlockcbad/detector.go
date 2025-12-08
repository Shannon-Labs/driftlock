//go:build cgo && !driftlock_no_cbad

package driftlockcbad

/*
#cgo CFLAGS: -I../../cbad-core/src
#cgo LDFLAGS: -L../../cbad-core/target/release -lcbad_core
#include <stdint.h>
#include <stdlib.h>

// C-compatible configuration structure
typedef struct {
    size_t baseline_size;
    size_t window_size;
    size_t hop_size;
    size_t max_capacity;
    double p_value_threshold;
    double ncd_threshold;
    double compression_ratio_drop_threshold;
    double entropy_change_threshold;
    double composite_threshold;
    size_t permutation_count;
    uint64_t seed;
    int require_statistical_significance; // 0 = false, 1 = true
    const char* compression_algorithm; // "zlab", "zstd", "lz4", "gzip", "openzl"
} CBADConfig;

// Enhanced metrics structure with additional fields
typedef struct {
    double ncd;
    double p_value;
    double baseline_compression_ratio;
    double window_compression_ratio;
    double baseline_entropy;
    double window_entropy;
    int is_anomaly;
    double confidence_level;
    int is_statistically_significant;
    double compression_ratio_change;
    double entropy_change;
    const char* explanation; // Owned by Rust, must be freed
    double recommended_ncd_threshold;
    size_t recommended_window_size;
    double data_stability_score;
} CBADEnhancedMetrics;

// Opaque handle for AnomalyDetector instances
typedef void* CBADDetectorHandle;

// Enhanced FFI function declarations
extern CBADDetectorHandle cbad_detector_create(const CBADConfig* config_ptr);
extern void cbad_detector_destroy(CBADDetectorHandle handle);
extern int cbad_detector_add_data(CBADDetectorHandle handle, const uint8_t* data_ptr, size_t data_len);
extern int cbad_detector_is_ready(CBADDetectorHandle handle);
extern int cbad_detector_detect_anomaly(CBADDetectorHandle handle, CBADEnhancedMetrics* metrics_ptr);
extern void cbad_free_explanation(char* explanation_ptr);
extern int cbad_detector_get_stats(CBADDetectorHandle handle, size_t* stats_ptr);
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// Error codes from CBAD FFI
const (
	cbadErrPanic = -99 // Rust panic was caught at FFI boundary
)

// ErrCBADPanic indicates a Rust panic was caught in the CBAD library.
// This is a critical error indicating a bug in the Rust code.
var ErrCBADPanic = errors.New("CBAD internal panic - Rust panic caught at FFI boundary")

// Detector represents a streaming anomaly detector
type Detector struct {
	handle C.CBADDetectorHandle
	config DetectorConfig
}

// NewDetector creates a new anomaly detector with the given configuration
func NewDetector(config DetectorConfig) (*Detector, error) {
	// Set defaults
	if config.BaselineSize <= 0 {
		config.BaselineSize = 1000
	}
	if config.WindowSize <= 0 {
		config.WindowSize = 100
	}
	if config.HopSize <= 0 {
		config.HopSize = 50
	}
	if config.MaxCapacity <= 0 {
		config.MaxCapacity = 10000
	}
	if config.PValueThreshold <= 0 || config.PValueThreshold >= 1 {
		config.PValueThreshold = 0.05
	}
	if config.NCDThreshold <= 0 {
		config.NCDThreshold = 0.3
	}
	if config.CompressionRatioDropThreshold <= 0 {
		config.CompressionRatioDropThreshold = 0.15
	}
	if config.EntropyChangeThreshold <= 0 {
		config.EntropyChangeThreshold = 0.2
	}
	if config.CompositeThreshold <= 0 {
		config.CompositeThreshold = 0.6
	}
	if config.PermutationCount <= 0 {
		config.PermutationCount = 1000
	}
	if config.Seed == 0 {
		config.Seed = 42
	}
	if config.CompressionAlgorithm == "" {
		config.CompressionAlgorithm = "zstd"
	}

	// Create C configuration
	cConfig := C.CBADConfig{
		baseline_size:                    C.size_t(config.BaselineSize),
		window_size:                      C.size_t(config.WindowSize),
		hop_size:                         C.size_t(config.HopSize),
		max_capacity:                     C.size_t(config.MaxCapacity),
		p_value_threshold:                C.double(config.PValueThreshold),
		ncd_threshold:                    C.double(config.NCDThreshold),
		compression_ratio_drop_threshold: C.double(config.CompressionRatioDropThreshold),
		entropy_change_threshold:         C.double(config.EntropyChangeThreshold),
		composite_threshold:              C.double(config.CompositeThreshold),
		permutation_count:                C.size_t(config.PermutationCount),
		seed:                             C.uint64_t(config.Seed),
		require_statistical_significance: C.int(0),
	}

	if config.RequireStatisticalSignificance {
		cConfig.require_statistical_significance = C.int(1)
	}

	// Set compression algorithm
	cAlgo := C.CString(config.CompressionAlgorithm)
	defer C.free(unsafe.Pointer(cAlgo))
	cConfig.compression_algorithm = (*C.char)(cAlgo)

	// Create detector
	handle := C.cbad_detector_create(&cConfig)
	if handle == nil {
		return nil, errors.New("failed to create CBAD detector")
	}

	return &Detector{
		handle: handle,
		config: config,
	}, nil
}

// Close destroys the detector and frees its memory
func (d *Detector) Close() {
	if d.handle != nil {
		C.cbad_detector_destroy(d.handle)
		d.handle = nil
	}
}

// AddData adds data to the anomaly detector
func (d *Detector) AddData(data []byte) (bool, error) {
	if d.handle == nil {
		return false, errors.New("detector is closed")
	}
	if len(data) == 0 {
		return false, errors.New("data must not be empty")
	}

	result := C.cbad_detector_add_data(d.handle, (*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))

	switch result {
	case 1:
		return true, nil // Data added successfully
	case 0:
		return false, nil // Data added but dropped due to privacy compliance
	case -1:
		return false, errors.New("invalid parameters")
	case -2:
		return false, errors.New("internal error")
	case cbadErrPanic:
		return false, ErrCBADPanic
	default:
		return false, fmt.Errorf("unknown error code: %d", result)
	}
}

// IsReady checks if the detector has enough data for analysis
func (d *Detector) IsReady() (bool, error) {
	if d.handle == nil {
		return false, errors.New("detector is closed")
	}

	result := C.cbad_detector_is_ready(d.handle)

	switch result {
	case 1:
		return true, nil // Ready for analysis
	case 0:
		return false, nil // Not ready yet
	case -1:
		return false, errors.New("invalid detector handle")
	case -2:
		return false, errors.New("internal error")
	case cbadErrPanic:
		return false, ErrCBADPanic
	default:
		return false, fmt.Errorf("unknown error code: %d", result)
	}
}

// DetectAnomaly performs anomaly detection and returns results
func (d *Detector) DetectAnomaly() (bool, *EnhancedMetrics, error) {
	if d.handle == nil {
		return false, nil, errors.New("detector is closed")
	}

	var cMetrics C.CBADEnhancedMetrics
	result := C.cbad_detector_detect_anomaly(d.handle, &cMetrics)

	var metrics *EnhancedMetrics

	switch result {
	case 1: // Anomaly detected
		metrics = &EnhancedMetrics{
			NCD:                        float64(cMetrics.ncd),
			PValue:                     float64(cMetrics.p_value),
			BaselineCompressionRatio:   float64(cMetrics.baseline_compression_ratio),
			WindowCompressionRatio:     float64(cMetrics.window_compression_ratio),
			BaselineEntropy:            float64(cMetrics.baseline_entropy),
			WindowEntropy:              float64(cMetrics.window_entropy),
			IsAnomaly:                  cMetrics.is_anomaly != 0,
			ConfidenceLevel:            float64(cMetrics.confidence_level),
			IsStatisticallySignificant: cMetrics.is_statistically_significant != 0,
			CompressionRatioChange:     float64(cMetrics.compression_ratio_change),
			EntropyChange:              float64(cMetrics.entropy_change),
			RecommendedNCDThreshold:    float64(cMetrics.recommended_ncd_threshold),
			RecommendedWindowSize:      int(cMetrics.recommended_window_size),
			DataStabilityScore:         float64(cMetrics.data_stability_score),
		}

		// Get the explanation string (must be freed)
		if cMetrics.explanation != nil {
			cExplanation := (*C.char)(unsafe.Pointer(cMetrics.explanation))
			metrics.Explanation = C.GoString(cExplanation)
			C.cbad_free_explanation((*C.char)(cMetrics.explanation))
		}

		return true, metrics, nil

	case 0: // No anomaly detected
		metrics = &EnhancedMetrics{
			NCD:                        float64(cMetrics.ncd),
			PValue:                     float64(cMetrics.p_value),
			BaselineCompressionRatio:   float64(cMetrics.baseline_compression_ratio),
			WindowCompressionRatio:     float64(cMetrics.window_compression_ratio),
			BaselineEntropy:            float64(cMetrics.baseline_entropy),
			WindowEntropy:              float64(cMetrics.window_entropy),
			IsAnomaly:                  false,
			ConfidenceLevel:            float64(cMetrics.confidence_level),
			IsStatisticallySignificant: cMetrics.is_statistically_significant != 0,
			CompressionRatioChange:     float64(cMetrics.compression_ratio_change),
			EntropyChange:              float64(cMetrics.entropy_change),
			RecommendedNCDThreshold:    float64(cMetrics.recommended_ncd_threshold),
			RecommendedWindowSize:      int(cMetrics.recommended_window_size),
			DataStabilityScore:         float64(cMetrics.data_stability_score),
		}

		// Get the explanation string (must be freed)
		if cMetrics.explanation != nil {
			cExplanation := (*C.char)(unsafe.Pointer(cMetrics.explanation))
			metrics.Explanation = C.GoString(cExplanation)
			C.cbad_free_explanation((*C.char)(cMetrics.explanation))
		}

		return false, metrics, nil

	case -1: // Not enough data
		return false, nil, errors.New("not enough data for analysis")
	case -2: // Internal error
		return false, nil, errors.New("internal error during anomaly detection")
	case cbadErrPanic: // Rust panic caught
		return false, nil, ErrCBADPanic
	default:
		return false, nil, fmt.Errorf("unknown error code: %d", result)
	}
}

// GetStats returns current detector statistics
func (d *Detector) GetStats() (totalEvents uint64, memoryUsage int, isReady bool, err error) {
	if d.handle == nil {
		return 0, 0, false, errors.New("detector is closed")
	}

	var stats [3]C.size_t
	result := C.cbad_detector_get_stats(d.handle, &stats[0])

	switch result {
	case 0: // Success
		return uint64(stats[0]), int(stats[1]), stats[2] != 0, nil
	case -1:
		return 0, 0, false, errors.New("invalid parameters")
	case -2:
		return 0, 0, false, errors.New("internal error")
	case cbadErrPanic:
		return 0, 0, false, ErrCBADPanic
	default:
		return 0, 0, false, fmt.Errorf("unknown error code: %d", result)
	}
}

// GetAnomalyExplanation generates a human-readable explanation of the anomaly detection result
func (m *EnhancedMetrics) GetAnomalyExplanation() string {
	if !m.IsAnomaly {
		return fmt.Sprintf("No anomaly detected: NCD=%.3f, p=%.3f, compression ratios similar (baseline=%.2fx, window=%.2fx)",
			m.NCD, m.PValue, m.BaselineCompressionRatio, m.WindowCompressionRatio)
	}

	return fmt.Sprintf("Anomaly detected: NCD=%.3f, p=%.3f, compression ratio changed by %.1f%% due to pattern changes",
		m.NCD, m.PValue, m.CompressionRatioChange*100)
}

// GetDetailedExplanation provides a comprehensive explanation with statistical significance
func (m *EnhancedMetrics) GetDetailedExplanation() string {
	significance := "not statistically significant"
	if m.IsStatisticallySignificant {
		significance = "statistically significant"
	}

	return fmt.Sprintf("CBAD Analysis: NCD=%.3f (%s), confidence=%.1f%%, entropy change=%+.1f%%, %s",
		m.NCD, significance, m.ConfidenceLevel*100, m.EntropyChange*100, m.Explanation)
}
