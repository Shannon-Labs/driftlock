// Driftlock - Compression-Based Anomaly Detection
// Patent-Pending Technology - Shannon Labs, LLC
// Licensed under Apache 2.0. Commercial licenses available.
// See PATENTS.md and LICENSE for details.

//! # Driftlock CBAD Core
//!
//! Patent-pending compression-based anomaly detection technology.
//! Licensed under Apache 2.0. Commercial licenses available from Shannon Labs.

pub mod anomaly;
pub mod compression;
pub mod ffi;
pub mod ffi_simplified;
pub mod metrics;
pub mod performance;
pub mod window;

use serde::{Deserialize, Serialize};

/// Legacy metrics struct (deprecated - use metrics::AnomalyMetrics instead)
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Metrics {
    pub entropy: f64,
    pub compression_ratio: f64,
    pub ncd: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComputeConfig {
    pub window_size: usize,
    pub hop_size: usize,
    pub threshold: f64,
    pub deterministic_seed: u64,
    pub permutation_count: usize,
}

impl Default for ComputeConfig {
    fn default() -> Self {
        Self {
            window_size: 1000,
            hop_size: 100,
            threshold: 0.05,
            deterministic_seed: 42,
            permutation_count: 1000,
        }
    }
}

/// Compute anomaly detection metrics using compression-based analysis
///
/// This is the main entry point for the CBAD algorithm. It analyzes the
/// provided data and returns comprehensive metrics including compression
/// ratios, entropy, NCD scores, and statistical significance testing.
///
/// # Arguments
/// * `baseline` - Historical data representing normal patterns
/// * `window` - Current data to analyze for anomalies  
/// * `adapter` - Compression adapter to use (OpenZL recommended)
/// * `config` - Configuration for the analysis
///
/// # Returns
/// * `AnomalyMetrics` - Complete metrics with glass-box explanation
pub fn compute_metrics(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn compression::CompressionAdapter,
    config: &ComputeConfig,
) -> metrics::Result<metrics::AnomalyMetrics> {
    metrics::compute_metrics(
        baseline,
        window,
        adapter,
        config.permutation_count,
        config.deterministic_seed,
    )
}

/// Quick metrics computation with default configuration
pub fn compute_metrics_quick(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn compression::CompressionAdapter,
) -> metrics::Result<metrics::AnomalyMetrics> {
    let config = ComputeConfig::default();
    compute_metrics(baseline, window, adapter, &config)
}

/// Legacy function for backward compatibility (deprecated)
#[deprecated(
    since = "0.1.0",
    note = "Use compute_metrics with proper parameters instead"
)]
pub fn compute_metrics_legacy(_data: &[u8], _cfg: &ComputeConfig) -> Metrics {
    Metrics {
        entropy: 0.0,
        compression_ratio: 1.0,
        ncd: 0.0,
    }
}

// C FFI exports for Go integration
use std::os::raw::c_int;
use std::panic::{catch_unwind, AssertUnwindSafe};
use std::slice;

/// Maximum allowed data size for FFI calls (16MB)
const MAX_FFI_DATA_SIZE: usize = 16 * 1024 * 1024;

/// C-compatible metrics structure for FFI
#[repr(C)]
#[derive(Clone, Copy)]
pub struct CBADMetrics {
    pub ncd: f64,
    pub p_value: f64,
    pub baseline_compression_ratio: f64,
    pub window_compression_ratio: f64,
    pub baseline_entropy: f64,
    pub window_entropy: f64,
    pub is_anomaly: c_int, // 0 = false, 1 = true
    pub confidence_level: f64,
}

/// Initialize logging (called from Go)
#[no_mangle]
pub extern "C" fn cbad_init_logging() {
    let _ = env_logger::try_init();
}

