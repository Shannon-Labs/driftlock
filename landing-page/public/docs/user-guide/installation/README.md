# Installation Guide

This guide will help you install and set up Driftlock in various environments.

## SDK Installation (Recommended)

The fastest way to get started with Driftlock is using an official SDK in your application:

### Node.js / TypeScript

```bash
npm install @driftlock/node
```

**Quick Start:**

```typescript
import { DriftlockClient } from '@driftlock/node';

const client = new DriftlockClient({ apiKey: 'dlk_your-key' });
const result = await client.detect([
  { message: 'normal event' },
  { message: 'unusual event' }
]);
```

See [Node.js SDK documentation](../../sdk/nodejs.md) for full reference.

### Python

```bash
pip install driftlock
```

**Quick Start:**

```python
import asyncio
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient(api_key='dlk_your-key')
    result = await client.detect([
        {'message': 'normal event'},
        {'message': 'unusual event'}
    ])

asyncio.run(main())
```

See [Python SDK documentation](../../sdk/python.md) for full reference.

### Other Languages

Use the [REST API](../api/rest-api.md) directly with your preferred HTTP client library.

---

## Local Development (Rust / Self-Hosted)

If you want to build and run Driftlock from source, follow the instructions below.

## Prerequisites

### System Requirements

- **Operating System**: Linux, macOS, or Windows (WSL2)
- **Memory**: Minimum 4GB RAM, recommended 8GB+
- **Storage**: Minimum 10GB free space
- **Network**: Internet connection for dependencies

### Software Requirements

- **Rust**: 1.75 or later
- **Node.js**: 18 or later (for UI development)
- **Docker**: 24+ and Docker Compose 2.0+ (optional)
- **Git**: 2.30 or later

## Quick Start (Docker Compose)

The fastest way to get Driftlock running is with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Copy environment template and configure
cp .env.example .env
# Edit .env with your configuration

# Start all services
docker compose up -d

# Access the API
curl http://localhost:8080/healthz
```

### Services Started

- **API Server**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## Local Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
```

### 2. Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit environment variables
nano .env
```

Key environment variables:

```bash
# Database Configuration
DATABASE_URL=postgres://driftlock:driftlock@localhost:5432/driftlock

# API Configuration
PORT=8080
RUST_LOG=info

# Authentication
FIREBASE_PROJECT_ID=your-firebase-project

# Billing (if needed)
STRIPE_SECRET_KEY=sk_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Email (optional)
SENDGRID_API_KEY=SG...
```

### 3. Build the Project

```bash
# Build the Rust API
cargo build -p driftlock-api --release

# Build CBAD core (if modifying algorithms)
cargo build -p cbad-core --release
```

### 4. Database Setup

#### Using Docker (Recommended)

```bash
# Start PostgreSQL
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15
```

#### Migrations

Migrations run automatically on startup, or manually:

```bash
# Install sqlx-cli
cargo install sqlx-cli

# Run migrations
sqlx migrate run --database-url "$DATABASE_URL"
```

### 5. Start the API Server

```bash
# Run the API server
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  cargo run -p driftlock-api --release

# Or run the release binary
./target/release/driftlock-api
```

Server will be available at:
- API: http://localhost:8080
- Health check: http://localhost:8080/healthz
- Metrics: http://localhost:8080/metrics

### 6. Start Frontend (Optional)

```bash
cd landing-page
npm install
npm run dev
```

Dashboard available at http://localhost:3000

## Production Deployment

### Replit (Recommended for Quick Start)

1. Import repository into Replit
2. Set environment secrets
3. Build: `cargo build -p driftlock-api --release`
4. Run: `./target/release/driftlock-api`

### Docker Deployment

```bash
# Build the image
docker build -t driftlock-api -f Dockerfile .

# Run the container
docker run -d --name driftlock-api \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e STRIPE_SECRET_KEY="sk_..." \
  driftlock-api
```

### Google Cloud Run

```bash
# Build and push
gcloud builds submit --tag gcr.io/PROJECT_ID/driftlock-api

# Deploy
gcloud run deploy driftlock-api \
  --image gcr.io/PROJECT_ID/driftlock-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## Configuration Options

### API Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `DATABASE_URL` | - | PostgreSQL connection string |
| `RUST_LOG` | info | Log level (debug, info, warn, error) |
| `FIREBASE_PROJECT_ID` | - | Firebase project for auth |
| `STRIPE_SECRET_KEY` | - | Stripe API key |
| `STRIPE_WEBHOOK_SECRET` | - | Stripe webhook signing secret |

### Detection Configuration

| Parameter | Default | Range | Description |
|-----------|---------|-------|-------------|
| `ncd_threshold` | 0.3 | 0.0-1.0 | Normalized compression distance threshold |
| `p_value_threshold` | 0.05 | 0.0-1.0 | Statistical significance threshold |
| `baseline_size` | 400 | 50-1000 | Baseline window size |
| `window_size` | 50 | 10-500 | Analysis window size |
| `hop_size` | 10 | 1-100 | Window hop size |

## Verification

### Health Checks

```bash
# API Server Health
curl http://localhost:8080/healthz

# Readiness
curl http://localhost:8080/readyz

# Prometheus Metrics
curl http://localhost:8080/metrics
```

### Test Detection

```bash
# Test demo endpoint (no auth required)
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

## Troubleshooting

### Common Issues

#### Build Errors

```bash
# Clean and rebuild
cargo clean
cargo build -p driftlock-api --release

# Check Rust version
rustc --version  # Should be 1.75+
```

#### Database Connection Issues

```bash
# Test connection
psql "$DATABASE_URL" -c "SELECT 1"

# Check PostgreSQL is running
docker ps | grep postgres
```

#### Port Conflicts

```bash
# Check what's using port 8080
lsof -i :8080

# Kill process if needed
kill -9 <PID>
```

### Getting Help

- **Documentation**: [docs/](../docs/)
- **Issues**: [GitHub Issues](https://github.com/Shannon-Labs/driftlock/issues)
- **Support**: support@driftlock.io

## Next Steps

After installation:

1. [Getting Started Guide](../getting-started/GETTING_STARTED.md)
2. [API Reference](../../architecture/API.md)
3. [Detection Profiles](../guides/detection-profiles.md)
4. [Deployment Guide](../../deployment/DEPLOYMENT.md)
