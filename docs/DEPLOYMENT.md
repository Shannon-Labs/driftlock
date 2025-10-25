# Driftlock Phase 4: Deployment Guide

This guide covers deploying Driftlock API server, UI, and OpenTelemetry Collector to production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Docker** 20.10+ (for containerized deployment)
- **Kubernetes** 1.25+ (for production deployment)
- **Helm** 3.10+ (for Kubernetes package management)
- **PostgreSQL** 15+ (database)
- **Go** 1.22+ (for local development)
- **Rust** 1.75+ (for CBAD core builds)

### System Requirements

**Minimum (Development)**:
- 4 CPU cores
- 8 GB RAM
- 50 GB disk space

**Recommended (Production)**:
- 16 CPU cores
- 32 GB RAM
- 500 GB SSD storage (database)
- 100 GB SSD storage (application)

## Local Development

### 1. Setup PostgreSQL

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name driftlock-postgres \
  -e POSTGRES_PASSWORD=dev \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_DB=driftlock \
  -p 5432:5432 \
  postgres:15

# Apply migrations
psql -h localhost -U driftlock -d driftlock -f api-server/internal/storage/migrations/001_initial_schema.sql
```

### 2. Build CBAD Core

```bash
cd cbad-core
cargo build --release
```

### 3. Run API Server

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_DATABASE=driftlock
export DB_USER=driftlock
export DB_PASSWORD=dev
export PORT=8080

# Run server
cd api-server
go run cmd/api-server/main.go
```

### 4. Test API

```bash
# Health check
curl http://localhost:8080/healthz

# Get version
curl http://localhost:8080/v1/version

# List anomalies (requires auth)
curl -H "Authorization: Bearer dev-key-12345" \
  http://localhost:8080/v1/anomalies
```

## Docker Deployment

### 1. Build Images

```bash
# Build API server
docker build -f deploy/docker/Dockerfile.api-server -t driftlock/api-server:latest .

# Build OTel Collector
docker build -f deploy/docker/Dockerfile.collector -t driftlock/collector:latest .

# Build UI
docker build -f deploy/docker/Dockerfile.ui -t driftlock/ui:latest .
```

### 2. Run with Docker Compose

```bash
# Create docker-compose.yml
cat > docker-compose.yml <<EOF
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: driftlock
      POSTGRES_PASSWORD: \${DB_PASSWORD:-changeme}
      POSTGRES_DB: driftlock
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U driftlock"]
      interval: 10s
      timeout: 5s
      retries: 5

  api-server:
    image: driftlock/api-server:latest
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_DATABASE: driftlock
      DB_USER: driftlock
      DB_PASSWORD: \${DB_PASSWORD:-changeme}
      PORT: 8080
    ports:
      - "8080:8080"
      - "9090:9090"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3

  ui:
    image: driftlock/ui:latest
    depends_on:
      - api-server
    environment:
      NEXT_PUBLIC_API_URL: http://api-server:8080
    ports:
      - "3000:3000"

  collector:
    image: driftlock/collector:latest
    depends_on:
      - api-server
    ports:
      - "4317:4317"  # OTLP gRPC
      - "4318:4318"  # OTLP HTTP
      - "8888:8888"  # Metrics
    volumes:
      - ./collector-config.yaml:/etc/otelcol/config.yaml

volumes:
  postgres-data:
EOF

# Start services
docker-compose up -d

# View logs
docker-compose logs -f api-server
```

## Kubernetes Deployment

### 1. Prerequisites

```bash
# Add Helm repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create namespace
kubectl create namespace driftlock
```

### 2. Create Secrets

```bash
# Database password
kubectl create secret generic driftlock-db-credentials \
  --from-literal=password='YOUR_STRONG_PASSWORD' \
  -n driftlock

# API keys (optional)
kubectl create secret generic driftlock-api-keys \
  --from-literal=admin-key='YOUR_ADMIN_API_KEY' \
  -n driftlock
```

### 3. Configure values.yaml

```bash
# Create custom values
cat > driftlock-values.yaml <<EOF
apiServer:
  replicaCount: 3
  image:
    repository: your-registry/driftlock/api-server
    tag: "1.0.0"

  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10

postgresql:
  auth:
    username: driftlock
    existingSecret: driftlock-db-credentials
    database: driftlock

  primary:
    persistence:
      size: 100Gi
      storageClass: fast-ssd

ui:
  ingress:
    enabled: true
    hosts:
      - host: driftlock.your-company.com
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: driftlock-tls
        hosts:
          - driftlock.your-company.com
EOF
```

