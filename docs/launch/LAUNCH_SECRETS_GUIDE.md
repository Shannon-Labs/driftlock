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

## 2. Frontend Configuration

The frontend requires Firebase and Stripe configuration. The Stripe Publishable Key is public and can be stored in `landing-page/.env.production`. The Firebase configuration, including the API key, is now fetched securely at runtime.

### Storing the Firebase API Key (VITE_FIREBASE_API_KEY)

The `VITE_FIREBASE_API_KEY` must be stored in **Google Secret Manager**. This is a critical security measure to prevent exposing the key publicly.

1.  **Create the secret in Google Secret Manager:**
    ```bash
    # Replace YOUR_NEW_API_KEY with your actual key
    echo -n "YOUR_NEW_API_KEY" | gcloud secrets versions add VITE_FIREBASE_API_KEY --data-file=- --project="your-project-id"
    ```

### How it Works

1.  A Firebase Cloud Function named `getFirebaseConfig` has been created.
2.  This function securely reads the `VITE_FIREBASE_API_KEY` from Google Secret Manager.
3.  The `firebase.json` file is configured to rewrite requests from the frontend at the path `/getFirebaseConfig` to this Cloud Function.
4.  The frontend application calls this endpoint to fetch the full Firebase configuration, including the API key, when it initializes.

### Stripe Public Key

Create `landing-page/.env.production` for your Stripe key:

```ini
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





