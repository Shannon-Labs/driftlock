//! Tenant repository

use crate::error::DbError;
use crate::models::Tenant;
use sqlx::PgPool;
use uuid::Uuid;

/// Get tenant by ID (alias)
pub async fn get_tenant(pool: &PgPool, id: Uuid) -> Result<Option<Tenant>, DbError> {
    get_tenant_by_id(pool, id).await
}

/// Get tenant by ID
pub async fn get_tenant_by_id(pool: &PgPool, id: Uuid) -> Result<Option<Tenant>, DbError> {
    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        SELECT id, name, slug, status, plan, default_compressor, rate_limit_rps,
               email, firebase_uid, stripe_customer_id, stripe_subscription_id,
               stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
               trial_ends_at, verification_token, verified_at, verification_token_expires_at,
               created_at, updated_at
        FROM tenants
        WHERE id = $1
        "#,
    )
    .bind(id)
    .fetch_optional(pool)
    .await?;

    Ok(tenant)
}

/// Get tenant by Firebase UID
pub async fn get_tenant_by_firebase_uid(
    pool: &PgPool,
    uid: &str,
) -> Result<Option<Tenant>, DbError> {
    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        SELECT id, name, slug, status, plan, default_compressor, rate_limit_rps,
               email, firebase_uid, stripe_customer_id, stripe_subscription_id,
               stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
               trial_ends_at, verification_token, verified_at, verification_token_expires_at,
               created_at, updated_at
        FROM tenants
        WHERE firebase_uid = $1
        "#,
    )
    .bind(uid)
    .fetch_optional(pool)
    .await?;

    Ok(tenant)
}

/// Get tenant by email
pub async fn get_tenant_by_email(pool: &PgPool, email: &str) -> Result<Option<Tenant>, DbError> {
    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        SELECT id, name, slug, status, plan, default_compressor, rate_limit_rps,
               email, firebase_uid, stripe_customer_id, stripe_subscription_id,
               stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
               trial_ends_at, verification_token, verified_at, verification_token_expires_at,
               created_at, updated_at
        FROM tenants
        WHERE email = $1
        "#,
    )
    .bind(email)
    .fetch_optional(pool)
    .await?;

    Ok(tenant)
}

/// Get tenant by Stripe customer ID
pub async fn get_tenant_by_stripe_customer(
    pool: &PgPool,
    stripe_customer_id: &str,
) -> Result<Option<Tenant>, DbError> {
    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        SELECT id, name, slug, status, plan, default_compressor, rate_limit_rps,
               email, firebase_uid, stripe_customer_id, stripe_subscription_id,
               stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
               trial_ends_at, verification_token, verified_at, verification_token_expires_at,
               created_at, updated_at
        FROM tenants
        WHERE stripe_customer_id = $1
        "#,
    )
    .bind(stripe_customer_id)
    .fetch_optional(pool)
    .await?;

    Ok(tenant)
}

/// Get tenant by verification token
pub async fn get_tenant_by_verification_token(
    pool: &PgPool,
    token: &str,
) -> Result<Option<Tenant>, DbError> {
    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        SELECT id, name, slug, status, plan, default_compressor, rate_limit_rps,
               email, firebase_uid, stripe_customer_id, stripe_subscription_id,
               stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
               trial_ends_at, verification_token, verified_at, verification_token_expires_at,
               created_at, updated_at
        FROM tenants
        WHERE verification_token = $1
          AND (verification_token_expires_at IS NULL OR verification_token_expires_at > NOW())
          AND verified_at IS NULL
        "#,
    )
    .bind(token)
    .fetch_optional(pool)
    .await?;

    Ok(tenant)
}

/// Set verification token for a tenant
pub async fn set_verification_token(
    pool: &PgPool,
    tenant_id: Uuid,
    token: &str,
) -> Result<(), DbError> {
    sqlx::query(
        r#"
        UPDATE tenants
        SET verification_token = $2,
            verification_token_expires_at = NOW() + INTERVAL '24 hours',
            updated_at = NOW()
        WHERE id = $1
        "#,
    )
    .bind(tenant_id)
    .bind(token)
    .execute(pool)
    .await?;

    Ok(())
}

/// Mark tenant as verified
pub async fn verify_tenant(pool: &PgPool, tenant_id: Uuid) -> Result<(), DbError> {
    sqlx::query(
        r#"
        UPDATE tenants
        SET verified_at = NOW(),
            verification_token = NULL,
            verification_token_expires_at = NULL,
            updated_at = NOW()
        WHERE id = $1
        "#,
    )
    .bind(tenant_id)
    .execute(pool)
    .await?;

    Ok(())
}

/// Create a new tenant with trial
pub async fn create_tenant(
    pool: &PgPool,
    firebase_uid: &str,
    email: &str,
    company_name: &str,
) -> Result<Tenant, DbError> {
    // Generate slug from company name
    let slug = generate_slug(company_name);

    let tenant = sqlx::query_as::<_, Tenant>(
        r#"
        INSERT INTO tenants (
            name,
            slug,
            email,
            firebase_uid,
            plan,
            status,
            default_compressor,
            rate_limit_rps,
            trial_ends_at,
            payment_failure_count
        )
        VALUES ($1, $2, $3, $4, 'free', 'trialing', 'zstd', 60, NOW() + INTERVAL '14 days', 0)
        RETURNING id, name, slug, status, plan, default_compressor, rate_limit_rps,
                  email, firebase_uid, stripe_customer_id, stripe_subscription_id,
                  stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
                  trial_ends_at, verification_token, verified_at, verification_token_expires_at,
                  created_at, updated_at
        "#,
    )
    .bind(company_name)
    .bind(&slug)
    .bind(email)
    .bind(firebase_uid)
    .fetch_one(pool)
    .await?;

    Ok(tenant)
}

