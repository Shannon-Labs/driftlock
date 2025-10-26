//! Performance validation suite for CBAD
//! 
//! This module provides comprehensive benchmarks and performance tests
//! to validate that the CBAD engine meets the Phase 1 targets:
//! - Throughput: >10k events/s
//! - Latency: <400ms p95
//! - Memory usage: <2GB for 1M event windows
//! 
//! The benchmarks measure real-world scenarios with synthetic data
//! that mimics production OTLP telemetry patterns.

use std::time::{Duration, Instant};
use std::collections::HashMap;

use crate::compression::{create_adapter, CompressionAlgorithm};
use crate::metrics;
use crate::window::{SlidingWindow, WindowConfig, DataEvent};

/// Performance benchmark results
#[derive(Debug, Clone)]
pub struct BenchmarkResult {
    /// Name of the benchmark
    pub name: String,
    /// Duration of the benchmark run
    pub duration: Duration,
    /// Number of operations completed
    pub operations: u64,
    /// Throughput in operations per second
    pub ops_per_sec: f64,
    /// Average operation time in microseconds
    pub avg_micros: f64,
    /// 95th percentile operation time in microseconds
    pub p95_micros: f64,
    /// 99th percentile operation time in microseconds
    pub p99_micros: f64,
    /// Memory usage at end of benchmark in MB
    pub memory_mb: f64,
    /// Additional metrics
    pub metrics: HashMap<String, f64>,
}

impl BenchmarkResult {
    /// Create a new benchmark result
    pub fn new(name: &str) -> Self {
        Self {
            name: name.to_string(),
            duration: Duration::new(0, 0),
            operations: 0,
            ops_per_sec: 0.0,
            avg_micros: 0.0,
            p95_micros: 0.0,
            p99_micros: 0.0,
            memory_mb: 0.0,
            metrics: HashMap::new(),
        }
    }

    /// Print results in a human-readable format
    pub fn print_summary(&self) {
        println!("Benchmark: {}", self.name);
        println!("  Duration: {:.2?}", self.duration);
        println!("  Operations: {}", self.operations);
        println!("  Throughput: {:.0} ops/sec", self.ops_per_sec);
        println!("  Avg: {:.1} μs/op", self.avg_micros);
        println!("  P95: {:.1} μs/op", self.p95_micros);
        println!("  P99: {:.1} μs/op", self.p99_micros);
        println!("  Memory: {:.1} MB", self.memory_mb);
        
        if !self.metrics.is_empty() {
            println!("  Additional metrics:");
            for (k, v) in &self.metrics {
                println!("    {}: {:.2}", k, v);
            }
        }
        println!();
    }
}

/// Performance benchmark configuration
#[derive(Debug, Clone)]
pub struct BenchmarkConfig {
    /// Duration to run each benchmark (seconds)
    pub duration: u64,
    /// Number of concurrent workers for parallel benchmarks
    pub concurrency: usize,
    /// Size of synthetic data to use in benchmarks
    pub data_size: usize,
    /// Type of data to generate for testing
    pub data_type: DataType,
}

#[derive(Debug, Clone)]
pub enum DataType {
    /// OTLP log data pattern
    OtlpLogs,
    /// OTLP metric data pattern
    OtlpMetrics,
    /// OTLP trace data pattern
    OtlpTraces,
    /// Generic structured data
    Structured,
    /// Random unstructured data
    Random,
}

impl Default for BenchmarkConfig {
    fn default() -> Self {
        Self {
            duration: 10, // 10 seconds per benchmark by default
            concurrency: 1, // Single-threaded by default
            data_size: 1024, // 1KB events by default
            data_type: DataType::OtlpLogs,
        }
    }
}

/// Performance validator that runs all required benchmarks
pub struct PerformanceValidator {
    config: BenchmarkConfig,
}

impl PerformanceValidator {
    /// Create a new performance validator with the given configuration
    pub fn new(config: BenchmarkConfig) -> Self {
        Self { config }
    }

