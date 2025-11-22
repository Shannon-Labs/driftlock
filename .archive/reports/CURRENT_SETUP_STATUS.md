# Driftlock Current Setup Status

**Generated**: November 19, 2025  
**Project**: driftlock  
**GCP Account**: hunter@shannonlabs.dev  
**Note**: Some checks require GCP authentication refresh (`gcloud auth login`)

## ✅ Verified: What's Already Set Up

### Code & Configuration (✅ VERIFIED)
- ✅ **Backend Code**: Firebase Auth + Stripe billing fully integrated
  - `collector-processor/cmd/driftlock-http/auth.go` - Firebase Admin SDK
  - `collector-processor/cmd/driftlock-http/billing.go` - Stripe checkout/portal
  - `collector-processor/cmd/driftlock-http/dashboard.go` - Protected dashboard endpoints
  
- ✅ **Frontend Code**: Firebase client + dashboard UI complete
  - `landing-page/src/firebase.ts` - Firebase client initialized (uses env vars)
  - `landing-page/src/views/DashboardView.vue` - Dashboard with API keys
  - `landing-page/src/views/LoginView.vue` - Magic link login
  - `landing-page/src/stores/` - Pinia store for auth state

- ✅ **Database Migrations**: 3 SQL migrations present
  - `api/migrations/20250301120000_initial_schema.sql` - Initial schema
  - `api/migrations/20250302000000_onboarding.sql` - Onboarding tables
  - `api/migrations/20251119000000_add_stripe_fields.sql` - Stripe fields

- ✅ **Frontend Build**: Production build exists
  - `landing-page/dist/index.html` exists (2.0KB, built Nov 18, 2025)

- ✅ **Firebase Configuration**: Hosting config ready
  - `landing-page/firebase.json` - Configured with Cloud Run rewrites
  - `landing-page/.firebaserc` - Project set to "driftlock"
  - `landing-page/.env.production` - Exists (67 bytes, 1 line - contains `VITE_API_BASE_URL`)

- ✅ **Cloud Build Config**: Deployment pipeline ready
  - `cloudbuild.yaml` - Configured with all secrets
  - Expected secrets (from cloudbuild.yaml):
    - `driftlock-api-key`
    - `driftlock-db-url`
    - `driftlock-license-key`
    - `firebase-service-account-key`
    - `sendgrid-api-key`
    - `stripe-price-id-pro`
    - `stripe-secret-key`

- ✅ **Documentation**: Setup guides present
  - `docs/LAUNCH_SECRETS_GUIDE.md` - Complete secrets reference
  - `SAAS_SETUP_README.md` - Automated setup guide
  - `SETUP_GUIDE.md` - Detailed instructions

### Infrastructure Configuration (✅ VERIFIED)
- ✅ **GCP Project**: Set to "driftlock"
- ✅ **Firebase Project**: Linked to "driftlock" (via `.firebaserc`)
- ✅ **Cloud Run Service**: Backend URL found in code: `https://driftlock-api-o6kjgrsowq-uc.a.run.app`
- ✅ **Region**: us-central1 (from cloudbuild.yaml)
- ✅ **Frontend Deployed**: `https://driftlock.web.app` returns HTTP 200

## ⚠️ What Needs Verification/Setup

### GCP Secrets (⚠️ REQUIRES AUTH REFRESH)
**Status**: Cannot verify - GCP auth tokens need refreshing

To check secrets, first run:
```bash
gcloud auth login
```

Then verify secrets:
```bash
gcloud secrets list --project=driftlock
```

**Expected secrets** (from `cloudbuild.yaml`):
- [ ] `driftlock-db-url` - Database connection string
- [ ] `driftlock-license-key` - License key
- [ ] `firebase-service-account-key` - Firebase Admin SDK JSON
- [ ] `sendgrid-api-key` - SendGrid API key (optional)
- [ ] `stripe-secret-key` - Stripe secret key
- [ ] `stripe-price-id-pro` - Stripe Price ID
- [ ] `driftlock-api-key` - API key (optional)

