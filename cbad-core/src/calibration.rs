//! Calibration module for CBAD threshold tuning
//!
//! This module provides configurable weights and calibration state for
//! auto-tuning anomaly detection thresholds based on benchmark results
//! or streaming data characteristics.

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Configurable weights for composite score calculation.
///
/// The composite score blends multiple signals:
/// `composite = ncd_weight * NCD + p_value_weight * (1 - p_value) + compression_weight * compression_signal`
///
/// Default weights (0.5, 0.25, 0.25) were validated on PaySim fraud dataset achieving AUPRC=1.0.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CompositeWeights {
    /// Weight for NCD (Normalized Compression Distance) signal
    pub ncd_weight: f64,
    /// Weight for statistical significance (1 - p_value)
    pub p_value_weight: f64,
    /// Weight for compression ratio change signal
    pub compression_weight: f64,
}

impl Default for CompositeWeights {
    fn default() -> Self {
        Self {
            ncd_weight: 0.5,
            p_value_weight: 0.25,
            compression_weight: 0.25,
        }
    }
}

impl CompositeWeights {
    /// Validate that weights sum to approximately 1.0
    pub fn validate(&self) -> Result<(), String> {
        let sum = self.ncd_weight + self.p_value_weight + self.compression_weight;
        if (sum - 1.0).abs() > 0.01 {
            return Err(format!(
                "Weights must sum to 1.0, got {} (ncd={}, p_value={}, compression={})",
                sum, self.ncd_weight, self.p_value_weight, self.compression_weight
            ));
        }
        if self.ncd_weight < 0.0 || self.p_value_weight < 0.0 || self.compression_weight < 0.0 {
            return Err("All weights must be non-negative".to_string());
        }
        Ok(())
    }

    /// Compute composite score using these weights
    pub fn compute_score(&self, ncd: f64, p_value: f64, compression_signal: f64) -> f64 {
        self.ncd_weight * ncd
            + self.p_value_weight * (1.0 - p_value)
            + self.compression_weight * compression_signal
    }

    /// Create weights tuned for high-precision detection (fewer false positives)
    pub fn high_precision() -> Self {
        Self {
            ncd_weight: 0.4,
            p_value_weight: 0.4, // Weight statistical significance more
            compression_weight: 0.2,
        }
    }

    /// Create weights tuned for high-recall detection (catch more anomalies)
    pub fn high_recall() -> Self {
        Self {
            ncd_weight: 0.6, // Weight NCD more
            p_value_weight: 0.2,
            compression_weight: 0.2,
        }
    }
}

/// Method for calibrating detection thresholds
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CalibrationMethod {
    /// Target a specific false positive rate on normal data (unsupervised)
    /// Example: FprTarget(0.01) sets threshold at 99th percentile of normal scores
    FprTarget(f64),

    /// Maximize F1 score on labeled data (supervised)
    /// Requires ground truth labels during calibration
    F1Max,

    /// Use a manually specified threshold
    Manual(f64),
}

impl Default for CalibrationMethod {
    fn default() -> Self {
        CalibrationMethod::FprTarget(0.01)
    }
}

/// State for threshold calibration during warmup
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CalibrationState {
    /// Collected composite scores during warmup
    pub warmup_scores: Vec<f64>,

    /// Optional labels if provided (true = anomaly, false = normal)
    pub warmup_labels: Vec<Option<bool>>,

    /// Per-stream scores for stream-specific calibration
    #[serde(skip)]
    pub stream_scores: HashMap<String, Vec<(f64, Option<bool>)>>,

    /// Whether calibration has been completed
    pub is_calibrated: bool,

    /// The calibrated threshold (if calibration completed)
    pub calibrated_threshold: Option<f64>,

    /// Per-stream calibrated thresholds
    pub stream_thresholds: HashMap<String, f64>,

    /// Method used for calibration
    pub calibration_method: CalibrationMethod,

    /// Minimum samples required before calibration
    pub min_samples: usize,
}

impl Default for CalibrationState {
    fn default() -> Self {
        Self {
            warmup_scores: Vec::new(),
            warmup_labels: Vec::new(),
            stream_scores: HashMap::new(),
            is_calibrated: false,
            calibrated_threshold: None,
            stream_thresholds: HashMap::new(),
            calibration_method: CalibrationMethod::default(),
            min_samples: 100,
        }
    }
}

impl CalibrationState {
    /// Create a new calibration state with specified method
    pub fn new(method: CalibrationMethod, min_samples: usize) -> Self {
        Self {
            calibration_method: method,
            min_samples,
            ..Default::default()
        }
    }

