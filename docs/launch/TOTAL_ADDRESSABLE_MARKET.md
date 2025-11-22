# Driftlock: Total Addressable Market (TAM) Strategy

**Status:** Confidential Internal Strategy
**Date:** November 2025
**Thesis:** "The Volatility Radar for Everything"

---

## 1. Executive Summary

Driftlock is currently positioned as a compliance tool for FinTech. This is a **limiting frame**. The core technology—Compression-Based Anomaly Detection (CBAD)—is a universal **entropy sensor**.

**The Pitch:** "If a system has a 'normal' state, Driftlock can detect when it breaks, without training data."

This document outlines the **Four Horizons** of market expansion, moving from our beachhead (FinTech) to ubiquitous infrastructure sensing.

---

## 2. The Four Horizons

### Horizon 1: FinTech & RegTech (The Beachhead)
*   **Use Case:** Fraud detection, Algo-trading volatility alerts, DORA compliance.
*   **Why:** High willingness to pay, regulatory mandates.
*   **Data Proxy:** `neharoychoudhury/credit-card-fraud-data`
*   **Status:** Active.

### Horizon 2: Critical Infrastructure (The "Delta" Moment)
*   **Use Case:** Operational resilience for airlines, power grids, and logistics.
*   **Problem:** "Silent failures" where sensors drift before crashing (e.g., CrowdStrike update loops, Delta crew scheduling collapse).
*   **Driftlock Advantage:** Detects the *entropy shift* of a system entering a chaotic state before the hard crash.
*   **Data Proxy:** `ryanjt/airline-delay-cause` (Airline Ops), `pcbreviglieri/smart-grid-stability` (Power).

### Horizon 3: Cloud & Cyber (The "Cloudflare" Moment)
*   **Use Case:** DDoS detection, CDN outage prediction, API abuse.
*   **Problem:** "Low-and-slow" attacks that evade rate limits but change traffic texture.
*   **Driftlock Advantage:** Low-entropy attacks (botnets repeating requests) trigger massive compression spikes.
*   **Data Proxy:** `dhoogla/unswnb15` (Network Intrusion), `boltzmannbrain/nab` (Web Traffic).

### Horizon 4: AI Safety & Observability (The Frontier)
*   **Use Case:** Model collapse detection, Prompt Injection, Jailbreak attempts.
*   **Problem:** LLM inputs/outputs are unstructured text. Regex fails.
*   **Driftlock Advantage:** Adversarial prompts (base64 injections, weird unicode) have vastly different compression profiles than normal English.
*   **Data Proxy:** `marycamilainfo/prompt-injection-malignant`.

---

## 3. Competitive Moat (The "Anti-Wrapper" Strategy)

Competitors (Datadog, Splunk, AWS) will try to wrap this. We beat them by:

1.  **The "Auditor" Persona:** We don't just flag; we *explain* via the Gemini 3 Pro integration. "Verdict: CRITICAL."
2.  **Ungameable Telemetry:** Our randomized windowing (documented in `STREAMING_VARIABILITY_PROMPT.md`) makes us the only "security-grade" anomaly detector.
3.  **Local-First:** We run in the sidecar/agent. They run in the cloud. We are faster and cheaper.

---

## 4. Strategic Action Items

1.  **Validate Horizon 4:** Run benchmarks on `prompt-injection-malignant`. If Driftlock can catch jailbreaks purely via compression, we have a massive new product line ("Driftlock for AI").
2.  **Publish "The Chaos Report":** A whitepaper analyzing the Delta/Cloudflare crashes using our "What If" methodology.
3.  **Build the "Universal Sensor":** A generic MCP server that any agent can use to "taste" a data stream for chaos.
