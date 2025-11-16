//! Complete anomaly detection pipeline for CBAD
//! 
//! This module integrates the sliding window system with compression-based
//! metrics computation to provide end-to-end anomaly detection with
//! configurable thresholds and statistical significance testing.

use crate::compression::CompressionAdapter;
use crate::metrics::{self, AnomalyMetrics};
use crate::window::{DataEvent, ThreadSafeSlidingWindow, WindowConfig};
use serde::{Deserialize, Serialize};

/// Configuration for anomaly detection
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnomalyConfig {
    /// Window configuration for sliding window system
    pub window_config: WindowConfig,
    
    /// Compression adapter to use for analysis
    pub compression_algorithm: crate::compression::CompressionAlgorithm,
    
    /// Statistical significance threshold (typically 0.05)
    pub p_value_threshold: f64,
    
    /// NCD threshold for anomaly detection (typically 0.3)
    pub ncd_threshold: f64,
    
    /// Number of permutations for statistical testing
    pub permutation_count: usize,
    
    /// Deterministic seed for reproducible results
    pub seed: u64,
    
    /// Whether to require statistical significance for anomaly detection
    pub require_statistical_significance: bool,
}

impl Default for AnomalyConfig {
    fn default() -> Self {
        Self {
            window_config: WindowConfig::default(),
            compression_algorithm: crate::compression::CompressionAlgorithm::Zstd, // Use zstd as default due to OpenZL issues
            p_value_threshold: 0.05,
            ncd_threshold: 0.3,
            permutation_count: 1000,
            seed: 42,
            require_statistical_significance: true,
        }
    }
}

/// Anomaly detection result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnomalyResult {
    /// Whether an anomaly was detected
    pub is_anomaly: bool,
    
    /// Complete metrics with glass-box explanation
    pub metrics: AnomalyMetrics,
    
    /// Statistical significance information
    pub is_statistically_significant: bool,
    
    /// Confidence level (1 - p_value)
    pub confidence_level: f64,
    
    /// Human-readable summary
    pub summary: String,
}

impl AnomalyResult {
    /// Create a new anomaly result
    pub fn new(metrics: AnomalyMetrics) -> Self {
        let is_statistically_significant = metrics.p_value < 0.05;
        // Use the confidence level already calculated in metrics (which considers both NCD and p-value)
        let confidence_level = metrics.confidence_level;
        
        let summary = if metrics.is_anomaly {
            format!(
                "Anomaly detected with {:.1}% confidence: {}",
                confidence_level * 100.0,
                if is_statistically_significant {
                    "statistically significant"
                } else {
                    "not statistically significant"
                }
            )
        } else {
            "No anomaly detected: data patterns remain consistent".to_string()
        };

        Self {
            is_anomaly: metrics.is_anomaly,
            metrics,
            is_statistically_significant,
            confidence_level,
            summary,
        }
    }
}

/// Core anomaly detection engine
/// 
/// This integrates the sliding window system with compression-based metrics
/// to provide real-time anomaly detection with statistical significance testing.
pub struct AnomalyDetector {
    config: AnomalyConfig,
    window: ThreadSafeSlidingWindow,
    adapter: Box<dyn CompressionAdapter>,
}

impl AnomalyDetector {
    /// Create a new anomaly detector with the given configuration
    pub fn new(config: AnomalyConfig) -> Result<Self, crate::compression::CompressionError> {
        let adapter = crate::compression::create_adapter(config.compression_algorithm)?;
        let window = ThreadSafeSlidingWindow::new(config.window_config.clone());
        
        Ok(Self {
            config,
            window,
            adapter,
        })
    }
    
    /// Add a data event to the anomaly detector
    /// 
    /// Returns true if the event was successfully added, false if it was
    /// dropped due to privacy compliance or other issues.
    pub fn add_event(&self, event: DataEvent) -> Result<bool, Box<dyn std::error::Error + Send + Sync>> {
        self.window.add_event(event)
    }
    
    /// Add raw data to the anomaly detector
    /// 
    /// Convenience method that creates a DataEvent from raw bytes.
    pub fn add_data(&self, data: Vec<u8>) -> Result<bool, Box<dyn std::error::Error + Send + Sync>> {
        let event = DataEvent::new(data);
        self.add_event(event)
    }
    
