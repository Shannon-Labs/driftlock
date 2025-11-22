# Driftlock SaaS Setup Guide

This guide walks you through setting up Driftlock as a complete SaaS product with GCP, Firebase, Supabase, and Stripe integration.

## 🚀 Quick Start

The entire setup can be completed in 4 main steps:

```bash
# 1. Set up GCP secrets
./scripts/setup-gcp-secrets.sh

# 2. Set up database
./scripts/setup-supabase.sh

# 3. Set up payment processing
./scripts/setup-stripe.sh

# 4. Deploy to production
./scripts/deploy-production.sh
```

## 📋 Prerequisites

### Required Accounts
- **Google Cloud Platform** account with project `driftlock`
- **Supabase** account (free tier is fine)
- **Stripe** account (test mode for development)
- **SendGrid** account (for email delivery)

### Required Tools
```bash
# Install required CLI tools
brew install gcloud-cli firebase-cli stripe-cli supabase/tap/supabase

# Or manually install:
# - Google Cloud SDK: https://cloud.google.com/sdk/docs/install
# - Firebase CLI: npm install -g firebase-tools
# - Stripe CLI: https://stripe.com/docs/stripe-cli
# - Supabase CLI: https://supabase.com/docs/guides/cli
```

## 🔐 Step 1: GCP Authentication & Secrets

### 1.1 Authenticate with Google Cloud
```bash
gcloud auth login
gcloud config set project driftlock
```

### 1.2 Create Required Secrets
Run the automated script:
```bash
./scripts/setup-gcp-secrets.sh
```

This will create these secrets in GCP Secret Manager:
- `driftlock-db-url` - Supabase connection string
- `driftlock-license-key` - License key (use "dev-mode" for testing)
- `sendgrid-api-key` - SendGrid API key for emails
- `stripe-secret-key` - Stripe secret key
- `stripe-price-id-pro` - Stripe price ID for Pro plan
- `admin-key` - Admin dashboard access key

## 🗄️ Step 2: Database Setup (Supabase)

### Option A: Local Development
```bash
./scripts/setup-supabase.sh
# Choose option 1 for local development
```

### Option B: Production Supabase
```bash
./scripts/setup-supabase.sh
# Choose option 2 for cloud setup
```

### 2.1 Manual Supabase Setup
1. Go to [supabase.com/dashboard](https://supabase.com/dashboard)
2. Create new project named `driftlock`
3. Wait for project creation (2-3 minutes)
4. Get connection string from Settings → Database → Connection string
5. Use the "Transaction pooler" connection string format:
   ```
   postgresql://postgres.[ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require
   ```

### 2.2 Run Database Migrations
```bash
# For local Supabase
supabase db push

# For cloud Supabase, copy-paste migrations from:
api/migrations/20250301120000_initial_schema.sql
api/migrations/20251126212452_ee73e3fb-7c81-4e1e-a919-6b5fa722aee1.sql
```

## 💳 Step 3: Stripe Setup

### 3.1 Automated Setup
```bash
./scripts/setup-stripe.sh
# Choose option 1 for automated setup
```

### 3.2 Manual Stripe Setup
1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Create a product:
   - Name: "Driftlock Pro"
   - Description: "Professional anomaly detection"
   - Price: $99/month (or your preferred pricing)

3. Create API keys at [stripe.com/apikeys](https://dashboard.stripe.com/apikeys)
   - Copy your secret key (starts with `sk_test_` for testing)

4. Set up webhooks at [stripe.com/webhooks](https://dashboard.stripe.com/webhooks)
   - Endpoint URL: `[your-api-url]/stripe/webhook`
   - Events:
     - `customer.subscription.created`
     - `customer.subscription.updated`
     - `customer.subscription.deleted`
     - `invoice.payment_succeeded`
     - `invoice.payment_failed`

## 📧 Step 4: SendGrid Setup

1. Sign up at [sendgrid.com](https://sendgrid.com)
2. Verify your sending domain
3. Create an API key at [sendgrid.com/settings/api_keys](https://app.sendgrid.com/settings/api_keys)
4. The API key should look like: `SG.xxxxxxxx...`

## 🚀 Step 5: Deployment

### 5.1 Deploy to Production
```bash
./scripts/deploy-production.sh
```

This will:
- Build and deploy frontend to Firebase Hosting
- Build and deploy backend to Google Cloud Run
- Configure monitoring and logging
- Test the deployment

### 5.2 Manual Deployment Commands

**Frontend (Firebase):**
```bash
cd landing-page
npm install
npm run build
firebase deploy --project=driftlock
```

**Backend (Cloud Run):**
```bash
gcloud builds submit --config=cloudbuild.yaml
```

## 🌐 Accessing Your Deployment

After successful deployment:

- **Frontend**: https://driftlock.web.app
- **Backend API**: https://driftlock-api-[hash]-uc.a.run.app
- **API Health Check**: [backend-url]/healthz
- **API Documentation**: [backend-url]/api/v1

## 🧪 Testing the Complete Flow

1. **User Signup**:
   ```bash
   curl -X POST https://your-api-url/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Email Verification**:
   - Check email inbox for verification link
   - Click the link to verify email

3. **API Usage**:
   ```bash
   curl -X POST https://your-api-url/api/v1/detect \
     -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"data":[{"event":"payment","amount":100.00,"timestamp":"2024-01-01T00:00:00Z"}]}'
   ```

## 📊 Monitoring and Maintenance

### View Logs
```bash
# Cloud Run logs
gcloud logs tail "resource.type=cloud_run" --project=driftlock

# Specific service logs
gcloud logs tail "resource.type=cloud_run resource.labels.service_name=driftlock-api" --project=driftlock
```

### Monitor Performance
- **Cloud Run Console**: https://console.cloud.google.com/run
- **Firebase Analytics**: https://console.firebase.google.com/project/driftlock/analytics
- **Error Reporting**: https://console.cloud.google.com/errors

## 🔧 Troubleshooting

### Common Issues

**Authentication Issues**:
```bash
# Re-authenticate GCP
gcloud auth login

# Re-authenticate Firebase
firebase login
```

**Database Connection Errors**:
- Verify Supabase is running and accessible
- Check connection string format in GCP secrets
- Ensure SSL mode is enabled

**Stripe Webhook Issues**:
- Verify webhook URL is accessible
- Check webhook signing secret
- Ensure webhook events are selected

**Build Failures**:
- Check Cloud Build logs in GCP Console
- Verify all required secrets exist
- Check Dockerfile and cloudbuild.yaml syntax

## 🎯 Next Steps

1. **Custom Domain**: Set up custom domain in Firebase Hosting
2. **Monitoring**: Configure alerts for errors and performance
3. **Scaling**: Adjust Cloud Run resources based on usage
4. **Security**: Set up additional authentication and rate limiting
5. **Email Templates**: Customize transactional emails
6. **Analytics**: Implement user analytics and tracking

## 📞 Support

For setup issues:
1. Check the logs in GCP Console
2. Review the [troubleshooting section](#-troubleshooting)
3. Verify all secrets are correctly configured
4. Test with the provided curl commands

The Driftlock SaaS is now ready for production usage! 🎉