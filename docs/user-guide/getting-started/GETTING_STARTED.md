# Getting Started with Driftlock

Welcome to Driftlock! This guide will help you get up and running with explainable anomaly detection in minutes.

## Quick Start

### Option 1: Sign Up for Free Trial

1. Visit [driftlock.net](https://driftlock.net)
2. Click "Get Started" and create an account
3. Save your API key (it won't be shown again!)
4. Make your first API call

### Option 2: Self-Hosted Setup

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Build the Rust core
cd cbad-core && cargo build --release && cd ..

# Run the API demo (includes Docker + Postgres)
DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
```

## Your First API Call

Once you have your API key, you can start detecting anomalies:

```bash
# Replace YOUR_API_KEY with your actual key
curl -X POST https://driftlock.net/api/v1/detect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2024-01-01T00:00:00Z", "level": "info", "message": "Normal operation"},
      {"timestamp": "2024-01-01T00:00:01Z", "level": "info", "message": "Normal operation"},
      // ... add more events (minimum ~450 for baseline)
      {"timestamp": "2024-01-01T00:10:00Z", "level": "error", "message": "CRITICAL: Unusual pattern!"}
    ]
  }'
```

## Understanding the Response

```json
{
  "success": true,
  "batch_id": "550e8400-e29b-41d4-a716-446655440000",
  "stream_id": "default",
  "total_events": 451,
  "anomaly_count": 1,
  "processing_time": "245.3ms",
  "compression_algo": "zstd",
  "anomalies": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "index": 450,
      "metrics": {
        "ncd": 0.89,
        "p_value": 0.002,
        "confidence": 0.95,
        "compression_ratio": 1.45
      },
      "why": "High NCD (0.89 > 0.30) indicates significant deviation from baseline. Statistical significance confirmed (p=0.002 < 0.05).",
      "detected": true
    }
  ]
}
```

To make `zlab` the default compressor without Docker, set an environment variable before running the Go binary directly:

```bash
DEFAULT_ALGO=zlab go run ./collector-processor/cmd/driftlock-http
```

### Key Metrics Explained

- **NCD (Normalized Compression Distance)**: 0-1 scale. Higher values = more anomalous
- **p_value**: Statistical significance. Lower = more confident the anomaly is real
- **confidence**: Overall confidence level (0-1)
- **compression_ratio**: How compressible the data is

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Health check |
| `/v1/detect` | POST | Detect anomalies in events |
| `/v1/anomalies` | GET | List detected anomalies |
| `/v1/anomalies/{id}` | GET | Get anomaly details |
| `/v1/anomalies/export` | POST | Export anomalies |

## Configuration Options

You can customize detection settings per request:

```json
{
  "stream_id": "default",
  "events": [...],
  "config_override": {
    "baseline_size": 400,      // Number of events for baseline
    "window_size": 50,         // Window size for comparison
    "hop_size": 10,            // Slide window by this amount
    "ncd_threshold": 0.3,      // NCD threshold for anomaly
    "p_value_threshold": 0.05, // Statistical significance threshold
    "compressor": "zstd"       // zlab, zstd, lz4, or gzip
  }
}
```

To default the HTTP service to the new `zlab` compressor without Docker, set the environment variable before starting the Go binary directly:

```bash
DEFAULT_ALGO=zlab go run ./collector-processor/cmd/driftlock-http
```

## Best Practices

### 1. Provide Enough Baseline Data

For accurate detection, provide at least 400 baseline events before the window you want to analyze.

### 2. Use Consistent Event Structure

Keep your event structure consistent across baseline and test data:

```json
{
  "timestamp": "ISO 8601 format",
  "level": "info|warn|error",
  "message": "Description",
  "custom_field": "value"
}
```

### 3. Stream Organization

Organize data by type:
- **logs**: Application logs
- **metrics**: Numerical metrics
- **traces**: Distributed traces
- **llm**: LLM responses/prompts

### 4. Tune Thresholds

Start with defaults and adjust based on your data:
- Higher `ncd_threshold` = fewer false positives
- Lower `p_value_threshold` = stricter statistical test

## Free Trial Limits

- **Events**: 10,000/month
- **Duration**: 14 days
- **Streams**: 1
- **Retention**: 14 days

Need more? [View pricing](https://driftlock.net/#pricing)

## Common Issues

### "Invalid API key"

- Ensure you're using the full key including the `dlk_` prefix
- Check that the key hasn't been revoked

### "events required"

- Provide at least one event in the `events` array
- Events must be valid JSON objects

### "no stream configured"

- Use `stream_id: "default"` or the specific stream ID from your account

### Low anomaly detection

- Ensure you have enough baseline events (400+)
- Check that your baseline represents "normal" behavior
- Adjust `ncd_threshold` if needed

## Next Steps

1. **Explore the API**: Check out the [API Documentation](API.md)
2. **Run the Demo**: Try the [API Demo Walkthrough](API-DEMO-WALKTHROUGH.md)
3. **Deploy to Production**: Follow the [Deployment Guide](COMPLETE_DEPLOYMENT_PLAN.md)
4. **Join the Community**: Star us on [GitHub](https://github.com/Shannon-Labs/driftlock)

## Support

- **Email**: hunter@shannonlabs.dev
- **GitHub Issues**: [Report a bug](https://github.com/Shannon-Labs/driftlock/issues)
- **Documentation**: [Full docs](https://driftlock.net/docs)

---

Built by [Shannon Labs](https://shannonlabs.dev) | [Apache 2.0 License](../LICENSE)
