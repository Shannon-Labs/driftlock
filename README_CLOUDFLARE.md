# DriftLock - Production-Ready Anomaly Detection Platform 🚀

## Quick Start

### Deploy to Cloudflare (Recommended for Production)

**One-command deployment:**
```bash
# Deploy API to Cloudflare Workers
cd /Volumes/VIXinSSD/driftlock/cloudflare-api-worker
bash deploy.sh

# Deploy Web Frontend to Cloudflare Pages
cd /Volumes/VIXinSSD/driftlock/web-frontend
bash deploy-pages.sh
```

**Full deployment guide:** [CLOUDFLARE_DEPLOYMENT.md](./CLOUDFLARE_DEPLOYMENT.md)

### Run Locally (Development)

```bash
# Start all services with Docker
docker-compose -f docker-compose.yml up -d

# Or start individually:
# Web Frontend (http://localhost:3000)
cd web-frontend && npm run dev

# API Server (http://localhost:8080)
cd api-server && go run ./cmd/driftlock-api
```

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     DriftLock Platform                      │
├─────────────────────────────────────────────────────────────┤
│  Web Frontend (React)  │  API Server (Go)                  │
│  - Auth via Supabase   │  - CBAD Algorithm                 │
│  - Real-time Dashboard │  - Anomaly Detection              │
│  - Billing Management  │  - Supabase Integration           │
│  - Usage Tracking      │  - Prometheus Metrics             │
└───────────────┬─────────┴───────────────────┬───────────────┘
                │                           │
                ▼                           ▼
┌─────────────────────────────────────────────────────────────┐
│                 Supabase Backend                            │
│  - PostgreSQL Database  - Edge Functions                    │
│  - Authentication       - Real-time Subscriptions           │
│  - Row-Level Security   - Stripe Integration                │
└─────────────────────────────────────────────────────────────┘
```

**Alternative Cloudflare Deployment:**
```
┌────────────────────────┐     ┌──────────────────────────┐
│  Cloudflare Pages      │     │  Cloudflare Workers      │
│  (React Frontend)      │────▶│  (TypeScript API)        │
│                        │     │                          │
│  - Fast CDN            │     │  - Edge Computing        │
│  - Auto SSL            │     │  - Low Latency           │
│  - Global Distribution │     │  - 100K req/day free     │
└─────────────┬──────────┘     └──────────────┬───────────┘
              │                               │
              └───────────────┬───────────────┘
                              ▼
                    ┌─────────────────────┐
                    │   Supabase Backend  │
                    │   (Managed DB +     │
                    │   Edge Functions)   │
                    └─────────────────────┘
