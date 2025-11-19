# Driftlock SaaS Setup & Deployment Guide

🎉 **Your Driftlock SaaS infrastructure is now ready for deployment!**

This guide provides everything you need to deploy and operate Driftlock as a complete SaaS product with GCP, Firebase, Supabase, and Stripe integration.

## 🚀 Quick Start (30 minutes)

For immediate deployment, run these scripts in order:

```bash
# 1. Set up all your secrets and API keys
./scripts/setup-gcp-secrets.sh

# 2. Configure your database
./scripts/setup-supabase.sh

# 3. Set up payment processing
./scripts/setup-stripe.sh

# 4. Deploy to production
./scripts/deploy-production.sh

# 5. Test everything works
./scripts/test-deployment-complete.sh
```

## 📋 What's Included

✅ **Production-Ready Infrastructure**
- GCP Cloud Run for scalable API hosting
- Firebase Hosting for frontend deployment
- Supabase PostgreSQL for managed database
- GCP Secret Manager for secure configuration
- Stripe integration for payments
- SendGrid integration for emails

✅ **Multi-Environment Support**
- Production environment (`driftlock`)
- Staging environment (`driftlock-staging`)
- Local development setup

✅ **Complete Automation**
- One-command deployment
- Automated testing suite
- Database migrations
- CI/CD pipeline ready

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Firebase      │    │   GCP Cloud     │    │   Supabase      │
│   Hosting       │◄──►│   Run (API)     │◄──►│   PostgreSQL    │
│                 │    │                 │    │                 │
│ Landing Page    │    │ Anomaly         │    │ User Data       │
│ & SPA           │    │ Detection API   │    │ & API Keys      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   GCP Secret    │              │
         └──────────────►│   Manager       │◄─────────────┘
                        │                 │
                        │ • API Keys      │
                        │ • DB URLs       │
                        │ • Stripe Keys   │
                        └─────────────────┘
                                   │
                        ┌─────────────────┐
                        │   External      │
                        │   Services      │
                        │                 │
                        │ • Stripe        │
                        │ • SendGrid      │
                        └─────────────────┘
```

## 🔧 Environment Setup

### Prerequisites

**Required Accounts:**
- Google Cloud Platform (with `driftlock` project)
- Supabase account
- Stripe account
- SendGrid account

**Required Tools:**
```bash
# Install all required tools
brew install gcloud-cli firebase-cli stripe-cli supabase/tap/supabase

# Or check individual installation:
gcloud --version
firebase --version
stripe --version
supabase --version
```

### Step 1: GCP Authentication

```bash
# Authenticate with Google Cloud
gcloud auth login
gcloud config set project driftlock

# Authenticate with Firebase
firebase login
```

## 🚀 Deployment Options

### 1. Production Deployment

```bash
# Complete production setup
./scripts/deploy-production.sh
```

**What this does:**
- Builds and deploys frontend to Firebase Hosting
- Builds and deploys backend to Google Cloud Run
- Configures all environment variables and secrets
- Sets up monitoring and logging
- Tests the deployment

**URLs after deployment:**
- Frontend: `https://driftlock.web.app`
- Backend: `https://driftlock-api-xxxxx-uc.a.run.app`

### 2. Staging Environment

```bash
# Set up isolated testing environment
./scripts/setup-staging.sh
```

**Staging URLs:**
- Frontend: `https://driftlock-staging.web.app`
- Backend: `https://driftlock-api-staging-xxxxx-uc.a.run.app`

### 3. Local Development

```bash
# Complete local development setup
./scripts/setup-local-dev.sh

# Start services manually
./scripts/start-api.sh      # Start backend on :8080
./scripts/start-frontend.sh # Start frontend on :5173
```

## 📊 Monitoring & Management

### View Logs

```bash
# Production logs
gcloud logs tail "resource.type=cloud_run" --project=driftlock

# Staging logs
gcloud logs tail "resource.type=cloud_run" --project=driftlock-staging

# Specific service logs
gcloud logs tail "resource.type=cloud_run resource.labels.service_name=driftlock-api" --project=driftlock
```

