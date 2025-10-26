//! Shannon entropy calculator for CBAD
//! 
//! Shannon entropy measures the randomness or unpredictability of data.
//! High entropy indicates random/unstructured data, while low entropy
//! indicates predictable/structured data.
//! 
//! Entropy Formula: H(X) = -Σ p(x) * log2(p(x))
//! Where p(x) is the probability of symbol x appearing in the data.
//! 
//! For byte data, entropy ranges from 0.0 (completely predictable) 
//! to 8.0 (completely random - all bytes equally likely).

use crate::metrics::Result;

/// Entropy calculation result with detailed metrics
#[derive(Debug, Clone)]
pub struct EntropyMetrics {
    /// Shannon entropy in bits per byte (0.0 to 8.0)
    pub entropy: f64,
    /// Most frequent byte value
    pub most_frequent_byte: u8,
    /// Frequency of the most frequent byte (0.0 to 1.0)
    pub max_frequency: f64,
    /// Number of unique byte values observed
    pub unique_bytes: usize,
    /// Total number of bytes analyzed
    pub total_bytes: usize,
    /// Byte frequency distribution
    pub byte_frequencies: [u64; 256],
}

impl EntropyMetrics {
    /// Create new entropy metrics from frequency distribution
    pub fn new(byte_frequencies: [u64; 256], total_bytes: usize) -> Self {
        if total_bytes == 0 {
            return Self {
                entropy: 0.0,
                most_frequent_byte: 0,
                max_frequency: 0.0,
                unique_bytes: 0,
                total_bytes: 0,
                byte_frequencies: [0; 256],
            };
        }

        // Calculate entropy using Shannon's formula
        let mut entropy = 0.0;
        let mut max_freq = 0u64;
        let mut most_frequent = 0u8;
        let mut unique_count = 0usize;

        for (byte, &count) in byte_frequencies.iter().enumerate() {
            if count > 0 {
                unique_count += 1;
                let probability = count as f64 / total_bytes as f64;
                entropy -= probability * probability.log2();
                
                if count > max_freq {
                    max_freq = count;
                    most_frequent = byte as u8;
                }
            }
        }

        let max_frequency = max_freq as f64 / total_bytes as f64;

        Self {
            entropy,
            most_frequent_byte: most_frequent,
            max_frequency,
            unique_bytes: unique_count,
            total_bytes,
            byte_frequencies,
        }
    }

    /// Get human-readable interpretation of entropy value
    pub fn interpretation(&self) -> &'static str {
        if self.entropy < 2.0 {
            "Very low entropy - highly structured/predictable data"
        } else if self.entropy < 4.0 {
            "Low entropy - structured data with patterns"
        } else if self.entropy < 6.0 {
            "Moderate entropy - semi-structured data"
        } else if self.entropy < 7.0 {
            "High entropy - mostly random data"
        } else {
            "Very high entropy - highly random/unstructured data"
        }
    }

    /// Check if entropy indicates an anomaly based on threshold
    pub fn is_anomaly(&self, threshold: f64) -> bool {
        self.entropy >= threshold
    }

    /// Calculate entropy change from baseline
    pub fn entropy_change(&self, baseline_entropy: f64) -> f64 {
        if baseline_entropy > 0.0 {
            (self.entropy - baseline_entropy) / baseline_entropy
        } else {
            0.0
        }
    }
}

/// Compute Shannon entropy for a byte sequence
/// 
/// Returns entropy in bits per byte (0.0 to 8.0)
pub fn compute_entropy(data: &[u8]) -> f64 {
    let frequencies = calculate_byte_frequencies(data);
    let metrics = EntropyMetrics::new(frequencies, data.len());
    metrics.entropy
}

/// Calculate byte frequency distribution
/// 
/// Returns array where index i contains count of byte value i
pub fn calculate_byte_frequencies(data: &[u8]) -> [u64; 256] {
    let mut frequencies = [0u64; 256];
    
    for &byte in data {
        frequencies[byte as usize] += 1;
    }
    
    frequencies
}

