# Driftlock Documentation Map

Everything you need to ship Driftlock in production—quickstart, core concepts, API reference, tutorials, operations, and compliance.

## Start Here

- **[Quickstart (5 minutes)](./getting-started/quickstart.md)**
- **[Authentication](./getting-started/authentication.md)** — API keys, headers, environment setup, troubleshooting
- **[Core Concepts](./getting-started/concepts.md)** — baselines, windows, p-values, sensitivity profiles, auto-tuning
- **Dashboard Upload & Analyze (in-app)** — Upload JSON/NDJSON and call `/v1/detect` from the dashboard (login required)

## API Reference

- **[REST API overview](./api/rest-api.md)** — base URL, headers, rate limits, pagination
- **Endpoints:** [POST /v1/detect](./api/endpoints/detect.md), [POST /v1/demo/detect](./api/endpoints/demo.md), [GET /v1/anomalies](./api/endpoints/anomalies.md)
- **Errors:** [Standard error codes](./api/errors.md)
- **Examples:** [cURL](./api/examples/curl-examples.md), [Python](./api/examples/python-examples.md)

## Tutorials & Examples

- **[Log monitoring](./tutorials/log-monitoring.md)** — parse JSON logs, flag suspicious entries
- **[Metrics spike detection](./tutorials/metrics-spike-detection.md)** — catch latency/cpu spikes
- **[LLM output drift](./tutorials/llm-output-drift.md)** — monitor hallucinations and format drift
- **[Webhook + alerting](./tutorials/webhook-alerting.md)** — send anomalies to Slack/PagerDuty
- **[Feedback loop](./tutorials/feedback-loop.md)** — capture false positives/confirmations
- **[Profiles tutorial](./tutorials/profiles-tutorial.md)** — choose and adjust sensitivity

## Operations & SRE

- **[Operations & Runbooks](./guides/operations.md)** — deployment modes, health/readiness checks, scaling, backups
- **[Troubleshooting](./guides/troubleshooting.md)** — "no anomalies", JSON errors, performance, auth issues
- **[Deployment runbooks](../deployment/RUNBOOKS.md)** — operational playbooks
- **Cloud Run/Firebase:** [Cloud Run setup](../deployment/cloud-run-setup.md), [Firebase hosting](../deployment/firebase-hosting-setup.md)

## Compliance & Assurance

- **[DORA](../compliance/COMPLIANCE_DORA.md)**
- **[NIS2](../compliance/COMPLIANCE_NIS2.md)**
- **[AI Act runtime monitoring](../compliance/COMPLIANCE_RUNTIME_AI.md)**

## Architecture & Concepts

- **[Algorithms](../architecture/ALGORITHMS.md)** — compression/NCD math
- **[API architecture](../architecture/API.md)** — data flow and components

## Support

- **Email:** [support@driftlock.io](mailto:support@driftlock.io)
- **GitHub:** [Shannon-Labs/driftlock](https://github.com/Shannon-Labs/driftlock)
- **Community:** [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)
