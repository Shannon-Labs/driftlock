//! Test CBAD detection across multiple datasets
//!
//! Run with: cargo run --example test_datasets --release

use cbad_core::anomaly::{AnomalyConfig, AnomalyDetector};
use cbad_core::window::WindowConfig;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::path::Path;
use std::time::Instant;

fn main() {
    println!("\n=== CBAD Anomaly Detection Test Suite ===\n");

    // Test 1: Financial transactions (mixed normal + fraud)
    test_jsonl_dataset(
        "Financial Transactions (Mixed)",
        "../test-data/mixed-transactions.jsonl",
        400, // baseline size
        50,  // window size
    );

    // Test 2: Fraud dataset
    test_jsonl_dataset(
        "Fraud Detection Dataset",
        "../test-data/fraud/fraud_sample.ndjson",
        200,
        30,
    );

    // Test 3: Anomalous transactions only
    test_jsonl_dataset(
        "Anomalous Transactions Only",
        "../test-data/anomalous-transactions.jsonl",
        50,
        20,
    );

    // Test 4: Normal transactions (should have low anomaly rate)
    test_jsonl_dataset(
        "Normal Transactions Only",
        "../test-data/normal-transactions.jsonl",
        300,
        50,
    );

    // Test 5: Terra Luna crash data (crypto) - uses pre-processed JSON
    test_csv_timeseries(
        "Terra Luna Crash (Crypto)",
        "../test-data/terra_luna/driftlock_ready.json",
    );

    // Test 6: AWS CloudWatch metrics
    test_csv_timeseries(
        "AWS CloudWatch CPU",
        "../test-data/web_traffic/realAWSCloudwatch/realAWSCloudwatch/ec2_cpu_utilization_825cc2.csv",
    );

    // Test 7: Terra Luna raw CSV
    test_csv_timeseries(
        "Terra Luna CSV (Raw)",
        "../test-data/terra_luna/terra-luna.csv",
    );

    // Test synthetic anomaly types
    println!("\n=== Synthetic Anomaly Detection Tests ===\n");
    test_synthetic_anomalies();

    // Test tokenizer impact on detection quality
    println!("\n=== Tokenizer Comparison Tests ===\n");
    test_tokenizer_comparison();

    println!("\n=== All Tests Complete ===\n");
}

fn test_jsonl_dataset(name: &str, path: &str, baseline_size: usize, window_size: usize) {
    println!("📊 Testing: {}", name);
    println!("   Path: {}", path);

    let full_path = Path::new(env!("CARGO_MANIFEST_DIR")).join(path);
    if !full_path.exists() {
        println!("   ⚠️  File not found, skipping\n");
        return;
    }

    let file = match File::open(&full_path) {
        Ok(f) => f,
        Err(e) => {
            println!("   ⚠️  Error opening file: {}\n", e);
            return;
        }
    };
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().filter_map(Result::ok).collect();

    println!("   Events loaded: {}", lines.len());

    if lines.len() < baseline_size + window_size {
        println!("   ⚠️  Not enough events for analysis\n");
        return;
    }

    // Configure detector
    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size,
            window_size,
            hop_size: window_size / 2,
            max_capacity: baseline_size + window_size * 3,
            ..Default::default()
        },
        permutation_count: 100, // Faster for testing
        ncd_threshold: 0.25,
        p_value_threshold: 0.05,
        require_statistical_significance: true,
        ..Default::default()
    };

    let detector = match AnomalyDetector::new(config) {
        Ok(d) => d,
        Err(e) => {
            println!("   ⚠️  Error creating detector: {:?}\n", e);
            return;
        }
    };
    let start = Instant::now();
    let mut anomalies_detected = 0;
    let mut detection_count = 0;

    // Feed events into detector
    for line in &lines {
        let _ = detector.add_data(line.as_bytes().to_vec());

        // Check for anomalies periodically
        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(result)) = detector.detect_anomaly() {
                detection_count += 1;
                if result.is_anomaly {
                    anomalies_detected += 1;
                    if anomalies_detected <= 3 {
                        println!(
                            "   🔴 Anomaly #{}: NCD={:.3}, p={:.4}, conf={:.1}%",
                            anomalies_detected,
                            result.metrics.ncd,
                            result.metrics.p_value,
                            result.confidence_level * 100.0
                        );
                    }
                }
            }
        }
    }

    let elapsed = start.elapsed();
    let anomaly_rate = if detection_count > 0 {
        anomalies_detected as f64 / detection_count as f64
    } else {
        0.0
    };

    println!("   ✅ Results:");
    println!("      - Detection cycles: {}", detection_count);
    println!("      - Anomalies detected: {}", anomalies_detected);
    println!("      - Anomaly rate: {:.2}%", anomaly_rate * 100.0);
    println!("      - Processing time: {:?}", elapsed);
    println!(
        "      - Throughput: {:.0} events/sec\n",
        lines.len() as f64 / elapsed.as_secs_f64()
    );
}

