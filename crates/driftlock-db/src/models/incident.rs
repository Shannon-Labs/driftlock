//! DORA-compliant incident model

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

/// DORA Article 10 compliant incident types
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum IncidentType {
    MajorIctIncident,
    SignificantCyberThreat,
    TransactionAnomaly,
    DataBreach,
    ServiceDisruption,
    UnauthorizedAccess,
    DataIntegrity,
    ComplianceViolation,
}

impl IncidentType {
    pub fn as_str(&self) -> &'static str {
        match self {
            Self::MajorIctIncident => "major_ict_incident",
            Self::SignificantCyberThreat => "significant_cyber_threat",
            Self::TransactionAnomaly => "transaction_anomaly",
            Self::DataBreach => "data_breach",
            Self::ServiceDisruption => "service_disruption",
            Self::UnauthorizedAccess => "unauthorized_access",
            Self::DataIntegrity => "data_integrity",
            Self::ComplianceViolation => "compliance_violation",
        }
    }

    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "major_ict_incident" => Some(Self::MajorIctIncident),
            "significant_cyber_threat" => Some(Self::SignificantCyberThreat),
            "transaction_anomaly" => Some(Self::TransactionAnomaly),
            "data_breach" => Some(Self::DataBreach),
            "service_disruption" => Some(Self::ServiceDisruption),
            "unauthorized_access" => Some(Self::UnauthorizedAccess),
            "data_integrity" => Some(Self::DataIntegrity),
            "compliance_violation" => Some(Self::ComplianceViolation),
            _ => None,
        }
    }
}

/// Incident severity levels
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum IncidentSeverity {
    Critical,
    High,
    Medium,
    Low,
}

impl IncidentSeverity {
    pub fn as_str(&self) -> &'static str {
        match self {
            Self::Critical => "critical",
            Self::High => "high",
            Self::Medium => "medium",
            Self::Low => "low",
        }
    }

    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "critical" => Some(Self::Critical),
            "high" => Some(Self::High),
            "medium" => Some(Self::Medium),
            "low" => Some(Self::Low),
            _ => None,
        }
    }
}

/// Incident status workflow
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum IncidentStatus {
    Detected,
    Investigating,
    Classified,
    Reported,
    Mitigated,
    Resolved,
    Closed,
}

impl IncidentStatus {
    pub fn as_str(&self) -> &'static str {
        match self {
            Self::Detected => "detected",
            Self::Investigating => "investigating",
            Self::Classified => "classified",
            Self::Reported => "reported",
            Self::Mitigated => "mitigated",
            Self::Resolved => "resolved",
            Self::Closed => "closed",
        }
    }

    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "detected" => Some(Self::Detected),
            "investigating" => Some(Self::Investigating),
            "classified" => Some(Self::Classified),
            "reported" => Some(Self::Reported),
            "mitigated" => Some(Self::Mitigated),
            "resolved" => Some(Self::Resolved),
            "closed" => Some(Self::Closed),
            _ => None,
        }
    }
}

/// DORA-compliant incident record
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Incident {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub stream_id: Uuid,

    // Classification
    pub incident_type: String,
    pub severity: String,
    pub status: String,

    // Transaction context
    pub transaction_id: Option<String>,
    pub transaction_type: Option<String>,
    pub amount: Option<f64>,
    pub currency: Option<String>,
    pub sender_account: Option<String>,
    pub receiver_account: Option<String>,

    // Detection
    pub risk_score: f64,
    pub confidence: f64,
    pub detection_method: String,
    pub explanation: String,
    pub recommended_action: Option<String>,

    // DORA compliance
    pub regulatory_notification_required: bool,
    pub notification_deadline: Option<DateTime<Utc>>,
    pub notification_sent_at: Option<DateTime<Utc>>,
    pub notification_reference: Option<String>,

    // Impact
    pub impact_assessment: Option<serde_json::Value>,
    pub affected_clients_count: Option<i32>,
    pub financial_impact_eur: Option<f64>,

    // Audit
    pub raw_event: serde_json::Value,

    // Timestamps
    pub detected_at: DateTime<Utc>,
    pub classification_timestamp: Option<DateTime<Utc>>,
    pub mitigation_timestamp: Option<DateTime<Utc>>,
    pub resolution_timestamp: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Parameters for creating a new incident
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IncidentCreateParams {
    pub tenant_id: Uuid,
    pub stream_id: Uuid,
    pub incident_type: String,
    pub severity: String,

    // Transaction context (optional)
    pub transaction_id: Option<String>,
    pub transaction_type: Option<String>,
    pub amount: Option<f64>,
    pub currency: Option<String>,
    pub sender_account: Option<String>,
    pub receiver_account: Option<String>,

    // Detection
    pub risk_score: f64,
    pub confidence: f64,
    pub explanation: String,
    pub recommended_action: Option<String>,

    // DORA compliance
    pub regulatory_notification_required: bool,

    // Audit
    pub raw_event: serde_json::Value,
}

/// Incident-anomaly link
#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct IncidentAnomaly {
    pub incident_id: Uuid,
    pub anomaly_id: Uuid,
    pub correlation_type: String,
    pub correlation_score: Option<f64>,
    pub created_at: DateTime<Utc>,
}

/// Correlation type for incident-anomaly links
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum CorrelationType {
    Primary,
    Related,
    Correlated,
    Subsequent,
}

impl CorrelationType {
    pub fn as_str(&self) -> &'static str {
        match self {
            Self::Primary => "primary",
            Self::Related => "related",
            Self::Correlated => "correlated",
            Self::Subsequent => "subsequent",
        }
    }
}
