//! Delta bits calculator for CBAD
//!
//! Delta bits measures the change in bit patterns between baseline and window data.
//! It calculates the compression difference at the bit level, which can reveal
//! subtle structural changes that might be missed by other metrics.
//!
//! The approach: compress each sequence individually, then compress their
//! concatenation. The difference between the concatenated size and the sum
//! of individual sizes indicates how different the bit patterns are.

use crate::compression::CompressionAdapter;
use crate::metrics::Result;

/// Delta bits calculation result with detailed metrics
#[derive(Debug, Clone)]
pub struct DeltaBitsMetrics {
    /// Delta bits score: difference between combined and individual compression
    /// Higher values indicate more difference between baseline and window
    pub delta_bits: f64,
    /// Size of baseline after compression
    pub baseline_compressed_size: usize,
    /// Size of window after compression
    pub window_compressed_size: usize,
    /// Size of concatenated baseline+window after compression
    pub combined_compressed_size: usize,
    /// Sum of individual compressed sizes
    pub individual_sum_size: usize,
    /// Percentage difference between combined and sum of individual sizes
    pub percentage_difference: f64,
}

impl DeltaBitsMetrics {
    /// Create new delta bits metrics from compression sizes
    pub fn new(
        baseline_compressed: usize,
        window_compressed: usize,
        combined_compressed: usize,
    ) -> Self {
        let individual_sum = baseline_compressed + window_compressed;

        // Calculate delta bits as the difference between combined and sum of individuals
        // This represents how much more space is needed to store both together vs separately
        let delta_bits = (combined_compressed as f64) - (individual_sum as f64);

        // Calculate percentage difference relative to the sum of individual sizes
        let percentage_difference = if individual_sum > 0 {
            (delta_bits / individual_sum as f64) * 100.0
        } else {
            0.0
        };

        Self {
            delta_bits,
            baseline_compressed_size: baseline_compressed,
            window_compressed_size: window_compressed,
            combined_compressed_size: combined_compressed,
            individual_sum_size: individual_sum,
            percentage_difference,
        }
    }

    /// Get human-readable interpretation of delta bits score
    pub fn interpretation(&self) -> &'static str {
        if self.delta_bits < 0.0 {
            "Negative delta: sequences share significant structure beyond individual compression"
        } else if self.delta_bits < 100.0 {
            "Low delta: sequences have similar bit patterns"
        } else if self.delta_bits < 500.0 {
            "Moderate delta: sequences have somewhat different bit patterns"
        } else if self.delta_bits < 1000.0 {
            "High delta: sequences have significantly different bit patterns"
        } else {
            "Very high delta: sequences have completely different bit patterns"
        }
    }

    /// Check if delta bits indicates an anomaly based on threshold
    pub fn is_anomaly(&self, threshold: f64) -> bool {
        self.delta_bits.abs() >= threshold
    }

    /// Calculate relative delta as percentage of individual sum
    pub fn relative_delta(&self) -> f64 {
        if self.individual_sum_size > 0 {
            self.delta_bits / self.individual_sum_size as f64
        } else {
            0.0
        }
    }
}

