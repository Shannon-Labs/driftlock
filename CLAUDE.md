# CLAUDE.md

Project context for Claude Code. For detailed information, see the linked docs.

## Quick Reference

```bash
# Start backend (Rust)
cargo run -p driftlock-api --release

# Start frontend
cd landing-page && npm run dev

# Run tests
cargo test -p driftlock-api

# Build release
cargo build -p driftlock-api --release
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
| Backend | Rust (Axum, Tokio, sqlx), PostgreSQL |
| CBAD Core | Rust (cbad-core crate) |
| Frontend | Vue 3, TypeScript, Tailwind, Chart.js |
| Billing | Stripe |
| Auth | Firebase JWT + API keys |
| Hosting | Replit / GCP Cloud Run, Firebase Hosting |

### Core Components

1. **cbad-core** (Rust) - Compression-based anomaly detection algorithms
2. **driftlock-api** (Rust) - REST API with Axum, PostgreSQL, Stripe
3. **driftlock-db** (Rust) - Database models and repositories
4. **driftlock-auth** (Rust) - Firebase JWT and API key authentication
5. **landing-page** (Vue) - Dashboard and landing page

### Data Flow

```
HTTP REST ─────────┐
                   │
Kafka Topic ───────┼──▶ API Server ──▶ CBAD Core ──▶ Detection ──▶ UI
(--features kafka) │        │
                   │        ▼
OTLP gRPC ─────────┘   PostgreSQL
(--features otlp)
```

---

## Common Commands

### Development

```bash
cargo run -p driftlock-api              # Start API server (debug)
cargo run -p driftlock-api --release    # Start API server (optimized)
cargo test                              # Run all tests
cargo build --release                   # Build all components
cargo fmt                               # Format code
cargo clippy                            # Run linter
```

### Optional Features

```bash
# With Kafka consumer
cargo build -p driftlock-api --features kafka --release

# With OTLP gRPC server (OpenTelemetry)
cargo build -p driftlock-api --features otlp --release

# All features
cargo build -p driftlock-api --features kafka,otlp,webhooks --release
```

See: `docs/deployment/KAFKA_INTEGRATION.md`, `docs/deployment/OTLP_INGESTION.md`

### Docker

```bash
docker compose up -d                     # Start local stack
docker build -t driftlock-api -f Dockerfile .  # Build image
```

### Database

```bash
# Migrations are run automatically by sqlx on startup
# Or run migrations via sqlx-cli:
sqlx migrate run --database-url "$DATABASE_URL"
```

---

## API Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | Liveness check |
| `GET /readyz` | Readiness check |
| `POST /v1/detect` | Submit events for detection |
| `POST /v1/demo/detect` | Anonymous demo (no auth) |
| `GET /v1/anomalies` | Retrieve anomalies |
| `POST /v1/anomalies/{id}/feedback` | Mark false positive/confirm |
| `GET /v1/streams/{id}/profile` | Get detection profile |
| `PATCH /v1/streams/{id}/profile` | Update profile settings |
| `GET /v1/streams/{id}/tuning` | Get tuning history |
| `GET /v1/profiles` | List detection profiles |

### Adaptive Features

- **Detection Profiles:** sensitive, balanced, strict, custom
- **Auto-Tuning:** Automatic threshold adjustment based on feedback
- **Adaptive Windowing:** Automatic window sizing based on stream characteristics

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
crates/driftlock-api/src/                # Rust API server (Axum)
crates/driftlock-db/src/                 # Database models & repos (sqlx)
crates/driftlock-auth/src/               # Authentication (Firebase + API keys)
crates/driftlock-billing/src/            # Stripe billing integration
crates/driftlock-email/src/              # Email service (SendGrid)
cbad-core/src/                           # CBAD Rust algorithms
landing-page/src/                        # Vue frontend
archive/go-backend/                      # Legacy Go backend (reference only)
docs/                                    # Documentation
.claude/agents/                          # Claude Code agents
```

---

## Agent Routing (Summary)

| Task | Agent |
|------|-------|
| Rust backend | `dl-backend` |
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

**Status:** 100% Rust API ready for Replit deployment
**Last updated:** 2025-12-11