    /// Record a score during warmup
    pub fn record_score(&mut self, score: f64, label: Option<bool>, stream_key: Option<&str>) {
        self.warmup_scores.push(score);
        self.warmup_labels.push(label);

        if let Some(key) = stream_key {
            self.stream_scores
                .entry(key.to_string())
                .or_default()
                .push((score, label));
        }
    }

    /// Check if we have enough samples for calibration
    pub fn ready_for_calibration(&self) -> bool {
        self.warmup_scores.len() >= self.min_samples
    }

    /// Calibrate threshold based on collected scores
    pub fn calibrate(&mut self) -> Option<f64> {
        if !self.ready_for_calibration() {
            return None;
        }

        let threshold = match &self.calibration_method {
            CalibrationMethod::FprTarget(target_fpr) => self.calibrate_by_fpr(*target_fpr),
            CalibrationMethod::F1Max => self.calibrate_by_f1(),
            CalibrationMethod::Manual(threshold) => Some(*threshold),
        };

        if let Some(t) = threshold {
            self.calibrated_threshold = Some(t);
            self.is_calibrated = true;
        }

        threshold
    }

    /// Calibrate per-stream thresholds using FPR method
    pub fn calibrate_streams(&mut self, target_fpr: f64) -> HashMap<String, f64> {
        let mut thresholds = HashMap::new();

        for (stream_key, scores) in &self.stream_scores {
            if scores.len() < 20 {
                continue; // Skip streams with too few samples
            }

            // Extract normal scores (where label is None or Some(false))
            let normal_scores: Vec<f64> = scores
                .iter()
                .filter(|(_, label)| label.map(|l| !l).unwrap_or(true))
                .map(|(score, _)| *score)
                .collect();

            if normal_scores.len() < 10 {
                continue;
            }

            let threshold = Self::compute_fpr_threshold(&normal_scores, target_fpr);
            thresholds.insert(stream_key.clone(), threshold);
        }

        self.stream_thresholds = thresholds.clone();
        thresholds
    }

    /// Calibrate threshold targeting a false positive rate
    fn calibrate_by_fpr(&self, target_fpr: f64) -> Option<f64> {
        // Extract normal scores (where label is None or Some(false))
        let normal_scores: Vec<f64> = self
            .warmup_scores
            .iter()
            .zip(&self.warmup_labels)
            .filter(|(_, label)| label.map(|l| !l).unwrap_or(true))
            .map(|(score, _)| *score)
            .collect();

        if normal_scores.is_empty() {
            return None;
        }

        Some(Self::compute_fpr_threshold(&normal_scores, target_fpr))
    }

    /// Compute threshold at (1 - target_fpr) quantile
    fn compute_fpr_threshold(scores: &[f64], target_fpr: f64) -> f64 {
        let mut sorted: Vec<f64> = scores.to_vec();
        sorted.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));

        let quantile = 1.0 - target_fpr.clamp(0.001, 0.5);
        let idx = ((sorted.len() as f64 - 1.0) * quantile).round() as usize;
        sorted[idx.min(sorted.len() - 1)]
    }

    /// Calibrate threshold by maximizing F1 score (requires labels)
    fn calibrate_by_f1(&self) -> Option<f64> {
        // Need labels for F1 calibration
        let labeled_scores: Vec<(f64, bool)> = self
            .warmup_scores
            .iter()
            .zip(&self.warmup_labels)
            .filter_map(|(score, label)| label.map(|l| (*score, l)))
            .collect();

        if labeled_scores.len() < 10 {
            return None; // Not enough labeled data
        }

        // Try thresholds from 0.1 to 0.95 in steps of 0.01
        let mut best_threshold = 0.5;
        let mut best_f1 = 0.0;

        for threshold_pct in 10..=95 {
            let threshold = threshold_pct as f64 / 100.0;

            let (mut tp, mut fp, mut fn_) = (0, 0, 0);
            for (score, is_anomaly) in &labeled_scores {
                let predicted = *score >= threshold;
                match (predicted, *is_anomaly) {
                    (true, true) => tp += 1,
                    (true, false) => fp += 1,
                    (false, true) => fn_ += 1,
                    (false, false) => {}
                }
            }

            let precision = if tp + fp > 0 {
                tp as f64 / (tp + fp) as f64
            } else {
                0.0
            };
            let recall = if tp + fn_ > 0 {
                tp as f64 / (tp + fn_) as f64
            } else {
                0.0
            };
            let f1 = if precision + recall > 0.0 {
                2.0 * precision * recall / (precision + recall)
            } else {
                0.0
            };

            if f1 > best_f1 {
                best_f1 = f1;
                best_threshold = threshold;
            }
        }

        Some(best_threshold)
    }

    /// Reset calibration state
    pub fn reset(&mut self) {
        self.warmup_scores.clear();
        self.warmup_labels.clear();
        self.stream_scores.clear();
        self.is_calibrated = false;
        self.calibrated_threshold = None;
        self.stream_thresholds.clear();
    }

    /// Get the number of collected samples
    pub fn sample_count(&self) -> usize {
        self.warmup_scores.len()
    }

    /// Get statistics about collected scores
    pub fn score_statistics(&self) -> Option<ScoreStatistics> {
        if self.warmup_scores.is_empty() {
            return None;
        }

        let n = self.warmup_scores.len() as f64;
        let mean = self.warmup_scores.iter().sum::<f64>() / n;
        let variance = self
            .warmup_scores
            .iter()
            .map(|x| (x - mean).powi(2))
            .sum::<f64>()
            / n;
        let stddev = variance.sqrt();

        let mut sorted = self.warmup_scores.clone();
        sorted.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));

        let min = sorted[0];
        let max = sorted[sorted.len() - 1];
        let median = sorted[sorted.len() / 2];
        let p95 = sorted[(sorted.len() as f64 * 0.95) as usize];
        let p99 = sorted[(sorted.len() as f64 * 0.99).min(sorted.len() as f64 - 1.0) as usize];

        Some(ScoreStatistics {
            count: self.warmup_scores.len(),
            mean,
            stddev,
            min,
            max,
            median,
            p95,
            p99,
        })
    }
}

