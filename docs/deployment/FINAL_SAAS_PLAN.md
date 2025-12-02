# Driftlock SaaS: Path to Production

This document outlines the comprehensive plan to take Driftlock from a repository state to a fully operational SaaS business.

## 1. Prerequisite Checklist (Human Intervention Required)

Before running automated scripts, ensure the following accounts and CLI tools are ready:

- **Google Cloud Project**: Active project with billing enabled.
- **Firebase Project**: Linked to the Google Cloud Project.
- **Stripe Account**: For billing/subscriptions (obtain Secret Key and Price IDs).
- **Domain Name**: Purchased (e.g., via Google Domains) for production hosting.
- **CLI Tools**:
  - `gcloud` (Authenticated: `gcloud auth login`)
  - `firebase` (Authenticated: `firebase login`)
  - `go` (v1.21+)
  - `rust` (stable)
  - `node` (v18+)

## 2. Secrets Configuration (Security Critical)

You must manually populate Google Secret Manager with the following keys. The scripts will fail without them.

Run the following commands (replacing values with your actual secrets):

```bash
# Database Connection (Format: postgresql://user:pass@127.0.0.1:5432/db?host=/cloudsql/PROJECT:REGION:INSTANCE)
echo -n "postgresql://driftlock:PASSWORD@127.0.0.1:5432/driftlock?host=/cloudsql/YOUR_PROJECT:us-central1:driftlock-db" | gcloud secrets create driftlock-db-url --data-file=-

# License Key (Self-generated for your own instance)
echo -n "license_key_..." | gcloud secrets create driftlock-license-key --data-file=-

# Firebase Service Account (JSON content)
gcloud secrets create firebase-service-account-key --data-file=path/to/service-account.json

# Driftlock Admin API Key (You generate this, e.g., `openssl rand -hex 32`)
echo -n "your_admin_api_key" | gcloud secrets create driftlock-api-key --data-file=-

# Stripe Keys (Optional for MVP, required for Billing)
echo -n "sk_live_..." | gcloud secrets create stripe-secret-key --data-file=-
echo -n "price_..." | gcloud secrets create stripe-price-id-pro --data-file=-

# SendGrid/Email (Optional for MVP)
echo -n "SG. ..." | gcloud secrets create sendgrid-api-key --data-file=-

# Gemini API Key (For AI Analysis)
echo -n "your_gemini_key" | gcloud secrets create gemini-api-key --data-file=-

# Frontend Firebase Key (Publicly safe, but stored in secrets for injection)
echo -n "your_firebase_web_api_key" | gcloud secrets create VITE_FIREBASE_API_KEY --data-file=-
```

## 3. Infrastructure Provisioning

Use the existing automation to provision Cloud SQL and Firebase resources.

```bash
# Sets up Cloud SQL, creates database, enables APIs
./scripts/deployment/setup-gcp-cloudsql-firebase.sh
```

**Manual Verification:**
- Check Google Cloud Console: ensure `driftlock-db` instance is "Runnable".
- Check Firebase Console: ensure Authentication is enabled.

## 4. Application Deployment

Deploy the backend (Cloud Run) and frontend (Firebase Hosting) together.

```bash
# Deploys Cloud Run service 'driftlock-api' and Firebase Hosting site
./scripts/deployment/deploy-production-cloudsql-firebase.sh
```

This script performs the following:
1.  Builds the Go backend container.
2.  Submits it to Cloud Run (with Cloud SQL connection).
3.  Builds the Vue frontend.
4.  Deploys frontend + functions to Firebase.

## 5. Post-Deployment Configuration

### A. Firebase Authentication
1.  Go to **Firebase Console > Authentication > Sign-in method**.
2.  Enable **Email/Password**.
3.  Enable **Google** (optional).
4.  In **Settings > Authorized domains**, add your custom domain.

### B. DNS & Custom Domain
1.  Go to **Firebase Console > Hosting**.
2.  Click **Add Custom Domain**.
3.  Follow instructions to update DNS A records at your registrar.
4.  Wait for SSL certificate provisioning (auto-managed by Firebase).

### C. Stripe Webhooks
1.  Go to **Stripe Dashboard > Developers > Webhooks**.
2.  Add Endpoint: `https://YOUR_PROJECT.web.app/webhooks/stripe` (or your custom domain).
3.  Select events: `checkout.session.completed`, `customer.subscription.updated`.

## 6. Verification

Run the end-to-end test suite against the live environment:

```bash
./scripts/test-deployment.sh
```

## 7. Maintenance & Updates

- **Frontend Updates:** `cd landing-page && npm run build && firebase deploy --only hosting`
- **Backend Updates:** `gcloud builds submit --config=cloudbuild.yaml`
- **Database Migrations:** Run automatically on container startup, or manually via `scripts/migrate.sh`.

## 8. Business Launch
1.  **Pricing:** Update `landing-page/src/config/pricing.ts` to match Stripe Price IDs.
2.  **Support:** Update email addresses in `landing-page/src/views/ContactView.vue`.
3.  **Analytics:** Verify Google Analytics (GA4) ID in `firebase.json` or frontend config.

