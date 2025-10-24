// OpenZL compression adapter with C FFI bindings
// Wraps Meta's OpenZL format-aware compression framework

use super::{CompressionAdapter, CompressionError, Result};
use std::os::raw::{c_char, c_int, c_uint, c_void};

// OpenZL C FFI bindings
#[repr(C)]
struct ZL_CCtx {
    _private: [u8; 0],
}

#[repr(C)]
struct ZL_DCtx {
    _private: [u8; 0],
}

// ZL_Report is a complex struct type, not a simple usize
// From zl_errors.h: ZL_RESULT_DECLARE_TYPE(size_t) and typedef ZL_RESULT_OF(size_t) ZL_Report;
#[repr(C)]
struct ZL_Report {
    _code: ZL_ErrorCode,
    _value: usize,
}

#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum ZL_ErrorCode {
    NoError = 0,
    Generic = 1,
    SrcSizeTooSmall = 3,
    SrcSizeTooLarge = 4,
    DstCapacityTooSmall = 5,
    UserBufferAlignmentIncorrect = 6,
    DecompressionIncorrectApi = 7,
    UserBuffersInvalidNum = 8,
    InvalidName = 9,
    HeaderUnknown = 10,
    FrameParameterUnsupported = 11,
    Corruption = 12,
    CompressedChecksumWrong = 13,
    ContentChecksumWrong = 14,
    OutputsTooNumerous = 15,
    CompressionParameterInvalid = 20,
    ParameterInvalid = 21,
    OutputIdInvalid = 22,
    InvalidRequestSingleOutputFrameOnly = 23,
    OutputNotCommitted = 24,
    OutputNotReserved = 25,
    SegmenterInputNotConsumed = 26,
    GraphInvalid = 30,
    GraphNonserializable = 31,
    InvalidTransform = 32,
    GraphInvalidNumInputs = 33,
    SuccessorInvalid = 40,
    SuccessorAlreadySet = 41,
    SuccessorInvalidNumInputs = 42,
    InputTypeUnsupported = 43,
    GraphParameterInvalid = 44,
    NodeParameterInvalid = 50,
    NodeParameterInvalidValue = 51,
    TransformExecutionFailure = 52,
    CustomNodeDefinitionInvalid = 53,
    NodeUnexpectedInputType = 54,
    NodeInvalidInput = 55,
    NodeInvalid = 56,
    NodeExecutionInvalidOutputs = 57,
    NodeRegenCountIncorrect = 58,
    FormatVersionUnsupported = 60,
    FormatVersionNotSet = 61,
    NodeVersionMismatch = 62,
    Allocation = 70,
    InternalBufferTooSmall = 71,
    IntegerOverflow = 72,
    StreamWrongInit = 73,
    StreamTypeIncorrect = 74,
    StreamCapacityTooSmall = 75,
    StreamParameterInvalid = 76,
    LogicError = 80,
    TemporaryLibraryLimitation = 81,
}

// Compression parameters
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum ZL_CParam {
    StickyParameters = 1,
    CompressionLevel = 2,
    DecompressionLevel = 3,
    FormatVersion = 4,
    PermissiveCompression = 5,
    CompressedChecksum = 6,
    ContentChecksum = 7,
    MinStreamSize = 11,
}

// OpenZL C FFI bindings - only the functions available as symbols
#[link(name = "openzl", kind = "static")]
extern "C" {
    // Compression context
    fn ZL_CCtx_create() -> *mut ZL_CCtx;
    fn ZL_CCtx_free(cctx: *mut ZL_CCtx);
    
    // Compression parameter setting
    fn ZL_CCtx_setParameter(cctx: *mut ZL_CCtx, gcparam: ZL_CParam, value: c_int) -> ZL_Report;

    fn ZL_CCtx_compress(
        cctx: *mut ZL_CCtx,
        dst: *mut c_void,
        dst_capacity: usize,
        src: *const c_void,
        src_size: usize,
    ) -> ZL_Report;

    // Simple decompression API
    fn ZL_decompress(
        dst: *mut c_void,
        dst_capacity: usize,
        src: *const c_void,
        src_size: usize,
    ) -> ZL_Report;

    fn ZL_getDecompressedSize(compressed: *const c_void, c_size: usize) -> ZL_Report;

    // Version functions
    fn ZL_getDefaultEncodingVersion() -> c_uint;

    // Error handling functions
    fn ZL_ErrorCode_toString(code: c_int) -> *const c_char;
}

// Inline function implementations (since they're not available as symbols)
fn ZL_compressBound(total_src_size: usize) -> usize {
    // From zl_compress.h: #define ZL_COMPRESSBOUND(s) (((s) * 2) + 512 + 8)
    (total_src_size * 2) + 512 + 8
}

