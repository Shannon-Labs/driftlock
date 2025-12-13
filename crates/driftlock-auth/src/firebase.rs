//! Firebase JWT verification

use jsonwebtoken::{decode, decode_header, Algorithm, DecodingKey, Validation};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{Duration, SystemTime};
use thiserror::Error;
use tokio::sync::RwLock;
use tracing::{debug, error, warn};

const GOOGLE_CERTS_URL: &str =
    "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com";

#[derive(Error, Debug)]
pub enum FirebaseError {
    #[error("Invalid token")]
    InvalidToken,
    #[error("Token expired")]
    TokenExpired,
    #[error("Key fetch failed: {0}")]
    KeyFetchError(String),
    #[error("Missing key ID in token header")]
    MissingKeyId,
    #[error("Key not found: {0}")]
    KeyNotFound(String),
    #[error("Invalid issuer")]
    InvalidIssuer,
    #[error("Invalid audience")]
    InvalidAudience,
    #[error("Invalid subject")]
    InvalidSubject,
    #[error("JWT error: {0}")]
    JwtError(#[from] jsonwebtoken::errors::Error),
    #[error("HTTP error: {0}")]
    HttpError(#[from] reqwest::Error),
}

/// Firebase user information extracted from verified token
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FirebaseUser {
    /// Firebase user ID (UID)
    pub uid: String,
    /// User's email address (if available)
    pub email: Option<String>,
    /// Whether the email has been verified
    pub email_verified: bool,
}

/// JWT claims from Firebase ID token
#[derive(Debug, Clone, Serialize, Deserialize)]
struct FirebaseClaims {
    /// Subject (user ID)
    sub: String,
    /// Issuer
    iss: String,
    /// Audience (project ID)
    aud: String,
    /// Expiration time
    exp: u64,
    /// Issued at time
    iat: u64,
    /// Email (optional)
    email: Option<String>,
    /// Email verified flag (optional)
    email_verified: Option<bool>,
}

/// Cached public keys from Google
#[derive(Debug, Clone)]
struct CachedKeys {
    /// Map of key ID to PEM-encoded public key
    keys: HashMap<String, String>,
    /// When the cache expires
    expires_at: SystemTime,
}

/// Firebase authentication verifier
pub struct FirebaseAuth {
    project_id: String,
    http_client: Client,
    cached_keys: Arc<RwLock<Option<CachedKeys>>>,
}

impl FirebaseAuth {
    /// Create a new Firebase authenticator
    pub fn new(project_id: &str) -> Self {
        Self {
            project_id: project_id.to_string(),
            http_client: Client::new(),
            cached_keys: Arc::new(RwLock::new(None)),
        }
    }

    /// Verify a Firebase ID token and return user information
    pub async fn verify_token(&self, token: &str) -> Result<FirebaseUser, FirebaseError> {
        // Decode the header to get the key ID
        let header = decode_header(token)?;
        let kid = header.kid.ok_or(FirebaseError::MissingKeyId)?;

        debug!("Verifying token with key ID: {}", kid);

        // Get the public key for this key ID
        let public_key = self.get_public_key(&kid).await?;

        // Set up validation
        let mut validation = Validation::new(Algorithm::RS256);
        validation.set_issuer(&[format!(
            "https://securetoken.google.com/{}",
            self.project_id
        )]);
        validation.set_audience(&[&self.project_id]);

        // Decode and validate the token
        let decoding_key = DecodingKey::from_rsa_pem(public_key.as_bytes())?;
        let token_data = decode::<FirebaseClaims>(token, &decoding_key, &validation)?;

        let claims = token_data.claims;

        // Additional validation
        validate_firebase_claims(&claims, &self.project_id)?;

        // Extract user information
        Ok(FirebaseUser {
            uid: claims.sub,
            email: claims.email,
            email_verified: claims.email_verified.unwrap_or(false),
        })
    }

    /// Get a public key by key ID, fetching from Google if needed
    async fn get_public_key(&self, kid: &str) -> Result<String, FirebaseError> {
        // Check cache first
        {
            let cache = self.cached_keys.read().await;
            if let Some(cached) = cache.as_ref() {
                if SystemTime::now() < cached.expires_at {
                    if let Some(key) = cached.keys.get(kid) {
                        debug!("Using cached public key for kid: {}", kid);
                        return Ok(key.clone());
                    }
                } else {
                    debug!("Cached keys expired, will refresh");
                }
            }
        }

        // Cache miss or expired, fetch new keys
        self.fetch_and_cache_keys().await?;

        // Try to get the key again
        let cache = self.cached_keys.read().await;
        if let Some(cached) = cache.as_ref() {
            cached
                .keys
                .get(kid)
                .cloned()
                .ok_or_else(|| FirebaseError::KeyNotFound(kid.to_string()))
        } else {
            Err(FirebaseError::KeyFetchError(
                "Failed to cache keys".to_string(),
            ))
        }
    }

    /// Fetch public keys from Google and cache them
    async fn fetch_and_cache_keys(&self) -> Result<(), FirebaseError> {
        debug!("Fetching public keys from Google");

        let response = self.http_client.get(GOOGLE_CERTS_URL).send().await?;

        // Parse Cache-Control header to determine TTL
        let cache_duration = response
            .headers()
            .get("cache-control")
            .and_then(|v| v.to_str().ok())
            .and_then(parse_max_age)
            .unwrap_or(3600); // Default to 1 hour if not specified

        debug!("Cache duration: {} seconds", cache_duration);

        let keys: HashMap<String, String> = response.json().await?;

        if keys.is_empty() {
            error!("Received empty key set from Google");
            return Err(FirebaseError::KeyFetchError("Empty key set".to_string()));
        }

        debug!("Fetched {} public keys", keys.len());

        let cached = CachedKeys {
            keys,
            expires_at: SystemTime::now() + Duration::from_secs(cache_duration),
        };

        let mut cache = self.cached_keys.write().await;
        *cache = Some(cached);

        Ok(())
    }

    /// Force refresh the cached keys (useful for testing or manual refresh)
    pub async fn refresh_keys(&self) -> Result<(), FirebaseError> {
        self.fetch_and_cache_keys().await
    }
}

/// Parse max-age from Cache-Control header
fn parse_max_age(cache_control: &str) -> Option<u64> {
    for directive in cache_control.split(',') {
        let directive = directive.trim();
        if let Some(max_age_str) = directive.strip_prefix("max-age=") {
            if let Ok(seconds) = max_age_str.trim().parse::<u64>() {
                return Some(seconds);
            }
        }
    }
    None
}

/// Validate Firebase JWT claims
fn validate_firebase_claims(
    claims: &FirebaseClaims,
    project_id: &str,
) -> Result<(), FirebaseError> {
    // Verify issuer format
    let expected_issuer = format!("https://securetoken.google.com/{}", project_id);
    if claims.iss != expected_issuer {
        error!(
            "Invalid issuer: expected {}, got {}",
            expected_issuer, claims.iss
        );
        return Err(FirebaseError::InvalidIssuer);
    }

    // Verify audience
    if claims.aud != project_id {
        error!(
            "Invalid audience: expected {}, got {}",
            project_id, claims.aud
        );
        return Err(FirebaseError::InvalidAudience);
    }

    // Verify subject is non-empty
    if claims.sub.is_empty() {
        error!("Invalid subject: empty user ID");
        return Err(FirebaseError::InvalidSubject);
    }

    // Verify issued time is in the past
    let now = SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .unwrap()
        .as_secs();

    if claims.iat > now {
        warn!(
            "Token issued in the future: iat={}, now={}",
            claims.iat, now
        );
        return Err(FirebaseError::InvalidToken);
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_max_age() {
        assert_eq!(parse_max_age("max-age=3600"), Some(3600));
        assert_eq!(parse_max_age("public, max-age=7200"), Some(7200));
        assert_eq!(parse_max_age("max-age=3600, must-revalidate"), Some(3600));
        assert_eq!(parse_max_age("public"), None);
        assert_eq!(parse_max_age(""), None);
    }

    #[test]
    fn test_validate_claims() {
        let project_id = "test-project-id";

        let valid_claims = FirebaseClaims {
            sub: "user123".to_string(),
            iss: format!("https://securetoken.google.com/{}", project_id),
            aud: project_id.to_string(),
            exp: (SystemTime::now() + Duration::from_secs(3600))
                .duration_since(SystemTime::UNIX_EPOCH)
                .unwrap()
                .as_secs(),
            iat: SystemTime::now()
                .duration_since(SystemTime::UNIX_EPOCH)
                .unwrap()
                .as_secs(),
            email: Some("test@example.com".to_string()),
            email_verified: Some(true),
        };

        assert!(validate_firebase_claims(&valid_claims, project_id).is_ok());

        // Invalid issuer
        let invalid_issuer = FirebaseClaims {
            iss: "https://wrong.google.com/test-project-id".to_string(),
            ..valid_claims.clone()
        };
        assert!(matches!(
            validate_firebase_claims(&invalid_issuer, project_id),
            Err(FirebaseError::InvalidIssuer)
        ));

        // Invalid audience
        let invalid_aud = FirebaseClaims {
            aud: "wrong-project-id".to_string(),
            ..valid_claims.clone()
        };
        assert!(matches!(
            validate_firebase_claims(&invalid_aud, project_id),
            Err(FirebaseError::InvalidAudience)
        ));

        // Empty subject
        let empty_sub = FirebaseClaims {
            sub: "".to_string(),
            ..valid_claims.clone()
        };
        assert!(matches!(
            validate_firebase_claims(&empty_sub, project_id),
            Err(FirebaseError::InvalidSubject)
        ));

        // Token issued in the future
        let future_token = FirebaseClaims {
            iat: (SystemTime::now() + Duration::from_secs(3600))
                .duration_since(SystemTime::UNIX_EPOCH)
                .unwrap()
                .as_secs(),
            ..valid_claims.clone()
        };
        assert!(matches!(
            validate_firebase_claims(&future_token, project_id),
            Err(FirebaseError::InvalidToken)
        ));
    }
}
