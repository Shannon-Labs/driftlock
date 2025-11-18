# Handoff: SaaS Transformation & "Antigravity" Platform

**Status**: Phase 1 (Core Architecture) Complete. Ready for Phase 2 (Cloud Deployment & Polish).
**Target Audience**: AI Agent (Gemini 1.5 Pro / 3) or Senior DevOps Engineer.

## 1. Mission Overview

We are transforming Driftlock from a GitHub open-source demo into a **proprietary SaaS platform** called "The Safety Layer for Autonomous Agents."

**The Innovation**: Two APIs working in sync.
1.  **Streaming Telemetry**: Real-time ingest of agent thought traces/tool calls.
2.  **Explainable Anomaly Detection**: Deterministic, compression-based math (NCD) to flag drift.

**The Goal**: "Antigravity" — Ease of use. Users shouldn't manage servers. They just send data.

## 2. Architecture State

The codebase has been refactored to support a modern Serverless architecture.

*   **Frontend**: Vue 3 + Vite (located in `landing-page/`).
    *   **State Management**: Pinia store implemented (`stores/index.ts`).
    *   **Streaming**: Native `EventSource` (SSE) integration for real-time dashboard.
    *   **Hosting Target**: Firebase Hosting.
*   **Backend**: Go HTTP Service (located in `collector-processor/cmd/driftlock-http/`).
    *   **New Endpoint**: `/v1/stream/anomalies` (SSE) added for real-time push.
    *   **Hosting Target**: Google Cloud Run (stateless, autoscaling).
    *   **Proxy**: Firebase Hosting rewrites configured to proxy `/v1/*` to Cloud Run.
*   **Documentation**: OpenAPI 3.0 spec updated (`docs/api/openapi.yaml`).

## 3. Completed Work (Manifest)

The following files have been created or significantly modified and are **ready for deployment**:

### Infrastructure & Config
*   `deploy.sh`: Master script. Builds frontend, checks for gcloud/firebase CLI, deploys both.
*   `cloudbuild.yaml`: Google Cloud Build config for the backend container.
*   `service.yaml`: Knative/Cloud Run service definition.
*   `firebase.json`: Configured with Rewrite rules to route `/v1/` to the Cloud Run service.
*   `.gcloudignore`: Optimization for build context.

### Backend Code (`collector-processor/cmd/driftlock-http/`)
*   `main.go`: Added `StreamManager`, SSE handler, and broadcasting logic. The server now pushes anomalies to connected clients instantly.

### Frontend Code (`landing-page/`)
*   `src/views/HomeView.vue`: Rebranded with "Safety Layer" and "Streaming + Math" messaging.
*   `src/views/DashboardView.vue`: Now consumes a Pinia store. Displays a live, animating feed of anomalies.
*   `src/stores/index.ts`: New Pinia store for robust WebSocket/SSE state management.
*   `package.json` & `main.ts`: Added Pinia and TanStack Query.

### Documentation
*   `docs/deployment/cloud-run-setup.md`: Detailed guide for GCP setup.
*   `docs/deployment/firebase-hosting-setup.md`: Detailed guide for Firebase.
*   `docs/STREAMING.md`: Guide for the new SSE endpoint.
*   `docs/api/DISTRIBUTION.md`: Strategy for Client SDKs vs. PyPI.

## 4. Remaining Tasks (The "To-Do" for the Next Agent)

The infrastructure code is written, but the **actual cloud environment connection** needs to be verified/executed by an agent with CLI access or by guiding the user through it.

### Priority 1: Cloud Connection & Secret Management
- [ ] **Database**: The current `service.yaml` assumes a generic `DATABASE_URL` env var.
    -   **Action**: Provision a Cloud SQL (Postgres) instance or a Supabase project.
    -   **Action**: Set the `DATABASE_URL` secret in Google Cloud Secret Manager and reference it in `service.yaml`.
- [ ] **Auth Integration**: The frontend currently uses a placeholder "Demo Mode" (client-side keys).
    -   **Action**: Integrate **Firebase Authentication** (Google Sign-In).
    -   **Action**: Update `driftlock-http` middleware to verify Firebase ID Tokens (JWT) in the `Authorization` header.

### Priority 2: "Antigravity" SDKs
We decided *not* to ship a heavy "PyPI equivalent" of the engine, but we *do* need thin client wrappers.
- [ ] **Python Client**: Use `openapi-generator` (as documented in `DISTRIBUTION.md`) to generate a Python client.
    -   **Task**: Create a new directory `sdk/python`.
    -   **Task**: Publish a simple `pip install driftlock-client` package that wraps the generated code for better DX (e.g., `driftlock.monitor(stream)`).

### Priority 3: Production Polish
- [ ] **CORS**: Verify `CORS_ALLOW_ORIGINS` in `service.yaml` matches the final Firebase Hosting domain.
- [ ] **Domain**: Map `api.driftlock.net` (if using) to the Cloud Run service mapping or keep using the Firebase proxy.

## 5. Execution Instructions for Next Agent

1.  **Read** `docs/deployment/cloud-run-setup.md` to understand the deployment topology.
2.  **Read** `service.yaml` to see the environment variable expectations.
3.  **Ask the User**: "Do you have a Postgres database ready (Cloud SQL or Supabase)? I need the connection string to configure the backend secrets."
4.  **Action**: Once DB is ready, run `./deploy.sh` to push the initial stack.
5.  **Action**: tackle **Priority 1 (Auth)**. This is the biggest missing piece for a "real" SaaS.

**Codebase Constraints**:
- **Do NOT** rewrite the core detection logic in `cbad-core` (Rust).
- **Do NOT** change the `driftlock-http` FFI signatures.
- **Maintain** the "Demo Mode" fallback in the frontend until Auth is fully working.

