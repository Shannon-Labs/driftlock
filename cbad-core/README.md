# CBAD Core (Rust)

Compression-Based Anomaly Detection primitives implemented in Rust with optional FFI bindings for Go and WASM.

This crate exposes production-facing APIs for dataset ingestion, streaming, persistence, OTLP processing, and Prometheus metrics export while keeping the compression/NCD core intact.

## Quick Start

```rust
use cbad_core::{
    CbadDetector, DetectionProfile, CsvConfig, StreamManager,
    InMemoryAnomalyStore,
};
use std::sync::Arc;

// Build a detector with the balanced profile
let detector = CbadDetector::builder()
    .with_profile(DetectionProfile::Balanced)
    .with_store(Arc::new(InMemoryAnomalyStore::default()))
    .build()?;

// Batch over an in-memory dataset
let results = detector.analyze_batch(vec![
    b"{\"msg\":\"ok\"}".to_vec(),
    b"{\"msg\":\"boom\"}".to_vec()
])?;

// Streaming from JSONL/NDJSON/CSV files
let jsonl = detector.analyze_file("data/events.jsonl")?;
let csv = detector.analyze_csv("metrics.csv", CsvConfig::default())?;

// Save/restore detector state
detector.save_state("detector.cbad")?;
let restored = CbadDetector::load_state("detector.cbad")?;

// Export Prometheus metrics
let metrics = detector.metrics();
println!("{}", metrics.to_prometheus_text());
```

## Key Capabilities

### Detection Profiles & Auto-Tuning

```rust
// Use preset profiles for different sensitivity levels
let strict_detector = CbadDetector::builder()
    .with_profile(DetectionProfile::Strict)
    .build()?;

let sensitive_detector = CbadDetector::builder()
    .with_profile(DetectionProfile::Sensitive)
    .build()?;

// Auto-tune based on feedback
detector.mark_false_positive(anomaly_id)?;  // Raises threshold
detector.confirm_anomaly(anomaly_id)?;       // Lowers threshold
```

### Multi-Stream Management

```rust
use cbad_core::{StreamManager, DetectionProfile};

let manager = StreamManager::new();

// Each stream has independent baseline/window
manager.create_stream("api-logs", DetectionProfile::Balanced, None)?;
manager.create_stream("db-metrics", DetectionProfile::Strict, None)?;

// Route events to appropriate stream
manager.ingest("api-logs", event)?;

// Cross-stream correlation
let correlated = manager.correlate_anomalies(Duration::from_secs(300))?;
```

### Prometheus Metrics Export

```rust
let metrics = detector.metrics();

// Plain text format for /metrics endpoint
let text = metrics.to_prometheus_text();

// With custom labels
let text = metrics.to_prometheus_text_with_labels(&[
    ("stream", "api-logs"),
    ("env", "production"),
]);
```

Exported metrics:
- `cbad_events_processed_total` - Total events ingested
- `cbad_anomalies_detected_total` - Total anomalies detected
- `cbad_detection_cycles_total` - Detection cycles run
- `cbad_detection_latency_seconds` - Average detection latency
- `cbad_baseline_size` / `cbad_window_size` - Current window sizes
- `cbad_tokenizer_replacements_total{type="uuid|hash|jwt|base64|json"}` - Tokenizer stats
- `cbad_tokenizer_bytes_saved_total` - Bytes saved by tokenization

### Dataset Ingestion

```rust
// CSV files
let csv_results = detector.analyze_csv("metrics.csv", CsvConfig {
    has_headers: true,
    delimiter: b',',
    column: Some(2), // Extract specific column
})?;

// JSONL/NDJSON files
let jsonl_results = detector.analyze_file("events.jsonl")?;

// Parquet files (requires `parquet` feature)
let parquet_results = detector.analyze_parquet("logs.parquet", ParquetConfig {
    column: Some("message".to_string()),
    batch_size: 1024,
})?;
```

### Async Streaming (Tokio)

```rust
// Requires `runtime-tokio` feature
use futures::StreamExt;
use tokio::io::BufReader;

let file = tokio::fs::File::open("events.jsonl").await?;
let reader = BufReader::new(file);

let mut stream = detector.stream_analyze(reader);
while let Some(result) = stream.next().await {
    let record = result?;
    if record.result.is_anomaly {
        println!("Anomaly: {}", record.result.summary);
    }
}
```

### Persistence

```rust
// Save detector state (config + window data)
detector.save_state("detector.cbad")?;

// Restore later
let restored = CbadDetector::load_state("detector.cbad")?;

// Custom storage backends
use cbad_core::storage::{AnomalyStore, AnomalyFilter};

// Query stored anomalies
let anomalies = store.query(AnomalyFilter {
    stream: Some("api-logs".to_string()),
    only_anomalies: true,
    ..Default::default()
}).await?;
```

### OTLP Integration

```rust
// Requires `otlp` feature
use cbad_core::otlp::{LogProcessor, CbadOtlpReceiver};

// Standalone OTLP gRPC receiver
let receiver = CbadOtlpReceiver::bind("0.0.0.0:4317")?
    .with_detector_config(config)?
    .on_anomaly(|anomaly| {
        println!("Detected: {}", anomaly.result.summary);
    })
    .start()
    .await?;
```

## Feature Flags

| Feature | Description |
|---------|-------------|
| `tracing` (default) | Enable `tracing` spans on hot paths |
| `runtime-tokio` | Async streaming ingestion with Tokio |
| `runtime-async-std` | Async helpers for async-std |
| `otlp` | OTLP gRPC receiver/processors (implies `runtime-tokio`) |
| `parquet` | Parquet/Arrow file ingestion |

## Storage Backends

Currently implemented:
- **InMemoryAnomalyStore** - Fast, ephemeral storage for dev/test

Planned (future releases):
- **SqliteAnomalyStore** - Single-file persistence with ACID guarantees
- **PostgresAnomalyStore** - Production-grade distributed persistence

See `storage` module documentation for implementing custom backends.

## Performance

| Operation | Latency |
|-----------|---------|
| Tokenization | 2.19μs/event |
| Compression (LZ4) | 319ns/event |
| JSON canonicalization | 398ns/event |
| Full detection cycle | 7.7ms (50 permutations) |

## Running Tests

```bash
cargo test --release              # All tests
cargo run --example test_datasets --release  # Dataset validation
cargo bench                       # Benchmarks
```

## License

Apache 2.0. Commercial licenses available from Shannon Labs.
