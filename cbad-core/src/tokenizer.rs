//! High-Entropy Tokenizer for preprocessing data before compression.
//!
//! This module replaces high-entropy fields (UUIDs, hashes, JWTs, Base64) with
//! fixed tokens, allowing the compression algorithm to focus on structural patterns
//! rather than random noise.
//!
//! Pattern order matters - most specific patterns must be applied first:
//! 1. JWT (contains Base64-like segments)
//! 2. UUID (specific hex pattern with dashes)
//! 3. Hash (32-64 hex chars)
//! 4. Base64 (most general pattern)

use lazy_static::lazy_static;
use regex::bytes::Regex;
use serde::{Deserialize, Serialize};
use std::sync::atomic::{AtomicU64, Ordering};

lazy_static! {
    /// JWT: Three Base64URL segments separated by dots
    /// Must match before Base64 to avoid partial matches
    static ref JWT_PATTERN: Regex = Regex::new(
        r"eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+"
    ).expect("JWT regex compilation failed");

    /// UUID: 8-4-4-4-12 hex pattern (case-insensitive)
    /// Must match before Hash to capture full UUIDs
    static ref UUID_PATTERN: Regex = Regex::new(
        r"(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
    ).expect("UUID regex compilation failed");

    /// Hex hashes: 32-64 hex characters (MD5, SHA-1, SHA-256, SHA-512)
    /// Word boundaries to avoid matching substrings
    static ref HASH_PATTERN: Regex = Regex::new(
        r"(?i)\b[0-9a-f]{32,64}\b"
    ).expect("Hash regex compilation failed");

    /// Base64: 20+ chars with optional padding
    /// Most general pattern, applied last
    static ref BASE64_PATTERN: Regex = Regex::new(
        r"[A-Za-z0-9+/]{20,}={0,2}"
    ).expect("Base64 regex compilation failed");

    /// IPv4 and IPv6 addresses
    static ref IP_PATTERN: Regex = Regex::new(
        r"(?i)\b((?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)|([0-9a-f]{0,4}:){2,7}[0-9a-f]{1,4})\b"
    ).expect("IP regex compilation failed");

    /// URLs (simple heuristic)
    static ref URL_PATTERN: Regex = Regex::new(
        r"https?://[A-Za-z0-9._~:/?#\[\]@!$&'()*+,;=%-]+"
    ).expect("URL regex compilation failed");

    /// Domains (fallback when URL not matched)
    static ref DOMAIN_PATTERN: Regex = Regex::new(
        r"\b([A-Za-z0-9-]+\.)+[A-Za-z]{2,}\b"
    ).expect("Domain regex compilation failed");

    /// Emails
    static ref EMAIL_PATTERN: Regex = Regex::new(
        r"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}"
    ).expect("Email regex compilation failed");

    /// ISO 8601 timestamps (UTC or with offset)
    static ref TIMESTAMP_PATTERN: Regex = Regex::new(
        r"\b\d{4}-\d{2}-\d{2}[Tt ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?\b"
    ).expect("Timestamp regex compilation failed");

    /// Long numeric sequences (IDs, counters)
    static ref NUMERIC_PATTERN: Regex = Regex::new(
        r"\b\d{4,}\b"
    ).expect("Numeric regex compilation failed");

    /// Cloud resource identifiers (AWS ARN style)
    static ref CLOUD_ID_PATTERN: Regex = Regex::new(
        r#"\barn:[A-Za-z0-9-]+:[A-Za-z0-9-]*:[A-Za-z0-9-]*:\d{12}:[^\s"']+\b"#
    ).expect("Cloud ID regex compilation failed");
}

/// Token replacements (constant byte slices)
const JWT_TOKEN: &[u8] = b"<JWT>";
const UUID_TOKEN: &[u8] = b"<UUID>";
const HASH_TOKEN: &[u8] = b"<HASH>";
const BASE64_TOKEN: &[u8] = b"<B64>";
const IP_TOKEN: &[u8] = b"<IP>";
const URL_TOKEN: &[u8] = b"<URL>";
const DOMAIN_TOKEN: &[u8] = b"<DOMAIN>";
const EMAIL_TOKEN: &[u8] = b"<EMAIL>";
const TIMESTAMP_TOKEN: &[u8] = b"<TS>";
const NUMERIC_TOKEN: &[u8] = b"<NUM>";
const CLOUD_ID_TOKEN: &[u8] = b"<CLOUD>";

