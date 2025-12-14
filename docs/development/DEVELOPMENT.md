# Driftlock Development Guide

This document provides a comprehensive guide for developers working on the Driftlock project.

## Table of Contents
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Development Environment Setup](#development-environment-setup)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Contribution Workflow](#contribution-workflow)
- [Project Structure](#project-structure)

## Project Overview

Driftlock is an explainable, deterministic anomaly detection toolkit designed for regulated and audit-conscious teams. It leverages compression-based anomaly detection (CBAD) to identify deviations from expected patterns in telemetry data.

### Core Features
- **Deterministic**: 100% reproducible results across runs
- **Explainable**: Glass-box explanations for detected anomalies
- **Compliance-Ready**: Built-in DORA, NIS2, and Runtime AI compliance
- **Real-time**: Sub-second anomaly detection and alerting
- **Privacy-First**: On-premises deployment with configurable data redaction

### Technology Stack
- **CBAD Engine**: Rust (cbad-core crate)
- **API Layer**: Rust (Axum, Tokio, sqlx) with PostgreSQL
- **Frontend**: Vue 3 with TypeScript, Tailwind CSS
- **Auth**: Firebase JWT + API keys (Argon2)
- **Billing**: Stripe
- **Metrics**: Prometheus

## Architecture

### High-Level Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Sources  │    │   Driftlock API │    │  UI Dashboard   │
│                 │    │     (Rust)      │    │     (Vue)       │
│  - Logs         │───▶│  - Detection    │───▶│  - Visualize    │
│  - Metrics      │    │  - Streams      │    │    anomalies    │
│  - Traces       │    │  - Billing      │    │  - Manage       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │                 │
                       │  - Anomalies    │
                       │  - Streams      │
                       │  - Tenants      │
                       └─────────────────┘
```

### Data Flow
1. Data sources send events via HTTP API
2. Events processed through CBAD detection engine
3. Compression ratios and statistical measures computed
4. Anomalies detected based on NCD and p-value thresholds
5. Anomaly details stored in PostgreSQL
6. Real-time notifications via webhook
7. UI dashboard displays anomalies with explanations

## Development Environment Setup

### Prerequisites
- Rust 1.75+ (install via rustup)
- PostgreSQL 15+ (or Docker for testing)
- Node.js 18+ (for UI development)
- Docker and Docker Compose (optional)

### Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
```

2. **Set up environment variables:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Build the CBAD core:**
```bash
cargo build -p cbad-core --release
```

4. **Build the API server:**
```bash
cargo build -p driftlock-api --release
```

5. **Start PostgreSQL with Docker:**
```bash
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15
```

6. **Run database migrations:**
```bash
# Migrations run automatically on startup, or manually:
sqlx migrate run --database-url "$DATABASE_URL"
```

7. **Start the API server:**
```bash
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  cargo run -p driftlock-api --release
```

### Running with Docker Compose
For a complete local development setup:

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

## Running Tests

### API Server Tests
```bash
# Run all tests
cargo test -p driftlock-api

# Run tests with output
cargo test -p driftlock-api -- --nocapture

# Run specific test
cargo test -p driftlock-api test_health_check
```

### CBAD Core Tests
```bash
# Run CBAD core tests
cargo test -p cbad-core

# Run with release optimizations
cargo test -p cbad-core --release
```

### Database Tests
```bash
# Run database layer tests
cargo test -p driftlock-db
```

### All Tests
```bash
# Run entire workspace
cargo test --workspace
```

## Code Style

### Rust Code Style
- Follow the [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)
- Use `rustfmt` for formatting: `cargo fmt`
- Use `clippy` for linting: `cargo clippy`
- Write comprehensive tests for all business logic

### Import Grouping
```rust
// Standard library
use std::sync::Arc;
use std::time::Duration;

// External crates
use axum::{Router, routing::get};
use sqlx::PgPool;
use tracing::info;

// Internal crates
use driftlock_db::models::Tenant;
use driftlock_auth::ApiAuth;
```

### Error Handling
- Use `thiserror` for custom error types
- Use `anyhow` for application-level errors
- Always provide context with errors
- Return structured API errors

### Naming Conventions
- Use descriptive names: `tenant_id` instead of `tid`
- Be consistent with terminology throughout the codebase
- Follow Rust naming conventions (snake_case for functions/variables, PascalCase for types)

## Contribution Workflow

### Branching Strategy
- Use feature branches off `main`
- Name branches descriptively: `feature/anomaly-detection-improvements` or `bugfix/api-rate-limiting`
- Submit pull requests to `main`

### Pull Request Process
1. Create an issue describing the feature or bug
2. Fork the repository and create a feature branch
3. Make your changes following the code style
4. Write or update tests as appropriate
5. Update documentation if needed
6. Submit a pull request with a clear description
7. Address review comments
8. Your PR will be merged after approval

### Commit Messages
- Use imperative mood: "Add feature" not "Added feature"
- First line should be 50 characters or less
- Include a body if needed to explain the "why"
- Link to issues when applicable

Example:
```
Add support for custom anomaly severity thresholds

Enables users to configure different severity levels based on
NCD scores and p-values. This allows for more granular alerting
strategies in different environments.

Resolves #123
```

## Project Structure

```
driftlock/
├── Cargo.toml                    # Workspace configuration
├── Dockerfile                    # API container definition
├── cbad-core/                    # Rust CBAD engine
│   ├── src/
│   │   ├── lib.rs               # Main library entry
│   │   ├── anomaly.rs           # Anomaly detection
│   │   ├── window.rs            # Sliding windows
│   │   └── metrics/             # Statistical metrics
│   └── Cargo.toml
├── crates/
│   ├── driftlock-api/           # Axum HTTP server
│   │   ├── src/
│   │   │   ├── main.rs          # Server entry point
│   │   │   ├── routes/          # HTTP handlers
│   │   │   ├── middleware/      # Auth, rate limiting
│   │   │   ├── state.rs         # Application state
│   │   │   └── errors.rs        # Error types
│   │   └── Cargo.toml
│   ├── driftlock-db/            # Database layer
│   │   ├── src/
│   │   │   ├── lib.rs           # Module exports
│   │   │   ├── models/          # Data models
│   │   │   ├── repos/           # Repository implementations
│   │   │   └── pool.rs          # Connection pooling
│   │   └── Cargo.toml
│   ├── driftlock-auth/          # Authentication
│   │   ├── src/
│   │   │   ├── lib.rs
│   │   │   ├── firebase.rs      # Firebase JWT
│   │   │   └── api_key.rs       # API key auth
│   │   └── Cargo.toml
│   ├── driftlock-billing/       # Stripe integration
│   └── driftlock-email/         # Email service
├── landing-page/                 # Vue frontend
├── docs/                         # Documentation
├── archive/go-backend/          # Legacy Go (reference only)
└── .github/                      # GitHub configuration
```

### Key Directories
- `crates/driftlock-api/src/routes/`: HTTP route handlers
- `crates/driftlock-db/src/repos/`: Database operations
- `crates/driftlock-auth/src/`: Authentication logic
- `cbad-core/src/`: CBAD algorithm implementation

## Performance Considerations

### API Server
- Uses async Tokio runtime for concurrent requests
- Connection pooling via sqlx
- In-memory rate limiting for single-instance deployment

### Database
- Create proper indexes for common query patterns
- Use prepared statements (sqlx compile-time checking)
- Connection pool configured for optimal throughput

### CBAD Engine
- Optimized compression algorithms
- Efficient sliding window implementation
- Memory-efficient data structures

## Security Best Practices

### Input Validation
- All inputs validated via serde + custom validators
- SQL injection prevented by sqlx parameterized queries
- Rate limiting on public endpoints

### Authentication & Authorization
- Firebase JWT tokens with proper verification
- API keys hashed with Argon2
- Role-based access (admin, stream)

### Data Protection
- HTTPS required in production
- Sensitive data encrypted at rest
- Audit logging for compliance

## Troubleshooting

### Common Issues
- **Database connection failures**: Check `DATABASE_URL` and ensure PostgreSQL is running
- **Build failures**: Ensure Rust 1.75+ is installed
- **Rate limiting**: Check rate limiter configuration

### Debugging Tips
- Set `RUST_LOG=debug` for detailed logging
- Use `cargo test -- --nocapture` to see test output
- Check PostgreSQL logs for query issues
- Monitor metrics via the `/metrics` endpoint

## Contact and Support

- For questions about the codebase: [team@driftlock.io]
- For security issues: [security@driftlock.io]
- For general support: Open an issue in the GitHub repository

For more detailed information, see the individual README files in each crate directory.
