# Driftlock Setup Guide - Cloud SQL + Firebase Auth

🎉 **Production-ready setup using Google Cloud SQL and Firebase Authentication**

This guide walks you through setting up Driftlock with the recommended Google Cloud ecosystem: Cloud SQL for database and Firebase Auth for authentication.

## 🚀 Quick Start

The complete setup can be completed in 3 main steps:

```bash
# 1. Set up Cloud SQL + Firebase Auth infrastructure
./scripts/setup-gcp-cloudsql-firebase.sh

# 2. Deploy to production
./scripts/deploy-production-cloudsql-firebase.sh

# 3. Test the complete setup
./scripts/test-deployment-complete.sh
```

## 📋 Prerequisites

### Required Accounts
- **Google Cloud Platform** account with project `driftlock`
- **Firebase** project (created with the GCP project)
- **Stripe** account (optional, for payments)
- **SendGrid** account (optional, for emails)

### Required Tools
```bash
# Install required CLI tools
brew install gcloud-cli firebase-cli stripe-cli

# Verify installation
gcloud --version
firebase --version
```

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Firebase      │    │   GCP Cloud     │    │   GCP Cloud     │
│   Hosting       │◄──►│   Run (API)     │◄──►│   SQL           │
│                 │    │                 │    │                 │
│ Landing Page    │    │ Anomaly         │    │ PostgreSQL      │
│ & SPA           │    │ Detection API   │    │ Database        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   GCP Secret    │              │
         └──────────────►│   Manager       │◄─────────────┘
                        │                 │
                        │ • DB URLs       │
                        │ • Firebase Keys │
                        │ • API Keys      │
                        └─────────────────┘
                                   │
                        ┌─────────────────┐
                        │   Firebase      │
                        │   Auth          │
                        │                 │
                        │ • User Auth     │
                        │ • JWT Tokens    │
                        │ • Email Verify  │
                        └─────────────────┘
```

## 🔧 Step 1: Cloud SQL + Firebase Setup

### 1.1 Run the Setup Script

```bash
./scripts/setup-gcp-cloudsql-firebase.sh
```

This script will:
- ✅ Enable all required GCP APIs
- ✅ Create Cloud SQL PostgreSQL instance
- ✅ Set up database and user
- ✅ Configure Firebase Authentication
- ✅ Generate Firebase service account key
- ✅ Create all required GCP secrets
- ✅ Run database migrations

### 1.2 Manual Firebase Auth Configuration

After the script completes, configure Firebase Auth:

1. **Go to Firebase Console**: https://console.firebase.google.com/project/driftlock/authentication
2. **Enable Sign-in Methods**:
   - Email/Password
   - Google (optional)
   - Other providers as needed
3. **Configure Email Templates**:
   - Email verification
   - Password reset
   - Email change verification
4. **Add Authorized Domains**:
   - `https://driftlock.web.app`
   - `https://driftlock-staging.web.app`
   - `localhost:5173` (for development)

### 1.3 Verify Cloud SQL Setup

```bash
# Check instance status
gcloud sql instances describe driftlock-db --project=driftlock

# Connect to database (optional)
gcloud sql connect driftlock-db --user=driftlock_user --project=driftlock
```

## 🚀 Step 2: Deployment

### 2.1 Deploy to Production

```bash
./scripts/deploy-production-cloudsql-firebase.sh
```

This will:
- ✅ Build and deploy backend to Cloud Run
- ✅ Configure Cloud SQL connectivity
- ✅ Deploy frontend to Firebase Hosting
- ✅ Set up IAM permissions
- ✅ Configure monitoring
- ✅ Test the deployment

### 2.2 Deployment URLs

After successful deployment:
- **Frontend**: `https://driftlock.web.app`
- **Backend**: `https://driftlock-api-xxxxx-uc.a.run.app`
- **API Health**: `[backend-url]/healthz`

## 🔐 Authentication Integration

### Firebase Auth Flow

1. **User Registration**:
   ```javascript
   import { getAuth, createUserWithEmailAndPassword } from "firebase/auth";

   const auth = getAuth();
   await createUserWithEmailAndPassword(auth, email, password);
   ```

2. **Token Verification (Backend)**:
   ```go
   // Using Firebase Admin SDK
   token, err := authClient.VerifyIDToken(ctx, idToken)
   ```

3. **API Authentication**:
   ```javascript
   // Include Firebase ID token in API requests
   const token = await user.getIdToken();

   fetch('/api/v1/detect', {
     headers: {
       'Authorization': `Bearer ${token}`,
       'Content-Type': 'application/json'
     }
   });
   ```

### Backend Configuration