/// Configuration for which high-entropy patterns to replace
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
pub struct TokenizerConfig {
    /// Replace UUIDs with `<UUID>`
    pub enable_uuid: bool,
    /// Replace hex hashes (32-64 chars) with `<HASH>`
    pub enable_hash: bool,
    /// Replace Base64 strings (20+ chars) with `<B64>`
    pub enable_base64: bool,
    /// Replace JWTs with `<JWT>`
    pub enable_jwt: bool,
    /// Replace IPv4/IPv6 addresses with `<IP>`
    pub enable_ip: bool,
    /// Replace URLs with `<URL>`
    pub enable_url: bool,
    /// Replace bare domains with `<DOMAIN>`
    pub enable_domain: bool,
    /// Replace emails with `<EMAIL>`
    pub enable_email: bool,
    /// Replace ISO 8601 timestamps with `<TS>`
    pub enable_timestamp: bool,
    /// Bucket long numeric sequences to `<NUM>`
    pub enable_numeric: bool,
    /// Replace cloud resource identifiers (e.g., AWS ARN) with `<CLOUD>`
    pub enable_cloud_id: bool,
    /// Canonicalize JSON by sorting keys alphabetically
    /// This ensures `{"b":1,"a":2}` == `{"a":2,"b":1}` for compression
    pub enable_json_canonicalization: bool,
}

impl Default for TokenizerConfig {
    fn default() -> Self {
        Self {
            enable_uuid: true,
            enable_hash: true,
            enable_base64: true,
            enable_jwt: true,
            enable_ip: true,
            enable_url: true,
            enable_domain: true,
            enable_email: true,
            enable_timestamp: true,
            enable_numeric: true,
            enable_cloud_id: true,
            enable_json_canonicalization: true,
        }
    }
}

impl TokenizerConfig {
    /// Create a config with all patterns disabled
    pub fn none() -> Self {
        Self {
            enable_uuid: false,
            enable_hash: false,
            enable_base64: false,
            enable_jwt: false,
            enable_ip: false,
            enable_url: false,
            enable_domain: false,
            enable_email: false,
            enable_timestamp: false,
            enable_numeric: false,
            enable_cloud_id: false,
            enable_json_canonicalization: false,
        }
    }

    /// Create a config with all patterns enabled
    pub fn all() -> Self {
        Self::default()
    }

    /// Check if any pattern is enabled
    pub fn has_any_enabled(&self) -> bool {
        self.enable_uuid
            || self.enable_hash
            || self.enable_base64
            || self.enable_jwt
            || self.enable_ip
            || self.enable_url
            || self.enable_domain
            || self.enable_email
            || self.enable_timestamp
            || self.enable_numeric
            || self.enable_cloud_id
            || self.enable_json_canonicalization
    }
}

/// Statistics about tokenization operations
#[derive(Debug, Clone, Copy, Default)]
pub struct TokenizerStats {
    /// Number of JWT tokens replaced
    pub jwt_count: u64,
    /// Number of UUID tokens replaced
    pub uuid_count: u64,
    /// Number of hash tokens replaced
    pub hash_count: u64,
    /// Number of Base64 tokens replaced
    pub base64_count: u64,
    /// Number of IP tokens replaced
    pub ip_count: u64,
    /// Number of URL tokens replaced
    pub url_count: u64,
    /// Number of domain tokens replaced
    pub domain_count: u64,
    /// Number of email tokens replaced
    pub email_count: u64,
    /// Number of timestamp tokens replaced
    pub timestamp_count: u64,
    /// Number of numeric buckets replaced
    pub numeric_count: u64,
    /// Number of cloud identifiers replaced
    pub cloud_id_count: u64,
    /// Number of JSON objects canonicalized
    pub json_canonicalized_count: u64,
    /// Total bytes saved by tokenization
    pub bytes_saved: u64,
}

/// Atomic stats for thread-safe updates
struct AtomicStats {
    jwt_count: AtomicU64,
    uuid_count: AtomicU64,
    hash_count: AtomicU64,
    base64_count: AtomicU64,
    ip_count: AtomicU64,
    url_count: AtomicU64,
    domain_count: AtomicU64,
    email_count: AtomicU64,
    timestamp_count: AtomicU64,
    numeric_count: AtomicU64,
    cloud_id_count: AtomicU64,
    json_canonicalized_count: AtomicU64,
    bytes_saved: AtomicU64,
}

