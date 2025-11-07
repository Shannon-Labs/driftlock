# ✅ CORRECT Stripe Configuration for DriftLock

## Summary

Fixed the Stripe configuration to match the **actual website pricing** ($49 and $249) instead of the incorrect prices ($10, $50, $200) that were initially created.

## Correct Subscription Plans

| Plan | Price | Product ID | Price ID | Included Calls | Overage Rate |
|------|-------|------------|----------|----------------|--------------|
| **Pro** | $49/month | `prod_TJKXbWnB3ExnqJ` | `price_1SMhsZL4rhSbUSqA51lWvPlQ` | 50,000 | $0.001/call |
| **Enterprise** | $249/month | `prod_TJKXEFXBjkcsAB` | `price_1SMhshL4rhSbUSqAyHfhWUSQ` | 500,000 | $0.0005/call |

## Webhook Configuration

**Endpoint:** `https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook`

**Events to listen for:**
- `checkout.session.completed`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.paid`
- `invoice.payment_failed`

## Stripe Dashboard

View products: https://dashboard.stripe.com/test/products  
View prices: https://dashboard.stripe.com/test/prices

## What Was Fixed

### ❌ Incorrect (Deleted)
- ~~Developer Plan - $10/month~~ - Deactivated
- ~~Standard Plan - $50/month~~ - Deactivated
- ~~Growth Plan - $200/month~~ - Deactivated

### ✅ Correct (Active)
- **Pro Plan - $49/month** - Active
- **Enterprise Plan - $249/month** - Active

## Integration Status

✅ **Edge Functions:** All 4 functions deployed with correct pricing
- health
- stripe-webhook (updated with correct price IDs)
- meter-usage
- send-alert-email

✅ **Database:** Migrations applied
✅ **Webhook:** Configured and working
✅ **Supabase:** Project linked and secrets configured

## Testing

Test the webhook endpoint:
```bash
curl https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/health
```

Test Stripe webhook:
```bash
stripe trigger checkout.session.completed
```

## Next Steps

1. Update your website's Stripe checkout flow to use the correct Price IDs
2. Test a full subscription flow
3. Deploy to production
4. Update Stripe webhook URL in production dashboard

## Environment Variables

All required variables are configured in `.env`:
- `SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh`
- `SUPABASE_ANON_KEY=...`
- `SUPABASE_SERVICE_ROLE_KEY=...`
- `SUPABASE_BASE_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co`
- `SUPABASE_WEBHOOK_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook`