### Monitor Performance

- **GCP Console**: https://console.cloud.google.com/run
- **Firebase Analytics**: https://console.firebase.google.com/project/driftlock/analytics
- **Error Reporting**: https://console.cloud.google.com/errors
- **Database**: https://supabase.com/dashboard

### Database Management

```bash
# Local database access
psql -h localhost -p 7543 -U driftlock -d driftlock

# Or use Supabase CLI
supabase db reset    # Reset local database
supabase db push     # Apply migrations
```

## 🔐 Configuration Management

### GCP Secrets

All sensitive data is stored in GCP Secret Manager:

```bash
# List all secrets
gcloud secrets list --project=driftlock

# View a secret
gcloud secrets versions access latest --secret=stripe-secret-key --project=driftlock

# Update a secret
echo -n "new-value" | gcloud secrets versions add stripe-secret-key --data-file=- --project=driftlock
```

### Environment Variables

| Environment | Database URL | Stripe Mode | Features |
|-------------|-------------|-------------|----------|
| Production | Supabase | Live/Test | Full SaaS |
| Staging | Local/Supabase | Test | Full SaaS |
| Local | PostgreSQL | Test | Full SaaS |

## 🧪 Testing

### Run Complete Test Suite

```bash
# Test production deployment
./scripts/test-deployment-complete.sh

# Choose option 1 for complete testing
```

### Manual API Testing

```bash
# Health check
curl https://your-api-url/healthz

# Test anomaly detection
curl -X POST https://your-api-url/api/v1/detect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "data": [
      {"event": "payment", "amount": 100.00, "timestamp": "2024-01-01T00:00:00Z"}
    ]
  }'
```

## 🔄 CI/CD Pipeline

The setup includes ready-to-use CI/CD configurations:

### GitHub Actions (Optional)

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy Driftlock

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to GCP
        run: ./scripts/deploy-production.sh
```

### Manual Deployment

```bash
# Deploy when needed
./scripts/deploy-production.sh

# Test deployment
./scripts/test-deployment-complete.sh
```

## 🎯 Production Checklist

Before going live, verify:

- [ ] All GCP secrets are configured with production values
- [ ] Stripe is configured with live API keys (when ready)
- [ ] SendGrid domain is verified
- [ ] Custom domain is configured (optional)
- [ ] Monitoring alerts are set up
- [ ] Database backups are enabled
- [ ] SSL certificates are active
- [ ] Rate limiting is configured
- [ ] User testing is complete

## 🔧 Troubleshooting

### Common Issues

**Build Failures:**
```bash
# Check build logs
gcloud builds list --project=driftlock

# View specific build
gcloud builds describe BUILD_ID --project=driftlock
```

**Database Connection Issues:**
```bash
# Test database connection
docker-compose exec driftlock-postgres pg_isready -U driftlock

# Check database URL secret
gcloud secrets versions access latest --secret=driftlock-db-url --project=driftlock
```

**API Issues:**
```bash
# Check API logs
gcloud logs tail "resource.type=cloud_run resource.labels.service_name=driftlock-api" --project=driftlock

# Test API locally
./scripts/start-api.sh
curl http://localhost:8080/healthz
```

## 📞 Support & Documentation

### Documentation Files
- `SETUP_GUIDE.md` - Detailed setup instructions
- `docs/GCP_SECRETS_CHECKLIST.md` - GCP secrets reference
- `docs/LAUNCH_SUMMARY.md` - Product launch information
- `README.md` - Project overview and demo information

### Getting Help
1. Check the logs in GCP Console
2. Review this setup guide
3. Run the test suite for diagnostics
4. Check individual script outputs

## 🎉 You're Ready!

With these scripts and configurations, your Driftlock SaaS is:

- ✅ Production-ready
- ✅ Fully automated
- ✅ Multi-environment
- ✅ Monitored and logged
- ✅ Secure and scalable

**Go live with confidence! 🚀**