fn test_csv_timeseries(name: &str, path: &str) {
    println!("📈 Testing: {}", name);
    println!("   Path: {}", path);

    let full_path = Path::new(env!("CARGO_MANIFEST_DIR")).join(path);
    if !full_path.exists() {
        println!("   ⚠️  File not found, skipping\n");
        return;
    }

    // For JSON payload files, parse and extract events
    if path.ends_with(".json") {
        let contents = match std::fs::read_to_string(&full_path) {
            Ok(c) => c,
            Err(e) => {
                println!("   ⚠️  Error reading file: {}\n", e);
                return;
            }
        };

        // Parse JSON and convert to events
        let json: serde_json::Value = match serde_json::from_str(&contents) {
            Ok(j) => j,
            Err(e) => {
                println!("   ⚠️  Error parsing JSON: {}\n", e);
                return;
            }
        };

        let events: Vec<Vec<u8>> = if let Some(arr) = json.as_array() {
            arr.iter()
                .map(|v| serde_json::to_vec(v).unwrap_or_default())
                .collect()
        } else if let Some(events) = json.get("events").and_then(|e| e.as_array()) {
            events
                .iter()
                .map(|v| serde_json::to_vec(v).unwrap_or_default())
                .collect()
        } else {
            vec![serde_json::to_vec(&json).unwrap_or_default()]
        };

        println!("   Events extracted: {}", events.len());
        run_timeseries_detection(events);
        return;
    }

    // For CSV files
    let file = match File::open(&full_path) {
        Ok(f) => f,
        Err(e) => {
            println!("   ⚠️  Error opening file: {}\n", e);
            return;
        }
    };
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().filter_map(Result::ok).collect();

    println!("   Data points: {}", lines.len());

    // Convert CSV rows to bytes
    let events: Vec<Vec<u8>> = lines.iter().map(|l| l.as_bytes().to_vec()).collect();
    run_timeseries_detection(events);
}

fn run_timeseries_detection(events: Vec<Vec<u8>>) {
    if events.len() < 100 {
        println!("   ⚠️  Not enough data points\n");
        return;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 50,
            window_size: 20,
            hop_size: 10,
            max_capacity: 200,
            ..Default::default()
        },
        permutation_count: 50,
        ncd_threshold: 0.20,
        ..Default::default()
    };

    let detector = match AnomalyDetector::new(config) {
        Ok(d) => d,
        Err(e) => {
            println!("   ⚠️  Error creating detector: {:?}\n", e);
            return;
        }
    };
    let start = Instant::now();
    let mut anomalies: Vec<(usize, f64, f64)> = Vec::new();

    for (i, event) in events.iter().enumerate() {
        let _ = detector.add_data(event.clone());

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(result)) = detector.detect_anomaly() {
                if result.is_anomaly {
                    anomalies.push((i, result.metrics.ncd, result.confidence_level));
                }
            }
        }
    }

    let elapsed = start.elapsed();

    println!("   ✅ Results:");
    println!("      - Anomalies detected: {}", anomalies.len());
    if !anomalies.is_empty() {
        println!("      - First anomaly at index: {}", anomalies[0].0);
        println!(
            "      - Peak NCD: {:.3}",
            anomalies
                .iter()
                .map(|(_, ncd, _)| *ncd)
                .fold(0.0_f64, f64::max)
        );
    }
    println!("      - Processing time: {:?}\n", elapsed);
}

