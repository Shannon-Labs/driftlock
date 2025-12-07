# CBAD (Compression-Based Anomaly Detection) System

This document describes how the CBAD system works in Driftlock, including baseline management, detection flow, and AI integration.

## Architecture Overview

```text
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   HTTP Client   │────▶│  driftlock-http │────▶│   cbad-core     │
│  (curl, SDK)    │     │   (Go server)   │     │  (Rust FFI)     │
└─────────────────┘     └────────┬────────┘     └─────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
              ┌──────────┐ ┌──────────┐ ┌──────────┐
              │ Postgres │ │  Redis   │ │  Ollama  │
              │ (state)  │ │(baseline)│ │   (AI)   │
              └──────────┘ └──────────┘ └──────────┘
```

## Core Components

### 1. CBAD Rust Library (`cbad-core/`)

The core anomaly detection algorithm is implemented in Rust and exposed via FFI:

- **Location**: `cbad-core/target/release/libcbad_core.a`
- **Go bindings**: `collector-processor/driftlockcbad/`
- **Build tags**: `cgo && !driftlock_no_cbad`

Key functions:
- `cbad_detector_create()` - Create a new detector instance
- `cbad_detector_add_data()` - Add events to the detector
- `cbad_detector_detect_anomaly()` - Run anomaly detection
- `cbad_compute_metrics()` - Compute NCD and p-value metrics

### 2. HTTP Server (`collector-processor/cmd/driftlock-http/`)

The main API server that exposes CBAD functionality:

- **Production endpoint**: `POST /v1/detect` (requires API key)
- **Demo endpoint**: `POST /v1/demo/detect` (rate-limited, no auth)

### 3. Baseline Persistence (Redis)

Baselines are persisted to Redis to maintain detection state across API requests:

- **Key format**: `cbad:baseline:{stream_id}`
- **TTL**: 24 hours (configurable via `BASELINE_TTL_HOURS`)
- **Format**: Newline-delimited JSON events

## Detection Flow

### Demo Endpoint (`/v1/demo/detect`)

```text
1. Receive events (max 50)
2. Create fresh Detector with demo settings:
   - baseline_size: 40
   - window_size: 10
   - hop_size: 5
3. Add all events to detector
4. Run detection after baseline is filled
5. If anomalies detected AND AI enabled:
   - Generate AI explanation via Ollama
6. Return response with metrics + AI analysis
```

### Production Endpoint (`/v1/detect`)

```text
1. Authenticate via API key
2. Resolve stream (from stream_id or default)
3. Check calibration status
4. Load persisted baseline from Redis (if available)
5. Create Detector with stream settings
6. Prime detector with persisted baseline
7. Add new events and run detection
8. Persist updated baseline to Redis (async)
9. Generate AI explanations (async)
10. Return response
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | required | PostgreSQL connection string |
| `BASELINE_REDIS_ADDR` | "" | Redis address for baseline persistence |
| `BASELINE_TTL_HOURS` | 24 | Baseline TTL in Redis |
| `AI_PROVIDER` | "" | AI provider: "ollama", "anthropic", "openai" |
| `OLLAMA_BASE_URL` | http://localhost:11434 | Ollama API endpoint |
| `OLLAMA_MODEL` | mistral | Default Ollama model |
| `DEFAULT_BASELINE` | 400 | Default baseline size |
| `DEFAULT_WINDOW` | 50 | Default window size |
| `PVALUE_THRESHOLD` | 0.05 | P-value threshold for anomaly detection |
| `NCD_THRESHOLD` | 0.3 | NCD threshold for anomaly detection |

### Detection Parameters

| Parameter | Demo | Production | Description |
|-----------|------|------------|-------------|
| `baseline_size` | 40 | 400 | Number of events in baseline corpus |
| `window_size` | 10 | 50 | Number of events in sliding window |
| `hop_size` | 5 | 10 | Events to advance window by |
| `permutation_count` | 1000 | 1000 | Permutations for p-value calculation |

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

## Running the System

### Local Development

```bash
# Terminal 1: Start dependencies
docker compose -f deploy/docker-compose.yml up -d driftlock-postgres driftlock-redis

# Terminal 2: Start Ollama (for AI)
ollama serve

# Terminal 3: Start the API server
cd collector-processor
DATABASE_URL="postgresql://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable" \
BASELINE_REDIS_ADDR="localhost:6379" \
AI_PROVIDER="ollama" \
OLLAMA_MODEL="ministral-3:3b" \
DRIFTLOCK_DEV_MODE="true" \
CGO_ENABLED=1 go run ./cmd/driftlock-http

# Terminal 4: Test the demo endpoint
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": [...]}'
```

### Docker Compose

```bash
cd deploy
docker compose up -d
```

## API Response Example

```json
{
  "success": true,
  "total_events": 50,
  "anomaly_count": 1,
  "processing_time": "178ms",
  "compression_algo": "zstd",
  "anomalies": [
    {
      "id": "uuid",
      "index": 49,
      "metrics": {
        "NCD": 0.499,
        "PValue": 0.51,
        "BaselineCompressionRatio": 6.6,
        "WindowCompressionRatio": 2.9,
        "ConfidenceLevel": 0.49,
        "CompressionRatioChange": -0.56
      },
      "event": {"level": "ERROR", "msg": "..."},
      "why": "CBAD Analysis: NCD=0.499..."
    }
  ],
  "ai_analysis": {
    "provider": "ollama",
    "model": "ministral-3:3b",
    "explanation": "This anomaly indicates...",
    "latency": "10.5s"
  }
}
```

## Troubleshooting

### "cbad: Rust core not available"

The stub implementation is being used. Ensure:
1. CGO is enabled: `CGO_ENABLED=1`
2. The Rust library is built: `cd cbad-core && cargo build --release`
3. No `driftlock_no_cbad` build tag

### "not enough data for analysis"

The detector needs `baseline_size + window_size` events before detection works:
- Demo: 40 + 10 = 50 events minimum
- Production: 400 + 50 = 450 events minimum

### "baseline redis unavailable"

Redis is optional. Without it, baselines are not persisted across requests.
Set `BASELINE_REDIS_ADDR` to enable persistence.

### AI analysis unavailable

Check:
1. `AI_PROVIDER` is set (e.g., "ollama")
2. Ollama is running: `curl http://localhost:11434/api/tags`
3. Model is available: `ollama pull ministral-3:3b`
