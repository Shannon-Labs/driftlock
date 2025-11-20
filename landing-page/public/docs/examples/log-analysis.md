# Log Analysis

Detect anomalies in your application logs without writing complex regex rules.

## Scenario

You have a microservice generating logs. You want to know when the logs start looking "weird" (e.g., new error messages, stack traces, or unusual access patterns).

## Implementation

### 1. Log Structure

We will analyze structured JSON logs from an Nginx server.

```json
{
  "timestamp": "2025-01-01T12:00:00Z",
  "remote_addr": "192.168.1.1",
  "request": "GET /api/v1/users HTTP/1.1",
  "status": 200,
  "body_bytes_sent": 1024,
  "http_referer": "https://example.com",
  "http_user_agent": "Mozilla/5.0..."
}
```

### 2. Streaming Logs to Driftlock

We can use the Driftlock CLI to pipe logs directly from a file or process.

```bash
# Tail the log file and pipe to Driftlock
tail -f /var/log/nginx/access.log | driftlock stream --stream-id nginx-logs --format nginx
```

Or programmatically via a log shipper (e.g., Fluentd, Logstash) or a simple script:

```python
import time
import json
from driftlock import DriftlockClient

client = DriftlockClient(api_key="...")

def tail_logs(filepath):
    with open(filepath, "r") as f:
        f.seek(0, 2) # Go to end
        while True:
            line = f.readline()
            if not line:
                time.sleep(0.1)
                continue
            yield json.loads(line)

async def monitor_logs():
    for log_entry in tail_logs("/var/log/app.json"):
        # Detect anomalies
        result = await client.detect(
            stream_id="app-logs", 
            events=[{"type": "log", "body": log_entry}]
        )
        
        if result.anomalies:
            print(f"[ALERT] Anomaly in logs: {result.anomalies[0].why}")
```

## What Driftlock Detects

### New Error Patterns
If your application normally logs `INFO: User logged in`, and suddenly starts logging `ERROR: Database connection failed`, Driftlock will flag this immediately because the compression distance between the new error and the baseline is high.

### Unusual Payloads
If a specific field in your logs (e.g., `response_time`) suddenly jumps from `50ms` to `5000ms`, or if a `user_agent` string looks like a SQL injection attack, it will be detected.

### Volume Spikes
While Driftlock focuses on *content* anomalies, significant changes in log volume or structure (e.g., a flood of similar requests) will also manifest as anomalies in the aggregate stream.

## Alerting

Connect Driftlock to a [Webhook](../tools/webhooks.md) to send these alerts to Slack.

**Slack Message Example:**
> 🚨 **Log Anomaly Detected**
> **Stream**: nginx-logs
> **Confidence**: 99%
> **Reason**: Unexpected status code pattern (500)
> **Log Entry**: `{"status": 500, "request": "POST /login"}`
