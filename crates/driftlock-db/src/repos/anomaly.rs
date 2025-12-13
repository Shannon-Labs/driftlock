//! Anomaly repository

use crate::error::DbError;
use crate::models::{Anomaly, AnomalyCreateParams};
use sqlx::PgPool;
use uuid::Uuid;

/// Insert a new anomaly
pub async fn insert_anomaly(
    pool: &PgPool,
    params: &AnomalyCreateParams,
) -> Result<Anomaly, DbError> {
    let anomaly = sqlx::query_as::<_, Anomaly>(
        r#"
        INSERT INTO anomalies (id, tenant_id, stream_id, ncd, compression_ratio, entropy_change, p_value, confidence, explanation, status, details)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'new', $10)
        RETURNING *
        "#,
    )
    .bind(Uuid::new_v4())
    .bind(params.tenant_id)
    .bind(params.stream_id)
    .bind(params.ncd)
    .bind(params.compression_ratio)
    .bind(params.entropy_change)
    .bind(params.p_value)
    .bind(params.confidence)
    .bind(&params.explanation)
    .bind(&params.details)
    .fetch_one(pool)
    .await?;

    Ok(anomaly)
}

/// Get anomaly by ID
pub async fn get_anomaly(pool: &PgPool, id: Uuid) -> Result<Option<Anomaly>, DbError> {
    let anomaly = sqlx::query_as::<_, Anomaly>("SELECT * FROM anomalies WHERE id = $1")
        .bind(id)
        .fetch_optional(pool)
        .await?;

    Ok(anomaly)
}

/// List anomalies for a tenant
pub async fn list_anomalies(
    pool: &PgPool,
    tenant_id: Uuid,
    stream_id: Option<Uuid>,
    limit: i64,
    offset: i64,
) -> Result<Vec<Anomaly>, DbError> {
    let anomalies = if let Some(sid) = stream_id {
        sqlx::query_as::<_, Anomaly>(
            r#"
            SELECT * FROM anomalies
            WHERE tenant_id = $1 AND stream_id = $2
            ORDER BY detected_at DESC
            LIMIT $3 OFFSET $4
            "#,
        )
        .bind(tenant_id)
        .bind(sid)
        .bind(limit)
        .bind(offset)
        .fetch_all(pool)
        .await?
    } else {
        sqlx::query_as::<_, Anomaly>(
            r#"
            SELECT * FROM anomalies
            WHERE tenant_id = $1
            ORDER BY detected_at DESC
            LIMIT $2 OFFSET $3
            "#,
        )
        .bind(tenant_id)
        .bind(limit)
        .bind(offset)
        .fetch_all(pool)
        .await?
    };

    Ok(anomalies)
}

/// Update anomaly status
pub async fn update_anomaly_status(
    pool: &PgPool,
    id: Uuid,
    status: &str,
) -> Result<Option<Anomaly>, DbError> {
    let anomaly = sqlx::query_as::<_, Anomaly>(
        r#"
        UPDATE anomalies SET status = $2
        WHERE id = $1
        RETURNING *
        "#,
    )
    .bind(id)
    .bind(status)
    .fetch_optional(pool)
    .await?;

    Ok(anomaly)
}

/// Count anomalies for a tenant
pub async fn count_anomalies(
    pool: &PgPool,
    tenant_id: Uuid,
    stream_id: Option<Uuid>,
) -> Result<i64, DbError> {
    let count: (i64,) = if let Some(sid) = stream_id {
        sqlx::query_as("SELECT COUNT(*) FROM anomalies WHERE tenant_id = $1 AND stream_id = $2")
            .bind(tenant_id)
            .bind(sid)
            .fetch_one(pool)
            .await?
    } else {
        sqlx::query_as("SELECT COUNT(*) FROM anomalies WHERE tenant_id = $1")
            .bind(tenant_id)
            .fetch_one(pool)
            .await?
    };

    Ok(count.0)
}