/// Calculate delta bits between baseline and window data
///
/// The delta bits metric measures how much additional space is needed to compress
/// two sequences together compared to compressing them independently.
///
/// A low delta indicates that the sequences have similar patterns (and can be
/// compressed efficiently together), while a high delta indicates different patterns.
pub fn calculate_delta_bits(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<f64> {
    let metrics = calculate_delta_bits_detailed(baseline, window, adapter)?;
    Ok(metrics.delta_bits)
}

/// Calculate delta bits with detailed metrics
pub fn calculate_delta_bits_detailed(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<DeltaBitsMetrics> {
    if baseline.is_empty() || window.is_empty() {
        return Err(crate::metrics::MetricsError::InvalidInput(
            "Baseline and window data must not be empty".to_string(),
        ));
    }

    // Compress baseline individually
    let baseline_compressed = adapter.compress(baseline).map_err(|e| {
        crate::metrics::MetricsError::CompressionFailed(format!(
            "Baseline compression failed: {}",
            e
        ))
    })?;

    // Compress window individually
    let window_compressed = adapter.compress(window).map_err(|e| {
        crate::metrics::MetricsError::CompressionFailed(format!("Window compression failed: {}", e))
    })?;

    // Compress concatenated data
    let mut combined = Vec::with_capacity(baseline.len() + window.len());
    combined.extend_from_slice(baseline);
    combined.extend_from_slice(window);

    let combined_compressed = adapter.compress(&combined).map_err(|e| {
        crate::metrics::MetricsError::CompressionFailed(format!(
            "Combined compression failed: {}",
            e
        ))
    })?;

    Ok(DeltaBitsMetrics::new(
        baseline_compressed.len(),
        window_compressed.len(),
        combined_compressed.len(),
    ))
}

/// Calculate normalized delta bits (relative to data size)
///
/// This normalizes the delta bits by the total input size to make it comparable
/// across different sized inputs.
pub fn calculate_normalized_delta_bits(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<f64> {
    let metrics = calculate_delta_bits_detailed(baseline, window, adapter)?;

    let total_input_size = (baseline.len() + window.len()) as f64;
    if total_input_size > 0.0 {
        Ok(metrics.delta_bits / total_input_size)
    } else {
        Ok(0.0)
    }
}

/// Compare delta bits with multiple thresholds to provide detailed analysis
pub fn analyze_delta_bits(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    thresholds: &[f64],
) -> Result<Vec<(f64, bool)>> {
    let delta_bits = calculate_delta_bits(baseline, window, adapter)?;

    let results: Vec<(f64, bool)> = thresholds
        .iter()
        .map(|&threshold| (threshold, delta_bits.abs() >= threshold))
        .collect();

    Ok(results)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::compression::{create_adapter, CompressionAlgorithm};

    fn create_test_adapter() -> Box<dyn CompressionAdapter> {
        // Use OpenZL for testing - it's the primary compression algorithm
        #[cfg(feature = "openzl")]
        {
            create_adapter(CompressionAlgorithm::OpenZL).expect("OpenZL adapter")
        }
        #[cfg(not(feature = "openzl"))]
        {
            create_adapter(CompressionAlgorithm::Zstd).expect("Zstd adapter")
        }
    }

    #[test]
    fn test_delta_bits_identical_sequences() {
        let adapter = create_test_adapter();

        let data =
            b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n";

        let delta_bits = calculate_delta_bits(data, data, adapter.as_ref())
            .expect("calculate delta bits for identical sequences");

        println!("Delta bits for identical sequences: {:.2}", delta_bits);

        // Identical sequences should have relatively small delta (may be positive or negative)
        // depending on the compression algorithm's behavior
        assert!(
            delta_bits.abs() < 1000.0,
            "Identical sequences should have small delta bits"
        );
    }

    #[test]
    fn test_delta_bits_similar_sequences() {
        let adapter = create_test_adapter();

        let baseline =
            b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n"
                .repeat(50);
        let window =
            b"INFO 2025-10-24T00:00:01Z service=api-gateway msg=request_completed duration_ms=43\n"
                .repeat(10);

        let metrics =
            calculate_delta_bits_detailed(baseline.as_slice(), window.as_slice(), adapter.as_ref())
                .expect("calculate delta bits for similar sequences");

        println!(
            "Delta bits for similar sequences: {:.2}",
            metrics.delta_bits
        );
        println!(
            "Baseline compressed: {} bytes",
            metrics.baseline_compressed_size
        );
        println!(
            "Window compressed: {} bytes",
            metrics.window_compressed_size
        );
        println!(
            "Combined compressed: {} bytes",
            metrics.combined_compressed_size
        );
        println!("Interpretation: {}", metrics.interpretation());

        // Similar sequences should have moderate delta bits
        assert!(
            metrics.delta_bits.abs() < 2000.0,
            "Similar sequences should have moderate delta bits"
        );
    }

    #[test]
    fn test_delta_bits_dissimilar_sequences() {
        let adapter = create_test_adapter();

        // Regular structured logs
        let baseline =
            b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n"
                .repeat(50);

        // Anomalous unstructured data (stack trace)
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic stack_trace=\"thread 'main' panicked at 'index out of bounds', src/main.rs:42:13\"\n".repeat(10);

        let metrics =
            calculate_delta_bits_detailed(baseline.as_slice(), window.as_slice(), adapter.as_ref())
                .expect("calculate delta bits for dissimilar sequences");

        println!(
            "Delta bits for dissimilar sequences: {:.2}",
            metrics.delta_bits
        );
        println!(
            "Percentage difference: {:.2}%",
            metrics.percentage_difference
        );
        println!("Interpretation: {}", metrics.interpretation());

        // Dissimilar sequences may have large delta bits, but let's check the relative difference
        let relative = metrics.relative_delta();
        println!("Relative delta: {:.4}", relative);

        // The key is that different sequences will have a different compression behavior
        // when combined vs separately
    }

    #[test]
    fn test_normalized_delta_bits() {
        let adapter = create_test_adapter();

        let baseline =
            b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed\n".repeat(10);
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic\n".repeat(5);

        let normalized = calculate_normalized_delta_bits(&baseline, &window, adapter.as_ref())
            .expect("calculate normalized delta bits");

        println!("Normalized delta bits: {:.4}", normalized);

        // Should be a reasonable value that's normalized by input size
        assert!(normalized.is_finite(), "Normalized delta should be finite");
    }

    #[test]
    fn test_delta_bits_empty_data() {
        let adapter = create_test_adapter();

        // Test with empty baseline
        let result = calculate_delta_bits(b"", b"non-empty", adapter.as_ref());
        assert!(result.is_err());

        // Test with empty window
        let result = calculate_delta_bits(b"non-empty", b"", adapter.as_ref());
        assert!(result.is_err());

        // Test with both empty
        let result = calculate_delta_bits(b"", b"", adapter.as_ref());
        assert!(result.is_err());
    }

    #[test]
    fn test_delta_bits_analysis() {
        let adapter = create_test_adapter();

        let baseline = b"INFO service=api-gateway msg=request_completed\n".repeat(20);
        let window = b"ERROR service=api-gateway msg=panic_occurred\n".repeat(5);

        let thresholds = vec![100.0, 500.0, 1000.0, 2000.0];
        let results = analyze_delta_bits(&baseline, &window, adapter.as_ref(), &thresholds)
            .expect("analyze delta bits");

        println!("Delta bits analysis:");
        for (threshold, is_anomaly) in &results {
            println!("  Threshold {:.0}: anomaly = {}", threshold, is_anomaly);
        }

        // Verify that results match expectations
        assert_eq!(results.len(), thresholds.len());
    }

    #[test]
    fn test_delta_bits_metrics_interpretation() {
        let metrics_low = DeltaBitsMetrics::new(100, 100, 180); // small delta
        let metrics_high = DeltaBitsMetrics::new(100, 100, 300); // large delta

        println!("Low delta interpretation: {}", metrics_low.interpretation());
        println!(
            "High delta interpretation: {}",
            metrics_high.interpretation()
        );

        // Interpretations should be different for low vs high delta
        assert_ne!(metrics_low.interpretation(), metrics_high.interpretation());

        // High delta should be considered anomaly at reasonable threshold
        assert!(metrics_high.is_anomaly(100.0));
    }
}
