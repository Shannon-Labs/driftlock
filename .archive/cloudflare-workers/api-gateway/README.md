# DriftLock API Gateway

Cloudflare Workers-based API gateway for the DriftLock anomaly detection platform with integrated Stripe webhook handling.

## Overview

This API gateway serves as the entry point for all DriftLock API requests, providing:
- API key authentication
- Usage metering and billing
- Request routing to the Go backend
- Stripe webhook processing for subscription management
- Audit logging
- Rate limiting (coming soon)
- Security filtering

## Architecture

```
Client → Cloudflare Workers (API Gateway) → Go Backend
              ↓              ↑
         Supabase ← ← ← ← ← Stripe Webhooks
         (Auth, Billing, Audit)
```

## Setup

### Prerequisites

- Node.js 18+
- Wrangler CLI installed (`npm install -g wrangler`)
- Cloudflare account
- Supabase project (with auth, billing, and metering tables)
- Stripe account with webhook endpoint configured

### Installation

1. Install dependencies:
```bash
npm install
```

2. Set up environment variables in `wrangler.toml`:
- Update `ENV_GO_BACKEND_URL` to point to your deployed Go backend
- Uncomment and set `account_id` if deploying to a custom route

3. Set secrets:
```bash
# Set your Supabase service role key
wrangler secret put ENV_SUPABASE_SERVICE_ROLE_KEY

# Set your JWT secret
wrangler secret put ENV_JWT_SECRET

# Set your Stripe webhook signing secret
wrangler secret put ENV_STRIPE_WEBHOOK_SECRET

# Set your Supabase URL (if not in wrangler.toml)
wrangler secret put ENV_SUPABASE_URL
```

### Development

To run in development mode:
```bash
wrangler dev
```

This will start a local development server that proxies requests to your local Go backend (default: http://localhost:8081).

### Deployment

To deploy to Cloudflare:
```bash
wrangler deploy
```

## API Key Authentication

All requests to the API must include an API key in the Authorization header:

```
Authorization: Bearer YOUR_API_KEY
```

## Stripe Webhook Handling

The gateway automatically processes Stripe webhooks at the root path (`/`). Configure your Stripe webhook to point to:

```
https://your-subdomain.your-username.workers.dev/
```

### Supported Stripe Events:
- `checkout.session.completed` - New subscription completed
- `customer.subscription.created` - New subscription
- `customer.subscription.updated` - Subscription plan changed
- `customer.subscription.deleted` - Subscription canceled
- `invoice.payment_succeeded` - Payment received
- `invoice.payment_failed` - Payment failed
- Other events as configured in your Stripe dashboard

## Billing Model

The gateway implements a pooled usage model:
- Only anomaly detection requests are metered
- Stream/ingestion requests are free
- Usage is pooled across both APIs
- Overage charges apply after included quota is exceeded

## Endpoints Mapped

The gateway proxies these endpoints to the Go backend:
- `GET /api/v1/anomalies` - List detected anomalies
- `POST /api/v1/anomalies/detect` - Trigger anomaly detection
- `GET /api/v1/anomalies/:id` - Get specific anomaly
- `PUT /api/v1/anomalies/:id/resolve` - Resolve anomaly
- `POST /api/v1/analyze` - Analyze data for anomalies
- `GET /api/v1/health` - Health check

**Note**: The root path (`/`) is reserved for Stripe webhooks, so API routes start with `/api/`.

## Environment Variables

Required environment variables:

- `ENV_SUPABASE_URL` - Your Supabase project URL
- `ENV_SUPABASE_SERVICE_ROLE_KEY` - Your Supabase service role key (stored as secret)
- `ENV_GO_BACKEND_URL` - Your Go backend URL
- `ENV_JWT_SECRET` - JWT secret for authentication (stored as secret)
- `ENV_STRIPE_WEBHOOK_SECRET` - Stripe webhook signing secret (stored as secret)

## Security

- All API keys are validated against Supabase
- Stripe webhook signatures are verified
- Usage is metered and billed automatically
- Audit logs are maintained for compliance
- Requests are rate-limited to prevent abuse

## Error Handling

Common error responses:
- `401 Unauthorized` - Invalid or missing API key
- `402 Payment Required` - Invalid subscription
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Gateway or backend error

## Local Development

For local development, make sure your Go backend is running:
```bash
cd /path/to/driftlock
PORT=8081 go run productized/cmd/server/main.go
```

Then start the worker in dev mode:
```bash
wrangler dev
```

The worker will be available at `http://localhost:8787` by default.