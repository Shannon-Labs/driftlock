# DriftLock Billing System Documentation

This document describes the billing system implementation in DriftLock, including subscription management, payment processing, and usage tracking.

## Overview

The billing system provides subscription management and payment processing powered by Stripe. It supports multiple pricing tiers with both fixed and usage-based billing components.

## Architecture

```
Frontend (React) -- Stripe Elements
         |
         v
DriftLock API (Go) -- Stripe API
         |
         v
Database (PostgreSQL) -- Subscription Records
         |
         v
Webhooks -- Event Processing
```

## Components

### Subscription Plans

The system includes three predefined plans:

1. **Free Plan** ($0/month)
   - 10,000 events/month
   - 7-day retention
   - 100 alerts/month
   - Basic support

2. **Pro Plan** ($29/month)
   - 100,000 events/month
   - 30-day retention
   - 1,000 alerts/month
   - Advanced support
   - Custom dashboards
   - Anomaly prediction

3. **Business Plan** ($99/month)
   - Unlimited events
   - 90-day retention
   - Unlimited alerts
   - Priority support
   - SSO integration
   - Dedicated account manager

### Database Schema

The billing system uses the following database tables:

- `subscription_plans`: Plan definitions
- `user_subscriptions`: Active subscriptions
- `features`: Features available in each plan
- `usage_records`: Usage tracking for metered billing

### API Endpoints

#### GET /billing/plans
Retrieve available subscription plans

#### POST /billing/checkout
Create a Stripe checkout session for a plan

#### GET /billing/subscription
Get current user's subscription

#### DELETE /billing/subscription
Cancel the current subscription

#### GET /billing/usage
Get current usage information

#### POST /billing/usage
Record usage (typically called internally)

## Stripe Integration

### Webhooks

The system handles these Stripe webhooks:
- `customer.subscription.created` - New subscription
- `customer.subscription.updated` - Subscription changed
- `customer.subscription.deleted` - Subscription canceled
- `invoice.payment_succeeded` - Payment processed
- `invoice.payment_failed` - Payment failed

### Security

- Webhook signatures are verified using Stripe's signing secret
- All payment information is handled directly by Stripe
- No sensitive payment data is stored on our servers

## Usage-Based Billing

The system tracks usage for metered billing:
- Events ingested
- API requests
- Storage used
- Other feature-specific metrics

Usage is reset at the start of each billing cycle.

## Customer Portal

Users can access the Stripe customer portal to:
- Update payment methods
- Change subscription plans
- Update billing information
- Download invoices

## Testing

For testing, use Stripe's test API keys and test cards:
- 4242424242424242 (valid card)
- 4000000000000002 (declined card)

## Error Handling

The billing system handles these error conditions:
- Payment failures
- Insufficient plan quotas
- Expired payment methods
- Failed webhooks
- API rate limits