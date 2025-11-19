# Quickstart Guide

Get started with Driftlock in 5 minutes. This guide will walk you through signing up, getting your API key, and running your first anomaly detection.

## Step 1: Create Your Account

1. Visit [https://driftlock.web.app](https://driftlock.web.app)
2. Click **"Sign Up"** in the navigation
3. Sign up with your email or Google account (Firebase Auth)
4. Verify your email if required

**Developer Plan**: You'll start on the free Developer plan with 10,000 events/month.

## Step 2: Get Your API Key

1. After logging in, go to your [Dashboard](https://driftlock.web.app/dashboard)
2. Navigate to **"API Keys"** section
3. Click **"Create API Key"**
4. Give it a name (e.g., "My First Key")
5. Copy your API key - you'll need this for authentication

> ⚠️ **Important**: Store your API key securely. It won't be shown again after you close the dialog.

## Step 3: Run Your First Detection

Let's detect anomalies in a simple dataset. Copy this cURL command and replace `YOUR_API_KEY`:

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"latency": 120}},
      {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"latency": 125}},
      {"timestamp": "2025-01-01T10:02:00Z", "type": "metric", "body": {"latency": 118}},
      {"timestamp": "2025-01-01T10:03:00Z", "type": "metric", "body": {"latency": 950}}
    ]
  }'
```

### Expected Response

You'll receive a JSON response with detected anomalies:

```json
{
  "success": true,
  "batch_id": "batch_abc123",
  "stream_id": "...",
  "total_events": 4,
  "anomaly_count": 1,
  "processing_time": "125ms",
  "compression_algo": "zstd",
  "anomalies": [
    {
      "id": "anom_xyz789",
      "index": 3,
      "metrics": {
        "ncd": 0.72,
        "compression_ratio": 1.41,
        "entropy_change": 0.13,
        "p_value": 0.004,
        "confidence": 0.96
      },
      "event": {"latency": 950},
      "why": "Significant latency spike detected",
      "detected": true
    }
  ],
  "request_id": "req_123"
}
```

## Step 4: Understanding the Results

Key metrics in the response:

- **`ncd`** (Normalized Compression Distance): How different this event is from the baseline (0-1, higher = more anomalous)
- **`p_value`**: Statistical significance (< 0.05 is typically significant)
- **`confidence`**: How confident we are this is an anomaly (0-1, higher = more confident)
- **`why`**: Plain English explanation of the anomaly

## Step 5: View in Dashboard

1. Go back to your [Dashboard](https://driftlock.web.app/dashboard)
2. Click **"Anomalies"** in the sidebar
3. You'll see all detected anomalies with their metrics
4. Click on an anomaly to see detailed evidence and explanations

## What's Next?

Now that you've run your first detection, you can:

- **[Understand Core Concepts](./concepts.md)** - Learn about NCD, baselines, and how the detection works
- **[Explore the REST API](../api/rest-api.md)** - Full API reference with all endpoints
- **[Set Up Authentication](./authentication.md)** - Manage API keys and Firebase auth
- **[Try GraphQL](../graphql/overview.md)** - Use Firebase Data Connect for richer queries
- **[View Code Examples](../api/examples/python-examples.md)** - Python, Node.js, and other language examples

## Common Issues

### "unauthorized" error
- Check that your API key is correct
- Verify you're including the `X-Api-Key` header
- Ensure your API key hasn't been revoked

### "rate_limit_exceeded" error
- Developer plan is limited to 60 requests/minute
- Implement exponential backoff
- Consider upgrading to Starter plan ($25/mo) for higher limits

### Need Help?

- **Documentation**: [Full API Reference](../api/rest-api.md)
- **Support**: support@driftlock.io
- **Community**: [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)

---

**Ready to build?** Check out our [tutorials](../tutorials/) for step-by-step guides on specific use cases like financial monitoring, log analysis, and IoT telemetry.
