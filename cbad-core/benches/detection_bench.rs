//! Performance benchmarks for CBAD anomaly detection
//!
//! Run with: cargo bench
//!
//! Target performance metrics:
//! - Tokenization: <100μs per event
//! - Compression: <200μs per event
//! - NCD calculation: <500μs per event
//! - Full detection cycle: <1ms per event

use cbad_core::anomaly::{AnomalyConfig, AnomalyDetector};
use cbad_core::compression::{create_adapter, CompressionAlgorithm};
use cbad_core::metrics::compute_metrics_with_tokenizer;
use cbad_core::tokenizer::{Tokenizer, TokenizerConfig};
use cbad_core::window::WindowConfig;
use criterion::{black_box, criterion_group, criterion_main, BenchmarkId, Criterion, Throughput};

/// Generate a sample JSON event with high-entropy fields
fn generate_sample_event(seed: u64) -> Vec<u8> {
    format!(
        r#"{{"timestamp":"2024-01-01T00:00:{:02}Z","user_id":"{:08x}-{:04x}-{:04x}-{:04x}-{:012x}","action":"login","ip":"192.168.1.{}","hash":"{:064x}","data":"event_{}"}}"#,
        seed % 60,
        seed,
        (seed >> 8) & 0xFFFF,
        (seed >> 16) & 0xFFFF,
        (seed >> 24) & 0xFFFF,
        seed,
        seed % 255,
        seed,
        seed
    )
    .into_bytes()
}

/// Generate multiple events for batch benchmarks
fn generate_events(count: usize, seed: u64) -> Vec<Vec<u8>> {
    (0..count)
        .map(|i| generate_sample_event(seed.wrapping_add(i as u64)))
        .collect()
}

/// Benchmark tokenization performance
fn bench_tokenization(c: &mut Criterion) {
    let mut group = c.benchmark_group("tokenization");
    let tokenizer = Tokenizer::new(TokenizerConfig::default());

    // Single event tokenization
    let event = generate_sample_event(42);
    group.throughput(Throughput::Bytes(event.len() as u64));
    group.bench_function("single_event", |b| {
        b.iter(|| tokenizer.tokenize(black_box(&event)))
    });

    // Batch tokenization
    let events = generate_events(100, 42);
    let total_bytes: u64 = events.iter().map(|e| e.len() as u64).sum();
    group.throughput(Throughput::Bytes(total_bytes));
    group.bench_function("batch_100_events", |b| {
        b.iter(|| {
            for event in &events {
                let _ = tokenizer.tokenize(black_box(event));
            }
        })
    });

    group.finish();
}

/// Benchmark compression performance with different algorithms
fn bench_compression(c: &mut Criterion) {
    let mut group = c.benchmark_group("compression");

    let event = generate_sample_event(42);
    group.throughput(Throughput::Bytes(event.len() as u64));

    // Test each compression algorithm
    for algo in [
        CompressionAlgorithm::Zstd,
        CompressionAlgorithm::Lz4,
        CompressionAlgorithm::Gzip,
        CompressionAlgorithm::Zlab,
    ] {
        let adapter = create_adapter(algo).unwrap();
        group.bench_with_input(
            BenchmarkId::from_parameter(format!("{:?}", algo)),
            &event,
            |b, event| b.iter(|| adapter.compress(black_box(event))),
        );
    }

    group.finish();
}

