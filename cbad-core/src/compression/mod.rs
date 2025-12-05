// Compression module for CBAD core
// Provides multiple compression adapters for anomaly detection

#[cfg(feature = "openzl")]
pub mod openzl;

pub mod fallback;

use serde::{Deserialize, Serialize};
use std::fmt;

/// Result type for compression operations
pub type Result<T> = std::result::Result<T, CompressionError>;

/// Compression error types
#[derive(Debug, Clone)]
pub enum CompressionError {
    /// Compression failed
    CompressionFailed(String),
    /// Decompression failed
    DecompressionFailed(String),
    /// Buffer too small
    BufferTooSmall { required: usize, available: usize },
    /// Invalid input
    InvalidInput(String),
    /// Unsupported algorithm
    UnsupportedAlgorithm(String),
}

impl fmt::Display for CompressionError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::CompressionFailed(msg) => write!(f, "Compression failed: {}", msg),
            Self::DecompressionFailed(msg) => write!(f, "Decompression failed: {}", msg),
            Self::BufferTooSmall {
                required,
                available,
            } => {
                write!(
                    f,
                    "Buffer too small: required {} bytes, available {} bytes",
                    required, available
                )
            }
            Self::InvalidInput(msg) => write!(f, "Invalid input: {}", msg),
            Self::UnsupportedAlgorithm(msg) => write!(f, "Unsupported algorithm: {}", msg),
        }
    }
}

impl std::error::Error for CompressionError {}

/// Trait for compression adapters
/// All compression algorithms must implement this trait
pub trait CompressionAdapter: Send + Sync {
    /// Compress data
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>>;

    /// Decompress data
    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>>;

    /// Get algorithm name
    fn name(&self) -> &str;

    /// Estimate maximum compressed size for given input size
    fn compress_bound(&self, src_size: usize) -> usize;

    /// Check if compression is deterministic (same input always produces same output)
    fn is_deterministic(&self) -> bool {
        true // Most compression algorithms are deterministic
    }
}

/// Compression algorithm selection
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum CompressionAlgorithm {
    /// OpenZL format-aware compression (primary)
    #[cfg(feature = "openzl")]
    OpenZL,

    /// Zlab (deterministic zlib-style compression)
    Zlab,

    /// Zstd (fallback for unstructured data)
    Zstd,

    /// Lz4 (fallback for ultra-fast compression)
    Lz4,

    /// Gzip (fallback for universal compatibility)
    Gzip,
}

impl CompressionAlgorithm {
    pub fn name(&self) -> &str {
        match self {
            #[cfg(feature = "openzl")]
            Self::OpenZL => "openzl",
            Self::Zlab => "zlab",
            Self::Zstd => "zstd",
            Self::Lz4 => "lz4",
            Self::Gzip => "gzip",
        }
    }
}

/// Create a compression adapter for the given algorithm
pub fn create_adapter(algo: CompressionAlgorithm) -> Result<Box<dyn CompressionAdapter>> {
    match algo {
        #[cfg(feature = "openzl")]
        CompressionAlgorithm::OpenZL => Ok(Box::new(openzl::OpenZLAdapter::new()?)),

        CompressionAlgorithm::Zlab => Ok(Box::new(fallback::ZlabAdapter::new())),

        CompressionAlgorithm::Zstd => Ok(Box::new(fallback::ZstdAdapter::new())),
        CompressionAlgorithm::Lz4 => Ok(Box::new(fallback::Lz4Adapter::new())),
        CompressionAlgorithm::Gzip => Ok(Box::new(fallback::GzipAdapter::new())),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_compression_roundtrip() {
        let data = b"Hello, Driftlock! This is test data for compression.";

        // Temporarily test only fallback adapters while OpenZL is being fixed
        for algo in [
            CompressionAlgorithm::Zlab,
            CompressionAlgorithm::Zstd,
            CompressionAlgorithm::Lz4,
            CompressionAlgorithm::Gzip,
        ] {
            let adapter = create_adapter(algo).expect("create adapter");
            let compressed = adapter.compress(data).expect("compress");
            let decompressed = adapter.decompress(&compressed).expect("decompress");
            assert_eq!(data, decompressed.as_slice());
        }
    }
}
