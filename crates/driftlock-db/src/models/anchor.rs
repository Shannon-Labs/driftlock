//! StreamAnchor model for drift detection baseline

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

/// Frozen baseline snapshot for drift detection
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct StreamAnchor {
    pub id: Uuid,
    pub stream_id: Uuid,
    pub anchor_data: Vec<u8>,
    pub compressor: String,
    pub event_count: i32,
    pub calibration_completed_at: DateTime<Utc>,
    pub is_active: bool,
    pub baseline_entropy: Option<f64>,
    pub baseline_compression_ratio: Option<f64>,
    pub baseline_ncd_self: Option<f64>,
    pub drift_ncd_threshold: f64,
    pub created_at: DateTime<Utc>,
    pub superseded_at: Option<DateTime<Utc>>,
    pub superseded_by: Option<Uuid>,
}

/// Parameters for creating a new anchor
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StreamAnchorCreateParams {
    pub stream_id: Uuid,
    pub anchor_data: Vec<u8>,
    pub compressor: String,
    pub event_count: i32,
    pub calibration_completed_at: DateTime<Utc>,
    pub baseline_entropy: Option<f64>,
    pub baseline_compression_ratio: Option<f64>,
    pub baseline_ncd_self: Option<f64>,
    pub drift_ncd_threshold: f64,
}
