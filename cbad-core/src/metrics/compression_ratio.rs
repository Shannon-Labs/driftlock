//! Compression ratio calculator for CBAD
//!
//! This module provides compression ratio metrics that demonstrate
//! how well data compresses, which is fundamental to anomaly detection.
//! Anomalous data typically compresses poorly compared to regular patterns.

use crate::compression::CompressionAdapter;
use crate::metrics::{MetricsError, Result};

/// Compression ratio metrics for baseline vs window comparison
#[derive(Debug, Clone)]
pub struct CompressionRatioMetrics {
    /// Compression ratio for baseline data (original_size / compressed_size)
    pub baseline_ratio: f64,
    /// Compression ratio for window data (original_size / compressed_size)  
    pub window_ratio: f64,
    /// Absolute change in compression ratio (window - baseline)
    pub absolute_change: f64,
    /// Percentage change in compression ratio
    pub percentage_change: f64,
    /// Size of baseline data before compression
    pub baseline_original_size: usize,
    /// Size of baseline data after compression
    pub baseline_compressed_size: usize,
    /// Size of window data before compression
    pub window_original_size: usize,
    /// Size of window data after compression
    pub window_compressed_size: usize,
}

impl CompressionRatioMetrics {
    /// Create new compression ratio metrics
    pub fn new(
        baseline_original: usize,
        baseline_compressed: usize,
        window_original: usize,
        window_compressed: usize,
    ) -> Self {
        let baseline_ratio = if baseline_compressed > 0 {
            baseline_original as f64 / baseline_compressed as f64
        } else {
            1.0
        };

        let window_ratio = if window_compressed > 0 {
            window_original as f64 / window_compressed as f64
        } else {
            1.0
        };

        let absolute_change = window_ratio - baseline_ratio;
        let percentage_change = if baseline_ratio > 0.0 {
            (absolute_change / baseline_ratio) * 100.0
        } else {
            0.0
        };

        Self {
            baseline_ratio,
            window_ratio,
            absolute_change,
            percentage_change,
            baseline_original_size: baseline_original,
            baseline_compressed_size: baseline_compressed,
            window_original_size: window_original,
            window_compressed_size: window_compressed,
        }
    }

    /// Check if compression efficiency degraded significantly
    /// Returns true if window compresses at least threshold% worse than baseline
    pub fn is_significant_degradation(&self, threshold_percent: f64) -> bool {
        self.percentage_change < -threshold_percent
    }

    /// Get human-readable interpretation of compression change
    pub fn interpretation(&self) -> String {
        if self.percentage_change > 10.0 {
            format!("Compression improved by {:.0}%", self.percentage_change)
        } else if self.percentage_change < -10.0 {
            format!("Compression degraded by {:.0}%", -self.percentage_change)
        } else {
            "Compression efficiency similar to baseline".to_string()
        }
    }
}

