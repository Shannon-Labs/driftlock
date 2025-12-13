//! API Key repository

use crate::error::DbError;
use crate::models::ApiKeyRecord;
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};
use rand::RngCore;
use sqlx::PgPool;
use uuid::Uuid;

/// Generate a new API key
pub fn generate_api_key() -> Result<(String, Uuid, String), DbError> {
    let id = Uuid::new_v4();
    let mut secret = [0u8; 20];
    rand::thread_rng().fill_bytes(&mut secret);
    let secret_encoded = base32::encode(base32::Alphabet::Rfc4648Lower { padding: false }, &secret);
    let key = format!("dlk_{}.{}", id, secret_encoded);

    let hash = hash_api_key(&key)?;
    Ok((key, id, hash))
}

/// Hash an API key using Argon2id
pub fn hash_api_key(key: &str) -> Result<String, DbError> {
    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::new(
        argon2::Algorithm::Argon2id,
        argon2::Version::V0x13,
        argon2::Params::new(65536, 1, 1, Some(32)).unwrap(),
    );
    let hash = argon2
        .hash_password(key.as_bytes(), &salt)
        .map_err(|e| DbError::Crypto(e.to_string()))?;
    Ok(hash.to_string())
}

/// Verify an API key against a hash
pub fn verify_api_key(hash: &str, candidate: &str) -> bool {
    let parsed = match PasswordHash::new(hash) {
        Ok(h) => h,
        Err(_) => return false,
    };
    Argon2::default()
        .verify_password(candidate.as_bytes(), &parsed)
        .is_ok()
}

/// Resolve API key by ID and verify
pub async fn resolve_api_key(
    pool: &PgPool,
    key_id: Uuid,
    key: &str,
) -> Result<Option<ApiKeyRecord>, DbError> {
    let record = sqlx::query_as::<_, ApiKeyRecord>(
        r#"
        SELECT ak.id, ak.tenant_id, ak.role, ak.key_hash, ak.stream_id,
               COALESCE(ak.rate_limit_rps, t.rate_limit_rps) as rate_limit_rps,
               t.name as tenant_name, t.slug as tenant_slug
        FROM api_keys ak
        JOIN tenants t ON ak.tenant_id = t.id
        WHERE ak.id = $1 AND ak.revoked_at IS NULL
        "#,
    )
    .bind(key_id)
    .fetch_optional(pool)
    .await?;

    match record {
        Some(r) if verify_api_key(&r.key_hash, key) => Ok(Some(r)),
        _ => Ok(None),
    }
}

/// Result of creating a new API key (includes plaintext key)
pub struct ApiKeyWithSecret {
    pub id: Uuid,
    pub key: String, // The plaintext key - only returned on creation
    pub name: String,
    pub role: String,
    pub stream_id: Option<Uuid>,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

/// Create a new API key for a tenant
pub async fn create_api_key(
    pool: &PgPool,
    tenant_id: Uuid,
    name: &str,
    role: &str,
    stream_id: Option<Uuid>,
) -> Result<ApiKeyWithSecret, DbError> {
    let (key, id, hash) = generate_api_key()?;

    let created_at: chrono::DateTime<chrono::Utc> = sqlx::query_scalar(
        r#"
        INSERT INTO api_keys (id, tenant_id, name, key_hash, role, stream_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at
        "#,
    )
    .bind(id)
    .bind(tenant_id)
    .bind(name)
    .bind(&hash)
    .bind(role)
    .bind(stream_id)
    .fetch_one(pool)
    .await?;

    Ok(ApiKeyWithSecret {
        id,
        key,
        name: name.to_string(),
        role: role.to_string(),
        stream_id,
        created_at,
    })
}

/// List API keys for a tenant (without secrets)
pub async fn list_api_keys(
    pool: &PgPool,
    tenant_id: Uuid,
) -> Result<Vec<crate::models::ApiKey>, DbError> {
    let keys = sqlx::query_as::<_, crate::models::ApiKey>(
        r#"
        SELECT id, tenant_id, name, key_hash, role, stream_id, rate_limit_rps, created_at, revoked_at
        FROM api_keys
        WHERE tenant_id = $1
        ORDER BY created_at DESC
        "#,
    )
    .bind(tenant_id)
    .fetch_all(pool)
    .await?;

    Ok(keys)
}

/// Revoke an API key (soft delete)
pub async fn revoke_api_key(pool: &PgPool, id: Uuid, tenant_id: Uuid) -> Result<bool, DbError> {
    let result = sqlx::query(
        "UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND tenant_id = $2 AND revoked_at IS NULL"
    )
    .bind(id)
    .bind(tenant_id)
    .execute(pool)
    .await?;

    Ok(result.rows_affected() > 0)
}

/// Regenerate an API key (new secret, same ID)
pub async fn regenerate_api_key(
    pool: &PgPool,
    id: Uuid,
    tenant_id: Uuid,
) -> Result<Option<ApiKeyWithSecret>, DbError> {
    // Generate new key with same structure but new secret
    let new_id = Uuid::new_v4();
    let mut secret = [0u8; 20];
    rand::thread_rng().fill_bytes(&mut secret);
    let secret_encoded = base32::encode(base32::Alphabet::Rfc4648Lower { padding: false }, &secret);
    let key = format!("dlk_{}.{}", new_id, secret_encoded);
    let hash = hash_api_key(&key)?;

    // Update the key in database
    let result =
        sqlx::query_as::<_, (String, String, Option<Uuid>, chrono::DateTime<chrono::Utc>)>(
            r#"
        UPDATE api_keys
        SET id = $3, key_hash = $4, updated_at = NOW()
        WHERE id = $1 AND tenant_id = $2 AND revoked_at IS NULL
        RETURNING name, role, stream_id, created_at
        "#,
        )
        .bind(id)
        .bind(tenant_id)
        .bind(new_id)
        .bind(&hash)
        .fetch_optional(pool)
        .await?;

    match result {
        Some((name, role, stream_id, created_at)) => Ok(Some(ApiKeyWithSecret {
            id: new_id,
            key,
            name,
            role,
            stream_id,
            created_at,
        })),
        None => Ok(None),
    }
}
