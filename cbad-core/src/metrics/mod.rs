//! Metrics module for CBAD (Compression-Based Anomaly Detection)
//! 
//! This module provides the core anomaly detection metrics that give users
//! confidence in anomaly detection through glass-box explanations.
//! 
//! All metrics are computed deterministically and provide mathematical
//! evidence for anomaly detection that can be audited and reproduced.

use crate::compression::{CompressionAdapter, CompressionError};
use serde::{Deserialize, Serialize};
use std::fmt;

/// Complete anomaly detection metrics with glass-box explanation
/// 
/// This struct contains all the evidence needed to understand why
/// an anomaly was detected, providing the transparency required
/// for regulatory compliance and user confidence.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnomalyMetrics {
    /// Primary anomaly detection
    pub ncd: f64,                    // Normalized Compression Distance (0.0-1.0)
    pub p_value: f64,                // Statistical significance (0.0-1.0)
    pub is_anomaly: bool,            // Final determination
    
    /// Compression-based evidence
    pub baseline_compression_ratio: f64,   // How well baseline compresses
    pub window_compression_ratio: f64,     // How well window compresses  
    pub compression_ratio_change: f64,     // Percentage change
    
    /// Entropy evidence
    pub baseline_entropy: f64,       // Randomness of baseline (0.0-8.0)
    pub window_entropy: f64,         // Randomness of window (0.0-8.0)
    pub entropy_change: f64,         // Change in randomness
    
    /// Statistical confidence
    pub permutation_count: usize,    // How many permutations tested
    pub confidence_level: f64,       // 1.0 - p_value
    
    /// Human-readable explanation
    pub explanation: String,         // Generated glass-box explanation
}

impl AnomalyMetrics {
    /// Create a new metrics instance with default values
    pub fn new() -> Self {
        Self {
            ncd: 0.0,
            p_value: 1.0,
            is_anomaly: false,
            baseline_compression_ratio: 1.0,
            window_compression_ratio: 1.0,
            compression_ratio_change: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            entropy_change: 0.0,
            permutation_count: 0,
            confidence_level: 0.0,
            explanation: String::new(),
        }
    }

    /// Generate human-readable explanation based on metrics
    pub fn generate_explanation(&mut self) {
        let confidence_percent = (self.confidence_level * 100.0).round();
        let ratio_change_percent = (self.compression_ratio_change * 100.0).round();
        let entropy_change_percent = (self.entropy_change * 100.0).round();

        let mut explanation = format!(
            "Anomaly {} with {:.1}% confidence (p={:.3}):\n\n",
            if self.is_anomaly { "DETECTED" } else { "NOT DETECTED" },
            confidence_percent,
            self.p_value
        );

        explanation.push_str("COMPRESSION EVIDENCE:\n");
        explanation.push_str(&format!(
            "- Baseline compression ratio: {:.1}x (normal pattern)\n",
            self.baseline_compression_ratio
        ));
        explanation.push_str(&format!(
            "- Window compression ratio: {:.1}x (current pattern)\n",
            self.window_compression_ratio
        ));
        explanation.push_str(&format!(
            "- Change: {:+.0}% compression efficiency\n",
            ratio_change_percent
        ));

        explanation.push_str("\nENTROPY EVIDENCE:\n");
        explanation.push_str(&format!(
            "- Baseline entropy: {:.1} bits/byte (structured data)\n",
            self.baseline_entropy
        ));
        explanation.push_str(&format!(
            "- Window entropy: {:.1} bits/byte (current randomness)\n",
            self.window_entropy
        ));
        explanation.push_str(&format!(
            "- Change: {:+.0}% randomness\n",
            entropy_change_percent
        ));

        explanation.push_str(&format!("\nNCD SCORE: {:.2} (", self.ncd));
        if self.ncd < 0.3 {
            explanation.push_str("low dissimilarity");
        } else if self.ncd < 0.7 {
            explanation.push_str("moderate dissimilarity");
        } else {
            explanation.push_str("high dissimilarity");
        }
        explanation.push_str(")\n");

        // Add interpretation
        if self.is_anomaly {
            explanation.push_str("\nINTERPRETATION: ");
            if self.compression_ratio_change < -0.5 {
                explanation.push_str("Significant degradation in compression efficiency indicates unstructured or anomalous data patterns. ");
            }
            if self.entropy_change > 0.5 {
                explanation.push_str("Increased randomness suggests introduction of unexpected data structures. ");
            }
            if self.ncd > 0.7 {
                explanation.push_str("High NCD score indicates substantial dissimilarity from baseline patterns.");
            }
        } else {
            explanation.push_str("\nINTERPRETATION: Data patterns remain consistent with baseline expectations.");
        }

        self.explanation = explanation;
    }
}

impl Default for AnomalyMetrics {
    fn default() -> Self {
        Self::new()
    }
}

impl fmt::Display for AnomalyMetrics {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.explanation)
    }
}

/// Result type for metrics computation
pub type Result<T> = std::result::Result<T, MetricsError>;

