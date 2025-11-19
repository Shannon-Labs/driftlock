# Driftlock Launch Guide: Secrets & Env Variables

Before deploying the full SaaS stack (Auth, Billing, Dashboard), you must configure the following secrets in Google Cloud Secret Manager.

## 1. Required GCP Secrets

Run these commands to create the secrets referenced in `cloudbuild.yaml`:

```bash
export PROJECT_ID="your-project-id"

# 1. Firebase Service Account (for Backend Auth)
# Download this JSON from Firebase Console -> Project Settings -> Service Accounts
gcloud secrets create firebase-service-account-key --data-file=./path/to/service-account.json

# 2. Stripe Secrets (for Billing)
gcloud secrets create stripe-secret-key --data-file=- <<< "sk_live_..."
gcloud secrets create stripe-webhook-secret --data-file=- <<< "whsec_..."
gcloud secrets create stripe-price-id-pro --data-file=- <<< "price_..."

# 3. Database URL (Cloud SQL / Supabase)
gcloud secrets create driftlock-db-url --data-file=- <<< "postgres://user:pass@host:5432/driftlock?sslmode=disable"

# 4. License Key
gcloud secrets create driftlock-license-key --data-file=- <<< "dev-mode" # or real key
```

## 2. Frontend Environment Variables

The frontend needs public keys to initialize Firebase and Stripe.
Create `landing-page/.env.production`:

```ini
# Firebase Public Config (from Firebase Console)
VITE_FIREBASE_API_KEY=AIzaSy...
VITE_FIREBASE_AUTH_DOMAIN=your-app.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-app
VITE_FIREBASE_STORAGE_BUCKET=your-app.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=123456789
VITE_FIREBASE_APP_ID=1:123456789:web:abc123456

# Stripe Public Key
VITE_STRIPE_PUBLISHABLE_KEY=pk_live_...
```

## 3. Deployment Steps

1. **Deploy Backend (Cloud Run):**
   ```bash
   gcloud builds submit --config=cloudbuild.yaml
   ```

2. **Deploy Frontend (Cloudflare Pages / Firebase Hosting):**
   ```bash
   cd landing-page
   npm run build
   # Deploy 'dist' folder to your hosting provider
   ```

3. **Verify:**
   - Go to `https://your-domain.com/login`
   - Log in via Magic Link.
   - Access `/dashboard`.
   - Click "Manage Billing" to test Stripe Portal.


