# Installation Guide

This guide will help you install and set up DriftLock in various environments.

## Prerequisites

### System Requirements

- **Operating System**: Linux, macOS, or Windows (WSL2)
- **Memory**: Minimum 4GB RAM, recommended 8GB+
- **Storage**: Minimum 10GB free space
- **Network**: Internet connection for dependencies

### Software Requirements

- **Go**: 1.24 or later
- **Rust**: 1.70 or later
- **Node.js**: 18 or later
- **Docker**: 20.10+ and Docker Compose 2.0+
- **Git**: 2.30 or later

## Quick Start (Docker Compose)

The fastest way to get DriftLock running is with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/shannon-labs/driftlock.git
cd driftlock

# Copy environment template and configure
cp .env.example .env
# Edit .env to set your API key: DEFAULT_API_KEY=your-secret-key

# Start all services
docker compose up -d

# Access the dashboard
open http://localhost:3000

# Log in with your API key
```

### Services Started

- **API Server**: http://localhost:8080
- **Dashboard**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379 (optional)
- **Kafka**: localhost:9092 (optional, for streaming)

**Note**: DriftLock now runs standalone without external dependencies like Supabase. All functionality is available through the REST API and dashboard.

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

Key environment variables to configure:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=driftlock
DB_USER=postgres
DB_PASSWORD=your_secure_password_here

# API Configuration
PORT=8080
LOG_LEVEL=info

# Authentication (required for dashboard access)
AUTH_TYPE=apikey
DEFAULT_API_KEY=your_api_key_here_for_dashboard_access
DEFAULT_ORG_ID=default

# Development API Key (optional)
DRIFTLOCK_DEV_API_KEY=

# Optional: Supabase (only needed for advanced compliance features)
# Leave empty for standalone OSS deployment
SUPABASE_PROJECT_ID=
SUPABASE_ANON_KEY=
SUPABASE_SERVICE_ROLE_KEY=
SUPABASE_BASE_URL=
```

**Note**: For OSS deployments, set `DEFAULT_API_KEY` (and optionally override `DRIFTLOCK_DEV_API_KEY`) to access the dashboard.

### 3. Install Dependencies

Use the automated setup:

```bash
make setup
```

Or install manually:

#### Rust Dependencies

```bash
cd cbad-core
cargo build --release
```

#### Go Dependencies

```bash
cd api-server
go mod download
go build -o driftlock-api ./cmd/api-server
```

#### Node.js Dependencies

```bash
cd web-frontend
npm install
```

### 4. Database Setup

#### Using Docker (Recommended)

```bash
# Start PostgreSQL
docker compose up -d postgres

# Run migrations (once database is ready)
make migrate
```

#### Local PostgreSQL

```bash
# Create database
createdb driftlock

# Run migrations
make migrate
```

The migration tool will:
- Create the `schema_migrations` tracking table
- Apply pending migrations in order
- Rollback migrations if needed (`make migrate-down`)

### 5. Start Services

```bash
# Start all services with Docker Compose
make dev

# Or start individually:
cd api-server && go run ./cmd/api-server &
cd web-frontend && npm run dev &
```

Services will be available at:
- API Server: http://localhost:8080
- Dashboard: http://localhost:3000
- Health check: http://localhost:8080/healthz

## OpenTelemetry Collector Setup

### 1. Download Collector

```bash
# Download OpenTelemetry Collector Contrib
curl -sSL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.88.0/otelcol-contrib_0.88.0_linux_amd64.tar.gz | tar xz
```

### 2. Configuration

Create `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  driftlock/anomaly:
    endpoint: http://localhost:8080
    thresholds:
      compression_ratio: 0.7
      ncd_threshold: 0.3
    explanation:
      enabled: true
      detail_level: "detailed"

  batch:

exporters:
  otlp:
    endpoint: http://localhost:8080
    tls:
      insecure: true

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [driftlock/anomaly, batch]
      exporters: [otlp]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]

    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

### 3. Start Collector

```bash
./otelcol-contrib --config=otel-collector-config.yaml
```

## Kubernetes Deployment

### 1. Create Namespace

```bash
kubectl create namespace driftlock
```

### 2. Deploy Components

```bash
# Deploy database
kubectl apply -f k8s/postgres.yaml

