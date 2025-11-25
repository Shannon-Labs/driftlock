# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Available Claude Models

This repository is configured to work with the following Claude models:
- **Claude Opus 4.5** - Most capable model for complex reasoning and code generation
- **Claude Sonnet 4.5** - Balanced performance and speed (default)
- **Claude Haiku 4.5** - Fastest model for quick tasks and simple queries

## Project Overview

**Driftlock** is a compression-based anomaly detection (CBAD) platform for OpenTelemetry data, powered by Meta's OpenZL format-aware compression framework. It provides explainable anomaly detection for regulated industries through advanced compression analysis of logs, metrics, traces, and LLM I/O.

## High-Level Architecture

### Core Components

The system follows a streaming architecture with these main components:

1. **cbad-core (Rust)** - Compression-based anomaly detection algorithms with FFI bindings
   - Implements compression ratio analysis, Normalized Compression Distance (NCD), permutation testing
   - Provides `libcbad_core.a` static library for Go integration
   - Supports OpenZL format-aware compression alongside zstd/lz4/gzip

2. **collector-processor (Go)** - OpenTelemetry Collector processor
   - Integrates with cbad-core via CGO
   - Processes OTLP logs, metrics, traces in real-time
   - Computes CBAD metrics and flags anomalies

3. **api-server (Go)** - REST API service
   - Handles event ingestion, anomaly retrieval, evidence export
   - Integrates with PostgreSQL for storage, Kafka for streaming
   - Provides SSE/WebSocket for real-time anomaly streaming
   - Built-in OpenTelemetry instrumentation

4. **llm-receivers (Go)** - OpenTelemetry Collector receivers
   - Specialized receivers for LLM prompts, responses, tool calls
   - Enables AI/ML system monitoring

5. **exporters (Go)** - Evidence bundle generation
   - Creates JSON/PDF reports for regulatory compliance
   - DORA, NIS2, and AI Act compliant evidence bundles

6. **ui (Next.js)** - Minimal dashboard
   - Browse anomalies, streams, and artifacts
   - Real-time monitoring interface

### Data Flow

```
OTLP Sources → Collector Receivers → driftlock_cbad Processor → API Server → UI
                     ↓
                Evidence Bundles (exporters)
                     ↓
                Kafka Stream → Anomaly Storage (PostgreSQL)
```

### Technology Stack

- **Go 1.22+** for API server, collectors, and exporters
- **Rust 1.70+** for CBAD core algorithms
- **PostgreSQL** for primary storage
- **Kafka** for streaming anomalies
- **Redis** for state management
- **OpenTelemetry** for observability
- **Next.js** for UI

## Common Commands

### Development Workflow

```bash
# Start API server in development mode (port 8080)
make run

# Run all tests
make test

# Build all components
make build

# Clean build artifacts
make clean
```

### Building Specific Components

```bash
# Build API server binary
make api

# Build OpenTelemetry collector with CBAD integration
make collector

# Build Rust CBAD static library
make cbad-core-lib

# Build development tools (synthetic data generator)
make tools

# Run full CI validation locally
make ci-check
```

### Code Quality

```bash
# Format Go and Rust code
make fmt

# Run linters
make lint

# Update dependencies
make tidy
```

### Testing

```bash
# Run all Go tests with verbose output
go test ./... -v

# Run Rust tests
cd cbad-core && cargo test

# Run Go benchmarks
go test -bench=. ./...

# Run Rust benchmarks
cd cbad-core && cargo bench

# Run end-to-end tests
go test ./tests/e2e/... -v
```

### Docker & Deployment

```bash
# Build Docker image
make docker

# Start local development stack
docker compose -f deploy/docker-compose.yml up

# Start with observability
docker compose -f deploy/docker-compose.yml \
  -f deploy/docker-compose.observability.yml up

# Build release binaries for multiple platforms
make release
```

## Key Configuration

### Environment Variables

- `PORT` - API server port (default: 8080)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTEL endpoint for telemetry export
- `OTEL_SERVICE_NAME` - Service name for tracing (default: driftlock-api)
- `OTEL_ENV` - Environment name (default: dev)
- `DRIFTLOCK_VERSION` - Version string override
- `CBAD_CORE_PROFILE` - Rust build profile: release or dev

### API Endpoints

- `GET /healthz` - Liveness check
- `GET /readyz` - Readiness check
- `GET /v1/version` - Build/version information
- `POST /v1/events` - Ingest JSON event payloads
- `GET /v1/anomalies` - Retrieve anomalies
- `GET /v1/evidence` - Export evidence bundles

