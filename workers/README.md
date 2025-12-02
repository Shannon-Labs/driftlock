# Driftlock Cloudflare Workers

Cloudflare Workers implementation replacing Firebase Functions for the Driftlock SaaS platform.

## Quick Start

```bash
# Install dependencies
npm install

# Set up environment variables
wrangler secret put GEMINI_API_KEY
wrangler secret put VITE_FIREBASE_API_KEY
# ... (see deployment guide for full list)

# Run locally
npm run dev

# Deploy
npm run deploy
```

## Endpoints

- `GET /healthz` - Health check
- `POST /api/v1/onboard/signup` - User signup
- `POST /api/analyze` - Anomaly analysis with Gemini
- `POST /api/compliance` - Compliance report generation
- `GET /getFirebaseConfig` - Firebase config for frontend
- `* /api/v1/**` - Proxy to Cloud Run backend
- `POST /webhooks/stripe` - Stripe webhook proxy

## Environment Variables

Set via `wrangler secret put <NAME>`:

- `GEMINI_API_KEY` - Google Gemini API key
- `VITE_FIREBASE_API_KEY` - Firebase API key
- `VITE_FIREBASE_AUTH_DOMAIN` - Firebase auth domain
- `VITE_FIREBASE_PROJECT_ID` - Firebase project ID
- `VITE_FIREBASE_STORAGE_BUCKET` - Firebase storage bucket
- `VITE_FIREBASE_MESSAGING_SENDER_ID` - Firebase messaging sender ID
- `VITE_FIREBASE_APP_ID` - Firebase app ID
- `VITE_FIREBASE_MEASUREMENT_ID` - Firebase measurement ID (optional)
- `CLOUD_RUN_API_URL` - Cloud Run backend URL (optional, defaults to production)

## Development

```bash
# Start local dev server
npm run dev

# Test endpoints
curl http://localhost:8787/healthz
curl http://localhost:8787/getFirebaseConfig
```

## Deployment

```bash
# Deploy to production
npm run deploy

# Deploy to preview environment
npm run deploy:preview
```

## Architecture

Workers act as a proxy layer between the frontend and Cloud Run backend:

1. **API Proxy**: Routes `/api/v1/**` to Cloud Run backend
2. **AI Features**: Handles Gemini AI integration for anomaly analysis
3. **Config**: Serves Firebase config to frontend
4. **CORS**: Handles CORS headers for all requests

## See Also

- [Cloudflare Migration Guide](../../docs/deployment/CLOUDFLARE_MIGRATION.md)
- [Cloudflare Workers Docs](https://developers.cloudflare.com/workers/)


