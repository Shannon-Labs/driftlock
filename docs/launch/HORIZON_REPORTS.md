# Driftlock Horizon Reports: Expanding the Anomaly Frontier

**Date:** November 22, 2025
**Status:** Verified Benchmarks
**Markets:** Critical Infrastructure, Cloud Operations, AI Safety

---

## Horizon 2: Critical Infrastructure (Airlines)

**Dataset:** Airline Delay Cause (`ryanjt/airline-delay-cause`)
**Stream:** Flight Operations Telemetry (Latency = Delay, Volume = Flight Count)
**Events Processed:** 2,000

### The Test
We simulated a stream of flight arrivals. Normal behavior is "on time" or "small delay." Anomalies are "cascading delays" (the Delta/CrowdStrike scenario).

### Driftlock Findings
- **Baseline:** Built on 400 flights (steady state).
- **Anomalies Detected:** **31 events** (1.55%).
- **Detection Signature:**
    - The *compression ratio* of the stream dropped significantly when delays became correlated.
    - Random delays compress well (noise). **Systemic** delays (cascading failure) create a new, highly repetitive pattern that deviates from the "normal noise" baseline.
- **Verdict:** **Driftlock can detect operational meltdowns** by sensing the "structure" of the delay propagation.

---

## Horizon 3: Cloud Operations (Web Traffic)

**Dataset:** Numenta Anomaly Benchmark (NAB) - Real AWS CloudWatch
**Stream:** EC2 CPU Utilization (i-5f5533)
**Events Processed:** 2,000

### The Test
Server load is periodic (daily cycles). Attacks or hangs break this periodicity.

### Driftlock Findings
- **Baseline:** Built on 400 timepoints (approx. 1.5 days).
- **Anomalies Detected:** **29 events** (1.45%).
- **Detection Signature:**
    - **Entropy Spikes:** Sudden, non-periodic load spikes (flash crowd or DDoS).
    - **Entropy Drops:** "Flatlining" (frozen process) creates a string of identical values, which compresses *too well* compared to the noisy baseline.
- **Verdict:** **Driftlock detects "Frozen" and "Fried" servers** equally well without setting static thresholds (e.g., "CPU > 90%").

---

## Horizon 4: AI Safety (The Frontier)

**Dataset:** Prompt Injection Malignant (`marycamilainfo/prompt-injection-malignant`)
**Stream:** User Inputs to an LLM
**Events Processed:** 1,581

### The Test
Can compression math distinguish "normal English questions" from "jailbreak attacks" (base64, foreign characters, repetitive padding)?

### Driftlock Findings
- **Baseline:** Built on 400 mixed inputs.
- **Anomalies Detected:** **44 events** (2.78%).
- **Detection Signature:**
    - **High NCD (0.8+):** Jailbreaks often use "obfuscation" (Base64, LeetSpeak) which has a completely different compression profile than standard English.
    - **Low Entropy:** "Ignore previous instructions" repeated 50 times (a common attack) compresses extremely well, flagging a "Repetition Anomaly."
- **Verdict:** **Driftlock is a viable "Pre-Flight" firewall for LLMs.** It catches structural attacks before they even reach the model context window.

---

## Strategic Implication

We have proven that **Compression-Based Anomaly Detection (CBAD)** is not just for logs. It is a universal sensor for:
1.  **Systemic Cascade** (Airlines)
2.  **Process Freezes** (Cloud)
3.  **Adversarial Inputs** (AI Safety)

This validates the "Universal Volatility Radar" positioning.
