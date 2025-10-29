Driftlock

A compression-based anomaly detection (CBAD) platform for OpenTelemetry data. Driftlock provides explainable anomaly detection for regulated industries through advanced compression analysis of logs, metrics, traces, and LLM I/O.

## Architecture

Driftlock consists of two main services:

1. **Go API Server** (`api-server/`): Core anomaly detection pipeline
   - Processes events and detects anomalies using CBAD algorithm
   - Exposes REST API for anomaly management
   - Integrates with Supabase for web-frontend data synchronization and usage metering

2. **React Web Frontend** (`web-frontend/`): Customer dashboard and billing
   - Built with React/TypeScript and shadcn-ui
   - Uses Supabase as backend (PostgreSQL + Edge Functions)
   - Handles user authentication, billing, and anomaly visualization

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local development)
- Supabase account with project created

### Environment Setup

1. Copy environment template:
   ```bash
   cp .env.example .env
   ```

2. Update Supabase configuration in `.env`:
   ```bash
   SUPABASE_PROJECT_ID=your_actual_project_id
   SUPABASE_ANON_KEY=your_actual_anon_key
   SUPABASE_SERVICE_ROLE_KEY=your_actual_service_role_key
   SUPABASE_BASE_URL=https://your_project_id.supabase.co
   SUPABASE_WEBHOOK_URL=https://your_project_id.supabase.co/functions/v1/webhook
   ```

### Database Migration (local Postgres)

If you have a local Postgres or a Supabase connection string available in `.env`, you can apply the bundled SQL schema:

```bash
make migrate
```

### Running with Docker Compose

1. Start all services:
   ```bash
   ./start.sh
   ```

2. Services will be available at:
   - Web Frontend: http://localhost:3000
   - Go API: http://localhost:8080
   - API Documentation: http://localhost:8080/healthz

### Local Development

#### Web Frontend
```bash
cd web-frontend
npm install
npm run dev
```

#### Go API Server
```bash
# Quick run
make run

# Or run directly
cd api-server
go run ./cmd/driftlock-api
```

Health and readiness:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Ingest an event (optionally using API key auth):

```bash
# Optional: set DEFAULT_API_KEY and DEFAULT_ORG_ID in .env to enable API key auth
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${DEFAULT_API_KEY:-testkey}" \
  -d '{
    "organization_id": "org_123",  
    "event_type": "log",
    "data": {"message": "Unusual pattern detected", "level": "WARN"}
  }'
```

SSE stream of anomalies:

```bash
curl -N http://localhost:8080/v1/stream/anomalies
```

## Integration Details

### Data Flow

1. **Event Ingestion**: Events are sent to Go API `/v1/events` endpoint
2. **Anomaly Detection**: Go API processes events using CBAD algorithm
3. **Data Synchronization**: Anomalies are stored in PostgreSQL and synced to Supabase
4. **Real-time Updates**: Web frontend receives updates via Supabase subscriptions
5. **Billing**: Usage is tracked in Supabase for subscription management

### API Endpoints (Go API)

```
GET  /healthz                         # Liveness
GET  /readyz                          # Readiness (DB ping; Supabase best-effort)
GET  /v1/version                      # Version
POST /v1/events                       # Ingest events (optionally API key auth)

GET  /v1/anomalies                    # List anomalies
POST /v1/anomalies                    # Create anomaly
GET  /v1/anomalies/{id}               # Get anomaly by ID
PATCH /v1/anomalies/{id}/status       # Update anomaly status

GET  /v1/stream/anomalies             # Server-Sent Events stream
```

Auth for ingestion (optional):
- Set `DEFAULT_API_KEY` and `DEFAULT_ORG_ID` in `.env` to enable API key auth for `/v1/events`. The key’s organization determines metering and sync context.

### Supabase Integration

The Go API server integrates with Supabase through:

- **Anomaly Synchronization**: Anomalies created in Go API are also stored in Supabase via REST
- **Status Updates**: Anomaly status changes are synchronized between systems
- **Usage Tracking**: API usage is tracked via Supabase Edge Function `meter-usage`
- **Webhook Notifications**: Go API can trigger Supabase Edge Functions

## Documentation

- [Integration Guide](INTEGRATION_README.md) - Detailed setup and integration instructions (legacy)
- [API Documentation](docs/API.md) - REST API reference
- [Architecture](docs/ARCHITECTURE.md) - System design and components
 - [CLOUDFLARE_DEPLOYMENT.md](CLOUDFLARE_DEPLOYMENT.md) - Production deployment
 - [README_CLOUDFLARE.md](README_CLOUDFLARE.md) - Cloudflare quick start
