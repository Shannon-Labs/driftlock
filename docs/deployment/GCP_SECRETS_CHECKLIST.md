# GCP Secrets Checklist for Driftlock

Ensure the following secrets are created in Google Secret Manager before deployment.

## Required Secrets

| Secret Name | Description | Example Value |
|---|---|---|
| `driftlock-db-url` | PostgreSQL Connection String (Cloud SQL) | `postgresql://username:password@127.0.0.1:5432/driftlock?host=/cloudsql/driftlock:us-central1:driftlock-db` |
| `driftlock-license-key` | License key for Driftlock (or "dev-mode") | `dev-mode` |
| `firebase-service-account-key` | Firebase Admin SDK Service Account Key | JSON service account key |
| `sendgrid-api-key` | SendGrid API Key for email delivery | `SG.xxxxxxxx...` |
| `stripe-secret-key` | Stripe Secret Key (Test or Live) | `sk_test_xxxxxxxx...` |
| `stripe-price-id-pro` | Stripe Price ID for the "Pro" plan | `price_1xxxxxxxx...` |
| `driftlock-api-key` | Admin API Key for service management | `sk-drltlck-xxxxxxxx...` |

## Optional / Future Secrets

| Secret Name | Description | Example Value |
|---|---|---|
| `admin-key` | Static key for admin dashboard access | `long-random-string` |
| `firebase-webhook-secret` | Firebase Auth webhook secret for token verification | `firebase-webhook-secret` |

## Creation Commands

```bash
export PROJECT_ID="driftlock"

# Database (Cloud SQL)
echo -n "postgresql://username:password@127.0.0.1:5432/driftlock?host=/cloudsql/driftlock:us-central1:driftlock-db" | gcloud secrets create driftlock-db-url --data-file=-

# License
echo -n "dev-mode" | gcloud secrets create driftlock-license-key --data-file=-

# Firebase Service Account
echo -n "$(cat firebase-service-account-key.json)" | gcloud secrets create firebase-service-account-key --data-file=-

# SendGrid
echo -n "SG.xxx" | gcloud secrets create sendgrid-api-key --data-file=-

# Stripe
echo -n "sk_test_xxx" | gcloud secrets create stripe-secret-key --data-file=-
echo -n "price_1xxx" | gcloud secrets create stripe-price-id-pro --data-file=-

# Admin API Key
echo -n "sk-drltlck-$(openssl rand -hex 16)" | gcloud secrets create driftlock-api-key --data-file=-
```

## Cloud SQL Setup

Before creating secrets, set up Cloud SQL:

```bash
# Enable Cloud SQL Admin API
gcloud services enable sqladmin.googleapis.com --project=$PROJECT_ID

# Create Cloud SQL instance
gcloud sql instances create driftlock-db \
    --project=$PROJECT_ID \
    --database-version=POSTGRES_15 \
    --tier=db-custom-4-16384 \
    --region=us-central1 \
    --storage-size=100GB \
    --storage-type=SSD \
    --backup-start-time=02:00

# Create database
gcloud sql databases create driftlock --instance=driftlock-db --project=$PROJECT_ID

# Create database user
gcloud sql users create driftlock_user --instance=driftlock-db --password=your-secure-password --project=$PROJECT_ID
```

## Firebase Auth Setup

```bash
# Enable Firebase Auth
firebase projects:enable auth driftlock

# Configure authentication providers in Firebase Console
# https://console.firebase.google.com/project/driftlock/authentication
```

