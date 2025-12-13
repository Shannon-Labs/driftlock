//! API Key model

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct ApiKey {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub name: String,
    pub key_hash: String,
    pub role: String,
    pub stream_id: Option<Uuid>,
    pub rate_limit_rps: Option<i32>,
    pub created_at: DateTime<Utc>,
    pub revoked_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, FromRow)]
pub struct ApiKeyRecord {
    pub id: Uuid,
    pub tenant_id: Uuid,
    pub tenant_name: String,
    pub tenant_slug: String,
    pub role: String,
    pub key_hash: String,
    pub stream_id: Option<Uuid>,
    pub rate_limit_rps: i32,
}
