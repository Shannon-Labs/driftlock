#!/bin/bash
set -e

# Set the specific service account invoker to avoid "allUsers" org policy violation
# Including all potential service accounts
export FUNCTIONS_INVOKERS="firebase-hosting@system.gserviceaccount.com,service-131489574303@gcp-sa-firebasehosting.iam.gserviceaccount.com,driftlock@appspot.gserviceaccount.com"

echo "Deploying Firebase Functions with restricted invokers: $FUNCTIONS_INVOKERS"

# Deploy only functions
npx -y firebase-tools@latest deploy --only functions --project driftlock