```

## What's Included

### ✅ Complete Full-Stack Application

1. **Web Frontend** (`/web-frontend/`)
   - React 18 + TypeScript + Vite
   - Tailwind CSS + shadcn/ui components
   - Real-time anomaly dashboard
   - Billing management
   - Usage tracking with alerts
   - Supabase authentication

2. **API Server** (`/api-server/`)
   - Go-based anomaly detection
   - CBAD algorithm implementation
   - OpenTelemetry tracing
   - Prometheus metrics
   - Supabase integration
   - Kafka streaming support

3. **Supabase Backend**
   - PostgreSQL database (multi-tenant)
   - 4 Edge Functions deployed:
     - `health` - Health check
     - `stripe-webhook` - Stripe event handler
     - `meter-usage` - Usage tracking
     - `send-alert-email` - Email notifications
   - 6 database migrations
   - Row-Level Security enabled
   - Real-time subscriptions

4. **Billing System**
   - Stripe integration
   - 3 subscription plans
   - Pay-for-anomalies model
   - Usage alerts and soft caps
   - Invoice management
   - Promotion code support

5. **Cloudflare Deployment**
   - Workers API server (TypeScript/Hono)
   - Pages web frontend
   - Automated deployment scripts
   - Environment configuration
   - Custom domain support

### 📁 Project Structure

```
driftlock/
├── api-server/                      # Go API server
│   ├── cmd/driftlock-api/          # Entry point
│   ├── internal/
│   │   ├── supabase/              # Supabase client
│   │   ├── handlers/              # HTTP handlers
│   │   ├── cbad/                  # CBAD algorithm
│   │   └── config/                # Configuration
│   └── go.mod
│
├── web-frontend/                    # React web app
│   ├── src/
│   │   ├── components/            # React components
│   │   │   └── dashboard/         # Dashboard UI
│   │   ├── pages/                 # Page components
│   │   ├── lib/                   # Utilities
│   │   └── integrations/          # Supabase client
│   ├── supabase/
│   │   ├── functions/             # Edge functions (4)
│   │   └── migrations/            # Database schema (6)
│   └── package.json
│
├── cloudflare-api-worker/          # Cloudflare Workers API
│   ├── src/
│   │   └── index.ts               # Hono-based API
│   ├── wrangler.toml              # Workers config
│   └── deploy.sh                  # Deployment script
│
├── .env                            # Environment variables
├── docker-compose.yml              # Local development
├── test-integration-simple.sh      # Integration tests (26 checks)
├── CLOUDFLARE_DEPLOYMENT.md        # Deployment guide
└── SYSTEM_ARCHITECTURE.md          # Architecture docs
```

## Key Features

### 🔍 Anomaly Detection
- **CBAD Algorithm**: Compression-Based Anomaly Detection
- **Real-time Processing**: Stream events through API
- **Multiple Event Types**: Logs, metrics, traces, LLM I/O
- **Severity Levels**: Low, medium, high
- **Metadata Support**: Rich context for each anomaly

### 💰 Billing System
- **Transparent Pricing**: Pay only for anomaly detections
- **Multiple Tiers**:
  - Developer (Free): 1,000 detections/month
  - Standard ($49/month): 50,000 detections + overage
  - Growth ($249/month): 500,000 detections + overage
- **Usage Alerts**: 70%, 90%, 100% thresholds
- **Soft Caps**: Warn at 120%, prevent overage surprises
- **Invoice Management**: Automatic Stripe integration

### 📊 Real-time Dashboard
- **Live Anomaly Feed**: See anomalies as they're detected
- **Usage Tracking**: Monitor consumption in real-time
- **Billing Overview**: Current plan, usage, invoices
- **Sensitivity Control**: Adjust detection sensitivity (0.0-1.0)
- **Cost Calculator**: Estimate charges before they happen

### 🔐 Security
- **Row-Level Security**: Database-level tenant isolation
- **JWT Authentication**: Secure API access
- **CORS Protection**: Properly configured
- **Webhook Verification**: Stripe signature validation
- **Audit Trails**: Complete billing action logging
- **No PII in Logs**: Privacy-first design

### 🌐 Cloud Deployment
- **Cloudflare Workers**: Edge-based API
- **Cloudflare Pages**: Global CDN for frontend
- **Supabase**: Managed PostgreSQL + Auth
- **Automatic SSL**: Zero configuration
- **Global Edge**: <50ms latency worldwide
- **Auto-scaling**: Handles traffic spikes

## Documentation

### Core Documentation
- **[SYSTEM_ARCHITECTURE.md](./SYSTEM_ARCHITECTURE.md)** - Complete system architecture
- **[CLOUDFLARE_DEPLOYMENT.md](./CLOUDFLARE_DEPLOYMENT.md)** - Production deployment guide
- **[web-frontend/IMPLEMENTATION_COMPLETE.md](./web-frontend/IMPLEMENTATION_COMPLETE.md)** - Web frontend implementation details
- **[web-frontend/PRODUCTION_RUNBOOK.md](./web-frontend/PRODUCTION_RUNBOOK.md)** - Operational runbook

### Testing
- **Integration Tests**: `bash test-integration-simple.sh` (26 tests)
- **Deployment Verification**: `bash verify-deployment.sh`
- **Local Testing**: `docker-compose -f docker-compose.yml up -d`

### API Documentation

#### Go API Server (`http://localhost:8080`)
```
GET  /healthz                      # Health check
GET  /v1/version                   # Version info
POST /v1/events                    # Ingest events
GET  /v1/anomalies                 # Get anomalies
GET  /v1/anomalies/{id}            # Get single anomaly
PATCH /v1/anomalies/{id}/status    # Update anomaly status
GET  /v1/stream/anomalies          # SSE stream
GET  /readyz                       # Readiness (DB ping; Supabase best-effort)
```

#### Cloudflare Worker API
```
GET  /health                    # Health check
GET  /                          # API info
POST /anomalies                 # Create anomaly
GET  /anomalies                 # Get anomalies
POST /usage                     # Track usage
GET  /usage                     # Get usage stats
GET  /subscription              # Get subscription
POST /stripe-webhook            # Stripe webhook
```

#### Supabase Edge Functions
```
GET  /functions/v1/health              # Health check
POST /functions/v1/meter-usage         # Usage tracking
POST /functions/v1/send-alert-email    # Send alerts
POST /functions/v1/stripe-webhook      # Stripe events
```

## Quick Commands

### Development
```bash
# Start everything locally
docker-compose up -d

# Run tests
bash test-integration-simple.sh

# Check logs
docker-compose logs -f api
docker-compose logs -f web-frontend
```

### Cloudflare Deployment
```bash
# Deploy API
cd cloudflare-api-worker && bash deploy.sh

# Deploy Frontend
cd web-frontend && bash deploy-pages.sh

# Verify deployment
bash verify-deployment.sh
```

