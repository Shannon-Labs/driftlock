# Driftlock Documentation

Compression-based anomaly detection for your logs, metrics, traces, and LLM I/O.

## Quick Links

| I want to… | Go to… |
| --- | --- |
| Try it now (no signup) | [Demo endpoint](./user-guide/api/endpoints/demo.md) |
| Get started in 5 minutes | [Quickstart Guide](./user-guide/getting-started/quickstart.md) |
| See API endpoints | [API Reference](./user-guide/api/rest-api.md) |
| Handle errors | [Error Codes](./user-guide/api/errors.md) |
| Understand how it works | [Core Concepts](./user-guide/getting-started/concepts.md) |
| See pricing | [Pricing](./user-guide/guides/pricing-tiers.md) |

## What is Driftlock?

Driftlock detects anomalies by analyzing how well your data compresses. When new data compresses poorly against a learned baseline, it’s likely anomalous.

### Key features

- Works immediately on any JSON data (no model training)
- Deterministic and explainable results with plain-English reasons
- Low false positives via statistical significance testing (p-values)
- Compliance-ready evidence bundles (DORA, NIS2, AI Act)

## Try It Now (no signup)

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
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

## Documentation Sections

- **Getting Started:** [Quickstart](./user-guide/getting-started/quickstart.md), [Core Concepts](./user-guide/getting-started/concepts.md), [Authentication](./user-guide/getting-started/authentication.md)
- **API Reference:** [POST /v1/detect](./user-guide/api/endpoints/detect.md), [POST /v1/demo/detect](./user-guide/api/endpoints/demo.md), [GET /v1/anomalies](./user-guide/api/endpoints/anomalies.md), [Error Codes](./user-guide/api/errors.md)
- **Examples:** [cURL](./user-guide/api/examples/curl-examples.md), [Python](./user-guide/api/examples/python-examples.md), [Node.js](./user-guide/api/examples/node-examples.md)
- **Compliance:** [DORA](./compliance/COMPLIANCE_DORA.md), [NIS2](./compliance/COMPLIANCE_NIS2.md), [AI Act](./compliance/COMPLIANCE_RUNTIME_AI.md)
- **How it works:** [Algorithms](./architecture/ALGORITHMS.md), [Architecture](./architecture/API.md)

## Pricing (snapshot)

| Plan | Price | Events/month | Best for |
| --- | --- | --- | --- |
| Pilot | Free | 10,000 | Testing & prototyping |
| Radar | $20/mo | 500,000 | Production monitoring |
| Lock | $200/mo | 5,000,000 | Enterprise with compliance |

## Support

- Email: [support@driftlock.io](mailto:support@driftlock.io)
- GitHub: [Shannon-Labs/driftlock](https://github.com/Shannon-Labs/driftlock)
- Community: [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)
