# DriftLock Production Deployment with Cloudflare

This document outlines the production deployment strategy for DriftLock using Cloudflare services.

## Architecture Overview

```
Internet
   |
Cloudflare (WAF, DDoS Protection)
   |
   ├── Frontend (Cloudflare Pages)
   ├── Edge Workers (Rate Limiting, Caching)
   └── Backend API (Self-hosted or on Cloudflare)
```

## Deployment Steps

1. **Database Setup**
   - Deploy PostgreSQL to a cloud provider (AWS RDS, GCP Cloud SQL, etc.)
   - Set up read replicas for scalability
   - Configure backup and disaster recovery

2. **Backend API Deployment**
   - Containerize the Go API using the provided Dockerfile
   - Deploy to a container service (AWS ECS, GKE, etc.)
   - Configure autoscaling based on demand

3. **Cloudflare Configuration**
   - Configure DNS to point to backend API
   - Set up WAF rules and security policies
   - Configure page rules for caching and performance
   - Set up custom domains and SSL certificates

4. **Frontend Deployment**
   - Deploy React frontend to Cloudflare Pages
   - Configure environment variables for production
   - Set up preview branches for staging

5. **Monitoring and Observability**
   - Set up logging and metrics collection
   - Configure alerts and notifications
   - Implement health checks

## Environment Configuration

### Production Environment
- Domain: `app.driftlock.com`
- API: `api.driftlock.com`
- Environment: `production`

### Staging Environment
- Domain: `staging.driftlock.com`
- API: `staging-api.driftlock.com`
- Environment: `staging`

## SSL Certificate Management
- Cloudflare automatically provisions SSL certificates
- Force HTTPS with automatic redirects
- Use TLS 1.2+ for all connections

## Security Configuration
- Enable all Cloudflare security features
- Implement rate limiting at the edge
- Use WAF rules to block common attacks
- Enable DDoS protection

## Backup and Recovery
- Database backups enabled with point-in-time recovery
- Backup retention: 30 days for production, 7 days for staging
- Test backup restoration procedures monthly

## Deployment Pipeline
- CI/CD pipeline triggers on main branch commits
- Automated tests run before deployment
- Blue-green deployment strategy to minimize downtime
- Rollback procedures defined for critical issues