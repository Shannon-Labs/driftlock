// Fallback compression adapters
// These are simple wrappers around standard compression libraries
// Used when OpenZL is not suitable (e.g., truly unstructured binary data)

use super::{CompressionAdapter, CompressionError, Result};

/// Zstd compression adapter (fallback)
///
/// TODO: Implement actual zstd compression using the `zstd` crate
/// For now, this is a placeholder that shows the interface
pub struct ZstdAdapter {
    level: i32,
}

impl ZstdAdapter {
    pub fn new() -> Self {
        Self { level: 3 } // Default compression level
    }

    pub fn with_level(level: i32) -> Self {
        Self { level }
    }
}

impl CompressionAdapter for ZstdAdapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual zstd compression
        // For now, return error to indicate not implemented
        Err(CompressionError::UnsupportedAlgorithm(
            "Zstd adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual zstd decompression
        Err(CompressionError::UnsupportedAlgorithm(
            "Zstd adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn name(&self) -> &str {
        "zstd"
    }

    fn compress_bound(&self, src_size: usize) -> usize {
        // Zstd bound formula
        src_size + (src_size >> 8) + (if src_size < 128 * 1024 { (128 * 1024 - src_size) >> 11 } else { 0 })
    }
}

/// Lz4 compression adapter (fallback)
///
/// TODO: Implement actual lz4 compression using the `lz4` crate
pub struct Lz4Adapter;

impl Lz4Adapter {
    pub fn new() -> Self {
        Self
    }
}

impl CompressionAdapter for Lz4Adapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual lz4 compression
        Err(CompressionError::UnsupportedAlgorithm(
            "Lz4 adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual lz4 decompression
        Err(CompressionError::UnsupportedAlgorithm(
            "Lz4 adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn name(&self) -> &str {
        "lz4"
    }

    fn compress_bound(&self, src_size: usize) -> usize {
        // Lz4 bound formula: src_size + (src_size / 255) + 16
        src_size + (src_size / 255) + 16
    }
}

/// Gzip compression adapter (fallback)
///
/// TODO: Implement actual gzip compression using the `flate2` crate
pub struct GzipAdapter;

impl GzipAdapter {
    pub fn new() -> Self {
        Self
    }
}

impl CompressionAdapter for GzipAdapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual gzip compression
        Err(CompressionError::UnsupportedAlgorithm(
            "Gzip adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // TODO: Replace with actual gzip decompression
        Err(CompressionError::UnsupportedAlgorithm(
            "Gzip adapter not yet implemented - use OpenZL".to_string(),
        ))
    }

    fn name(&self) -> &str {
        "gzip"
    }

    fn compress_bound(&self, src_size: usize) -> usize {
        // Gzip bound (conservative estimate)
        src_size + (src_size / 1000) + 12
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_adapters_return_not_implemented() {
        let data = b"test data";

        // Zstd
        let zstd = ZstdAdapter::new();
        assert!(zstd.compress(data).is_err());
        assert_eq!(zstd.name(), "zstd");

        // Lz4
        let lz4 = Lz4Adapter::new();
        assert!(lz4.compress(data).is_err());
        assert_eq!(lz4.name(), "lz4");

        // Gzip
        let gzip = GzipAdapter::new();
        assert!(gzip.compress(data).is_err());
        assert_eq!(gzip.name(), "gzip");
    }

    #[test]
    fn test_compress_bounds() {
        let src_size = 1024 * 1024; // 1MB

        let zstd = ZstdAdapter::new();
        assert!(zstd.compress_bound(src_size) > src_size);

        let lz4 = Lz4Adapter::new();
        assert!(lz4.compress_bound(src_size) > src_size);

        let gzip = GzipAdapter::new();
        assert!(gzip.compress_bound(src_size) > src_size);
    }
}