# Deploy API server
kubectl apply -f k8s/api-server.yaml

# Deploy dashboard
kubectl apply -f k8s/dashboard.yaml

# Deploy OpenTelemetry Collector
kubectl apply -f k8s/otel-collector.yaml
```

### 3. Verify Deployment

```bash
kubectl get pods -n driftlock
kubectl port-forward -n driftlock svc/driftlock-api 8080:8080
kubectl port-forward -n driftlock svc/driftlock-dashboard 3000:3000
```

## Helm Installation

### 1. Add Repository

```bash
helm repo add shannon-labs https://charts.shannonlabs.ai
helm repo update
```

### 2. Install Chart

```bash
helm install driftlock shannon-labs/driftlock \
  --namespace driftlock \
  --create-namespace \
  --set apiServer.replicas=2 \
  --set dashboard.replicas=1
```

### 3. Custom Values

Create `values.yaml`:

```yaml
apiServer:
  replicas: 3
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"

dashboard:
  replicas: 2
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "200m"

postgresql:
  enabled: true
  auth:
    password: "your-secure-password"
```

Install with custom values:

```bash
helm install driftlock Shannon-Labs/driftlock \
  --namespace driftlock \
  --create-namespace \
  -f values.yaml
```

## Configuration Options

### API Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `LOG_LEVEL` | info | Log level (debug, info, warn, error) |
| `CBAD_NCD_THRESHOLD` | 0.3 | NCD threshold for anomaly detection |
| `CBAD_P_VALUE_THRESHOLD` | 0.05 | P-value threshold |

### CBAD Algorithm Configuration

| Parameter | Default | Range | Description |
|-----------|---------|-------|-------------|
| `compression_ratio_threshold` | 0.7 | 0.0-1.0 | Compression ratio threshold |
| `ncd_threshold` | 0.3 | 0.0-1.0 | Normalized compression distance |
| `baseline_size` | 100 | 50-1000 | Baseline window size |
| `window_size` | 50 | 10-500 | Analysis window size |
| `hop_size` | 10 | 1-100 | Window hop size |

## Verification

### Health Checks

```bash
# API Server Health
curl http://localhost:8080/healthz

# Dashboard Access
curl http://localhost:3000

# Database Connection
curl http://localhost:8080/readyz
```

### Test Anomaly Detection

```bash
# Send test data
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "log",
    "data": {
      "message": "This is a test log entry",
      "level": "INFO"
    }
  }'

# Check for anomalies
curl http://localhost:8080/v1/anomalies
```

## Troubleshooting

### Common Issues

#### Port Conflicts

```bash
# Check what's using ports
lsof -i :8080
lsof -i :3000

# Kill processes if needed
kill -9 <PID>
```

#### Database Connection Issues

```bash
# Check PostgreSQL status
docker compose ps postgres

# View logs
docker compose logs postgres

# Test connection
psql -h localhost -U postgres -d driftlock
```

#### Build Errors

```bash
# Clean and rebuild
make clean
make build

# Update dependencies
cd api-server && go mod tidy
cd web-frontend && npm install
cd cbad-core && cargo update
```

#### Permission Issues

```bash
# Fix Docker permissions
sudo chown -R $USER:$USER .

# Fix Go module permissions
chmod +x api-server/driftlock-api
```

### Getting Help

- **Documentation**: [docs/](../docs/)
- **Issues**: [GitHub Issues](https://github.com/shannon-labs/driftlock/issues)
- **Discussions**: [GitHub Discussions](https://github.com/shannon-labs/driftlock/discussions)
- **Security**: security@shannonlabs.ai

## Next Steps

After installation:

1. [Configure OpenTelemetry Collector](../docs/otel-collector.md)
2. [Review API Documentation](../docs/api-reference.md)
3. [Check Examples](../examples/)
4. [Set up Compliance Integration](../docs/compliance.md)
