# DriftLock Integration Guide

## Overview
This guide will help you prepare DriftLock for integration with your Supabase + Stripe + Cloudflare Workers stack.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Frontend (Lovable/Next.js)                            │
│  - User Interface                                      │
│  - Dashboard & Analytics                               │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│  Cloudflare Workers (api.driftlock.net)                │
│  - API Gateway & Auth                                   │
│  - Usage Metering                                       │
│  - Request Proxying                                     │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│  DriftLock API (backend.driftlock.net)                 │
│  - Event Ingestion                                      │
│  - Anomaly Detection                                    │
│  - Billing Integration                                  │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│  Data Layer                                            │
│  - Supabase (PostgreSQL)                                │
│  - Kafka (Event Streaming)                              │
│  - Redis (Caching)                                      │
└──────────────────────────────────────────────────────────┘
```

## Pre-Integration Checklist

### 1. Supabase Setup (Already Done on Lovable)
Your Supabase should have these tables:
- `api_keys` - For API key validation
- `subscriptions` - Stripe subscription status
- `billing_customers` - Stripe customer mapping
- `audit_logs` - API call tracking
- `users` - User authentication
- `organizations` - Multi-tenancy

### 2. Stripe Configuration

#### Create Products & Prices
```bash
# Free Plan
stripe products create --name="DriftLock Free" --description="10K events/mo, 7 days retention"
stripe prices create --product=<FREE_PRODUCT_ID> --unit-amount=0 --currency=usd --recurring[interval]=month

# Pro Plan
stripe products create --name="DriftLock Pro" --description="100K events/mo, 30 days retention"
stripe prices create --product=<PRO_PRODUCT_ID> --unit-amount=2900 --currency=usd --recurring[interval]=month

# Business Plan
stripe products create --name="DriftLock Business" --description="Unlimited events, 90 days retention"
stripe prices create --product=<BUSINESS_PRODUCT_ID> --unit-amount=9900 --currency=usd --recurring[interval]=month
```

#### Configure Webhooks
Set webhook endpoint to: `https://api.driftlock.net/api/v1/webhooks/stripe`

Enable these events:
- `checkout.session.completed`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.payment_succeeded`
- `invoice.payment_failed`

### 3. Cloudflare Workers Deployment

#### Step 1: Install Dependencies
```bash
cd cloudflare-workers/api-gateway
npm install
```

#### Step 2: Configure Secrets
```bash
# Set Supabase credentials
wrangler secret put ENV_SUPABASE_URL
# Enter: https://your-project.supabase.co

wrangler secret put ENV_SUPABASE_SERVICE_ROLE_KEY
# Enter: your-supabase-service-role-key

# Set backend URL
wrangler secret put ENV_GO_BACKEND_URL
# Enter: https://backend.driftlock.net

# Set JWT secret (must match backend)
wrangler secret put ENV_JWT_SECRET
# Enter: your-super-secret-jwt-key-min-32-chars

# Set Stripe webhook secret
wrangler secret put ENV_STRIPE_WEBHOOK_SECRET
# Enter: whsec_...
```

#### Step 3: Update wrangler.toml
```toml
[env.production]
name = "driftlock-api-gateway-production"

[env.production.vars]
ENV_SUPABASE_URL = "https://your-actual-project.supabase.co"
ENV_GO_BACKEND_URL = "https://backend.driftlock.net"

# Uncomment and configure routes
[[env.production.routes]]
pattern = "api.driftlock.net/*"
zone_name = "driftlock.net"
```

#### Step 4: Deploy
```bash
wrangler deploy --env production
```

### 4. DriftLock API Server Setup

#### Step 1: Configure Environment
Create `/productized/.env`:
```bash
# Server
SERVER_ADDRESS=:8080
DEBUG=false

# Database (Your Supabase connection string)
DATABASE_URL=postgres://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres

# Authentication
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
SESSION_SECRET=your-super-secret-session-key-min-32-chars

# Kafka (if using event streaming)
KAFKA_BROKERS=your-kafka-broker:9092

# CORS (add your frontend domains)
ALLOWED_ORIGINS=https://driftlock.net,https://app.driftlock.net

# Stripe
STRIPE_SECRET_KEY=sk_live_...
STRIPE_PUBLISHABLE_KEY=pk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_SUCCESS_URL=https://driftlock.net/billing/success
STRIPE_CANCEL_URL=https://driftlock.net/billing/cancel
STRIPE_CURRENCY=usd
STRIPE_FREE_PLAN_ID=price_xxx
STRIPE_PRO_PLAN_ID=price_xxx
STRIPE_BUSINESS_PLAN_ID=price_xxx
STRIPE_USAGE_ENABLED=true

# Email (Choose one)
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@driftlock.com

# OR use SendGrid
# EMAIL_PROVIDER=sendgrid
# SENDGRID_API_KEY=SG.xxx

# Analytics
GA4_MEASUREMENT_ID=G-XXXXXXXXXX
GA4_API_KEY=your-ga4-api-key
ENABLE_ANALYTICS=true

# Audit & Compliance
AUDIT_LOG_ENABLED=true
MAX_EVENT_RETENTION=30
```

#### Step 2: Build & Deploy
```bash
# Local testing
cd productized
go run cmd/server/main.go

# Production build
docker build -t driftlock-api:latest .
docker push your-registry/driftlock-api:latest
```

## Testing the Integration

### Step 1: Test API Health
```bash
curl https://backend.driftlock.net/health
# Expected: {"status":"ok","time":"..."}
```

### Step 2: Test Registration (via Cloudflare Gateway)
```bash
curl -X POST https://api.driftlock.net/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "SecurePassword123!"
  }'