## Core Algorithms & Math

### Compression-Based Anomaly Detection (CBAD)

The system uses several mathematical approaches:

1. **Compression Ratio**: `CR = compressed_bytes / raw_bytes`
2. **Delta Bits**: `Δ_bits = ((C_B - C_W) × 8) / |W|`
3. **Shannon Entropy**: `H = -Σ p_i log₂ p_i`
4. **Normalized Compression Distance**: `NCD = (C_{BW} - min(C_B, C_W)) / max(C_B, C_W)`

### Statistical Significance

Permutation testing with deterministic seeds (ChaCha20 RNG) for p-value calculation:
```
p-value = (1 + #{permutations where metric ≥ observed}) / (1 + total_permutations)
```

### Deterministic Rules

- All compression operations use deterministic seeds
- Window sizes (baseline/window/hop) configured explicitly
- Avoid non-deterministic concurrency in core algorithms
- Regex-based PII detection with configurable policies

## Code Organization

### Main Directories

- `api-server/` - Go API service (main entry point: `api-server/cmd/driftlock-api/main.go`)
- `cbad-core/` - Rust crate implementing CBAD algorithms
- `collector-processor/` - OpenTelemetry Collector processor
- `llm-receivers/` - Collector receivers for LLM I/O
- `exporters/` - Evidence bundle exporters
- `pkg/version/` - Version information
- `docs/` - Comprehensive documentation
- `deploy/` - Docker Compose and Kubernetes manifests

### API Server Structure

```
api-server/
├── cmd/driftlock-api/       # Entry point
├── internal/
│   ├── api/                 # HTTP handlers and routing
│   ├── auth/                # Authentication/authorization
│   ├── billing/             # Stripe integration
│   ├── cbad/                # CBAD integration layer
│   ├── compression/         # Compression adapters
│   ├── handlers/            # Business logic handlers
│   ├── middleware/          # HTTP middleware (logging, rate limiting)
│   ├── models/              # Data models
│   ├── services/            # Business services
│   ├── storage/             # Database abstraction (PostgreSQL, Redis)
│   ├── streaming/           # Kafka integration
│   └── telemetry/           # OpenTelemetry setup
└── migrations/              # Database migrations
```

### Key Files

- `Makefile` - Build and development commands
- `go.mod` - Go dependencies and module replacements
- `.env.example` - Environment variable template
- `docs/BUILD.md` - Detailed build instructions
- `docs/ARCHITECTURE.md` - Architecture overview
- `docs/ALGORITHMS.md` - Mathematical foundations

## Development Guidelines

### Coding Standards

**Go Standards:**
- Target Go 1.22+
- Follow OpenTelemetry Collector conventions
- No `panic` - return errors with context
- Enforce gofmt, goimports, staticcheck
- Table-driven tests with deterministic seeds

**Rust Standards:**
- Edition 2021
- No `unsafe` without documented justification
- Use `cargo fmt` and `cargo clippy -D warnings`
- Property tests for compression math
- Deterministic RNG seeds

**Testing Requirements:**
- ≥80% line coverage on core math packages
- Deterministic seeds in benchmarks
- Edge cases and error paths covered

### Build Tags

- `-tags driftlock_cbad_cgo` - Enable CGO integration with CBAD core

### Important Integration Points

1. **Rust-Go FFI Boundary**
   - Located in `collector-processor/`
   - CGO linking with `libcbad_core.a`
   - Use `make collector` to build

2. **OpenTelemetry Integration**
   - All components instrumented with OTEL
   - Trace context propagation across boundaries
   - Configurable exporters via environment variables

3. **Storage Layer**
   - PostgreSQL for persistent storage
   - Redis for ephemeral state
   - Kafka for streaming anomalies
   - Tiered storage pattern in `api-server/internal/storage/`

## Testing & Debugging

### Local Testing

```bash
# Start local stack
docker compose -f deploy/docker-compose.yml up

# Generate synthetic test data
go run ./tools/synthetic

# Check endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/version
```

### Debugging Tips

- Set `RUST_LOG=debug` for detailed Rust logging
- Use `go run -tags driftlock_cbad_cgo` for testing FFI integration
- Monitor `/metrics` endpoint for Prometheus metrics
- Check `backend.log`, `firebase-debug.log`, `pipeline.log` for diagnostics

### Performance Issues

- Verify CBAD window sizing configuration
- Monitor Go heap and Rust allocations
- Check for memory leaks at FFI boundary
- Profile compression algorithm selection