/// Statistics about collected calibration scores
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ScoreStatistics {
    pub count: usize,
    pub mean: f64,
    pub stddev: f64,
    pub min: f64,
    pub max: f64,
    pub median: f64,
    pub p95: f64,
    pub p99: f64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_composite_weights_default() {
        let weights = CompositeWeights::default();
        assert!((weights.ncd_weight - 0.5).abs() < 0.001);
        assert!((weights.p_value_weight - 0.25).abs() < 0.001);
        assert!((weights.compression_weight - 0.25).abs() < 0.001);
        assert!(weights.validate().is_ok());
    }

    #[test]
    fn test_composite_weights_validation() {
        let invalid = CompositeWeights {
            ncd_weight: 0.5,
            p_value_weight: 0.5,
            compression_weight: 0.5, // Sum = 1.5
        };
        assert!(invalid.validate().is_err());

        let negative = CompositeWeights {
            ncd_weight: -0.1,
            p_value_weight: 0.6,
            compression_weight: 0.5,
        };
        assert!(negative.validate().is_err());
    }

    #[test]
    fn test_composite_score_computation() {
        let weights = CompositeWeights::default();
        // NCD=0.5, p_value=0.1, compression_signal=0.2
        // Expected: 0.5*0.5 + 0.25*(1-0.1) + 0.25*0.2 = 0.25 + 0.225 + 0.05 = 0.525
        let score = weights.compute_score(0.5, 0.1, 0.2);
        assert!((score - 0.525).abs() < 0.001);
    }

    #[test]
    fn test_calibration_by_fpr() {
        let mut state = CalibrationState::new(CalibrationMethod::FprTarget(0.1), 10);

        // Add 100 normal scores from 0.0 to 0.99
        for i in 0..100 {
            state.record_score(i as f64 / 100.0, Some(false), None);
        }

        let threshold = state.calibrate().unwrap();
        // At 10% FPR, threshold should be at ~90th percentile = 0.9
        assert!(threshold > 0.85 && threshold < 0.95, "threshold = {}", threshold);
    }

    #[test]
    fn test_calibration_by_f1() {
        let mut state = CalibrationState::new(CalibrationMethod::F1Max, 10);

        // Add labeled samples: normals have low scores, anomalies have high scores
        for i in 0..50 {
            state.record_score(i as f64 / 100.0, Some(false), None); // normal: 0.0-0.49
        }
        for i in 50..100 {
            state.record_score(i as f64 / 100.0, Some(true), None); // anomaly: 0.5-0.99
        }

        let threshold = state.calibrate().unwrap();
        // Perfect separation at 0.5, threshold should be around 0.5
        assert!(threshold > 0.4 && threshold < 0.6, "threshold = {}", threshold);
    }

    #[test]
    fn test_score_statistics() {
        let mut state = CalibrationState::default();
        for i in 0..100 {
            state.record_score(i as f64 / 100.0, None, None);
        }

        let stats = state.score_statistics().unwrap();
        assert_eq!(stats.count, 100);
        assert!((stats.mean - 0.495).abs() < 0.01);
        assert!((stats.min - 0.0).abs() < 0.001);
        assert!((stats.max - 0.99).abs() < 0.001);
    }
}
