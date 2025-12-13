//! Calibration repository - DB operations for threshold calibration

use crate::error::DbError;
use crate::models::{
    FeedbackStatistics, FeedbackStatisticsInput, ProfileCalibration, ProfileCalibrationInput,
    StreamCalibration, StreamCalibrationInput,
};
use chrono::{DateTime, Utc};
use sqlx::PgPool;
use uuid::Uuid;

// ==================== PROFILE CALIBRATIONS ====================

/// Get all profile calibrations
pub async fn list_profile_calibrations(pool: &PgPool) -> Result<Vec<ProfileCalibration>, DbError> {
    let profiles = sqlx::query_as::<_, ProfileCalibration>(
        "SELECT * FROM profile_calibrations ORDER BY profile_name",
    )
    .fetch_all(pool)
    .await?;

    Ok(profiles)
}

/// Get a profile calibration by name
pub async fn get_profile_calibration(
    pool: &PgPool,
    profile_name: &str,
) -> Result<Option<ProfileCalibration>, DbError> {
    let profile = sqlx::query_as::<_, ProfileCalibration>(
        "SELECT * FROM profile_calibrations WHERE profile_name = $1",
    )
    .bind(profile_name)
    .fetch_optional(pool)
    .await?;

    Ok(profile)
}

/// Upsert a profile calibration
pub async fn upsert_profile_calibration(
    pool: &PgPool,
    input: &ProfileCalibrationInput,
) -> Result<ProfileCalibration, DbError> {
    let profile = sqlx::query_as::<_, ProfileCalibration>(
        r#"
        INSERT INTO profile_calibrations (
            id, profile_name, ncd_threshold, p_value_threshold, composite_threshold,
            ncd_weight, p_value_weight, compression_weight,
            baseline_size, window_size, permutation_count,
            adaptive_target_fpr, require_statistical_significance,
            source, benchmark_auprc, benchmark_f1, benchmark_dataset,
            updated_by
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
        )
        ON CONFLICT (profile_name) DO UPDATE SET
            ncd_threshold = EXCLUDED.ncd_threshold,
            p_value_threshold = EXCLUDED.p_value_threshold,
            composite_threshold = EXCLUDED.composite_threshold,
            ncd_weight = EXCLUDED.ncd_weight,
            p_value_weight = EXCLUDED.p_value_weight,
            compression_weight = EXCLUDED.compression_weight,
            baseline_size = EXCLUDED.baseline_size,
            window_size = EXCLUDED.window_size,
            permutation_count = EXCLUDED.permutation_count,
            adaptive_target_fpr = EXCLUDED.adaptive_target_fpr,
            require_statistical_significance = EXCLUDED.require_statistical_significance,
            source = EXCLUDED.source,
            benchmark_auprc = COALESCE(EXCLUDED.benchmark_auprc, profile_calibrations.benchmark_auprc),
            benchmark_f1 = COALESCE(EXCLUDED.benchmark_f1, profile_calibrations.benchmark_f1),
            benchmark_dataset = COALESCE(EXCLUDED.benchmark_dataset, profile_calibrations.benchmark_dataset),
            updated_by = EXCLUDED.updated_by,
            updated_at = NOW()
        RETURNING *
        "#,
    )
    .bind(&input.profile_name)
    .bind(input.ncd_threshold)
    .bind(input.p_value_threshold)
    .bind(input.composite_threshold)
    .bind(input.ncd_weight.unwrap_or(0.5))
    .bind(input.p_value_weight.unwrap_or(0.25))
    .bind(input.compression_weight.unwrap_or(0.25))
    .bind(input.baseline_size)
    .bind(input.window_size)
    .bind(input.permutation_count.unwrap_or(100))
    .bind(input.adaptive_target_fpr)
    .bind(input.require_statistical_significance.unwrap_or(true))
    .bind(input.source.as_deref().unwrap_or("manual"))
    .bind(input.benchmark_auprc)
    .bind(input.benchmark_f1)
    .bind(&input.benchmark_dataset)
    .bind(&input.updated_by)
    .fetch_one(pool)
    .await?;

    Ok(profile)
}

/// Apply benchmark results to a profile
pub async fn apply_benchmark_to_profile(
    pool: &PgPool,
    profile_name: &str,
    composite_threshold: f64,
    auprc: f64,
    f1: Option<f64>,
    dataset: &str,
    updated_by: Option<&str>,
) -> Result<ProfileCalibration, DbError> {
    let profile = sqlx::query_as::<_, ProfileCalibration>(
        r#"
        UPDATE profile_calibrations SET
            composite_threshold = $2,
            benchmark_auprc = $3,
            benchmark_f1 = $4,
            benchmark_dataset = $5,
            source = 'benchmark',
            updated_by = $6,
            updated_at = NOW()
        WHERE profile_name = $1
        RETURNING *
        "#,
    )
    .bind(profile_name)
    .bind(composite_threshold)
    .bind(auprc)
    .bind(f1)
    .bind(dataset)
    .bind(updated_by)
    .fetch_one(pool)
    .await?;

    Ok(profile)
}

