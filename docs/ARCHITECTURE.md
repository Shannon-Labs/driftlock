# Architecture Overview

## Core Components

- **cbad-core (Rust)**: The mathematical core. Implements compression-based algorithms (NCD) with FFI for Go.
- **driftlock-http (Go)**: The primary API service.
    - **Detection Engine**: Wraps `cbad-core` via CGO.
    - **Streaming Manager**: Handles real-time SSE connections for live anomalies.
    - **API**: REST endpoints for detection and history.
- **Landing Page / Dashboard (Vue 3)**:
    - **Frontend**: Single-page application for marketing and dashboard.
    - **Proxy**: Cloudflare Pages Functions / Firebase Hosting rewrites for API access.

## SaaS Deployment Architecture

Driftlock is architected for serverless, auto-scaling deployment on Google Cloud Platform.

### 1. Frontend Layer (Firebase Hosting)
- Hosts the static assets (HTML/JS/CSS).
- Provides global CDN distribution.
- Handles SSL termination.
- Proxies `/api/*` requests to the backend (via rewrites).

### 2. Application Layer (Cloud Run)
- **Stateless Container**: The `driftlock-http` service runs as a stateless container.
- **Auto-scaling**: Scales from 0 to N based on request load.
- **Streaming**: Supports long-lived HTTP connections for SSE.

### 3. Data Layer (Cloud SQL / Postgres)
- **Persistence**: Stores anomaly records, tenant configuration, and API keys.
- **Migration**: Schema managed via Go migrations.

### 4. Streaming Flow
1. **Agent** sends event to `/v1/detect` (HTTP POST).
2. **driftlock-http** calculates NCD against baseline.
3. If anomalous:
    - Saves record to Postgres.
    - **Broadcasts** to `StreamManager`.
    - Pushes to connected SSE clients (Dashboard).

## Data Flow

1. **Ingest**: Agents/Apps send JSON events to API.
2. **Process**: API calls Rust core to compute compression distance.
3. **Detect**: If distance > threshold, flag as anomaly.
4. **Alert**: Push to real-time stream and store evidence.
5. **Explain**: Dashboard renders NCD explanation.

## Determinism

- Use deterministic seeds in permutation testing.
- Configure windows (baseline/window/hop) and thresholds explicitly.
- Avoid non-deterministic concurrency paths in the core algorithm.
