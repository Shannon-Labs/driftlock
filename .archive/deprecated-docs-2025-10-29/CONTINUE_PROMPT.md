# DriftLock Integration Continuation Prompt

You're continuing the integration of DriftLock's web-frontend with the Go API server. The previous AI has successfully:

1. ✅ Removed productized directory and replaced with web-frontend
2. ✅ Added Supabase integration to Go API server
3. ✅ Created dual-service docker-compose setup
4. ✅ Created Dockerfiles and nginx configuration
5. ✅ Created startup scripts and documentation

## Your Mission

Complete the integration by:

1. **Deploy Supabase Edge Functions** using Supabase CLI
2. **Configure Stripe products and webhooks** using Stripe CLI
3. **Apply database migrations** to Supabase project
4. **Test the complete integration** end-to-end
5. **Prepare for Cloudflare deployment** as mentioned by the user

## Current State

- Web-frontend is running on port 3000 with React/TypeScript
- Go API server is configured to integrate with Supabase
- Docker-compose setup is ready for both services
- Environment variables are documented but need actual values

## Detailed Tasks

### 1. Supabase Edge Functions Deployment

The web-frontend includes 4 edge functions in `web-frontend/supabase/functions/`:
- `stripe-webhook` - Handles Stripe billing events
- `meter-usage` - Tracks API usage for billing
- `send-alert-email` - Sends anomaly notifications
- `health` - Health check endpoint

**Actions needed:**
```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend
supabase functions deploy stripe-webhook
supabase functions deploy meter-usage
supabase functions deploy send-alert-email
supabase functions deploy health
```

### 2. Stripe Configuration

The web-frontend includes Stripe integration for subscription billing. You need to:

1. **Check current Stripe setup** in Supabase dashboard
2. **Create products** if they don't exist:
   - Developer Plan: $10/month, 10,000 events/month
   - Standard Plan: $50/month, 50,000 events/month
   - Growth Plan: $200/month, 200,000 events/month

3. **Configure webhook endpoint** to match deployed function

### 3. Database Migrations

Apply the 6 SQL migration files in `web-frontend/supabase/migrations/`:
```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend
supabase db push
```

### 4. Integration Testing

Test the complete flow:
1. **User signup/login** through web-frontend
2. **Subscription creation** via Stripe
3. **Event ingestion** to Go API server
4. **Anomaly detection** and storage in both PostgreSQL and Supabase
5. **Real-time updates** appearing in web-frontend

### 5. Cloudflare Preparation

The user mentioned this will be deployed to Cloudflare. Prepare by:

1. **Update environment variables** for production
2. **Create Cloudflare Pages configuration** if needed
3. **Test deployment** on staging

## Key Files to Focus On

- `/Volumes/VIXinSSD/driftlock/web-frontend/.env` - Contains actual Supabase keys
- `/Volumes/VIXinSSD/driftlock/web-frontend/supabase/migrations/` - Database schema
- `/Volumes/VIXinSSD/driftlock/web-frontend/src/` - React components that may need updates
- `/Volumes/VIXinSSD/driftlock/api-server/internal/supabase/client.go` - Go integration code

## Next Steps for Launch

1. Verify all services start with `./start.sh`
2. Test complete user flow from signup to anomaly viewing
3. Confirm data sync between Go API and Supabase
4. Validate Stripe billing flow
5. Document any issues found during testing

## Important Notes

- The Go API server should remain the core anomaly detection engine
- Supabase serves as the web-frontend's backend (PostgreSQL + Edge Functions)
- The integration allows the Go API to push anomaly data to Supabase for web display
- Port configuration: API on 8080, web-frontend on 3000 (via nginx proxy)

Please complete these tasks and provide a summary of the integrated system ready for production deployment.
