#!/bin/bash
# Test Stripe webhooks locally
# Usage: ./test-stripe-webhooks.sh [command]
#
# Commands:
#   listen   - Start Stripe CLI listener (default)
#   checkout - Trigger checkout.session.completed
#   payment  - Trigger invoice.payment_succeeded
#   failed   - Trigger invoice.payment_failed
#   trial    - Trigger customer.subscription.trial_will_end
#   all      - Trigger all test events
#
# Prerequisites:
# - Stripe CLI installed: brew install stripe/stripe-cli/stripe
# - Logged into Stripe CLI: stripe login
# - Local API running on port 8080

set -e

LOCAL_URL="${LOCAL_URL:-http://localhost:8080/api/v1/billing/webhook}"
COMMAND="${1:-listen}"

case "$COMMAND" in
  listen)
    echo "Starting Stripe webhook listener..."
    echo "Forwarding to: $LOCAL_URL"
    echo ""
    echo "In another terminal, trigger test events:"
    echo "  ./test-stripe-webhooks.sh checkout"
    echo "  ./test-stripe-webhooks.sh payment"
    echo "  ./test-stripe-webhooks.sh failed"
    echo "  ./test-stripe-webhooks.sh trial"
    echo ""
    stripe listen --forward-to "$LOCAL_URL"
    ;;

  checkout)
    echo "Triggering checkout.session.completed..."
    stripe trigger checkout.session.completed
    ;;

  payment)
    echo "Triggering invoice.payment_succeeded..."
    stripe trigger invoice.payment_succeeded
    ;;

  failed)
    echo "Triggering invoice.payment_failed..."
    stripe trigger invoice.payment_failed
    ;;

  trial)
    echo "Triggering customer.subscription.trial_will_end..."
    stripe trigger customer.subscription.trial_will_end
    ;;

  all)
    echo "Triggering all test events..."
    echo ""
    echo "1/4: checkout.session.completed"
    stripe trigger checkout.session.completed
    sleep 2
    echo ""
    echo "2/4: invoice.payment_succeeded"
    stripe trigger invoice.payment_succeeded
    sleep 2
    echo ""
    echo "3/4: invoice.payment_failed"
    stripe trigger invoice.payment_failed
    sleep 2
    echo ""
    echo "4/4: customer.subscription.trial_will_end"
    stripe trigger customer.subscription.trial_will_end
    echo ""
    echo "All events triggered!"
    ;;

  *)
    echo "Unknown command: $COMMAND"
    echo "Usage: $0 [listen|checkout|payment|failed|trial|all]"
    exit 1
    ;;
esac
