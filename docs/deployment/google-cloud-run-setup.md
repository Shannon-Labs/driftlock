# Google Cloud Run Deployment Guide

This guide covers deploying the Driftlock HTTP API service to Google Cloud Run.

## Prerequisites

1. Google Cloud Project with billing enabled
2. `gcloud` CLI installed and authenticated
3. Cloud Build API enabled
4. Cloud Run API enabled
5. Artifact Registry or Container Registry API enabled

## Initial Setup

### 1. Set Project ID

```bash
export PROJECT_ID="your-gcp-project-id"
gcloud config set project $PROJECT_ID
```

### 2. Enable Required APIs

```bash
gcloud services enable \
  cloudbuild.googleapis.com \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com
```

### 3. Create Secrets

Store your database URL and license key in Secret Manager:

```bash
# Database URL (PostgreSQL connection string)
echo -n "postgres://user:pass@host:5432/dbname?sslmode=require" | \
  gcloud secrets create driftlock-db-url --data-file=-

# License key
echo -n "your-license-key-here" | \
  gcloud secrets create driftlock-license-key --data-file=-
```

Grant Cloud Run service account access to secrets:

```bash
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
SERVICE_ACCOUNT="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

gcloud secrets add-iam-policy-binding driftlock-db-url \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding driftlock-license-key \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/secretmanager.secretAccessor"
```

### 4. Set Up Database (Supabase Recommended)

**We recommend Supabase over Cloud SQL** - see `docs/deployment/supabase-setup.md` for full guide.

Quick setup:
1. Create project at https://supabase.com
2. Get connection string from Settings → Database
3. Use **pooler URL** (port 6543) for Cloud Run
4. Store in Secret Manager as `driftlock-db-url`

**Alternative: Cloud SQL** (if you prefer GCP-native):

```bash
# Create Cloud SQL instance
gcloud sql instances create driftlock-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --root-password=YOUR_ROOT_PASSWORD

# Create database
gcloud sql databases create driftlock --instance=driftlock-db

# Create user
gcloud sql users create driftlock \
  --instance=driftlock-db \
  --password=YOUR_USER_PASSWORD

# Get connection name for Cloud Run
CONNECTION_NAME=$(gcloud sql instances describe driftlock-db \
  --format="value(connectionName)")
```

## Deployment Methods

### Option 1: Cloud Build (Recommended for CI/CD)

This uses the `cloudbuild.yaml` configuration:

```bash
# Submit build
gcloud builds submit --config=cloudbuild.yaml

# Or trigger from GitHub (set up Cloud Build GitHub app)
```

The build will:
1. Build Rust core library
2. Build Go HTTP service
3. Build Docker image
4. Push to Container Registry
5. Deploy to Cloud Run

### Option 2: Manual Docker Build and Deploy

```bash
# Build Docker image locally
docker build \
  -f collector-processor/cmd/driftlock-http/Dockerfile \
  -t gcr.io/$PROJECT_ID/driftlock-api:latest \
  --build-arg USE_OPENZL=false \
  .

# Push to Container Registry
docker push gcr.io/$PROJECT_ID/driftlock-api:latest

# Deploy to Cloud Run
gcloud run deploy driftlock-api \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars PORT=8080,CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net \
  --set-secrets DATABASE_URL=driftlock-db-url:latest,DRIFTLOCK_LICENSE_KEY=driftlock-license-key:latest \
  --memory 2Gi \
  --cpu 2 \
  --min-instances 1 \
  --max-instances 10 \
  --timeout 300 \
  --concurrency 80
```

### Option 3: Using service.yaml (Knative)

```bash
# Update PROJECT_ID in service.yaml first
sed -i "s/PROJECT_ID/$PROJECT_ID/g" service.yaml

# Deploy
gcloud run services replace service.yaml --region=us-central1
```

## Post-Deployment

### 1. Run Migrations

Connect to your Cloud Run service and run migrations:

```bash
# Get service URL
SERVICE_URL=$(gcloud run services describe driftlock-api \
  --region us-central1 \
  --format 'value(status.url)')

# Run migrations (requires direct database access or exec into container)
# Option A: Use Cloud SQL Proxy
cloud-sql-proxy $CONNECTION_NAME &
export DATABASE_URL="postgres://user:pass@127.0.0.1:5432/driftlock?sslmode=disable"
./bin/driftlock-http migrate up

# Option B: Use Cloud Run Jobs (recommended)
# Create a Cloud Run Job for migrations
gcloud run jobs create driftlock-migrate \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=driftlock-db-url:latest \
  --command /usr/local/bin/driftlock-http \
  --args migrate,up \
  --max-retries 1

# Execute migration job
gcloud run jobs execute driftlock-migrate --region us-central1
```

### 2. Create Initial Tenant

```bash
# Use Cloud Run Jobs or exec into container
gcloud run jobs create driftlock-create-tenant \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=driftlock-db-url:latest,DRIFTLOCK_LICENSE_KEY=driftlock-license-key:latest \
  --command /usr/local/bin/driftlock-http \
  --args create-tenant,--name,Demo,--key-role,admin,--json \
  --max-retries 1

gcloud run jobs execute driftlock-create-tenant --region us-central1
```

### 3. Verify Deployment

```bash
SERVICE_URL=$(gcloud run services describe driftlock-api \
  --region us-central1 \
  --format 'value(status.url)')

# Health check
curl $SERVICE_URL/healthz | jq

# Should return:
# {
#   "success": true,
#   "library_status": "healthy",
#   "database": "connected",
#   "license": { ... }
# }
```

## Environment Variables

| Variable | Source | Description |
|----------|--------|-------------|
| `PORT` | Env var | HTTP port (default: 8080) |
| `DATABASE_URL` | Secret | PostgreSQL connection string |
| `DRIFTLOCK_LICENSE_KEY` | Secret | License key for production |
| `CORS_ALLOW_ORIGINS` | Env var | Comma-separated allowed origins |
| `LOG_LEVEL` | Env var | Logging level (info/debug) |

## Scaling Configuration

Current settings:
- **Min instances**: 1 (always warm)
- **Max instances**: 10
- **CPU**: 2 vCPU
- **Memory**: 2Gi
- **Concurrency**: 80 requests per instance
- **Timeout**: 300 seconds

Adjust in `service.yaml` or via `gcloud run services update`.

## Monitoring

### View Logs

```bash
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=driftlock-api" \
  --limit 50 \
  --format json
```

### Metrics

View in Cloud Console:
- Cloud Run → driftlock-api → Metrics
- Key metrics: Request count, Latency, Error rate, Instance count

## Troubleshooting

### Service won't start

1. Check logs: `gcloud run services logs read driftlock-api --region us-central1`
2. Verify secrets are accessible
3. Check database connectivity
4. Verify Docker image builds successfully

### Database connection issues

1. Ensure Cloud SQL instance is running
2. Check firewall rules allow Cloud Run to connect
3. Verify connection name and credentials in secret
4. Use Cloud SQL Proxy for testing

### High latency

1. Increase CPU allocation
2. Increase min instances for always-warm
3. Check database performance
4. Review application logs for bottlenecks

## Cost Optimization

- Use `db-f1-micro` for development (free tier eligible)
- Set `min-instances: 0` for dev/staging (cold starts OK)
- Use `min-instances: 1` for production (always warm)
- Monitor Cloud SQL costs (consider smaller instance tiers)

## Next Steps

- Set up Cloud Build triggers for automatic deployments
- Configure custom domain mapping
- Set up Cloud Armor for DDoS protection
- Configure Cloud CDN if needed
- Set up alerting policies in Cloud Monitoring

