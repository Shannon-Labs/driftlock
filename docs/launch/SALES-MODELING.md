# Driftlock Sales & Pricing Modeling (Draft)

> This document is an internal modeling aid, not a public price sheet. Numbers are directional and should be refined with real usage and cost data.

## 1. Core Value Props and Pricing Anchors

- **Explainable anomaly detection for regulated workloads** (EU DORA, NYDFS, CFPB, FFIEC).
- **Developer-first HTTP API** that can be dropped into any payment gateway, risk engine, or AI training pipeline.
- **Compliance-grade evidence**: NCD, p-values, confidence, and human-readable explanations per anomaly.

Anchors from market comparables (approximate):

- Observability/logs:
  - New Relic: ~$0.40/GB beyond 100 GB free.
  - Sentry Logs: ~$0.50/GB.
  - Azure Monitor Analytics Logs: ~$1.5–$2.3/GB at low volume.
- ML/AI operations:
  - AWS DevOps Guru: ~$40 per 1M API calls (+ underlying CloudWatch/log costs).
- Serverless infra reference:
  - Upstash: ~$2 per 1M commands + ~$0.25/GB storage.

**Driftlock target:** stay meaningfully below these per-event/per-GB rates while maintaining 70–85% gross margin on variable costs.

## 2. Internal Cost Assumptions (Rough)

These are directional for planning; refine with real infra data.

- **Variable cost per 1M anomaly checks** (Rust core + Go API + Postgres + storage + bandwidth):
  - Low/medium scale: $0.20–$0.60 per 1M calls.
  - Optimized scale: could trend toward $0.10–$0.30 per 1M.
- **Semi-fixed per-tenant costs (enterprise)**:
  - Dedicated DB slice, monitoring, backups, compliance tooling, and support:
  - $200–$1,000/month effective cost depending on tenant size and SLA.

We should aim for **70–85% gross margin** on the variable component and allow fixed costs to amortize over volume.

## 3. Product Surfaces and Meters

Driftlock exposes three primary monetization surfaces:

1. **Developer Anomaly API**
   - Meter: anomaly checks (events) and/or GB of payload.
   - Usage pattern: SaaS apps, payment gateways, AI training/monitoring pipelines.

2. **Enterprise Compliance API / Platform**
   - Meter: events/GB + account-level features (retention, SSO, regional routing, support).
   - Usage pattern: banks, PSPs, critical-infrastructure operators.

3. **Consulting & Enablement**
   - Meter: project-based pricing (DORA readiness reviews, model governance workshops, integration projects).

## 4. Developer Anomaly API (Internal Model)

### 4.1. Free Tier (Acquisition)

- Included:
  - Up to **5M anomaly checks/month** _or_ **5 GB/month** of payload (choose 1 canonical meter for implementation).
  - 90-day retention.
- Target profile: individual developers, small teams, experimentation, open-source usage.
- Expected cost:
  - 5M events × $0.30/1M = **$1.50/month** per active free tenant at medium scale.

### 4.2. Paid Usage – Event-Based

Proposed internal ladder:

- **0–100M anomaly checks/month**: **$1.00 per 1M**.
- **100–500M anomaly checks/month**: **$0.80 per 1M**.
- **500M+ anomaly checks/month**: **$0.60 per 1M**.

Margin sanity check:

- If cost ≈ $0.30 per 1M, revenue at $1.00 per 1M → **70% gross margin**.
- If cost improves to $0.20 per 1M, margin at $1.00 per 1M → **80%**.

### 4.3. Alternative: GB-Based

If we prefer GB-based pricing (to align with observability benchmarks):

- **0–10 TB/month**: **$0.20 per GB**.
- **10 TB+**: **$0.15 per GB**.

Benchmark comparison:

- New Relic: $0.40/GB → Driftlock is ~50% cheaper.
- Azure Monitor Analytics Logs: ~$1.5–$2.3/GB → Driftlock is significantly cheaper.

### 4.4. Add-Ons

- **Extended retention (>90 days)**:
  - Charge: $0.03–$0.05 per GB-month.
  - Cost: map to chosen cold storage / archive solution.
- **Compliance/SLA add-on** for dev tenants:
  - Charge: $100–$200/month.
  - Includes: faster response targets, audit log export, guidance on using Driftlock evidence in audits.

## 5. Enterprise Compliance Plans (Internal Model)

Enterprise pricing combines a **monthly minimum** (for features and support) plus **usage-based charges** for events/GB.

### 5.1. Tier A – Compliance Starter

- Target customers: fintechs, regional banks, high-value B2B SaaS.
- Internal list price: **$2,500–$5,000/month** minimum.
- Included:
  - Up to **500M anomaly checks/month** or **5 TB/month**.
  - 180-day retention.
  - Basic SSO (SAML/OIDC) and role-based access.
  - Email support with defined response windows.

### 5.2. Tier B – Growth Bank / Multi-Region

- Target customers: mid-size banks, PSPs, multi-region SaaS platforms.
- Internal list price: **$8,000–$15,000/month** minimum.
- Included:
  - **1–2B anomaly checks/month** or **10–20 TB/month**.
  - 1-year retention.
  - Regional data residency options.
  - Priority support and change-control reviews.

### 5.3. Tier C – Tier-1 Bank / Critical Infrastructure

- Target customers: global banks, major PSPs, critical infrastructure networks.
- Internal target range: **$25,000–$75,000+/month** (contracted).
- Included:
  - Multi-region deployment, custom SLAs, and on-call integration.
  - Co-authored audit runbooks and regulator-facing documentation.
  - Tight integration with client SIEM/GRC systems.

### 5.4. Effective Unit Price and Margin

Example for Tier B:

- 1B events/month, effective charge ≈ $0.50 per 1M → $500 variable revenue.
- Monthly minimum of $10,000 covers:
  - Infra overhead.
  - Support / compliance labor.
  - Significant margin for product development.

Even if variable charges compress at scale, most enterprise revenue should come from the minimum.

## 6. Consulting & Enablement Levers

Consulting packages can be priced independently but should reinforce the platform:

- **DORA/NYDFS Readiness Assessment**
  - Indicative range: $25k–$75k per engagement.
  - Includes: inventory of current fraud models, gap analysis, recommended Driftlock integration.

- **Model Governance & Explainability Workshop**
  - Indicative range: $10k–$30k.
  - Audience: risk, compliance, data science.

- **Integration Accelerator (API + pipelines)**
  - Indicative range: $20k–$100k depending on scope.

These projects are high-margin and should act as accelerants for platform adoption.

## 7. Public Messaging vs. Internal Numbers

For README and landing page, we should **avoid locking in exact prices** at this stage. Recommended external language:

- "Developer API with a generous free tier and usage-based pricing around **$1 per million anomaly checks**, with volume discounts."
- "Data-based plans are **meaningfully cheaper per GB** than full-stack observability tools, because Driftlock focuses on explainable anomaly signals."
- "Enterprise compliance plans start in the **low thousands per month** and are designed to stay below the combined cost of fines, existing monitoring, and manual audit work."

Internal numbers in this doc can evolve independently of the public story.

## 8. Next Steps

- Validate cost assumptions with real infra metrics.
- Choose the primary meter (events vs GB) for implementation.
- Add a small, non-interactive ROI example to the landing page.
- Later: consider an interactive calculator once the numbers stabilize.