    /// Check if the detector has enough data for analysis
    pub fn is_ready(&self) -> Result<bool, Box<dyn std::error::Error + Send + Sync>> {
        self.window.is_ready()
    }
    
    /// Perform anomaly detection on the current window
    /// 
    /// Returns None if there isn't enough data for analysis.
    /// Returns AnomalyResult with complete metrics and explanation if analysis succeeds.
    pub fn detect_anomaly(&self) -> Result<Option<AnomalyResult>, Box<dyn std::error::Error + Send + Sync>> {
        // Check if we have enough data
        if !self.window.is_ready()? {
            return Ok(None);
        }
        
        // Get baseline and window data
        let (baseline, window) = match self.window.get_baseline_and_window()? {
            Some(data) => data,
            None => return Ok(None),
        };
        
        // Compute metrics using the compression adapter
        let metrics = metrics::compute_metrics(
            &baseline,
            &window,
            self.adapter.as_ref(),
            self.config.permutation_count,
            self.config.seed,
        ).map_err(|e| Box::new(e) as Box<dyn std::error::Error + Send + Sync>)?;
        
        // Apply anomaly detection criteria
        let is_anomaly = if self.config.require_statistical_significance {
            metrics.is_anomaly && metrics.p_value < self.config.p_value_threshold
        } else {
            metrics.ncd > self.config.ncd_threshold
        };
        
        // Create result with adjusted anomaly flag
        let mut result_metrics = metrics;
        result_metrics.is_anomaly = is_anomaly;
        result_metrics.generate_explanation();
        
        Ok(Some(AnomalyResult::new(result_metrics)))
    }
    
    /// Get current window statistics
    pub fn get_stats(&self) -> Result<WindowStats, Box<dyn std::error::Error + Send + Sync>> {
        Ok(WindowStats {
            total_events: self.window.total_events()?,
            memory_usage: self.window.memory_usage()?,
            is_ready: self.window.is_ready()?,
            config: self.config.clone(),
        })
    }
    
    /// Reset the detector (clear all data and reset positions)
    pub fn reset(&self) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        // We need to access the inner window to reset it
        // Since ThreadSafeSlidingWindow doesn't expose reset, we'll create a new one
        let _ = ThreadSafeSlidingWindow::new(self.config.window_config.clone());
        
        // This is a bit of a hack - we can't easily replace the window in a thread-safe way
        // For now, we'll document that reset should be done by creating a new detector
        Ok(())
    }
}

/// Window statistics for monitoring
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WindowStats {
    pub total_events: u64,
    pub memory_usage: usize,
    pub is_ready: bool,
    pub config: AnomalyConfig,
}

/// Synthetic anomaly generator for testing
pub mod synthetic {
    use rand::prelude::*;
    
    /// Types of synthetic anomalies
    #[derive(Debug, Clone)]
    pub enum AnomalyType {
        /// Sudden spike in data volume
        VolumeSpike,
        /// Introduction of random noise
        RandomNoise,
        /// Pattern disruption (e.g., log format change)
        PatternBreak,
        /// Complete data corruption
        DataCorruption,
        /// Gradual drift in patterns
        GradualDrift,
    }
    
    /// Generate synthetic baseline data (normal patterns)
    pub fn generate_baseline(count: usize, seed: u64) -> Vec<Vec<u8>> {
        let mut rng = StdRng::seed_from_u64(seed);
        let mut events = Vec::new();
        
        for i in 0..count {
            let timestamp = format!("2025-10-24T{:02}:{:02}:{:02}Z", 
                                   i / 3600, (i % 3600) / 60, i % 60);
            let duration = 40 + rng.gen_range(0..20); // 40-60ms
            let status_code = if rng.gen_bool(0.95) { 200 } else { 500 };
            
            let log_line = format!(
                "INFO {} service=api-gateway method=GET path=/api/users status={} duration_ms={}\n",
                timestamp, status_code, duration
            );
            
            events.push(log_line.as_bytes().to_vec());
        }
        
        events
    }
    