/// Compute CBAD metrics via C FFI
///
/// # Safety
/// This function is unsafe because it deals with raw pointers from C.
/// Callers must ensure:
/// - baseline_ptr and window_ptr are valid pointers to baseline_len and window_len bytes respectively
/// - The memory is properly allocated and accessible
/// - baseline_len and window_len are accurate and <= 16MB each
///
/// # Returns
/// CBADMetrics with default values if:
/// - Pointers are null
/// - Lengths are 0 or exceed 16MB
/// - Compression adapter creation fails
/// - Computation fails
/// - A panic is caught
#[no_mangle]
pub unsafe extern "C" fn cbad_compute_metrics(
    baseline_ptr: *const u8,
    baseline_len: usize,
    window_ptr: *const u8,
    window_len: usize,
    seed: u64,
    permutations: usize,
) -> CBADMetrics {
    // Default metrics for error/failure cases
    let default_metrics = CBADMetrics {
        ncd: 0.0,
        p_value: 1.0,
        baseline_compression_ratio: 1.0,
        window_compression_ratio: 1.0,
        baseline_entropy: 0.0,
        window_entropy: 0.0,
        is_anomaly: 0,
        confidence_level: 0.0,
    };

    // Wrap entire function body in catch_unwind to prevent panics from unwinding into C
    let result = catch_unwind(AssertUnwindSafe(|| {
        // Initialize logger if not already done
        cbad_init_logging();

        // Validate pointers
        if baseline_ptr.is_null() || window_ptr.is_null() {
            return default_metrics;
        }

        // Validate lengths (bounds checking)
        if baseline_len == 0
            || window_len == 0
            || baseline_len > MAX_FFI_DATA_SIZE
            || window_len > MAX_FFI_DATA_SIZE
        {
            log::warn!(
                "cbad_compute_metrics: invalid lengths baseline={} window={}",
                baseline_len,
                window_len
            );
            return default_metrics;
        }

        // Create slices from raw pointers
        let baseline = slice::from_raw_parts(baseline_ptr, baseline_len);
        let window = slice::from_raw_parts(window_ptr, window_len);

        // Use OpenZL adapter if available, otherwise fall back to zstd
        #[cfg(feature = "openzl")]
        let adapter = match compression::create_adapter(compression::CompressionAlgorithm::OpenZL) {
            Ok(adapter) => adapter,
            Err(_) => match compression::create_adapter(compression::CompressionAlgorithm::Zstd) {
                Ok(adapter) => adapter,
                Err(_) => {
                    return default_metrics;
                }
            },
        };

        #[cfg(not(feature = "openzl"))]
        let adapter = match compression::create_adapter(compression::CompressionAlgorithm::Zstd) {
            Ok(adapter) => adapter,
            Err(_) => {
                return default_metrics;
            }
        };

        // Compute metrics
        let config = ComputeConfig {
            deterministic_seed: seed,
            permutation_count: permutations,
            ..Default::default()
        };

        match compute_metrics(baseline, window, adapter.as_ref(), &config) {
            Ok(metrics) => CBADMetrics {
                ncd: metrics.ncd,
                p_value: metrics.p_value,
                baseline_compression_ratio: metrics.baseline_compression_ratio,
                window_compression_ratio: metrics.window_compression_ratio,
                baseline_entropy: metrics.baseline_entropy,
                window_entropy: metrics.window_entropy,
                is_anomaly: if metrics.is_anomaly { 1 } else { 0 },
                confidence_level: metrics.confidence_level,
            },
            Err(_) => default_metrics,
        }
    }));

    // Return default metrics if panic was caught
    match result {
        Ok(metrics) => metrics,
        Err(e) => {
            // Log the panic if possible
            if let Some(s) = e.downcast_ref::<&str>() {
                log::error!("cbad_compute_metrics panic caught: {}", s);
            } else if let Some(s) = e.downcast_ref::<String>() {
                log::error!("cbad_compute_metrics panic caught: {}", s);
            } else {
                log::error!("cbad_compute_metrics panic caught: unknown");
            }
            default_metrics
        }
    }
}

/// Legacy C FFI function (kept for backward compatibility)
#[no_mangle]
pub extern "C" fn cbad_compute_metrics_len(_: *const u8, len: usize) -> f64 {
    len as f64
}
