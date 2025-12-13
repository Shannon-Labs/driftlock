//! Stream repository

use crate::error::DbError;
use crate::models::Stream;
use sqlx::PgPool;
use uuid::Uuid;

/// Get stream by ID
pub async fn get_stream(pool: &PgPool, id: Uuid) -> Result<Option<Stream>, DbError> {
    let stream = sqlx::query_as::<_, Stream>("SELECT * FROM streams WHERE id = $1")
        .bind(id)
        .fetch_optional(pool)
        .await?;

    Ok(stream)
}

/// Get stream by slug for a tenant
pub async fn get_stream_by_slug(
    pool: &PgPool,
    tenant_id: Uuid,
    slug: &str,
) -> Result<Option<Stream>, DbError> {
    let stream =
        sqlx::query_as::<_, Stream>("SELECT * FROM streams WHERE tenant_id = $1 AND slug = $2")
            .bind(tenant_id)
            .bind(slug)
            .fetch_optional(pool)
            .await?;

    Ok(stream)
}

/// List streams for a tenant
pub async fn list_streams(
    pool: &PgPool,
    tenant_id: Uuid,
    limit: i64,
    offset: i64,
) -> Result<Vec<Stream>, DbError> {
    let streams = sqlx::query_as::<_, Stream>(
        r#"
        SELECT * FROM streams
        WHERE tenant_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(tenant_id)
    .bind(limit)
    .bind(offset)
    .fetch_all(pool)
    .await?;

    Ok(streams)
}

/// Create a new stream
pub async fn create_stream(
    pool: &PgPool,
    tenant_id: Uuid,
    slug: &str,
    stream_type: &str,
) -> Result<Stream, DbError> {
    let stream = sqlx::query_as::<_, Stream>(
        r#"
        INSERT INTO streams (id, tenant_id, slug, type, seed, compressor, retention_days, events_ingested, is_calibrated, min_baseline_size, detection_profile, auto_tune_enabled)
        VALUES ($1, $2, $3, $4, $5, 'zstd', 30, 0, false, 50, 'auto', true)
        RETURNING *
        "#,
    )
    .bind(Uuid::new_v4())
    .bind(tenant_id)
    .bind(slug)
    .bind(stream_type)
    .bind(rand::random::<i64>().abs())
    .fetch_one(pool)
    .await?;

    Ok(stream)
}

/// Update stream settings
pub async fn update_stream_profile(
    pool: &PgPool,
    id: Uuid,
    detection_profile: &str,
    auto_tune_enabled: bool,
) -> Result<Option<Stream>, DbError> {
    let stream = sqlx::query_as::<_, Stream>(
        r#"
        UPDATE streams SET
            detection_profile = $2,
            auto_tune_enabled = $3,
            updated_at = NOW()
        WHERE id = $1
        RETURNING *
        "#,
    )
    .bind(id)
    .bind(detection_profile)
    .bind(auto_tune_enabled)
    .fetch_optional(pool)
    .await?;

    Ok(stream)
}

/// Increment events ingested counter
pub async fn increment_events(pool: &PgPool, id: Uuid, count: i64) -> Result<(), DbError> {
    sqlx::query(
        "UPDATE streams SET events_ingested = events_ingested + $2, updated_at = NOW() WHERE id = $1"
    )
    .bind(id)
    .bind(count)
    .execute(pool)
    .await?;

    Ok(())
}

/// Mark stream as calibrated
pub async fn mark_calibrated(pool: &PgPool, id: Uuid) -> Result<(), DbError> {
    sqlx::query("UPDATE streams SET is_calibrated = true, updated_at = NOW() WHERE id = $1")
        .bind(id)
        .execute(pool)
        .await?;

    Ok(())
}

/// Increment events and optionally mark as calibrated if threshold reached
pub async fn increment_events_and_check_calibration(
    pool: &PgPool,
    id: Uuid,
    count: i64,
) -> Result<bool, DbError> {
    // Atomic update that also checks calibration threshold
    let result = sqlx::query_scalar::<_, bool>(
        r#"
        UPDATE streams
        SET
            events_ingested = events_ingested + $2,
            is_calibrated = CASE
                WHEN NOT is_calibrated AND (events_ingested + $2) >= min_baseline_size
                THEN true
                ELSE is_calibrated
            END,
            updated_at = NOW()
        WHERE id = $1
        RETURNING is_calibrated
        "#,
    )
    .bind(id)
    .bind(count)
    .fetch_one(pool)
    .await?;

    Ok(result)
}

/// Count streams for a tenant
pub async fn count_streams(pool: &PgPool, tenant_id: Uuid) -> Result<i64, DbError> {
    let count: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM streams WHERE tenant_id = $1")
        .bind(tenant_id)
        .fetch_one(pool)
        .await?;

    Ok(count.0)
}
