# Driftlock Auth

Authentication library for Driftlock providing Firebase JWT verification and API key validation.

## Features

- **Firebase JWT Verification**: Verify Firebase ID tokens with automatic public key fetching and caching
- **API Key Management**: Generate, validate, and manage API keys with scopes

## Firebase JWT Verification

### Overview

The Firebase JWT verification implementation:

1. Fetches Google public keys from the Firebase certificate endpoint
2. Caches keys with automatic TTL parsing from Cache-Control headers
3. Verifies JWT signatures using the RS256 algorithm
4. Validates all required claims (issuer, audience, expiration, etc.)
5. Returns structured user information

### Usage

```rust
use driftlock_auth::{FirebaseAuth, FirebaseUser};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize Firebase authenticator
    let firebase = FirebaseAuth::new("your-project-id");

    // Verify a token
    match firebase.verify_token("eyJhbGciOiJSUzI1NiIsInR5cCI6...").await {
        Ok(user) => {
            println!("User ID: {}", user.uid);
            println!("Email: {:?}", user.email);
            println!("Email verified: {}", user.email_verified);
        }
        Err(e) => {
            eprintln!("Authentication failed: {}", e);
        }
    }

    Ok(())
}
```

### Implementation Details

#### Key Caching

Public keys are automatically fetched from:
```
https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com
```

The cache TTL is parsed from the `Cache-Control: max-age=<seconds>` header (typically 1 hour).

Keys are cached in memory using an `Arc<RwLock<...>>` for thread-safe concurrent access.

#### JWT Validation

The following claims are validated:

| Claim | Validation |
|-------|-----------|
| `iss` | Must be `https://securetoken.google.com/{project_id}` |
| `aud` | Must match the project ID |
| `exp` | Token must not be expired |
| `iat` | Issued time must be in the past |
| `sub` | Must be non-empty (this is the user ID) |

#### Error Handling

```rust
use driftlock_auth::FirebaseError;

match firebase.verify_token(token).await {
    Err(FirebaseError::InvalidToken) => {
        // Malformed or invalid token
    }
    Err(FirebaseError::TokenExpired) => {
        // Token has expired
    }
    Err(FirebaseError::InvalidIssuer) => {
        // Token not from Firebase
    }
    Err(FirebaseError::InvalidAudience) => {
        // Token for different project
    }
    Err(FirebaseError::KeyFetchError(msg)) => {
        // Failed to fetch public keys
    }
    // ... other errors
}
```

### API Reference

#### `FirebaseAuth`

```rust
impl FirebaseAuth {
    /// Create a new Firebase authenticator
    pub fn new(project_id: &str) -> Self;

    /// Verify a Firebase ID token and return user information
    pub async fn verify_token(&self, token: &str) -> Result<FirebaseUser, FirebaseError>;

    /// Force refresh the cached keys (useful for testing)
    pub async fn refresh_keys(&self) -> Result<(), FirebaseError>;
}
```

#### `FirebaseUser`

```rust
pub struct FirebaseUser {
    /// Firebase user ID (UID)
    pub uid: String,
    /// User's email address (if available)
    pub email: Option<String>,
    /// Whether the email has been verified
    pub email_verified: bool,
}
```

#### `FirebaseError`

```rust
pub enum FirebaseError {
    InvalidToken,
    TokenExpired,
    KeyFetchError(String),
    MissingKeyId,
    KeyNotFound(String),
    InvalidIssuer,
    InvalidAudience,
    InvalidSubject,
    JwtError(jsonwebtoken::errors::Error),
    HttpError(reqwest::Error),
}
```

## Testing

Run tests:
```bash
cargo test -p driftlock-auth
```

Run the example:
```bash
# Get a real Firebase token from your app first
cargo run --example verify_token -- your-project-id "eyJhbGciOiJSUzI1NiIs..."
```

## Dependencies

- `jsonwebtoken` - JWT encoding/decoding with RS256 support
- `reqwest` - HTTP client for fetching public keys
- `tokio` - Async runtime
- `serde` / `serde_json` - Serialization
- `thiserror` - Error handling

## License

Apache-2.0
