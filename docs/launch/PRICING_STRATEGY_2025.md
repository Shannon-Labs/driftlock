# Driftlock SaaS Pricing Strategy (2025)

**Status:** Draft (November 2025)
**Target Audience:** DevOps Engineers, SREs, Compliance Officers in Regulated Industries
**Goal:** Maximize developer adoption while capturing enterprise value through compliance and AI upsells.

---

## Executive Summary

Based on late 2025 SaaS trends, the market has shifted definitively toward **Hybrid Usage-Based Pricing**. Pure seat-based models are dead for developer tools, and pure usage-based models create procurement friction.

Driftlock will adopt a **"Compliance-Core, AI-Flex"** model:
1.  **Core Platform (SaaS Subscription):** Predictable monthly fee for data retention, users, and compliance reporting (DORA/NIS2).
2.  **AI Insights (Usage-Based):** "Pay-as-you-ask" model for LLM-based anomaly explanations and root cause analysis.

---

## 1. The Pricing Tiers

### **Developer (Free)**
*For individuals and local testing.*
- **Price:** $0/mo
- **Deployment:** Local CLI / Docker only
- **Features:**
  - Unlimited local detection
  - Basic CLI reporting
  - Community support
- **Goal:** Ubiquity. Every developer should have `driftlock` in their local toolchain.

### **Pro (SaaS Pilot)**
*For startups and individual teams.*
- **Price:** $99/month (includes 5 seats)
- **Usage:** Up to 10GB logs/day
- **Features:**
  - Hosted Dashboard (`driftlock.net`)
  - 14-day retention
  - Email alerts
  - **AI Credits:** 50/month (generative explanations)
- **Goal:** Low friction adoption. "Put it on the corporate card."

### **Enterprise (Compliance)**
*For regulated organizations (FinTech, HealthTech).*
- **Price:** Starting at $1,500/month
- **Usage:** Custom volume
- **Features:**
  - SSO / SAML
  - **DORA & NIS2 Compliance Bundles** (The "Killer Feature")
  - 90+ day retention
  - **Private AI:** Zero-retention LLM contracts or BYO-Key
  - Priority Support
- **Goal:** High ACV. Selling "Audit Insurance," not just logging.

---

## 2. The 2025 "Hook": AI Cost Transparency

Market research indicates deep fatigue with opaque AI upcharges. Driftlock differentiates by:

1.  **Deterministic First:** We emphasize that our core detection is *mathematical* (NCD, Entropy) and costs near-zero to run.
2.  **AI Optional:** AI is positioned as an "Explainer Layer," not the detection engine.
3.  **Pass-Through+ Pricing:** We charge a transparent markup on LLM tokens (e.g., cost + 20%) for explanation queries, or allow Enterprise customers to bring their own API keys (BYOK).

## 3. Acquisition Strategy

1.  **The "Trojan Horse":** The open-source CLI (`driftlock-cli`) is the primary marketing channel. It includes a `driftlock login` command that frictionlessly upgrades local users to the SaaS pilot.
2.  **Compliance Fear:** Marketing materials focus heavily on **DORA** (Digital Operational Resilience Act) deadlines. "Is your anomaly detection audit-ready?"
3.  **Usage Transparency:** The dashboard features a prominent "Cost Forecast" widget (Brutalist design), showing exactly what the bill will be, updated in real-time.

---

## 4. Implementation Roadmap

### Phase 1: Infrastructure (Current)
- [x] Stripe Integration (Subscriptions)
- [x] Usage Metering (Postgres/Redis counters)
- [ ] "Cost Forecast" UI Widget

### Phase 2: Packaging (Q1 2026)
- [ ] Self-Service Checkout for Pro Tier
- [ ] Automated DORA Report generation (PDF) for Enterprise

### Phase 3: Enterprise (Q2 2026)
- [ ] SSO Integration (Okta/Auth0)
- [ ] BYO-Key for AI models