### Database Management
```bash
# View Supabase dashboard
open https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh

# Check edge functions
open https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions

# View logs
open https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/logs
```

## Environment Variables

### Required Configuration

**`.env` (API Server):**
```bash
SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh
SUPABASE_ANON_KEY=...
SUPABASE_SERVICE_ROLE_KEY=...
SUPABASE_BASE_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co
SUPABASE_WEBHOOK_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook
PORT=8080
```

**`web-frontend/.env`:**
```bash
VITE_SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh
VITE_SUPABASE_PUBLISHABLE_KEY=...
VITE_SUPABASE_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co
```

## Integration Guide

### Using the API

**Create an Anomaly:**
```bash
curl -X POST http://localhost:8080/anomalies \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "ORG_123",
    "event_type": "error",
    "severity": "high",
    "metadata": {"source": "application.log"}
  }'
```

**Track Usage (for billing):**
```bash
curl -X POST http://localhost:8080/usage \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "ORG_123",
    "count": 1,
    "is_anomaly": true
  }'
```

### Frontend Integration

**React Component:**
```typescript
import { createClient } from '@supabase/supabase-js'

const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_PUBLISHABLE_KEY
)

// Create anomaly
const { data, error } = await supabase
  .from('anomalies')
  .insert({
    organization_id: 'ORG_123',
    event_type: 'error',
    severity: 'high'
  })

// Track usage
await supabase.functions.invoke('meter-usage', {
  body: {
    organization_id: 'ORG_123',
    count: 1,
    anomaly: true
  }
})
```

## Monitoring & Observability

### Dashboards
- **Supabase**: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh
- **Stripe**: https://dashboard.stripe.com
- **Cloudflare**: https://dash.cloudflare.com

### Metrics
- **API Server**: Prometheus at `/metrics`
- **Cloudflare Workers**: Built-in analytics
- **Database**: Supabase metrics dashboard
- **Billing**: Stripe Dashboard + Supabase invoices

### Logs
- **API Server**: Structured JSON logs
- **Edge Functions**: Supabase Dashboard → Logs
- **Workers**: `wrangler tail driftlock-api`

## Production Checklist

### Pre-Launch
- [ ] Run integration tests: `bash test-integration-simple.sh`
- [ ] Deploy to Cloudflare Workers
- [ ] Deploy to Cloudflare Pages
- [ ] Configure environment variables
- [ ] Set Cloudflare secrets
- [ ] Test Stripe webhook
- [ ] Verify email delivery
- [ ] Run load tests

### Post-Launch
- [ ] Update Stripe webhook URL
- [ ] Configure custom domains
- [ ] Set up monitoring alerts
- [ ] Verify SSL certificates
- [ ] Test end-to-end flow
- [ ] Document incident response
- [ ] Set up backup strategy
- [ ] Configure log retention

## Support & Resources

### Documentation
- [Cloudflare Workers Docs](https://developers.cloudflare.com/workers/)
- [Cloudflare Pages Docs](https://developers.cloudflare.com/pages/)
- [Supabase Docs](https://supabase.com/docs)
- [Stripe Docs](https://stripe.com/docs)

### Community
- [Cloudflare Discord](https://discord.gg/cloudflaredev)
- [Supabase Discord](https://discord.supabase.com/)
- [Stripe Discord](https://discord.gg/stripe)

### Commercial Support
- Cloudflare Pro: $20/month
- Supabase Pro: $25/month
- Stripe: 2.9% + $0.30 per transaction

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please read our contributing guidelines before submitting PRs.

## Roadmap

- [ ] Advanced anomaly detection models
- [ ] Additional data sources (Prometheus, Datadog)
- [ ] Custom detection rules
- [ ] Multi-region deployment
- [ ] Team collaboration features
- [ ] API rate limiting tiers
- [ ] White-label solution

---

## Summary

**DriftLock is a production-ready, cloud-native anomaly detection platform with:**

✅ **Complete Full-Stack Application**
- React frontend with real-time dashboard
- Go API server with CBAD algorithm
- Supabase backend with Edge Functions
- Stripe billing integration

✅ **Cloudflare Deployment Ready**
- Workers API (TypeScript/Hono)
- Pages frontend
- Automated deployment scripts
- Custom domain support

✅ **Production Features**
- Real-time anomaly detection
- Transparent billing (pay-for-anomalies)
- Usage alerts and soft caps
- Security (RLS, Auth, Webhooks)
- Monitoring and observability

✅ **Developer Experience**
- Local development with Docker
- Comprehensive test suite
- Clear documentation
- Easy deployment

**Ready to launch!** 🚀

---

For questions or support, please open an issue on GitHub.