// Helper to safely convert C string to Rust String
unsafe fn c_string_to_rust(c_str: *const c_char) -> String {
    if c_str.is_null() {
        return "Unknown error".to_string();
    }
    match std::ffi::CStr::from_ptr(c_str).to_str() {
        Ok(s) => s.to_string(),
        Err(_) => "Invalid UTF-8 in error message".to_string(),
    }
}

// Helper to check if ZL_Report is an error
fn is_error(report: &ZL_Report) -> bool {
    report._code != ZL_ErrorCode::NoError
}

// Helper to get the value from a successful ZL_Report
fn get_value(report: &ZL_Report) -> usize {
    if is_error(report) {
        0
    } else {
        report._value
    }
}

// Helper to get error message from ZL_Report
fn get_error_message(report: &ZL_Report) -> String {
    if !is_error(report) {
        return "No error".to_string();
    }
    
    unsafe {
        let error_code = report._code as c_int;
        let name_ptr = ZL_ErrorCode_toString(error_code);
        c_string_to_rust(name_ptr)
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

        // Set the format version parameter (required by OpenZL)
        let format_version = unsafe { ZL_getDefaultEncodingVersion() } as c_int;
        println!("Debug: Setting format version to {}", format_version);
        let result = unsafe {
            ZL_CCtx_setParameter(cctx, ZL_CParam::FormatVersion, format_version)
        };
        if is_error(&result) {
            unsafe { ZL_CCtx_free(cctx) };
            return Err(CompressionError::CompressionFailed(format!(
                "Failed to set format version parameter: {}",
                get_error_message(&result)
            )));
        }

        // Verify the parameter was set by trying to compress a small test buffer
        let test_data = b"test";
        let mut test_output = vec![0u8; 100];
        let test_result = unsafe {
            ZL_CCtx_compress(
                cctx,
                test_output.as_mut_ptr() as *mut c_void,
                test_output.len(),
                test_data.as_ptr() as *const c_void,
                test_data.len(),
            )
        };
        
        if is_error(&test_result) {
            unsafe { ZL_CCtx_free(cctx) };
            return Err(CompressionError::CompressionFailed(format!(
                "OpenZL context validation failed: {}",
                get_error_message(&test_result)
            )));
        }

        println!("Debug: OpenZL adapter created successfully");
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

        // Ensure format version is set before compression (OpenZL contexts may need reconfiguration)
        let format_version = unsafe { ZL_getDefaultEncodingVersion() } as c_int;
        let param_result = unsafe {
            ZL_CCtx_setParameter(self.cctx, ZL_CParam::FormatVersion, format_version)
        };
        if is_error(&param_result) {
            return Err(CompressionError::CompressionFailed(format!(
                "Failed to set format version parameter before compression: {}",
                get_error_message(&param_result)
            )));
        }

        println!("Debug: Starting compression of {} bytes", data.len());
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

        if is_error(&result) {
            let error_msg = get_error_message(&result);
            println!("Debug: Compression failed: {}", error_msg);
            return Err(CompressionError::CompressionFailed(error_msg));
        }

        let compressed_size = get_value(&result);
        dst.truncate(compressed_size);
        println!("Debug: Compression successful: {} -> {} bytes", src_size, compressed_size);
        Ok(dst)
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        if data.is_empty() {
            return Ok(Vec::new());
        }

        // Get decompressed size first
        let decompressed_size_report = unsafe {
            ZL_getDecompressedSize(data.as_ptr() as *const c_void, data.len())
        };

        if is_error(&decompressed_size_report) {
            return Err(CompressionError::DecompressionFailed(format!(
                "Failed to get decompressed size: {}",
                get_error_message(&decompressed_size_report)
            )));
        }

        let decompressed_size = get_value(&decompressed_size_report);
        let mut dst = vec![0u8; decompressed_size];

        let result = unsafe {
            ZL_decompress(
                dst.as_mut_ptr() as *mut c_void,
                decompressed_size,
                data.as_ptr() as *const c_void,
                data.len(),
            )
        };

        if is_error(&result) {
            return Err(CompressionError::DecompressionFailed(get_error_message(
                &result,
            )));
        }

        let actual_size = get_value(&result);
        if actual_size != decompressed_size {
            dst.truncate(actual_size);
        }

        Ok(dst)
    }

    fn name(&self) -> &str {
        "openzl"
    }

    fn compress_bound(&self, src_size: usize) -> usize {
        ZL_compressBound(src_size)
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
