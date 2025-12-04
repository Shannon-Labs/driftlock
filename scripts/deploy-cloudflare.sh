#!/bin/bash
# Deploy Driftlock to Cloudflare (Workers + Pages)

set -e

echo "🚀 Deploying Driftlock to Cloudflare..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if wrangler is installed
if ! command -v wrangler &> /dev/null; then
    echo -e "${YELLOW}⚠️  Wrangler CLI not found. Installing...${NC}"
    npm install -g wrangler
fi

# Step 1: Deploy Workers
echo -e "\n${BLUE}📦 Step 1: Deploying Cloudflare Workers...${NC}"
cd workers

if [ ! -d "node_modules" ]; then
    echo "Installing Workers dependencies..."
    npm install
fi

echo "Deploying Workers..."
npm run deploy

cd ..

# Step 2: Build and deploy Pages
echo -e "\n${BLUE}🌐 Step 2: Building and deploying Cloudflare Pages...${NC}"
cd landing-page

echo "Installing Pages dependencies..."
bun install --frozen-lockfile

echo "Building frontend..."
bun run build

echo "Deploying to Cloudflare Pages..."
bun run deploy:cloudflare

cd ..

echo -e "\n${GREEN}✅ Deployment complete!${NC}"
echo -e "\n${YELLOW}📝 Next steps:${NC}"
echo "1. Configure Workers routes in Cloudflare Dashboard"
echo "2. Set environment variables in Pages settings"
echo "3. Update DNS if needed"
echo "4. Test endpoints: https://api.driftlock.net/healthz"
echo ""
echo "See docs/deployment/CLOUDFLARE_MIGRATION.md for details"



