#!/bin/bash

# Driftlock Stripe Setup Script
# Creates products and prices for the new tier structure:
#   Starter ($29), Pro ($99), Team ($249), Scale ($499), Enterprise (optional)
#
# Usage:
#   ./scripts/setup-stripe.sh           # Interactive mode
#   ./scripts/setup-stripe.sh --auto    # Non-interactive (requires STRIPE_API_KEY env var)

set -e

# =============================================================================
# Configuration
# =============================================================================

TIERS=(
    "starter:Starter:29:250K events/mo, 50 streams, 30-day retention"
    "pro:Pro:99:1.5M events/mo, 200 streams, 180-day retention"
    "team:Team:249:10M events/mo, 1,000 streams, 1-year retention"
    "scale:Scale:499:50M events/mo, 5,000 streams, 2-year retention"
)

# Optional: Enterprise tier (commented out by default - manual setup recommended)
# ENTERPRISE_TIER="enterprise:Enterprise:1500:Committed volume, custom limits, dedicated support"

OUTPUT_FILE="/tmp/driftlock-stripe-ids.env"

# =============================================================================
# Helper Functions
# =============================================================================

check_stripe_cli() {
    if ! command -v stripe &> /dev/null; then
        echo "❌ Stripe CLI not found. Install it:"
        echo "   macOS: brew install stripe/stripe-cli/stripe"
        echo "   Linux: See https://stripe.com/docs/stripe-cli"
        exit 1
    fi
    echo "✅ Stripe CLI found"
}

check_jq() {
    if ! command -v jq &> /dev/null; then
        echo "❌ jq not found. Install it:"
        echo "   macOS: brew install jq"
        echo "   Linux: apt install jq"
        exit 1
    fi
}

login_stripe() {
    if [ -n "$STRIPE_API_KEY" ]; then
        echo "✅ Using STRIPE_API_KEY from environment"
        return
    fi

    echo "🔐 Stripe API Key Required"
    echo "Get your key from: https://dashboard.stripe.com/apikeys"
    echo ""
    read -sp "Enter your Stripe secret key (sk_test_... or sk_live_...): " stripe_key
    echo ""

    if [[ ! $stripe_key == sk_test_* && ! $stripe_key == sk_live_* ]]; then
        echo "❌ Invalid key format. Should start with sk_test_ or sk_live_"
        exit 1
    fi

    export STRIPE_API_KEY="$stripe_key"
    echo "✅ Stripe API key configured"

    # Detect environment
    if [[ $stripe_key == sk_live_* ]]; then
        echo "⚠️  LIVE MODE DETECTED - This will create real billable products!"
        read -p "Continue? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            echo "Aborted."
            exit 1
        fi
    fi
}

# =============================================================================
# Product/Price Creation
# =============================================================================

create_product_and_price() {
    local tier_key=$1
    local tier_name=$2
    local price_cents=$3
    local description=$4

    echo ""
    echo "📦 Creating Driftlock $tier_name..."

    # Create product
    local product_json=$(stripe products create \
        --name="Driftlock $tier_name" \
        --description="$description" \
        --type="service" \
        --metadata="tier=$tier_key" \
        --metadata="service=driftlock" \
        2>&1)

    local product_id=$(echo "$product_json" | jq -r '.id // empty')
    if [ -z "$product_id" ]; then
        echo "❌ Failed to create product: $product_json"
        return 1
    fi
    echo "   Product: $product_id"

    # Create monthly price
    local price_json=$(stripe prices create \
        --product="$product_id" \
        --currency="usd" \
        --unit-amount="${price_cents}00" \
        --recurring-interval="month" \
        --metadata="tier=$tier_key" \
        2>&1)

    local price_id=$(echo "$price_json" | jq -r '.id // empty')
    if [ -z "$price_id" ]; then
        echo "❌ Failed to create price: $price_json"
        return 1
    fi
    echo "   Price:   $price_id (\$$price_cents/month)"

    # Save to output file
    local var_name="STRIPE_PRICE_ID_${tier_key^^}"
    echo "$var_name=$price_id" >> "$OUTPUT_FILE"

    return 0
}