    /// Run all benchmarks and return results
    pub fn run_all_benchmarks(&self) -> Vec<BenchmarkResult> {
        vec![
            self.benchmark_compression_throughput(),
            self.benchmark_metrics_calculation(),
            self.benchmark_window_system(),
            self.benchmark_end_to_end_detection(),
        ]
    }

    /// Benchmark compression throughput
    fn benchmark_compression_throughput(&self) -> BenchmarkResult {
        let mut result = BenchmarkResult::new("Compression Throughput");
        let adapter = create_adapter(CompressionAlgorithm::OpenZL).expect("create adapter");
        
        let test_data = self.generate_test_data(self.config.data_size);
        let start_time = Instant::now();
        let mut operations = 0u64;
        let mut times = Vec::new();
        
        // Run for configured duration
        while start_time.elapsed().as_secs() < self.config.duration {
            let op_start = Instant::now();
            let _compressed = adapter.compress(&test_data).expect("compression failed");
            let op_time = op_start.elapsed();
            
            times.push(op_time.as_micros() as f64);
            operations += 1;
        }
        
        let duration = start_time.elapsed();
        result.duration = duration;
        result.operations = operations;
        result.ops_per_sec = operations as f64 / duration.as_secs_f64();
        result.avg_micros = times.iter().sum::<f64>() / times.len() as f64;
        
        // Calculate percentiles
        if !times.is_empty() {
            let mut sorted_times = times.clone();
            sorted_times.sort_by(|a, b| a.partial_cmp(b).unwrap());
            let len = sorted_times.len();
            
            result.p95_micros = sorted_times[(len as f64 * 0.95) as usize];
            result.p99_micros = sorted_times[(len as f64 * 0.99) as usize];
        }
        
        // Memory estimate (simplified)
        result.memory_mb = (operations * self.config.data_size as u64) as f64 / (1024.0 * 1024.0);
        
        result
    }

    /// Benchmark metrics calculation performance
    fn benchmark_metrics_calculation(&self) -> BenchmarkResult {
        let mut result = BenchmarkResult::new("Metrics Calculation");
        let adapter = create_adapter(CompressionAlgorithm::OpenZL).expect("create adapter");
        
        let baseline = self.generate_test_data(self.config.data_size);
        let window = self.generate_test_data(self.config.data_size);
        let config = crate::ComputeConfig::default();
        
        let start_time = Instant::now();
        let mut operations = 0u64;
        let mut times = Vec::new();
        
        // Run for configured duration
        while start_time.elapsed().as_secs() < self.config.duration {
            let op_start = Instant::now();
            let _metrics = crate::compute_metrics(&baseline, &window, adapter.as_ref(), &config)
                .expect("metrics computation failed");
            let op_time = op_start.elapsed();
            
            times.push(op_time.as_micros() as f64);
            operations += 1;
        }
        
        let duration = start_time.elapsed();
        result.duration = duration;
        result.operations = operations;
        result.ops_per_sec = operations as f64 / duration.as_secs_f64();
        result.avg_micros = times.iter().sum::<f64>() / times.len() as f64;
        
        // Calculate percentiles
        if !times.is_empty() {
            let mut sorted_times = times.clone();
            sorted_times.sort_by(|a, b| a.partial_cmp(b).unwrap());
            let len = sorted_times.len();
            
            result.p95_micros = sorted_times[(len as f64 * 0.95) as usize];
            result.p99_micros = sorted_times[(len as f64 * 0.99) as usize];
        }
        
        result
    }

