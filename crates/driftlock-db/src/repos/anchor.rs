//! StreamAnchor repository for drift detection

use crate::error::DbError;
use crate::models::{StreamAnchor, StreamAnchorCreateParams};
use sqlx::PgPool;
use uuid::Uuid;

/// Get the active anchor for a stream
pub async fn get_active_anchor(
    pool: &PgPool,
    stream_id: Uuid,
) -> Result<Option<StreamAnchor>, DbError> {
    let anchor = sqlx::query_as::<_, StreamAnchor>(
        "SELECT * FROM stream_anchors WHERE stream_id = $1 AND is_active = true",
    )
    .bind(stream_id)
    .fetch_optional(pool)
    .await?;

    Ok(anchor)
}

/// Create a new anchor (automatically deactivates existing active anchor)
pub async fn create_anchor(
    pool: &PgPool,
    params: &StreamAnchorCreateParams,
) -> Result<StreamAnchor, DbError> {
    let new_id = Uuid::new_v4();

    // Deactivate existing active anchor first (within transaction)
    let mut tx = pool.begin().await?;

    // Get existing active anchor ID for superseded_by reference
    let existing: Option<(Uuid,)> =
        sqlx::query_as("SELECT id FROM stream_anchors WHERE stream_id = $1 AND is_active = true")
            .bind(params.stream_id)
            .fetch_optional(&mut *tx)
            .await?;

    // Deactivate existing
    if existing.is_some() {
        sqlx::query(
            r#"
            UPDATE stream_anchors
            SET is_active = false, superseded_at = NOW(), superseded_by = $2
            WHERE stream_id = $1 AND is_active = true
            "#,
        )
        .bind(params.stream_id)
        .bind(new_id)
        .execute(&mut *tx)
        .await?;
    }

    // Insert new anchor
    let anchor = sqlx::query_as::<_, StreamAnchor>(
        r#"
        INSERT INTO stream_anchors (
            id, stream_id, anchor_data, compressor, event_count,
            calibration_completed_at, is_active,
            baseline_entropy, baseline_compression_ratio, baseline_ncd_self,
            drift_ncd_threshold
        )
        VALUES ($1, $2, $3, $4, $5, $6, true, $7, $8, $9, $10)
        RETURNING *
        "#,
    )
    .bind(new_id)
    .bind(params.stream_id)
    .bind(&params.anchor_data)
    .bind(&params.compressor)
    .bind(params.event_count)
    .bind(params.calibration_completed_at)
    .bind(params.baseline_entropy)
    .bind(params.baseline_compression_ratio)
    .bind(params.baseline_ncd_self)
    .bind(params.drift_ncd_threshold)
    .fetch_one(&mut *tx)
    .await?;

    tx.commit().await?;

    Ok(anchor)
}

/// Deactivate the current active anchor for a stream
pub async fn deactivate_anchor(pool: &PgPool, stream_id: Uuid) -> Result<(), DbError> {
    sqlx::query(
        r#"
        UPDATE stream_anchors
        SET is_active = false, superseded_at = NOW()
        WHERE stream_id = $1 AND is_active = true
        "#,
    )
    .bind(stream_id)
    .execute(pool)
    .await?;

    Ok(())
}

/// List anchor history for a stream (most recent first)
pub async fn list_anchor_history(
    pool: &PgPool,
    stream_id: Uuid,
    limit: i64,
) -> Result<Vec<StreamAnchor>, DbError> {
    let anchors = sqlx::query_as::<_, StreamAnchor>(
        r#"
        SELECT * FROM stream_anchors
        WHERE stream_id = $1
        ORDER BY created_at DESC
        LIMIT $2
        "#,
    )
    .bind(stream_id)
    .bind(limit)
    .fetch_all(pool)
    .await?;

    Ok(anchors)
}
