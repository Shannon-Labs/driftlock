// OpenZL compression adapter with C FFI bindings
// Wraps Meta's OpenZL format-aware compression framework

use super::{CompressionAdapter, CompressionError, Result};
use std::os::raw::{c_char, c_int, c_void};
use std::ptr;

// OpenZL C FFI bindings
#[repr(C)]
struct ZL_CCtx {
    _private: [u8; 0],
}

#[repr(C)]
struct ZL_DCtx {
    _private: [u8; 0],
}

// ZL_Report is a size_t that can be either a success (size) or error (negative)
type ZL_Report = usize;

#[link(name = "openzl", kind = "static")]
extern "C" {
    // Compression context
    fn ZL_CCtx_create() -> *mut ZL_CCtx;
    fn ZL_CCtx_free(cctx: *mut ZL_CCtx);

    fn ZL_CCtx_compress(
        cctx: *mut ZL_CCtx,
        dst: *mut c_void,
        dst_capacity: usize,
        src: *const c_void,
        src_size: usize,
    ) -> ZL_Report;

    fn ZL_compressBound(total_src_size: usize) -> usize;

    // Decompression context
    fn ZL_DCtx_create() -> *mut ZL_DCtx;
    fn ZL_DCtx_free(dctx: *mut ZL_DCtx);

    // Simple decompression API
    fn ZL_decompress(
        dst: *mut c_void,
        dst_capacity: usize,
        src: *const c_void,
        src_size: usize,
    ) -> ZL_Report;

    fn ZL_getDecompressedSize(compressed: *const c_void, c_size: usize) -> ZL_Report;

    // Error checking
    fn ZL_isError(report: ZL_Report) -> c_int;
    fn ZL_getErrorName(report: ZL_Report) -> *const c_char;
}

// Helper to check if ZL_Report is an error
fn is_error(report: ZL_Report) -> bool {
    unsafe { ZL_isError(report) != 0 }
}

// Helper to get error message from ZL_Report
fn get_error_message(report: ZL_Report) -> String {
    unsafe {
        let name_ptr = ZL_getErrorName(report);
        if name_ptr.is_null() {
            return format!("OpenZL error code: {}", report);
        }
        let c_str = std::ffi::CStr::from_ptr(name_ptr);
        c_str
            .to_str()
            .unwrap_or("Invalid UTF-8 in error message")
            .to_string()
    }
}

/// OpenZL compression adapter
pub struct OpenZLAdapter {
    cctx: *mut ZL_CCtx,
}

impl OpenZLAdapter {
    /// Create a new OpenZL adapter
    pub fn new() -> Result<Self> {
        let cctx = unsafe { ZL_CCtx_create() };
        if cctx.is_null() {
            return Err(CompressionError::CompressionFailed(
                "Failed to create OpenZL compression context".to_string(),
            ));
        }

        Ok(Self { cctx })
    }
}

impl Drop for OpenZLAdapter {
    fn drop(&mut self) {
        if !self.cctx.is_null() {
            unsafe {
                ZL_CCtx_free(self.cctx);
            }
        }
    }
}

unsafe impl Send for OpenZLAdapter {}
unsafe impl Sync for OpenZLAdapter {}

impl CompressionAdapter for OpenZLAdapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        let src_size = data.len();
        let dst_capacity = self.compress_bound(src_size);
        let mut dst = vec![0u8; dst_capacity];

        let result = unsafe {
            ZL_CCtx_compress(
                self.cctx,
                dst.as_mut_ptr() as *mut c_void,
                dst_capacity,
                data.as_ptr() as *const c_void,
                src_size,
            )
        };

        if is_error(result) {
            return Err(CompressionError::CompressionFailed(get_error_message(
                result,
            )));
        }

        let compressed_size = result;
        dst.truncate(compressed_size);
        Ok(dst)
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        // Get decompressed size first
        let decompressed_size = unsafe {
            ZL_getDecompressedSize(data.as_ptr() as *const c_void, data.len())
        };

        if is_error(decompressed_size) {
            return Err(CompressionError::DecompressionFailed(format!(
                "Failed to get decompressed size: {}",
                get_error_message(decompressed_size)
            )));
        }

        let mut dst = vec![0u8; decompressed_size];

        let result = unsafe {
            ZL_decompress(
                dst.as_mut_ptr() as *mut c_void,
                decompressed_size,
                data.as_ptr() as *const c_void,
                data.len(),
            )
        };

        if is_error(result) {
            return Err(CompressionError::DecompressionFailed(get_error_message(
                result,
            )));
        }

        let actual_size = result;
        if actual_size != decompressed_size {
            dst.truncate(actual_size);
        }

        Ok(dst)
    }

    fn name(&self) -> &str {
        "openzl"
    }

    fn compress_bound(&self, src_size: usize) -> usize {
        unsafe { ZL_compressBound(src_size) }
    }

    fn is_deterministic(&self) -> bool {
        // OpenZL is deterministic when using fixed compression plans
        true
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_openzl_create() {
        let adapter = OpenZLAdapter::new();
        assert!(adapter.is_ok(), "Should create OpenZL adapter");
    }

    #[test]
    fn test_openzl_compress_decompress() {
        let adapter = OpenZLAdapter::new().expect("create adapter");

        let original = b"Hello, OpenZL! This is a test of format-aware compression.";
        let compressed = adapter.compress(original).expect("compress");
        let decompressed = adapter.decompress(&compressed).expect("decompress");

        assert_eq!(original, decompressed.as_slice());
        println!(
            "OpenZL compression ratio: {:.2}x ({} bytes -> {} bytes)",
            original.len() as f64 / compressed.len() as f64,
            original.len(),
            compressed.len()
        );
    }

    #[test]
    fn test_openzl_empty_data() {
        let adapter = OpenZLAdapter::new().expect("create adapter");

        let empty: &[u8] = &[];
        let compressed = adapter.compress(empty).expect("compress empty");
        assert_eq!(compressed.len(), 0);

        let decompressed = adapter.decompress(empty).expect("decompress empty");
        assert_eq!(decompressed.len(), 0);
    }

    #[test]
    fn test_openzl_large_data() {
        let adapter = OpenZLAdapter::new().expect("create adapter");

        // Generate 1MB of structured data (simulating OTLP logs)
        let mut original = Vec::new();
        for i in 0..10000 {
            let log_entry = format!(
                r#"{{"timestamp":"2025-10-24T00:00:00Z","level":"info","service":"api-gateway","msg":"request completed","duration_ms":{},"request_id":"req-{}"}}"#,
                42 + (i % 100),
                i
            );
            original.extend_from_slice(log_entry.as_bytes());
            original.push(b'\n');
        }

        let compressed = adapter.compress(&original).expect("compress large data");
        let decompressed = adapter.decompress(&compressed).expect("decompress large data");

        assert_eq!(original, decompressed);
        println!(
            "Large data compression ratio: {:.2}x ({} bytes -> {} bytes)",
            original.len() as f64 / compressed.len() as f64,
            original.len(),
            compressed.len()
        );
    }
}
