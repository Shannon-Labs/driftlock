# Quickstart Guide

Get started with Driftlock in 5 minutes. This guide will walk you through signing up, getting your API key, and running your first anomaly detection.

## Try Without Signing Up (Optional)

Want to test immediately? Use our [demo endpoint](../api/endpoints/demo.md) - no signup required:

```bash
curl -X POST https://api.driftlock.net/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"body": {"latency": 120}},
      {"body": {"latency": 125}},
      {"body": {"latency": 118}},
      {"body": {"latency": 950}}
    ]
  }'
```

Ready for full access? Continue below.

---

## Step 1: Create Your Account

1. Visit [https://driftlock.net](https://driftlock.net)
2. Scroll to the **Sign Up** section or click "Get Started"
3. Enter your email and company name
4. Check your email for a verification link
5. Click the link to verify and activate your account

**Pilot Plan**: You'll start on the free Pilot plan with 10,000 events/month.

## Step 2: Get Your API Key

After verifying your email, you'll receive your API key immediately on the verification page. You can also find it in your dashboard:

1. Log in to your [Dashboard](https://driftlock.net/dashboard)
2. Your API key is displayed in the **Quick Start** section
3. Click the copy button to copy it

> **Important**: Your API key is shown once during signup. Store it securely! You can regenerate a new one from the dashboard if needed.

## Step 3: Run Your First Detection

Let's detect anomalies in a simple dataset. Copy this cURL command and replace `YOUR_API_KEY`:

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
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

1. Go back to your [Dashboard](https://driftlock.net/dashboard)
2. View your usage statistics and recent activity
3. Check your detected anomalies in the feed
4. Click on an anomaly to see detailed metrics and explanations

## What's Next?

Now that you've run your first detection, you can:

- **[Understand Core Concepts](./concepts.md)** - Learn about NCD, baselines, and how the detection works
- **[Explore the REST API](../api/endpoints/detect.md)** - Full API reference with all endpoints
- **[Try the Demo Endpoint](../api/endpoints/demo.md)** - Test without authentication
- **[View Code Examples](../api/examples/python-examples.md)** - Python and cURL examples
- **[Error Codes Reference](../api/errors.md)** - Handle errors properly

## Common Issues

### "unauthorized" error
- Check that your API key is correct
- Verify you're including the `X-Api-Key` header
- Ensure your API key hasn't been revoked
- See [Error Codes](../api/errors.md#unauthorized) for details

### "rate_limit_exceeded" error
- Pilot plan is limited to 60 requests/minute
- Implement exponential backoff
- Consider upgrading to Radar ($20/mo) for higher limits
- See [Error Codes](../api/errors.md#rate_limit_exceeded) for retry logic

### Need Help?

- **Documentation**: [Full API Reference](../api/endpoints/detect.md)
- **Support**: support@driftlock.io

---

**Ready to build?** Check out our [Python examples](../api/examples/python-examples.md) for integration patterns.
