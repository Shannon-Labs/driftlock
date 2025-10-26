//! Enhanced C FFI exports for Go integration
//! 
//! This module provides production-ready C FFI exports that enable
//! Go-based applications to use the complete CBAD anomaly detection
//! engine with streaming capabilities and proper error handling.

use crate::anomaly::{AnomalyConfig, AnomalyDetector};
use crate::compression::CompressionAlgorithm;
use crate::window::{DataEvent, WindowConfig};
use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int};
use std::ptr;
use std::slice;

/// Opaque handle for AnomalyDetector instances
pub type CBADDetectorHandle = *mut AnomalyDetector;

/// Configuration for anomaly detection (C-compatible)
#[repr(C)]
pub struct CBADConfig {
    pub baseline_size: usize,
    pub window_size: usize,
    pub hop_size: usize,
    pub max_capacity: usize,
    pub p_value_threshold: f64,
    pub ncd_threshold: f64,
    pub permutation_count: usize,
    pub seed: u64,
    pub require_statistical_significance: c_int, // 0 = false, 1 = true
    pub compression_algorithm: *const c_char, // "zstd", "lz4", "gzip", "openzl"
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
}

/// Create a new anomaly detector with configuration
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - config_ptr is a valid pointer to a CBADConfig struct
/// - compression_algorithm string is valid UTF-8 and null-terminated
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_create(
    config_ptr: *const CBADConfig,
) -> CBADDetectorHandle {
    if config_ptr.is_null() {
        return ptr::null_mut();
    }

    let config = &*config_ptr;
    
    // Parse compression algorithm
    let algo = if config.compression_algorithm.is_null() {
        CompressionAlgorithm::Zstd
    } else {
        match CStr::from_ptr(config.compression_algorithm).to_str() {
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

    // Create anomaly configuration
    let anomaly_config = AnomalyConfig {
        window_config,
        compression_algorithm: algo,
        p_value_threshold: config.p_value_threshold,
        ncd_threshold: config.ncd_threshold,
        permutation_count: config.permutation_count,
        seed: config.seed,
        require_statistical_significance: config.require_statistical_significance != 0,
    };

    // Create detector
    match AnomalyDetector::new(anomaly_config) {
        Ok(detector) => Box::into_raw(Box::new(detector)),
        Err(_) => ptr::null_mut(),
    }
}

/// Destroy an anomaly detector and free its memory
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_destroy(handle: CBADDetectorHandle) {
    if !handle.is_null() {
        let _ = Box::from_raw(handle); // This will drop the detector
    }
}

/// Add data to the anomaly detector
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - data_ptr is a valid pointer to data_len bytes
/// - data_len is accurate
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_add_data(
    handle: CBADDetectorHandle,
    data_ptr: *const u8,
    data_len: usize,
) -> c_int {
    if handle.is_null() || data_ptr.is_null() {
        return -1; // Invalid parameters
    }

    let detector = &*handle;
    let data = slice::from_raw_parts(data_ptr, data_len);
    
    let event = DataEvent::new(data.to_vec());
    
    match detector.add_event(event) {
        Ok(added) => if added { 1 } else { 0 },
        Err(_) => -2, // Internal error
    }
}

/// Check if the detector has enough data for analysis
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_is_ready(handle: CBADDetectorHandle) -> c_int {
    if handle.is_null() {
        return -1;
    }

    let detector = &*handle;
    
    match detector.is_ready() {
        Ok(ready) => if ready { 1 } else { 0 },
        Err(_) => -2, // Internal error
    }
}

/// Perform anomaly detection and get results
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - metrics_ptr points to a valid CBADEnhancedMetrics struct that can be written to
/// 
/// Returns:
/// - 1 if anomaly detected
/// - 0 if no anomaly detected
/// - -1 if not enough data or error
/// - -2 if internal error
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_detect_anomaly(
    handle: CBADDetectorHandle,
    metrics_ptr: *mut CBADEnhancedMetrics,
) -> c_int {
    if handle.is_null() || metrics_ptr.is_null() {
        return -1;
    }

    let detector = &*handle;
    
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
            metrics.is_statistically_significant = if result.is_statistically_significant { 1 } else { 0 };
            metrics.compression_ratio_change = result.metrics.compression_ratio_change;
            metrics.entropy_change = result.metrics.entropy_change;
            
            // Create explanation string (caller must free this)
            let explanation = CString::new(result.metrics.explanation.clone())
                .unwrap_or_else(|_| CString::new("Error generating explanation").unwrap());
            metrics.explanation = explanation.into_raw();
            
            if result.is_anomaly { 1 } else { 0 }
        },
        Ok(None) => -1, // Not enough data
        Err(_) => -2,   // Internal error
    }
}

/// Free the explanation string returned by cbad_detector_detect_anomaly
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - explanation_ptr was returned by cbad_detector_detect_anomaly
#[no_mangle]
pub unsafe extern "C" fn cbad_free_explanation(explanation_ptr: *mut c_char) {
    if !explanation_ptr.is_null() {
        let _ = CString::from_raw(explanation_ptr); // This will free the string
    }
}

/// Get current detector statistics
/// 
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - handle is a valid pointer returned by cbad_detector_create
/// - stats_ptr points to writable memory for 3 usize values: total_events, memory_usage, is_ready
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_get_stats(
    handle: CBADDetectorHandle,
    stats_ptr: *mut usize,
) -> c_int {
    if handle.is_null() || stats_ptr.is_null() {
        return -1;
    }

    let detector = &*handle;
    
    match detector.get_stats() {
        Ok(stats) => {
            let stats_array = slice::from_raw_parts_mut(stats_ptr, 3);
            stats_array[0] = stats.total_events as usize;
            stats_array[1] = stats.memory_usage;
            stats_array[2] = if stats.is_ready { 1 } else { 0 };
            0 // Success
        },
        Err(_) => -2, // Internal error
    }
}

/// Legacy function for backward compatibility (keeps existing Go code working)
/// This function is deprecated - use the new detector-based API instead
/// # Safety
/// - `baseline_ptr` and `window_ptr` must be valid for reads of `baseline_len` and `window_len` bytes respectively.
/// - The pointers must remain valid for the duration of the call.
/// - This function shares semantics with `cbad_compute_metrics` and inherits its safety requirements.
#[no_mangle]
pub unsafe extern "C" fn cbad_compute_metrics_legacy(
    baseline_ptr: *const u8,
    baseline_len: usize,
    window_ptr: *const u8,
    window_len: usize,
    seed: u64,
    permutations: usize,
) -> crate::CBADMetrics {
    crate::cbad_compute_metrics(baseline_ptr, baseline_len, window_ptr, window_len, seed, permutations)
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
