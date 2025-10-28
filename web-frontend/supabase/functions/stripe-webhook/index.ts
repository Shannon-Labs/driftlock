import { serve } from "https://deno.land/std@0.224.0/http/server.ts";
import Stripe from 'https://esm.sh/stripe@14.21.0';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2.45.0';

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type, stripe-signature',
};

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response(null, { headers: corsHeaders });
  }

  try {
    const stripe = new Stripe(Deno.env.get('STRIPE_SECRET_KEY')!, {
      apiVersion: '2023-10-16',
    });

    const supabase = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
    );

    const signature = req.headers.get('stripe-signature');
    if (!signature) {
      console.error('Missing stripe-signature header');
      return new Response('Missing signature', { status: 400, headers: corsHeaders });
    }

    const body = await req.text();
    const webhookSecret = Deno.env.get('STRIPE_WEBHOOK_SECRET')!;
    
    let event: Stripe.Event;
    try {
      event = stripe.webhooks.constructEvent(body, signature, webhookSecret);
    } catch (err: any) {
      console.error('Webhook signature verification failed:', err.message);
      return new Response(`Webhook Error: ${err.message}`, { status: 400, headers: corsHeaders });
    }

    console.log('Processing webhook event:', event.type, event.id);

    // Store webhook event for idempotency
    const { error: eventError } = await supabase
      .from('billing_events')
      .insert({
        stripe_event_id: event.id,
        event_type: event.type,
        payload: event.data.object,
        processed_at: new Date().toISOString(),
      })
      .select()
      .single();

    // If event already processed (duplicate), return success
    if (eventError?.code === '23505') {
      console.log('Event already processed:', event.id);
      return new Response(JSON.stringify({ received: true, duplicate: true }), {
        status: 200,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    // Handle different event types
    switch (event.type) {
      case 'checkout.session.completed':
        await handleCheckoutCompleted(event.data.object as Stripe.Checkout.Session, stripe, supabase);
        break;
      
      case 'customer.subscription.created':
      case 'customer.subscription.updated':
        await handleSubscriptionUpdated(event.data.object as Stripe.Subscription, supabase);
        break;
      
      case 'customer.subscription.deleted':
        await handleSubscriptionDeleted(event.data.object as Stripe.Subscription, supabase);
        break;
      
      case 'invoice.paid':
        await handleInvoicePaid(event.data.object as Stripe.Invoice, supabase);
        break;
      
      case 'invoice.payment_failed':
        await handleInvoicePaymentFailed(event.data.object as Stripe.Invoice, supabase);
        break;
      
      default:
        console.log('Unhandled event type:', event.type);
    }

    return new Response(JSON.stringify({ received: true }), {
      status: 200,
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  } catch (error: any) {
    console.error('Error processing webhook:', error);
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500,
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  }
});

async function handleCheckoutCompleted(
  session: Stripe.Checkout.Session,
  stripe: Stripe,
  supabase: any
) {
  console.log('Processing checkout.session.completed:', session.id);

  const customerId = session.customer as string;
  const subscriptionId = session.subscription as string;
  
  // Get organization_id from metadata
  const organizationId = session.metadata?.organization_id;
  if (!organizationId) {
    console.error('No organization_id in session metadata');
    return;
  }

  // Get full subscription details
  const subscription = await stripe.subscriptions.retrieve(subscriptionId);
  
  // Determine plan from subscription
  const plan = getPlanFromSubscription(subscription);

  // Create or update billing customer
  await supabase
    .from('billing_customers')
    .upsert({
      organization_id: organizationId,
      stripe_customer_id: customerId,
      billing_email: session.customer_details?.email,
      company_name: session.customer_details?.name,
      updated_at: new Date().toISOString(),
    });

  // Create or update subscription
  await supabase
    .from('subscriptions')
    .upsert({
      organization_id: organizationId,
      stripe_subscription_id: subscriptionId,
      plan: plan.tier,
      status: subscription.status,
      current_period_start: new Date(subscription.current_period_start * 1000).toISOString(),
      current_period_end: new Date(subscription.current_period_end * 1000).toISOString(),
      included_calls: plan.included_calls,
      overage_rate_per_call: plan.overage_rate,
      has_launch50_promo: session.metadata?.promo_code === 'LAUNCH50',
      updated_at: new Date().toISOString(),
    });

  // Create usage counter for this period
  await supabase
    .from('usage_counters')
    .insert({
      organization_id: organizationId,
      period_start: new Date(subscription.current_period_start * 1000).toISOString(),
      period_end: new Date(subscription.current_period_end * 1000).toISOString(),
      total_calls: 0,
      included_calls_used: 0,
      overage_calls: 0,
      estimated_charges_cents: 0,
    });

  console.log('Checkout completed successfully for org:', organizationId);
}

async function handleSubscriptionUpdated(subscription: Stripe.Subscription, supabase: any) {
  console.log('Processing subscription update:', subscription.id);

  // Get organization from customer
  const { data: customer } = await supabase
    .from('billing_customers')
    .select('organization_id')
    .eq('stripe_customer_id', subscription.customer)
    .single();

  if (!customer) {
    console.error('No organization found for customer:', subscription.customer);
    return;
  }

  const plan = getPlanFromSubscription(subscription);

  await supabase
    .from('subscriptions')
    .update({
      plan: plan.tier,
      status: subscription.status,
      current_period_start: new Date(subscription.current_period_start * 1000).toISOString(),
      current_period_end: new Date(subscription.current_period_end * 1000).toISOString(),
      included_calls: plan.included_calls,
      overage_rate_per_call: plan.overage_rate,
      updated_at: new Date().toISOString(),
    })
    .eq('stripe_subscription_id', subscription.id);

  console.log('Subscription updated for org:', customer.organization_id);
}

async function handleSubscriptionDeleted(subscription: Stripe.Subscription, supabase: any) {
  console.log('Processing subscription deletion:', subscription.id);

  await supabase
    .from('subscriptions')
    .update({
      status: 'canceled',
      updated_at: new Date().toISOString(),
    })
    .eq('stripe_subscription_id', subscription.id);
}

async function handleInvoicePaid(invoice: Stripe.Invoice, supabase: any) {
  console.log('Processing invoice.paid:', invoice.id);

  const { data: customer } = await supabase
    .from('billing_customers')
    .select('organization_id')
    .eq('stripe_customer_id', invoice.customer)
    .single();

  if (!customer) return;

  await supabase
    .from('invoices_mirror')
    .upsert({
      organization_id: customer.organization_id,
      stripe_invoice_id: invoice.id,
      status: invoice.status || 'paid',
      amount_due_cents: invoice.amount_due,
      amount_paid_cents: invoice.amount_paid,
      hosted_invoice_url: invoice.hosted_invoice_url,
      invoice_pdf_url: invoice.invoice_pdf,
      finalized_at: invoice.status_transitions?.finalized_at 
        ? new Date(invoice.status_transitions.finalized_at * 1000).toISOString()
        : null,
      paid_at: invoice.status_transitions?.paid_at
        ? new Date(invoice.status_transitions.paid_at * 1000).toISOString()
        : null,
      period_start: invoice.period_start
        ? new Date(invoice.period_start * 1000).toISOString()
        : null,
      period_end: invoice.period_end
        ? new Date(invoice.period_end * 1000).toISOString()
        : null,
    });
}

async function handleInvoicePaymentFailed(invoice: Stripe.Invoice, supabase: any) {
  console.log('Processing invoice.payment_failed:', invoice.id);

  const { data: customer } = await supabase
    .from('billing_customers')
    .select('organization_id')
    .eq('stripe_customer_id', invoice.customer)
    .single();

  if (!customer) return;

  // Update subscription status
  if (invoice.subscription) {
    await supabase
      .from('subscriptions')
      .update({
        status: 'past_due',
        updated_at: new Date().toISOString(),
      })
      .eq('stripe_subscription_id', invoice.subscription);
  }

  // Mirror invoice
  await supabase
    .from('invoices_mirror')
    .upsert({
      organization_id: customer.organization_id,
      stripe_invoice_id: invoice.id,
      status: 'open',
      amount_due_cents: invoice.amount_due,
      amount_paid_cents: invoice.amount_paid,
      hosted_invoice_url: invoice.hosted_invoice_url,
    });
}

function getPlanFromSubscription(subscription: Stripe.Subscription): {
  tier: string;
  included_calls: number;
  overage_rate: number;
} {
  // Map based on price IDs or product metadata
  const priceId = subscription.items.data[0]?.price.id;
  
  // Website pricing: $49 Pro, $249 Enterprise
  const planMapping: Record<string, any> = {
    'price_1SMhsZL4rhSbUSqA51lWvPlQ': { tier: 'pro', included_calls: 50000, overage_rate: 0.001 },
    'price_1SMhshL4rhSbUSqAyHfhWUSQ': { tier: 'enterprise', included_calls: 500000, overage_rate: 0.0005 },
  };

  // Check product metadata for tier info
  const product = subscription.items.data[0]?.price.product;
  if (typeof product === 'object' && product.metadata?.tier) {
    const tier = product.metadata.tier;
    if (tier === 'pro') return { tier: 'pro', included_calls: 50000, overage_rate: 0.001 };
    if (tier === 'enterprise') return { tier: 'enterprise', included_calls: 500000, overage_rate: 0.0005 };
  }

  // Check by price ID
  if (planMapping[priceId]) {
    return planMapping[priceId];
  }

  // Default to pro
  return { tier: 'pro', included_calls: 50000, overage_rate: 0.001 };
}
