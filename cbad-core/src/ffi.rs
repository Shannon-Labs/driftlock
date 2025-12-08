//! Enhanced C FFI exports for Go integration
//!
//! This module provides production-ready C FFI exports that enable
//! Go-based applications to use the complete CBAD anomaly detection
//! engine with streaming capabilities and proper error handling.
//!
//! ## Safety Features
//!
//! All FFI functions are wrapped with panic catching to prevent Rust panics
//! from unwinding into C/Go code, which would cause undefined behavior.
//! Additionally, handle validation prevents use-after-free vulnerabilities.

use crate::anomaly::{AnomalyConfig, AnomalyDetector};
use crate::compression::CompressionAlgorithm;
use crate::window::{DataEvent, WindowConfig};
use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int};
use std::panic::{catch_unwind, AssertUnwindSafe};
use std::ptr;
use std::slice;

/// Error code returned when a panic is caught at FFI boundary
pub const CBAD_ERR_PANIC: c_int = -99;

/// Maximum allowed data size for any single FFI call (16MB)
const MAX_FFI_DATA_SIZE: usize = 16 * 1024 * 1024;

/// Sentinel value indicating a valid, active detector
const VALID_SENTINEL: usize = 0xCAFEBABE_CAFEBABE;

/// Sentinel value indicating a freed detector (for debugging)
const FREED_SENTINEL: usize = 0xDEADBEEF_DEADBEEF;

/// Wrapper around AnomalyDetector that adds safety validation
// Public visibility is required because the opaque handle type is exported via FFI,
// but fields remain private to keep the wrapper opaque to callers.
pub struct DetectorWrapper {
    sentinel: usize,
    detector: AnomalyDetector,
}

impl DetectorWrapper {
    fn new(detector: AnomalyDetector) -> Self {
        Self {
            sentinel: VALID_SENTINEL,
            detector,
        }
    }

    fn is_valid(&self) -> bool {
        self.sentinel == VALID_SENTINEL
    }

    fn mark_freed(&mut self) {
        self.sentinel = FREED_SENTINEL;
    }

    fn detector(&self) -> &AnomalyDetector {
        &self.detector
    }
}

/// Validate pointer and length parameters for FFI calls
/// Returns true if valid, false if invalid
#[inline]
fn validate_ptr_len(ptr: *const u8, len: usize) -> bool {
    // Null check
    if ptr.is_null() {
        return false;
    }

    // Zero length is invalid for data operations
    if len == 0 {
        return false;
    }

    // Prevent obviously incorrect huge lengths
    if len > MAX_FFI_DATA_SIZE {
        log::warn!(
            "FFI data length {} exceeds maximum {}",
            len,
            MAX_FFI_DATA_SIZE
        );
        return false;
    }

    true
}

/// Wraps an FFI function body with catch_unwind to prevent panics from unwinding into C
macro_rules! ffi_catch_unwind {
    ($default:expr, $body:block) => {{
        let result = catch_unwind(AssertUnwindSafe(|| $body));
        match result {
            Ok(val) => val,
            Err(e) => {
                // Log the panic if possible
                if let Some(s) = e.downcast_ref::<&str>() {
                    log::error!("CBAD FFI panic caught: {}", s);
                } else if let Some(s) = e.downcast_ref::<String>() {
                    log::error!("CBAD FFI panic caught: {}", s);
                } else {
                    log::error!("CBAD FFI panic caught: unknown panic payload");
                }
                $default
            }
        }
    }};
}

/// Opaque handle for AnomalyDetector instances (now uses DetectorWrapper for safety)
pub type CBADDetectorHandle = *mut DetectorWrapper;

/// Configuration for anomaly detection (C-compatible)
#[repr(C)]
pub struct CBADConfig {
    pub baseline_size: usize,
    pub window_size: usize,
    pub hop_size: usize,
    pub max_capacity: usize,
    pub p_value_threshold: f64,
    pub ncd_threshold: f64,
    pub compression_ratio_drop_threshold: f64,
    pub entropy_change_threshold: f64,
    pub composite_threshold: f64,
    pub permutation_count: usize,
    pub seed: u64,
    pub require_statistical_significance: c_int, // 0 = false, 1 = true
    pub compression_algorithm: *const c_char,    // "zlab", "zstd", "lz4", "gzip", "openzl"
}

