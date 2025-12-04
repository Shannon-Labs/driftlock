# General Anomaly Detection with Driftlock

Driftlock uses **Compression-Based Anomaly Detection (CBAD)** to find unusual patterns in any JSON event stream. Unlike rule-based systems, CBAD works without predefined thresholds.

## How It Works

1. **Baseline Window**: Recent events form a "normal" baseline
2. **Compression**: Events are compressed using zstd/lz4
3. **NCD Calculation**: Normalized Compression Distance measures similarity
4. **Statistical Testing**: Permutation testing determines significance
5. **Anomaly Flagging**: Events with high NCD and low p-value are anomalous

## Use Case 1: Log Anomaly Detection

Find unusual log entries without defining patterns in advance.

### The Problem
Your application generates thousands of logs. An unusual error appears, but it's buried in noise. Traditional regex rules miss novel patterns.

### The Solution

```bash
curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "app-logs",
    "events": [
      {"level": "info", "message": "Request processed", "service": "api"},
      {"level": "info", "message": "Request processed", "service": "api"},
      {"level": "info", "message": "Request processed", "service": "api"},
      {"level": "error", "message": "PANIC: nil pointer dereference", "service": "api"}
    ]
  }'
```

### What You Get

```json
{
  "anomaly_count": 1,
  "anomalies": [{
    "index": 3,
    "why": "High NCD (0.68) - log structure differs significantly from baseline",
    "metrics": {
      "ncd": 0.68,
      "p_value": 0.002,
      "confidence_level": 0.998
    }
  }]
}
```

## Use Case 2: Metric Spike Detection

Determine if a metric change is statistically significant or just noise.

### The Problem
CPU jumps from 45% to 85%. Is this a real problem or normal variance?

### The Solution

```json
{
  "stream_id": "system-metrics",
  "events": [
    {"cpu": 45, "memory": 62, "timestamp": "2024-01-15T10:00:00Z"},
    {"cpu": 48, "memory": 63, "timestamp": "2024-01-15T10:01:00Z"},
    {"cpu": 47, "memory": 61, "timestamp": "2024-01-15T10:02:00Z"},
    {"cpu": 85, "memory": 78, "timestamp": "2024-01-15T10:03:00Z"}
  ]
}
```

### Interpretation

- **NCD = 0.45**: Moderate deviation
- **p-value = 0.03**: Statistically significant (< 0.05 threshold)
- **Verdict**: This is a real anomaly, not noise

## Use Case 3: API Traffic Analysis

Detect unusual API request patterns that might indicate abuse or attacks.

### The Problem
Someone is probing your API with unusual requests. Traditional WAF rules don't catch it.

### The Solution

```json
{
  "stream_id": "api-requests",
  "events": [
    {"endpoint": "/api/users", "method": "GET", "status": 200, "ip": "192.168.1.1"},
    {"endpoint": "/api/users", "method": "GET", "status": 200, "ip": "192.168.1.2"},
    {"endpoint": "/api/admin/delete-all", "method": "DELETE", "status": 403, "ip": "10.0.0.99"},
    {"endpoint": "/.env", "method": "GET", "status": 404, "ip": "10.0.0.99"}
  ]
}
```

### What Gets Flagged

- Unusual endpoint patterns (`/.env`, `/admin/delete-all`)
- Unusual HTTP methods for endpoints
- Unusual status code patterns

## Use Case 4: Time Series Anomaly Detection

Detect anomalies in any time-ordered data.

### Financial Data

```json
{
  "events": [
    {"symbol": "AAPL", "price": 185.50, "volume": 1000000},
    {"symbol": "AAPL", "price": 186.20, "volume": 1200000},
    {"symbol": "AAPL", "price": 145.00, "volume": 50000000}
  ]
}
```

### Sensor Data

```json
{
  "events": [
    {"sensor_id": "temp-01", "value": 72.5, "unit": "F"},
    {"sensor_id": "temp-01", "value": 73.1, "unit": "F"},
    {"sensor_id": "temp-01", "value": 150.0, "unit": "F"}
  ]
}
```

## Configuration Tips

### For High-Frequency Data
```json
{
  "config_override": {
    "baseline_size": 1000,
    "window_size": 100,
    "hop_size": 50
  }
}
```

### For Sensitive Detection
```json
{
  "config_override": {
    "ncd_threshold": 0.2,
    "p_value_threshold": 0.01
  }
}
```

### For Noisy Data
```json
{
  "config_override": {
    "ncd_threshold": 0.5,
    "p_value_threshold": 0.1,
    "permutation_count": 500
  }
}
```

## When to Use Driftlock

**Good Fit:**
- Log analysis where patterns aren't known in advance
- Metric monitoring where thresholds are hard to set
- Security monitoring for novel attack patterns
- Any JSON event stream with temporal ordering

**Not Ideal For:**
- Single-value threshold monitoring (use Prometheus)
- Image/video anomaly detection
- Unstructured text without JSON wrapping

## Getting Started

1. **Try the demo** (no signup): `POST /v1/demo/detect`
2. **Sign up** for an API key
3. **Stream your events** to `/v1/detect`
4. **Query anomalies** at `/v1/anomalies`

## API Reference

- Base URL: `https://driftlock.net/api`
- OpenAPI: `https://driftlock.net/docs/architecture/api/openapi.yaml`
- Full docs: `https://driftlock.net/docs`
