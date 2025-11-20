# Feature Comparison

How Driftlock stacks up against traditional observability tools.

## The Landscape

| Feature | Driftlock | Datadog / New Relic | Splunk / ELK |
|---------|-----------|---------------------|--------------|
| **Core Tech** | Compression (NCD) | Statistical / ML | Search Index |
| **Setup Time** | Minutes | Hours/Days | Days/Weeks |
| **Maintenance** | Zero | High (Rules/Configs) | High (Index Mgmt) |
| **Data Type** | Any JSON Event | Metrics / Traces | Logs |
| **Pricing** | Per Event | Per Host + Extras | Ingestion Volume |

## Why Driftlock?

### 1. "It Just Works"
Traditional tools require you to know what you're looking for. You have to set thresholds ("Alert if CPU > 80%"). Driftlock learns what is normal for *your* specific data stream without manual configuration.

### 2. Explainable AI
Black-box ML tools often give you a "magic score" with no context. Driftlock tells you *why* something is anomalous by pointing to the specific fields that caused the compression distance to spike.

### 3. Developer First
We are built for developers.
- **CLI**: Test locally.
- **API**: Integrate into your code.
- **Webhooks**: Automate your response.

## When NOT to use Driftlock

Driftlock is specialized for **Anomaly Detection**. It is not a replacement for:
- **Long-term Log Storage**: Use Splunk/S3/CloudWatch.
- **Distributed Tracing**: Use Jaeger/Zipkin/Datadog APM.
- **Infrastructure Metrics**: Use Prometheus/Grafana.

**Best Practice**: Use Driftlock as an intelligent layer *on top* of your existing stack to catch the "unknown unknowns".
