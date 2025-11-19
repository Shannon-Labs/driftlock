#!/bin/bash

# Driftlock Stripe Setup Script
# This script helps set up Stripe products and pricing for Driftlock

set -e

echo "🚀 Setting up Stripe for Driftlock..."
echo ""

# Function to check if stripe CLI is installed
check_stripe_cli() {
    if ! command -v stripe &> /dev/null; then
        echo "❌ Stripe CLI not found. Please install it first:"
        echo "   macOS: brew install stripe/stripe-cli/stripe"
        echo "   Linux: curl -s https://packages.stripe.com/api/security/keypairs/stripe-cli-gpg/public | gpg --dearmor | sudo tee /usr/share/keyrings/stripe.gpg"
        echo "          echo \"deb [signed-by=/usr/share/keyrings/stripe.gpg] https://packages.stripe.com/stripe-cli-debian-local stable main\" | sudo tee -a /etc/apt/sources.list.d/stripe.list"
        echo "          sudo apt update"
        echo "          sudo apt install stripe"
        echo "   Download: https://stripe.com/docs/stripe-cli"
        exit 1
    fi
    echo "✅ Stripe CLI found"
}

# Function to login to Stripe
login_stripe() {
    echo "🔐 Logging into Stripe..."
    echo "Please enter your Stripe API keys:"
    echo "You can find these at https://dashboard.stripe.com/apikeys"
    echo ""

    read -p "Enter your Stripe secret key (sk_test_...): " stripe_key

    if [[ ! $stripe_key == sk_test_* && ! $stripe_key == sk_live_* ]]; then
        echo "❌ Invalid Stripe key format. Should start with sk_test_ or sk_live_"
        exit 1
    fi

    echo "$stripe_key" > /tmp/stripe-key.txt
    chmod 600 /tmp/stripe-key.txt
    export STRIPE_API_KEY="$stripe_key"

    echo "✅ Stripe API key configured"
    echo ""
}

# Function to create products and prices
create_stripe_products() {
    echo "🛍️  Creating Stripe products and prices..."
    echo ""

    # Create Pro product
    echo "Creating 'Pro' product..."
    product_id=$(stripe products create \
        --name="Driftlock Pro" \
        --description="Professional anomaly detection for production workloads" \
        --type="service" \
        --metadata="tier=pro,service=driftlock" \
        --json | jq -r '.id')

    echo "✅ Created product: $product_id"

    # Create price for Pro plan ($99/month)
    echo "Creating price for Pro plan..."
    price_id=$(stripe prices create \
        --product="$product_id" \
        --currency="usd" \
        --unit-amount=9900 \
        --recurring-interval="month" \
        --metadata="tier=pro,service=driftlock" \
        --json | jq -r '.id')

    echo "✅ Created price: $price_id ($99/month)"
    echo ""

    # Save to file for GCP secrets
    echo "$price_id" > /tmp/stripe-price-id-pro.txt

    echo "📋 Stripe Configuration Summary:"
    echo "Product ID: $product_id"
    echo "Price ID (Pro): $price_id"
    echo "Amount: $99/month"
    echo ""

    echo "🔐 To add these to GCP Secret Manager, run:"
    echo "echo -n '$stripe_key' | gcloud secrets create stripe-secret-key --data-file=- --project=driftlock"
    echo "echo -n '$price_id' | gcloud secrets create stripe-price-id-pro --data-file=- --project=driftlock"
    echo ""
}

# Function to create webhooks
setup_webhooks() {
    echo "🪝 Setting up Stripe webhooks..."
    echo ""

    read -p "Enter your deployed API URL (e.g., https://driftlock-api-xxxxx-uc.a.run.app): " api_url

    if [ -z "$api_url" ]; then
        echo "⚠️  Skipping webhook setup (no API URL provided)"
        echo "You can set up webhooks later at: https://dashboard.stripe.com/webhooks"
        return
    fi

    webhook_url="$api_url/stripe/webhook"

    echo "Creating webhook endpoint: $webhook_url"
    webhook_id=$(stripe webhooks create \
        --url="$webhook_url" \
        --enabled-events="customer.subscription.created,customer.subscription.updated,customer.subscription.deleted,invoice.payment_succeeded,invoice.payment_failed" \
        --json | jq -r '.id')

    echo "✅ Created webhook: $webhook_id"
    echo ""

    echo "⚠️  IMPORTANT: Webhook signing secret will be displayed in Stripe dashboard"
    echo "Go to https://dashboard.stripe.com/webhooks to copy the signing secret"
    echo "Then add it as a GCP secret: stripe-webhook-secret"
    echo ""
}

# Function to create test data
create_test_data() {
    echo "🧪 Creating test customer and subscription..."
    echo ""

    # Create test customer
    customer_id=$(stripe customers create \
        --name="Test Customer" \
        --email="test@driftlock.dev" \
        --metadata="test=true" \
        --json | jq -r '.id')

    echo "✅ Created test customer: $customer_id"

    # Get the price ID we created earlier
    if [ -f /tmp/stripe-price-id-pro.txt ]; then
        price_id=$(cat /tmp/stripe-price-id-pro.txt)

        # Create test subscription
        subscription_id=$(stripe subscriptions create \
            --customer="$customer_id" \
            --items="price=$price_id" \
            --metadata="test=true" \
            --json | jq -r '.id')

        echo "✅ Created test subscription: $subscription_id"
        echo ""
        echo "🧪 Test data created successfully!"
    else
        echo "⚠️  No price ID found, skipping test subscription"
    fi
}

# Function to provide manual setup instructions
manual_setup_instructions() {
    echo "📋 Manual Stripe Setup Instructions"
    echo "================================="
    echo ""
    echo "1. Go to https://dashboard.stripe.com/products"
    echo "2. Click 'Add product'"
    echo "3. Configure product:"
    echo "   - Name: Driftlock Pro"
    echo "   - Description: Professional anomaly detection for production workloads"
    echo "   - Pricing: $99/month (or your desired price)"
    echo "4. Save the product and note the Price ID"
    echo ""
    echo "5. Go to https://dashboard.stripe.com/apikeys"
    echo "6. Copy your secret key (sk_test_... for testing)"
    echo ""
    echo "7. Go to https://dashboard.stripe.com/webhooks"
    echo "8. Add webhook endpoint: [your-api-url]/stripe/webhook"
    echo "9. Select events:"
    echo "   - customer.subscription.created"
    echo "   - customer.subscription.updated"
    echo "   - customer.subscription.deleted"
    echo "   - invoice.payment_succeeded"
    echo "   - invoice.payment_failed"
    echo ""
    echo "10. Copy the webhook signing secret"
    echo ""
    echo "Then add these to GCP Secret Manager:"
    echo "- stripe-secret-key"
    echo "- stripe-price-id-pro"
    echo "- stripe-webhook-secret"
}

# Main execution
check_stripe_cli

echo "Choose setup option:"
echo "1) Automated setup with Stripe CLI"
echo "2) Manual setup instructions"

read -p "Enter your choice (1-2): " choice

case $choice in
    1)
        login_stripe
        create_stripe_products
        setup_webhooks
        create_test_data
        ;;
    2)
        manual_setup_instructions
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "🎉 Stripe setup complete!"
echo ""
echo "Next steps:"
echo "1. Add your Stripe secrets to GCP Secret Manager"
echo "2. Deploy your API to get webhook URL"
echo "3. Test the integration with a real payment"