```

### Step 3: Test Event Ingestion
```bash
# First, login to get JWT token
TOKEN=$(curl -s -X POST https://api.driftlock.net/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePassword123!"}' \
  | jq -r '.token')

# Ingest an event
curl -X POST https://api.driftlock.net/api/v1/events/ingest \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{
    "timestamp": "2025-01-15T12:00:00Z",
    "stream_type": "logs",
    "data": "Test log message",
    "metadata": {"source": "test"}
  }]'
```

### Step 4: Run Automated Tests
```bash
# Make sure API is running
cd /home/user/driftlock

# Run the test script
./test-api.sh
```

### Step 5: Generate Test Data
```bash
# Get JWT token first
export JWT_TOKEN="your-jwt-token-here"

# Generate test data
./generate-test-data.sh
```

## Critical Fixes Applied

### 1. Database Integration ✅
- **Fixed**: Services now query the actual database instead of returning mock data
- **Added**: Proper user registration with bcrypt password hashing
- **Added**: Automatic tenant creation on user registration
- **Fixed**: Authentication now validates against database with bcrypt

### 2. JWT Configuration ✅
- **Fixed**: JWT secret now reads from environment variables
- **Fixed**: Middleware properly injects user_id and tenant_id into context

### 3. Tenant Management ✅
- **Fixed**: AddTenantToContext now looks up user's actual tenant from database
- **Added**: Proper tenant isolation for all anomaly queries

## API Endpoints Reference

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh

### User Management
- `GET /api/v1/user` - Get profile
- `PUT /api/v1/user` - Update profile

### Anomalies
- `GET /api/v1/anomalies` - List anomalies (paginated)
- `GET /api/v1/anomalies/:id` - Get specific anomaly
- `PUT /api/v1/anomalies/:id/resolve` - Mark as resolved
- `DELETE /api/v1/anomalies/:id` - Delete anomaly

### Events
- `POST /api/v1/events/ingest` - Ingest events
- `GET /api/v1/events` - List events

### Dashboard
- `GET /api/v1/dashboard/stats` - Dashboard statistics
- `GET /api/v1/dashboard/recent` - Recent anomalies

### Billing (Stripe)
- `GET /api/v1/billing/plans` - List plans
- `POST /api/v1/billing/checkout` - Create checkout session
- `GET /api/v1/billing/portal` - Get customer portal URL
- `GET /api/v1/billing/subscription` - Get current subscription
- `DELETE /api/v1/billing/subscription` - Cancel subscription
- `GET /api/v1/billing/usage` - Get usage metrics
- `POST /api/v1/billing/usage` - Record usage

### Webhooks
- `POST /api/v1/webhooks/stripe` - Stripe webhook handler

## Monitoring & Observability

### Health Checks
- Backend: `GET /health`
- Metrics: Prometheus-compatible endpoint available

### Logging
- All API calls are logged via audit middleware
- Cloudflare Workers logs available via `wrangler tail`

### Alerts
Configure email alerts in tenant settings for:
- Critical anomalies detected
- Billing thresholds exceeded
- System errors

## Troubleshooting

### Common Issues

#### 1. Database Connection Fails
```
Error: Failed to connect to database
```
**Solution**: Check DATABASE_URL in .env. Ensure Supabase allows connections from your server IP.

#### 2. Stripe Webhooks Not Working
```
Error: Webhook signature verification failed
```
**Solution**: Ensure STRIPE_WEBHOOK_SECRET matches your Stripe webhook endpoint secret.

#### 3. CORS Errors
```
Error: CORS policy blocked
```
**Solution**: Add your frontend domain to ALLOWED_ORIGINS in .env

#### 4. JWT Token Invalid
```
Error: Invalid token
```
**Solution**: Ensure JWT_SECRET matches between API server and Cloudflare Workers

## Security Checklist

- [ ] Change all default secrets in production
- [ ] Use HTTPS/TLS for all endpoints
- [ ] Enable rate limiting on Cloudflare
- [ ] Configure Supabase RLS policies
- [ ] Set up database backups
- [ ] Enable audit logging
- [ ] Configure CORS properly
- [ ] Use strong JWT secrets (min 32 chars)
- [ ] Enable Stripe webhook signature verification
- [ ] Set up monitoring and alerts

## Next Steps

1. **Deploy to Production**
   - Set up production Supabase database
   - Configure production Stripe account
   - Deploy Cloudflare Workers
   - Deploy DriftLock API to hosting platform

2. **Configure DNS**
   - Point `api.driftlock.net` to Cloudflare Workers
   - Point `backend.driftlock.net` to API server
   - Configure SSL certificates

3. **Test End-to-End**
   - Run integration tests
   - Test billing flows
   - Verify anomaly detection
   - Test email notifications

4. **Monitor & Optimize**
   - Set up application monitoring
   - Configure error tracking
   - Monitor API performance
   - Analyze usage patterns

## Support & Documentation

- **API Documentation**: See `docs/api.md` (to be generated)
- **Architecture Overview**: See this guide
- **Deployment Guide**: See `docs/deployment.md` (to be created)
- **Issues**: https://github.com/Shannon-Labs/driftlock/issues

---

**Status**: Ready for integration testing and deployment
**Last Updated**: 2025-01-26
**Version**: 1.0.0