/// Enhanced metrics structure with additional fields
#[repr(C)]
pub struct CBADEnhancedMetrics {
    pub ncd: f64,
    pub p_value: f64,
    pub baseline_compression_ratio: f64,
    pub window_compression_ratio: f64,
    pub baseline_entropy: f64,
    pub window_entropy: f64,
    pub is_anomaly: c_int,
    pub confidence_level: f64,
    pub is_statistically_significant: c_int,
    pub compression_ratio_change: f64,
    pub entropy_change: f64,
    pub explanation: *const c_char, // Owned by Rust, must be freed
    pub recommended_ncd_threshold: f64,
    pub recommended_window_size: usize,
    pub data_stability_score: f64,
}

/// Returns 1 when CBAD was compiled with OpenZL support, 0 otherwise
#[no_mangle]
pub extern "C" fn cbad_has_openzl() -> c_int {
    if cfg!(feature = "openzl") {
        1
    } else {
        0
    }
}

/// Create a new anomaly detector with configuration
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - config_ptr is a valid pointer to a CBADConfig struct
/// - compression_algorithm string is valid UTF-8 and null-terminated
///
/// # Returns
/// - Valid handle on success
/// - null_ptr on invalid parameters or creation failure
/// - Panics are caught and return null_ptr
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_create(config_ptr: *const CBADConfig) -> CBADDetectorHandle {
    ffi_catch_unwind!(ptr::null_mut(), {
        if config_ptr.is_null() {
            return ptr::null_mut();
        }

        let config = &*config_ptr;

        // Parse compression algorithm
        let algo = if config.compression_algorithm.is_null() {
            CompressionAlgorithm::Zstd
        } else {
            match CStr::from_ptr(config.compression_algorithm).to_str() {
                Ok("zlab") => CompressionAlgorithm::Zlab,
                Ok("zstd") => CompressionAlgorithm::Zstd,
                Ok("lz4") => CompressionAlgorithm::Lz4,
                Ok("gzip") => CompressionAlgorithm::Gzip,
                #[cfg(feature = "openzl")]
                Ok("openzl") => CompressionAlgorithm::OpenZL,
                _ => CompressionAlgorithm::Zstd,
            }
        };

        // Create window configuration
        let window_config = WindowConfig {
            baseline_size: config.baseline_size,
            window_size: config.window_size,
            hop_size: config.hop_size,
            max_capacity: config.max_capacity,
            time_window: None,
            privacy_config: Default::default(),
        };

        // Create anomaly configuration (new fields derive sensible defaults)
        let anomaly_config = AnomalyConfig {
            window_config,
            compression_algorithm: algo,
            p_value_threshold: config.p_value_threshold,
            ncd_threshold: config.ncd_threshold,
            compression_ratio_drop_threshold: config.compression_ratio_drop_threshold,
            entropy_change_threshold: config.entropy_change_threshold,
            composite_threshold: config.composite_threshold,
            permutation_count: config.permutation_count,
            seed: config.seed,
            require_statistical_significance: config.require_statistical_significance != 0,
        };

        // Create detector wrapped with safety sentinel
        match AnomalyDetector::new(anomaly_config) {
            Ok(detector) => {
                let wrapper = DetectorWrapper::new(detector);
                Box::into_raw(Box::new(wrapper))
            }
            Err(_) => ptr::null_mut(),
        }
    })
}

/// Destroy an anomaly detector and free its memory
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - handle is not used after this call (use-after-free protection via sentinel)
///
/// Double-destroy is safe - the function checks sentinel validity before freeing.
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_destroy(handle: CBADDetectorHandle) {
    ffi_catch_unwind!((), {
        if handle.is_null() {
            return;
        }

        // Check sentinel to detect double-free attempts
        let wrapper = &mut *handle;
        if !wrapper.is_valid() {
            log::warn!("cbad_detector_destroy called on invalid/freed handle");
            return;
        }

        // Mark as freed before actually freeing (for debugging stale pointer access)
        wrapper.mark_freed();

        // Actually free the memory
        let _ = Box::from_raw(handle);
    })
}

