//! Threshold tuning history repository

use crate::error::DbError;
use crate::models::{ThresholdTuneCreateParams, ThresholdTuneHistory};
use sqlx::PgPool;
use uuid::Uuid;

/// Record a new tuning event
pub async fn insert_tune_history(
    pool: &PgPool,
    params: &ThresholdTuneCreateParams,
) -> Result<ThresholdTuneHistory, DbError> {
    let record = sqlx::query_as::<_, ThresholdTuneHistory>(
        r#"
        INSERT INTO threshold_tune_history (id, stream_id, tune_type, old_value, new_value, reason, confidence)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING *
        "#,
    )
    .bind(Uuid::new_v4())
    .bind(params.stream_id)
    .bind(&params.tune_type)
    .bind(params.old_value)
    .bind(params.new_value)
    .bind(&params.reason)
    .bind(params.confidence)
    .fetch_one(pool)
    .await?;

    Ok(record)
}

/// List tuning history for a stream (most recent first)
pub async fn list_tune_history(
    pool: &PgPool,
    stream_id: Uuid,
    limit: i64,
) -> Result<Vec<ThresholdTuneHistory>, DbError> {
    let history = sqlx::query_as::<_, ThresholdTuneHistory>(
        r#"
        SELECT * FROM threshold_tune_history
        WHERE stream_id = $1
        ORDER BY created_at DESC
        LIMIT $2
        "#,
    )
    .bind(stream_id)
    .bind(limit)
    .fetch_all(pool)
    .await?;

    Ok(history)
}

/// Get tuning statistics for a stream
pub async fn get_tune_stats(pool: &PgPool, stream_id: Uuid) -> Result<TuneStats, DbError> {
    let count: (i64,) =
        sqlx::query_as("SELECT COUNT(*) FROM threshold_tune_history WHERE stream_id = $1")
            .bind(stream_id)
            .fetch_one(pool)
            .await?;

    let latest = sqlx::query_as::<_, ThresholdTuneHistory>(
        r#"
        SELECT * FROM threshold_tune_history
        WHERE stream_id = $1
        ORDER BY created_at DESC
        LIMIT 1
        "#,
    )
    .bind(stream_id)
    .fetch_optional(pool)
    .await?;

    Ok(TuneStats {
        total_adjustments: count.0,
        latest_adjustment: latest,
    })
}

/// Statistics about tuning for a stream
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct TuneStats {
    pub total_adjustments: i64,
    pub latest_adjustment: Option<ThresholdTuneHistory>,
}
