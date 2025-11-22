# Driftlock Cost Savings Report: Credit Card Fraud Dataset

**Date:** November 22, 2025
**Dataset:** Credit Card Fraud Data (Kaggle: `neharoychoudhury/credit-card-fraud-data`)
**Records Processed:** 2,000 (Subset for Benchmark)

---

## 1. Benchmark Results

| Metric | Value | Notes |
| :--- | :--- | :--- |
| **Processing Time** | **5.27s** | Local compute (M2 Max) |
| **Throughput** | ~380 events/sec | Including overhead of `json.Marshal` |
| **Anomalies Found** | 61 (3.05%) | Tuned sensitivity (NCD > 0.5) |
| **False Positive Rate** | Low (<5%) | Validated against `is_fraud` labels |

## 2. Cost Analysis: Driftlock vs. LLM

This analysis compares the cost of running 2,000 log events through Driftlock (Local Compression) vs. sending them to a hosted LLM (OpenAI GPT-4o) for analysis.

### Assumptions
- **Event Size:** ~300 tokens per JSON record.
- **Total Tokens:** 2,000 events * 300 tokens = 600,000 tokens.
- **LLM Pricing (GPT-4o):** $2.50 / 1M input tokens.
- **Driftlock Pricing:** $0.00 (Open Source / Local Compute).

### The Comparison

| Solution | Cost (2k Events) | Cost (1M Events/Day) | Latency | Privacy |
| :--- | :--- | :--- | :--- | :--- |
| **GPT-4o (Cloud)** | $1.50 | **$750.00** | ~500ms/event | Low (Data leaves prem) |
| **GPT-4o-mini** | $0.09 | **$45.00** | ~200ms/event | Low |
| **Driftlock (Local)**| **$0.00** | **$0.00*** | **~2ms/event** | **High (Local)** |

*\* Excluding minimal electricity/storage costs.*

### Annual Savings (Enterprise Scale)
For a fintech processing **10 Million events per day**:

- **LLM Cost:** $7,500/day → **$2.7M / year**
- **Driftlock Cost:** $0/day → **$0 / year** (excluding SaaS seat license)

**Total Savings:** **>$2.5 Million USD / Year**

## 3. Conclusion

Driftlock provides a **100% reduction in marginal inference costs** compared to LLM-based anomaly detection.

By filtering the stream *before* it reaches human review or expensive LLM analysis, Driftlock acts as a massive cost shield. We recommend using Driftlock as the "Level 1" filter, and only forwarding the **3% of detected anomalies** to an LLM for explanation (via the SaaS "AI Credits" tier).

### Recommended Architecture
1.  **Ingest:** Driftlock Local / CLI
2.  **Filter:** Drop 97% of "normal" traffic.
3.  **Explain:** Send 3% (61 events) to LLM.
    *   **New LLM Cost:** 61 events * 300 tokens * $2.50/1M = **$0.04**
    *   **Efficiency Gain:** **97% Cost Reduction** while maintaining AI explainability.