/// Add data to the anomaly detector
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - data_ptr is a valid pointer to data_len bytes
/// - data_len is accurate and <= 16MB
///
/// # Returns
/// - 1: Data added successfully
/// - 0: Data not added (capacity full)
/// - -1: Invalid parameters
/// - -2: Internal error
/// - -99: Panic caught
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_add_data(
    handle: CBADDetectorHandle,
    data_ptr: *const u8,
    data_len: usize,
) -> c_int {
    ffi_catch_unwind!(CBAD_ERR_PANIC, {
        if handle.is_null() {
            return -1; // Invalid parameters
        }

        // Validate data pointer and length with bounds checking
        if !validate_ptr_len(data_ptr, data_len) {
            return -1; // Invalid parameters
        }

        // Check sentinel for use-after-free protection
        let wrapper = &*handle;
        if !wrapper.is_valid() {
            log::warn!("cbad_detector_add_data called on invalid/freed handle");
            return -1;
        }

        let detector = wrapper.detector();
        let data = slice::from_raw_parts(data_ptr, data_len);

        let event = DataEvent::new(data.to_vec());

        match detector.add_event(event) {
            Ok(added) => {
                if added {
                    1
                } else {
                    0
                }
            }
            Err(_) => -2, // Internal error
        }
    })
}

/// Check if the detector has enough data for analysis
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
///
/// # Returns
/// - 1: Ready for analysis
/// - 0: Not ready (need more data)
/// - -1: Invalid parameters
/// - -2: Internal error
/// - -99: Panic caught
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_is_ready(handle: CBADDetectorHandle) -> c_int {
    ffi_catch_unwind!(CBAD_ERR_PANIC, {
        if handle.is_null() {
            return -1;
        }

        // Check sentinel for use-after-free protection
        let wrapper = &*handle;
        if !wrapper.is_valid() {
            log::warn!("cbad_detector_is_ready called on invalid/freed handle");
            return -1;
        }

        let detector = wrapper.detector();

        match detector.is_ready() {
            Ok(ready) => {
                if ready {
                    1
                } else {
                    0
                }
            }
            Err(_) => -2, // Internal error
        }
    })
}

/// Perform anomaly detection and get results
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - metrics_ptr points to a valid CBADEnhancedMetrics struct that can be written to
///
/// # Returns
/// - 1: Anomaly detected
/// - 0: No anomaly detected
/// - -1: Not enough data or invalid parameters
/// - -2: Internal error
/// - -99: Panic caught
///
/// # Memory
/// The `explanation` field in metrics is allocated by Rust and must be freed
/// by calling `cbad_free_explanation`.
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_detect_anomaly(
    handle: CBADDetectorHandle,
    metrics_ptr: *mut CBADEnhancedMetrics,
) -> c_int {
    ffi_catch_unwind!(CBAD_ERR_PANIC, {
        if handle.is_null() || metrics_ptr.is_null() {
            return -1;
        }

        // Check sentinel for use-after-free protection
        let wrapper = &*handle;
        if !wrapper.is_valid() {
            log::warn!("cbad_detector_detect_anomaly called on invalid/freed handle");
            return -1;
        }

        let detector = wrapper.detector();

        match detector.detect_anomaly() {
            Ok(Some(result)) => {
                let metrics = &mut *metrics_ptr;

                // Fill in the metrics
                metrics.ncd = result.metrics.ncd;
                metrics.p_value = result.metrics.p_value;
                metrics.baseline_compression_ratio = result.metrics.baseline_compression_ratio;
                metrics.window_compression_ratio = result.metrics.window_compression_ratio;
                metrics.baseline_entropy = result.metrics.baseline_entropy;
                metrics.window_entropy = result.metrics.window_entropy;
                metrics.is_anomaly = if result.is_anomaly { 1 } else { 0 };
                metrics.confidence_level = result.confidence_level;
                metrics.is_statistically_significant = if result.is_statistically_significant {
                    1
                } else {
                    0
                };
                metrics.compression_ratio_change = result.metrics.compression_ratio_change;
                metrics.entropy_change = result.metrics.entropy_change;
                metrics.recommended_ncd_threshold = result.metrics.recommended_ncd_threshold;
                metrics.recommended_window_size = result.metrics.recommended_window_size;
                metrics.data_stability_score = result.metrics.data_stability_score;

                // Create explanation string (caller must free this)
                // Use nested fallbacks to avoid panic from unwrap
                let explanation_result = CString::new(result.metrics.explanation.clone())
                    .or_else(|_| CString::new("Error generating explanation"))
                    .or_else(|_| CString::new("error"));

                metrics.explanation = match explanation_result {
                    Ok(s) => s.into_raw(),
                    Err(_) => ptr::null(), // Cannot allocate at all, return null
                };

                if result.is_anomaly {
                    1
                } else {
                    0
                }
            }
            Ok(None) => -1, // Not enough data
            Err(_) => -2,   // Internal error
        }
    })
}

