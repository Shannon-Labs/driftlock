//! Waitlist repository

use crate::error::DbError;
use crate::models::Waitlist;
use sqlx::PgPool;

/// Add an email to the waitlist (idempotent - ignores duplicates)
pub async fn add_to_waitlist(
    pool: &PgPool,
    email: &str,
    source: &str,
    ip_address: Option<&str>,
) -> Result<(), DbError> {
    sqlx::query(
        r#"
        INSERT INTO waitlist (email, source, ip_address)
        VALUES ($1, $2, $3)
        ON CONFLICT (email) DO NOTHING
        "#,
    )
    .bind(email.to_lowercase())
    .bind(source)
    .bind(ip_address)
    .execute(pool)
    .await?;

    Ok(())
}

/// Get a waitlist entry by email
pub async fn get_waitlist_entry(pool: &PgPool, email: &str) -> Result<Option<Waitlist>, DbError> {
    let entry = sqlx::query_as::<_, Waitlist>("SELECT * FROM waitlist WHERE email = $1")
        .bind(email.to_lowercase())
        .fetch_optional(pool)
        .await?;

    Ok(entry)
}

/// Count total waitlist entries
pub async fn count_waitlist(pool: &PgPool) -> Result<i64, DbError> {
    let count: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM waitlist")
        .fetch_one(pool)
        .await?;

    Ok(count.0)
}
