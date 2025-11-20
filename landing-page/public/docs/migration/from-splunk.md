# Migrating from Splunk

Moving from Splunk's heavy log analysis to Driftlock's real-time anomaly detection?

## The Paradigm Shift

Splunk is a **Search Engine** for logs. You index everything and search later.
Driftlock is a **Detection Engine**. You stream events, and we tell you if they are weird immediately.

## Migration Steps

### 1. Data Ingestion

Instead of sending logs to a Splunk Forwarder, send them to the Driftlock API.

**Splunk (inputs.conf):**
```ini
[monitor:///var/log/app.log]
index = main
```

**Driftlock (CLI):**
```bash
tail -f /var/log/app.log | driftlock stream --stream-id app-logs
```

### 2. Replacing Alerts

Splunk alerts are based on search queries (SPL).

**Splunk SPL:**
```splunk
index=main status=500 | stats count by host | where count > 10
```

**Driftlock:**
You don't write queries. Driftlock automatically detects if the count of 500 errors is higher than normal for that host.

### 3. Cost Savings

Splunk licensing is typically based on **Ingestion Volume (GB/day)**. This encourages you to log less.

Driftlock is based on **Events Processed**. We encourage you to send rich, structured data because it helps our compression models learn better.

## Common Patterns

### Hybrid Approach
Keep Splunk for long-term compliance storage and forensic search. Use Driftlock for real-time operational intelligence and alerting.

1. App writes to Log File.
2. Splunk Forwarder reads Log File -> Splunk Indexer.
3. Driftlock CLI reads Log File -> Driftlock API.

This gives you the best of both worlds without doubling your Splunk license cost.