/// Benchmark NCD calculation
fn bench_ncd_calculation(c: &mut Criterion) {
    let mut group = c.benchmark_group("ncd_calculation");

    let adapter = create_adapter(CompressionAlgorithm::Zstd).unwrap();

    // Generate baseline and window data
    let baseline: Vec<u8> = generate_events(50, 1).into_iter().flatten().collect();
    let window: Vec<u8> = generate_events(20, 100).into_iter().flatten().collect();

    let total_bytes = (baseline.len() + window.len()) as u64;
    group.throughput(Throughput::Bytes(total_bytes));

    // NCD without tokenizer
    group.bench_function("without_tokenizer", |b| {
        b.iter(|| {
            compute_metrics_with_tokenizer(
                black_box(&baseline),
                black_box(&window),
                adapter.as_ref(),
                10, // Reduced permutations for faster benchmark
                42,
                None,
            )
        })
    });

    // NCD with tokenizer
    let tokenizer = Tokenizer::new(TokenizerConfig::default());
    group.bench_function("with_tokenizer", |b| {
        b.iter(|| {
            compute_metrics_with_tokenizer(
                black_box(&baseline),
                black_box(&window),
                adapter.as_ref(),
                10,
                42,
                Some(&tokenizer),
            )
        })
    });

    group.finish();
}

/// Benchmark full detection cycle
fn bench_full_detection(c: &mut Criterion) {
    let mut group = c.benchmark_group("full_detection");

    // Configure detector
    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 50,
            window_size: 20,
            hop_size: 10,
            max_capacity: 200,
            ..Default::default()
        },
        permutation_count: 50, // Reduced for benchmarking
        ..Default::default()
    };

    let detector = AnomalyDetector::new(config).unwrap();

    // Fill with baseline data
    let baseline_events = generate_events(50, 1);
    for event in &baseline_events {
        let _ = detector.add_data(event.clone());
    }

    // Benchmark detection on new events
    let test_events = generate_events(20, 100);
    let avg_event_size = test_events.iter().map(|e| e.len()).sum::<usize>() / test_events.len();
    group.throughput(Throughput::Bytes(avg_event_size as u64));

    group.bench_function("single_detection_cycle", |b| {
        b.iter(|| {
            // Add event and detect
            if let Some(event) = test_events.first() {
                let _ = detector.add_data(event.clone());
                if detector.is_ready().unwrap_or(false) {
                    let _ = detector.detect_anomaly();
                }
            }
        })
    });

    group.finish();
}

/// Benchmark JSON canonicalization
fn bench_json_canonicalization(c: &mut Criterion) {
    let mut group = c.benchmark_group("json_canonicalization");

    // Create tokenizer with only JSON canonicalization
    let config_canon_only = TokenizerConfig {
        enable_uuid: false,
        enable_hash: false,
        enable_base64: false,
        enable_jwt: false,
        enable_json_canonicalization: true,
    };
    let tokenizer_canon = Tokenizer::new(config_canon_only);

    // Create tokenizer with all patterns
    let tokenizer_full = Tokenizer::new(TokenizerConfig::default());

    // Simple JSON
    let simple_json = br#"{"z":1,"a":2,"m":3}"#;
    group.throughput(Throughput::Bytes(simple_json.len() as u64));
    group.bench_function("simple_json_canon_only", |b| {
        b.iter(|| tokenizer_canon.tokenize(black_box(simple_json)))
    });

    // Complex JSON with nested objects
    let complex_json = br#"{"user":{"profile":{"z":1,"a":2},"settings":{"theme":"dark"}},"items":[{"z":1,"a":2}]}"#;
    group.throughput(Throughput::Bytes(complex_json.len() as u64));
    group.bench_function("complex_json_canon_only", |b| {
        b.iter(|| tokenizer_canon.tokenize(black_box(complex_json)))
    });

    // Full tokenization (canonicalization + pattern replacement)
    let event_with_uuid = generate_sample_event(42);
    group.throughput(Throughput::Bytes(event_with_uuid.len() as u64));
    group.bench_function("full_tokenization", |b| {
        b.iter(|| tokenizer_full.tokenize(black_box(&event_with_uuid)))
    });

    group.finish();
}

criterion_group!(
    benches,
    bench_tokenization,
    bench_compression,
    bench_ncd_calculation,
    bench_full_detection,
    bench_json_canonicalization,
);

criterion_main!(benches);
