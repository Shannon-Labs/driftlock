// Fallback compression adapters
// These are simple wrappers around standard compression libraries
// Used when OpenZL is not suitable (e.g., truly unstructured binary data)

use super::{CompressionAdapter, CompressionError, Result};
use std::io::Write;

/// Zstd compression adapter (fallback)
///
/// TODO: Implement actual zstd compression using the `zstd` crate
/// For now, this is a placeholder that shows the interface
pub struct ZstdAdapter {
    level: i32,
}

impl Default for ZstdAdapter {
    fn default() -> Self {
        Self::new()
    }
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
        if data.is_empty() {
            return Ok(Vec::new());
        }

        zstd::encode_all(data, self.level)
            .map_err(|e| CompressionError::CompressionFailed(format!("Zstd compression error: {}", e)))
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        zstd::decode_all(data)
            .map_err(|e| CompressionError::DecompressionFailed(format!("Zstd decompression error: {}", e)))
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

impl Default for Lz4Adapter {
    fn default() -> Self {
        Self::new()
    }
}

impl Lz4Adapter {
    pub fn new() -> Self {
        Self
    }
}

impl CompressionAdapter for Lz4Adapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        let compressed = lz4::block::compress(data, None, false)
            .map_err(|e| CompressionError::CompressionFailed(format!("Lz4 compression error: {}", e)))?;
        
        Ok(compressed)
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        // For lz4 block compression, we need to know the uncompressed size
        // Since we don't store it separately, we'll use a reasonable estimate
        let uncompressed_size = data.len() * 4; // Conservative estimate
        lz4::block::decompress(data, Some(uncompressed_size as i32))
            .map_err(|e| CompressionError::DecompressionFailed(format!("Lz4 decompression error: {}", e)))
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

impl Default for GzipAdapter {
    fn default() -> Self {
        Self::new()
    }
}

impl GzipAdapter {
    pub fn new() -> Self {
        Self
    }
}

impl CompressionAdapter for GzipAdapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        let mut encoder = flate2::write::GzEncoder::new(Vec::new(), flate2::Compression::default());
        encoder.write_all(data)
            .map_err(|e| CompressionError::CompressionFailed(format!("Gzip write error: {}", e)))?;
        
        encoder.finish()
            .map_err(|e| CompressionError::CompressionFailed(format!("Gzip finish error: {}", e)))
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        let mut decoder = flate2::write::GzDecoder::new(Vec::new());
        decoder.write_all(data)
            .map_err(|e| CompressionError::DecompressionFailed(format!("Gzip write error: {}", e)))?;
        
        decoder.finish()
            .map_err(|e| CompressionError::DecompressionFailed(format!("Gzip finish error: {}", e)))
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
    fn test_adapters_work_correctly() {
        let data = b"test data for compression";

        // Zstd
        let zstd = ZstdAdapter::new();
        let compressed = zstd.compress(data).expect("zstd compress should work");
        let decompressed = zstd.decompress(&compressed).expect("zstd decompress should work");
        assert_eq!(data, decompressed.as_slice());
        assert_eq!(zstd.name(), "zstd");

        // Lz4
        let lz4 = Lz4Adapter::new();
        let compressed = lz4.compress(data).expect("lz4 compress should work");
        let decompressed = lz4.decompress(&compressed).expect("lz4 decompress should work");
        assert_eq!(data, decompressed.as_slice());
        assert_eq!(lz4.name(), "lz4");

        // Gzip
        let gzip = GzipAdapter::new();
        let compressed = gzip.compress(data).expect("gzip compress should work");
        let decompressed = gzip.decompress(&compressed).expect("gzip decompress should work");
        assert_eq!(data, decompressed.as_slice());
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