## Regulatory Compliance

The system includes enterprise-ready compliance features:

- **DORA Compliance** - Digital Operational Resilience Act evidence bundles
- **NIS2 Compliance** - EU cybersecurity incident reporting
- **AI Act Compliance** - Runtime AI monitoring for LLM/ML systems
- **Audit Trails** - Cryptographically signed evidence packages

Evidence bundles generated in `exporters/` directory with JSON/PDF formats.

## OpenZL Integration

**Key Competitive Advantage:** Meta's OpenZL format-aware compression framework provides:
- 1.5-2x better compression ratios than zstd on structured data
- 20-40% faster compression/decompression
- Format-aware transforms (struct-of-arrays, delta encoding, tokenization)
- Glass-box compression with explainable anomaly detection

Located at `deps/openzl/` as a git submodule. Provides deterministic training and embedded decode recipes.

## Database Schema

Initial schema in `api-server/migrations/`:
- `api-server/migrations/001_initial_schema.up.sql` - Creates tables
- `api-server/internal/storage/migrations/` - Migration management

Primary entities:
- Tenants (multi-tenant architecture)
- Anomalies (detected anomalies with metadata)
- Evidence bundles (compliance reports)
- Telemetry streams (raw data references)

## Deployment

### Options

1. **Kubernetes with Helm**
   ```bash
   helm install driftlock ./deploy/helm/driftlock
   ```

2. **Docker Compose**
   ```bash
   docker compose -f deploy/docker-compose.yml up
   ```

3. **Manual Binary Deployment**
   ```bash
   make release
   scp bin/driftlock-api-* target-host:/opt/driftlock/
   ```

### Production Readiness

- Health checks: `/healthz`, `/readyz`, `/metrics`
- Prometheus metrics endpoint
- CBAD performance metrics
- Anomaly detection rate monitoring

## Documentation

Comprehensive documentation in `docs/`:
- `ARCHITECTURE.md` - System design
- `ALGORITHMS.md` - Mathematical foundations
- `BUILD.md` - Build and deployment guide
- `CODING_STANDARDS.md` - Development guidelines
- `CONTRIBUTING.md` - Contribution workflow
- `API.md` - API reference
- `DEPLOYMENT.md` - Production deployment guide

See also phase summaries and roadmaps for implementation progress tracking.

---

# LAUNCH DEVELOPMENT GUIDE

> **For AI Assistants:** Start here! This section tracks what's done, what's next, and how to continue development.

## Current Status: ~95% Launch Ready

**Last Updated:** 2025-11-25
**Target:** Public launch with self-serve signup, working billing, demo/playground, API access

### What's Done
- All backend features implemented and tested
- All frontend features implemented
- Documentation updated
- E2E tests written
- Security audit passed (no SQL injection, all endpoints auth-protected)

## Quick Start for New AI Sessions

```
Hey Claude! Check the LAUNCH DEVELOPMENT GUIDE section in CLAUDE.md for:
1. Current progress checklist
2. What to work on next
3. Key files and context
```

---

## LAUNCH CHECKLIST

### Phase 1: User Onboarding (COMPLETE)

- [x] Database migration for verification flow (`api/migrations/20251124000000_launch_enhancements.sql`)
- [x] Email verification backend (`collector-processor/cmd/driftlock-http/onboarding.go`)
- [x] Verify endpoint implementation (`GET /v1/onboard/verify?token=xxx`)
- [x] API key regeneration endpoint (`POST /v1/me/keys/regenerate`)
- [x] API key create endpoint (`POST /v1/me/keys/create`)
- [x] API key revoke endpoint (`POST /v1/me/keys/revoke`)
- [x] SignupForm.vue - pending verification state
- [x] VerifyEmailView.vue - verification landing page
- [x] Router update for `/verify` route

### Phase 2: Stripe Billing (COMPLETE)

- [x] **Add 14-day trial to checkout** (`collector-processor/cmd/driftlock-http/billing.go`)
  - Add `TrialPeriodDays: 14` to checkout session params
  - Track `trial_ends_at` in tenant record

- [x] **Complete webhook handlers** (`billing.go`)
  - [x] `customer.subscription.trial_will_end` - send reminder email
  - [x] `invoice.payment_failed` - enter grace period
  - [x] `invoice.payment_succeeded` - clear grace flags

- [x] **Grace period logic** (in `billing.go` and `db.go`)
  - 7-day grace period after payment failure
  - Auto-downgrade to free tier after grace expiry

