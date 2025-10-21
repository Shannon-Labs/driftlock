use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Metrics {
    pub entropy: f64,
    pub compression_ratio: f64,
    pub ncd: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComputeConfig {
    pub window_size: usize,
    pub hop_size: usize,
    pub threshold: f64,
    pub deterministic_seed: u64,
}

// Placeholder implementation. Replace with real compression/entropy/NCD.
pub fn compute_metrics(_data: &[u8], _cfg: &ComputeConfig) -> Metrics {
    Metrics { entropy: 0.0, compression_ratio: 1.0, ncd: 0.0 }
}

// Optional C FFI stubs
#[no_mangle]
pub extern "C" fn cbad_compute_metrics_len(_: *const u8, len: usize) -> f64 { len as f64 }