impl Default for AtomicStats {
    fn default() -> Self {
        Self {
            jwt_count: AtomicU64::new(0),
            uuid_count: AtomicU64::new(0),
            hash_count: AtomicU64::new(0),
            base64_count: AtomicU64::new(0),
            ip_count: AtomicU64::new(0),
            url_count: AtomicU64::new(0),
            domain_count: AtomicU64::new(0),
            email_count: AtomicU64::new(0),
            timestamp_count: AtomicU64::new(0),
            numeric_count: AtomicU64::new(0),
            cloud_id_count: AtomicU64::new(0),
            json_canonicalized_count: AtomicU64::new(0),
            bytes_saved: AtomicU64::new(0),
        }
    }
}

impl AtomicStats {
    fn load(&self) -> TokenizerStats {
        TokenizerStats {
            jwt_count: self.jwt_count.load(Ordering::Relaxed),
            uuid_count: self.uuid_count.load(Ordering::Relaxed),
            hash_count: self.hash_count.load(Ordering::Relaxed),
            base64_count: self.base64_count.load(Ordering::Relaxed),
            ip_count: self.ip_count.load(Ordering::Relaxed),
            url_count: self.url_count.load(Ordering::Relaxed),
            domain_count: self.domain_count.load(Ordering::Relaxed),
            email_count: self.email_count.load(Ordering::Relaxed),
            timestamp_count: self.timestamp_count.load(Ordering::Relaxed),
            numeric_count: self.numeric_count.load(Ordering::Relaxed),
            cloud_id_count: self.cloud_id_count.load(Ordering::Relaxed),
            json_canonicalized_count: self.json_canonicalized_count.load(Ordering::Relaxed),
            bytes_saved: self.bytes_saved.load(Ordering::Relaxed),
        }
    }

    fn reset(&self) {
        self.jwt_count.store(0, Ordering::Relaxed);
        self.uuid_count.store(0, Ordering::Relaxed);
        self.hash_count.store(0, Ordering::Relaxed);
        self.base64_count.store(0, Ordering::Relaxed);
        self.ip_count.store(0, Ordering::Relaxed);
        self.url_count.store(0, Ordering::Relaxed);
        self.domain_count.store(0, Ordering::Relaxed);
        self.email_count.store(0, Ordering::Relaxed);
        self.timestamp_count.store(0, Ordering::Relaxed);
        self.numeric_count.store(0, Ordering::Relaxed);
        self.cloud_id_count.store(0, Ordering::Relaxed);
        self.json_canonicalized_count.store(0, Ordering::Relaxed);
        self.bytes_saved.store(0, Ordering::Relaxed);
    }
}

/// Recursively sort JSON object keys alphabetically.
/// Returns the canonical JSON as a byte vector.
/// If parsing fails, returns None (caller should use original data).
fn canonicalize_json(data: &[u8]) -> Option<Vec<u8>> {
    // Try to parse as JSON
    let value: serde_json::Value = serde_json::from_slice(data).ok()?;

    // Recursively sort object keys
    fn sort_value(v: serde_json::Value) -> serde_json::Value {
        match v {
            serde_json::Value::Object(map) => {
                // Convert to BTreeMap for sorted keys, recursively sort values
                let sorted: serde_json::Map<String, serde_json::Value> = map
                    .into_iter()
                    .map(|(k, v)| (k, sort_value(v)))
                    .collect::<std::collections::BTreeMap<_, _>>()
                    .into_iter()
                    .collect();
                serde_json::Value::Object(sorted)
            }
            serde_json::Value::Array(arr) => {
                serde_json::Value::Array(arr.into_iter().map(sort_value).collect())
            }
            other => other,
        }
    }

    let sorted = sort_value(value);
    // Serialize without extra whitespace for consistent compression
    serde_json::to_vec(&sorted).ok()
}

/// Tokenizer preprocesses data to replace high-entropy fields with fixed tokens.
///
/// This improves compression-based anomaly detection by normalizing random-looking
/// data (UUIDs, hashes, tokens) so the algorithm can focus on structural patterns.
///
/// # Example
///
/// ```
/// use cbad_core::tokenizer::{Tokenizer, TokenizerConfig};
///
/// let tokenizer = Tokenizer::new(TokenizerConfig::default());
/// let input = b"user_id: 550e8400-e29b-41d4-a716-446655440000";
/// let output = tokenizer.tokenize(input);
/// assert_eq!(&output, b"user_id: <UUID>");
/// ```
pub struct Tokenizer {
    config: TokenizerConfig,
    stats: AtomicStats,
}

impl Tokenizer {
    /// Create a new tokenizer with the given configuration
    pub fn new(config: TokenizerConfig) -> Self {
        Self {
            config,
            stats: AtomicStats::default(),
        }
    }

