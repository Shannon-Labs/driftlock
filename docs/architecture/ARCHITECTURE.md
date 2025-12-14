# Architecture Overview

## Core Components

- **cbad-core** (Rust): Compression-based anomaly detection algorithms
- **driftlock-api** (Rust/Axum): REST API server with PostgreSQL, Stripe billing
- **driftlock-db** (Rust/sqlx): Database models and repository layer
- **driftlock-auth** (Rust): Firebase JWT and API key authentication
- **driftlock-billing** (Rust): Stripe integration for subscriptions
- **landing-page** (Vue 3): Dashboard and landing page

## Data Flow

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│   Client    │────▶│  Driftlock API  │────▶│  PostgreSQL │
│  (HTTP/s)   │     │   (Axum/Rust)   │     │   Database  │
└─────────────┘     └────────┬────────┘     └─────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   CBAD Core     │
                    │   (Detection)   │
                    └─────────────────┘
```

1. Clients send events via HTTP POST to `/v1/detect` or `/v1/demo/detect`
2. API server authenticates request (API key or Firebase JWT)
3. Events processed through CBAD detection engine
4. Results stored in PostgreSQL (anomalies, streams, feedback)
5. Response returned with detection results

## Crate Structure

```
driftlock/
├── Cargo.toml              # Workspace root
├── cbad-core/              # CBAD algorithms
│   └── src/
│       ├── lib.rs          # Main entry
│       ├── anomaly.rs      # Detection logic
│       ├── window.rs       # Sliding windows
│       └── metrics/        # Statistical measures
├── crates/
│   ├── driftlock-api/      # HTTP server (Axum)
│   │   └── src/
│   │       ├── main.rs     # Entry point
│   │       ├── routes/     # HTTP handlers
│   │       ├── middleware/ # Auth middleware
│   │       └── state.rs    # App state
│   ├── driftlock-db/       # Database layer
│   │   └── src/
│   │       ├── models/     # Data models
│   │       └── repos/      # Repository pattern
│   ├── driftlock-auth/     # Authentication
│   │   └── src/
│   │       ├── firebase.rs # JWT verification
│   │       └── api_key.rs  # API key auth
│   ├── driftlock-billing/  # Stripe billing
│   └── driftlock-email/    # SendGrid emails
└── landing-page/           # Vue frontend
```

## Authentication Flow

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐
│  User   │────▶│  Firebase   │────▶│  Driftlock  │
│ (Login) │     │  Auth (JWT) │     │  API        │
└─────────┘     └─────────────┘     └──────┬──────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │  API Key    │
                                    │  Generated  │
                                    └─────────────┘
```

1. User authenticates with Firebase (Google, email/password)
2. Firebase JWT sent to `/v1/auth/signup` or `/v1/auth/me`
3. Tenant created/retrieved, API key generated
4. Subsequent API calls use API key in `X-Api-Key` header

## Detection Flow

```
Events → Tokenization → Compression → NCD Calculation → Anomaly Decision
                            │
                            ▼
                     ┌─────────────┐
                     │  Baseline   │
                     │  Reference  │
                     └─────────────┘
```

1. **Tokenization**: Events normalized and tokenized
2. **Compression**: Events compressed using selected algorithm (zstd, lz4, gzip)
3. **NCD Calculation**: Normalized Compression Distance computed against baseline
4. **P-Value**: Statistical significance via permutation testing
5. **Anomaly Decision**: Threshold comparison + confidence scoring

## Stream Management

Each tenant can have multiple streams:
- Different data types (logs, metrics, traces)
- Independent baselines and detection profiles
- Auto-tuning based on feedback

## Detection Profiles

| Profile | NCD Threshold | P-Value | Use Case |
|---------|---------------|---------|----------|
| sensitive | 0.20 | 0.10 | High-security, early warning |
| balanced | 0.30 | 0.05 | General purpose (default) |
| strict | 0.45 | 0.01 | Low noise, high confidence |
| custom | User-defined | User-defined | Fine-tuned settings |

## Drift Detection (Anchors)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Baseline   │────▶│   Current   │────▶│   Drift?    │
│   Anchor    │     │   Events    │     │   (NCD)     │
└─────────────┘     └─────────────┘     └─────────────┘
```

- Anchors freeze a baseline for long-term drift detection
- Automatic reset on significant drift (configurable)
- Historical anchor tracking

## Determinism

- Deterministic seeds for reproducible results
- Configurable windows (baseline/window/hop)
- Explicit threshold configuration
- No non-deterministic concurrency in core algorithm

## Ingestion Methods

DriftLock supports multiple ingestion methods via optional features:

### HTTP REST API (Default)

Primary ingestion via HTTP endpoints:
- `POST /v1/detect`: Authenticated detection
- `POST /v1/demo/detect`: Anonymous demo (rate limited)

### Kafka Consumer (`--features kafka`)

High-throughput streaming ingestion from Kafka topics:
- Consumes from configured topic
- Backpressure via semaphore
- Stream ID from headers or JSON fields

See: [Kafka Integration](../deployment/KAFKA_INTEGRATION.md)

### OTLP gRPC Server (`--features otlp`)

Native OpenTelemetry Protocol support:
- Accepts logs, metrics, and traces
- Standard port 4317
- Stream ID from attributes

See: [OTLP Ingestion](../deployment/OTLP_INGESTION.md)

### Building with Features

```bash
# HTTP only (default)
cargo build -p driftlock-api --release

# With Kafka
cargo build -p driftlock-api --features kafka --release

# With OTLP
cargo build -p driftlock-api --features otlp --release

# All features
cargo build -p driftlock-api --features kafka,otlp,webhooks --release
```

## Scalability

Current architecture supports:
- Single-instance deployment (Replit, Cloud Run)
- PostgreSQL connection pooling
- In-memory rate limiting
- Kafka consumer for high-volume streams
- OTLP gRPC for OpenTelemetry pipelines

Horizontal scaling:
- Stateless API servers behind load balancer
- Kafka consumer groups for parallel processing
- PostgreSQL handles concurrent writes
