//! Metrics module for CBAD (Compression-Based Anomaly Detection)
//!
//! This module provides the core anomaly detection metrics that give users
//! confidence in anomaly detection through glass-box explanations.
//!
//! All metrics are computed deterministically and provide mathematical
//! evidence for anomaly detection that can be audited and reproduced.

use crate::calibration::CompositeWeights;
use crate::compression::{CompressionAdapter, CompressionError};
use crate::tokenizer::Tokenizer;
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
    pub ncd: f64, // Normalized Compression Distance (0.0-1.0)
    pub p_value: f64,     // Statistical significance (0.0-1.0)
    pub is_anomaly: bool, // Final determination

    #[serde(default)]
    pub composite_score: f64,

    #[serde(default)]
    pub conditional_novelty: f64,

    /// Compression-based evidence
    pub baseline_compression_ratio: f64, // How well baseline compresses
    pub window_compression_ratio: f64, // How well window compresses
    pub compression_ratio_change: f64, // Percentage change

    /// Entropy evidence
    pub baseline_entropy: f64, // Randomness of baseline (0.0-8.0)
    pub window_entropy: f64, // Randomness of window (0.0-8.0)
    pub entropy_change: f64, // Change in randomness

    /// Statistical confidence
    pub permutation_count: usize, // How many permutations tested
    pub confidence_level: f64, // 1.0 - p_value

    /// Human-readable explanation
    pub explanation: String, // Generated glass-box explanation

    /// Recommendations derived from observed metrics
    pub recommended_ncd_threshold: f64,
    pub recommended_window_size: usize,
    pub data_stability_score: f64, // 0.0 (unstable/noisy) to 1.0 (stable/consistent)
}

impl AnomalyMetrics {
    /// Create a new metrics instance with default values
    pub fn new() -> Self {
        Self {
            ncd: 0.0,
            p_value: 1.0,
            is_anomaly: false,

            composite_score: 0.0,
            conditional_novelty: 0.0,
            baseline_compression_ratio: 1.0,
            window_compression_ratio: 1.0,
            compression_ratio_change: 0.0,
            baseline_entropy: 0.0,
            window_entropy: 0.0,
            entropy_change: 0.0,
            permutation_count: 0,
            confidence_level: 0.0,
            explanation: String::new(),
            recommended_ncd_threshold: 0.0,
            recommended_window_size: 0,
            data_stability_score: 0.0,
        }
    }

