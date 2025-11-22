#!/bin/bash
# Create a test API key using the CLI and store it in .env
# This uses the production database to create a tenant programmatically

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$REPO_ROOT/.env"

echo "🔑 Creating Test API Key..."
echo ""

# Check if binary exists
BINARY="$REPO_ROOT/bin/driftlock-http"
if [ ! -f "$BINARY" ]; then
    BINARY="$REPO_ROOT/collector-processor/driftlock-http"
    if [ ! -f "$BINARY" ]; then
        echo "❌ Error: driftlock-http binary not found"
        echo "   Please build it first:"
        echo "   cd collector-processor/cmd/driftlock-http && go build -o ../../driftlock-http ."
        exit 1
    fi
fi

# Get DATABASE_URL from Secret Manager or env
if [ -z "$DATABASE_URL" ]; then
    echo "📥 Fetching DATABASE_URL from Secret Manager..."
    DB_URL_RAW=$(gcloud secrets versions access latest --secret="driftlock-db-url" 2>/dev/null || echo "")
    
    if [ -z "$DB_URL_RAW" ]; then
        echo "❌ Error: DATABASE_URL not found in Secret Manager"
        echo ""
        echo "For local execution, you need Cloud SQL Proxy or use Cloud Run Jobs instead:"
        echo "  ./scripts/create-test-api-key-cloudrun.sh"
        exit 1
    fi
    
    # Check if it's a Cloud SQL Unix socket (needs proxy)
    if echo "$DB_URL_RAW" | grep -q "/cloudsql/"; then
        echo "⚠️  Database URL uses Cloud SQL Unix socket"
        echo "   For local execution, use Cloud Run Jobs instead:"
        echo "   ./scripts/create-test-api-key-cloudrun.sh"
        echo ""
        echo "   Or set up Cloud SQL Proxy and use a direct connection string"
        exit 1
    fi
    
    DATABASE_URL="$DB_URL_RAW"
    echo "✅ Got DATABASE_URL from Secret Manager"
else
    echo "✅ Using DATABASE_URL from environment"
fi

export DATABASE_URL

# Create tenant with unique name
TIMESTAMP=$(date +%s)
TENANT_NAME="Crypto Test Runner $TIMESTAMP"
TENANT_SLUG="crypto-test-${TIMESTAMP}"

echo "📝 Creating tenant: $TENANT_NAME"
echo ""

# Create tenant and get API key
TENANT_JSON=$(
    "$BINARY" create-tenant \
        --name "$TENANT_NAME" \
        --slug "$TENANT_SLUG" \
        --plan "trial" \
        --key-role "admin" \
        --key-name "crypto-test-key" \
        --json 2>&1
)

if [ $? -ne 0 ]; then
    echo "❌ Failed to create tenant:"
    echo "$TENANT_JSON"
    exit 1
fi

# Extract API key
API_KEY=$(echo "$TENANT_JSON" | grep -o '"api_key":"[^"]*"' | cut -d'"' -f4 || echo "")

if [ -z "$API_KEY" ]; then
    echo "❌ Failed to extract API key from response:"
    echo "$TENANT_JSON"
    exit 1
fi

echo "✅ Created tenant successfully!"
echo "   API Key: ${API_KEY:0:30}..."
echo ""

# Save to .env file
if [ -f "$ENV_FILE" ]; then
    # Update existing DRIFTLOCK_API_KEY if present
    if grep -q "^DRIFTLOCK_API_KEY=" "$ENV_FILE"; then
        sed -i.bak "s|^DRIFTLOCK_API_KEY=.*|DRIFTLOCK_API_KEY=$API_KEY|" "$ENV_FILE"
        echo "✅ Updated DRIFTLOCK_API_KEY in .env"
    else
        echo "DRIFTLOCK_API_KEY=$API_KEY" >> "$ENV_FILE"
        echo "✅ Added DRIFTLOCK_API_KEY to .env"
    fi
else
    # Create new .env file
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
echo "You can now run:"
echo "  source .env"
echo "  ./scripts/start_crypto_test.sh"
echo ""

