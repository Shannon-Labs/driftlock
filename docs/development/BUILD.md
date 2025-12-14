# Driftlock Build and Deployment Guide

This document describes how to build the Driftlock Rust API and CBAD core.

## Prerequisites

### Rust Development
- Rust 1.75+ with stable toolchain
- Cargo with workspace support

### System Dependencies
- Linux: `build-essential`, `pkg-config`, `libssl-dev`
- macOS: Xcode command line tools
- Windows: Visual Studio Build Tools

### Database
- PostgreSQL 15+ (or Docker)

### Node.js Development
- Node.js 18+ with npm (for dashboard development)

## Quick Start

### Local Development

```bash
# Clone and setup
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Build the API server
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

# Test endpoints
curl http://localhost:8080/healthz
```

## Build Targets

### Cargo Workspace

The project uses a Cargo workspace with multiple crates:

```bash
# Build all crates
cargo build --workspace

# Build specific crates
cargo build -p driftlock-api --release    # API server
cargo build -p cbad-core --release        # CBAD algorithms
cargo build -p driftlock-db               # Database layer
cargo build -p driftlock-auth             # Authentication
cargo build -p driftlock-billing          # Stripe billing
cargo build -p driftlock-email            # Email sending
```

### Development Workflow

```bash
# Run API server (development)
cargo run -p driftlock-api

# Run with release optimizations
cargo run -p driftlock-api --release

# Run all tests
cargo test --workspace

# Run specific tests
cargo test -p driftlock-api
cargo test -p cbad-core --release

# Format code
cargo fmt --all

# Lint code
cargo clippy --all-targets -- -D warnings
```

## CBAD Core

### Building the Algorithm Library

```bash
# Build CBAD core
cargo build -p cbad-core --release

# Run CBAD tests
cargo test -p cbad-core --release

# Run benchmarks
cargo bench -p cbad-core
```

### CBAD Development

```bash
cd cbad-core

# Format and lint
cargo fmt
cargo clippy --all-targets -- -D warnings

# Test
cargo test

# Benchmark
cargo bench

# Run examples
cargo run --example crypto_stream
```

## Database Operations

### Migrations

Migrations run automatically on API server startup via sqlx.

For manual migration management:

```bash
# Install sqlx-cli
cargo install sqlx-cli

# Run migrations
sqlx migrate run --database-url "$DATABASE_URL"

# Check migration status
sqlx migrate info --database-url "$DATABASE_URL"

# Revert last migration
sqlx migrate revert --database-url "$DATABASE_URL"
```

## Testing

### Unit Tests

```bash
# All tests
cargo test --workspace

# Specific crate
cargo test -p driftlock-api
cargo test -p cbad-core --release

# With output
cargo test --workspace -- --nocapture

# Single test
cargo test -p driftlock-api test_health_check
```

### Integration Testing

```bash
# Start local stack
docker compose up -d

# Run API tests
cargo test -p driftlock-api

# Check API endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["test event"]}'
```

### Benchmarking

```bash
# CBAD benchmarks
cargo bench -p cbad-core

# Custom benchmark
cargo run --example test_datasets --release
```

## Docker Deployment

### Building Images

```bash
# Build API server image
docker build -t driftlock-api:latest -f Dockerfile .

# Build with specific tag
docker build -t driftlock-api:v1.0.0 -f Dockerfile .
```

### Running with Docker

```bash
# Run API server
docker run -d --name driftlock-api \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e STRIPE_SECRET_KEY="sk_..." \
  -e STRIPE_WEBHOOK_SECRET="whsec_..." \
  -e FIREBASE_PROJECT_ID="driftlock" \
  driftlock-api:latest
```

### Docker Compose

```bash
# Start full stack
docker compose up -d

# View logs
docker compose logs -f driftlock-api

# Stop
docker compose down
```

## Environment Configuration

### Required Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | 8080 |
| `DATABASE_URL` | PostgreSQL connection string | required |
| `STRIPE_SECRET_KEY` | Stripe API key | required |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | required |
| `FIREBASE_PROJECT_ID` | Firebase project for auth | required |

### Optional Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `RUST_LOG` | Log level (debug, info, warn, error) | info |
| `SENDGRID_API_KEY` | SendGrid API key | optional |
| `DETECTOR_CLEANUP_INTERVAL_SECONDS` | Detector cleanup interval | 300 |
| `DETECTOR_TTL_SECONDS` | Detector time-to-live | 3600 |

### Development Setup

```bash
# Copy example environment
cp .env.example .env

# Edit configuration
$EDITOR .env

# Run with environment
source .env
cargo run -p driftlock-api
```

## CI/CD Pipeline

### Local CI Validation

```bash
# Run full CI checks
cargo fmt --all --check && \
cargo clippy --all-targets -- -D warnings && \
cargo build --release && \
cargo test --workspace
```

### GitHub Actions

The repository includes workflows for:
- Continuous integration on pull requests
- Code formatting and linting checks
- Build verification
- Test execution
- Automated releases

## Troubleshooting

### Common Build Issues

1. **Rust Compilation Errors**
   ```bash
   # Update Rust toolchain
   rustup update stable

   # Clean and rebuild
   cargo clean
   cargo build -p driftlock-api --release
   ```

2. **Database Connection Errors**
   ```bash
   # Test connection
   psql "$DATABASE_URL" -c "SELECT 1"

   # Check PostgreSQL is running
   docker ps | grep postgres
   ```

3. **Docker Build Issues**
   ```bash
   # Clean Docker build cache
   docker builder prune

   # Rebuild without cache
   docker build --no-cache -t driftlock-api .
   ```

### Driftlog (Debug Logging)

```bash
# Run with debug logging
RUST_LOG=debug cargo run -p driftlock-api

# Filter to specific modules
RUST_LOG=driftlock_api=debug,driftlock_db=info cargo run -p driftlock-api

# Trace all SQL queries
RUST_LOG=sqlx=trace cargo run -p driftlock-api
```

### Performance Issues

1. **Memory Usage**
   - Monitor via `/metrics` endpoint
   - Check active detector count
   - Review detector TTL settings

2. **Throughput Issues**
   - Profile compression algorithm selection
   - Use `--release` builds for production
   - Check database connection pool settings

## Production Deployment

### Prerequisites

- Cloud provider account (GCP, AWS, Replit)
- PostgreSQL database (managed recommended)
- Domain and TLS certificates

### Deployment Options

1. **Replit (Recommended for Quick Start)**
   ```bash
   # Build
   cargo build -p driftlock-api --release

   # Run
   ./target/release/driftlock-api
   ```

2. **Google Cloud Run**
   ```bash
   # Build and push
   gcloud builds submit --tag gcr.io/PROJECT_ID/driftlock-api

   # Deploy
   gcloud run deploy driftlock-api \
     --image gcr.io/PROJECT_ID/driftlock-api \
     --platform managed \
     --region us-central1
   ```

3. **Docker Compose (Simple)**
   ```bash
   docker compose -f docker-compose.prod.yml up -d
   ```

### Health Checking

Production deployments should monitor:
- `GET /healthz` - Basic liveness check
- `GET /readyz` - Readiness with dependency validation
- `GET /metrics` - Prometheus metrics endpoint

For detailed production deployment guidance, see [DEPLOYMENT.md](../deployment/DEPLOYMENT.md).