/// Update tenant subscription info
pub async fn update_subscription(
    pool: &PgPool,
    id: Uuid,
    stripe_customer_id: &str,
    plan: &str,
    status: &str,
) -> Result<(), DbError> {
    sqlx::query(
        r#"
        UPDATE tenants
        SET stripe_customer_id = $2,
            plan = $3,
            stripe_status = $4,
            updated_at = NOW()
        WHERE id = $1
        "#,
    )
    .bind(id)
    .bind(stripe_customer_id)
    .bind(plan)
    .bind(status)
    .execute(pool)
    .await?;

    Ok(())
}

/// Update tenant general info
pub async fn update_tenant(
    pool: &PgPool,
    id: Uuid,
    name: Option<&str>,
    slug: Option<&str>,
) -> Result<Tenant, DbError> {
    // Build query based on which fields are provided
    match (name, slug) {
        (Some(n), Some(s)) => {
            // Update both name and slug
            let tenant = sqlx::query_as::<_, Tenant>(
                r#"
                UPDATE tenants
                SET name = $2,
                    slug = $3,
                    updated_at = NOW()
                WHERE id = $1
                RETURNING id, name, slug, status, plan, default_compressor, rate_limit_rps,
                          email, firebase_uid, stripe_customer_id, stripe_subscription_id,
                          stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
                          trial_ends_at, verification_token, verified_at, verification_token_expires_at,
                          created_at, updated_at
                "#,
            )
            .bind(id)
            .bind(n)
            .bind(s)
            .fetch_one(pool)
            .await?;

            Ok(tenant)
        }
        (Some(n), None) => {
            // Update only name
            let tenant = sqlx::query_as::<_, Tenant>(
                r#"
                UPDATE tenants
                SET name = $2,
                    updated_at = NOW()
                WHERE id = $1
                RETURNING id, name, slug, status, plan, default_compressor, rate_limit_rps,
                          email, firebase_uid, stripe_customer_id, stripe_subscription_id,
                          stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
                          trial_ends_at, verification_token, verified_at, verification_token_expires_at,
                          created_at, updated_at
                "#,
            )
            .bind(id)
            .bind(n)
            .fetch_one(pool)
            .await?;

            Ok(tenant)
        }
        (None, Some(s)) => {
            // Update only slug
            let tenant = sqlx::query_as::<_, Tenant>(
                r#"
                UPDATE tenants
                SET slug = $2,
                    updated_at = NOW()
                WHERE id = $1
                RETURNING id, name, slug, status, plan, default_compressor, rate_limit_rps,
                          email, firebase_uid, stripe_customer_id, stripe_subscription_id,
                          stripe_status, current_period_end, grace_period_ends_at, payment_failure_count,
                          trial_ends_at, verification_token, verified_at, verification_token_expires_at,
                          created_at, updated_at
                "#,
            )
            .bind(id)
            .bind(s)
            .fetch_one(pool)
            .await?;

            Ok(tenant)
        }
        (None, None) => {
            // No fields to update, just fetch the current tenant
            get_tenant_by_id(pool, id).await?.ok_or(DbError::NotFound)
        }
    }
}

/// Set grace period for failed payment
pub async fn set_grace_period(pool: &PgPool, id: Uuid, days: i32) -> Result<(), DbError> {
    sqlx::query(
        r#"
        UPDATE tenants
        SET grace_period_ends_at = NOW() + INTERVAL '1 day' * $2,
            payment_failure_count = payment_failure_count + 1,
            updated_at = NOW()
        WHERE id = $1
        "#,
    )
    .bind(id)
    .bind(days)
    .execute(pool)
    .await?;

    Ok(())
}

/// Get total events ingested across all streams for a tenant
pub async fn get_total_events_for_tenant(pool: &PgPool, tenant_id: Uuid) -> Result<i64, DbError> {
    let result = sqlx::query_scalar::<_, Option<i64>>(
        r#"
        SELECT COALESCE(SUM(events_ingested), 0)
        FROM streams
        WHERE tenant_id = $1
        "#,
    )
    .bind(tenant_id)
    .fetch_one(pool)
    .await?;

    Ok(result.unwrap_or(0))
}

/// Generate a URL-safe slug from a company name
fn generate_slug(company_name: &str) -> String {
    company_name
        .to_lowercase()
        .chars()
        .map(|c| match c {
            'a'..='z' | '0'..='9' => c,
            ' ' | '-' | '_' => '-',
            _ => '-',
        })
        .collect::<String>()
        .split('-')
        .filter(|s| !s.is_empty())
        .collect::<Vec<_>>()
        .join("-")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_generate_slug() {
        assert_eq!(generate_slug("Acme Corp"), "acme-corp");
        assert_eq!(generate_slug("Test Company, Inc."), "test-company-inc");
        assert_eq!(generate_slug("My-Cool_Company"), "my-cool-company");
        assert_eq!(generate_slug("123 Industries"), "123-industries");
        assert_eq!(
            generate_slug("Special!@#$%Characters"),
            "special-characters"
        );
        assert_eq!(generate_slug("   Multiple   Spaces   "), "multiple-spaces");
        assert_eq!(generate_slug("Übertech GmbH"), "bertech-gmbh");
    }
}
