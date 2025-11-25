#!/bin/bash
# Create a test API key using Cloud Run Jobs
# This creates a tenant in the production database

set -e

PROJECT_ID="${GCP_PROJECT:-driftlock}"
REGION="us-central1"
JOB_NAME="driftlock-create-test-key-$(date +%s)"

echo "🔑 Creating Test API Key via Cloud Run Job..."
echo "   Project: $PROJECT_ID"
echo "   Region: $REGION"
echo ""

# Check if job already exists (cleanup old ones)
EXISTING_JOBS=$(gcloud run jobs list --region="$REGION" --filter="name:driftlock-create-test-key" --format="value(name)" 2>/dev/null || echo "")
if [ -n "$EXISTING_JOBS" ]; then
    echo "🧹 Cleaning up old test key jobs..."
    echo "$EXISTING_JOBS" | while read job; do
        gcloud run jobs delete "$job" --region="$REGION" --quiet 2>/dev/null || true
    done
fi

# Get the latest image
IMAGE="gcr.io/$PROJECT_ID/driftlock-api:latest"
echo "📦 Using image: $IMAGE"
echo ""

# Create the job
echo "📝 Creating Cloud Run Job..."
TIMESTAMP=$(date +%s)
TENANT_NAME="Crypto Test Runner ${TIMESTAMP}"
TENANT_SLUG="crypto-test-${TIMESTAMP}"
KEY_NAME="crypto-test-key-${TIMESTAMP}"
gcloud run jobs create "$JOB_NAME" \
    --image="$IMAGE" \
    --region="$REGION" \
    --set-cloudsql-instances="$PROJECT_ID:us-central1:driftlock-db" \
    --set-secrets="DATABASE_URL=driftlock-db-url:latest" \
    --set-env-vars="DRIFTLOCK_DEV_MODE=true" \
    --command="/usr/local/bin/driftlock-http" \
    --args="create-tenant,--name,${TENANT_NAME},--slug,${TENANT_SLUG},--plan,trial,--key-role,admin,--key-name,${KEY_NAME},--json" \
    --max-retries=1 \
    --task-timeout=60 \
    --memory=512Mi \
    --cpu=1 \
    > /dev/null

echo "✅ Job created: $JOB_NAME"
echo ""

# Execute the job
echo "🚀 Executing job to create tenant..."
EXECUTION_NAME=$(gcloud run jobs execute "$JOB_NAME" --region="$REGION" --format="value(name)" 2>/dev/null || echo "")

if [ -z "$EXECUTION_NAME" ]; then
    echo "❌ Failed to execute job"
    exit 1
fi

echo "⏳ Waiting for job to complete..."
sleep 5

# Get logs and extract API key
echo "📥 Fetching API key from job logs..."
MAX_ATTEMPTS=12
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    # Get all log entries and extract API key pattern
    LOGS=$(gcloud logging read "resource.type=cloud_run_job AND resource.labels.job_name=$JOB_NAME AND resource.labels.location=$REGION" \
        --limit=50 \
        --format="value(textPayload)" \
        --freshness=5m 2>/dev/null || echo "")
    
    # Look for API key pattern: dlk_ followed by UUID and token
    API_KEY=$(echo "$LOGS" | grep -oE 'dlk_[a-f0-9-]+\.[a-z0-9]+' | head -1 || echo "")
    
    if [ -n "$API_KEY" ]; then
        break
    fi
    
    ATTEMPT=$((ATTEMPT + 1))
    if [ $ATTEMPT -lt $MAX_ATTEMPTS ]; then
        echo "   Waiting for job output... ($ATTEMPT/$MAX_ATTEMPTS)"
        sleep 5
    fi
done

# Cleanup job
echo ""
echo "🧹 Cleaning up job..."
gcloud run jobs delete "$JOB_NAME" --region="$REGION" --quiet 2>/dev/null || true

if [ -z "$API_KEY" ]; then
    echo "❌ Failed to extract API key from job logs"
    echo ""
    echo "You can check logs manually:"
    echo "  gcloud logging read \"resource.type=cloud_run_job AND resource.labels.job_name=$JOB_NAME\" --limit=50"
    exit 1
fi

echo "✅ API key created successfully!"
echo "   Key: ${API_KEY:0:30}..."
echo ""

# Save to .env
ENV_FILE=".env"
if [ -f "$ENV_FILE" ]; then
    if grep -q "^DRIFTLOCK_API_KEY=" "$ENV_FILE"; then
        sed -i.bak "s|^DRIFTLOCK_API_KEY=.*|DRIFTLOCK_API_KEY=$API_KEY|" "$ENV_FILE"
        echo "✅ Updated DRIFTLOCK_API_KEY in .env"
    else
        echo "DRIFTLOCK_API_KEY=$API_KEY" >> "$ENV_FILE"
        echo "✅ Added DRIFTLOCK_API_KEY to .env"
    fi
else
    cat > "$ENV_FILE" << EOF
# Driftlock Test API Key
# Generated: $(date)
DRIFTLOCK_API_KEY=$API_KEY
DRIFTLOCK_API_URL=https://driftlock.web.app/api/v1
EOF
    echo "✅ Created .env file with API key"
fi

echo ""
echo "🎉 Test API key created and saved to .env"
echo ""
