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

# Build the Rust API
cargo build -p driftlock-api --release

# Start PostgreSQL (Docker)
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15

# Run the API server
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  ./target/release/driftlock-api
```

## Your First API Call

Once you have your API key, you can start detecting anomalies:

```bash
# Replace YOUR_API_KEY with your actual key
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      "2024-01-01T00:00:00Z INFO Normal operation",
      "2024-01-01T00:00:01Z INFO Normal operation",
      "2024-01-01T00:10:00Z ERROR CRITICAL: Unusual pattern!"
    ]
  }'
```

### Try the Demo (No Auth Required)

```bash
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "normal log entry 1",
      "normal log entry 2",
      "ERROR: unusual event detected"
    ]
  }'
```

## Understanding the Response

```json
{
  "anomalies": [
    {
      "id": "anom_abc123",
      "ncd": 0.72,
      "compression_ratio": 1.41,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "explanation": "Significant deviation from baseline pattern"
    }
  ],
  "metrics": {
    "processed": 100,
    "baseline": 400,
    "window": 50,
    "duration_ms": 42
  }
}
```

### Key Metrics Explained

- **NCD (Normalized Compression Distance)**: 0-1 scale. Higher values = more anomalous
- **p_value**: Statistical significance. Lower = more confident the anomaly is real
- **confidence**: Overall confidence level (0-1)
- **compression_ratio**: How compressible the data is

## API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/healthz` | GET | None | Liveness check |
| `/v1/demo/detect` | POST | IP Rate | Demo detection (10 req/hr) |
| `/v1/detect` | POST | API Key | Authenticated detection |
| `/v1/anomalies` | GET | API Key | List detected anomalies |
| `/v1/anomalies/:id` | GET | API Key | Get anomaly details |
| `/v1/streams` | GET | API Key | List streams |
| `/v1/profiles` | GET | API Key | List detection profiles |

## Configuration Options

You can customize detection settings per request:

```json
{
  "stream_id": "my-stream",
  "events": ["..."],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "hop_size": 10,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd"
  }
}
```

### Detection Profiles

Choose a pre-configured profile or use custom settings:

| Profile | NCD Threshold | P-Value | Use Case |
|---------|---------------|---------|----------|
| `sensitive` | 0.20 | 0.10 | Security-critical, early warning |
| `balanced` | 0.30 | 0.05 | General purpose (default) |
| `strict` | 0.45 | 0.01 | Low noise, high confidence |
| `custom` | User-defined | User-defined | Fine-tuned settings |

## Best Practices

### 1. Provide Enough Baseline Data

For accurate detection, provide at least 400 baseline events before the window you want to analyze.

### 2. Use Consistent Event Format

Keep your event format consistent:

```json
{
  "events": [
    "2024-01-01T00:00:00Z INFO User logged in",
    "2024-01-01T00:00:01Z INFO Page viewed",
    "2024-01-01T00:00:02Z ERROR Database connection failed"
  ]
}
```

### 3. Stream Organization

Organize data by type using different streams:
- **logs**: Application logs
- **metrics**: Numerical metrics
- **traces**: Distributed traces
- **llm**: LLM responses/prompts

### 4. Tune Thresholds

Start with defaults and adjust based on your data:
- Higher `ncd_threshold` = fewer false positives
- Lower `p_value_threshold` = stricter statistical test

## Plan Limits

### Free Tier
- **Events**: 10,000/month
- **Streams**: 5
- **Retention**: 14 days

### Pro ($99/mo)
- **Events**: 500,000/month
- **Streams**: 20
- **Retention**: 90 days

### Team ($199/mo)
- **Events**: 5,000,000/month
- **Streams**: 100
- **Retention**: 1 year

### Enterprise (Custom)
- **Events**: Unlimited
- **Streams**: 500+
- **Retention**: Custom
- **EU Data Residency**: Available
- **Self-Hosting**: Available

Need more? [View pricing](https://driftlock.net/#pricing)

## Common Issues

### "Invalid API key"

- Ensure you're using the full key including the `dlk_` prefix
- Check that the key hasn't been revoked
- Use `X-Api-Key` header (not `Authorization: Bearer`)

### "events required"

- Provide at least one event in the `events` array
- Events should be strings

### "no stream configured"

- Use `stream_id: "default"` or create a stream first via `POST /v1/streams`

### Low anomaly detection

- Ensure you have enough baseline events (400+)
- Check that your baseline represents "normal" behavior
- Adjust `ncd_threshold` if needed
- Try the `sensitive` detection profile

### Rate limit exceeded (429)

- Demo endpoint: 10 requests per hour per IP
- Authenticated endpoints: Based on your plan tier

## Driftlog (Debug Logging)

Enable detailed logging for troubleshooting:

```bash
# Run with debug logging
RUST_LOG=debug cargo run -p driftlock-api

# Filter to specific modules
RUST_LOG=driftlock_api=debug,driftlock_db=info cargo run -p driftlock-api
```

## Next Steps

1. **Explore the API**: Check out the [API Reference](../../architecture/API.md)
2. **Configure Detection**: Learn about [Detection Profiles](../guides/detection-profiles.md)
3. **Deploy to Production**: Follow the [Deployment Guide](../../deployment/DEPLOYMENT.md)
4. **Join the Community**: Star us on [GitHub](https://github.com/Shannon-Labs/driftlock)

## Support

- **Email**: hunter@shannonlabs.dev
- **GitHub Issues**: [Report a bug](https://github.com/Shannon-Labs/driftlock/issues)
- **Documentation**: [Full docs](https://docs.driftlock.io)

---

Built by [Shannon Labs](https://shannonlabs.dev) | [Apache 2.0 License](../../../LICENSE)