### 4. Deploy with Helm

```bash
# Install Driftlock
helm install driftlock ./deploy/helm/driftlock \
  -f driftlock-values.yaml \
  -n driftlock

# Verify deployment
kubectl get pods -n driftlock
kubectl get svc -n driftlock

# Check logs
kubectl logs -f deployment/driftlock-api-server -n driftlock
```

### 5. Apply Database Migrations

```bash
# Create migration job
kubectl create job --from=cronjob/driftlock-migrations migration-001 -n driftlock

# Or run manually
kubectl run -it --rm driftlock-migrate \
  --image=postgres:15 \
  --restart=Never \
  -n driftlock \
  -- psql -h driftlock-postgresql -U driftlock -d driftlock -f /migrations/001_initial_schema.sql
```

## Configuration

### Environment Variables

**API Server**:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_DATABASE` | Database name | `driftlock` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | (required) |
| `DB_MAX_CONNECTIONS` | Max DB connections | `100` |
| `CBAD_NCD_THRESHOLD` | NCD detection threshold | `0.3` |
| `CBAD_P_VALUE_THRESHOLD` | P-value threshold | `0.05` |
| `AUTH_TYPE` | Auth type (apikey/oidc) | `apikey` |
| `LOG_LEVEL` | Log level | `info` |

### Tuning CBAD Detection

Adjust thresholds via API:

```bash
curl -X PATCH http://localhost:8080/v1/config \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "ncd_threshold": 0.25,
    "p_value_threshold": 0.01,
    "baseline_size": 150,
    "window_size": 75
  }'
```

## Monitoring

### Prometheus Metrics

Access metrics endpoint:
```bash
curl http://localhost:9090/metrics
```

Key metrics:
- `driftlock_anomalies_detected_total` - Total anomalies detected
- `driftlock_events_processed_total` - Total events processed
- `driftlock_api_request_duration_seconds` - API latency
- `driftlock_active_sse_connections` - Active SSE clients
- `driftlock_database_connections` - DB connection pool stats

### Grafana Dashboards

Import dashboards:

```bash
# Anomaly Detection dashboard
kubectl apply -f deploy/grafana/dashboards/anomaly-detection.json

# API Performance dashboard
kubectl apply -f deploy/grafana/dashboards/api-performance.json
```

### Health Checks

```bash
# Liveness probe
curl http://localhost:8080/healthz

# Readiness probe (checks DB connectivity)
curl http://localhost:8080/readyz
```

## Troubleshooting

### API Server Won't Start

**Check database connectivity**:
```bash
kubectl exec -it deployment/driftlock-api-server -n driftlock -- \
  sh -c 'echo "SELECT 1" | psql -h $DB_HOST -U $DB_USER'
```

**Check CBAD library**:
```bash
# Verify libcbad_core.a is linked
kubectl logs deployment/driftlock-api-server -n driftlock | grep CBAD
```

### High API Latency

**Check database queries**:
```sql
SELECT * FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Scale API server**:
```bash
kubectl scale deployment driftlock-api-server --replicas=6 -n driftlock
```

### SSE Connections Dropping

**Check connection limits**:
```bash
# Increase max connections
kubectl set env deployment/driftlock-api-server \
  MAX_SSE_CONNECTIONS=2000 \
  -n driftlock
```

### Database Running Out of Space

**Cleanup old metrics**:
```sql
-- Delete performance metrics older than 30 days
DELETE FROM performance_metrics
WHERE created_at < NOW() - INTERVAL '30 days';
```

## Production Checklist

- [ ] Database backups configured (daily)
- [ ] TLS certificates installed
- [ ] API keys rotated and secured
- [ ] Prometheus alerting rules configured
- [ ] Log aggregation setup (ELK/Loki)
- [ ] Resource limits configured
- [ ] Autoscaling tested under load
- [ ] Disaster recovery plan documented
- [ ] Security scanning enabled (Trivy/Snyk)
- [ ] Compliance evidence export tested

## Next Steps

1. **Load Testing**: Use tools like k6 or Locust to test 10k+ events/sec
2. **Security Hardening**: Enable OIDC, rotate keys, audit logs
3. **Compliance**: Configure evidence bundle exports for DORA/NIS2
4. **Optimization**: Tune PostgreSQL indexes for your query patterns

For support: https://github.com/your-org/driftlock/issues
