# Driftlock Pricing Study Report

**Date:** November 19, 2025
**Status:** Complete

## Executive Summary

Based on load testing with representative transaction data, the estimated infrastructure cost for Driftlock running on Google Cloud Run is approximately **$0.05 per 1 million events**. This includes compute and memory costs, assuming batching of 50 events per request.

## Methodology

- **Environment:** Local Docker stack (simulating Cloud Run latency profile)
- **Compute:** 1 vCPU, 0.5 GB Memory (Cloud Run Tier 1 pricing)
- **Load:** Concurrent requests (4 workers) processing JSONL transaction logs.
- **Batch Size:** 50 events per API request.

## Detailed Results

### 1. Normal Transactions (Baseline)
- **Dataset:** `normal-transactions.jsonl` (500 events)
- **Throughput:** ~2,100 events/sec
- **Average Latency:** 80ms per batch
- **Cost per 1k Events:** $0.000049
- **Projected Cost per 1M Events:** **$0.049**

### 2. Mixed Transactions (Streaming Simulation)
- **Dataset:** `mixed-transactions.jsonl` (1000 events, varied payload)
- **Throughput:** ~2,600 events/sec
- **Average Latency:** 73ms per batch
- **Cost per 1k Events:** $0.000045
- **Projected Cost per 1M Events:** **$0.045**

## Pricing Recommendations

Given the base infrastructure cost of ~$0.05/million events, we recommend the following pricing tiers to ensure healthy margins while remaining competitive:

| Tier | Price | Included Events | Effective Cost/M | Margin |
|---|---|---|---|---|
| **Developer** | Free | 10,000 / mo | $0 | N/A |
| **Starter** | $99 / mo | 500,000 / mo | $198 | >99% |
| **Growth** | $499 / mo | 5,000,000 / mo | $99.80 | >99% |

*Note: These margins do not account for storage costs (PostgreSQL), which will scale with data retention policies.*

## Next Steps

1. **Data Retention:** Analyze storage costs for 30-day vs 90-day retention to refine the "Growth" tier margin.
2. **Optimization:** Increasing batch size to 100-200 could further reduce request costs by 50%.

