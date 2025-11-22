#!/bin/bash
set -e

echo "🚀 Starting Driftlock Deployment..."

# 1. Build the Landing Page
echo "📦 Building landing page..."
cd landing-page
npm install
npm run build
cd ..

# 2. Check if Firebase CLI is installed
if ! command -v firebase &> /dev/null; then
    echo "❌ Firebase CLI not found. Please install it: npm install -g firebase-tools"
    echo "Then run 'firebase login' before running this script again."
    exit 1
fi

# 3. Deploy to Firebase Hosting (Frontend)
echo "🔥 Deploying frontend to Firebase..."
firebase deploy --only hosting

# 4. Build and Deploy Backend to Cloud Run
if command -v gcloud &> /dev/null; then
    echo "☁️  Deploying Backend to Cloud Run..."
    
    # Check if PROJECT_ID is set or get from gcloud config
    PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
    if [ -z "$PROJECT_ID" ]; then
        echo "❌ No Google Cloud project selected. Run 'gcloud config set project YOUR_PROJECT_ID' first."
        exit 1
    fi
    
    echo "   Using project: $PROJECT_ID"
    
    # Submit build to Cloud Build (builds, pushes, and deploys based on cloudbuild.yaml)
    # We use the root directory as context
    echo "   Submitting build to Cloud Build..."
    SHORT_SHA=$(git rev-parse --short HEAD)
    gcloud builds submit --config cloudbuild.yaml . --substitutions=SHORT_SHA=$SHORT_SHA
    
    echo "   Backend deployed successfully!"
else
    echo "⚠️  gcloud CLI not found. Skipping backend deployment."
fi

echo "✅ Deployment complete!"
echo "   Frontend should be live at your Firebase URL."

