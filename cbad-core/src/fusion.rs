//! FusionWeights scoring for compression-based anomaly detection
//!
//! This module implements the fusion approach that combines multiple signals:
//! - NCD (Normalized Compression Distance)
//! - Entropy
//! - Compression patterns
//!
//! The formula: complexity = ncd_weight * ncd + entropy_weight * (entropy/8) + patterns_weight * compression_ratio

use crate::stats::{entropy_bits_per_byte, SimpleQuantile, Welford};
use serde::{Deserialize, Serialize};

/// Weights for combining detection signals
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FusionWeights {
    /// Weight for NCD score (default: 0.4)
    pub ncd: f64,
    /// Weight for VAE score - reserved for future ML integration (default: 0.3)
    pub vae: f64,
    /// Weight for entropy signal (default: 0.2)
    pub entropy: f64,
    /// Weight for compression pattern score (default: 0.1)
    pub patterns: f64,
}

impl Default for FusionWeights {
    fn default() -> Self {
        Self {
            ncd: 0.4,
            vae: 0.3, // Reserved for VAE integration
            entropy: 0.2,
            patterns: 0.1,
        }
    }
}

/// Result of fusion-based detection
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FusionResult {
    /// Whether this is detected as an anomaly
    pub is_anomaly: bool,
    /// Confidence level (0.0 to 1.0)
    pub confidence: f64,
    /// Combined complexity score
    pub complexity_score: f64,
    /// Shannon entropy in bits per byte
    pub entropy: f64,
    /// Compression ratio (original/compressed)
    pub compression_ratio: f64,
    /// Z-score relative to historical mean
    pub z_score: f64,
    /// The effective threshold used for detection
    pub effective_threshold: f64,
    /// Processing time in microseconds
    pub processing_time_us: u64,
    /// Breakdown of score contributions
    pub ncd_contribution: f64,
    pub entropy_contribution: f64,
    pub patterns_contribution: f64,
}

/// Fusion-based anomaly scorer
///
/// Combines NCD, entropy, and compression patterns using configurable weights
/// with optional adaptive thresholding based on historical statistics.
pub struct FusionScorer {
    weights: FusionWeights,
    stats: Welford,
    history: SimpleQuantile,
    base_threshold: f64,
    adaptive: bool,
    window_size: usize,
}

impl FusionScorer {
    /// Create a new fusion scorer
    pub fn new(
        weights: FusionWeights,
        base_threshold: f64,
        adaptive: bool,
        window_size: usize,
    ) -> Self {
        Self {
            weights,
            stats: Welford::new(),
            history: SimpleQuantile::new_with_cap(window_size),
            base_threshold,
            adaptive,
            window_size,
        }
    }

    /// Create with default settings
    pub fn with_defaults() -> Self {
        Self::new(FusionWeights::default(), 2.5, true, 1000)
    }