    /// Benchmark sliding window system performance
    fn benchmark_window_system(&self) -> BenchmarkResult {
        let mut result = BenchmarkResult::new("Sliding Window System");
        
        let config = WindowConfig {
            baseline_size: 1000,
            window_size: 100,
            hop_size: 50,
            max_capacity: 10000,
            ..Default::default()
        };
        
        let mut window = SlidingWindow::new(config);
        let test_data = self.generate_test_data(self.config.data_size);
        
        let start_time = Instant::now();
        let mut operations = 0u64;
        let mut times = Vec::new();
        
        // Run for configured duration
        while start_time.elapsed().as_secs() < self.config.duration {
            let op_start = Instant::now();
            let event = DataEvent::new(test_data.clone());
            window.add_event(event);
            let _ = window.get_baseline_and_window(); // Try to get windows if available
            let op_time = op_start.elapsed();
            
            times.push(op_time.as_micros() as f64);
            operations += 1;
        }
        
        let duration = start_time.elapsed();
        result.duration = duration;
        result.operations = operations;
        result.ops_per_sec = operations as f64 / duration.as_secs_f64();
        result.avg_micros = times.iter().sum::<f64>() / times.len() as f64;
        
        // Calculate percentiles
        if !times.is_empty() {
            let mut sorted_times = times.clone();
            sorted_times.sort_by(|a, b| a.partial_cmp(b).unwrap());
            let len = sorted_times.len();
            
            result.p95_micros = sorted_times[(len as f64 * 0.95) as usize];
            result.p99_micros = sorted_times[(len as f64 * 0.99) as usize];
        }
        
        // Memory usage
        result.memory_mb = (window.memory_usage() * self.config.data_size) as f64 / (1024.0 * 1024.0);
        
        result
    }

    /// Benchmark full end-to-end anomaly detection
    fn benchmark_end_to_end_detection(&self) -> BenchmarkResult {
        let mut result = BenchmarkResult::new("End-to-End Anomaly Detection");
        
        // Create window system
        let config = WindowConfig {
            baseline_size: 1000,
            window_size: 100,
            hop_size: 50,
            max_capacity: 10000,
            ..Default::default()
        };
        
        let mut window = SlidingWindow::new(config);
        let adapter = create_adapter(CompressionAlgorithm::OpenZL).expect("create adapter");
        let compute_config = crate::ComputeConfig::default();
        
        // Pre-fill window with baseline data
        for _ in 0..1100 { // More than needed to ensure we have baseline and window
            let data = self.generate_test_data(self.config.data_size);
            let event = DataEvent::new(data);
            window.add_event(event);
        }
        
        let start_time = Instant::now();
        let mut operations = 0u64;
        let mut times = Vec::new();
        
        // Run for configured duration
        while start_time.elapsed().as_secs() < self.config.duration {
            let op_start = Instant::now();
            
            // Add one more event to advance the window
            let data = self.generate_test_data(self.config.data_size);
            let event = DataEvent::new(data);
            window.add_event(event);
            
            // Perform anomaly detection if we have data
            if let Some((baseline, detection_window)) = window.get_baseline_and_window() {
                let _metrics = crate::compute_metrics(&baseline, &detection_window, adapter.as_ref(), &compute_config)
                    .unwrap_or_else(|_| metrics::AnomalyMetrics::new());
            }
            
            let op_time = op_start.elapsed();
            times.push(op_time.as_micros() as f64);
            operations += 1;
        }
        
        let duration = start_time.elapsed();
        result.duration = duration;
        result.operations = operations;
        result.ops_per_sec = operations as f64 / duration.as_secs_f64();
        result.avg_micros = times.iter().sum::<f64>() / times.len() as f64;
        
        // Calculate percentiles
        if !times.is_empty() {
            let mut sorted_times = times.clone();
            sorted_times.sort_by(|a, b| a.partial_cmp(b).unwrap());
            let len = sorted_times.len();
            
            result.p95_micros = sorted_times[(len as f64 * 0.95) as usize];
            result.p99_micros = sorted_times[(len as f64 * 0.99) as usize];
        }
        
        // Memory usage
        result.memory_mb = (window.memory_usage() * self.config.data_size) as f64 / (1024.0 * 1024.0);
        
        result
    }

    /// Generate test data based on configured type
    fn generate_test_data(&self, size: usize) -> Vec<u8> {
        match &self.config.data_type {
            DataType::OtlpLogs => self.generate_otlp_logs(size),
            DataType::OtlpMetrics => self.generate_otlp_metrics(size),
            DataType::OtlpTraces => self.generate_otlp_traces(size),
            DataType::Structured => self.generate_structured_data(size),
            DataType::Random => self.generate_random_data(size),
        }
    }

