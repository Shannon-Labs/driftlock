# Driftlock SaaS Pricing Strategy (2025)

**Status:** Implemented (November 2025)
**Target Audience:** DevOps Engineers, SREs, Compliance Officers in Regulated Industries
**Goal:** Maximize developer adoption while capturing enterprise value through compliance and AI upsells.

---

## Executive Summary

Based on late 2025 SaaS trends, the market has shifted definitively toward **Hybrid Usage-Based Pricing**. Pure seat-based models are dead for developer tools, and pure usage-based models create procurement friction.

Driftlock adopts a **"Compliance-Core, AI-Flex"** model:
1.  **Core Platform (SaaS Subscription):** Predictable monthly fee for data retention, users, and compliance reporting (DORA/NIS2).
2.  **AI Insights (Included/Usage):** AI explanations are included in paid tiers with generous limits.

---

## 1. The Pricing Tiers

### **Developer (Free)**
*For individuals and local testing.*
- **Price:** $0/mo
- **Deployment:** Local CLI + Limited Hosted
- **Features:**
  - 10,000 events / month
  - Basic anomaly detection (zstd compression)
  - CLI demo & Playground access
  - JSON reports
  - Community support
  - 7-day data retention
- **Goal:** Ubiquity. Every developer should have `driftlock` in their local toolchain.

### **Basic ($20/mo)**
*For startups and individual teams.*
- **Price:** $20/month
- **Usage:** 500,000 events / month
- **Features:**
  - Hosted Dashboard (`driftlock.net`)
  - 30-day retention
  - Email alerts
  - **AI Insights:** Gemini summaries for detected anomalies
- **Goal:** Low friction adoption. "Put it on the corporate card."

### **Pro ($200/mo)**
*For growing companies and compliance-sensitive teams.*
- **Price:** $200/month
- **Usage:** 5,000,000 events / month
- **Features:**
  - **DORA & NIS2 Compliance Evidence Bundles**
  - 90-day retention
  - Priority Support
  - Advanced Reporting
  - Higher rate limits
- **Goal:** Compliance readiness for scale-ups.

### **Enterprise (Custom)**
*For regulated organizations (FinTech, HealthTech).*
- **Price:** Custom (Contact Sales)
- **Usage:** Unlimited / Custom volume
- **Features:**
  - SSO / SAML
  - **Adaptive Windowing:** "Ungameable" randomized audit intervals (Security Shield)
  - Custom Retention
  - **Private Cloud / On-Prem Options**
  - **Custom Compression Models (OpenZL)**
  - Dedicated Account Manager
- **Goal:** High ACV. Selling "Audit Insurance," not just logging.

---

## 2. The 2025 "Hook": AI Cost Transparency

Market research indicates deep fatigue with opaque AI upcharges. Driftlock differentiates by:

1.  **Deterministic First:** We emphasize that our core detection is *mathematical* (NCD, Entropy) and costs near-zero to run.
2.  **AI Optional:** AI is positioned as an "Explainer Layer," not the detection engine.
3.  **Value-Add:** We price the AI capability into the plan rather than charging per-token, simplifying procurement.

## 3. Acquisition Strategy

1.  **The "Trojan Horse":** The open-source CLI (`driftlock-cli`) is the primary marketing channel. It includes a `driftlock login` command that frictionlessly upgrades local users to the SaaS pilot.
2.  **Compliance Fear:** Marketing materials focus heavily on **DORA** (Digital Operational Resilience Act) deadlines. "Is your anomaly detection audit-ready?"
3.  **Usage Transparency:** The dashboard features a prominent "Cost Forecast" widget (Brutalist design), showing exactly what the bill will be, updated in real-time.

---

## 4. Implementation Roadmap

### Phase 1: Infrastructure (Current)
- [x] Stripe Integration (Subscriptions)
- [x] Usage Metering (Postgres/Redis counters)
- [x] Multi-tier support (Developer, Basic, Pro)

### Phase 2: Packaging (Q1 2026)
- [ ] Automated DORA Report generation (PDF) for Enterprise
- [ ] Self-Service upgrade flow refinement

### Phase 3: Enterprise (Q2 2026)
- [ ] SSO Integration (Okta/Auth0)
- [ ] BYO-Key for AI models
