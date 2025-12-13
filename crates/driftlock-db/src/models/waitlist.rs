//! Waitlist model for pre-launch email capture

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Waitlist {
    pub id: Uuid,
    pub email: String,
    pub source: String,
    pub ip_address: Option<String>,
    pub created_at: DateTime<Utc>,
}