### Frontend Environment Variables (⚠️ PARTIALLY CONFIGURED)
**Status**: `.env.production` exists but only contains 1 line (67 bytes)

Current content (masked):
- ✅ `VITE_API_BASE_URL=***` (points to backend)

**Missing variables** (required by `landing-page/src/firebase.ts`):
- [ ] `VITE_FIREBASE_API_KEY`
- [ ] `VITE_FIREBASE_AUTH_DOMAIN`
- [ ] `VITE_FIREBASE_PROJECT_ID`
- [ ] `VITE_FIREBASE_STORAGE_BUCKET`
- [ ] `VITE_FIREBASE_MESSAGING_SENDER_ID`
- [ ] `VITE_FIREBASE_APP_ID`
- [ ] `VITE_STRIPE_PUBLISHABLE_KEY` (for billing UI)

**To fix**: Add Firebase config from Firebase Console → Project Settings → General → Your apps → Web app config

### Backend Deployment Status (⚠️ REQUIRES AUTH REFRESH)
**Status**: Backend URL found in code, but health check returns 404

**Found URL**: `https://driftlock-api-o6kjgrsowq-uc.a.run.app`
- Health endpoint (`/healthz`) returns 404 HTML page
- Service may be deployed but endpoint path incorrect, or service needs redeployment

To verify:
```bash
gcloud auth login  # Refresh auth first
gcloud run services describe driftlock-api \
  --region=us-central1 \
  --project=driftlock \
  --format="value(status.url)"
```

### Database Status (⚠️ REQUIRES AUTH REFRESH)
**Status**: Migrations exist, but database connection not verified

**Migrations ready**: 3 SQL files in `api/migrations/`
- [ ] Database instance created (Cloud SQL or Supabase)
- [ ] Migrations applied to database
- [ ] Connection string stored in `driftlock-db-url` secret

To check Cloud SQL:
```bash
gcloud auth login
gcloud sql instances list --project=driftlock
```

### Firebase Auth Configuration (⚠️ NEEDS MANUAL CHECK)
**Status**: Firebase project linked, but auth config needs verification

- ✅ Firebase project linked: "driftlock" (via `.firebaserc`)
- [ ] Email/Password sign-in method enabled (check Firebase Console)
- [ ] Authorized domains configured (check Firebase Console)
- [ ] Service account key downloaded and stored in `firebase-service-account-key` secret

**To check**: Visit https://console.firebase.google.com/project/driftlock/authentication

### Stripe Configuration (⚠️ NEEDS MANUAL CHECK)
**Status**: Backend code ready, but Stripe setup needs verification

- ✅ Backend billing code integrated (`billing.go`)
- [ ] Product created ("Driftlock Pro") - check Stripe Dashboard
- [ ] Price created (monthly subscription) - check Stripe Dashboard
- [ ] Price ID stored in `stripe-price-id-pro` secret
- [ ] Secret key stored in `stripe-secret-key` secret
- [ ] Webhook endpoint configured (after backend is working)
- [ ] Webhook signing secret stored in `stripe-webhook-secret` secret

**To check**: Visit https://dashboard.stripe.com/products

## 🚀 Quick Commands to Check Status

**⚠️ IMPORTANT**: Most commands require GCP authentication refresh first:
```bash
gcloud auth login
firebase login --reauth
```

### 1. Check GCP Secrets (REQUIRES AUTH)
```bash
# List all secrets
gcloud secrets list --project=driftlock

# Check if specific secret exists
gcloud secrets describe driftlock-db-url --project=driftlock

# View secret value (for verification)
gcloud secrets versions access latest --secret=driftlock-db-url --project=driftlock
```

### 2. Check Frontend Env (✅ VERIFIED)
```bash
cd landing-page
cat .env.production
# Currently only has VITE_API_BASE_URL - needs Firebase + Stripe vars
```

### 3. Check Backend Deployment (REQUIRES AUTH)
```bash
# List Cloud Run services
gcloud run services list --project=driftlock --region=us-central1

# Get service URL
gcloud run services describe driftlock-api \
  --region=us-central1 \
  --project=driftlock \
  --format="value(status.url)"

# Check service logs
gcloud logs tail "resource.type=cloud_run resource.labels.service_name=driftlock-api" \
  --project=driftlock
```