// ==================== STREAM CALIBRATIONS ====================

/// Get stream calibration by stream ID
pub async fn get_stream_calibration(
    pool: &PgPool,
    stream_id: Uuid,
) -> Result<Option<StreamCalibration>, DbError> {
    let cal = sqlx::query_as::<_, StreamCalibration>(
        "SELECT * FROM stream_calibrations WHERE stream_id = $1",
    )
    .bind(stream_id)
    .fetch_optional(pool)
    .await?;

    Ok(cal)
}

/// Upsert stream calibration (from auto-calibration)
pub async fn upsert_stream_calibration(
    pool: &PgPool,
    input: &StreamCalibrationInput,
) -> Result<StreamCalibration, DbError> {
    let cal = sqlx::query_as::<_, StreamCalibration>(
        r#"
        INSERT INTO stream_calibrations (
            id, stream_id, calibration_method, target_fpr,
            calibrated_threshold, calibrated_ncd_threshold, calibrated_pvalue_threshold,
            ncd_weight, p_value_weight, compression_weight,
            warmup_sample_count, warmup_score_mean, warmup_score_stddev,
            warmup_score_p95, warmup_score_p99
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )
        ON CONFLICT (stream_id) DO UPDATE SET
            calibration_method = EXCLUDED.calibration_method,
            target_fpr = EXCLUDED.target_fpr,
            calibrated_threshold = EXCLUDED.calibrated_threshold,
            calibrated_ncd_threshold = EXCLUDED.calibrated_ncd_threshold,
            calibrated_pvalue_threshold = EXCLUDED.calibrated_pvalue_threshold,
            ncd_weight = EXCLUDED.ncd_weight,
            p_value_weight = EXCLUDED.p_value_weight,
            compression_weight = EXCLUDED.compression_weight,
            warmup_sample_count = EXCLUDED.warmup_sample_count,
            warmup_score_mean = EXCLUDED.warmup_score_mean,
            warmup_score_stddev = EXCLUDED.warmup_score_stddev,
            warmup_score_p95 = EXCLUDED.warmup_score_p95,
            warmup_score_p99 = EXCLUDED.warmup_score_p99,
            calibrated_at = NOW()
        RETURNING *
        "#,
    )
    .bind(input.stream_id)
    .bind(&input.calibration_method)
    .bind(input.target_fpr)
    .bind(input.calibrated_threshold)
    .bind(input.calibrated_ncd_threshold)
    .bind(input.calibrated_pvalue_threshold)
    .bind(input.ncd_weight)
    .bind(input.p_value_weight)
    .bind(input.compression_weight)
    .bind(input.warmup_sample_count)
    .bind(input.warmup_score_mean)
    .bind(input.warmup_score_stddev)
    .bind(input.warmup_score_p95)
    .bind(input.warmup_score_p99)
    .fetch_one(pool)
    .await?;

    Ok(cal)
}

/// Update observed performance metrics from feedback
pub async fn update_stream_calibration_metrics(
    pool: &PgPool,
    stream_id: Uuid,
    observed_fpr: Option<f64>,
    observed_f1: Option<f64>,
    observed_precision: Option<f64>,
    observed_recall: Option<f64>,
    validation_sample_count: i32,
) -> Result<Option<StreamCalibration>, DbError> {
    let cal = sqlx::query_as::<_, StreamCalibration>(
        r#"
        UPDATE stream_calibrations SET
            observed_fpr = COALESCE($2, observed_fpr),
            observed_f1 = COALESCE($3, observed_f1),
            observed_precision = COALESCE($4, observed_precision),
            observed_recall = COALESCE($5, observed_recall),
            validation_sample_count = $6,
            last_validated_at = NOW()
        WHERE stream_id = $1
        RETURNING *
        "#,
    )
    .bind(stream_id)
    .bind(observed_fpr)
    .bind(observed_f1)
    .bind(observed_precision)
    .bind(observed_recall)
    .bind(validation_sample_count)
    .fetch_optional(pool)
    .await?;

    Ok(cal)
}

/// Delete stream calibration
pub async fn delete_stream_calibration(pool: &PgPool, stream_id: Uuid) -> Result<bool, DbError> {
    let result = sqlx::query("DELETE FROM stream_calibrations WHERE stream_id = $1")
        .bind(stream_id)
        .execute(pool)
        .await?;

    Ok(result.rows_affected() > 0)
}

// ==================== FEEDBACK STATISTICS ====================

