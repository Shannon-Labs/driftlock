//! Threshold tuning history model for auto-tuning audit trail

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

/// History of threshold adjustments for a stream
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct ThresholdTuneHistory {
    pub id: Uuid,
    pub stream_id: Uuid,
    pub tune_type: String,
    pub old_value: Option<f64>,
    pub new_value: f64,
    pub reason: String,
    pub confidence: Option<f64>,
    pub created_at: DateTime<Utc>,
}

/// Parameters for recording a tuning event
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ThresholdTuneCreateParams {
    pub stream_id: Uuid,
    pub tune_type: String,
    pub old_value: Option<f64>,
    pub new_value: f64,
    pub reason: String,
    pub confidence: Option<f64>,
}
