# Driftlock Use Cases

Driftlock's compression-based anomaly detection (CBAD) is a versatile engine that works on any stream of structured data. While our roots are in financial compliance, the ability to mathematically prove "weirdness" without training data makes it ideal for a wide range of modern applications.

## 1. Financial & Regulatory Compliance (Core)

This is our primary enterprise use case, driven by regulations like **EU DORA** and **US FFIEC**.

*   **Problem**: Banks must explain *why* a transaction was flagged. " Black-box AI said so" is no longer a valid legal defense.
*   **Driftlock Solution**: We provide a mathematical proof (NCD score + p-value) for every decision.
*   **Examples**:
    *   **Algorithmic Trading**: Detect sudden shifts in market behavior or trading bot logic that deviate from the norm.
    *   **Fraud Detection**: Identify payment patterns that mathematically resemble known fraud or significantly differ from a user's baseline.
    *   **AML (Anti-Money Laundering)**: Spot structuring and layering attempts that statistical baselines miss.

## 2. Cybersecurity & Operations

Because Driftlock learns "normal" behavior from a small window of events, it adapts instantly to new operational environments.

*   **DDOS & Traffic Patterns**:
    *   Detect subtle changes in request headers or payload entropy that signal a sophisticated application-layer attack.
    *   Spot "low and slow" attacks that threshold-based WAFs miss.
*   **API Abuse**:
    *   Identify scraped data or non-human interaction patterns by analyzing the compression ratios of request sequences.
    *   Flag credential stuffing attacks where the payload entropy changes slightly.

## 3. AI Agent Monitoring (Future-Proofing)

As autonomous AI agents become more common, monitoring their behavior is critical. Agents can hallucinate, loop, or drift into unsafe states.

*   **Behavioral Drift**:
    *   Detect when an agent's output stream (logs, actions, text) statistically deviates from its instructions or previous successful runs.
    *   "Circuit Breaker": Automatically pause an agent if its NCD score spikes, preventing runaway costs or damage.
*   **Safety & Alignment**:
    *   Monitor the internal thought traces or tool usage of an agent. If the complexity (entropy) of its reasoning suddenly drops or explodes, it might be hallucinating or stuck.

## 4. IoT & Smart Home

The "Internet of Things" generates massive streams of sensor data. Driftlock runs efficiently (even on edge devices via Rust) to detect hardware failures or intrusions.

*   **Smart Home Security**:
    *   Analyze network traffic from smart bulbs and cameras. If a device starts sending data to a new IP or with a different encryption pattern, Driftlock flags it instantly.
*   **Predictive Maintenance**:
    *   Monitor vibration or temperature sensors on industrial machinery.
    *   Detect the mathematical signature of a failing bearing *before* it breaks, without needing terabytes of failure training data.

## Why Driftlock?

For all these use cases, the advantage is the same:
1.  **No Training**: You don't need months of labeled "bad" data.
2.  **Explainable**: You get a math score, not a black-box guess.
3.  **Fast**: Sub-millisecond detection suitable for real-time blocking.

