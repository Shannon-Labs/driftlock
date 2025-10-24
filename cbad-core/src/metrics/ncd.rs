//! Normalized Compression Distance (NCD) calculator for CBAD
//! 
//! NCD is the primary anomaly detection metric in Driftlock. It measures
//! the dissimilarity between two data sequences using compression.
//! 
//! The key insight: if two sequences are similar, compressing them together
//! won't be much larger than compressing the larger one individually.
//! If they're different, the combined compression will be significantly larger.
//! 
//! NCD Formula: NCD(x,y) = (C(xy) - min(C(x), C(y))) / max(C(x), C(y))
//! Where C(x) is the compressed size of sequence x, and C(xy) is the
//! compressed size of the concatenation of x and y.

use crate::compression::CompressionAdapter;
use crate::metrics::{MetricsError, Result};

/// NCD calculation result with detailed metrics
#[derive(Debug, Clone)]
pub struct NCDMetrics {
    /// Normalized Compression Distance (0.0 to 1.0)
    /// 0.0 = identical patterns, 1.0 = completely different
    pub ncd: f64,
    /// Compressed size of baseline data
    pub baseline_compressed_size: usize,
    /// Compressed size of window data
    pub window_compressed_size: usize,
    /// Compressed size of concatenated baseline + window
    pub combined_compressed_size: usize,
    /// Which sequence was larger after compression
    pub larger_sequence: SequenceType,
}

#[derive(Debug, Clone, PartialEq)]
pub enum SequenceType {
    Baseline,
    Window,
    Equal,
}

impl NCDMetrics {
    /// Create new NCD metrics from compression sizes
    pub fn new(
        baseline_compressed: usize,
        window_compressed: usize,
        combined_compressed: usize,
    ) -> Self {
        let min_compressed = baseline_compressed.min(window_compressed) as f64;
        let max_compressed = baseline_compressed.max(window_compressed) as f64;
        
        // NCD formula: (C(xy) - min(C(x), C(y))) / max(C(x), C(y))
        let ncd = if max_compressed > 0.0 {
            ((combined_compressed as f64) - min_compressed) / max_compressed
        } else {
            0.0
        };

        // Clamp to valid range [0.0, 1.0] to handle floating point errors
        let ncd = ncd.clamp(0.0, 1.0);

        let larger_sequence = if baseline_compressed > window_compressed {
            SequenceType::Baseline
        } else if window_compressed > baseline_compressed {
            SequenceType::Window
        } else {
            SequenceType::Equal
        };

        Self {
            ncd,
            baseline_compressed_size: baseline_compressed,
            window_compressed_size: window_compressed,
            combined_compressed_size: combined_compressed,
            larger_sequence,
        }
    }

    /// Get human-readable interpretation of NCD score
    pub fn interpretation(&self) -> &'static str {
        if self.ncd < 0.3 {
            "Low dissimilarity - sequences are very similar"
        } else if self.ncd < 0.5 {
            "Moderate dissimilarity - sequences share some patterns"
        } else if self.ncd < 0.7 {
            "High dissimilarity - sequences are quite different"
        } else {
            "Very high dissimilarity - sequences are completely different"
        }
    }

    /// Check if NCD indicates an anomaly based on threshold
    pub fn is_anomaly(&self, threshold: f64) -> bool {
        self.ncd >= threshold
    }

    /// Get compression efficiency of concatenation
    /// Returns how much larger the combined data is vs the larger individual sequence
    pub fn concatenation_overhead(&self) -> f64 {
        let min_individual = self.baseline_compressed_size.min(self.window_compressed_size) as f64;
        if min_individual > 0.0 {
            (self.combined_compressed_size as f64 - min_individual) / min_individual
        } else {
            0.0
        }
    }
}

