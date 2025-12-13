//! Calibration models for DB-driven threshold tuning

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

/// Profile calibration settings - default thresholds per detection profile
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct ProfileCalibration {
    pub id: Uuid,
    pub profile_name: String,

    // Detection thresholds
    pub ncd_threshold: f64,
    pub p_value_threshold: f64,
    pub composite_threshold: f64,

    // Composite score weights
    pub ncd_weight: f64,
    pub p_value_weight: f64,
    pub compression_weight: f64,

    // Window configuration
    pub baseline_size: i32,
    pub window_size: i32,
    pub permutation_count: i32,

    // Adaptive settings
    pub adaptive_target_fpr: Option<f64>,
    pub require_statistical_significance: bool,

    // Source tracking
    pub source: String,
    pub benchmark_auprc: Option<f64>,
    pub benchmark_f1: Option<f64>,
    pub benchmark_dataset: Option<String>,

    // Metadata
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub updated_by: Option<String>,
}

/// Data for creating/updating a profile calibration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProfileCalibrationInput {
    pub profile_name: String,
    pub ncd_threshold: f64,
    pub p_value_threshold: f64,
    pub composite_threshold: f64,
    pub ncd_weight: Option<f64>,
    pub p_value_weight: Option<f64>,
    pub compression_weight: Option<f64>,
    pub baseline_size: i32,
    pub window_size: i32,
    pub permutation_count: Option<i32>,
    pub adaptive_target_fpr: Option<f64>,
    pub require_statistical_significance: Option<bool>,
    pub source: Option<String>,
    pub benchmark_auprc: Option<f64>,
    pub benchmark_f1: Option<f64>,
    pub benchmark_dataset: Option<String>,
    pub updated_by: Option<String>,
}

impl Default for ProfileCalibrationInput {
    fn default() -> Self {
        Self {
            profile_name: "balanced".to_string(),
            ncd_threshold: 0.3,
            p_value_threshold: 0.05,
            composite_threshold: 0.6,
            ncd_weight: Some(0.5),
            p_value_weight: Some(0.25),
            compression_weight: Some(0.25),
            baseline_size: 400,
            window_size: 50,
            permutation_count: Some(100),
            adaptive_target_fpr: Some(0.01),
            require_statistical_significance: Some(true),
            source: Some("default".to_string()),
            benchmark_auprc: None,
            benchmark_f1: None,
            benchmark_dataset: None,
            updated_by: None,
        }
    }
}

/// Stream-specific calibration from auto-calibration during warmup
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct StreamCalibration {
    pub id: Uuid,
    pub stream_id: Uuid,

    // Calibration settings
    pub calibration_method: String,
    pub target_fpr: Option<f64>,

    // Calibrated thresholds
    pub calibrated_threshold: f64,
    pub calibrated_ncd_threshold: Option<f64>,
    pub calibrated_pvalue_threshold: Option<f64>,

    // Custom weights
    pub ncd_weight: Option<f64>,
    pub p_value_weight: Option<f64>,
    pub compression_weight: Option<f64>,

    // Calibration statistics
    pub warmup_sample_count: i32,
    pub warmup_score_mean: Option<f64>,
    pub warmup_score_stddev: Option<f64>,
    pub warmup_score_p95: Option<f64>,
    pub warmup_score_p99: Option<f64>,

    // Observed performance
    pub observed_fpr: Option<f64>,
    pub observed_f1: Option<f64>,
    pub observed_precision: Option<f64>,
    pub observed_recall: Option<f64>,

    // Metadata
    pub calibrated_at: DateTime<Utc>,
    pub last_validated_at: Option<DateTime<Utc>>,
    pub validation_sample_count: Option<i32>,
}

/// Data for creating/updating stream calibration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StreamCalibrationInput {
    pub stream_id: Uuid,
    pub calibration_method: String,
    pub target_fpr: Option<f64>,
    pub calibrated_threshold: f64,
    pub calibrated_ncd_threshold: Option<f64>,
    pub calibrated_pvalue_threshold: Option<f64>,
    pub ncd_weight: Option<f64>,
    pub p_value_weight: Option<f64>,
    pub compression_weight: Option<f64>,
    pub warmup_sample_count: i32,
    pub warmup_score_mean: Option<f64>,
    pub warmup_score_stddev: Option<f64>,
    pub warmup_score_p95: Option<f64>,
    pub warmup_score_p99: Option<f64>,
}

/// Aggregated feedback statistics for learning
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct FeedbackStatistics {
    pub id: Uuid,
    pub stream_id: Option<Uuid>,
    pub profile_name: Option<String>,
    pub tenant_id: Option<Uuid>,

    // Time period
    pub period_start: DateTime<Utc>,
    pub period_end: DateTime<Utc>,

    // Counts
    pub total_detections: i32,
    pub confirmed_count: i32,
    pub false_positive_count: i32,
    pub dismissed_count: i32,

    // Performance metrics
    pub observed_precision: Option<f64>,
    pub observed_recall: Option<f64>,

    // Score distribution
    pub avg_composite_score: Option<f64>,
    pub score_stddev: Option<f64>,
    pub score_min: Option<f64>,
    pub score_max: Option<f64>,

    // Recommendations
    pub recommended_threshold: Option<f64>,
    pub recommendation_confidence: Option<f64>,

    pub created_at: DateTime<Utc>,
}

/// Data for recording feedback statistics
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FeedbackStatisticsInput {
    pub stream_id: Option<Uuid>,
    pub profile_name: Option<String>,
    pub tenant_id: Option<Uuid>,
    pub period_start: DateTime<Utc>,
    pub period_end: DateTime<Utc>,
    pub total_detections: i32,
    pub confirmed_count: i32,
    pub false_positive_count: i32,
    pub dismissed_count: i32,
    pub observed_precision: Option<f64>,
    pub observed_recall: Option<f64>,
    pub avg_composite_score: Option<f64>,
    pub score_stddev: Option<f64>,
    pub score_min: Option<f64>,
    pub score_max: Option<f64>,
    pub recommended_threshold: Option<f64>,
    pub recommendation_confidence: Option<f64>,
}

/// Composite weights for score calculation
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CompositeWeights {
    pub ncd_weight: f64,
    pub p_value_weight: f64,
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
    /// Compute composite score using these weights
    pub fn compute_score(&self, ncd: f64, p_value: f64, compression_signal: f64) -> f64 {
        self.ncd_weight * ncd
            + self.p_value_weight * (1.0 - p_value)
            + self.compression_weight * compression_signal
    }
}

impl ProfileCalibration {
    /// Get composite weights from this profile
    pub fn weights(&self) -> CompositeWeights {
        CompositeWeights {
            ncd_weight: self.ncd_weight,
            p_value_weight: self.p_value_weight,
            compression_weight: self.compression_weight,
        }
    }
}

impl StreamCalibration {
    /// Get composite weights from this calibration (if custom, otherwise None)
    pub fn weights(&self) -> Option<CompositeWeights> {
        match (self.ncd_weight, self.p_value_weight, self.compression_weight) {
            (Some(ncd), Some(pv), Some(comp)) => Some(CompositeWeights {
                ncd_weight: ncd,
                p_value_weight: pv,
                compression_weight: comp,
            }),
            _ => None,
        }
    }
}