fn test_synthetic_anomalies() {
    use cbad_core::anomaly::synthetic::{generate_anomaly, generate_baseline, AnomalyType};

    let anomaly_types = [
        ("Volume Spike", AnomalyType::VolumeSpike),
        ("Random Noise", AnomalyType::RandomNoise),
        ("Pattern Break", AnomalyType::PatternBreak),
        ("Data Corruption", AnomalyType::DataCorruption),
        ("Gradual Drift", AnomalyType::GradualDrift),
    ];

    for (name, anomaly_type) in anomaly_types {
        println!("🧪 Testing: {} detection", name);

        let config = AnomalyConfig {
            window_config: WindowConfig {
                baseline_size: 80,
                window_size: 30,
                hop_size: 15,
                max_capacity: 200,
                ..Default::default()
            },
            permutation_count: 100,
            ncd_threshold: 0.20,
            require_statistical_significance: false, // More sensitive for synthetic tests
            ..Default::default()
        };

        let detector = match AnomalyDetector::new(config) {
            Ok(d) => d,
            Err(e) => {
                println!("   ⚠️  Error creating detector: {:?}", e);
                continue;
            }
        };

        // Feed baseline (normal) data
        let baseline = generate_baseline(100, 42);
        for event in baseline {
            let _ = detector.add_data(event);
        }

        // Feed anomalous data
        let anomalies = generate_anomaly(40, anomaly_type, 43);
        let mut detected = false;
        let mut max_ncd: f64 = 0.0;

        for event in anomalies {
            let _ = detector.add_data(event);
            if detector.is_ready().unwrap_or(false) {
                if let Ok(Some(result)) = detector.detect_anomaly() {
                    max_ncd = max_ncd.max(result.metrics.ncd);
                    if result.is_anomaly {
                        detected = true;
                    }
                }
            }
        }

        if detected {
            println!("   ✅ Detected! Max NCD: {:.3}", max_ncd);
        } else {
            println!("   ⚠️  Not detected. Max NCD: {:.3}", max_ncd);
        }
    }
}

/// Test the impact of tokenization on detection quality
///
/// Uses direct NCD comparison to demonstrate how tokenization reduces
/// the distance between structurally similar data with different high-entropy values.
fn test_tokenizer_comparison() {
    use cbad_core::compression::{create_adapter, CompressionAlgorithm};
    use cbad_core::metrics::compute_metrics_with_tokenizer;
    use cbad_core::Tokenizer;

    println!("🔬 Testing tokenizer impact on NCD scores\n");

    // Create adapter
    let adapter = create_adapter(CompressionAlgorithm::Zstd).expect("create adapter");

    // Create baseline and test data with same structure but different UUIDs
    let baseline = generate_uuid_heavy_data_seeded(100, 42);
    let normal_different_uuids = generate_uuid_heavy_data_seeded(30, 100);
    let anomalous = generate_anomalous_pattern(30);

    // Concatenate for NCD calculation
    let baseline_bytes: Vec<u8> = baseline.iter().flat_map(|e| e.clone()).collect();
    let normal_bytes: Vec<u8> = normal_different_uuids
        .iter()
        .flat_map(|e| e.clone())
        .collect();
    let anomalous_bytes: Vec<u8> = anomalous.iter().flat_map(|e| e.clone()).collect();

    // Test 1: NCD WITHOUT tokenizer
    println!("   📊 Test 1: NCD scores WITHOUT Tokenizer");
    let metrics_normal_no_tok = compute_metrics_with_tokenizer(
        &baseline_bytes,
        &normal_bytes,
        adapter.as_ref(),
        50,
        42,
        None,
    )
    .expect("compute metrics");

    let metrics_anom_no_tok = compute_metrics_with_tokenizer(
        &baseline_bytes,
        &anomalous_bytes,
        adapter.as_ref(),
        50,
        42,
        None,
    )
    .expect("compute metrics");

    println!(
        "      - Normal data (different UUIDs): NCD = {:.3}",
        metrics_normal_no_tok.ncd
    );
    println!(
        "      - Anomalous data:                 NCD = {:.3}",
        metrics_anom_no_tok.ncd
    );
    println!(
        "      - Separation gap:                 {:.3}",
        metrics_anom_no_tok.ncd - metrics_normal_no_tok.ncd
    );

    // Test 2: NCD WITH tokenizer
    println!("\n   📊 Test 2: NCD scores WITH Tokenizer");
    let tokenizer = Tokenizer::default();

    let metrics_normal_tok = compute_metrics_with_tokenizer(
        &baseline_bytes,
        &normal_bytes,
        adapter.as_ref(),
        50,
        42,
        Some(&tokenizer),
    )
    .expect("compute metrics");

    let metrics_anom_tok = compute_metrics_with_tokenizer(
        &baseline_bytes,
        &anomalous_bytes,
        adapter.as_ref(),
        50,
        42,
        Some(&tokenizer),
    )
    .expect("compute metrics");

    println!(
        "      - Normal data (UUIDs normalized): NCD = {:.3}",
        metrics_normal_tok.ncd
    );
    println!(
        "      - Anomalous data:                 NCD = {:.3}",
        metrics_anom_tok.ncd
    );
    println!(
        "      - Separation gap:                 {:.3}",
        metrics_anom_tok.ncd - metrics_normal_tok.ncd
    );

    // Summary
    println!("\n   📈 Analysis:");

    let ncd_drop = metrics_normal_no_tok.ncd - metrics_normal_tok.ncd;
    if ncd_drop > 0.0 {
        println!(
            "      ✅ Tokenizer reduced NCD on normal data by {:.3}",
            ncd_drop
        );
        println!("         (Lower NCD = more similar to baseline = fewer false positives)");
    } else {
        println!("      ➖ Tokenizer did not significantly change NCD on normal data");
    }

    let gap_without = metrics_anom_no_tok.ncd - metrics_normal_no_tok.ncd;
    let gap_with = metrics_anom_tok.ncd - metrics_normal_tok.ncd;
    if gap_with > gap_without {
        println!(
            "      ✅ Separation gap improved: {:.3} → {:.3}",
            gap_without, gap_with
        );
        println!("         (Larger gap = easier to distinguish anomalies)");
    }

    // Show tokenizer stats
    let stats = tokenizer.stats();
    println!("\n   📊 Tokenizer Statistics:");
    println!("      - UUIDs replaced: {}", stats.uuid_count);
    println!("      - Hashes replaced: {}", stats.hash_count);
    println!("      - Bytes saved: {}", stats.bytes_saved);
}