The backend automatically:
- Loads Firebase service account from GCP Secret Manager
- Verifies JWT tokens on authenticated endpoints
- Maps Firebase users to internal user data
- Handles email verification requirements

## 🗄️ Database Schema

### Cloud SQL Configuration

- **Instance**: `driftlock-db`
- **Database**: `driftlock`
- **User**: `driftlock_user`
- **Connection**: Uses Cloud SQL Proxy via Cloud Run

### Connection String Format
```
postgresql://driftlock_user:password@127.0.0.1:5432/driftlock?host=/cloudsql/driftlock:us-central1:driftlock-db
```

### Key Tables
- `users` - Firebase user mappings
- `tenants` - Organization/tenant data
- `api_keys` - Generated API keys
- `detections` - Anomaly detection results
- `subscriptions` - Stripe subscription data

## 📊 Monitoring & Management

### View Logs
```bash
# Cloud Run logs
gcloud logs tail "resource.type=cloud_run" --project=driftlock

# Cloud SQL logs
gcloud sql instances logs list driftlock-db --project=driftlock
```

### Console Links
- **Cloud Run**: https://console.cloud.google.com/run
- **Cloud SQL**: https://console.cloud.google.com/sql
- **Firebase Auth**: https://console.firebase.google.com/project/driftlock/authentication
- **Secret Manager**: https://console.cloud.google.com/security/secret-manager

### Monitoring Metrics
- Request latency and error rates
- Database connection count
- Authentication success/failure rates
- Cloud SQL CPU and memory usage

## 🧪 Testing

### Run Complete Test Suite
```bash
./scripts/test-deployment-complete.sh
```

### Manual Testing Commands
```bash
# Test backend health
curl https://your-backend-url/healthz

# Test with Firebase token (get from frontend)
curl -H "Authorization: Bearer FIREBASE_TOKEN" \
     https://your-backend-url/api/v1/user/profile
```

### Local Development
```bash
# Start Firebase emulators
firebase emulators:start --project=driftlock

# Start backend with local config
./scripts/start-api.sh

# Start frontend
./scripts/start-frontend.sh
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Source |
|----------|-------------|--------|
| `DATABASE_URL` | Cloud SQL connection string | GCP Secret Manager |
| `FIREBASE_SERVICE_ACCOUNT_KEY` | Firebase Admin SDK key | GCP Secret Manager |
| `FIREBASE_PROJECT_ID` | Firebase project ID | Environment variable |
| `FIREBASE_AUTH_ENABLED` | Enable Firebase Auth | Environment variable |
| `DRIFTLOCK_API_KEY` | Admin API key | GCP Secret Manager |

### GCP Secrets

All sensitive data is stored in GCP Secret Manager:
- `driftlock-db-url` - Database connection string
- `firebase-service-account-key` - Firebase service account JSON
- `driftlock-api-key` - Admin API key
- `stripe-secret-key` - Stripe API key (optional)
- `sendgrid-api-key` - SendGrid API key (optional)

## 🚨 Troubleshooting

### Common Issues

**Cloud SQL Connection Failed**:
```bash
# Check IAM permissions
gcloud projects get-iam-policy driftlock

# Verify service account has Cloud SQL Client role
gcloud projects add-iam-policy-binding driftlock \
  --member="serviceAccount:driftlock-compute@developer.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

**Firebase Auth Not Working**:
```bash
# Check Firebase project configuration
firebase projects:list

# Verify service account key format
gcloud secrets versions access latest \
  --secret=firebase-service-account-key --project=driftlock | jq .
```

**Deployment Failures**:
```bash
# Check build logs
gcloud builds list --project=driftlock

# Check Cloud Run logs
gcloud logs tail "resource.type=cloud_run" --project=driftlock
```

## 🎯 Best Practices

### Security
- Use least privilege IAM roles
- Regularly rotate secrets
- Enable VPC Serverless access for Cloud SQL
- Use Cloud Armor for additional protection

### Performance
- Configure appropriate Cloud SQL instance size
- Use connection pooling
- Enable Cloud SQL automatic backups
- Monitor Cloud Run request patterns

### Scalability
- Configure Cloud Run min/max instances
- Use Cloud SQL read replicas for read-heavy workloads
- Implement proper caching strategies
- Set up alerting for metric thresholds

## 📞 Support

### Getting Help
1. Check the GCP Console for error details
2. Review the logs using provided commands
3. Verify all prerequisites are met
4. Run the test suite for diagnostics

### Documentation Files
- `docs/GCP_SECRETS_CHECKLIST.md` - Detailed secrets reference
- `cloudbuild.yaml` - Build and deployment configuration
- `landing-page/firebase.json` - Firebase hosting configuration

---

**Your Driftlock SaaS with Cloud SQL + Firebase Auth is ready for production! 🎉**