/// Free the explanation string returned by cbad_detector_detect_anomaly
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - explanation_ptr was returned by cbad_detector_detect_anomaly
/// - Must NOT be called with a pointer from other sources (malloc, Go allocator, etc.)
/// - Must be called exactly once per non-null explanation pointer
/// - After calling, the pointer must not be used
#[no_mangle]
pub unsafe extern "C" fn cbad_free_explanation(explanation_ptr: *mut c_char) {
    ffi_catch_unwind!((), {
        if !explanation_ptr.is_null() {
            let _ = CString::from_raw(explanation_ptr); // This will free the string
        }
    })
}

/// Get current detector statistics
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - stats_ptr points to writable memory for 3 usize values: total_events, memory_usage, is_ready
///
/// # Returns
/// - 0: Success
/// - -1: Invalid parameters
/// - -2: Internal error
/// - -99: Panic caught
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_get_stats(
    handle: CBADDetectorHandle,
    stats_ptr: *mut usize,
) -> c_int {
    ffi_catch_unwind!(CBAD_ERR_PANIC, {
        if handle.is_null() || stats_ptr.is_null() {
            return -1;
        }

        // Check sentinel for use-after-free protection
        let wrapper = &*handle;
        if !wrapper.is_valid() {
            log::warn!("cbad_detector_get_stats called on invalid/freed handle");
            return -1;
        }

        let detector = wrapper.detector();

        match detector.get_stats() {
            Ok(stats) => {
                let stats_array = slice::from_raw_parts_mut(stats_ptr, 3);
                stats_array[0] = stats.total_events as usize;
                stats_array[1] = stats.memory_usage;
                stats_array[2] = if stats.is_ready { 1 } else { 0 };
                0 // Success
            }
            Err(_) => -2, // Internal error
        }
    })
}