/// Generate data with high-entropy fields (UUIDs, hashes)
fn generate_uuid_heavy_data_seeded(count: usize, seed: u64) -> Vec<Vec<u8>> {
    use rand::prelude::*;
    let mut rng = StdRng::seed_from_u64(seed);
    let mut events = Vec::new();

    for i in 0..count {
        // Generate random UUID
        let uuid = format!(
            "{:08x}-{:04x}-{:04x}-{:04x}-{:012x}",
            rng.gen::<u32>(),
            rng.gen::<u16>(),
            rng.gen::<u16>(),
            rng.gen::<u16>(),
            rng.gen::<u64>() & 0xFFFFFFFFFFFF
        );

        // Generate random hash
        let hash: String = (0..64)
            .map(|_| format!("{:x}", rng.gen::<u8>() % 16))
            .collect();

        let log = format!(
            r#"{{"timestamp":"2025-10-24T{:02}:{:02}:{:02}Z","user_id":"{}","request_hash":"{}","action":"login","status":"success"}}"#,
            i / 3600,
            (i % 3600) / 60,
            i % 60,
            uuid,
            hash
        );

        events.push(log.into_bytes());
    }

    events
}

/// Generate anomalous pattern (different structure, no UUIDs)
fn generate_anomalous_pattern(count: usize) -> Vec<Vec<u8>> {
    use rand::prelude::*;
    let mut rng = StdRng::seed_from_u64(43);
    let mut events = Vec::new();

    for _ in 0..count {
        // Completely different log structure with stack traces
        let stack_trace = format!(
            "panic at line {}: index out of bounds (len={}, index={})",
            rng.gen_range(1..500),
            rng.gen_range(10..100),
            rng.gen_range(100..1000)
        );

        let log = format!(
            r#"{{"severity":"PANIC","error":"{}","binary_dump":"{}"}}"#,
            stack_trace,
            (0..100)
                .map(|_| format!("{:02x}", rng.gen_range(0u8..255)))
                .collect::<String>()
        );

        events.push(log.into_bytes());
    }

    events
}
