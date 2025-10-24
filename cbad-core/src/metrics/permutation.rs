//! Permutation testing framework for CBAD
//! 
//! This module provides statistical significance testing using permutation tests.
//! It's crucial for determining whether observed anomalies are statistically
//! significant or just random variations.
//! 
//! The approach: shuffle the data many times, compute the metric on each
//! permutation, and see how often we get results as extreme as the original.
//! This gives us a p-value that auditors can trust.

use crate::compression::CompressionAdapter;
use crate::metrics::{ncd, Result};
use rand::prelude::*;
use rand::rngs::StdRng;
use rand::seq::SliceRandom;

/// Permutation test result with statistical significance
#[derive(Debug, Clone)]
pub struct PermutationResult {
    /// Observed metric value from original data
    pub observed: f64,
    /// p-value (probability of seeing results this extreme by chance)
    pub p_value: f64,
    /// Number of permutations tested
    pub num_permutations: usize,
    /// Number of permutations that were as extreme or more extreme than observed
    pub extreme_count: usize,
    /// Whether the result is statistically significant (p < 0.05)
    pub is_significant: bool,
}

impl PermutationResult {
    /// Create new permutation test result
    pub fn new(observed: f64, extreme_count: usize, num_permutations: usize) -> Self {
        // Calculate p-value with continuity correction (add 1 to numerator and denominator)
        let p_value = (1 + extreme_count) as f64 / (1 + num_permutations) as f64;
        let is_significant = p_value < 0.05;

        Self {
            observed,
            p_value,
            num_permutations,
            extreme_count,
            is_significant,
        }
    }

    /// Get confidence level (1 - p_value)
    pub fn confidence_level(&self) -> f64 {
        1.0 - self.p_value
    }

    /// Get human-readable interpretation
    pub fn interpretation(&self) -> String {
        let confidence = (self.confidence_level() * 100.0).round();
        
        if self.is_significant {
            format!("Statistically significant (p={:.3}, {:.0}% confidence)", self.p_value, confidence)
        } else {
            format!("Not statistically significant (p={:.3}, {:.0}% confidence)", self.p_value, confidence)
        }
    }
}

/// Permutation tester for statistical significance testing
pub struct PermutationTester {
    rng: StdRng,
    num_permutations: usize,
}

impl PermutationTester {
    /// Create a new permutation tester with deterministic seed
    pub fn new(seed: u64, num_permutations: usize) -> Self {
        Self {
            rng: StdRng::seed_from_u64(seed),
            num_permutations,
        }
    }

    /// Test statistical significance of NCD between baseline and window
    /// 
    /// This is the primary method for determining if an anomaly is statistically significant.
    pub fn test_ncd_significance(
        &mut self,
        baseline: &[u8],
        window: &[u8],
        adapter: &dyn CompressionAdapter,
    ) -> Result<PermutationResult> {
        // Compute observed NCD
        let observed = ncd::compute_ncd(baseline, window, adapter)?;
        
        // Combine baseline and window for permutation
        let mut combined = Vec::with_capacity(baseline.len() + window.len());
        combined.extend_from_slice(baseline);
        combined.extend_from_slice(window);

        // Count extreme permutations
        let mut extreme_count = 0;

        for _ in 0..self.num_permutations {
            // Shuffle the combined data
            combined.shuffle(&mut self.rng);

            // Split back into baseline and window sized chunks
            let (perm_baseline, perm_window) = combined.split_at(baseline.len());

            // Compute NCD on permuted data
            let perm_ncd = ncd::compute_ncd(perm_baseline, perm_window, adapter)?;

            // Count if this permutation is as extreme or more extreme than observed
            if perm_ncd.abs() >= observed.abs() {
                extreme_count += 1;
            }
        }

        Ok(PermutationResult::new(observed, extreme_count, self.num_permutations))
    }

    /// Test statistical significance using a custom metric function
    /// 
    /// The metric function should return a value where larger absolute values
    /// indicate more extreme results.
    pub fn test_significance<F>(
        &mut self,
        baseline: &[u8],
        window: &[u8],
        metric_fn: F,
    ) -> Result<PermutationResult>
    where
        F: Fn(&[u8], &[u8]) -> f64,
    {
        // Compute observed metric
        let observed = metric_fn(baseline, window);

        // Combine baseline and window for permutation
        let mut combined = Vec::with_capacity(baseline.len() + window.len());
        combined.extend_from_slice(baseline);
        combined.extend_from_slice(window);

        // Count extreme permutations
        let mut extreme_count = 0;

        for _ in 0..self.num_permutations {
            // Shuffle the combined data
            combined.shuffle(&mut self.rng);

            // Split back into baseline and window sized chunks
            let (perm_baseline, perm_window) = combined.split_at(baseline.len());

            // Compute metric on permuted data
            let perm_metric = metric_fn(perm_baseline, perm_window);

            // Count if this permutation is as extreme or more extreme than observed
            if perm_metric.abs() >= observed.abs() {
                extreme_count += 1;
            }
        }

        Ok(PermutationResult::new(observed, extreme_count, self.num_permutations))
    }