    /// Create a tokenizer with all patterns enabled (default)
    pub fn default_all() -> Self {
        Self::new(TokenizerConfig::default())
    }

    /// Tokenize the input data, replacing high-entropy patterns with tokens.
    ///
    /// Returns a new Vec<u8> - the original data is not modified.
    /// If no patterns are enabled, returns a copy of the input.
    pub fn tokenize(&self, data: &[u8]) -> Vec<u8> {
        // Fast path: if nothing enabled, return copy of original
        if !self.config.has_any_enabled() {
            return data.to_vec();
        }

        let original_len = data.len();
        let mut result = data.to_vec();

        // JSON canonicalization FIRST (before other patterns)
        // This ensures {"b":1,"a":2} == {"a":2,"b":1} for consistent compression
        if self.config.enable_json_canonicalization {
            if let Some(canonical) = canonicalize_json(&result) {
                result = canonical;
                self.stats
                    .json_canonicalized_count
                    .fetch_add(1, Ordering::Relaxed);
            }
            // If parsing fails, continue with original data
        }

        // Apply patterns in order (most specific first)
        // JWT must come before Base64 (JWTs contain Base64-like segments)
        if self.config.enable_timestamp {
            let matches = TIMESTAMP_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = TIMESTAMP_PATTERN
                    .replace_all(&result, TIMESTAMP_TOKEN)
                    .into_owned();
                self.stats
                    .timestamp_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        if self.config.enable_jwt {
            let matches = JWT_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = JWT_PATTERN.replace_all(&result, JWT_TOKEN).into_owned();
                self.stats
                    .jwt_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // UUID must come before Hash (UUIDs are specific hex patterns)
        if self.config.enable_uuid {
            let matches = UUID_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = UUID_PATTERN.replace_all(&result, UUID_TOKEN).into_owned();
                self.stats
                    .uuid_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Hash patterns (after UUID to avoid conflicts)
        if self.config.enable_hash {
            let matches = HASH_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = HASH_PATTERN.replace_all(&result, HASH_TOKEN).into_owned();
                self.stats
                    .hash_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Cloud IDs (e.g., AWS ARN)
        if self.config.enable_cloud_id {
            let matches = CLOUD_ID_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = CLOUD_ID_PATTERN
                    .replace_all(&result, CLOUD_ID_TOKEN)
                    .into_owned();
                self.stats
                    .cloud_id_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Emails before URLs to avoid partial replacements
        if self.config.enable_email {
            let matches = EMAIL_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = EMAIL_PATTERN.replace_all(&result, EMAIL_TOKEN).into_owned();
                self.stats
                    .email_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // URLs before domains
        if self.config.enable_url {
            let matches = URL_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = URL_PATTERN.replace_all(&result, URL_TOKEN).into_owned();
                self.stats
                    .url_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        if self.config.enable_domain {
            let matches = DOMAIN_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = DOMAIN_PATTERN
                    .replace_all(&result, DOMAIN_TOKEN)
                    .into_owned();
                self.stats
                    .domain_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        if self.config.enable_ip {
            let matches = IP_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = IP_PATTERN.replace_all(&result, IP_TOKEN).into_owned();
                self.stats
                    .ip_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Base64 last (most general pattern)
        if self.config.enable_base64 {
            let matches = BASE64_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = BASE64_PATTERN
                    .replace_all(&result, BASE64_TOKEN)
                    .into_owned();
                self.stats
                    .base64_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Numeric buckets last to avoid touching structured tokens
        if self.config.enable_numeric {
            let matches = NUMERIC_PATTERN.find_iter(&result).count();
            if matches > 0 {
                result = NUMERIC_PATTERN
                    .replace_all(&result, NUMERIC_TOKEN)
                    .into_owned();
                self.stats
                    .numeric_count
                    .fetch_add(matches as u64, Ordering::Relaxed);
            }
        }

        // Track bytes saved
        let bytes_saved = original_len.saturating_sub(result.len());
        if bytes_saved > 0 {
            self.stats
                .bytes_saved
                .fetch_add(bytes_saved as u64, Ordering::Relaxed);
        }

        result
    }

    /// Get current tokenization statistics
    pub fn stats(&self) -> TokenizerStats {
        self.stats.load()
    }

    /// Reset all statistics to zero
    pub fn reset_stats(&self) {
        self.stats.reset();
    }

    /// Get the current configuration
    pub fn config(&self) -> TokenizerConfig {
        self.config
    }
}

impl Default for Tokenizer {
    fn default() -> Self {
        Self::new(TokenizerConfig::default())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_uuid_replacement() {
        let tokenizer = Tokenizer::default();
        let input = b"user: 550e8400-e29b-41d4-a716-446655440000 logged in";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"user: <UUID> logged in");
        assert_eq!(tokenizer.stats().uuid_count, 1);
    }

    #[test]
    fn test_uuid_case_insensitive() {
        let tokenizer = Tokenizer::default();
        let input = b"ID: 550E8400-E29B-41D4-A716-446655440000";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"ID: <UUID>");
    }

    #[test]
    fn test_hash_replacement() {
        let tokenizer = Tokenizer::default();
        // SHA-256 hash (64 chars)
        let input = b"hash: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"hash: <HASH>");
        assert_eq!(tokenizer.stats().hash_count, 1);
    }

    #[test]
    fn test_jwt_replacement() {
        let tokenizer = Tokenizer::default();
        let jwt = b"token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
        let output = tokenizer.tokenize(jwt);
        assert_eq!(&output, b"token: <JWT>");
        assert_eq!(tokenizer.stats().jwt_count, 1);
    }

    #[test]
    fn test_base64_replacement() {
        let tokenizer = Tokenizer::default();
        let input = b"data: SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB0ZXN0";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"data: <B64>");
        assert_eq!(tokenizer.stats().base64_count, 1);
    }

    #[test]
    fn test_multiple_patterns() {
        let tokenizer = Tokenizer::default();
        let input = b"user: 550e8400-e29b-41d4-a716-446655440000, hash: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"user: <UUID>, hash: <HASH>");
        let stats = tokenizer.stats();
        assert_eq!(stats.uuid_count, 1);
        assert_eq!(stats.hash_count, 1);
    }

    #[test]
    fn test_ip_and_url_replacement() {
        let tokenizer = Tokenizer::default();
        let input = b"client=192.168.0.1 url=https://example.com/path";
        let output = tokenizer.tokenize(input);
        let out_str = String::from_utf8_lossy(&output);
        assert!(out_str.contains("<IP>"));
        assert!(out_str.contains("<URL>"));
    }

    #[test]
    fn test_email_and_domain_replacement() {
        let tokenizer = Tokenizer::default();
        let input = b"user=john.doe@example.com host=api.example.com";
        let output = tokenizer.tokenize(input);
        let out_str = String::from_utf8_lossy(&output);
        assert!(out_str.contains("<EMAIL>"));
        assert!(out_str.contains("<DOMAIN>"));
    }

    #[test]
    fn test_timestamp_and_numeric_bucket() {
        let tokenizer = Tokenizer::default();
        let input = b"time=2024-10-01T12:34:56Z id=12345678";
        let output = tokenizer.tokenize(input);
        let out_str = String::from_utf8_lossy(&output);
        assert!(out_str.contains("<TS>"));
        assert!(out_str.contains("<NUM>"));
    }

    #[test]
    fn test_cloud_id_replacement() {
        let tokenizer = Tokenizer::default();
        let input = b"arn:aws:lambda:us-east-1:123456789012:function:my-func";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, b"<CLOUD>");
        assert_eq!(tokenizer.stats().cloud_id_count, 1);
    }

    #[test]
    fn test_disabled_patterns() {
        let config = TokenizerConfig {
            enable_uuid: true,
            enable_hash: false,
            enable_base64: false,
            enable_jwt: false,
            enable_json_canonicalization: false,
            ..TokenizerConfig::none()
        };
        let tokenizer = Tokenizer::new(config);
        let input = b"uuid: 550e8400-e29b-41d4-a716-446655440000, hash: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";
        let output = tokenizer.tokenize(input);
        // UUID replaced, hash NOT replaced
        assert!(output.windows(6).any(|w| w == b"<UUID>"));
        assert!(output.windows(6).any(|w| w == b"<HASH>") == false);
    }

    #[test]
    fn test_no_patterns_enabled() {
        let tokenizer = Tokenizer::new(TokenizerConfig::none());
        let input = b"uuid: 550e8400-e29b-41d4-a716-446655440000";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, input);
    }

    #[test]
    fn test_bytes_saved() {
        let tokenizer = Tokenizer::default();
        let uuid = "550e8400-e29b-41d4-a716-446655440000"; // 36 chars
        let input = format!("id: {}", uuid);
        let _output = tokenizer.tokenize(input.as_bytes());
        // <UUID> is 6 chars, UUID is 36 chars, saved 30 bytes
        assert_eq!(tokenizer.stats().bytes_saved, 30);
    }

    #[test]
    fn test_json_with_mixed_fields() {
        let tokenizer = Tokenizer::default();
        let json = br#"{"user_id":"550e8400-e29b-41d4-a716-446655440000","action":"login","timestamp":"2024-01-01T00:00:00Z"}"#;
        let output = tokenizer.tokenize(json);
        let output_str = String::from_utf8_lossy(&output);
        assert!(output_str.contains("<UUID>"));
        assert!(output_str.contains("login")); // Non-entropy fields preserved
        assert!(output_str.contains("timestamp")); // Field names preserved
    }

    #[test]
    fn test_stats_reset() {
        let tokenizer = Tokenizer::default();
        let input = b"550e8400-e29b-41d4-a716-446655440000";
        let _ = tokenizer.tokenize(input);
        assert_eq!(tokenizer.stats().uuid_count, 1);
        tokenizer.reset_stats();
        assert_eq!(tokenizer.stats().uuid_count, 0);
    }

    #[test]
    fn test_json_canonicalization_sorts_keys() {
        let tokenizer = Tokenizer::default();
        // Keys in non-alphabetical order
        let input = br#"{"zebra":1,"apple":2,"mango":3}"#;
        let output = tokenizer.tokenize(input);
        // Keys should be sorted alphabetically
        let expected = br#"{"apple":2,"mango":3,"zebra":1}"#;
        assert_eq!(&output, expected);
        assert_eq!(tokenizer.stats().json_canonicalized_count, 1);
    }

    #[test]
    fn test_json_canonicalization_nested_objects() {
        let tokenizer = Tokenizer::default();
        // Nested objects with unsorted keys
        let input = br#"{"z":{"b":1,"a":2},"a":{"d":3,"c":4}}"#;
        let output = tokenizer.tokenize(input);
        // All levels should be sorted
        let expected = br#"{"a":{"c":4,"d":3},"z":{"a":2,"b":1}}"#;
        assert_eq!(&output, expected);
    }

    #[test]
    fn test_json_canonicalization_with_arrays() {
        let tokenizer = Tokenizer::default();
        // Arrays preserve order, but object keys inside are sorted
        let input = br#"{"items":[{"z":1,"a":2},{"b":3}],"name":"test"}"#;
        let output = tokenizer.tokenize(input);
        let expected = br#"{"items":[{"a":2,"z":1},{"b":3}],"name":"test"}"#;
        assert_eq!(&output, expected);
    }

    #[test]
    fn test_json_canonicalization_and_tokenization_combined() {
        let tokenizer = Tokenizer::default();
        // JSON with unsorted keys AND a UUID
        let input = br#"{"id":"550e8400-e29b-41d4-a716-446655440000","action":"login"}"#;
        let output = tokenizer.tokenize(input);
        let output_str = String::from_utf8_lossy(&output);
        // Keys should be sorted (action before id)
        assert!(output_str.starts_with(r#"{"action":"login","id":"#));
        // UUID should be tokenized
        assert!(output_str.contains("<UUID>"));
    }

    #[test]
    fn test_json_canonicalization_disabled() {
        let config = TokenizerConfig {
            enable_json_canonicalization: false,
            ..TokenizerConfig::default()
        };
        let tokenizer = Tokenizer::new(config);
        let input = br#"{"z":1,"a":2}"#;
        let output = tokenizer.tokenize(input);
        // Keys should NOT be sorted
        assert_eq!(&output, input);
        assert_eq!(tokenizer.stats().json_canonicalized_count, 0);
    }

    #[test]
    fn test_json_canonicalization_invalid_json_passthrough() {
        let tokenizer = Tokenizer::default();
        // Not valid JSON - should pass through unchanged (except other tokenization)
        let input = b"This is not JSON at all";
        let output = tokenizer.tokenize(input);
        assert_eq!(&output, input);
        assert_eq!(tokenizer.stats().json_canonicalized_count, 0);
    }

    #[test]
    fn test_json_canonicalization_equivalent_outputs() {
        let tokenizer = Tokenizer::default();
        // Two JSONs with same content but different key orders
        let json1 = br#"{"b":1,"a":2}"#;
        let json2 = br#"{"a":2,"b":1}"#;

        let output1 = tokenizer.tokenize(json1);
        tokenizer.reset_stats();
        let output2 = tokenizer.tokenize(json2);

        // Both should produce the same canonical output
        assert_eq!(output1, output2);
    }
}
