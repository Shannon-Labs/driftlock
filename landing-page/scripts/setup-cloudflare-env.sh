#!/bin/bash

# Script to help set up Cloudflare Pages environment variables for Firebase
# Usage: ./scripts/setup-cloudflare-env.sh

set -e

echo "🔧 Setting up Cloudflare Pages environment variables for Firebase"
echo "=================================================="

# Check if wrangler is installed
if ! command -v wrangler &> /dev/null; then
    echo "❌ Wrangler CLI not found. Please install with: npm install -g wrangler"
    exit 1
fi

# Project name from wrangler.toml
PROJECT_NAME="driftlock"

echo "📝 Please provide your Firebase configuration values:"
echo "You can find these in your Firebase Console > Project Settings > General > Your apps"
echo ""

read -p "Firebase API Key: " FIREBASE_API_KEY
read -p "Firebase Auth Domain (usually project-id.firebaseapp.com): " FIREBASE_AUTH_DOMAIN
read -p "Firebase Project ID: " FIREBASE_PROJECT_ID
read -p "Firebase Storage Bucket (usually project-id.appspot.com): " FIREBASE_STORAGE_BUCKET
read -p "Firebase Messaging Sender ID: " FIREBASE_MESSAGING_SENDER_ID
read -p "Firebase App ID: " FIREBASE_APP_ID
read -p "Firebase Measurement ID (optional, for Analytics): " FIREBASE_MEASUREMENT_ID

echo ""
echo "🔐 Setting environment variables in Cloudflare Pages..."

# Set environment variables for both production and preview environments
wrangler pages secret put VITE_FIREBASE_API_KEY --project-name="$PROJECT_NAME" <<< "$FIREBASE_API_KEY"
wrangler pages secret put VITE_FIREBASE_AUTH_DOMAIN --project-name="$PROJECT_NAME" <<< "$FIREBASE_AUTH_DOMAIN"
wrangler pages secret put VITE_FIREBASE_PROJECT_ID --project-name="$PROJECT_NAME" <<< "$FIREBASE_PROJECT_ID"
wrangler pages secret put VITE_FIREBASE_STORAGE_BUCKET --project-name="$PROJECT_NAME" <<< "$FIREBASE_STORAGE_BUCKET"
wrangler pages secret put VITE_FIREBASE_MESSAGING_SENDER_ID --project-name="$PROJECT_NAME" <<< "$FIREBASE_MESSAGING_SENDER_ID"
wrangler pages secret put VITE_FIREBASE_APP_ID --project-name="$PROJECT_NAME" <<< "$FIREBASE_APP_ID"

if [ -n "$FIREBASE_MEASUREMENT_ID" ]; then
    wrangler pages secret put VITE_FIREBASE_MEASUREMENT_ID --project-name="$PROJECT_NAME" <<< "$FIREBASE_MEASUREMENT_ID"
fi

echo ""
echo "✅ Environment variables set successfully!"
echo ""
echo "🚀 Next steps:"
echo "1. Trigger a new deployment in Cloudflare Pages (push to your connected git branch)"
echo "2. Or deploy manually with: npm run deploy:cloudflare"
echo ""
echo "📋 Alternative: Set variables manually in Cloudflare Dashboard:"
echo "   Dashboard → Pages → $PROJECT_NAME → Settings → Environment Variables"
echo ""
echo "Variables that were set:"
echo "- VITE_FIREBASE_API_KEY"
echo "- VITE_FIREBASE_AUTH_DOMAIN" 
echo "- VITE_FIREBASE_PROJECT_ID"
echo "- VITE_FIREBASE_STORAGE_BUCKET"
echo "- VITE_FIREBASE_MESSAGING_SENDER_ID"
echo "- VITE_FIREBASE_APP_ID"
[ -n "$FIREBASE_MEASUREMENT_ID" ] && echo "- VITE_FIREBASE_MEASUREMENT_ID"