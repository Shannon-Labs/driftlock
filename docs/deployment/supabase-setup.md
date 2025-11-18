# Supabase Setup for Google Cloud Run

**Recommendation: Use Supabase instead of Cloud SQL**

## Why Supabase?

✅ **Faster setup** - Managed Postgres ready in minutes  
✅ **Connection pooling** - Built-in pgBouncer (better for serverless)  
✅ **Free tier** - Generous for development/testing  
✅ **Better DX** - Dashboard, SQL editor, real-time subscriptions  
✅ **PostgreSQL compatible** - Works with your existing code (no changes needed)  
✅ **Can migrate later** - Easy to move to Cloud SQL if needed  

## Setup Steps

### 1. Create Supabase Project

1. Go to https://supabase.com
2. Create new project
3. Choose region closest to your Cloud Run region (e.g., `us-central1`)
4. Wait for database to provision (~2 minutes)

### 2. Get Connection Strings

In Supabase Dashboard → Settings → Database:

**Direct connection** (for migrations):
```
postgresql://postgres:[YOUR-PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
```

**Connection pooling** (for Cloud Run - recommended):
```
postgresql://postgres.[PROJECT-REF]:[YOUR-PASSWORD]@aws-0-us-central-1.pooler.supabase.com:6543/postgres
```

**Note**: Use the **pooler** URL for Cloud Run (port 6543) - it handles connection pooling better for serverless.

### 3. Store in Google Secret Manager

```bash
export PROJECT_ID="your-gcp-project-id"
export SUPABASE_DB_URL="postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-us-central-1.pooler.supabase.com:6543/postgres"

# Create secret
echo -n "$SUPABASE_DB_URL" | \
  gcloud secrets create driftlock-db-url --data-file=-

# Grant Cloud Run access
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
SERVICE_ACCOUNT="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

gcloud secrets add-iam-policy-binding driftlock-db-url \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/secretmanager.secretAccessor"
```

### 4. Run Migrations

Your migrations are in `api/migrations/20250301120000_initial_schema.sql`.

**Option A: Using Supabase Dashboard**
1. Go to SQL Editor in Supabase Dashboard
2. Copy contents of `api/migrations/20250301120000_initial_schema.sql`
3. Paste and run

**Option B: Using psql locally**
```bash
export DATABASE_URL="postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres"
cd /path/to/driftlock
./bin/driftlock-http migrate up
```

**Option C: Using Cloud Run Job** (after first deployment)
```bash
gcloud run jobs create driftlock-migrate \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=driftlock-db-url:latest \
  --command /usr/local/bin/driftlock-http \
  --args migrate,up \
  --max-retries 1

gcloud run jobs execute driftlock-migrate --region us-central1
```

### 5. Update Cloud Run Service

The `service.yaml` already references `driftlock-db-url` secret - just make sure it's set to your Supabase connection string.

Deploy:
```bash
gcloud run services replace service.yaml --region=us-central1
```

## Connection Pooling Notes

Supabase provides two connection methods:

1. **Direct** (port 5432): For migrations, admin tasks
2. **Pooler** (port 6543): For application connections (recommended for Cloud Run)

**Always use the pooler URL** (`pooler.supabase.com:6543`) for Cloud Run because:
- Serverless functions can't maintain persistent connections
- Pooler handles connection lifecycle automatically
- Better performance under load
- Prevents "too many connections" errors

## Environment Variables

Your `DATABASE_URL` should look like:
```
postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-us-central-1.pooler.supabase.com:6543/postgres?sslmode=require
```

**Important**: Add `?sslmode=require` for production (Supabase requires SSL).

## Monitoring

Supabase Dashboard provides:
- Database metrics (connections, queries, size)
- Query performance insights
- Real-time logs
- Backup management

## Cost Comparison

**Supabase Free Tier:**
- 500MB database
- 2GB bandwidth
- Unlimited API requests
- Perfect for development/testing

**Supabase Pro ($25/month):**
- 8GB database
- 50GB bandwidth
- Daily backups
- Better for production

**Cloud SQL:**
- db-f1-micro: ~$7/month (shared CPU, 0.6GB RAM)
- db-n1-standard-1: ~$50/month (1 vCPU, 3.75GB RAM)

**Recommendation**: Start with Supabase Pro ($25/month) - better value than Cloud SQL for most use cases.

## Migration Path

If you need to migrate to Cloud SQL later:

1. Export from Supabase: `pg_dump` via Supabase dashboard
2. Import to Cloud SQL: `psql` to Cloud SQL instance
3. Update connection string in Secret Manager
4. Redeploy Cloud Run service

Your application code doesn't change - it's just PostgreSQL!

## Troubleshooting

**Connection timeout:**
- Use pooler URL (port 6543), not direct (port 5432)
- Check Supabase project is active
- Verify SSL mode is set (`sslmode=require`)

**Too many connections:**
- Use pooler URL (handles this automatically)
- Check connection pooling settings in Supabase dashboard

**Migration errors:**
- Run migrations using direct connection (port 5432)
- Check Supabase logs in dashboard
- Verify `pgcrypto` extension is enabled (your migration creates it)

