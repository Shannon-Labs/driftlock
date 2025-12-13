//! Driftlog audit trail model
//!
//! The Driftlog is a complete audit record of all detection decisions,
//! required for DORA compliance and regulatory audits.

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

/// Detection decision types
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum DetectionDecision {
    Normal,
    Anomaly,
    Escalated,
    Suppressed,
    Skipped,
}

impl DetectionDecision {
    pub fn as_str(&self) -> &'static str {
        match self {
            Self::Normal => "normal",
            Self::Anomaly => "anomaly",
            Self::Escalated => "escalated",
            Self::Suppressed => "suppressed",
            Self::Skipped => "skipped",
        }
    }

    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "normal" => Some(Self::Normal),
            "anomaly" => Some(Self::Anomaly),
            "escalated" => Some(Self::Escalated),
            "suppressed" => Some(Self::Suppressed),
            "skipped" => Some(Self::Skipped),
            _ => None,
        }
    }
}

/// Driftlog entry - complete audit record of a detection decision
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct DriftlogEntry {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub stream_id: Uuid,

    // Event identification
    pub event_hash: String,
    pub transaction_id: Option<String>,

    // Decision
    pub decision: String,

    // Metrics at decision time
    pub ncd: Option<f64>,
    pub compression_ratio: Option<f64>,
    pub entropy: Option<f64>,
    pub p_value: Option<f64>,
    pub confidence: Option<f64>,

    // Thresholds
    pub ncd_threshold_applied: Option<f64>,
    pub profile_applied: Option<String>,

    // References
    pub anomaly_id: Option<Uuid>,
    pub incident_id: Option<Uuid>,

    // Audit metadata
    pub processing_time_us: Option<i32>,
    pub api_key_id: Option<Uuid>,
    pub client_ip: Option<String>,

    pub created_at: DateTime<Utc>,
}

/// Parameters for creating driftlog entry
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DriftlogCreateParams {
    pub tenant_id: Uuid,
    pub stream_id: Uuid,
    pub event_hash: String,
    pub transaction_id: Option<String>,
    pub decision: String,
    pub ncd: Option<f64>,
    pub compression_ratio: Option<f64>,
    pub entropy: Option<f64>,
    pub p_value: Option<f64>,
    pub confidence: Option<f64>,
    pub ncd_threshold_applied: Option<f64>,
    pub profile_applied: Option<String>,
    pub anomaly_id: Option<Uuid>,
    pub incident_id: Option<Uuid>,
    pub processing_time_us: Option<i32>,
    pub api_key_id: Option<Uuid>,
    pub client_ip: Option<String>,
}

/// Query filters for driftlog
#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct DriftlogQuery {
    pub stream_id: Option<Uuid>,
    pub decision: Option<String>,
    pub transaction_id: Option<String>,
    pub from: Option<DateTime<Utc>>,
    pub to: Option<DateTime<Utc>>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

/// Driftlog statistics
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DriftlogStats {
    pub total_events: i64,
    pub normal_count: i64,
    pub anomaly_count: i64,
    pub escalated_count: i64,
    pub suppressed_count: i64,
    pub skipped_count: i64,
    pub avg_processing_us: Option<i32>,
}
