# Driftlock Documentation

Compression-based anomaly detection for logs, metrics, traces, and LLM I/O.

## Quick Links

| I want to… | Go to… |
| --- | --- |
| Try it now (no signup) | [Demo endpoint](./user-guide/api/endpoints/demo.md) |
| Get started in 5 minutes | [Quickstart Guide](./user-guide/getting-started/quickstart.md) |
| See API endpoints | [REST API Reference](./user-guide/api/rest-api.md) |
| Understand how it works | [Core Concepts](./user-guide/getting-started/concepts.md) |
| Run a tutorial | [Tutorials & Examples](./user-guide/tutorials/README.md) |
| Handle errors | [Error Codes](./user-guide/api/errors.md) |
| Operate in production | [Operations & Runbooks](./user-guide/guides/operations.md) |
| Check compliance | [DORA](./compliance/COMPLIANCE_DORA.md) |

## What is Driftlock?

Driftlock detects anomalies by analyzing how well your data compresses. When new data compresses poorly against a learned baseline, it's likely anomalous.

**Key features**
- No model training required—works immediately on any JSON payload
- Deterministic, explainable results with plain-language reasons
- Low false positives via p-value significance testing
- Compliance-ready evidence bundles (DORA, NIS2, AI Act)

## Try It Now (no signup)

```bash
curl -X POST https://api.driftlock.net/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"body": {"cpu": 45, "memory": 2048}},
      {"body": {"cpu": 47, "memory": 2100}},
      {"body": {"cpu": 44, "memory": 2050}},
      {"body": {"cpu": 99, "memory": 8000}}
    ]
  }'
```

The last event (CPU spike + memory spike) will be flagged as anomalous.

## Documentation Map

- **Getting Started:** [Quickstart](./user-guide/getting-started/quickstart.md), [Core Concepts](./user-guide/getting-started/concepts.md), [Authentication](./user-guide/getting-started/authentication.md)
- **API Reference:** [REST API](./user-guide/api/rest-api.md), [POST /v1/detect](./user-guide/api/endpoints/detect.md), [POST /v1/demo/detect](./user-guide/api/endpoints/demo.md), [GET /v1/anomalies](./user-guide/api/endpoints/anomalies.md), [Error Codes](./user-guide/api/errors.md)
- **Tutorials & Examples:** [Log monitoring](./user-guide/tutorials/log-monitoring.md), [Metrics spike detection](./user-guide/tutorials/metrics-spike-detection.md), [LLM output drift](./user-guide/tutorials/llm-output-drift.md), [Webhook/alerting](./user-guide/tutorials/webhook-alerting.md), [cURL examples](./user-guide/api/examples/curl-examples.md), [Python examples](./user-guide/api/examples/python-examples.md)
- **Operations & SRE:** [Operations & Runbooks](./user-guide/guides/operations.md), [Deployment runbooks](./deployment/RUNBOOKS.md), [Cloud Run setup](./deployment/cloud-run-setup.md)
- **Compliance:** [DORA](./compliance/COMPLIANCE_DORA.md), [NIS2](./compliance/COMPLIANCE_NIS2.md), [AI Act](./compliance/COMPLIANCE_RUNTIME_AI.md)
- **Architecture:** [Algorithms](./architecture/ALGORITHMS.md), [API architecture](./architecture/API.md)

## Pricing (snapshot)

| Plan | Price | Events/month | Streams | Best for |
| --- | --- | --- | --- | --- |
| Free | $0/mo | 10,000 | 5 | Experimentation |
| Pro | $99/mo | 500,000 | 20 | Startups & active monitoring |
| Team | $199/mo | 5,000,000 | 100 | High-volume production |
| Enterprise | Custom | Unlimited | 500+ | Large-scale compliance, EU data residency |

## Support

- Email: [support@driftlock.io](mailto:support@driftlock.io)
- GitHub: [Shannon-Labs/driftlock](https://github.com/Shannon-Labs/driftlock)
- Community: [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)
