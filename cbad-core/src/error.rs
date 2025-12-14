use crate::compression::CompressionError;
use thiserror::Error;

/// Canonical result type for production-ready CBAD APIs.
pub type Result<T> = std::result::Result<T, CbadError>;

/// Top-level error type to provide actionable messages to callers.
#[derive(Debug, Error)]
pub enum CbadError {
    /// Detector does not have enough baseline/window data yet.
    #[error("detector not ready: need {events_needed} events, have {events_have}")]
    NotReady {
        events_needed: usize,
        events_have: usize,
    },

    /// Configuration is invalid or inconsistent.
    #[error("invalid configuration: {0}")]
    InvalidConfig(String),

    /// Compression failed in the underlying adapter.
    #[error(transparent)]
    Compression(#[from] CompressionError),

    /// CSV parsing/formatting error.
    #[error(transparent)]
    Csv(#[from] csv::Error),

    /// A storage backend reported an error.
    #[error("storage error: {0}")]
    StorageError(String),

    /// Resource limits were exceeded (memory, channel size, etc.).
    #[error("resource exhausted: {resource} (limit {limit})")]
    ResourceExhausted {
        resource: &'static str,
        limit: usize,
    },

    /// IO failure while reading/writing datasets or state.
    #[error(transparent)]
    Io(#[from] std::io::Error),

    /// Serialization/deserialization failure.
    #[error("serialization error: {0}")]
    Serialization(String),

    /// Unsupported file or protocol format.
    #[error("unsupported format: {0}")]
    UnsupportedFormat(String),
}

impl CbadError {
    /// Helper to construct a not-ready error given baseline/window sizes.
    pub fn not_ready(events_needed: usize, events_have: usize) -> Self {
        Self::NotReady {
            events_needed,
            events_have,
        }
    }
}
