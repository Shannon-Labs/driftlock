# Driftlock Documentation

Welcome to Driftlock - **compression-based anomaly detection** for your logs, metrics, traces, and LLM I/O.

## Quick Links

| I want to... | Go to... |
|--------------|----------|
| Try it now (no signup) | [Demo Endpoint](./user-guide/api/endpoints/demo.md) |
| Get started in 5 minutes | [Quickstart Guide](./user-guide/getting-started/quickstart.md) |
| See API endpoints | [REST API Reference](./user-guide/api/rest-api.md) |
| Understand how it works | [Core Concepts](./user-guide/getting-started/concepts.md) |
| Run a tutorial | [Tutorials & Examples](./user-guide/tutorials/README.md) |
| Operate in production | [Operations & Runbooks](./user-guide/guides/operations.md) |
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
| **Free** | $0/mo | 10,000 | Experimentation |
| **Standard** | $15/mo | 500,000 | Startups & active monitoring |
| **Pro** | $100/mo | 5,000,000 | High-volume production |
| **Enterprise** | $299/mo | 25,000,000 | Large-scale compliance |

[Sign up free →](https://driftlock.net/#signup)

---

## Documentation Sections

### Getting Started

- [Quickstart Guide](./user-guide/getting-started/quickstart.md) - First detection in 5 minutes
- [Core Concepts](./user-guide/getting-started/concepts.md) - How compression-based detection works
- [Authentication](./user-guide/getting-started/authentication.md) - API keys and access
- Dashboard Upload & Analyze (in-app) - Run `/v1/detect` from the dashboard after login

### API Reference

- [REST API Overview](./user-guide/api/rest-api.md) - Base URL, headers, rate limits
- [POST /v1/detect](./user-guide/api/endpoints/detect.md) - Main detection endpoint
- [POST /v1/demo/detect](./user-guide/api/endpoints/demo.md) - Try without auth
- [GET /v1/anomalies](./user-guide/api/endpoints/anomalies.md) - Query your anomalies
- [Error Codes](./user-guide/api/errors.md) - Error handling reference

### Tutorials & Examples

- [Log monitoring](./user-guide/tutorials/log-monitoring.md) - Detect suspicious log lines
- [Metrics spike detection](./user-guide/tutorials/metrics-spike-detection.md) - Catch latency/CPU spikes
- [LLM output drift](./user-guide/tutorials/llm-output-drift.md) - Monitor model output changes
- [Webhook/alerting](./user-guide/tutorials/webhook-alerting.md) - Send anomalies to Slack/PagerDuty
- [Feedback loop](./user-guide/tutorials/feedback-loop.md) - Improve sensitivity with feedback
- [Profiles tutorial](./user-guide/tutorials/profiles-tutorial.md) - Choose the right sensitivity
- [cURL examples](./user-guide/api/examples/curl-examples.md) - Command line usage
- [Python examples](./user-guide/api/examples/python-examples.md) - Python integration

### Operations & SRE

- [Operations & Runbooks](./user-guide/guides/operations.md) - Health checks, scaling, backup/restore
- [Deployment runbooks](./deployment/RUNBOOKS.md) - Operational playbooks

### Compliance

- [DORA Compliance](./compliance/COMPLIANCE_DORA.md) - EU financial regulations
- [NIS2 Compliance](./compliance/COMPLIANCE_NIS2.md) - EU cybersecurity directive
- [AI Act Compliance](./compliance/COMPLIANCE_RUNTIME_AI.md) - Runtime AI monitoring

---

## Support

- **Email**: [support@driftlock.io](mailto:support@driftlock.io)
- **Documentation**: You're here!

---

**Ready to get started?** [Sign up free](https://driftlock.net/#signup) or [try the demo](./user-guide/api/endpoints/demo.md).