    /// Generate human-readable explanation based on metrics
    pub fn generate_explanation(&mut self) {
        let confidence_percent = (self.confidence_level * 100.0).round();
        let ratio_change_percent = (self.compression_ratio_change * 100.0).round();
        let entropy_change_percent = (self.entropy_change * 100.0).round();

        let mut explanation = format!(
            "Anomaly {} with {:.1}% confidence (p={:.3}):\n\n",
            if self.is_anomaly {
                "DETECTED"
            } else {
                "NOT DETECTED"
            },
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

        explanation.push_str(&format!(
            "\nCOMPOSITE SCORE: {:.3} (ncd/p-value/compression fusion)\n",
            self.composite_score
        ));

        explanation.push_str(&format!(
            "\nCONDITIONAL NOVELTY: {:.3} ((C(b+w)-C(b))/C(w)) lower=more expected)\n",
            self.conditional_novelty
        ));

        // Add interpretation
        if self.is_anomaly {
            explanation.push_str("\nINTERPRETATION: ");
            if self.compression_ratio_change < -0.5 {
                explanation.push_str("Significant degradation in compression efficiency indicates unstructured or anomalous data patterns. ");
            }
            if self.entropy_change > 0.5 {
                explanation.push_str(
                    "Increased randomness suggests introduction of unexpected data structures. ",
                );
            }
            if self.ncd > 0.7 {
                explanation.push_str(
                    "High NCD score indicates substantial dissimilarity from baseline patterns.",
                );
            }
        } else {
            explanation.push_str(
                "\nINTERPRETATION: Data patterns remain consistent with baseline expectations.",
            );
        }

        self.explanation = explanation;
    }

    /// Compute tuning recommendations based on observed stability and current configuration.
    pub fn apply_recommendations(
        &mut self,
        current_ncd_threshold: f64,
        current_window_size: usize,
        baseline_size: usize,
    ) {
        // Estimate stability from entropy/compression variability and confidence.
        let entropy_instability = self.entropy_change.abs().min(1.0);
        let compression_instability = self.compression_ratio_change.abs().min(1.0);
        let confidence_instability = (1.0 - self.confidence_level).abs().min(1.0);

        let stability = (1.0
            - (0.4 * entropy_instability
                + 0.4 * compression_instability
                + 0.2 * confidence_instability))
            .clamp(0.0, 1.0);
        self.data_stability_score = stability;

        // Recommend NCD threshold: raise slightly on noisy streams, lower on stable ones.
        let base_ncd = if current_ncd_threshold > 0.0 {
            current_ncd_threshold
        } else {
            0.3
        };
        self.recommended_ncd_threshold = (base_ncd + (0.5 - stability) * 0.2).clamp(0.1, 0.8);

        // Recommend window size: larger for unstable/high-entropy streams, smaller for very stable ones.
        let base_window = if current_window_size > 0 {
            current_window_size
        } else {
            50
        };
        let baseline_cap = if baseline_size > 0 {
            baseline_size as f64 * 0.75
        } else {
            (base_window as f64) * 3.0
        };
        let window_scale = 1.0 + (1.0 - stability) * 0.5 + entropy_instability * 0.25;
        let proposed = (base_window as f64 * window_scale).clamp(10.0, baseline_cap.max(10.0));
        self.recommended_window_size = proposed.round() as usize;
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
///
/// If a tokenizer is provided, high-entropy fields (UUIDs, hashes, JWTs, Base64)
/// will be replaced with fixed tokens before compression. This improves detection
/// accuracy across different data types by normalizing random-looking noise.
pub fn compute_metrics(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    permutation_count: usize,
    seed: u64,
) -> Result<AnomalyMetrics> {
    compute_metrics_with_tokenizer(baseline, window, adapter, permutation_count, seed, None)
}

/// Metrics computation with optional tokenization preprocessing
///
/// This variant allows passing a tokenizer to normalize high-entropy fields
/// before compression analysis. Use this for type-agnostic anomaly detection
/// across different data formats.
pub fn compute_metrics_with_tokenizer(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    permutation_count: usize,
    seed: u64,
    tokenizer: Option<&Tokenizer>,
) -> Result<AnomalyMetrics> {
    let mut metrics = AnomalyMetrics::new();

    // Apply tokenization if enabled (normalizes UUIDs, hashes, JWTs, Base64)
    let (baseline_data, window_data): (std::borrow::Cow<[u8]>, std::borrow::Cow<[u8]>) =
        match tokenizer {
            Some(t) => (
                std::borrow::Cow::Owned(t.tokenize(baseline)),
                std::borrow::Cow::Owned(t.tokenize(window)),
            ),
            None => (
                std::borrow::Cow::Borrowed(baseline),
                std::borrow::Cow::Borrowed(window),
            ),
        };

    // Compute compression ratios (reuse compressed sizes for NCD to avoid double work)
    let baseline_compressed = adapter.compress(&baseline_data)?;
    let window_compressed = adapter.compress(&window_data)?;

    metrics.baseline_compression_ratio =
        baseline_data.len() as f64 / baseline_compressed.len() as f64;
    metrics.window_compression_ratio = window_data.len() as f64 / window_compressed.len() as f64;
    metrics.compression_ratio_change = (metrics.window_compression_ratio
        - metrics.baseline_compression_ratio)
        / metrics.baseline_compression_ratio;

    // Compute entropy (on tokenized data if applicable)
    metrics.baseline_entropy = entropy::compute_entropy(&baseline_data);
    metrics.window_entropy = entropy::compute_entropy(&window_data);
    metrics.entropy_change =
        (metrics.window_entropy - metrics.baseline_entropy) / metrics.baseline_entropy.max(0.001); // Avoid division by zero

    // Compute NCD using already compressed sizes to avoid recompressing baseline/window
    let mut combined = Vec::with_capacity(baseline_data.len() + window_data.len());
    combined.extend_from_slice(&baseline_data);
    combined.extend_from_slice(&window_data);
    let combined_compressed = adapter.compress(&combined)?;
    metrics.ncd = ncd::compute_ncd_from_sizes(
        baseline_compressed.len(),
        window_compressed.len(),
        combined_compressed.len(),
    );

    let delta_conditional = combined_compressed
        .len()
        .saturating_sub(baseline_compressed.len()) as f64;
    metrics.conditional_novelty = if window_compressed.is_empty() {
        0.0
    } else {
        (delta_conditional / (window_compressed.len() as f64)).max(0.0)
    };

    // Perform permutation testing for statistical significance (on tokenized data)
    let perm_result = {
        let mut tester = permutation::PermutationTester::new(seed, permutation_count);
        tester.test_ncd_significance(&baseline_data, &window_data, adapter)?
    };

    metrics.p_value = perm_result.p_value;
    // Calculate confidence level that considers both NCD and p-value
    // Strategy:
    // 1. If statistically significant (p < 0.05): use p-value based confidence (high confidence)
    // 2. If not significant but NCD is high (> 0.5): use NCD as confidence (there's a clear difference)
    // 3. Otherwise: use weighted combination
    let p_value_confidence = 1.0 - perm_result.p_value;
    let ncd_confidence = metrics.ncd.min(1.0); // NCD is already 0-1, use it directly

    if perm_result.p_value < 0.05 {
        // Statistically significant: use p-value based confidence
        metrics.confidence_level = p_value_confidence;
    } else if metrics.ncd > 0.5 {
        // Not statistically significant but high NCD indicates clear difference
        // Use NCD as confidence, but cap it to show it's not statistically proven
        metrics.confidence_level = (ncd_confidence * 0.75).min(0.85); // Cap at 85% when not statistically significant
    } else {
        // Low NCD and not significant: use weighted combination
        metrics.confidence_level =
            (p_value_confidence * 0.6 + ncd_confidence * 0.4).clamp(0.0, 1.0);
    }
    metrics.permutation_count = permutation_count;

    let compression_signal = (-metrics.compression_ratio_change).max(0.0);
    // Use default weights (validated on PaySim: AUPRC=1.0)
    let default_weights = CompositeWeights::default();
    metrics.composite_score =
        default_weights.compute_score(metrics.ncd, metrics.p_value, compression_signal);

    // Apply a conservative default anomaly heuristic for callers that don't supply config.
    let default_p = 0.05;
    let default_ncd = 0.3;
    let default_drop = 0.15;
    let default_entropy = 0.2;
    metrics.is_anomaly = (metrics.ncd >= default_ncd && metrics.p_value <= default_p)
        || (metrics.compression_ratio_change <= -default_drop && metrics.ncd >= 0.2)
        || metrics.entropy_change >= default_entropy;

    // Generate human-readable explanation
    metrics.generate_explanation();

    Ok(metrics)
}

/// Metrics computation with configurable composite weights
///
/// This variant allows passing custom weights for the composite score calculation,
/// enabling profile-specific or calibrated scoring.
pub fn compute_metrics_with_weights(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    permutation_count: usize,
    seed: u64,
    tokenizer: Option<&Tokenizer>,
    weights: &CompositeWeights,
) -> Result<AnomalyMetrics> {
    let mut metrics = AnomalyMetrics::new();

    // Apply tokenization if enabled (normalizes UUIDs, hashes, JWTs, Base64)
    let (baseline_data, window_data): (std::borrow::Cow<[u8]>, std::borrow::Cow<[u8]>) =
        match tokenizer {
            Some(t) => (
                std::borrow::Cow::Owned(t.tokenize(baseline)),
                std::borrow::Cow::Owned(t.tokenize(window)),
            ),
            None => (
                std::borrow::Cow::Borrowed(baseline),
                std::borrow::Cow::Borrowed(window),
            ),
        };

    // Compute compression ratios (reuse compressed sizes for NCD to avoid double work)
    let baseline_compressed = adapter.compress(&baseline_data)?;
    let window_compressed = adapter.compress(&window_data)?;

    metrics.baseline_compression_ratio =
        baseline_data.len() as f64 / baseline_compressed.len() as f64;
    metrics.window_compression_ratio = window_data.len() as f64 / window_compressed.len() as f64;
    metrics.compression_ratio_change = (metrics.window_compression_ratio
        - metrics.baseline_compression_ratio)
        / metrics.baseline_compression_ratio;

    // Compute entropy (on tokenized data if applicable)
    metrics.baseline_entropy = entropy::compute_entropy(&baseline_data);
    metrics.window_entropy = entropy::compute_entropy(&window_data);
    metrics.entropy_change =
        (metrics.window_entropy - metrics.baseline_entropy) / metrics.baseline_entropy.max(0.001);

    // Compute NCD using already compressed sizes to avoid recompressing baseline/window
    let mut combined = Vec::with_capacity(baseline_data.len() + window_data.len());
    combined.extend_from_slice(&baseline_data);
    combined.extend_from_slice(&window_data);
    let combined_compressed = adapter.compress(&combined)?;
    metrics.ncd = ncd::compute_ncd_from_sizes(
        baseline_compressed.len(),
        window_compressed.len(),
        combined_compressed.len(),
    );

    let delta_conditional = combined_compressed
        .len()
        .saturating_sub(baseline_compressed.len()) as f64;
    metrics.conditional_novelty = if window_compressed.is_empty() {
        0.0
    } else {
        (delta_conditional / (window_compressed.len() as f64)).max(0.0)
    };

    // Perform permutation testing for statistical significance (on tokenized data)
    let perm_result = {
        let mut tester = permutation::PermutationTester::new(seed, permutation_count);
        tester.test_ncd_significance(&baseline_data, &window_data, adapter)?
    };

    metrics.p_value = perm_result.p_value;
    // Calculate confidence level that considers both NCD and p-value
    let p_value_confidence = 1.0 - perm_result.p_value;
    let ncd_confidence = metrics.ncd.min(1.0);

    if perm_result.p_value < 0.05 {
        metrics.confidence_level = p_value_confidence;
    } else if metrics.ncd > 0.5 {
        metrics.confidence_level = (ncd_confidence * 0.75).min(0.85);
    } else {
        metrics.confidence_level =
            (p_value_confidence * 0.6 + ncd_confidence * 0.4).clamp(0.0, 1.0);
    }
    metrics.permutation_count = permutation_count;

    // Compute composite score using provided weights
    let compression_signal = (-metrics.compression_ratio_change).max(0.0);
    metrics.composite_score = weights.compute_score(metrics.ncd, metrics.p_value, compression_signal);

    // Apply a conservative default anomaly heuristic
    let default_p = 0.05;
    let default_ncd = 0.3;
    let default_drop = 0.15;
    let default_entropy = 0.2;
    metrics.is_anomaly = (metrics.ncd >= default_ncd && metrics.p_value <= default_p)
        || (metrics.compression_ratio_change <= -default_drop && metrics.ncd >= 0.2)
        || metrics.entropy_change >= default_entropy;

    // Generate human-readable explanation
    metrics.generate_explanation();

    Ok(metrics)
}

/// Compute only the composite score for a given set of raw metrics.
///
/// Useful for batch scoring during calibration when you already have
/// the individual metric values.
pub fn compute_composite_score(
    ncd: f64,
    p_value: f64,
    compression_ratio_change: f64,
    weights: &CompositeWeights,
) -> f64 {
    let compression_signal = (-compression_ratio_change).max(0.0);
    weights.compute_score(ncd, p_value, compression_signal)
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

        let baseline_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#;
        let anomalous_log = r#"{"timestamp":"2025-10-24T00:00:01Z","severity":"ERROR","service":"api-gateway","message":"Panic occurred","attributes":{"stack_trace":"0x3fa8d1b2c9e47f56::panic::trace::[ns=923847923847923847923847]::random_payload=9fjK2L1pQwZ8xT4rB7nC6Mv0HdYG5s2tR1uQ3w8yAaEeIiOo","binary_blob":"Q29tcHJlc3NlZEJsb2I6ZGV0ZXJtaW5pc3RpY1Nob3J0c0FuZFJhbmRvbVVuaWNvZGVEYXRh"}}"#;

        // Create test data - baseline with regular pattern
        let baseline_entry = {
            let mut entry = baseline_log.to_owned();
            entry.push('\n');
            entry.into_bytes()
        };
        let baseline = baseline_entry.repeat(200);

        // Window with anomalous data (stack trace + binary blob)
        let window_entry = {
            let mut entry = anomalous_log.to_owned();
            entry.push('\n');
            entry.into_bytes()
        };
        let window = window_entry.repeat(200);

        let mut metrics = compute_metrics(
            &baseline,
            &window,
            adapter.as_ref(),
            100, // permutation count
            42,  // seed
        )
        .expect("compute metrics");

        metrics.generate_explanation();

        // Verify metrics are computed
        assert!(metrics.baseline_compression_ratio > 1.0);
        assert!(metrics.window_compression_ratio > 1.0);
        assert!(metrics.baseline_entropy >= 0.0 && metrics.baseline_entropy <= 8.0);
        assert!(metrics.window_entropy >= 0.0 && metrics.window_entropy <= 8.0);
        assert!(metrics.ncd >= 0.0 && metrics.ncd <= 1.0);
        assert!(metrics.p_value >= 0.0 && metrics.p_value <= 1.0);

        println!("{}", metrics.explanation);
        println!("ncd={:.3} p_value={:.3}", metrics.ncd, metrics.p_value);

        // Should be strong anomaly signals
        assert!(metrics.ncd > 0.5);
        assert!(metrics.p_value < 0.05);
    }

    #[test]
    fn test_similar_data_not_anomaly() {
        let adapter = create_test_adapter();

        let baseline_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#;
        let similar_log = r#"{"timestamp":"2025-10-24T00:00:01Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":45}}"#;

        let baseline_entry = {
            let mut entry = baseline_log.to_owned();
            entry.push('\n');
            entry.into_bytes()
        };
        let baseline = baseline_entry.repeat(200);

        // Similar pattern - should not be anomaly
        let window_entry = {
            let mut entry = similar_log.to_owned();
            entry.push('\n');
            entry.into_bytes()
        };
        let window = window_entry.repeat(200);

        let mut metrics = compute_metrics(&baseline, &window, adapter.as_ref(), 100, 42)
            .expect("compute metrics");

        metrics.generate_explanation();

        // Similar data should yield low NCD and low confidence
        assert!(metrics.ncd < 0.5);
        assert!(metrics.p_value >= 0.01);
    }
}