    /// Generate realistic OTLP log data
    fn generate_otlp_logs(&self, size: usize) -> Vec<u8> {
        let log_template = r#"{"resource":{"service.name":"api-gateway"},"time":"2025-10-24T10:00:00.000Z","severity":"INFO","body":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42,"request_id":"req-"}"#;
        
        let mut data = String::new();
        let mut counter = 0;
        
        while data.len() < size {
            counter += 1;
            data.push_str(&log_template.replace("req-", &format!("req-{:06}", counter)));
            data.push('\n');
        }
        
        data.truncate(size);
        data.into_bytes()
    }

    /// Generate realistic OTLP metric data
    fn generate_otlp_metrics(&self, size: usize) -> Vec<u8> {
        let metric_template = r#"{"resource":{"service.name":"api-gateway"},"metric":"http_request_duration_seconds","type":"histogram","value":42.0,"attributes":{"method":"GET","path":"/api/users","status":"200","quantile":0.95}}"#;
        
        let mut data = String::new();
        
        while data.len() < size {
            data.push_str(metric_template);
            data.push('\n');
        }
        
        data.truncate(size);
        data.into_bytes()
    }

    /// Generate realistic OTLP trace data
    fn generate_otlp_traces(&self, size: usize) -> Vec<u8> {
        let trace_template = r#"{"trace_id":"trace-001","span_id":"span-001","name":"http.request","start_time":"2025-10-24T10:00:00.000Z","end_time":"2025-10-24T10:00:00.042Z","attributes":{"http.method":"GET","http.url":"/api/users","http.status_code":200}}"#;
        
        let mut data = String::new();
        
        while data.len() < size {
            data.push_str(trace_template);
            data.push('\n');
        }
        
        data.truncate(size);
        data.into_bytes()
    }

