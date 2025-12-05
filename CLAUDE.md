# CLAUDE.md

Project context for Claude Code. For detailed information, see the linked docs.

## Quick Reference

```bash
# Start backend
cd collector-processor && go run ./cmd/driftlock-http

# Start frontend
cd landing-page && npm run dev

# Run tests
go test ./... -v

# Format code
make fmt
```

**Agent routing:** See `docs/AI_ROUTING.md`
**Launch status:** See `docs/LAUNCH_STATUS.md`
**Key files:** See `docs/KEY_FILES.md`

---

## Project Overview

**Driftlock** is a compression-based anomaly detection (CBAD) platform for OpenTelemetry data. It provides explainable anomaly detection for regulated industries through compression analysis of logs, metrics, traces, and LLM I/O.

### Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+, PostgreSQL, Kafka, Redis |
| CBAD Core | Rust (compression algorithms via FFI) |
| Frontend | Vue 3, TypeScript, Tailwind, Chart.js |
| Billing | Stripe |
| Auth | Firebase JWT + API keys |
| Hosting | GCP Cloud Run, Firebase Hosting |

### Core Components

1. **cbad-core** (Rust) - Compression-based anomaly detection algorithms
2. **collector-processor** (Go) - OpenTelemetry Collector processor
3. **api-server** (Go) - REST API with PostgreSQL, Kafka, Stripe
4. **landing-page** (Vue) - Dashboard and landing page

### Data Flow

```
OTLP Sources → Collector → CBAD Processor → API Server → UI
                                    ↓
                            PostgreSQL / Kafka
```

---

## Common Commands

### Development

```bash
make run          # Start API server
make test         # Run all tests
make build        # Build all components
make fmt          # Format code
make lint         # Run linters
```

### Docker

```bash
docker compose up -d                     # Start local stack
make docker                              # Build image
gcloud builds submit --tag gcr.io/driftlock/driftlock-http  # Deploy
```

### Database

```bash
goose -dir api/migrations postgres "$DATABASE_URL" up      # Run migrations
goose -dir api/migrations postgres "$DATABASE_URL" status  # Check status
```

---

## API Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | Liveness check |
| `GET /readyz` | Readiness check |
| `POST /v1/events` | Ingest events |
| `GET /v1/anomalies` | Retrieve anomalies |
| `POST /v1/demo/detect` | Anonymous demo (no auth) |

Full spec: `docs/architecture/api/openapi.yaml`

---

## Environment Variables

```bash
DATABASE_URL=postgresql://postgres:postgres@localhost:7543/driftlock
STRIPE_SECRET_KEY=<your-key>
STRIPE_WEBHOOK_SECRET=<your-key>
SENDGRID_API_KEY=<your-key>
FIREBASE_PROJECT_ID=driftlock
```

---

## Directory Structure

```
collector-processor/cmd/driftlock-http/  # Main API server
landing-page/src/                        # Vue frontend
api/migrations/                          # Database migrations
cbad-core/                               # Rust algorithms
docs/                                    # Documentation
.claude/agents/                          # Claude Code agents
```

---

## Agent Routing (Summary)

| Task | Agent |
|------|-------|
| Go backend | `dl-backend` |
| Vue frontend | `dl-frontend` |
| Database | `dl-db` |
| Testing | `dl-testing` |
| Deploy/infra | `dl-devops` |
| Docs | `dl-docs` |

Full matrix: `docs/AI_ROUTING.md`

---

## Linear Integration

```
"Show me issues in the driftlock project"
"Create a Linear issue: Bug - ..."
"Mark DRI-123 as Done"
```

---

## Documentation Index

- `docs/AI_ROUTING.md` - Which agent for which task
- `docs/LAUNCH_STATUS.md` - Launch checklist and progress
- `docs/KEY_FILES.md` - Important file reference
- `docs/ARCHITECTURE.md` - System design
- `docs/ALGORITHMS.md` - CBAD math foundations
- `docs/integrations/MCP_SETUP.md` - MCP server configuration

---

**Status:** ~98% launch ready
**Last updated:** 2025-12-05
