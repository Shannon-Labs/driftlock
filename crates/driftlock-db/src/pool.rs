//! Database connection pool

use crate::error::DbError;
use sqlx::postgres::{PgPool, PgPoolOptions};
use std::time::Duration;

pub type DbPool = PgPool;

/// Create a new PostgreSQL connection pool
pub async fn create_pool(database_url: &str) -> Result<PgPool, DbError> {
    PgPoolOptions::new()
        .max_connections(25)
        .min_connections(5)
        .acquire_timeout(Duration::from_secs(30))
        .idle_timeout(Duration::from_secs(600))
        .max_lifetime(Duration::from_secs(1800))
        .connect(database_url)
        .await
        .map_err(DbError::from)
}

/// Health check - verify database connection
pub async fn ping(pool: &PgPool) -> Result<(), DbError> {
    sqlx::query("SELECT 1").execute(pool).await?;
    Ok(())
}