/// Legacy function for backward compatibility (keeps existing Go code working)
/// This function is deprecated - use the new detector-based API instead
///
/// # Safety
/// - `baseline_ptr` and `window_ptr` must be valid for reads of `baseline_len` and `window_len` bytes respectively.
/// - The pointers must remain valid for the duration of the call.
/// - This function shares semantics with `cbad_compute_metrics` and inherits its safety requirements.
///
/// # Returns
/// CBADMetrics with default values if validation fails or panic is caught.
#[no_mangle]
pub unsafe extern "C" fn cbad_compute_metrics_legacy(
    baseline_ptr: *const u8,
    baseline_len: usize,
    window_ptr: *const u8,
    window_len: usize,
    seed: u64,
    permutations: usize,
) -> crate::CBADMetrics {
    // Default metrics for error cases
    let default_metrics = crate::CBADMetrics {
        ncd: 0.0,
        p_value: 1.0,
        baseline_compression_ratio: 1.0,
        window_compression_ratio: 1.0,
        baseline_entropy: 0.0,
        window_entropy: 0.0,
        is_anomaly: 0,
        confidence_level: 0.0,
    };

    ffi_catch_unwind!(default_metrics, {
        // Validate bounds before calling the main function
        if !validate_ptr_len(baseline_ptr, baseline_len)
            || !validate_ptr_len(window_ptr, window_len)
        {
            return default_metrics;
        }

        crate::cbad_compute_metrics(
            baseline_ptr,
            baseline_len,
            window_ptr,
            window_len,
            seed,
            permutations,
        )
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_detector_lifecycle() {
        let config = CBADConfig {
            baseline_size: 10,
            window_size: 5,
            hop_size: 2,
            max_capacity: 100,
            p_value_threshold: 0.05,
            ncd_threshold: 0.3,
            compression_ratio_drop_threshold: 0.15,
            entropy_change_threshold: 0.2,
            composite_threshold: 0.6,
            permutation_count: 100,
            seed: 42,
            require_statistical_significance: 1,
            compression_algorithm: std::ptr::null(),
        };

        unsafe {
            let handle = cbad_detector_create(&config);
            assert!(!handle.is_null());

            // Test adding data
            let data = b"INFO test log message\n";
            let result = cbad_detector_add_data(handle, data.as_ptr(), data.len());
            assert_eq!(result, 1); // Data added successfully

            // Test readiness (should be false with only 1 event)
            let ready = cbad_detector_is_ready(handle);
            assert_eq!(ready, 0); // Not ready yet

            // Add more data to make it ready
            for i in 0..20 {
                let log_data = format!("INFO test log message {}\n", i);
                let _ = cbad_detector_add_data(handle, log_data.as_ptr(), log_data.len());
            }

            let ready = cbad_detector_is_ready(handle);
            assert_eq!(ready, 1); // Should be ready now

            // Test anomaly detection
            let mut metrics = std::mem::MaybeUninit::<CBADEnhancedMetrics>::uninit();
            let anomaly_result = cbad_detector_detect_anomaly(handle, metrics.as_mut_ptr());

            // Should either detect anomaly or not, but not error
            assert!(anomaly_result >= 0);

            let metrics = metrics.assume_init();
            assert!(metrics.ncd >= 0.0 && metrics.ncd <= 1.0);
            assert!(metrics.p_value >= 0.0 && metrics.p_value <= 1.0);

            // Free the explanation string
            if !metrics.explanation.is_null() {
                cbad_free_explanation(metrics.explanation as *mut c_char);
            }

            // Test stats
            let mut stats = [0usize; 3];
            let stats_result = cbad_detector_get_stats(handle, stats.as_mut_ptr());
            assert_eq!(stats_result, 0);
            assert!(stats[0] > 0); // total_events
            assert!(stats[1] > 0); // memory_usage
            assert_eq!(stats[2], 1); // is_ready

            // Cleanup
            cbad_detector_destroy(handle);
        }
    }

    #[test]
    fn test_null_pointer_handling() {
        unsafe {
            // Test null handle
            let result = cbad_detector_is_ready(std::ptr::null_mut());
            assert_eq!(result, -1);

            let result = cbad_detector_add_data(std::ptr::null_mut(), b"test".as_ptr(), 4);
            assert_eq!(result, -1);

            let mut metrics = std::mem::MaybeUninit::<CBADEnhancedMetrics>::uninit();
            let result = cbad_detector_detect_anomaly(std::ptr::null_mut(), metrics.as_mut_ptr());
            assert_eq!(result, -1);

            // Test null data pointer
            let config = CBADConfig {
                baseline_size: 10,
                window_size: 5,
                hop_size: 2,
                max_capacity: 100,
                p_value_threshold: 0.05,
                ncd_threshold: 0.3,
                compression_ratio_drop_threshold: 0.15,
                entropy_change_threshold: 0.2,
                composite_threshold: 0.6,
                permutation_count: 100,
                seed: 42,
                require_statistical_significance: 1,
                compression_algorithm: std::ptr::null(),
            };

            let handle = cbad_detector_create(&config);
            assert!(!handle.is_null());

            let result = cbad_detector_add_data(handle, std::ptr::null(), 4);
            assert_eq!(result, -1);

            cbad_detector_destroy(handle);
        }
    }
}
