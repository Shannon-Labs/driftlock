//! Incident repository - DORA-compliant incident management

use crate::error::DbError;
use crate::models::{Incident, IncidentAnomaly, IncidentCreateParams};
use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use uuid::Uuid;

/// Create a new incident
pub async fn create_incident(
    pool: &PgPool,
    params: &IncidentCreateParams,
) -> Result<Incident, DbError> {
    // Calculate notification deadline if required (24 hours per DORA Article 19)
    let notification_deadline = if params.regulatory_notification_required {
        Some(Utc::now() + Duration::hours(24))
    } else {
        None
    };

    let incident = sqlx::query_as::<_, Incident>(
        r#"
        INSERT INTO incidents (
            id, tenant_id, stream_id, incident_type, severity, status,
            transaction_id, transaction_type, amount, currency, sender_account, receiver_account,
            risk_score, confidence, detection_method, explanation, recommended_action,
            regulatory_notification_required, notification_deadline, raw_event, detected_at
        )
        VALUES (
            $1, $2, $3, $4, $5, 'detected',
            $6, $7, $8, $9, $10, $11,
            $12, $13, 'CBAD', $14, $15,
            $16, $17, $18, now()
        )
        RETURNING *
        "#,
    )
    .bind(Uuid::new_v4())
    .bind(params.tenant_id)
    .bind(params.stream_id)
    .bind(&params.incident_type)
    .bind(&params.severity)
    .bind(&params.transaction_id)
    .bind(&params.transaction_type)
    .bind(params.amount)
    .bind(&params.currency)
    .bind(&params.sender_account)
    .bind(&params.receiver_account)
    .bind(params.risk_score)
    .bind(params.confidence)
    .bind(&params.explanation)
    .bind(&params.recommended_action)
    .bind(params.regulatory_notification_required)
    .bind(notification_deadline)
    .bind(&params.raw_event)
    .fetch_one(pool)
    .await?;

    Ok(incident)
}

/// Get incident by ID
pub async fn get_incident(pool: &PgPool, incident_id: Uuid) -> Result<Option<Incident>, DbError> {
    let incident = sqlx::query_as::<_, Incident>("SELECT * FROM incidents WHERE id = $1")
        .bind(incident_id)
        .fetch_optional(pool)
        .await?;

    Ok(incident)
}

/// List incidents for a tenant with optional filters
pub async fn list_incidents(
    pool: &PgPool,
    tenant_id: Uuid,
    stream_id: Option<Uuid>,
    status: Option<&str>,
    severity: Option<&str>,
    from: Option<DateTime<Utc>>,
    to: Option<DateTime<Utc>>,
    limit: i64,
    offset: i64,
) -> Result<Vec<Incident>, DbError> {
    let incidents = sqlx::query_as::<_, Incident>(
        r#"
        SELECT * FROM incidents
        WHERE tenant_id = $1
          AND ($2::uuid IS NULL OR stream_id = $2)
          AND ($3::text IS NULL OR status = $3)
          AND ($4::text IS NULL OR severity = $4)
          AND ($5::timestamptz IS NULL OR detected_at >= $5)
          AND ($6::timestamptz IS NULL OR detected_at <= $6)
        ORDER BY detected_at DESC
        LIMIT $7 OFFSET $8
        "#,
    )
    .bind(tenant_id)
    .bind(stream_id)
    .bind(status)
    .bind(severity)
    .bind(from)
    .bind(to)
    .bind(limit)
    .bind(offset)
    .fetch_all(pool)
    .await?;

    Ok(incidents)
}

/// Count incidents matching filters
pub async fn count_incidents(
    pool: &PgPool,
    tenant_id: Uuid,
    stream_id: Option<Uuid>,
    status: Option<&str>,
    severity: Option<&str>,
) -> Result<i64, DbError> {
    let count: (i64,) = sqlx::query_as(
        r#"
        SELECT COUNT(*) FROM incidents
        WHERE tenant_id = $1
          AND ($2::uuid IS NULL OR stream_id = $2)
          AND ($3::text IS NULL OR status = $3)
          AND ($4::text IS NULL OR severity = $4)
        "#,
    )
    .bind(tenant_id)
    .bind(stream_id)
    .bind(status)
    .bind(severity)
    .fetch_one(pool)
    .await?;

    Ok(count.0)
}

/// Update incident status
pub async fn update_incident_status(
    pool: &PgPool,
    incident_id: Uuid,
    status: &str,
) -> Result<Option<Incident>, DbError> {
    // Set appropriate timestamp based on status
    let incident = sqlx::query_as::<_, Incident>(
        r#"
        UPDATE incidents SET
            status = $2,
            classification_timestamp = CASE WHEN $2 = 'classified' THEN now() ELSE classification_timestamp END,
            mitigation_timestamp = CASE WHEN $2 = 'mitigated' THEN now() ELSE mitigation_timestamp END,
            resolution_timestamp = CASE WHEN $2 IN ('resolved', 'closed') THEN now() ELSE resolution_timestamp END,
            updated_at = now()
        WHERE id = $1
        RETURNING *
        "#,
    )
    .bind(incident_id)
    .bind(status)
    .fetch_optional(pool)
    .await?;

    Ok(incident)
}

