# Driftlock → Stripe Setup Prompt (for an AI with Stripe access)

You have access to Stripe Dashboard / Stripe API. Configure Driftlock’s self-serve subscriptions to match the published “Costco pricing” tiers, then return the exact IDs and settings needed to wire production.

## Context
- Driftlock backend creates Stripe Checkout Sessions by `plan` and maps that to Stripe price IDs via env vars:
  - `STRIPE_PRICE_ID_STARTER`
  - `STRIPE_PRICE_ID_PRO`
  - `STRIPE_PRICE_ID_TEAM`
  - `STRIPE_PRICE_ID_SCALE`
  - (optional) `STRIPE_PRICE_ID_ENTERPRISE`
- Stripe webhooks are received at: `https://driftlock.net/webhooks/stripe`
  - (Firebase rewrites → backend `/v1/billing/webhook`)
- Required webhook events:
  - `checkout.session.completed`
  - `customer.subscription.updated`
  - `customer.subscription.deleted`
  - `invoice.paid`
  - `invoice.payment_failed`

## What to build in Stripe (USD, monthly recurring)
Create **Products** and **Prices**:
1. **Driftlock Starter** — `$29/mo`
2. **Driftlock Pro** — `$99/mo`
3. **Driftlock Team** — `$249/mo`
4. **Driftlock Scale** — `$499/mo`

Optional (only if we want a self-serve Enterprise checkout):
5. **Driftlock Enterprise** — choose a published “starting at” price (e.g. `$1,500/mo`) and mark it clearly as “Starting at (annual contract preferred)”.

Add-on (optional, but recommended to prepare):
- **EU Data Residency** add-on — `$150/mo` recurring add-on price, attachable to Starter/Pro/Team/Scale subscriptions.

## Portal + Checkout configuration
1. Configure **Customer Portal** so customers can:
   - Upgrade/downgrade between Starter/Pro/Team/Scale
   - Add/remove EU Data Residency add-on
2. Ensure proration behavior is reasonable (Stripe default is fine unless specified).

## Webhook
1. Create a webhook endpoint pointing to `https://driftlock.net/webhooks/stripe`
2. Subscribe to the events listed above.
3. Capture and return the webhook signing secret (`whsec_...`).

## Deliverable (must return these exact values)
Return a structured response with:
- `STRIPE_SECRET_KEY` (live key only if you have it; otherwise confirm where it will be set)
- `STRIPE_WEBHOOK_SECRET` (`whsec_...`)
- `STRIPE_PRICE_ID_STARTER`
- `STRIPE_PRICE_ID_PRO`
- `STRIPE_PRICE_ID_TEAM`
- `STRIPE_PRICE_ID_SCALE`
- `STRIPE_PRICE_ID_ENTERPRISE` (if created)
- EU add-on: product ID + price ID (if created)
- Confirmation that the webhook events are enabled on the endpoint.

## Notes / guardrails
- Keep pricing flat monthly (no usage-based metering).
- Avoid automatic overages; Driftlock enforces caps in-app.
- If you create annual prices too, include their price IDs, but monthly is required.