/// Compute Normalized Compression Distance between baseline and window
/// 
/// This is the core anomaly detection algorithm in Driftlock.
/// High NCD values indicate that the window contains patterns
/// that are significantly different from the baseline.
pub fn compute_ncd(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<f64> {
    let metrics = compute_ncd_detailed(baseline, window, adapter)?;
    Ok(metrics.ncd)
}

/// Compute NCD with detailed metrics for comprehensive analysis
pub fn compute_ncd_detailed(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
) -> Result<NCDMetrics> {
    if baseline.is_empty() || window.is_empty() {
        return Err(MetricsError::InvalidInput(
            "Baseline and window data must not be empty".to_string()
        ));
    }

    // Compress baseline individually
    let baseline_compressed = adapter.compress(baseline)
        .map_err(|e| MetricsError::CompressionFailed(format!("Baseline compression failed: {}", e)))?;

    // Compress window individually
    let window_compressed = adapter.compress(window)
        .map_err(|e| MetricsError::CompressionFailed(format!("Window compression failed: {}", e)))?;

    // Compress concatenated data
    let mut combined = Vec::with_capacity(baseline.len() + window.len());
    combined.extend_from_slice(baseline);
    combined.extend_from_slice(window);
    
    let combined_compressed = adapter.compress(&combined)
        .map_err(|e| MetricsError::CompressionFailed(format!("Combined compression failed: {}", e)))?;

    Ok(NCDMetrics::new(
        baseline_compressed.len(),
        window_compressed.len(),
        combined_compressed.len(),
    ))
}

/// Compute NCD for multiple sequences (useful for clustering analysis)
/// 
/// Returns a matrix of NCD values between all pairs of sequences
pub fn compute_ncd_matrix(
    sequences: &[&[u8]],
    adapter: &dyn CompressionAdapter,
) -> Result<Vec<Vec<f64>>> {
    let n = sequences.len();
    let mut matrix = vec![vec![0.0; n]; n];

    for i in 0..n {
        for j in i..n {
            if i == j {
                matrix[i][j] = 0.0; // NCD of sequence with itself is 0
            } else {
                let ncd = compute_ncd(sequences[i], sequences[j], adapter)?;
                matrix[i][j] = ncd;
                matrix[j][i] = ncd; // NCD is symmetric
            }
        }
    }

    Ok(matrix)
}

/// Quick NCD check for anomaly detection
/// 
/// Returns true if NCD exceeds threshold, false otherwise
/// More efficient than computing full detailed metrics when only
/// the anomaly decision is needed
pub fn is_anomaly_ncd(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    threshold: f64,
) -> Result<bool> {
    let ncd = compute_ncd(baseline, window, adapter)?;
    Ok(ncd >= threshold)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::compression::{create_adapter, CompressionAlgorithm};

    fn create_test_adapter() -> Box<dyn CompressionAdapter> {
        // Use OpenZL for testing - it's the primary compression algorithm
        #[cfg(feature = "openzl")]
        {
            create_adapter(CompressionAlgorithm::OpenZL).expect("OpenZL adapter")
        }
        #[cfg(not(feature = "openzl"))]
        {
            create_adapter(CompressionAlgorithm::Zstd).expect("Zstd adapter")
        }
    }

    #[test]
    fn test_ncd_identical_sequences() {
        let adapter = create_test_adapter();

        let data = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(10);
        
        let metrics = compute_ncd_detailed(
            data.as_slice(),
            data.as_slice(),
            adapter.as_ref(),
        ).expect("compute NCD for identical sequences");

        println!("NCD of identical sequences: {:.4}", metrics.ncd);
        
        // NCD of identical sequences should be very close to 0
        assert!(metrics.ncd < 0.1, "Identical sequences should have NCD ≈ 0");
        assert_eq!(metrics.larger_sequence, SequenceType::Equal);
    }

    #[test]
    fn test_ncd_similar_sequences() {
        let adapter = create_test_adapter();

        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(100);
        let window = b"INFO 2025-10-24T00:00:01Z service=api-gateway msg=request_completed duration_ms=43\n".repeat(10);
        
        let metrics = compute_ncd_detailed(
            baseline.as_slice(),
            window.as_slice(),
            adapter.as_ref(),
        ).expect("compute NCD for similar sequences");

        println!("NCD of similar sequences: {:.4}", metrics.ncd);
        println!("Interpretation: {}", metrics.interpretation());
        
        // Similar sequences should have moderate NCD
        assert!(metrics.ncd < 0.5, "Similar sequences should have moderate NCD");
        assert!(metrics.ncd > 0.0, "Different sequences should have NCD > 0");
    }

    #[test]
    fn test_ncd_dissimilar_sequences() {
        let adapter = create_test_adapter();

        // Regular structured logs
        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed duration_ms=42\n".repeat(100);
        
        // Anomalous unstructured data (stack trace)
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic stack_trace=\"thread 'main' panicked at 'index out of bounds: the len is 10 but the index is 42', /rustc/.../src/libcore/slice/mod.rs:1234:5\"\n".repeat(10);

        let metrics = compute_ncd_detailed(
            baseline.as_slice(),
            window.as_slice(),
            adapter.as_ref(),
        ).expect("compute NCD for dissimilar sequences");

        println!("NCD of dissimilar sequences: {:.4}", metrics.ncd);
        println!("Baseline compressed: {} bytes", metrics.baseline_compressed_size);
        println!("Window compressed: {} bytes", metrics.window_compressed_size);
        println!("Combined compressed: {} bytes", metrics.combined_compressed_size);
        println!("Interpretation: {}", metrics.interpretation());
        println!("Concatenation overhead: {:.1}%", metrics.concatenation_overhead() * 100.0);
        
        // Dissimilar sequences should have high NCD
        assert!(metrics.ncd > 0.5, "Dissimilar sequences should have high NCD");
        
        // Should be detected as anomaly with reasonable threshold
        assert!(metrics.is_anomaly(0.4), "Should be detected as anomaly");
    }

    #[test]
    fn test_ncd_range() {
        let adapter = create_test_adapter();

        // Test that NCD always returns values in valid range [0.0, 1.0]
        let test_cases = vec![
            (b"short".as_slice(), b"very long sequence with lots of data".as_slice()),
            (b"identical".as_slice(), b"identical".as_slice()),
            (b"".as_slice(), b"non-empty".as_slice()), // Edge case: empty vs non-empty
        ];

        for (baseline, window) in test_cases {
            if !baseline.is_empty() && !window.is_empty() {
                let ncd = compute_ncd(baseline, window, adapter.as_ref())
                    .expect("compute NCD");
                
                println!("NCD for test case: {:.4}", ncd);
                assert!(ncd >= 0.0 && ncd <= 1.0, "NCD must be in range [0.0, 1.0]");
            }
        }
    }

    #[test]
    fn test_ncd_symmetry() {
        let adapter = create_test_adapter();

        let seq1 = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed\n".repeat(50);
        let seq2 = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=stack_trace_panic\n".repeat(20);

        let ncd_1_2 = compute_ncd(&seq1, &seq2, adapter.as_ref()).expect("compute NCD 1->2");
        let ncd_2_1 = compute_ncd(&seq2, &seq1, adapter.as_ref()).expect("compute NCD 2->1");

        println!("NCD(1,2): {:.4}, NCD(2,1): {:.4}", ncd_1_2, ncd_2_1);
        
        // NCD should be approximately symmetric (allowing for small floating point differences)
        assert!((ncd_1_2 - ncd_2_1).abs() < 0.01, "NCD should be symmetric");
    }

    #[test]
    fn test_ncd_matrix() {
        let adapter = create_test_adapter();

        let seq1 = b"INFO service=api-gateway msg=request_completed\n".repeat(10);
        let seq2 = b"INFO service=api-gateway msg=request_completed\n".repeat(10); // Identical to first
        let seq3 = b"ERROR service=api-gateway msg=stack_trace\n".repeat(5);      // Different

        let sequences = vec![
            seq1.as_slice(),
            seq2.as_slice(),
            seq3.as_slice(),
        ];

        let matrix = compute_ncd_matrix(&sequences, adapter.as_ref())
            .expect("compute NCD matrix");

        println!("NCD Matrix:");
        for row in &matrix {
            println!("{:?}", row.iter().map(|x| format!("{:.3}", x)).collect::<Vec<_>>());
        }

        // Matrix should be symmetric
        for i in 0..matrix.len() {
            for j in i+1..matrix.len() {
                assert!((matrix[i][j] - matrix[j][i]).abs() < 0.01, 
                        "Matrix should be symmetric");
            }
        }

        // Identical sequences should have NCD ≈ 0
        assert!(matrix[0][1] < 0.1, "Identical sequences should have low NCD");
        
        // Diagonal should be 0 (sequence with itself)
        for i in 0..matrix.len() {
            assert!(matrix[i][i] < 0.01, "Diagonal should be ~0");
        }
    }

    #[test]
    fn test_otlp_log_ncd() {
        let adapter = create_test_adapter();

        // Realistic OTLP log entries
        let baseline_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#;
        
        let similar_log = r#"{"timestamp":"2025-10-24T00:00:01Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":43}}"#;
        
        let anomalous_log = r#"{"timestamp":"2025-10-24T00:00:02Z","severity":"ERROR","service":"api-gateway","message":"Panic occurred","attributes":{"stack_trace":"thread 'main' panicked at 'index out of bounds', src/main.rs:42:13","error":"runtime panic"}}"#;

        let baseline_similar = compute_ncd(
            baseline_log.as_bytes(),
            similar_log.as_bytes(),
            adapter.as_ref(),
        ).expect("compute NCD baseline vs similar");

        let baseline_anomalous = compute_ncd(
            baseline_log.as_bytes(),
            anomalous_log.as_bytes(),
            adapter.as_ref(),
        ).expect("compute NCD baseline vs anomalous");

        println!("Baseline vs Similar: {:.4}", baseline_similar);
        println!("Baseline vs Anomalous: {:.4}", baseline_anomalous);

        // Similar logs should have lower NCD than anomalous logs
        assert!(baseline_similar < baseline_anomalous, 
                "Similar logs should have lower NCD than anomalous logs");
        
        // Anomalous log should be detectable as anomaly
        assert!(baseline_anomalous > 0.4, "Anomalous log should have high NCD");
    }

    #[test]
    fn test_quick_anomaly_check() {
        let adapter = create_test_adapter();

        let baseline = b"INFO 2025-10-24T00:00:00Z service=api-gateway msg=request_completed\n".repeat(50);
        let window = b"ERROR 2025-10-24T00:00:01Z service=api-gateway msg=panic\n".repeat(20);

        // Test with low threshold (should not detect)
        let is_anomaly_low = is_anomaly_ncd(&baseline, &window, adapter.as_ref(), 0.9)
            .expect("quick anomaly check");
        assert!(!is_anomaly_low, "Should not detect with very high threshold");

        // Test with high threshold (should detect)
        let is_anomaly_high = is_anomaly_ncd(&baseline, &window, adapter.as_ref(), 0.3)
            .expect("quick anomaly check");
        assert!(is_anomaly_high, "Should detect with low threshold");
    }
}