/// Metrics computation errors
#[derive(Debug, Clone)]
pub enum MetricsError {
    /// Compression failed during metrics calculation
    CompressionFailed(String),
    /// Invalid input data
    InvalidInput(String),
    /// Mathematical computation error
    ComputationError(String),
}

impl fmt::Display for MetricsError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::CompressionFailed(msg) => write!(f, "Compression failed: {}", msg),
            Self::InvalidInput(msg) => write!(f, "Invalid input: {}", msg),
            Self::ComputationError(msg) => write!(f, "Computation error: {}", msg),
        }
    }
}

impl std::error::Error for MetricsError {}

impl From<CompressionError> for MetricsError {
    fn from(error: CompressionError) -> Self {
        Self::CompressionFailed(error.to_string())
    }
}

/// Individual metrics modules
pub mod compression_ratio;
pub mod delta_bits;
pub mod entropy;
pub mod ncd;

/// Main metrics computation function
/// 
/// Computes all anomaly detection metrics for baseline and window data.
/// This is the primary entry point for the CBAD algorithm.
pub fn compute_metrics(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    permutation_count: usize,
    seed: u64,
) -> Result<AnomalyMetrics> {
    let mut metrics = AnomalyMetrics::new();
    
    // Compute compression ratios
    let baseline_compressed = adapter.compress(baseline)?;
    let window_compressed = adapter.compress(window)?;
    
    metrics.baseline_compression_ratio = baseline.len() as f64 / baseline_compressed.len() as f64;
    metrics.window_compression_ratio = window.len() as f64 / window_compressed.len() as f64;
    metrics.compression_ratio_change = (metrics.window_compression_ratio - metrics.baseline_compression_ratio) / metrics.baseline_compression_ratio;

    // Compute entropy
    metrics.baseline_entropy = entropy::compute_entropy(baseline);
    metrics.window_entropy = entropy::compute_entropy(window);
    metrics.entropy_change = (metrics.window_entropy - metrics.baseline_entropy) / metrics.baseline_entropy.max(0.001); // Avoid division by zero

    // Compute NCD
    metrics.ncd = ncd::compute_ncd(baseline, window, adapter)?;

    // Perform permutation testing for statistical significance
    let perm_result = {
        let mut tester = permutation::PermutationTester::new(seed, permutation_count);
        tester.test_ncd_significance(baseline, window, adapter)?
    };
    
    metrics.p_value = perm_result.p_value;
    metrics.confidence_level = 1.0 - perm_result.p_value;
    metrics.permutation_count = permutation_count;
    
    // Determine if this is an anomaly (typically p < 0.05)
    metrics.is_anomaly = metrics.p_value < 0.05 && metrics.ncd > 0.3;

    // Generate human-readable explanation
    metrics.generate_explanation();

    Ok(metrics)
}

/// Permutation testing module (to be implemented)
pub mod permutation;

#[cfg(test)]
mod tests {
    use super::*;
    use crate::compression::create_adapter;
    use crate::compression::CompressionAlgorithm;

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
    fn test_metrics_computation() {
        let adapter = create_test_adapter();
        
        // Create test data - baseline with regular pattern
        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".to_vec();
        baseline.repeat(100).as_slice();
        
        // Window with anomalous data (stack trace)
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic stack_trace=\"thread 'main' panicked at 'index out of bounds', src/main.rs:42:13\"\n".to_vec();
        window.repeat(10).as_slice();

        let metrics = compute_metrics(
            &baseline,
            &window,
            adapter.as_ref(),
            100, // permutation count
            42,  // seed
        ).expect("compute metrics");

        // Verify metrics are computed
        assert!(metrics.baseline_compression_ratio > 1.0);
        assert!(metrics.window_compression_ratio > 1.0);
        assert!(metrics.baseline_entropy >= 0.0 && metrics.baseline_entropy <= 8.0);
        assert!(metrics.window_entropy >= 0.0 && metrics.window_entropy <= 8.0);
        assert!(metrics.ncd >= 0.0 && metrics.ncd <= 1.0);
        assert!(metrics.p_value >= 0.0 && metrics.p_value <= 1.0);
        
        // Should detect anomaly due to very different patterns
        assert!(metrics.is_anomaly);
        assert!(metrics.p_value < 0.05);
        
        println!("{}", metrics.explanation);
    }

    #[test]
    fn test_similar_data_not_anomaly() {
        let adapter = create_test_adapter();
        
        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".to_vec();
        baseline.repeat(100).as_slice();
        
        // Similar pattern - should not be anomaly
        let window = b"INFO 2025-10-24T00:00:01Z service=api-gateway msg=request_completed duration_ms=43\n".to_vec();
        window.repeat(10).as_slice();

        let metrics = compute_metrics(
            &baseline,
            &window,
            adapter.as_ref(),
            100,
            42,
        ).expect("compute metrics");

        // Similar data should not be flagged as anomaly
        assert!(!metrics.is_anomaly || metrics.p_value >= 0.05);
    }
}
