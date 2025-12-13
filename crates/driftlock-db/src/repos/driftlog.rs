//! Driftlog repository - audit trail for all detection decisions

use crate::error::DbError;
use crate::models::{DriftlogCreateParams, DriftlogEntry, DriftlogQuery, DriftlogStats};
use chrono::{DateTime, Utc};
use sqlx::PgPool;
use uuid::Uuid;

/// Insert a driftlog entry
pub async fn insert_driftlog(
    pool: &PgPool,
    params: &DriftlogCreateParams,
) -> Result<DriftlogEntry, DbError> {
    let entry = sqlx::query_as::<_, DriftlogEntry>(
        r#"
        INSERT INTO driftlog (
            id, tenant_id, stream_id, event_hash, transaction_id, decision,
            ncd, compression_ratio, entropy, p_value, confidence,
            ncd_threshold_applied, profile_applied,
            anomaly_id, incident_id, processing_time_us, api_key_id, client_ip
        )
        VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11,
            $12, $13,
            $14, $15, $16, $17, $18
        )
        RETURNING *
        "#,
    )
    .bind(Uuid::new_v4())
    .bind(params.tenant_id)
    .bind(params.stream_id)
    .bind(&params.event_hash)
    .bind(&params.transaction_id)
    .bind(&params.decision)
    .bind(params.ncd)
    .bind(params.compression_ratio)
    .bind(params.entropy)
    .bind(params.p_value)
    .bind(params.confidence)
    .bind(params.ncd_threshold_applied)
    .bind(&params.profile_applied)
    .bind(params.anomaly_id)
    .bind(params.incident_id)
    .bind(params.processing_time_us)
    .bind(params.api_key_id)
    .bind(&params.client_ip)
    .fetch_one(pool)
    .await?;

    Ok(entry)
}

/// Query driftlog with filters
pub async fn query_driftlog(
    pool: &PgPool,
    tenant_id: Uuid,
    query: &DriftlogQuery,
) -> Result<Vec<DriftlogEntry>, DbError> {
    let limit = query.limit.unwrap_or(100).min(1000);
    let offset = query.offset.unwrap_or(0);

    let entries = sqlx::query_as::<_, DriftlogEntry>(
        r#"
        SELECT * FROM driftlog
        WHERE tenant_id = $1
          AND ($2::uuid IS NULL OR stream_id = $2)
          AND ($3::text IS NULL OR decision = $3)
          AND ($4::text IS NULL OR transaction_id = $4)
          AND ($5::timestamptz IS NULL OR created_at >= $5)
          AND ($6::timestamptz IS NULL OR created_at <= $6)
        ORDER BY created_at DESC
        LIMIT $7 OFFSET $8
        "#,
    )
    .bind(tenant_id)
    .bind(query.stream_id)
    .bind(&query.decision)
    .bind(&query.transaction_id)
    .bind(query.from)
    .bind(query.to)
    .bind(limit)
    .bind(offset)
    .fetch_all(pool)
    .await?;

    Ok(entries)
}

/// Count driftlog entries
pub async fn count_driftlog(
    pool: &PgPool,
    tenant_id: Uuid,
    query: &DriftlogQuery,
) -> Result<i64, DbError> {
    let count: (i64,) = sqlx::query_as(
        r#"
        SELECT COUNT(*) FROM driftlog
        WHERE tenant_id = $1
          AND ($2::uuid IS NULL OR stream_id = $2)
          AND ($3::text IS NULL OR decision = $3)
          AND ($4::text IS NULL OR transaction_id = $4)
          AND ($5::timestamptz IS NULL OR created_at >= $5)
          AND ($6::timestamptz IS NULL OR created_at <= $6)
        "#,
    )
    .bind(tenant_id)
    .bind(query.stream_id)
    .bind(&query.decision)
    .bind(&query.transaction_id)
    .bind(query.from)
    .bind(query.to)
    .fetch_one(pool)
    .await?;

    Ok(count.0)
}

/// Get driftlog statistics for a stream
pub async fn get_driftlog_stats(
    pool: &PgPool,
    tenant_id: Uuid,
    stream_id: Option<Uuid>,
    since: DateTime<Utc>,
) -> Result<DriftlogStats, DbError> {
    let row: (i64, i64, i64, i64, i64, i64, Option<i32>) = sqlx::query_as(
        r#"
        SELECT
            COUNT(*) as total_events,
            COUNT(*) FILTER (WHERE decision = 'normal') as normal_count,
            COUNT(*) FILTER (WHERE decision = 'anomaly') as anomaly_count,
            COUNT(*) FILTER (WHERE decision = 'escalated') as escalated_count,
            COUNT(*) FILTER (WHERE decision = 'suppressed') as suppressed_count,
            COUNT(*) FILTER (WHERE decision = 'skipped') as skipped_count,
            AVG(processing_time_us)::int4 as avg_processing_us
        FROM driftlog
        WHERE tenant_id = $1
          AND ($2::uuid IS NULL OR stream_id = $2)
          AND created_at >= $3
        "#,
    )
    .bind(tenant_id)
    .bind(stream_id)
    .bind(since)
    .fetch_one(pool)
    .await?;

    Ok(DriftlogStats {
        total_events: row.0,
        normal_count: row.1,
        anomaly_count: row.2,
        escalated_count: row.3,
        suppressed_count: row.4,
        skipped_count: row.5,
        avg_processing_us: row.6,
    })
}

/// Get driftlog entry by ID
pub async fn get_driftlog_entry(pool: &PgPool, id: Uuid) -> Result<Option<DriftlogEntry>, DbError> {
    let entry = sqlx::query_as::<_, DriftlogEntry>("SELECT * FROM driftlog WHERE id = $1")
        .bind(id)
        .fetch_optional(pool)
        .await?;

    Ok(entry)
}
