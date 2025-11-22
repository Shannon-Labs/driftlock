# Driftlock Deployment Guide

## Quick Start with Docker Compose

The fastest way to get started is with the unified `docker-compose.yml`:

```bash
# Start HTTP API server
docker compose up -d driftlock-http

# Verify it's running
curl http://localhost:8080/healthz

# Start with Kafka collector (optional)
docker compose --profile kafka up -d
```

See [QUICKSTARTS.md](../QUICKSTARTS.md) for more examples.

## Prerequisites

- Go 1.22 or later
- PostgreSQL 15 or later (optional, for full API features)
- Docker 24+ (recommended for containerized deployment)
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

# Run migrations
export DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable"
cd api-server/migrations
migrate -path . -database "$DATABASE_URL" up
```

### 2. API Server

```bash
# Build
cd api-server
go build -o bin/driftlock-api ./cmd/driftlock-api

# Run
export DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable"
export PORT=8080
./bin/driftlock-api
```

### 3. UI Development

```bash
cd ui
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
| `LOG_LEVEL` | Log level (debug, info, warn, error) | info |
| `LOG_FORMAT` | Log format (json, text) | json |
| `RATE_LIMIT_RPS` | Rate limit requests per second | 100 |
| `RATE_LIMIT_BURST` | Rate limit burst size | 200 |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OpenTelemetry collector endpoint | - |

### UI

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | API server URL | http://localhost:8080 |
| `NODE_ENV` | Environment (development, production) | development |

## Production Deployment

### Docker Compose (Unified Setup)

The unified `docker-compose.yml` includes both APIs:

```bash
# Start HTTP API server (default)
docker compose up -d driftlock-http

# Start with Kafka collector (optional)
docker compose --profile kafka up -d

# View logs
docker compose logs -f driftlock-http

# Stop services
docker compose down
```

### Individual Docker Builds

```bash
# Build HTTP API server image
docker build -t driftlock-http:latest -f collector-processor/cmd/driftlock-http/Dockerfile .

# Build Kafka collector image
docker build -t driftlock-collector:latest -f collector-processor/cmd/driftlock-collector/Dockerfile .

# Run HTTP API
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e CORS_ALLOW_ORIGINS=https://play.driftlock.net \
  driftlock-http:latest
```

### Kubernetes with Helm

```bash
# Install Helm chart
helm install driftlock ./helm/driftlock \
  --set database.host=postgres.default.svc.cluster.local \
  --set database.password=SECURE_PASSWORD \
  --set api.replicas=3
```

## Monitoring

### Prometheus Metrics

API server exposes metrics at `/metrics`:

- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request duration histogram
- `anomalies_detected_total` - Total anomalies detected
- `anomalies_by_stream_type` - Anomalies by stream type

### Grafana Dashboard

Import the dashboard from `monitoring/grafana/dashboard.json`

## Troubleshooting

### Database Connection Issues

```bash
# Test connection
psql "$DATABASE_URL"

# Check migrations
migrate -path api-server/migrations -database "$DATABASE_URL" version
```

### High Memory Usage

- Check connection pool settings
- Monitor goroutine leaks
- Review SSE client connections

### Slow Queries

- Check PostgreSQL indexes
- Enable query logging: `SET log_min_duration_statement = 1000;`
- Use EXPLAIN ANALYZE for slow queries

## Security Checklist

- [ ] Use strong database passwords
- [ ] Enable TLS for database connections
- [ ] Use HTTPS in production
- [ ] Rotate API keys regularly
- [ ] Enable rate limiting
- [ ] Set up firewall rules
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
-- Increase connection pool
ALTER SYSTEM SET max_connections = 200;

-- Tune work_mem for queries
ALTER SYSTEM SET work_mem = '16MB';

-- Enable query plan caching
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
```

### API Server

- Increase `GOMAXPROCS` for high-traffic deployments
- Use connection pooling (default: 100 max, 10 idle)
- Enable HTTP/2
- Use CDN for static assets

## Support

For issues, contact: support@driftlock.io
Documentation: https://docs.driftlock.io
