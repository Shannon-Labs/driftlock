# Webhooks

Receive real-time notifications when anomalies are detected by configuring webhooks.

## Overview

Webhooks allow Driftlock to push anomaly data to your systems immediately upon detection. This is useful for:
- Triggering alerts in PagerDuty or Slack
- Automating remediation workflows
- Logging anomalies to external systems

## Configuration

1. Go to the **[Integrations](https://driftlock.web.app/dashboard/integrations)** section in your dashboard.
2. Click **"Add Webhook"**.
3. Enter your **Payload URL** (the endpoint on your server).
4. (Optional) Enter a **Secret** to sign requests.
5. Select the events you want to subscribe to (currently only `anomaly.detected`).

## Payload Format

When an anomaly is detected, Driftlock sends a `POST` request to your URL with a JSON payload:

```json
{
  "id": "evt_1234567890",
  "type": "anomaly.detected",
  "created_at": "2025-01-01T12:00:00Z",
  "data": {
    "stream_id": "payment-service",
    "anomaly_id": "anom_abc123",
    "metrics": {
      "ncd": 0.85,
      "confidence": 0.98,
      "p_value": 0.001
    },
    "event": {
      "timestamp": "2025-01-01T12:00:00Z",
      "type": "transaction",
      "body": {
        "amount": 1000000,
        "currency": "USD"
      }
    },
    "why": "Transaction amount is 500x higher than average"
  }
}
```

## Security

### Signature Verification

If you configured a secret, Driftlock includes an `X-Driftlock-Signature` header in every request. This allows you to verify that the request genuinely came from Driftlock.

The signature is an HMAC-SHA256 hash of the request body using your secret.

#### Node.js Example

```javascript
const crypto = require('crypto');

function verifySignature(req, secret) {
  const signature = req.headers['x-driftlock-signature'];
  const body = JSON.stringify(req.body);
  
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(body)
    .digest('hex');
    
  return signature === expectedSignature;
}
```

#### Python Example

```python
import hmac
import hashlib
import json

def verify_signature(request, secret):
    signature = request.headers.get('X-Driftlock-Signature')
    body = request.body
    
    expected_signature = hmac.new(
        secret.encode(), 
        body, 
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(signature, expected_signature)
```

## Retries

If your server returns a non-2xx response or times out, Driftlock will attempt to resend the webhook up to 5 times with exponential backoff.

## Testing

You can trigger a test webhook from the dashboard to verify your integration.
