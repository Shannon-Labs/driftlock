# Authentication Guide

Learn how to authenticate with Driftlock's APIs using Firebase Auth and API keys.

## Overview

Driftlock uses two authentication methods:

1. **Firebase Authentication** - For dashboard and user-level operations
2. **API Keys** - For programmatic access to detection APIs

## Firebase Authentication

Firebase Auth is used for:
- Logging into the web dashboard
- Managing your account settings
- Creating and managing API keys
- Viewing usage and billing

### Signing Up

1. Visit [https://driftlock.web.app](https://driftlock.web.app)
2. Click **"Sign Up"**
3. Choose your sign-in method:
   - **Email/Password**: Standard email registration
   - **Google**: Sign in with your Google account
   - **GitHub** (coming soon): Sign in with GitHub

4. Verify your email (if using email/password)
5. You'll be logged in and redirected to the dashboard

### Sign-In Methods

#### Email/Password
```javascript
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';

const auth = getAuth();
await signInWithEmailAndPassword(auth, email, password);
```

#### Google Sign-In
```javascript
import { getAuth, signInWithPopup, GoogleAuthProvider } from 'firebase/auth';

const auth = getAuth();
const provider = new GoogleAuthProvider();
await signInWithPopup(auth, provider);
```

### Managing Your Account

After logging in, go to [Dashboard → Account Settings](https://driftlock.web.app/dashboard/settings) to:
- Update your email  
- Change password
- Enable 2FA (coming soon)
- Delete your account

## API Keys

API keys are used for programmatic access to Driftlock's detection APIs.

### Creating an API Key

1. Log in to your [Dashboard](https://driftlock.web.app/dashboard)
2. Navigate to **"API Keys"** section
3. Click **"Create API Key"**
4. Fill in the details:
   - **Name**: Descriptive name (e.g., "Production API", "Development")
   - **Role**: 
     - `admin` - Full access to all operations
     - `stream` - Limited to detection and data ingestion
5. Click **"Create"**
6. **Copy your API key immediately** - it won't be shown again!

### API Key Roles

#### Admin Role
Full access to:
- All detection and anomaly endpoints
- Export operations
- Configuration management
- User management
- Billing operations

####  Stream Role
Limited to:
- `/v1/detect` - Run detections
- `/v1/streams/{id}/events` - Ingest events
- `/v1/anomalies` - View anomalies
- `/v1/anomalies/{id}` - View anomaly details

**Best Practice**: Use stream-role keys in production apps to limit blast radius if compromised.

### Using API Keys

Include your API key in the `X-Api-Key` header:

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{"stream_id": "default", "events": [...]}'
```

### Managing API Keys

In the dashboard:
- **View all keys**: See all API keys for your account
- **Revoke keys**: Delete compromised or unused keys
- **Rotate keys**: Create new key, update apps, revoke old key

> ⚠️ **Security Tip**: Rotate API keys every 90 days

### Key Binding (Optional)

You can bind a stream-role key to a specific stream for extra security:

1. When creating key, select **"Bind to stream"**
2. Choose the stream ID
3. Key will only work with that specific stream

## Security Best Practices

### Storing API Keys

**DO:**
- ✅ Store in environment variables
- ✅ Use secret management (AWS Secrets Manager, GCP Secret Manager)
- ✅ Encrypt at rest
- ✅ Rotate regularly

**DON'T:**
- ❌ Hardcode in source code
- ❌ Commit to version control
- ❌ Share in Slack/email
- ❌ Use the same key everywhere

### Example: Environment Variables

```bash
# .env file (DO NOT COMMIT)
DRIFTLOCK_API_KEY=sk_live_abc123...

# Load in your app
const apiKey = process.env.DRIFTLOCK_API_KEY;
```

### Using Multiple Keys

Use different API keys for different environments:

```bash
# Development
DRIFTLOCK_API_KEY=sk_test_dev_xyz...

# Staging
DRIFTLOCK_API_KEY=sk_test_staging_abc...

# Production
DRIFTLOCK_API_KEY=sk_live_prod_123...
```

## Rate Limiting

API keys are subject to rate limits based on your plan:

| Plan | Rate Limit |
|------|------------|
| Developer (Free) | 60 req/min |
| Starter ($25/mo) | 300 req/min |
| Pro (Custom) | 1,000+ req/min |

### Rate Limit Headers

Every API response includes rate limit information:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1672531200
```

### Handling Rate Limits

When you exceed the limit, you'll receive a `429 Too Many Requests` response:

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "request_id": "req_123",
    "retry_after_seconds": 30
  }
}
```

**Implementation**:

```python
import time
import requests

def detect_with_retry(api_key, payload, max_retries=3):
    for attempt in range(max_retries):
        response = requests.post(
            "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect",
            headers={"X-Api-Key": api_key},
            json=payload
        )
        
        if response.status_code == 429:
            retry_after = int(response.headers.get('Retry-After', 60))
            print(f"Rate limited. Retrying after {retry_after}s...")
            time.sleep(retry_after)
            continue
            
        return response.json()
    
    raise Exception("Max retries exceeded")
```

## Firebase Admin SDK (Advanced)

If you're building a backend that needs to verify Firebase tokens:

```javascript
const admin = require('firebase-admin');

// Initialize with service account
admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

// Verify ID token
async function verifyUser(idToken) {
  const decodedToken = await admin.auth().verifyIdToken(idToken);
  return decodedToken.uid;
}
```

## Troubleshooting

### "unauthorized" Error

**Possible causes**:
- API key is incorrect
- API key has been revoked
- Missing `X-Api-Key` header
- Using wrong header name (should be `X-Api-Key`, not `Authorization`)

**Solution**:
```bash
# Check your API key format
echo $DRIFTLOCK_API_KEY

# Verify header is included
curl -v -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect
```

### "forbidden" Error

**Possible causes**:
- Using stream-role key for admin operation
- Key is bound to different stream
- Account suspended or over quota

**Solution**:
- Check key role in dashboard
- Use admin-role key for admin operations
- Verify account status

### Firebase Auth Not Working

**Possible causes**:
- Email not verified
- Account disabled
- Wrong password
- Browser blocking cookies

**Solution**:
1. Check email for verification link
2. Try password reset
3. Clear browser cache/cookies
4. Try incognito mode

## Next Steps

- **[API Reference](../api/rest-api.md)**: Learn about all available endpoints
- **[Code Examples](../api/examples/python-examples.md)**: See authentication in action
- **[Billing Guide](../guides/billing.md)**: Manage subscriptions and payments

---

Need help? Contact support at support@driftlock.io