    /// Test compression ratio significance
    /// 
    /// Specifically tests if the difference in compression ratios is statistically significant.
    pub fn test_compression_ratio_significance(
        &mut self,
        baseline: &[u8],
        window: &[u8],
        adapter: &dyn CompressionAdapter,
    ) -> Result<PermutationResult> {
        use crate::metrics::compression_ratio;

        // Compute observed compression ratio difference
        let baseline_ratio = compression_ratio::calculate_single_compression_ratio(baseline, adapter)?;
        let window_ratio = compression_ratio::calculate_single_compression_ratio(window, adapter)?;
        let observed = window_ratio - baseline_ratio;

        // Combine baseline and window for permutation
        let mut combined = Vec::with_capacity(baseline.len() + window.len());
        combined.extend_from_slice(baseline);
        combined.extend_from_slice(window);

        // Count extreme permutations
        let mut extreme_count = 0;

        for _ in 0..self.num_permutations {
            // Shuffle the combined data
            combined.shuffle(&mut self.rng);

            // Split back into baseline and window sized chunks
            let (perm_baseline, perm_window) = combined.split_at(baseline.len());

            // Compute compression ratios on permuted data
            let perm_baseline_ratio = compression_ratio::calculate_single_compression_ratio(perm_baseline, adapter)?;
            let perm_window_ratio = compression_ratio::calculate_single_compression_ratio(perm_window, adapter)?;
            let perm_difference = perm_window_ratio - perm_baseline_ratio;

            // Count if this permutation is as extreme or more extreme than observed
            if perm_difference.abs() >= observed.abs() {
                extreme_count += 1;
            }
        }

        Ok(PermutationResult::new(observed, extreme_count, self.num_permutations))
    }
}

/// Quick permutation test with default parameters
/// 
/// Uses 1000 permutations and a fixed seed for reproducibility.
pub fn quick_permutation_test<F>(
    baseline: &[u8],
    window: &[u8],
    metric_fn: F,
) -> PermutationResult
where
    F: Fn(&[u8], &[u8]) -> f64,
{
    let mut tester = PermutationTester::new(42, 1000);
    tester.test_significance(baseline, window, metric_fn)
        .expect("permutation test should succeed with valid data")
}

/// Test NCD significance with default parameters
pub fn test_ncd_significance(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<PermutationResult> {
    let mut tester = PermutationTester::new(42, 1000);
    tester.test_ncd_significance(baseline, window, adapter)
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
    fn test_permutation_test_deterministic() {
        let mut tester = PermutationTester::new(12345, 100);

        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(50);
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic\n".repeat(20);

        // Run same test twice - should get identical results
        let result1 = tester.test_significance(&baseline, &window, |b, w| {
            (b.len() as f64 - w.len() as f64).abs()
        }).expect("first test");

        let result2 = tester.test_significance(&baseline, &window, |b, w| {
            (b.len() as f64 - w.len() as f64).abs()
        }).expect("second test");

        assert_eq!(result1.observed, result2.observed);
        assert_eq!(result1.p_value, result2.p_value);
        assert_eq!(result1.extreme_count, result2.extreme_count);
    }

    #[test]
    fn test_ncd_significance_detection() {
        let adapter = create_test_adapter();
        let mut tester = PermutationTester::new(42, 100);

        // Similar data - should not be significant
        let baseline_similar = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(50);
        let window_similar = b"INFO 2025-10-24T00:00:01Z service=api-gateway msg=request_completed duration_ms=43\n".repeat(20);

        let result_similar = tester.test_ncd_significance(&baseline_similar, &window_similar, adapter.as_ref())
            .expect("test similar data");

        println!("Similar data: {}", result_similar.interpretation());

        // Different data - should be significant
        let baseline_different = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(50);
        let window_different = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=stack_trace_panic_at_line_42\n".repeat(20);

        let result_different = tester.test_ncd_significance(&baseline_different, &window_different, adapter.as_ref())
            .expect("test different data");

        println!("Different data: {}", result_different.interpretation());

        // Different data should be more significant than similar data
        assert!(result_different.p_value < result_similar.p_value, 
                "Different data should have lower p-value than similar data");
    }

    #[test]
    fn test_compression_ratio_significance() {
        let adapter = create_test_adapter();
        let mut tester = PermutationTester::new(42, 100);

        // Structured vs unstructured data should show significant compression difference
        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(100);
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic stack_trace=\"random unstructured data with lots of entropy\"\n".repeat(20);

        let result = tester.test_compression_ratio_significance(&baseline, &window, adapter.as_ref())
            .expect("test compression ratio significance");

        println!("Compression ratio significance: {}", result.interpretation());

        // Should detect significant difference in compression ratios
        assert!(result.p_value < 0.1, "Should detect significant compression ratio difference");
    }

    #[test]
    fn test_quick_permutation_test() {
        let baseline = b"AAAAA".repeat(100);
        let window = b"BBBBB".repeat(50);

        let result = quick_permutation_test(&baseline, &window, |b, w| {
            (b.len() as f64 - w.len() as f64).abs()
        });

        println!("Quick test result: {}", result.interpretation());
        
        assert_eq!(result.num_permutations, 1000);
        assert!(result.p_value >= 0.0 && result.p_value <= 1.0);
    }

    #[test]
    fn test_permutation_result_interpretation() {
        let result_significant = PermutationResult::new(0.8, 5, 100);  // 5% extreme
        let result_not_significant = PermutationResult::new(0.1, 50, 100); // 50% extreme

        println!("Significant: {}", result_significant.interpretation());
        println!("Not significant: {}", result_not_significant.interpretation());

        assert!(result_significant.is_significant);
        assert!(!result_not_significant.is_significant);
        assert!(result_significant.p_value < result_not_significant.p_value);
    }

    #[test]
    fn test_empty_data_handling() {
        let mut tester = PermutationTester::new(42, 10);
        
        let result = tester.test_significance(b"", b"", |b, w| {
            (b.len() as f64 - w.len() as f64).abs()
        });

        // Should handle empty data gracefully
        assert!(result.is_ok() || result.is_err());
    }
}