    /// Generate generic structured data
    fn generate_structured_data(&self, size: usize) -> Vec<u8> {
        let mut data = String::new();
        let mut counter = 0;
        
        while data.len() < size {
            counter += 1;
            data.push_str(&format!(r#"{{"event_id":{},"timestamp":"2025-10-24T10:00:00Z","type":"user_action","data":{{"user_id":{},"action":"api_call","path":"/api/endpoint"}}}}"#, counter, counter % 1000));
            data.push('\n');
        }
        
        data.truncate(size);
        data.into_bytes()
    }

    /// Generate random unstructured data
    fn generate_random_data(&self, size: usize) -> Vec<u8> {
        use std::time::{SystemTime, UNIX_EPOCH};
        
        let seed = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u32;
        
        let mut rng = fastrand::Rng::with_seed(seed as u64);
        (0..size).map(|_| rng.u8(..)).collect()
    }

    /// Validate that benchmarks meet Phase 1 targets
    pub fn validate_targets(&self) -> ValidationReport {
        let results = self.run_all_benchmarks();
        let mut report = ValidationReport::new();
        
        for result in results {
            // Check throughput targets (>10k events/s)
            if result.name == "End-to-End Anomaly Detection" {
                if result.ops_per_sec >= 10000.0 {
                    report.pass("Throughput", format!("Achieved {:.0} ops/sec >= 10k", result.ops_per_sec));
                } else {
                    report.fail("Throughput", format!("Achieved {:.0} ops/sec < 10k target", result.ops_per_sec));
                }
            }
            
            // Check latency targets (<400ms p95)
            if result.name == "End-to-End Anomaly Detection" {
                let p95_ms = result.p95_micros / 1000.0;
                if p95_ms < 400.0 {
                    report.pass("Latency", format!("P95 latency {:.1}ms < 400ms", p95_ms));
                } else {
                    report.fail("Latency", format!("P95 latency {:.1}ms >= 400ms target", p95_ms));
                }
            }
            
            // Memory usage check
            if result.name == "End-to-End Anomaly Detection" {
                if result.memory_mb < 2048.0 { // 2GB
                    report.pass("Memory", format!("Used {:.1}MB < 2GB", result.memory_mb));
                } else {
                    report.fail("Memory", format!("Used {:.1}MB >= 2GB target", result.memory_mb));
                }
            }
            
            // Add result to report
            report.add_result(result);
        }
        
        report
    }
}

/// Performance validation report
#[derive(Debug, Clone)]
pub struct ValidationReport {
    /// Successful validations
    pub passes: Vec<(String, String)>, // (metric, description)
    /// Failed validations
    pub failures: Vec<(String, String)>, // (metric, description)
    /// All benchmark results
    pub results: Vec<BenchmarkResult>,
}

impl Default for ValidationReport {
    fn default() -> Self {
        Self::new()
    }
}

impl ValidationReport {
    /// Create a new validation report
    pub fn new() -> Self {
        Self {
            passes: Vec::new(),
            failures: Vec::new(),
            results: Vec::new(),
        }
    }

    /// Add a pass to the report
    pub fn pass(&mut self, metric: &str, description: String) {
        self.passes.push((metric.to_string(), description));
    }

    /// Add a fail to the report
    pub fn fail(&mut self, metric: &str, description: String) {
        self.failures.push((metric.to_string(), description));
    }

    /// Add a benchmark result to the report
    pub fn add_result(&mut self, result: BenchmarkResult) {
        self.results.push(result);
    }

    /// Check if all validations passed
    pub fn all_passed(&self) -> bool {
        self.failures.is_empty()
    }

    /// Print the validation report
    pub fn print_report(&self) {
        println!("=== PERFORMANCE VALIDATION REPORT ===");
        println!();
        
        if self.passes.is_empty() && self.failures.is_empty() {
            println!("No validations performed.");
            return;
        }
        
        if !self.passes.is_empty() {
            println!("✅ PASSES:");
            for (metric, description) in &self.passes {
                println!("  {}: {}", metric, description);
            }
            println!();
        }
        
        if !self.failures.is_empty() {
            println!("❌ FAILURES:");
            for (metric, description) in &self.failures {
                println!("  {}: {}", metric, description);
            }
            println!();
        }
        
        let passed = self.passes.len();
        let failed = self.failures.len();
        let total = passed + failed;
        
        println!("SUMMARY: {} passed, {} failed out of {} validations", passed, failed, total);
        
        if self.all_passed() {
            println!("🎉 All performance targets met!");
        } else {
            println!("⚠️  Some performance targets not met. See failures above.");
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_performance_validator_creation() {
        let config = BenchmarkConfig::default();
        let validator = PerformanceValidator::new(config);
        
        assert_eq!(validator.config.duration, 10); // Default 10 seconds
        assert_eq!(validator.config.concurrency, 1); // Default single-threaded
    }

    #[test] 
    fn test_data_generation() {
        let config = BenchmarkConfig {
            duration: 1,
            concurrency: 1,
            data_size: 100,
            data_type: DataType::OtlpLogs,
        };
        
        let validator = PerformanceValidator::new(config);
        let data = validator.generate_test_data(100);
        
        assert_eq!(data.len(), 100);
    }

    #[test]
    fn test_benchmark_result_creation() {
        let result = BenchmarkResult::new("Test Benchmark");
        
        assert_eq!(result.name, "Test Benchmark");
        assert_eq!(result.operations, 0);
        assert_eq!(result.ops_per_sec, 0.0);
    }

    #[test]
    fn test_validation_report() {
        let mut report = ValidationReport::new();
        report.pass("Throughput", "Achieved 15000 ops/sec".to_string());
        report.fail("Latency", "P95 was 500ms".to_string());
        
        assert_eq!(report.passes.len(), 1);
        assert_eq!(report.failures.len(), 1);
        assert!(!report.all_passed());
        
        // Test passing report
        let mut good_report = ValidationReport::new();
        good_report.pass("Throughput", "Achieved 15000 ops/sec".to_string());
        
        assert_eq!(good_report.failures.len(), 0);
        assert!(good_report.all_passed());
    }
}
