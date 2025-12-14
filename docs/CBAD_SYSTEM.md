# CBAD (Compression-Based Anomaly Detection) System

This document describes how the CBAD system works in Driftlock, including baseline management, detection flow, and configuration.

## Architecture Overview

```text
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   HTTP Client   │────▶│  driftlock-api  │────▶│   cbad-core     │
│  (curl, SDK)    │     │  (Rust/Axum)    │     │   (Rust lib)    │
└─────────────────┘     └────────┬────────┘     └─────────────────┘
                                │
                   ┌────────────┴────────────┐
                   ▼                         ▼
             ┌──────────┐              ┌──────────┐
             │ Postgres │              │  Stripe  │
             │ (state)  │              │(billing) │
             └──────────┘              └──────────┘
```

## Core Components

### 1. CBAD Rust Library (`cbad-core/`)

The core anomaly detection algorithm is implemented entirely in Rust:

- **Location**: `cbad-core/src/`
- **API integration**: Direct Rust crate dependency in `driftlock-api`
- **Key modules**:
  - `anomaly.rs` - Anomaly detection logic
  - `window.rs` - Sliding window management
  - `metrics/` - NCD, p-value, entropy calculations

Key functions:
- `Detector::new()` - Create a new detector instance
- `Detector::add_events()` - Add events to the detector
- `Detector::detect()` - Run anomaly detection
- `compute_ncd()` - Compute Normalized Compression Distance
- `compute_p_value()` - Compute statistical significance

### 2. HTTP Server (`crates/driftlock-api/`)

The main API server that exposes CBAD functionality via Axum:

- **Production endpoint**: `POST /v1/detect` (requires API key)
- **Demo endpoint**: `POST /v1/demo/detect` (rate-limited, no auth)

### 3. Baseline Persistence (PostgreSQL)

Stream anchors provide persistent baseline management:

- **Table**: `stream_anchors`
- **Features**: Anchor-based drift detection, historical tracking
- **Management**: `/v1/streams/:id/anchor` endpoints

## Detection Flow

### Demo Endpoint (`/v1/demo/detect`)

```text
1. Receive events (validated)
2. Create fresh Detector with demo settings:
   - baseline_size: 40
   - window_size: 10
   - hop_size: 5
3. Add all events to detector
4. Run detection after baseline is filled
5. Return response with metrics and anomaly details
```

### Production Endpoint (`/v1/detect`)

```text
1. Authenticate via API key
2. Resolve stream (from stream_id or default)
3. Load stream configuration and detection profile
4. Load persisted anchor/baseline (if available)
5. Create Detector with stream settings
6. Add new events and run detection
7. Persist anomalies to database
8. Return response
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | required | PostgreSQL connection string |
| `PORT` | 8080 | HTTP server port |
| `RUST_LOG` | info | Log level (debug, info, warn, error) |
| `DETECTOR_TTL_SECONDS` | 3600 | In-memory detector cache TTL |
| `DETECTOR_CLEANUP_INTERVAL_SECONDS` | 300 | Detector cleanup interval |

### Detection Parameters

| Parameter | Demo | Production | Description |
|-----------|------|------------|-------------|
| `baseline_size` | 40 | 400 | Number of events in baseline corpus |
| `window_size` | 10 | 50 | Number of events in sliding window |
| `hop_size` | 5 | 10 | Events to advance window by |
| `ncd_threshold` | 0.3 | 0.3 | NCD threshold for anomaly detection |
| `p_value_threshold` | 0.05 | 0.05 | P-value threshold for significance |

### Detection Profiles

| Profile | NCD Threshold | P-Value | Use Case |
|---------|---------------|---------|----------|
| `sensitive` | 0.20 | 0.10 | Security-critical, early warning |
| `balanced` | 0.30 | 0.05 | General purpose (default) |
| `strict` | 0.45 | 0.01 | Low noise, high confidence |
| `custom` | User-defined | User-defined | Fine-tuned settings |

## Metrics Explained

### NCD (Normalized Compression Distance)

Measures how different the window is from the baseline:
- **0.0**: Identical patterns (window compresses well with baseline)
- **1.0**: Completely different patterns
- **Threshold**: 0.3 (configurable)

### P-Value

Statistical significance of the anomaly:
- **< 0.05**: Statistically significant anomaly
- **> 0.05**: May be random variation

### Compression Ratio Change

How much the compression efficiency changed:
- **Negative**: Window is less compressible (more random/anomalous)
- **Positive**: Window is more compressible (more structured)

### Entropy Change

Change in information content:
- **Positive**: Increased randomness/complexity
- **Negative**: Decreased randomness/complexity

## Running the System

### Local Development

```bash
# Build the Rust API
cargo build -p driftlock-api --release

# Start PostgreSQL
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15

# Run the API server
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  ./target/release/driftlock-api

# Test the demo endpoint
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "normal log entry 1",
      "normal log entry 2",
      "ERROR: unusual event detected"
    ]
  }'
```

### Docker Compose

```bash
docker compose up -d
```

## API Response Example

```json
{
  "anomalies": [
    {
      "id": "anom_abc123",
      "ncd": 0.72,
      "compression_ratio": 1.41,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "explanation": "Significant deviation from baseline pattern"
    }
  ],
  "metrics": {
    "processed": 100,
    "baseline": 400,
    "window": 50,
    "duration_ms": 42
  }
}
```

## Troubleshooting

### "not enough data for analysis"

The detector needs `baseline_size + window_size` events before detection works:
- Demo: 40 + 10 = 50 events minimum
- Production: 400 + 50 = 450 events minimum

### Low anomaly detection rate

Check:
1. Baseline events represent "normal" behavior
2. Detection profile matches your use case
3. NCD threshold may need adjustment
4. Try the `sensitive` profile for more detections

### High false positive rate

Consider:
1. Increase `ncd_threshold` to reduce sensitivity
2. Decrease `p_value_threshold` for stricter significance
3. Use the `strict` profile
4. Provide more baseline data

### Driftlog (Debug Logging)

Enable detailed logging:

```bash
# Run with debug logging
RUST_LOG=debug cargo run -p driftlock-api

# Filter to specific modules
RUST_LOG=driftlock_api::routes::detection=debug cargo run -p driftlock-api

# Trace CBAD operations
RUST_LOG=cbad_core=trace cargo run -p driftlock-api
```

## CBAD Algorithm Details

For detailed information about the compression-based anomaly detection algorithm, see:

- [ALGORITHMS.md](architecture/ALGORITHMS.md) - Mathematical foundations
- `cbad-core/README.md` - Implementation details
- `cbad-core/src/anomaly.rs` - Detection logic
