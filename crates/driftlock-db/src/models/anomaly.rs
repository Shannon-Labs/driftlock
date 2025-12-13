//! Anomaly model

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Anomaly {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub stream_id: Uuid,
    pub ncd: f64,
    pub compression_ratio: f64,
    pub entropy_change: f64,
    pub p_value: f64,
    pub confidence: f64,
    pub explanation: Option<String>,
    pub status: String,
    pub detected_at: DateTime<Utc>,
    pub details: Option<serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnomalyCreateParams {
    pub tenant_id: Uuid,
    pub stream_id: Uuid,
    pub ncd: f64,
    pub compression_ratio: f64,
    pub entropy_change: f64,
    pub p_value: f64,
    pub confidence: f64,
    pub explanation: Option<String>,
    pub details: Option<serde_json::Value>,
}