/// Record feedback statistics for a period
pub async fn record_feedback_statistics(
    pool: &PgPool,
    input: &FeedbackStatisticsInput,
) -> Result<FeedbackStatistics, DbError> {
    let stats = sqlx::query_as::<_, FeedbackStatistics>(
        r#"
        INSERT INTO feedback_statistics (
            id, stream_id, profile_name, tenant_id,
            period_start, period_end,
            total_detections, confirmed_count, false_positive_count, dismissed_count,
            observed_precision, observed_recall,
            avg_composite_score, score_stddev, score_min, score_max,
            recommended_threshold, recommendation_confidence
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
        )
        RETURNING *
        "#,
    )
    .bind(input.stream_id)
    .bind(&input.profile_name)
    .bind(input.tenant_id)
    .bind(input.period_start)
    .bind(input.period_end)
    .bind(input.total_detections)
    .bind(input.confirmed_count)
    .bind(input.false_positive_count)
    .bind(input.dismissed_count)
    .bind(input.observed_precision)
    .bind(input.observed_recall)
    .bind(input.avg_composite_score)
    .bind(input.score_stddev)
    .bind(input.score_min)
    .bind(input.score_max)
    .bind(input.recommended_threshold)
    .bind(input.recommendation_confidence)
    .fetch_one(pool)
    .await?;

    Ok(stats)
}

/// Get feedback statistics for a stream
pub async fn get_stream_feedback_statistics(
    pool: &PgPool,
    stream_id: Uuid,
    since: Option<DateTime<Utc>>,
    limit: i64,
) -> Result<Vec<FeedbackStatistics>, DbError> {
    let stats = sqlx::query_as::<_, FeedbackStatistics>(
        r#"
        SELECT * FROM feedback_statistics
        WHERE stream_id = $1
        AND ($2::timestamptz IS NULL OR period_end >= $2)
        ORDER BY period_end DESC
        LIMIT $3
        "#,
    )
    .bind(stream_id)
    .bind(since)
    .bind(limit)
    .fetch_all(pool)
    .await?;

    Ok(stats)
}

/// Get feedback statistics for a profile
pub async fn get_profile_feedback_statistics(
    pool: &PgPool,
    profile_name: &str,
    since: Option<DateTime<Utc>>,
    limit: i64,
) -> Result<Vec<FeedbackStatistics>, DbError> {
    let stats = sqlx::query_as::<_, FeedbackStatistics>(
        r#"
        SELECT * FROM feedback_statistics
        WHERE profile_name = $1
        AND ($2::timestamptz IS NULL OR period_end >= $2)
        ORDER BY period_end DESC
        LIMIT $3
        "#,
    )
    .bind(profile_name)
    .bind(since)
    .bind(limit)
    .fetch_all(pool)
    .await?;

    Ok(stats)
}

/// Aggregate feedback from anomaly_feedback table for a stream in a time period
pub async fn aggregate_stream_feedback(
    pool: &PgPool,
    stream_id: Uuid,
    period_start: DateTime<Utc>,
    period_end: DateTime<Utc>,
) -> Result<FeedbackStatisticsInput, DbError> {
    #[derive(sqlx::FromRow)]
    struct FeedbackAgg {
        total: i64,
        confirmed: i64,
        false_positive: i64,
        dismissed: i64,
    }

    let agg = sqlx::query_as::<_, FeedbackAgg>(
        r#"
        SELECT
            COUNT(*) as total,
            COUNT(*) FILTER (WHERE feedback_type = 'confirmed') as confirmed,
            COUNT(*) FILTER (WHERE feedback_type = 'false_positive') as false_positive,
            COUNT(*) FILTER (WHERE feedback_type = 'dismissed') as dismissed
        FROM anomaly_feedback
        WHERE stream_id = $1
        AND created_at >= $2 AND created_at < $3
        "#,
    )
    .bind(stream_id)
    .bind(period_start)
    .bind(period_end)
    .fetch_one(pool)
    .await?;

    let confirmed = agg.confirmed as i32;
    let false_positive = agg.false_positive as i32;
    let precision = if confirmed + false_positive > 0 {
        Some(confirmed as f64 / (confirmed + false_positive) as f64)
    } else {
        None
    };

    Ok(FeedbackStatisticsInput {
        stream_id: Some(stream_id),
        profile_name: None,
        tenant_id: None,
        period_start,
        period_end,
        total_detections: agg.total as i32,
        confirmed_count: confirmed,
        false_positive_count: false_positive,
        dismissed_count: agg.dismissed as i32,
        observed_precision: precision,
        observed_recall: None,
        avg_composite_score: None,
        score_stddev: None,
        score_min: None,
        score_max: None,
        recommended_threshold: None,
        recommendation_confidence: None,
    })
}

// ==================== HELPER FUNCTIONS ====================

/// Get effective threshold for a stream (checks stream calibration, then profile)
pub async fn get_effective_threshold(
    pool: &PgPool,
    stream_id: Uuid,
    profile_name: &str,
) -> Result<f64, DbError> {
    // First check stream-specific calibration
    if let Some(stream_cal) = get_stream_calibration(pool, stream_id).await? {
        return Ok(stream_cal.calibrated_threshold);
    }

    // Fall back to profile default
    if let Some(profile) = get_profile_calibration(pool, profile_name).await? {
        return Ok(profile.composite_threshold);
    }

    // Default fallback
    Ok(0.6)
}
