# Driftlock Transition Plan: General-Purpose Explainable AI Platform

## 1. Value Proposition Shift
**From:** "DORA/EU Compliance Tool for Banks"
**To:** "Universal Explainable Anomaly Detection for Developers & Enterprise"

**Core Message:**
"Detect anomalies in any data stream (logs, metrics, payments) with mathematical proof. No black boxes, no training data required, instant integration."

**Key Differentiators:**
1. **Universal:** Works on any structured (JSON) or unstructured data.
2. **Explainable:** "Show your work" via NCD (Normalized Compression Distance) - glass box, not black box.
3. **Developer-First:** CLI, simple API, instant results (no ML training phase).
4. **Enterprise-Ready:** Scalable, secure, and deployable on standard cloud infra (GCP/Firebase).

## 2. Technical Architecture Transition (Google Cloud / Firebase)

To make this "seamless" and "developer friendly", we will move to a serverless/managed architecture using **Google Cloud Platform (GCP)** and **Firebase**.

### A. Backend & Hosting (GCP/Firebase)
- **Firebase Hosting:** For the frontend (Vue/Vite app). Fast, global CDN, easy SSL.
- **Firebase Authentication:** Frictionless social login (Google/GitHub) for developers.
- **Cloud Run (The "Antigravity" lift):**
  - Host the `driftlock-http` API as a stateless container.
  - Auto-scales to zero (cost-efficient) and scales up instantly.
  - **Why:** This removes the "manage your own server" friction.
- **Firestore:** Store user configurations, API keys, and anomaly reports (replacing or augmenting Postgres for the SaaS tier).
- **Cloud SQL (Postgres):** Retain for high-volume anomaly persistence if needed, or strictly use Firestore for the "easy" tier.

### B. "Antigravity" (Ease of Use)
- **Concept:** "Antigravity" here metaphorically means removing the weight of DevOps.
- **Implementation:**
  - **One-Click Deploy:** "Deploy to Google Cloud" button for self-hosters.
  - **Managed SaaS:** We host it; they just send JSON to an endpoint.
  - **CLI:** `driftlock login`, `driftlock watch logs.json` (pushes to cloud).

## 3. Frontend Buildout & Redesign

The frontend needs to look like a modern SaaS tool (think Vercel, Supabase, Linear style), not a compliance audit tool.

### A. New Components
1. **Dashboard/Console:**
   - Real-time stream view (WebSocket/SSE from Cloud Run).
   - "Anomaly Stream" graph (d3/chartjs) showing compression distance over time.
   - API Key management.
2. **Integrations Page:**
   - "Copy/Paste" snippets for Python, Node.js, Go, Bash (curl).
   - Webhook configuration (Slack, PagerDuty).
3. **Documentation Hub:**
   - Interactive API explorer.
   - "How it Works" interactive visualization (visualizing compression).

### B. Refactoring `landing-page`
- **HomeView.vue:**
  - **Hero:** "Stop Debugging Black Boxes. Start Explaining Anomalies."
  - **Social Proof:** Developer testimonials (simulated/placeholder for now).
  - **Demo:** The interactive terminal/video is good, keep it but reframe context.
  - **Pricing:** "Free for Developers" (SaaS tier), "Enterprise" (Self-hosted/Compliance).
- **Remove:** Specific "DORA" scare tactics (move to a "Compliance" solutions page).

## 4. Execution Plan

### Phase 1: Rebranding & Cleanup (Immediate)
- [ ] Update `landing-page` copy to remove hard DORA focus.
- [ ] Create new "Solutions" pages for specific verticals (Compliance, DevOps, Security).
- [ ] Rename/Refactor frontend components to be generic (e.g., `RegulatoryMap` -> `GlobalDeploymentMap`).

### Phase 2: Firebase/GCP Setup
- [ ] Initialize Firebase project.
- [ ] Configure Firebase Hosting for `landing-page`.
- [ ] Set up Github Action to deploy frontend to Firebase on merge.
- [ ] (Backend) Prepare `driftlock-http` Dockerfile for Cloud Run deployment.

### Phase 3: SaaS Features (The "Easy" Stuff)
- [ ] Add "Login with Google" (Firebase Auth).
- [ ] Create "Dashboard" view (authenticated route).
- [ ] Implement API Key generation (stored in Firestore).

### Phase 4: Launch
- [ ] "Launch Week" style announcement.
- [ ] Publish "Driftlock: The Antigravity Guide to Anomaly Detection".

## 5. "Antigravity" Python Library (Bonus Idea)
To truly lean into "Antigravity", we could publish a tiny Python wrapper:
```python
import driftlock
# It just works - like import antigravity
driftlock.watch(my_data_stream)
```
This sends data to our Cloud Run instance.

---
**Next Step:** Confirm this plan, and I will begin with **Phase 1: Rebranding & Cleanup** on the frontend.