create_all_products() {
    echo "🛍️  Creating Stripe Products and Prices"
    echo "========================================"

    # Clear output file
    > "$OUTPUT_FILE"
    echo "# Driftlock Stripe Configuration" >> "$OUTPUT_FILE"
    echo "# Generated: $(date -u +"%Y-%m-%dT%H:%M:%SZ")" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

    for tier_data in "${TIERS[@]}"; do
        IFS=':' read -r key name price desc <<< "$tier_data"
        create_product_and_price "$key" "$name" "$price" "$desc"
        if [ $? -ne 0 ]; then
            echo "❌ Failed to create tier: $name"
            exit 1
        fi
    done

    echo ""
    echo "✅ All products created successfully!"
}

# =============================================================================
# EU Data Residency Add-on (Optional)
# =============================================================================

create_eu_addon() {
    echo ""
    read -p "Create EU Data Residency add-on (\$150/mo)? (y/n): " create_eu
    if [ "$create_eu" != "y" ]; then
        return
    fi

    echo "📦 Creating EU Data Residency add-on..."

    local product_json=$(stripe products create \
        --name="Driftlock EU Data Residency" \
        --description="Host your Driftlock data in EU region (europe-west1) for GDPR/DORA compliance" \
        --type="service" \
        --metadata="addon=eu_residency" \
        --metadata="service=driftlock" \
        2>&1)

    local product_id=$(echo "$product_json" | jq -r '.id // empty')
    if [ -z "$product_id" ]; then
        echo "❌ Failed to create EU addon product"
        return 1
    fi

    local price_json=$(stripe prices create \
        --product="$product_id" \
        --currency="usd" \
        --unit-amount="15000" \
        --recurring-interval="month" \
        --metadata="addon=eu_residency" \
        2>&1)

    local price_id=$(echo "$price_json" | jq -r '.id // empty')
    echo "   EU Add-on Price: $price_id (\$150/month)"
    echo "" >> "$OUTPUT_FILE"
    echo "# Optional add-ons" >> "$OUTPUT_FILE"
    echo "STRIPE_PRICE_ID_EU_ADDON=$price_id" >> "$OUTPUT_FILE"
}

# =============================================================================
# Webhook Setup
# =============================================================================

setup_webhook() {
    echo ""
    echo "🪝 Webhook Configuration"
    echo "========================"
    echo ""
    echo "Your webhook endpoint should be:"
    echo "  https://driftlock.net/webhooks/stripe"
    echo ""
    echo "Required events:"
    echo "  - checkout.session.completed"
    echo "  - customer.subscription.updated"
    echo "  - customer.subscription.deleted"
    echo "  - invoice.paid"
    echo "  - invoice.payment_failed"
    echo ""

    read -p "Create webhook now? (y/n): " create_webhook
    if [ "$create_webhook" != "y" ]; then
        echo "⏭️  Skipping webhook creation. Create manually at:"
        echo "   https://dashboard.stripe.com/webhooks"
        return
    fi

    local webhook_url="https://driftlock.net/webhooks/stripe"
    read -p "Webhook URL [$webhook_url]: " custom_url
    webhook_url="${custom_url:-$webhook_url}"

    echo "Creating webhook endpoint: $webhook_url"
    local webhook_json=$(stripe webhook_endpoints create \
        --url="$webhook_url" \
        --enabled-events="checkout.session.completed,customer.subscription.updated,customer.subscription.deleted,invoice.paid,invoice.payment_failed" \
        2>&1)

    local webhook_id=$(echo "$webhook_json" | jq -r '.id // empty')
    local webhook_secret=$(echo "$webhook_json" | jq -r '.secret // empty')

    if [ -n "$webhook_id" ]; then
        echo "✅ Webhook created: $webhook_id"
        if [ -n "$webhook_secret" ]; then
            echo "" >> "$OUTPUT_FILE"
            echo "# Webhook" >> "$OUTPUT_FILE"
            echo "STRIPE_WEBHOOK_SECRET=$webhook_secret" >> "$OUTPUT_FILE"
            echo "   Secret saved to output file"
        else
            echo "⚠️  Copy webhook secret from dashboard: https://dashboard.stripe.com/webhooks"
        fi
    else
        echo "❌ Failed to create webhook: $webhook_json"
    fi
}

