//! Stream model

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Stream {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub slug: String,
    #[sqlx(rename = "type")]
    pub stream_type: String,
    pub seed: i64,
    pub compressor: String,
    pub retention_days: i32,
    pub events_ingested: i64,
    pub is_calibrated: bool,
    pub min_baseline_size: i32,
    pub detection_profile: String,
    pub auto_tune_enabled: bool,
    // Anchor/drift detection settings
    pub anchor_enabled: bool,
    pub drift_ncd_threshold: f64,
    pub anchor_reset_on_drift: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct StreamSettings {
    pub baseline_size: i32,
    pub window_size: i32,
    pub hop_size: i32,
    pub ncd_threshold: f64,
    pub p_value_threshold: f64,
    pub permutation_count: i32,
    pub compressor: String,
    pub detection_profile: String,
    pub auto_tune_enabled: bool,
}

impl StreamSettings {
    pub fn apply_defaults(&mut self) {
        if self.compressor.is_empty() {
            self.compressor = "zstd".to_string();
        }
        if self.detection_profile.is_empty() {
            self.detection_profile = "auto".to_string();
        }
    }
}
