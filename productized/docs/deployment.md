# DriftLock Deployment Documentation

This document provides detailed instructions for deploying DriftLock to production using Kubernetes and Cloudflare.

## Architecture Overview

```
Internet
   |
Cloudflare (WAF, Rate Limiting, CDN)
   |
   ├── Frontend (Cloudflare Pages)
   ├── Edge Workers (Caching, Additional Security)
   └── Backend API (Kubernetes)
         |
         ├── API Service (Load Balanced)
         ├── Database (PostgreSQL)
         ├── Cache (Redis)
         └── Message Queue (Kafka)
```

## Prerequisites

### Infrastructure
- Kubernetes cluster (EKS, GKE, AKS, or self-hosted)
- Container registry (Docker Hub, AWS ECR, GCR, etc.)
- PostgreSQL database (managed or self-hosted)
- Redis instance (for caching and sessions)
- Kafka cluster (for event streaming)

### Accounts and Services
- Cloudflare account with DNS management
- Stripe account for billing
- SendGrid or SMTP service for emails
- Google Analytics 4 property (optional)
- Domain names registered and pointed to Cloudflare

### Tools
- `kubectl` configured for your cluster
- `docker` for building images
- `helm` (optional, for Helm charts)
- `cf-terraforming` or Terraform (for infrastructure as code)

## Environment Setup

### Environment Variables
Create a `.env.production` file with the following variables:

```env
# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=driftlock
DB_USER=your-db-user
DB_PASSWORD=your-db-password

# Authentication
JWT_SECRET=your-jwt-secret
SESSION_SECRET=your-session-secret

# Stripe
STRIPE_SECRET_KEY=sk_live_...
STRIPE_PUBLISHABLE_KEY=pk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Email
SENDGRID_API_KEY=your-sendgrid-api-key
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
EMAIL_FROM=noreply@driftlock.com

# Analytics
GA4_MEASUREMENT_ID=G-...
GA4_API_KEY=your-ga4-api-key

# Cloudflare
CLOUDFLARE_API_KEY=your-cloudflare-api-key
CLOUDFLARE_ACCOUNT_ID=your-account-id

# Application
SERVER_ADDRESS=:8080
DEBUG=false
```

## Backend Deployment

### 1. Build and Push Docker Images

```bash
# Build the API image
docker build -t your-registry/driftlock-api:latest -f Dockerfile .

# Tag for production
docker tag your-registry/driftlock-api:latest your-registry/driftlock-api:prod-latest

# Push to registry
docker push your-registry/driftlock-api:latest
docker push your-registry/driftlock-api:prod-latest
```

### 2. Set Up Kubernetes Secrets

```bash
# Create namespace
kubectl create namespace driftlock

# Create database secret
kubectl create secret generic driftlock-db-secret \
  --from-literal=url="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" \
  --from-literal=user="$DB_USER" \
  --from-literal=password="$DB_PASSWORD" \
  -n driftlock

# Create auth secret
kubectl create secret generic driftlock-auth-secret \
  --from-literal=jwt="$JWT_SECRET" \
  -n driftlock

# Create billing secret
kubectl create secret generic driftlock-billing-secret \
  --from-literal=stripe="$STRIPE_SECRET_KEY" \
  -n driftlock

# Create email secret
kubectl create secret generic driftlock-email-secret \
  --from-literal=sendgrid="$SENDGRID_API_KEY" \
  -n driftlock

# Create analytics secret
kubectl create secret generic driftlock-analytics-secret \
  --from-literal=ga4="$GA4_API_KEY" \
  -n driftlock
```

### 3. Deploy to Kubernetes

```bash
# Apply database deployment
kubectl apply -f deployment/k8s/db-deployment.yaml

# Apply API deployment
kubectl apply -f deployment/k8s/api-deployment.yaml

# Verify deployment
kubectl get pods -n driftlock
kubectl get services -n driftlock
```

### 4. Set Up Ingress with SSL

If using a managed load balancer with SSL termination:

```bash
# If using cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml

# Apply ingress with TLS
kubectl apply -f deployment/k8s/api-ingress.yaml
```

