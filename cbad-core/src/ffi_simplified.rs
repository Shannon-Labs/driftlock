// Simplified FFI functions for the demo CLI

use crate::anomaly::{AnomalyConfig, AnomalyDetector};
use crate::compression::CompressionAlgorithm;
use crate::window::{DataEvent, WindowConfig, PrivacyConfig};
use std::ffi::CString;
use std::os::raw::{c_char, c_int};
use std::ptr;
use std::collections::HashMap;

/// Opaque handle for AnomalyDetector instances
pub type CBADDetectorHandle = *mut AnomalyDetector;

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
    pub explanation: *const c_char,
}

/// Create a simple anomaly detector with default configuration
/// Returns a handle to the detector, or null on error
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_create_simple() -> CBADDetectorHandle {
    let window_config = WindowConfig {
        baseline_size: 50,
        window_size: 25,
        hop_size: 10,
        max_capacity: 1000,
        time_window: None,
        privacy_config: PrivacyConfig::default(),
    };

    let config = AnomalyConfig {
        window_config,
        compression_algorithm: CompressionAlgorithm::Zstd,
        p_value_threshold: 0.2,
        ncd_threshold: 0.15,
        permutation_count: 1000,
        seed: 42,
        require_statistical_significance: true,
    };

    match AnomalyDetector::new(config) {
        Ok(detector) => Box::into_raw(Box::new(detector)),
        Err(_) => ptr::null_mut(),
    }
}

/// Add a transaction to the detector (ingestion phase)
/// Returns 1 on success, 0 on failure
#[no_mangle]
pub unsafe extern "C" fn cbad_add_transaction(
    handle: CBADDetectorHandle,
    data: *const c_char,
    len: usize,
) -> c_int {
    if handle.is_null() || data.is_null() {
        return 0;
    }

    let detector = &mut *handle;
    
    // Convert C string to Rust string
    let data_slice = std::slice::from_raw_parts(data as *const u8, len);
    let data_str = match std::str::from_utf8(data_slice) {
        Ok(s) => s,
        Err(_) => return 0,
    };

    // Create metadata hashmap
    let mut metadata = HashMap::new();
    metadata.insert("source".to_string(), "demo".to_string());

    // Create data event
    let event = DataEvent {
        timestamp: std::time::SystemTime::now(),
        data: data_str.as_bytes().to_vec(),
        metadata,
        integrity_hash: Some(0),
    };

    // Add data to detector
    match detector.add_event(event) {
        Ok(_) => 1,
        Err(_) => 0,
    }
}

/// Run anomaly detection on the current state
/// Returns metrics structure with detection results
#[no_mangle]
pub unsafe extern "C" fn cbad_detect(
    handle: CBADDetectorHandle,
) -> CBADEnhancedMetrics {
    if handle.is_null() {
        return CBADEnhancedMetrics {
            ncd: 0.0,
            p_value: 0.0,
            baseline_compression_ratio: 0.0,
            window_compression_ratio: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            is_anomaly: 0,
            confidence_level: 0.0,
            is_statistically_significant: 0,
            compression_ratio_change: 0.0,
            entropy_change: 0.0,
            explanation: ptr::null(),
        };
    }

    let detector = &mut *handle;

    // Check if detector is ready
    if !detector.is_ready().unwrap_or(false) {
        return CBADEnhancedMetrics {
            ncd: 0.0,
            p_value: 0.0,
            baseline_compression_ratio: 0.0,
            window_compression_ratio: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            is_anomaly: 0,
            confidence_level: 0.0,
            is_statistically_significant: 0,
            compression_ratio_change: 0.0,
            entropy_change: 0.0,
            explanation: ptr::null(),
        };
    }

    // Perform anomaly detection
    match detector.detect_anomaly() {
        Ok(Some(result)) => {
            let explanation_str = format!(
                "NCD={} (threshold={}), P-value={}. {}. Compression ratio changed by {}x, entropy changed by {}.",
                result.metrics.ncd,
                0.3, // Using threshold from config
                result.metrics.p_value,
                if result.is_statistically_significant {
                    "Statistically significant"
                } else {
                    "Not statistically significant"
                },
                result.metrics.compression_ratio_change,
                result.metrics.entropy_change
            );

            let explanation_cstring = CString::new(explanation_str).unwrap_or_else(|_| CString::new("Error").unwrap());
            let explanation_ptr = explanation_cstring.into_raw();

            CBADEnhancedMetrics {
                ncd: result.metrics.ncd,
                p_value: result.metrics.p_value,
                baseline_compression_ratio: result.metrics.baseline_compression_ratio,
                window_compression_ratio: result.metrics.window_compression_ratio,
                baseline_entropy: result.metrics.baseline_entropy,
                window_entropy: result.metrics.window_entropy,
                is_anomaly: if result.is_anomaly { 1 } else { 0 },
                confidence_level: result.confidence_level,
                is_statistically_significant: if result.is_statistically_significant { 1 } else { 0 },
                compression_ratio_change: result.metrics.compression_ratio_change,
                entropy_change: result.metrics.entropy_change,
                explanation: explanation_ptr,
            }
        }
        Ok(None) => CBADEnhancedMetrics {
            ncd: 0.0,
            p_value: 0.0,
            baseline_compression_ratio: 0.0,
            window_compression_ratio: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            is_anomaly: 0,
            confidence_level: 0.0,
            is_statistically_significant: 0,
            compression_ratio_change: 0.0,
            entropy_change: 0.0,
            explanation: ptr::null(),
        },
        Err(_) => CBADEnhancedMetrics {
            ncd: 0.0,
            p_value: 0.0,
            baseline_compression_ratio: 0.0,
            window_compression_ratio: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            is_anomaly: 0,
            confidence_level: 0.0,
            is_statistically_significant: 0,
            compression_ratio_change: 0.0,
            entropy_change: 0.0,
            explanation: ptr::null(),
        },
    }
}

/// Check if detector is ready for anomaly detection
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_ready(handle: CBADDetectorHandle) -> c_int {
    if handle.is_null() {
        return 0;
    }

    let detector = &*handle;
    
    match detector.is_ready() {
        Ok(ready) => if ready { 1 } else { 0 },
        Err(_) => 0,
    }
}

/// Free the detector
#[no_mangle]
pub unsafe extern "C" fn cbad_detector_free(handle: CBADDetectorHandle) {
    if !handle.is_null() {
        let _ = Box::from_raw(handle);
    }
}

/// Free a string allocated by Rust
#[no_mangle]
pub unsafe extern "C" fn cbad_free_string(s: *mut c_char) {
    if !s.is_null() {
        let _ = CString::from_raw(s);
    }
}