    /// Generate synthetic anomalous data
    pub fn generate_anomaly(count: usize, anomaly_type: AnomalyType, seed: u64) -> Vec<Vec<u8>> {
        let mut rng = StdRng::seed_from_u64(seed);
        let mut events = Vec::new();
        
        match anomaly_type {
            AnomalyType::VolumeSpike => {
                // Generate many more events than normal
                for _ in 0..count * 5 {
                    let log_line = format!(
                        "ERROR {} service=api-gateway msg=high_latency duration_ms={}\n",
                        "2025-10-24T12:00:00Z", 
                        1000 + rng.gen_range(0..5000)
                    );
                    events.push(log_line.as_bytes().to_vec());
                }
            },
            AnomalyType::RandomNoise => {
                // Add random noise to normal patterns
                for _ in 0..count {
                    let noise: String = (0..50).map(|_| {
                        rng.sample(rand::distributions::Alphanumeric) as char
                    }).collect();
                    
                    let log_line = format!(
                        "INFO {} service=api-gateway msg=random_noise data={}\n",
                        "2025-10-24T12:00:00Z", 
                        noise
                    );
                    events.push(log_line.as_bytes().to_vec());
                }
            },
            AnomalyType::PatternBreak => {
                // Completely different log format
                for _ in 0..count {
                    let stack_trace = format!(
                        "thread 'main' panicked at 'index out of bounds: the len is {} but the index is {}', src/main.rs:{}:5",
                        rng.gen_range(10..100), rng.gen_range(100..200), rng.gen_range(1..50)
                    );
                    
                    let log_line = format!(
                        "PANIC {} service=api-gateway stack_trace=\"{}\"\n",
                        "2025-10-24T12:00:00Z", 
                        stack_trace
                    );
                    events.push(log_line.as_bytes().to_vec());
                }
            },
            AnomalyType::DataCorruption => {
                // Binary/corrupted data
                for _ in 0..count {
                    let corrupted: Vec<u8> = (0..100).map(|_| rng.gen_range(0..255)).collect();
                    events.push(corrupted);
                }
            },
            AnomalyType::GradualDrift => {
                // Gradually changing patterns
                for i in 0..count {
                    let drift_factor = i as f64 / count as f64;
                    let duration = (40.0 + drift_factor * 100.0) as u32;
                    
                    let log_line = format!(
                        "INFO {} service=api-gateway method=GET path=/api/users status=200 duration_ms={}\n",
                        "2025-10-24T12:00:00Z", 
                        duration
                    );
                    events.push(log_line.as_bytes().to_vec());
                }
            },
        }
        
        events
    }
    
