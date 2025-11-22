# Cloud Run Deployment Guide

This guide describes how to deploy the Driftlock API backend to Google Cloud Run.

## Prerequisites

1.  **Google Cloud SDK (`gcloud`)**: Installed and authenticated (`gcloud auth login`).
2.  **Google Cloud Project**: A project created with billing enabled.
3.  **APIs Enabled**:
    *   Cloud Run API
    *   Cloud Build API
    *   Artifact Registry or Container Registry API

## Configuration

Driftlock uses `cloudbuild.yaml` for building and deploying.

### 1. Set Default Project

```bash
gcloud config set project YOUR_PROJECT_ID
```

### 2. Enable Services

```bash
gcloud services enable run.googleapis.com cloudbuild.googleapis.com containerregistry.googleapis.com
```

### 3. Deployment

Run the deployment script which handles the build and deploy process:

```bash
./deploy.sh
```

Or run manually:

```bash
gcloud builds submit --config cloudbuild.yaml .
```

## Service Configuration (`service.yaml`)

The default configuration in `service.yaml` sets up:
- **Autoscaling**: 0 to 100 instances (scales to zero when idle).
- **Resources**: Default CPU/Memory.
- **Environment Variables**:
    - `LOG_LEVEL`: info
    - `DRIFTLOCK_ENV`: production

## Environment Variables

To connect to a database or add license keys, update the environment variables in Cloud Run console or `service.yaml`:

```yaml
        env:
          - name: DATABASE_URL
            value: "postgres://user:pass@host:5432/db"
          - name: DRIFTLOCK_LICENSE_KEY
            value: "your-license-key"
```

## Verification

After deployment, get the service URL:

```bash
gcloud run services describe driftlock-api --format 'value(status.url)'
```

Test the health endpoint:

```bash
curl $(gcloud run services describe driftlock-api --format 'value(status.url)')/healthz
```