/// Calculate compression ratios for baseline and window data
///
/// This function demonstrates the core principle of CBAD: anomalous data
/// typically doesn't compress as well as regular, structured data.
pub fn calculate_compression_ratios(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<CompressionRatioMetrics> {
    if baseline.is_empty() || window.is_empty() {
        return Err(MetricsError::InvalidInput(
            "Baseline and window data must not be empty".to_string(),
        ));
    }

    // Compress baseline data
    let baseline_compressed = adapter.compress(baseline).map_err(|e| {
        MetricsError::CompressionFailed(format!("Baseline compression failed: {}", e))
    })?;

    // Compress window data
    let window_compressed = adapter.compress(window).map_err(|e| {
        MetricsError::CompressionFailed(format!("Window compression failed: {}", e))
    })?;

    Ok(CompressionRatioMetrics::new(
        baseline.len(),
        baseline_compressed.len(),
        window.len(),
        window_compressed.len(),
    ))
}

/// Calculate compression ratio for a single data buffer
///
/// Useful for baseline establishment or single-sample analysis
pub fn calculate_single_compression_ratio(
    data: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<f64> {
    if data.is_empty() {
        return Ok(1.0); // No compression for empty data
    }

    let compressed = adapter
        .compress(data)
        .map_err(|e| MetricsError::CompressionFailed(format!("Compression failed: {}", e)))?;

    Ok(data.len() as f64 / compressed.len() as f64)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::compression::{create_adapter, CompressionAlgorithm};

    fn create_test_adapter() -> Box<dyn CompressionAdapter> {
        // Temporarily use zstd while OpenZL is being fixed
        create_adapter(CompressionAlgorithm::Zstd).expect("Zstd adapter")
    }

    #[test]
    fn test_compression_ratio_calculation() {
        let adapter = create_test_adapter();

        // Regular structured data (should compress well)
        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42 request_id=req-123\n".repeat(50);

        // Similar structured data (should compress similarly) - use same size for fair comparison
        let window = b"INFO 2025-10-24T00:00:01Z service=api-gateway msg=request_completed duration_ms=43 request_id=req-124\n".repeat(50);

        let metrics =
            calculate_compression_ratios(baseline.as_slice(), window.as_slice(), adapter.as_ref())
                .expect("calculate compression ratios");

        println!("Baseline compression ratio: {:.2}x", metrics.baseline_ratio);
        println!("Window compression ratio: {:.2}x", metrics.window_ratio);
        println!("Change: {:.1}%", metrics.percentage_change);

        // Both should have reasonable compression ratios
        assert!(
            metrics.baseline_ratio > 1.5,
            "Baseline should compress well"
        );
        assert!(metrics.window_ratio > 1.5, "Window should compress well");

        // Similar data should have similar compression ratios (within 30% to account for minor differences)
        assert!(
            metrics.percentage_change.abs() < 30.0,
            "Similar data should compress similarly"
        );
    }

    #[test]
    fn test_anomalous_data_compression() {
        let adapter = create_test_adapter();

        // Regular structured data
        let baseline =
            b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n"
                .repeat(100);

        // Anomalous unstructured data (stack trace - should compress poorly)
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic stack_trace=\"thread 'main' panicked at 'index out of bounds: the len is 10 but the index is 42', /rustc/.../src/libcore/slice/mod.rs:1234:5\"\n".repeat(10);

        let metrics =
            calculate_compression_ratios(baseline.as_slice(), window.as_slice(), adapter.as_ref())
                .expect("calculate compression ratios");

        println!(
            "Baseline: {:.2}x, Window: {:.2}x, Change: {:.1}%",
            metrics.baseline_ratio, metrics.window_ratio, metrics.percentage_change
        );

        // Baseline should compress much better than anomalous window
        assert!(
            metrics.baseline_ratio > metrics.window_ratio,
            "Structured baseline should compress better than anomalous window"
        );

        // Should detect significant compression degradation
        assert!(
            metrics.is_significant_degradation(20.0),
            "Should detect significant compression degradation"
        );
    }

    #[test]
    fn test_single_compression_ratio() {
        let adapter = create_test_adapter();

        // Use larger, more repetitive data that should compress well
        let data =
            b"Hello, Driftlock! This is test data for compression ratio calculation. ".repeat(10);
        println!("Original data length: {} bytes", data.len());

        let compressed = adapter.compress(&data).expect("compress");
        println!("Compressed data length: {} bytes", compressed.len());

        let ratio = calculate_single_compression_ratio(&data, adapter.as_ref())
            .expect("calculate single ratio");

        println!("Single compression ratio: {:.2}x", ratio);
        assert!(ratio > 1.0, "Data should be compressible");
    }

    #[test]
    fn test_empty_data_handling() {
        let adapter = create_test_adapter();

        let empty_data = b"";
        let ratio = calculate_single_compression_ratio(empty_data, adapter.as_ref())
            .expect("calculate empty ratio");

        assert_eq!(ratio, 1.0, "Empty data should have ratio of 1.0");
    }

    #[test]
    fn test_otlp_log_compression() {
        let adapter = create_test_adapter();

        // Realistic OTLP log entry
        let otlp_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed successfully","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42,"request_id":"req-abc123","user_id":"user-456"}}"#;

        let ratio = calculate_single_compression_ratio(otlp_log.as_bytes(), adapter.as_ref())
            .expect("calculate OTLP ratio");

        println!("OTLP log compression ratio: {:.2}x", ratio);

        // OTLP JSON should compress reasonably well due to repeated structure
        assert!(ratio > 1.2, "OTLP JSON should be compressible");

        // Test with multiple similar logs
        let mut multiple_logs = String::new();
        for i in 0..10 {
            multiple_logs.push_str(&format!(
                r#"{{"timestamp":"2025-10-24T00:00:{:02}Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{{"method":"GET","path":"/api/users","status":200,"duration_ms":{},"request_id":"req-{:03}"}}}}"#,
                i, 40 + i, i
            ));
            multiple_logs.push('\n');
        }

        let multi_ratio =
            calculate_single_compression_ratio(multiple_logs.as_bytes(), adapter.as_ref())
                .expect("calculate multiple OTLP ratio");

        println!("Multiple OTLP logs compression ratio: {:.2}x", multi_ratio);

        // Multiple similar OTLP logs should compress even better
        assert!(
            multi_ratio > ratio,
            "Multiple similar logs should compress better than single log"
        );
    }
}