- [x] **Billing status endpoint** (`GET /v1/me/billing`)
  - Returns: status, plan, trial_days_remaining, grace_period_ends_at

- [x] **Frontend billing UI** (`DashboardView.vue`)
  - Escalating trial countdown banners (relaxed/reminder/urgent)
  - Grace period warning
  - Free tier upgrade prompt
  - Success toast after Stripe checkout
  - Error handling with retry

- [x] **Pricing page checkout** (`HomeView.vue`)
  - Connect pricing buttons to Stripe checkout
  - Handle unauthenticated users (scroll to signup)
  - Loading states and error toast

### Phase 3: Developer Experience (COMPLETE)

- [x] **Anonymous demo endpoint** (new file: `demo.go`)
  - `POST /v1/demo/detect` - no auth required
  - Rate limit: 10 requests/min per IP
  - Max 50 events per request
  - Response includes signup CTA

- [x] **Playground demo mode** (`PlaygroundShell.vue`)
  - Use demo endpoint when not authenticated
  - Show conversion banner
  - Demo mode indicator in header

- [x] **Usage dashboard** (`DashboardView.vue`)
  - Daily usage chart (Chart.js via vue-chartjs)
  - Stream breakdown table with anomaly rates
  - `GET /v1/me/usage/details` endpoint

### Phase 4: Polish & Testing (MOSTLY COMPLETE)

- [x] E2E test: signup → verify → detect → anomaly (`e2e_onboarding_test.go`)
- [x] E2E test: trial → checkout → subscription (`e2e_billing_test.go`)
- [x] Error code reference page (`landing-page/public/docs/user-guide/api/errors.md`)
- [ ] Load test Cloud Run deployment (optional pre-launch)
- [ ] Production deployment verification

---

## KEY FILES REFERENCE

### Backend (Go)
| File | Purpose |
|------|---------|
| `collector-processor/cmd/driftlock-http/main.go` | Route registration, server setup |
| `collector-processor/cmd/driftlock-http/onboarding.go` | Signup/verify handlers |
| `collector-processor/cmd/driftlock-http/billing.go` | Stripe integration |
| `collector-processor/cmd/driftlock-http/dashboard.go` | User dashboard endpoints |
| `collector-processor/cmd/driftlock-http/db.go` | Database operations |
| `collector-processor/cmd/driftlock-http/store_auth_ext.go` | API key management |
| `collector-processor/cmd/driftlock-http/email.go` | SendGrid email service |
| `collector-processor/cmd/driftlock-http/demo.go` | Anonymous demo endpoint with rate limiting |

### Frontend (Vue)
| File | Purpose |
|------|---------|
| `landing-page/src/components/cta/SignupForm.vue` | Main signup form |
| `landing-page/src/views/VerifyEmailView.vue` | Email verification page |
| `landing-page/src/views/DashboardView.vue` | User dashboard with usage chart |
| `landing-page/src/views/HomeView.vue` | Landing page with pricing |
| `landing-page/src/components/playground/PlaygroundShell.vue` | Interactive playground (demo mode support) |
| `landing-page/src/components/dashboard/UsageChart.vue` | Daily usage chart component |
| `landing-page/src/router/index.ts` | Route definitions |

### Database
| File | Purpose |
|------|---------|
| `api/migrations/20250301120000_initial_schema.sql` | Core tables |
| `api/migrations/20250302000000_onboarding.sql` | Onboarding fields |
| `api/migrations/20251119000000_add_stripe_fields.sql` | Stripe columns |
| `api/migrations/20251124000000_launch_enhancements.sql` | Verification & billing |

---

## MCP SERVERS TO CONFIGURE

For easier development, configure these MCP servers in your Claude Code settings:

### 1. PostgreSQL MCP Server
```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres"],
      "env": {
        "DATABASE_URL": "postgresql://postgres:postgres@localhost:7543/driftlock"
      }
    }
  }
}
```
**Use for:** Direct database queries, schema inspection, data verification

### 2. Filesystem MCP Server
```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/Volumes/VIXinSSD/driftlock"]
    }
  }
}
```
**Use for:** File operations, bulk reads

### 3. GitHub MCP Server
```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "<your-token>"
      }
    }
  }
}
```
**Use for:** PR creation, issue management, code review

### 4. Stripe MCP Server (if available)
For Stripe testing and webhook simulation

### 5. Firebase MCP Server (if available)
For Firebase Auth testing and user management

---

## PRICING & COST ANALYSIS

### Current Pricing Tiers

