# Driftlock Quickstarts

These quickstarts help you try Driftlock across all surfaces: REST API, Docker, and hosted API.

## Prerequisites

- **Rust** 1.75+ (for building)
- **Docker** 24+ and Docker Compose 2.0+ (for containerized deployment)
- **Node.js** 18+ (for UI development)
- **PostgreSQL** 15+ (or use Docker)

## 1) Local Development (Rust)

Build and run the Rust API server locally:

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Build the API
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

# Health check
curl -s http://localhost:8080/healthz
```

## 2) Docker Deployment

Build and run with Docker:

```bash
# Build the image
docker build -t driftlock-api:latest -f Dockerfile .

# Run with Docker
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e STRIPE_SECRET_KEY="sk_..." \
  driftlock-api:latest

# Health check
curl -s http://localhost:8080/healthz
```

### Docker Compose (Full Stack)

```bash
# Start all services
docker compose up -d

# Health check
curl -s http://localhost:8080/healthz

# Test demo endpoint
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["normal log 1", "normal log 2", "ERROR: anomaly"]}'
```

## 3) REST API Examples

### Demo Detection (No Auth Required)

```bash
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "2025-01-01T10:00:00Z INFO Normal operation",
      "2025-01-01T10:00:01Z INFO Normal operation",
      "2025-01-01T10:00:02Z ERROR CRITICAL: Unusual pattern!"
    ]
  }'
```

### Authenticated Detection

```bash
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "my-logs",
    "events": ["event 1", "event 2", "anomalous event"]
  }'
```

### Detection with Config Override

```bash
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "my-logs",
    "events": ["..."],
    "config_override": {
      "ncd_threshold": 0.25,
      "p_value_threshold": 0.05,
      "compressor": "zstd"
    }
  }'
```

## 4) Hosted API

### Replit Deployment (Recommended)

1. Import repository into Replit
2. Set environment secrets:
   - `DATABASE_URL`
   - `STRIPE_SECRET_KEY`
   - `STRIPE_WEBHOOK_SECRET`
   - `FIREBASE_PROJECT_ID`
3. Build: `cargo build -p driftlock-api --release`
4. Run: `./target/release/driftlock-api`

### Google Cloud Run

```bash
# Build and push
gcloud builds submit --tag gcr.io/PROJECT_ID/driftlock-api

# Deploy
gcloud run deploy driftlock-api \
  --image gcr.io/PROJECT_ID/driftlock-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=postgres://..."
```

See [Deployment Guide](../../deployment/DEPLOYMENT.md) for complete instructions.

## 5) Web Dashboard

```bash
cd landing-page
npm install
npm run dev  # opens http://localhost:3000
```

Set `VITE_API_URL` to point to your API server.

## 6) Run Tests

```bash
# Run all tests
cargo test --workspace

# Run API tests only
cargo test -p driftlock-api

# Run with output
cargo test -p driftlock-api -- --nocapture

# Run CBAD algorithm tests
cargo test -p cbad-core --release
```

## Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `DATABASE_URL` | PostgreSQL connection string | required |
| `RUST_LOG` | Log level (debug, info, warn, error) | info |
| `STRIPE_SECRET_KEY` | Stripe API key | required |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | required |
| `FIREBASE_PROJECT_ID` | Firebase project ID | required |

### Detection Parameters

| Parameter | Default | Range | Description |
|-----------|---------|-------|-------------|
| `ncd_threshold` | 0.3 | 0.0-1.0 | NCD threshold for anomaly |
| `p_value_threshold` | 0.05 | 0.0-1.0 | Statistical significance |
| `baseline_size` | 400 | 50-1000 | Baseline window size |
| `window_size` | 50 | 10-500 | Analysis window size |
| `compressor` | zstd | zstd/lz4/gzip | Compression algorithm |

## Next Steps

- [Getting Started Guide](./GETTING_STARTED.md)
- [API Reference](../../architecture/API.md)
- [Detection Profiles](../guides/detection-profiles.md)
- [Testing Guide](../../development/TESTING.md)

---

**Need help?** Contact [hunter@shannonlabs.dev](mailto:hunter@shannonlabs.dev)
