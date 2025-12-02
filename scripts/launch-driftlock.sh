#!/bin/bash

# Driftlock Launch Script
# One-command launch helper

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Clear screen
clear

echo ""
echo "══════════════════════════════════════════════════════════════"
echo "  🚀 Driftlock Launch Script"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "This script will help you launch Driftlock today."
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."
echo ""

# Check Firebase CLI
if command -v firebase &> /dev/null; then
    echo -e "${GREEN}✅ Firebase CLI installed${NC}"
else
    echo -e "${RED}❌ Firebase CLI not installed${NC}"
    echo "   Install: npm install -g firebase-tools"
    exit 1
fi

# Check Google Cloud CLI
if command -v gcloud &> /dev/null; then
    echo -e "${GREEN}✅ Google Cloud CLI installed${NC}"
else
    echo -e "${RED}❌ Google Cloud CLI not installed${NC}"
    echo "   Install: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Check Stripe CLI
if command -v stripe &> /dev/null; then
    echo -e "${GREEN}✅ Stripe CLI installed${NC}"
else
    echo -e "${YELLOW}⚠️  Stripe CLI not installed${NC}"
    echo "   Optional for webhook testing: brew install stripe/stripe-cli/stripe"
fi

echo ""

# Menu
echo "══════════════════════════════════════════════════════════════"
echo "  What would you like to do?"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "1) Run launch readiness tests"
echo "2) Set up Stripe webhooks"
echo "3) Re-authenticate Firebase & Google Cloud"
echo "4) Deploy everything"
echo "5) Test Stripe webhooks locally"
echo "6) Full launch sequence (tests → deploy → verify)"
echo "7) Exit"
echo ""
read -p "Enter your choice (1-7): " choice

case $choice in
    1)
        echo ""
        echo "Running launch readiness tests..."
        ./scripts/test-launch-readiness.sh
        ;;
    
    2)
        echo ""
        echo "Setting up Stripe webhooks..."
        echo ""
        echo "Step 1: Configure Stripe webhook in Dashboard"
        echo "  URL: https://driftlock.net/webhooks/stripe"
        echo "  Events: subscription.created, subscription.deleted, invoice.payment_succeeded, invoice.payment_failed"
        echo ""
        read -p "Have you configured the webhook in Stripe Dashboard? (y/n): " configured
        if [[ $configured == "y" ]]; then
            read -p "Enter webhook signing secret (whsec_...): " webhook_secret
            echo ""
            echo "Storing webhook secret..."
            cd functions
            firebase functions:config:set stripe.webhook_secret="$webhook_secret"
            cd ..
            echo -e "${GREEN}✅ Webhook secret configured${NC}"
        else
            echo "Please configure webhook first: https://dashboard.stripe.com/webhooks"
        fi
        ;;
    
    3)
        echo ""
        echo "Re-authenticating..."
        echo ""
        echo -e "${BLUE}Firebase authentication:${NC}"
        firebase login --reauth
        echo ""
        echo -e "${BLUE}Google Cloud authentication:${NC}"
        gcloud auth login
        gcloud config set project driftlock
        echo -e "${GREEN}✅ Re-authentication complete${NC}"
        ;;
    
    4)
        echo ""
        echo "Deploying everything..."
        echo ""
        
        # Build and deploy functions
        echo -e "${BLUE}Deploying Firebase Functions...${NC}"
        cd functions
        npm install
        npm run build
        firebase deploy --only functions
        cd ..
        
        # Build and deploy hosting
        echo -e "${BLUE}Deploying Firebase Hosting...${NC}"
        cd landing-page
        npm install
        npm run build
        cd ..
        firebase deploy --only hosting
        
        echo -e "${GREEN}✅ Deployment complete${NC}"
        ;;
    
    5)
        echo ""
        echo "Testing Stripe webhooks locally..."
        echo ""
        ./scripts/stripe-webhook-forward.sh
        ;;
    
    6)
        echo ""
        echo "🚀 Starting FULL LAUNCH SEQUENCE..."
        echo ""
        echo "Step 1: Running readiness tests..."
        if ./scripts/test-launch-readiness.sh; then
            echo ""
            echo "Step 2: Deploying..."
            
            # Deploy functions
            cd functions
            npm install
            npm run build
            firebase deploy --only functions
            cd ..
            
            # Deploy hosting
            cd landing-page
            npm install
            npm run build
            cd ..
            firebase deploy --only hosting
            
            echo ""
            echo "Step 3: Verifying deployment..."
            sleep 10
            ./scripts/test-launch-readiness.sh
            
            echo ""
            echo -e "${GREEN}🎉 LAUNCH SEQUENCE COMPLETE!${NC}"
            echo ""
            echo "Next steps:"
            echo "1. Test user signup flow"
            echo "2. Test Stripe checkout (use test card 4242 4242 4242 4242)"
            echo "3. Monitor logs for any issues"
            echo "4. Celebrate! 🎊"
        else
            echo -e "${RED}❌ Launch cancelled - tests failed${NC}"
            echo "Fix the issues and try again."
        fi
        ;;
    
    7)
        echo "Exiting..."
        exit 0
        ;;
    
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo "══════════════════════════════════════════════════════════════"
echo "  Done!"
echo "══════════════════════════════════════════════════════════════"
echo ""
