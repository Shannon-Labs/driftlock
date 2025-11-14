#!/bin/bash

# Cloudflare R2 Setup Script for Driftlock
# This script sets up the R2 bucket with proper configuration

set -e

# Configuration
BUCKET_NAME="driftlock-file-uploads"
ACCOUNT_ID="${CLOUDFLARE_ACCOUNT_ID}"

if [ -z "$ACCOUNT_ID" ]; then
    echo "Error: CLOUDFLARE_ACCOUNT_ID environment variable is required"
    echo "Please set it with: export CLOUDFLARE_ACCOUNT_ID=your_account_id"
    exit 1
fi

echo "Setting up Cloudflare R2 bucket: $BUCKET_NAME"

# Create the R2 bucket
echo "Creating R2 bucket..."
wrangler r2 bucket create "$BUCKET_NAME"

# Set up CORS configuration
echo "Configuring CORS..."
cat > /tmp/cors-config.json << 'EOF'
[
  {
    "allowedHeaders": ["*"],
    "allowedMethods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    "allowedOrigins": [
      "https://driftlock.net",
      "https://www.driftlock.net",
      "https://api.driftlock.net"
    ],
    "exposeHeaders": ["ETag", "Content-Length"],
    "maxAgeSeconds": 3600
  }
]
EOF

wrangler r2 bucket cors put "$BUCKET_NAME" --cors-file /tmp/cors-config.json

# Set up lifecycle policy (optional - for cleanup of old files)
echo "Setting up lifecycle policy..."
cat > /tmp/lifecycle-policy.json << 'EOF'
{
  "rules": [
    {
      "id": "cleanup-old-files",
      "status": "Enabled",
      "filter": {
        "prefix": "uploads/"
      },
      "expiration": {
        "days": 30
      }
    }
  ]
}
EOF

wrangler r2 bucket lifecycle put "$BUCKET_NAME" --lifecycle-file /tmp/lifecycle-policy.json

# Create API tokens documentation
echo "Creating R2 API token documentation..."
cat > /tmp/r2-tokens.md << 'EOF'
# Cloudflare R2 API Token Setup

For the Driftlock application to access R2 storage, you need to create API tokens with the following permissions:

## Required Token Permissions

1. **R2 Upload Token** (for file uploads):
   - Account permissions: `Cloudflare R2:Edit`
   - Zone permissions: None required
   - Account Resources: `*:*` (or specify your account ID)

2. **R2 Read Token** (for file downloads):
   - Account permissions: `Cloudflare R2:Read`
   - Zone permissions: None required
   - Account Resources: `*:*` (or specify your account ID)

## Environment Variables for the Application

```
CLOUDFLARE_ACCOUNT_ID=your_account_id
CLOUDFLARE_R2_ACCESS_KEY=your_r2_access_key
CLOUDFLARE_R2_SECRET_KEY=your_r2_secret_key
CLOUDFLARE_R2_BUCKET_NAME=driftlock-file-uploads
CLOUDFLARE_R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
```

## Creating the Tokens

1. Go to Cloudflare Dashboard → My Profile → API Tokens
2. Click "Create Token"
3. Use the "Custom token" template
4. Configure with the permissions specified above
5. Store the tokens securely and add them to your environment variables
EOF

echo "Setup complete!"
echo ""
echo "Important: Create the required API tokens as documented in /tmp/r2-tokens.md"
echo "Add the environment variables to your deployment configuration"
echo ""
echo "Bucket details:"
echo "  Name: $BUCKET_NAME"
echo "  Account ID: $ACCOUNT_ID"
echo "  Endpoint: https://$ACCOUNT_ID.r2.cloudflarestorage.com"

# Cleanup
rm -f /tmp/cors-config.json /tmp/lifecycle-policy.json /tmp/r2-tokens.md