# =============================================================================
# Customer Portal Setup
# =============================================================================

configure_portal() {
    echo ""
    echo "🚪 Customer Portal"
    echo "=================="
    echo "Configure the customer portal at:"
    echo "  https://dashboard.stripe.com/settings/billing/portal"
    echo ""
    echo "Recommended settings:"
    echo "  ✓ Allow customers to switch plans (upgrade/downgrade)"
    echo "  ✓ Allow customers to cancel subscriptions"
    echo "  ✓ Allow customers to update payment methods"
    echo "  ✓ Show invoice history"
    echo ""
}

# =============================================================================
# Summary
# =============================================================================

print_summary() {
    echo ""
    echo "=========================================="
    echo "🎉 Stripe Setup Complete!"
    echo "=========================================="
    echo ""
    echo "Generated configuration saved to: $OUTPUT_FILE"
    echo ""
    cat "$OUTPUT_FILE"
    echo ""
    echo "=========================================="
    echo ""
    echo "📋 Next Steps:"
    echo ""
    echo "1. Add these values to your .env or GCP Secret Manager:"
    echo "   cat $OUTPUT_FILE"
    echo ""
    echo "2. For GCP deployment, add secrets:"
    echo '   while IFS="=" read -r key value; do'
    echo '     gcloud secrets create "${key,,}" --data-file=- --project=driftlock <<< "$value"'
    echo '   done < /tmp/driftlock-stripe-ids.env'
    echo ""
    echo "3. Configure Customer Portal:"
    echo "   https://dashboard.stripe.com/settings/billing/portal"
    echo ""
    echo "4. Test the integration:"
    echo "   ./scripts/test-stripe-checkout.sh"
    echo ""
}

# =============================================================================
# Manual Instructions (for reference)
# =============================================================================

manual_instructions() {
    echo ""
    echo "📋 Manual Stripe Setup"
    echo "======================"
    echo ""
    echo "Go to https://dashboard.stripe.com/products and create:"
    echo ""
    echo "1. Driftlock Starter - \$29/month"
    echo "   Description: 250K events/mo, 50 streams, 30-day retention"
    echo ""
    echo "2. Driftlock Pro - \$99/month"
    echo "   Description: 1.5M events/mo, 200 streams, 180-day retention"
    echo ""
    echo "3. Driftlock Team - \$249/month"
    echo "   Description: 10M events/mo, 1,000 streams, 1-year retention"
    echo ""
    echo "4. Driftlock Scale - \$499/month"
    echo "   Description: 50M events/mo, 5,000 streams, 2-year retention"
    echo ""
    echo "5. (Optional) EU Data Residency Add-on - \$150/month"
    echo "   Description: Host data in EU region for GDPR/DORA compliance"
    echo ""
    echo "Then go to https://dashboard.stripe.com/webhooks and create:"
    echo "  URL: https://driftlock.net/webhooks/stripe"
    echo "  Events: checkout.session.completed, customer.subscription.*,"
    echo "          invoice.paid, invoice.payment_failed"
    echo ""
    echo "Copy the price IDs and webhook secret to your environment."
}

# =============================================================================
# Main
# =============================================================================

main() {
    echo "🚀 Driftlock Stripe Setup"
    echo "========================="
    echo ""

    check_stripe_cli
    check_jq

    if [ "$1" == "--manual" ]; then
        manual_instructions
        exit 0
    fi

    echo ""
    echo "This script will create Stripe products for:"
    echo "  • Starter (\$29/mo) - 250K events"
    echo "  • Pro (\$99/mo) - 1.5M events"
    echo "  • Team (\$249/mo) - 10M events"
    echo "  • Scale (\$499/mo) - 50M events"
    echo ""

    if [ "$1" != "--auto" ]; then
        read -p "Continue? (y/n): " proceed
        if [ "$proceed" != "y" ]; then
            echo "Aborted."
            exit 0
        fi
    fi

    login_stripe
    create_all_products
    create_eu_addon
    setup_webhook
    configure_portal
    print_summary
}

main "$@"
