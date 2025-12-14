# Launch Development Guide

> **For AI Assistants:** Start here! This section tracks what's done, what's next, and how to continue development.

## Current Status: 100% Rust Migration Complete

**Last Updated:** 2025-12-11
**Target:** Replit deployment with pure Rust backend

### What's Done
- All backend features migrated from Go to Rust (Axum/Tokio)
- All 28+ API routes implemented
- Database layer complete (sqlx with PostgreSQL)
- Authentication: Firebase JWT + API keys
- Stripe billing integration
- Prometheus metrics endpoint
- Detection profiles and auto-tuning
- Stream anchors for drift detection
- Rate limiting (IP-based for demo endpoint)

### Tech Stack (Rust)
- **Framework:** Axum 0.7 with Tower middleware
- **Runtime:** Tokio async runtime
- **Database:** sqlx with PostgreSQL
- **Auth:** Firebase JWT (jsonwebtoken), API keys (Argon2)
- **Billing:** stripe-rust
- **Metrics:** metrics + metrics-exporter-prometheus

### Production URLs
- **API**: Replit deployment (pending)
- **Website**: https://driftlock.net

---

## Project Management

**Primary**: [Linear - shannon-labs/driftlock](https://linear.app/shannon-labs/project/driftlock-a8c80503816c/overview)
**Automation Level**: Heavy (auto-triage, daily standups, velocity reports)

### Linear MCP Commands

```
"Show me issues in the driftlock project"
"Create a Linear issue: Bug - API returns 500 on empty payload"
"Mark DRI-123 as Done"
"What issues are blocked or need attention?"
```

### Workflow

1. Check Linear for current priorities
2. Pick an issue and move to "In Progress"
3. Create branch: `git checkout -b DRI-123-description`
4. Code & commit - PRs auto-link to Linear
5. Merge - Issue auto-closes when PR merges

---

## API Routes Summary

### Health & Metrics (Public)
- [x] `GET /healthz` - Liveness check
- [x] `GET /readyz` - Readiness check
- [x] `GET /v1/version` - Version info
- [x] `GET /metrics` - Prometheus metrics

### Demo (IP Rate Limited)
- [x] `POST /v1/demo/detect` - Anonymous demo detection
- [x] `POST /v1/waitlist` - Email capture

### Authentication (Firebase)
- [x] `POST /v1/auth/signup` - Firebase signup
- [x] `GET /v1/auth/me` - Get current user

### Detection (API Key Required)
- [x] `POST /v1/detect` - Authenticated detection
- [x] `GET /v1/anomalies` - List anomalies
- [x] `GET /v1/anomalies/:id` - Get anomaly details
- [x] `POST /v1/anomalies/:id/feedback` - Submit feedback

### Streams (API Key Required)
- [x] `GET /v1/streams` - List streams
- [x] `POST /v1/streams` - Create stream
- [x] `GET /v1/streams/:id` - Get stream
- [x] `GET /v1/streams/:id/profile` - Get detection profile
- [x] `PATCH /v1/streams/:id/profile` - Update profile
- [x] `GET /v1/streams/:id/tuning` - Tuning history

### Anchors (API Key Required)
- [x] `GET /v1/streams/:id/anchor` - Get anchor settings
- [x] `GET /v1/streams/:id/anchor/details` - Full anchor data
- [x] `POST /v1/streams/:id/reset-anchor` - Create new anchor
- [x] `DELETE /v1/streams/:id/anchor` - Deactivate anchor

### Account (API Key Required)
- [x] `GET /v1/account` - Get account info
- [x] `PATCH /v1/account` - Update account
- [x] `GET /v1/account/usage` - Usage summary

### API Keys (API Key Required)
- [x] `GET /v1/api-keys` - List keys
- [x] `POST /v1/api-keys` - Create key
- [x] `DELETE /v1/api-keys/:id` - Revoke key
- [x] `POST /v1/api-keys/:id/regenerate` - Regenerate key

### Billing (API Key Required)
- [x] `POST /v1/billing/checkout` - Create checkout session
- [x] `POST /v1/billing/portal` - Create portal session
- [x] `POST /v1/billing/webhook` - Stripe webhooks
- [x] `GET /v1/me/billing` - Billing status
- [x] `GET /v1/me/usage/details` - Daily usage
- [x] `GET /v1/me/usage/ai` - AI usage (stub)
- [x] `GET /v1/me/ai/config` - AI config (stub)

### Profiles (API Key Required)
- [x] `GET /v1/profiles` - List detection profiles

---

## Useful Commands

```bash
# Build and run Rust API
cargo build -p driftlock-api --release
cargo run -p driftlock-api --release

# Run tests
cargo test -p driftlock-api

# Frontend
cd landing-page && npm run dev

# Database (sqlx)
sqlx migrate run --database-url "$DATABASE_URL"

# Testing detection
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["test event 1", "test event 2"]}'

# Stripe webhooks (local)
stripe listen --forward-to localhost:8080/v1/billing/webhook
```

---

## Pricing Tiers

| Tier | Price | Events/Month |
|------|-------|--------------|
| Pulse (Free) | $0 | 10,000 |
| Radar | $15/mo | 500,000 |
| Tensor | $100/mo | 5,000,000 |
| Orbit | $499/mo | Unlimited |

---

## Crate Structure

```
driftlock/
├── Cargo.toml           # Workspace root
├── cbad-core/           # CBAD algorithms
├── crates/
│   ├── driftlock-api/   # Axum HTTP server
│   ├── driftlock-db/    # sqlx models/repos
│   ├── driftlock-auth/  # Firebase + API keys
│   ├── driftlock-billing/ # Stripe integration
│   └── driftlock-email/ # SendGrid service
├── landing-page/        # Vue frontend
└── archive/go-backend/  # Legacy Go (reference)
```

---

**Last updated:** 2025-12-11