/// Compare entropy between baseline and window
/// 
/// Returns both individual entropies and the change between them
pub fn compare_entropy(
    baseline: &[u8],
    window: &[u8],
) -> Result<(f64, f64, f64)> {
    let baseline_entropy = compute_entropy(baseline);
    let window_entropy = compute_entropy(window);
    
    let entropy_change = if baseline_entropy > 0.0 {
        (window_entropy - baseline_entropy) / baseline_entropy
    } else {
        0.0
    };
    
    Ok((baseline_entropy, window_entropy, entropy_change))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_entropy_uniform_distribution() {
        // Uniform distribution should have maximum entropy (8.0 bits/byte)
        let mut uniform_data = Vec::new();
        for i in 0..256 {
            uniform_data.push(i as u8);
        }
        // Repeat to get statistical significance
        uniform_data = uniform_data.repeat(100);

        let entropy = compute_entropy(&uniform_data);
        println!("Uniform distribution entropy: {:.2} bits/byte", entropy);
        
        // Should be very close to 8.0 for uniform distribution
        assert!((entropy - 8.0).abs() < 0.1, "Uniform distribution should have entropy ≈ 8.0");
    }

    #[test]
    fn test_entropy_single_byte() {
        // Single repeated byte should have minimum entropy (0.0 bits/byte)
        let single_byte = b"A".repeat(1000);

        let entropy = compute_entropy(&single_byte);
        println!("Single byte entropy: {:.2} bits/byte", entropy);
        
        // Should be very close to 0.0 for single repeated byte
        assert!(entropy < 0.01, "Single repeated byte should have entropy ≈ 0.0");
    }

    #[test]
    fn test_entropy_empty_data() {
        let empty_data = b"";
        
        let entropy = compute_entropy(empty_data);
        assert_eq!(entropy, 0.0, "Empty data should have entropy 0.0");
        
        let metrics = EntropyMetrics::new([0; 256], 0);
        assert_eq!(metrics.total_bytes, 0);
        assert_eq!(metrics.entropy, 0.0);
    }

    #[test]
    fn test_entropy_otlp_logs() {
        // Test with realistic OTLP log data
        let otlp_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#;
        
        let entropy = compute_entropy(otlp_log.as_bytes());
        println!("OTLP log entropy: {:.2} bits/byte", entropy);
        
        // OTLP JSON should have moderate entropy due to structured format
        assert!(entropy > 2.0 && entropy < 6.0, "OTLP JSON should have moderate entropy");
    }

    #[test]
    fn test_entropy_comparison() {
        let baseline_log = r#"{"timestamp":"2025-10-24T00:00:00Z","severity":"INFO","service":"api-gateway","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#;
        let high_entropy_log = r#"{"timestamp":"2025-10-24T00:00:01Z","severity":"ERROR","service":"api-gateway","message":"Panic occurred","attributes":{"stack_trace":"0x3fa8d1b2c9e47f56::panic::trace::[ns=923847923847923847923847]::random_payload=9fjK2L1pQwZ8xT4rB7nC6Mv0HdYG5s2tR1uQ3w8yAaEeIiOo","binary_blob":"Q29tcHJlc3NlZEJsb2I6ZGV0ZXJtaW5pc3RpY1Nob3J0c0FuZFJhbmRvbVVuaWNvZGVEYXRh"}}"#;

        let baseline = baseline_log.as_bytes().repeat(100);
        let window = high_entropy_log.as_bytes().repeat(20);

        let (baseline_entropy, window_entropy, entropy_change) = compare_entropy(&baseline, &window)
            .expect("compare entropy");

        println!("Baseline entropy: {:.2} bits/byte", baseline_entropy);
        println!("Window entropy: {:.2} bits/byte", window_entropy);
        println!("Entropy change: {:.1}%", entropy_change * 100.0);

        // Anomalous window (with stack trace) should have higher entropy
        assert!(window_entropy > baseline_entropy, "Anomalous data should have higher entropy");
    }
}
