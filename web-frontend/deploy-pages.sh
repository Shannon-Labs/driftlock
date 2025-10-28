#!/bin/bash

echo "========================================="
echo "DriftLock Web Frontend - Cloudflare Pages Deployment"
echo "========================================="
echo ""

# Check if Cloudflare CLI (cf) is installed
if ! command -v cf &> /dev/null; then
    echo "❌ Cloudflare CLI not found."
    echo "Please install it from: https://developers.cloudflare.com/pages/how-to/installing-the-cloudflare-cli/"
    exit 1
fi

echo "✅ Cloudflare CLI installed"
echo ""

# Login check
echo "Checking Cloudflare authentication..."
if ! cf pages whoami &> /dev/null; then
    echo "❌ Not logged in to Cloudflare"
    echo "Please run: cf pages login"
    exit 1
fi
echo "✅ Logged in to Cloudflare"
echo ""

# Check if wrangler is available for deployment
if ! command -v wrangler &> /dev/null; then
    echo "⚠️  Wrangler CLI not found. Installing..."
    npm install -g wrangler
fi

# Install dependencies
echo "Installing dependencies..."
npm install
echo ""

# Build the project
echo "Building project..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "✅ Build successful"
echo ""

# Ask for project name
echo "========================================="
echo "Choose deployment target:"
echo "========================================="
echo "1. Create new Pages project"
echo "2. Deploy to existing project"
echo ""
read -p "Select option (1-2): " -n 1 -r
echo ""

if [[ $REPLY =~ ^1$ ]]; then
    echo ""
    echo "========================================="
    echo "Creating new Pages project..."
    echo "========================================="
    echo ""

    read -p "Enter project name (e.g., driftlock-web): " PROJECT_NAME

    if [ -z "$PROJECT_NAME" ]; then
        echo "❌ Project name required"
        exit 1
    fi

    echo ""
    echo "Creating project: $PROJECT_NAME"
    cf pages project create $PROJECT_NAME --commit-source control

    if [ $? -ne 0 ]; then
        echo "❌ Failed to create project"
        exit 1
    fi

    echo ""
    echo "========================================="
    echo "Deploying to $PROJECT_NAME..."
    echo "========================================="
    cf pages deploy . --project-name $PROJECT_NAME

else
    echo ""
    echo "========================================="
    echo "Deploying to existing project..."
    echo "========================================="
    echo ""

    read -p "Enter existing project name: " PROJECT_NAME

    if [ -z "$PROJECT_NAME" ]; then
        echo "❌ Project name required"
        exit 1
    fi

    echo ""
    echo "Deploying to: $PROJECT_NAME"
    cf pages deploy . --project-name $PROJECT_NAME
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================="
    echo "🎉 DEPLOYMENT SUCCESSFUL!"
    echo "========================================="
    echo ""
    echo "Your site is live at:"
    cf pages domain list $PROJECT_NAME
    echo ""
    echo "Next steps:"
    echo "1. Configure environment variables in Pages dashboard"
    echo "2. Add custom domain (optional):"
    echo "   cf pages domain add $PROJECT_NAME www.driftlock.com"
    echo ""
    echo "3. Update Supabase CORS settings if using custom domain"
    echo ""
else
    echo "❌ Deployment failed"
    exit 1
fi