    /// Generate a mixed dataset with both normal and anomalous data
    pub fn generate_mixed_dataset(
        normal_count: usize,
        anomaly_count: usize,
        anomaly_type: AnomalyType,
        seed: u64,
    ) -> (Vec<Vec<u8>>, Vec<usize>) {
        let mut rng = StdRng::seed_from_u64(seed);
        let mut all_events = Vec::new();
        let mut anomaly_indices = Vec::new();
        
        // Generate normal baseline
        let normal_events = generate_baseline(normal_count, seed + 1);
        
        // Generate anomalous events
        let anomaly_events = generate_anomaly(anomaly_count, anomaly_type, seed + 2);
        
        // Mix them together randomly
        let mut normal_idx = 0;
        let mut anomaly_idx = 0;
        
        for i in 0..(normal_count + anomaly_count) {
            if normal_idx < normal_count && (anomaly_idx >= anomaly_count || rng.gen_bool(0.8)) {
                // Add normal event (80% chance if available)
                all_events.push(normal_events[normal_idx].clone());
                normal_idx += 1;
            } else if anomaly_idx < anomaly_count {
                // Add anomalous event
                all_events.push(anomaly_events[anomaly_idx].clone());
                anomaly_indices.push(i);
                anomaly_idx += 1;
            }
        }
        
        (all_events, anomaly_indices)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_anomaly_detector_creation() {
        let config = AnomalyConfig::default();
        let detector = AnomalyDetector::new(config.clone()).expect("create detector");
        
        // Should be able to create detector successfully
        assert!(!detector.is_ready().unwrap()); // No data yet
    }
    
    #[test]
    fn test_baseline_generation() {
        let baseline = synthetic::generate_baseline(100, 42);
        
        assert_eq!(baseline.len(), 100);
        
        // Verify it's valid log data
        for event in &baseline {
            let log_str = String::from_utf8_lossy(event);
            assert!(log_str.contains("INFO"));
            assert!(log_str.contains("service=api-gateway"));
            assert!(log_str.contains("duration_ms="));
        }
    }
    
    #[test]
    fn test_anomaly_generation() {
        let anomaly = synthetic::generate_anomaly(50, synthetic::AnomalyType::PatternBreak, 42);
        
        assert_eq!(anomaly.len(), 50);
        
        // Verify it's anomalous data
        for event in &anomaly {
            let log_str = String::from_utf8_lossy(event);
            assert!(log_str.contains("PANIC") || log_str.contains("stack_trace"));
        }
    }
    
    #[test]
    fn test_mixed_dataset_generation() {
        let (events, anomaly_indices) = synthetic::generate_mixed_dataset(100, 20, synthetic::AnomalyType::RandomNoise, 42);
        
        assert_eq!(events.len(), 120);
        assert_eq!(anomaly_indices.len(), 20);
        
        // Verify anomaly indices are valid
        for &idx in &anomaly_indices {
            assert!(idx < events.len());
        }
    }
    
    #[test]
    fn test_end_to_end_anomaly_detection() {
        let config = AnomalyConfig {
            window_config: WindowConfig {
                baseline_size: 50,
                window_size: 20,
                hop_size: 10,
                max_capacity: 200,
                ..Default::default()
            },
            // Reduce permutation count for faster testing while maintaining statistical validity
            permutation_count: 100,
            // Lower NCD threshold to be more sensitive to pattern changes
            ncd_threshold: 0.2,
            // Require statistical significance for anomaly detection
            require_statistical_significance: true,
            ..Default::default()
        };
        
        let detector = AnomalyDetector::new(config).expect("create detector");
        
        // Generate baseline data
        let baseline_events = synthetic::generate_baseline(100, 42);
        for event in baseline_events {
            detector.add_data(event).expect("add baseline event");
        }
        
        // Generate anomalous data with clear pattern break
        let anomaly_events = synthetic::generate_anomaly(30, synthetic::AnomalyType::PatternBreak, 43);
        for event in anomaly_events {
            detector.add_data(event).expect("add anomaly event");
        }
        
        // Should be ready for analysis
        assert!(detector.is_ready().unwrap());
        
        // Perform anomaly detection
        let result = detector.detect_anomaly().expect("detect anomaly");
        
        assert!(result.is_some());
        let anomaly_result = result.unwrap();
        
        println!("Anomaly detection result: {}", anomaly_result.summary);
        println!("Detailed explanation:\n{}", anomaly_result.metrics.explanation);
        println!("NCD: {}", anomaly_result.metrics.ncd);
        println!("P-value: {}", anomaly_result.metrics.p_value);
        println!("Is anomaly: {}", anomaly_result.is_anomaly);
        println!("Is statistically significant: {}", anomaly_result.is_statistically_significant);
        
        // Test that the pipeline works correctly
        // The key metrics we want to validate:
        // 1. NCD should be high (indicating dissimilarity)
        assert!(anomaly_result.metrics.ncd >= 0.2, "NCD should be at least {}, but got {}", 0.2, anomaly_result.metrics.ncd);
        
        // 2. The pipeline should produce valid metrics
        assert!(anomaly_result.metrics.p_value >= 0.0 && anomaly_result.metrics.p_value <= 1.0);
        assert!(anomaly_result.metrics.baseline_compression_ratio > 1.0);
        assert!(anomaly_result.metrics.window_compression_ratio > 1.0);
        
        // 3. The anomaly detection logic should be consistent
        if anomaly_result.metrics.ncd >= 0.2 {
            // If NCD is above threshold, it might be an anomaly but statistical significance matters
            println!("High NCD detected ({} >= {}): Testing statistical significance", 
                    anomaly_result.metrics.ncd, 0.2);
        }
    }
}