## Frontend Deployment

### 1. Update Environment Variables

Update the `cloudflare-pages.toml` file with production values:

```toml
[env.production]
  NODE_VERSION = "18"
  REACT_APP_API_URL = "https://api.driftlock.com"
  REACT_APP_STRIPE_PUBLIC_KEY = "pk_live_..."
```

### 2. Deploy to Cloudflare Pages

```bash
# Using Cloudflare Pages CLI (wrangler)
npm install -g wrangler
wrangler pages deploy ./build --project-name=driftlock-frontend

# Or through Cloudflare Dashboard
# Connect your GitHub repository to Cloudflare Pages
```

### 3. Configure Custom Domain

In the Cloudflare dashboard:
1. Go to Pages → Your Project → Settings → Custom Domains
2. Add your domain (e.g., app.driftlock.com)
3. Update DNS settings to point to Cloudflare

## Cloudflare Configuration

### 1. DNS Settings
Ensure your DNS records point to:
- Frontend: Cloudflare Pages (proxy enabled)
- API: Your Kubernetes load balancer IP (proxy enabled)

### 2. Page Rules
Set up page rules for performance:
- Cache static assets: `app.driftlock.com/*static/*`
- Skip cache for API calls: `api.driftlock.com/api/*`

### 3. Security Settings
- SSL/TLS: Full (encrypts to origin)
- Security Level: Medium or High
- Bots: Enable Bot Fight Mode
- Rate Limiting: Configure as needed

### 4. Edge Workers
Deploy rate limiting worker:

```javascript
// In Cloudflare dashboard or using Wrangler
// Apply rate limiting to API endpoints
```

## Monitoring and Logging

### 1. Set Up Monitoring
- Configure health checks for Kubernetes services
- Set up alerting for system metrics
- Monitor API response times and error rates

### 2. Logging
- Centralized logging (ELK stack, Cloudflare Logs, etc.)
- Audit log monitoring
- Error tracking

## Deployment Pipeline

### CI/CD Configuration

#### GitHub Actions Example:
```yaml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Build and Push Docker Image
      run: |
        docker build -t your-registry/driftlock-api:${{ github.sha }} .
        docker push your-registry/driftlock-api:${{ github.sha }}
    
    - name: Deploy to Kubernetes
      run: |
        kubectl set image deployment/driftlock-api driftlock-api=your-registry/driftlock-api:${{ github.sha }} -n driftlock
```

## Rollback Procedure

If issues occur after deployment:

1. Roll back Kubernetes deployment:
```bash
kubectl rollout undo deployment/driftlock-api -n driftlock
```

2. If needed, rollback database changes (using migrations)
3. Monitor system after rollback

## Post-Deployment Tasks

### 1. Verification
- Test API endpoints
- Verify frontend functionality
- Check billing system
- Validate audit logging

### 2. Performance Testing
- Load test the API
- Verify rate limiting
- Check database performance

### 3. Security Validation
- Penetration testing
- Security scan of deployed images
- Verify SSL certificates

## Maintenance

### 1. Regular Tasks
- Database backups and verification
- Dependency updates
- Security patching
- Performance optimization

### 2. Scaling
- Monitor resource utilization
- Adjust Kubernetes resource limits
- Scale database as needed
- Load test before scaling

## Troubleshooting

Common issues and solutions:

1. **API not responding**:
   - Check Kubernetes pods: `kubectl get pods -n driftlock`
   - Check service endpoints: `kubectl get svc -n driftlock`
   - Check logs: `kubectl logs -l app=driftlock-api -n driftlock`

2. **Database connection issues**:
   - Verify database is running
   - Check database secrets in Kubernetes
   - Confirm network connectivity

3. **SSL certificate issues**:
   - Verify DNS settings
   - Check Cloudflare proxy status
   - Confirm certificate status in Cloudflare dashboard

4. **Frontend not loading**:
   - Check Cloudflare Pages build status
   - Verify custom domain configuration
   - Confirm API URL in environment variables