| Tier | Price | Events/Month | Features |
|------|-------|--------------|----------|
| **Pilot** (Free) | $0 | 10,000 | Basic detection, 14-day retention |
| **Radar** | $20/mo | 500,000 | Email alerts, 30-day retention |
| **Lock** | $200/mo | 5,000,000 | DORA/NIS2 evidence, priority support |
| **Orbit** | Custom | Unlimited | Dedicated support, SLA |

### Infrastructure Costs (Estimated Monthly)

| Service | Provider | Est. Cost | Notes |
|---------|----------|-----------|-------|
| **Compute** | Cloud Run | $50-200 | Auto-scaling, pay-per-use |
| **Database** | Cloud SQL (Postgres) | $30-100 | Depends on instance size |
| **Email** | SendGrid | $15-50 | Based on volume |
| **Payments** | Stripe | 2.9% + $0.30 | Per transaction |
| **DNS/CDN** | Cloudflare | $0-20 | Free tier available |
| **Monitoring** | Built-in OTEL | $0 | Self-instrumented |
| **Total** | | **~$100-400/mo** | Before revenue |

### Unit Economics

| Metric | Calculation |
|--------|-------------|
| **Gross Margin (Radar)** | $20 - ~$2 infra = 90% |
| **Gross Margin (Lock)** | $200 - ~$10 infra = 95% |
| **Break-even** | ~5-10 Radar customers or 1-2 Lock customers |
| **Target CAC** | < $50 (content marketing, SEO, word-of-mouth) |

### Cost Optimization Opportunities

1. **Reserved instances** - 30-50% savings on Cloud SQL
2. **Committed use discounts** - Cloud Run pricing
3. **Batch processing** - Reduce real-time compute costs
4. **Tiered storage** - Move old data to cheaper storage
5. **OpenZL compression** - Reduce storage costs by 50%+

---

## DEVELOPMENT WORKFLOW

### Starting a New Session

1. **Check this file first** - See what's done and what's next
2. **Run health check:**
   ```bash
   curl http://localhost:8080/healthz
   ```
3. **Start local stack (if needed):**
   ```bash
   docker compose -f docker-compose.yml up -d
   ```
4. **Pick next unchecked item** from the checklist above
5. **Update checklist** when done (edit this file!)

### Before Committing

```bash
# Format code
make fmt

# Run tests
make test

# Check for issues
make lint
```

### Key Environment Variables

```bash
# Required for local dev
DATABASE_URL=postgresql://postgres:postgres@localhost:7543/driftlock
SENDGRID_API_KEY=<your-key>  # Or leave empty for mock emails
STRIPE_SECRET_KEY=<your-key>
STRIPE_WEBHOOK_SECRET=<your-key>
FIREBASE_PROJECT_ID=driftlock

# Optional
LOG_LEVEL=debug
DRIFTLOCK_DEV_MODE=true
```

---

## NEXT PRIORITY TASKS

**If you're an AI picking this up, work on these in order:**

### Pre-Launch Verification
1. **Run E2E tests** - `go test ./collector-processor/cmd/driftlock-http/... -v`
2. **Test production deployment** - Verify Cloud Run is working
3. **Verify Stripe webhooks** - Test with `stripe listen`

### Optional Enhancements
4. **Load testing** - Performance verification under load
5. **Recent anomalies feed** - Populate dashboard with real data
6. **Redis rate limiting** - Upgrade demo endpoint for multi-instance

### All Major Features Complete
- [x] User onboarding with email verification
- [x] Stripe billing with trials and grace periods
- [x] Anonymous demo endpoint with rate limiting
- [x] Dashboard with usage charts
- [x] Playground demo mode for unauthenticated users
- [x] API documentation
- [x] E2E tests for onboarding and billing
- [x] Security audit passed

---

## USEFUL COMMANDS

```bash
# Local development
docker compose up -d                    # Start Postgres
cd collector-processor && go run ./cmd/driftlock-http  # Run API

# Frontend development
cd landing-page && npm run dev          # Start Vue dev server

# Database
goose -dir api/migrations postgres "$DATABASE_URL" up  # Run migrations
goose -dir api/migrations postgres "$DATABASE_URL" status  # Check status

# Testing
curl -X POST http://localhost:8080/api/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","company_name":"Test Co","plan":"trial"}'

# Check Stripe webhooks locally
stripe listen --forward-to localhost:8080/api/v1/billing/webhook
```

---

**Last updated by:** Claude (Opus 4.5)
**Session:** Codebase audit, cleanup, and launch preparation (2025-11-25)