### 4. Check Firebase Hosting (✅ VERIFIED)
```bash
# Frontend is deployed and accessible
curl -I https://driftlock.web.app
# Returns: HTTP/2 200

# Check Firebase project
cd landing-page
firebase projects:list
```

### 5. Test Backend Health (⚠️ CURRENTLY FAILING)
```bash
# Backend URL from code
BACKEND_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app"
curl "$BACKEND_URL/healthz"
# Currently returns 404 - may need redeployment or correct path
```

## 📋 Next Steps Checklist (Priority Order)

### 1. Refresh Authentication (REQUIRED FIRST)
```bash
gcloud auth login
firebase login --reauth
```

### 2. Verify/Set Up GCP Secrets
```bash
# Check what exists
gcloud secrets list --project=driftlock

# If missing, run setup script
./scripts/setup-gcp-secrets.sh
```

**Required secrets to verify**:
- `driftlock-db-url` - Database connection string
- `driftlock-license-key` - License key (can use "dev-mode")
- `firebase-service-account-key` - Firebase Admin SDK JSON
- `stripe-secret-key` - Stripe secret key
- `stripe-price-id-pro` - Stripe Price ID

### 3. Complete Frontend Environment Variables
```bash
cd landing-page
# Edit .env.production to add Firebase config
```

**Add these from Firebase Console** (https://console.firebase.google.com/project/driftlock/settings/general):
```bash
VITE_FIREBASE_API_KEY=AIzaSy...
VITE_FIREBASE_AUTH_DOMAIN=driftlock.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=driftlock
VITE_FIREBASE_STORAGE_BUCKET=driftlock.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=123456789
VITE_FIREBASE_APP_ID=1:123456789:web:abc123456
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...  # From Stripe Dashboard
```

Then rebuild and redeploy:
```bash
npm run build
firebase deploy --only hosting --project=driftlock
```

### 4. Verify/Fix Backend Deployment
```bash
# Check if service exists and get correct URL
gcloud run services describe driftlock-api \
  --region=us-central1 \
  --project=driftlock

# If service doesn't exist or health check fails, redeploy
gcloud builds submit --config=cloudbuild.yaml --project=driftlock
```

### 5. Set Up Stripe (if not done)
- Create product in Stripe Dashboard
- Create price (monthly subscription)
- Store Price ID in `stripe-price-id-pro` secret
- Store secret key in `stripe-secret-key` secret

### 6. Set Up Stripe Webhook (after backend is working)
- Go to Stripe Dashboard → Webhooks
- Add endpoint: `[BACKEND_URL]/stripe/webhook`
- Copy signing secret and store in `stripe-webhook-secret` secret

### 7. Test Complete Flow
- Visit https://driftlock.web.app
- Sign up with email
- Check email for magic link
- Login → Dashboard → Billing

## 🔍 Automated Status Check

Use the provided script to check everything:

```bash
./scripts/check-setup-status.sh
```

This script will:
- Check GCP authentication
- List all secrets (if authenticated)
- Check Cloud Run service status
- Verify frontend build and env vars
- Check Firebase configuration
- Test frontend deployment
- List database migrations

## 📝 Summary Notes

**What's Working**:
- ✅ Frontend code complete and built
- ✅ Frontend deployed to Firebase Hosting (https://driftlock.web.app)
- ✅ Backend code complete with Firebase + Stripe integration
- ✅ Database migrations ready (3 SQL files)
- ✅ Firebase project linked
- ✅ Cloud Build config ready

**What Needs Action**:
- ⚠️ GCP authentication refresh required for secret verification
- ⚠️ Frontend `.env.production` missing Firebase + Stripe config
- ⚠️ Backend health check failing (may need redeployment)
- ⚠️ Database connection not verified
- ⚠️ Stripe product/price setup not verified

**Quick Win**: Start with refreshing auth, then verify secrets exist. Most likely only need to:
1. Add Firebase config to `.env.production`
2. Rebuild/redeploy frontend
3. Verify backend deployment

