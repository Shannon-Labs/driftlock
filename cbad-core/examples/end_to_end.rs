//! Example: End-to-End CBAD Anomaly Detection
//! 
//! This example demonstrates the complete CBAD workflow:
//! 1. Create synthetic baseline and anomalous data
//! 2. Use the sliding window system to manage data streams
//! 3. Perform anomaly detection with multiple metrics
//! 4. Display results with glass-box explanations

use cbad_core::{compute_metrics, ComputeConfig, compression::create_adapter, compression::CompressionAlgorithm};
use cbad_core::window::{SlidingWindow, WindowConfig, DataEvent};

fn main() {
    println!("🔍 Driftlock CBAD - End-to-End Example");
    println!("=====================================");
    
    // Step 1: Create a compression adapter (using OpenZL if available, otherwise zstd)
    #[cfg(feature = "openzl")]
    let adapter = create_adapter(CompressionAlgorithm::OpenZL)
        .or_else(|_| create_adapter(CompressionAlgorithm::Zstd))
        .expect("Failed to create compression adapter");

    #[cfg(not(feature = "openzl"))]
    let adapter = create_adapter(CompressionAlgorithm::Zstd)
        .expect("Failed to create compression adapter");
    
    println!("✅ Compression adapter created: {}", adapter.name());
    
    // Step 2: Create synthetic data - baseline (normal) and anomalous patterns
    let baseline_data = r#"{"service":"api-gateway","level":"INFO","message":"Request completed","attributes":{"method":"GET","path":"/api/users","status":200,"duration_ms":42}}"#.repeat(50);
    let anomalous_data = r#"{"service":"api-gateway","level":"ERROR","message":"PANIC OCCURRED","attributes":{"error_type":"stack_overflow","stack_trace":"thread 'main' panicked at 'index out of bounds', src/main.rs:42:13","exception":"critical"}}"#.repeat(10);
    
    println!("📊 Baseline data (normal): {} bytes", baseline_data.len());
    println!("🚨 Anomalous data: {} bytes", anomalous_data.len());
    
    // Step 3: Perform direct metrics computation
    println!("\n📈 Direct Metrics Computation");
    println!("----------------------------");
    
    let config = ComputeConfig::default();
    match compute_metrics(
        baseline_data.as_bytes(),
        anomalous_data.as_bytes(),
        adapter.as_ref(),
        &config,
    ) {
        Ok(metrics) => {
            println!("NCD Score: {:.3} (0.0 = identical, 1.0 = completely different)", metrics.ncd);
            println!("P-Value: {:.3} (statistical significance)", metrics.p_value);
            println!("Compression Ratio - Baseline: {:.2}x, Window: {:.2}x", 
                     metrics.baseline_compression_ratio, metrics.window_compression_ratio);
            println!("Entropy - Baseline: {:.2} bits/byte, Window: {:.2} bits/byte", 
                     metrics.baseline_entropy, metrics.window_entropy);
            println!("Anomaly Detected: {}", if metrics.is_anomaly { "✅ YES" } else { "❌ NO" });
            println!("\n📋 Glass-Box Explanation:");
            println!("{}", metrics.explanation);
        }
        Err(e) => {
            eprintln!("❌ Error computing metrics: {}", e);
        }
    }
    
    // Step 4: Demonstrate sliding window functionality
    println!("\n🔄 Sliding Window System Demo");
    println!("----------------------------");
    
    let window_config = WindowConfig {
        baseline_size: 10, // 10 events for baseline
        window_size: 5,    // 5 events for detection window
        max_capacity: 100, // Max 100 events in memory
        ..Default::default()
    };
    
    let mut sliding_window = SlidingWindow::new(window_config);
    println!(
        "✅ Sliding window initialized with baseline_size={}, window_size={}",
        sliding_window.config().baseline_size,
        sliding_window.config().window_size
    );
    
    // Add baseline events to the window
    for i in 0..15 {
        let event_data = if i < 10 {
            // Baseline events (normal pattern)
            format!(r#"{{"event_id":{},"level":"INFO","service":"api-gateway","action":"request_completed","duration_ms":{}}}"#, 
                    i, 40 + (i % 5))
        } else {
            // Some anomalous events
            format!(r#"{{"event_id":{},"level":"ERROR","service":"api-gateway","action":"critical_error","exception":"stack_overflow"}}"#, i)
        };
        
        let event = DataEvent::new(event_data.as_bytes().to_vec());
        let _ = sliding_window.add_event(event);
    }
    
    println!("📊 Added {} events to sliding window", sliding_window.total_events());
    println!("💾 Memory usage: {} events", sliding_window.memory_usage());
    
    // Check if we can get baseline and window data
    match sliding_window.get_baseline_and_window() {
        Some((baseline, window)) => {
            println!("✅ Retrieved baseline ({} bytes) and window ({} bytes) from sliding window", 
                     baseline.len(), window.len());
            
            // Perform another detection on the sliding window data
            match compute_metrics(&baseline, &window, adapter.as_ref(), &config) {
                Ok(metrics) => {
                    println!("🔍 Anomaly Detection on Window Data:");
                    println!("   NCD: {:.3}", metrics.ncd);
                    println!("   P-Value: {:.3}", metrics.p_value);
                    println!("   Anomaly: {}", if metrics.is_anomaly { "YES" } else { "NO" });
                }
                Err(e) => {
                    eprintln!("❌ Error computing window metrics: {}", e);
                }
            }
        }
        None => {
            println!("⚠️  Not enough data in window for analysis yet");
        }
    }
    
    // Step 5: Performance validation example
    println!("\n⏱️  Performance Validation Example");
    println!("--------------------------------");
    
    use cbad_core::performance::{PerformanceValidator, BenchmarkConfig, DataType};
    
    let perf_config = BenchmarkConfig {
        duration: 2,          // Run for 2 seconds to get a quick sample
        data_type: DataType::OtlpLogs,
        data_size: 512,       // 512 bytes per event
        ..Default::default()
    };
    
    let duration = perf_config.duration; // Store the duration before moving perf_config
    let validator = PerformanceValidator::new(perf_config);
    println!("✅ Running performance validation ({}s duration)...", duration);
    
    let report = validator.validate_targets();
    report.print_report();
    
    println!("\n🎯 CBAD Example Complete");
    println!("=======================");
    println!("This example demonstrated:");
    println!("• Direct anomaly detection with multiple metrics");
    println!("• Sliding window system for streaming data");
    println!("• Performance validation framework");
    println!("• Glass-box explanations for audit compliance");
}
