use thiserror::Error;

/// Errors that can occur when using AI services
#[derive(Error, Debug)]
pub enum AiError {
    /// HTTP request failed
    #[error("HTTP request failed: {0}")]
    Http(#[from] reqwest::Error),

    /// JSON serialization/deserialization failed
    #[error("JSON error: {0}")]
    Json(#[from] serde_json::Error),

    /// API returned an error response
    #[error("API error: {0}")]
    Api(String),

    /// Empty response from API
    #[error("Empty response from API")]
    EmptyResponse,

    /// Missing API key
    #[error("Missing API key: {0}")]
    MissingApiKey(String),

    /// Invalid configuration
    #[error("Invalid configuration: {0}")]
    InvalidConfig(String),

    /// Rate limit exceeded
    #[error("Rate limit exceeded: {0}")]
    RateLimit(String),

    /// Cost limit exceeded
    #[error("Cost limit exceeded: {0}")]
    CostLimit(String),

    /// Model not allowed for plan
    #[error("Model {0} not allowed for plan {1}")]
    ModelNotAllowed(String, String),

    /// Generic error
    #[error("AI error: {0}")]
    Other(String),
}

/// Type alias for Results with AiError
pub type Result<T> = std::result::Result<T, AiError>;