/// Link anomaly to incident
pub async fn link_anomaly_to_incident(
    pool: &PgPool,
    incident_id: Uuid,
    anomaly_id: Uuid,
    correlation_type: &str,
    correlation_score: Option<f64>,
) -> Result<IncidentAnomaly, DbError> {
    let link = sqlx::query_as::<_, IncidentAnomaly>(
        r#"
        INSERT INTO incident_anomalies (incident_id, anomaly_id, correlation_type, correlation_score)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (incident_id, anomaly_id) DO UPDATE SET
            correlation_type = EXCLUDED.correlation_type,
            correlation_score = EXCLUDED.correlation_score
        RETURNING *
        "#,
    )
    .bind(incident_id)
    .bind(anomaly_id)
    .bind(correlation_type)
    .bind(correlation_score)
    .fetch_one(pool)
    .await?;

    Ok(link)
}

/// Get anomalies linked to an incident
pub async fn get_incident_anomalies(
    pool: &PgPool,
    incident_id: Uuid,
) -> Result<Vec<IncidentAnomaly>, DbError> {
    let links = sqlx::query_as::<_, IncidentAnomaly>(
        "SELECT * FROM incident_anomalies WHERE incident_id = $1 ORDER BY created_at",
    )
    .bind(incident_id)
    .fetch_all(pool)
    .await?;

    Ok(links)
}

/// Get incidents requiring notification
pub async fn get_pending_notifications(pool: &PgPool) -> Result<Vec<Incident>, DbError> {
    let incidents = sqlx::query_as::<_, Incident>(
        r#"
        SELECT * FROM incidents
        WHERE regulatory_notification_required = TRUE
          AND notification_sent_at IS NULL
          AND notification_deadline IS NOT NULL
        ORDER BY notification_deadline ASC
        "#,
    )
    .fetch_all(pool)
    .await?;

    Ok(incidents)
}

/// Mark notification as sent
pub async fn mark_notification_sent(
    pool: &PgPool,
    incident_id: Uuid,
    reference: &str,
) -> Result<Option<Incident>, DbError> {
    let incident = sqlx::query_as::<_, Incident>(
        r#"
        UPDATE incidents SET
            notification_sent_at = now(),
            notification_reference = $2,
            updated_at = now()
        WHERE id = $1
        RETURNING *
        "#,
    )
    .bind(incident_id)
    .bind(reference)
    .fetch_optional(pool)
    .await?;

    Ok(incident)
}

/// DORA compliance score
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DoraComplianceScore {
    pub score: f64,
    pub resolved_incidents: i64,
    pub notifications_sent: i64,
    pub overdue_notifications: i64,
    pub open_incidents: i64,
    pub avg_resolution_hours: Option<f64>,
}

/// Calculate DORA compliance score for a tenant
pub async fn calculate_dora_score(
    pool: &PgPool,
    tenant_id: Uuid,
) -> Result<DoraComplianceScore, DbError> {
    // Raw stats query
    let row: (i64, i64, i64, i64, Option<f64>) = sqlx::query_as(
        r#"
        SELECT
            COUNT(*) FILTER (WHERE status IN ('resolved', 'closed')) as resolved_count,
            COUNT(*) FILTER (WHERE regulatory_notification_required AND notification_sent_at IS NOT NULL) as notified_count,
            COUNT(*) FILTER (WHERE regulatory_notification_required AND notification_sent_at IS NULL AND notification_deadline < now()) as overdue_notifications,
            COUNT(*) FILTER (WHERE status NOT IN ('resolved', 'closed')) as open_count,
            AVG(EXTRACT(EPOCH FROM (resolution_timestamp - detected_at)) / 3600)::float8 as avg_resolution_hours
        FROM incidents
        WHERE tenant_id = $1
          AND detected_at > now() - INTERVAL '30 days'
        "#,
    )
    .bind(tenant_id)
    .fetch_one(pool)
    .await?;

    // Calculate compliance score (0-100)
    let mut score = 100.0;

    // Penalty for overdue notifications
    if row.2 > 0 {
        score -= (row.2 * 10) as f64;
    }

    // Penalty for open incidents
    if row.3 > 5 {
        score -= ((row.3 - 5) * 2) as f64;
    }

    Ok(DoraComplianceScore {
        score: score.max(0.0),
        resolved_incidents: row.0,
        notifications_sent: row.1,
        overdue_notifications: row.2,
        open_incidents: row.3,
        avg_resolution_hours: row.4,
    })
}
