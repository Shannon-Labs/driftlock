# Driftlock Documentation

Welcome to Driftlock - **compression-based anomaly detection** for your logs, metrics, traces, and LLM I/O.

## Quick Links

| I want to... | Go to... |
|--------------|----------|
| Try it now (no signup) | [Demo Endpoint](./user-guide/api/endpoints/demo.md) |
| Get started in 5 minutes | [Quickstart Guide](./user-guide/getting-started/quickstart.md) |
| See API endpoints | [API Reference](./user-guide/api/endpoints/detect.md) |
| Understand how it works | [Core Concepts](./user-guide/getting-started/concepts.md) |
| Handle errors | [Error Codes](./user-guide/api/errors.md) |

---

## What is Driftlock?

Driftlock detects anomalies by analyzing how well your data compresses. When new data compresses poorly against a learned baseline, it's likely anomalous.

**Key features:**
- **No training required** - Works immediately on any JSON data
- **Explainable results** - Every anomaly includes metrics and plain English explanations
- **Low false positives** - Statistical significance testing (p-values) reduces noise
- **Compliance-ready** - Evidence bundles for DORA, NIS2, and AI Act

## Try It Now

No signup required. Run this in your terminal:

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

## Pricing

| Plan | Price | Events/month | Best for |
|------|-------|--------------|----------|
| **Pilot** | Free | 10,000 | Testing and prototyping |
| **Radar** | $20/mo | 500,000 | Production monitoring |
| **Lock** | $200/mo | 5,000,000 | Enterprise with compliance |

[Sign up free →](https://driftlock.net/#signup)

---

## Documentation Sections

### Getting Started
- [Quickstart Guide](./user-guide/getting-started/quickstart.md) - First detection in 5 minutes
- [Core Concepts](./user-guide/getting-started/concepts.md) - How compression-based detection works
- [Authentication](./user-guide/getting-started/authentication.md) - API keys and access

### API Reference
- [POST /v1/detect](./user-guide/api/endpoints/detect.md) - Main detection endpoint
- [POST /v1/demo/detect](./user-guide/api/endpoints/demo.md) - Try without auth
- [GET /v1/anomalies](./user-guide/api/endpoints/anomalies.md) - Query your anomalies
- [Error Codes](./user-guide/api/errors.md) - Error handling reference

### Code Examples
- [cURL Examples](./user-guide/api/examples/curl-examples.md) - Command line usage
- [Python Examples](./user-guide/api/examples/python-examples.md) - Python integration

### Compliance
- [DORA Compliance](./compliance/COMPLIANCE_DORA.md) - EU financial regulations
- [NIS2 Compliance](./compliance/COMPLIANCE_NIS2.md) - EU cybersecurity directive
- [AI Act Compliance](./compliance/COMPLIANCE_RUNTIME_AI.md) - Runtime AI monitoring

---

## Support

- **Email**: support@driftlock.io
- **Documentation**: You're here!

---

**Ready to get started?** [Sign up free](https://driftlock.net/#signup) or [try the demo](./user-guide/api/endpoints/demo.md).
