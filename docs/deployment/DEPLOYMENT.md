# Driftlock Deployment Guide

## Quick Start with Docker

The fastest way to get started:

```bash
# Build the Rust API image
docker build -t driftlock-api -f Dockerfile .

# Run with Docker
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e STRIPE_SECRET_KEY="sk_..." \
  driftlock-api

# Verify it's running
curl http://localhost:8080/healthz
```

## Prerequisites

- Rust 1.75+ (for building)
- PostgreSQL 15+
- Docker 24+ (for containerized deployment)
- Node.js 18+ (for UI/playground)

## Local Development

### 1. Database Setup

```bash
# Start PostgreSQL with Docker
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15

# Set database URL
export DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable"
```

### 2. Build and Run API Server

```bash
# Build release binary
cargo build -p driftlock-api --release

# Run server
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  ./target/release/driftlock-api
```

### 3. UI Development

```bash
cd landing-page
npm install
npm run dev
```

Access UI at: http://localhost:3000

## Environment Variables

### API Server

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | 8080 |
| `DATABASE_URL` | PostgreSQL connection string | required |
| `RUST_LOG` | Log level (debug, info, warn, error) | info |
| `STRIPE_SECRET_KEY` | Stripe API key | required |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | required |
| `FIREBASE_PROJECT_ID` | Firebase project ID | required |
| `SENDGRID_API_KEY` | SendGrid API key | optional |
| `DETECTOR_CLEANUP_INTERVAL_SECONDS` | Detector cleanup interval | 300 |
| `DETECTOR_TTL_SECONDS` | Detector time-to-live | 3600 |

### UI

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | API server URL | http://localhost:8080 |

## Production Deployment

### Replit Deployment

Recommended for quick deployment:

1. Create a new Replit from GitHub
2. Set environment secrets in Replit dashboard:
   - `DATABASE_URL`
   - `STRIPE_SECRET_KEY`
   - `STRIPE_WEBHOOK_SECRET`
   - `FIREBASE_PROJECT_ID`
3. Build command: `cargo build -p driftlock-api --release`
4. Run command: `./target/release/driftlock-api`

### Docker Deployment

```bash
# Build image
docker build -t driftlock-api:latest -f Dockerfile .

# Run container
docker run -d --name driftlock-api \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e STRIPE_SECRET_KEY="sk_..." \
  -e STRIPE_WEBHOOK_SECRET="whsec_..." \
  -e FIREBASE_PROJECT_ID="driftlock" \
  driftlock-api:latest
```

### Google Cloud Run

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT_ID/driftlock-api

# Deploy to Cloud Run
gcloud run deploy driftlock-api \
  --image gcr.io/PROJECT_ID/driftlock-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=postgres://..." \
  --set-secrets "STRIPE_SECRET_KEY=stripe-key:latest"
```

### Kubernetes with Helm

```bash
# Install Helm chart
helm install driftlock ./helm/driftlock \
  --set database.host=postgres.default.svc.cluster.local \
  --set database.password=SECURE_PASSWORD \
  --set api.replicas=3
```

## Database Migrations

Migrations are embedded in the binary and run automatically on startup via sqlx.

For manual migration management:

```bash
# Install sqlx-cli
cargo install sqlx-cli

# Run migrations
sqlx migrate run --database-url "$DATABASE_URL"

# Check migration status
sqlx migrate info --database-url "$DATABASE_URL"
```

## Monitoring

### Prometheus Metrics

API server exposes metrics at `/metrics`:

- `driftlock_http_requests_total` - Total HTTP requests
- `driftlock_events_processed_total` - Events processed
- `driftlock_anomalies_detected_total` - Anomalies detected
- `driftlock_detectors_active` - Active detector count

### Health Checks

- `GET /healthz` - Liveness probe (returns 200 if server is running)
- `GET /readyz` - Readiness probe (checks database connection)

## Troubleshooting

### Database Connection Issues

```bash
# Test connection
psql "$DATABASE_URL"

# Check PostgreSQL logs
docker logs driftlock-postgres
```

### Build Issues

```bash
# Clean build
cargo clean
cargo build -p driftlock-api --release

# Check Rust version
rustc --version  # Should be 1.75+
```

### High Memory Usage

- Check connection pool settings
- Monitor active detector count via `/metrics`
- Review detector TTL settings

### Slow Queries

- Check PostgreSQL indexes
- Enable query logging
- Use EXPLAIN ANALYZE for slow queries

## Security Checklist

- [ ] Use strong database passwords
- [ ] Enable TLS for database connections
- [ ] Use HTTPS in production
- [ ] Rotate API keys regularly
- [ ] Set up Stripe webhook signature verification
- [ ] Configure CORS appropriately
- [ ] Regular security audits
- [ ] Keep dependencies updated

## Backup and Recovery

### Database Backups

```bash
# Backup
pg_dump "$DATABASE_URL" > backup.sql

# Restore
psql "$DATABASE_URL" < backup.sql
```

### Disaster Recovery

- Daily automated backups
- Test restore procedures monthly
- Maintain backup retention for 30 days
- Document recovery procedures

## Performance Tuning

### Database

```sql
-- Check connection count
SELECT count(*) FROM pg_stat_activity;

-- Analyze slow queries
SET log_min_duration_statement = 1000;
```

### API Server

- Connection pool default: 10 connections
- Increase for high-traffic deployments
- Monitor via Prometheus metrics

## Support

For issues, contact: support@driftlock.io
Documentation: https://docs.driftlock.io