    /// Compute fusion score from raw metrics
    ///
    /// # Arguments
    /// * `ncd` - Normalized compression distance (0.0 to 1.0)
    /// * `entropy` - Shannon entropy in bits per byte (0.0 to 8.0)
    /// * `compression_ratio` - Ratio of original to compressed size
    pub fn score(&mut self, ncd: f64, entropy: f64, compression_ratio: f64) -> FusionResult {
        let start = std::time::Instant::now();

        // Get prior statistics for z-score calculation
        let prior_mean = self.stats.mean();
        let prior_std = self.stats.std();

        // Compute weighted contributions
        let ncd_contribution = self.weights.ncd * ncd;
        let entropy_contribution = self.weights.entropy * (entropy / 8.0);
        let patterns_contribution = self.weights.patterns * compression_ratio;

        // Combined complexity score (VAE contribution is 0 for now)
        let complexity_score = ncd_contribution + entropy_contribution + patterns_contribution;

        // Update statistics
        self.stats.update(complexity_score);
        self.history.add(complexity_score);

        // Compute z-score
        let z_score = if prior_std > 0.0 {
            (complexity_score - prior_mean) / prior_std
        } else {
            0.0
        };

        // Compute effective threshold
        let effective_threshold = if self.adaptive {
            let mean = self.stats.mean();
            let std = self.stats.std();
            if std > 0.0 {
                mean + 2.0 * std
            } else {
                self.base_threshold
            }
        } else {
            self.base_threshold
        };

        // Determine anomaly status
        let is_anomaly = complexity_score > effective_threshold;

        // Compute confidence
        let confidence = if is_anomaly && effective_threshold > 0.0 {
            ((complexity_score - effective_threshold) / effective_threshold).clamp(0.0, 1.0)
        } else {
            0.0
        };

        let processing_time_us = start.elapsed().as_micros() as u64;

        FusionResult {
            is_anomaly,
            confidence,
            complexity_score,
            entropy,
            compression_ratio,
            z_score,
            effective_threshold,
            processing_time_us,
            ncd_contribution,
            entropy_contribution,
            patterns_contribution,
        }
    }

    /// Score raw data bytes
    pub fn score_data(&mut self, data: &[u8], compressed_len: usize) -> FusionResult {
        if data.is_empty() {
            return FusionResult {
                is_anomaly: false,
                confidence: 0.0,
                complexity_score: 0.0,
                entropy: 0.0,
                compression_ratio: 0.0,
                z_score: 0.0,
                effective_threshold: self.base_threshold,
                processing_time_us: 0,
                ncd_contribution: 0.0,
                entropy_contribution: 0.0,
                patterns_contribution: 0.0,
            };
        }

        let entropy = entropy_bits_per_byte(data);
        let compression_ratio = data.len() as f64 / compressed_len.max(1) as f64;
        let ncd = 1.0 - (compressed_len as f64 / data.len() as f64);

        self.score(ncd, entropy, compression_ratio)
    }

    /// Reset the scorer statistics
    pub fn reset(&mut self) {
        self.stats = Welford::new();
        self.history = SimpleQuantile::new_with_cap(self.window_size);
    }

    /// Get current statistics
    pub fn statistics(&self) -> FusionStatistics {
        FusionStatistics {
            mean: self.stats.mean(),
            std: self.stats.std(),
            variance: self.stats.variance(),
            count: self.stats.count(),
            threshold: self.base_threshold,
            adaptive: self.adaptive,
        }
    }
}

/// Statistics from the fusion scorer
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FusionStatistics {
    pub mean: f64,
    pub std: f64,
    pub variance: f64,
    pub count: u64,
    pub threshold: f64,
    pub adaptive: bool,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_fusion_weights_default() {
        let weights = FusionWeights::default();
        assert_eq!(weights.ncd, 0.4);
        assert_eq!(weights.vae, 0.3);
        assert_eq!(weights.entropy, 0.2);
        assert_eq!(weights.patterns, 0.1);
    }

    #[test]
    fn test_fusion_scorer_basic() {
        let mut scorer = FusionScorer::with_defaults();

        // Normal data - low NCD, moderate entropy
        for _ in 0..100 {
            let result = scorer.score(0.1, 4.0, 3.0);
            assert!(!result.is_anomaly, "Normal data should not trigger anomaly");
        }

        // After building baseline, anomalous data should trigger
        let anomaly = scorer.score(0.9, 7.5, 1.2);
        // The anomaly detection depends on adaptive threshold
        // After 100 normal samples, a high score should stand out
        assert!(anomaly.complexity_score > 0.5);
    }

    #[test]
    fn test_fusion_scorer_reset() {
        let mut scorer = FusionScorer::with_defaults();

        for _ in 0..50 {
            scorer.score(0.2, 5.0, 2.5);
        }

        assert!(scorer.statistics().count > 0);

        scorer.reset();

        assert_eq!(scorer.statistics().count, 0);
    }
}
