//! Tenant model

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Tenant {
    pub id: Uuid,
    pub name: String,
    pub slug: String,
    pub status: String,
    pub plan: String,
    pub default_compressor: String,
    pub rate_limit_rps: i32,
    pub email: Option<String>,
    pub firebase_uid: Option<String>,
    pub stripe_customer_id: Option<String>,
    pub stripe_subscription_id: Option<String>,
    pub stripe_status: Option<String>,
    pub current_period_end: Option<DateTime<Utc>>,
    pub grace_period_ends_at: Option<DateTime<Utc>>,
    pub payment_failure_count: i32,
    pub trial_ends_at: Option<DateTime<Utc>>,
    pub verification_token: Option<String>,
    pub verified_at: Option<DateTime<Utc>>,
    pub verification_token_expires_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Default)]
pub struct TenantCreateParams {
    pub name: String,
    pub slug: Option<String>,
    pub plan: String,
    pub email: Option<String>,
    pub firebase_uid: Option<String>,
    pub default_compressor: String,
    pub rate_limit_rps: i32,
}
