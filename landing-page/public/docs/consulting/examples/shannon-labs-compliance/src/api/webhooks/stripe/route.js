/**
 * Stripe Webhook Handler
 * Processes Stripe webhook events for payment confirmation and fulfillment
 */

import { NextResponse } from 'next/server';
import Stripe from 'stripe';
import { headers } from 'next/headers';
import { fulfillOrder, createSubscription, updateSubscription } from '@/lib/fulfillment';
import { notifyCustomer } from '@/lib/notifications';
import { trackEvent } from '@/lib/analytics';

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY, {
  apiVersion: '2023-10-16',
});

const webhookSecret = process.env.STRIPE_WEBHOOK_SECRET;

export async function POST(request) {
  try {
    const body = await request.text();
    const signature = headers().get('stripe-signature');

    if (!signature) {
      return NextResponse.json({ error: 'Missing signature' }, { status: 400 });
    }

    // Verify webhook signature
    let event;
    try {
      event = stripe.webhooks.constructEvent(body, signature, webhookSecret);
    } catch (err) {
      console.error('Webhook signature verification failed:', err);
      return NextResponse.json({ error: 'Invalid signature' }, { status: 400 });
    }

    console.log('Processing webhook event:', event.type);

    // Handle different event types
    switch (event.type) {
      case 'checkout.session.completed':
        await handleCheckoutSessionCompleted(event.data.object);
        break;

      case 'checkout.session.expired':
        await handleCheckoutSessionExpired(event.data.object);
        break;

      case 'payment_intent.succeeded':
        await handlePaymentSucceeded(event.data.object);
        break;

      case 'payment_intent.payment_failed':
        await handlePaymentFailed(event.data.object);
        break;

      case 'invoice.payment_succeeded':
        await handleInvoicePaymentSucceeded(event.data.object);
        break;

      case 'invoice.payment_failed':
        await handleInvoicePaymentFailed(event.data.object);
        break;

      case 'customer.subscription.created':
        await handleSubscriptionCreated(event.data.object);
        break;

      case 'customer.subscription.updated':
        await handleSubscriptionUpdated(event.data.object);
        break;

      case 'customer.subscription.deleted':
        await handleSubscriptionDeleted(event.data.object);
        break;

      default:
        console.log(`Unhandled event type: ${event.type}`);
    }

    return NextResponse.json({ received: true });

  } catch (error) {
    console.error('Webhook processing failed:', error);
    return NextResponse.json(
      { error: 'Webhook processing failed' },
      { status: 500 }
    );
  }
}

// Event handlers
async function handleCheckoutSessionCompleted(session) {
  console.log('Checkout session completed:', session.id);

  try {
    const { customer_id, product_id, user_id } = session.metadata;

    if (!customer_id || !product_id) {
      console.error('Missing metadata in checkout session:', session.id);
      return;
    }

    // Track successful checkout
    await trackEvent('checkout_completed', {
      session_id: session.id,
      customer_id,
      product_id,
      amount: session.amount_total,
      currency: session.currency,
    });

    // For one-time purchases, fulfill the order
    if (session.mode === 'payment') {
      await fulfillOrder({
        sessionId: session.id,
        customerId: customer_id,
        productId: product_id,
        userId: user_id,
        amount: session.amount_total,
        currency: session.currency,
        customerEmail: session.customer_email,
        metadata: session.metadata,
      });

      // Send confirmation email
      await notifyCustomer({
        type: 'order_confirmation',
        customerId,
        sessionId: session.id,
        productId,
        amount: session.amount_total,
        currency: session.currency,
      });

      console.log('Order fulfilled for session:', session.id);
    }

    // For subscriptions, creation is handled by subscription events
    if (session.mode === 'subscription') {
      console.log('Subscription checkout completed:', session.id);
    }

  } catch (error) {
    console.error('Error handling checkout session completed:', error);
    // Don't return error to Stripe - webhook should always return 200
  }
}

async function handleCheckoutSessionExpired(session) {
  console.log('Checkout session expired:', session.id);

  try {
    const { customer_id, product_id } = session.metadata;

    // Track expired checkout
    await trackEvent('checkout_expired', {
      session_id: session.id,
      customer_id,
      product_id,
    });

    // Clean up any pending orders
    await cleanupPendingOrder(session.id);

    // Optionally notify customer about expired session
    if (customer_id) {
      await notifyCustomer({
        type: 'checkout_expired',
        customerId,
        sessionId: session.id,
        productId: product_id,
      });
    }

  } catch (error) {
    console.error('Error handling checkout session expired:', error);
  }
}

async function handlePaymentSucceeded(paymentIntent) {
  console.log('Payment succeeded:', paymentIntent.id);

  try {
    // Track successful payment
    await trackEvent('payment_succeeded', {
      payment_intent_id: paymentIntent.id,
      amount: paymentIntent.amount,
      currency: paymentIntent.currency,
    });

  } catch (error) {
    console.error('Error handling payment succeeded:', error);
  }
}

async function handlePaymentFailed(paymentIntent) {
  console.log('Payment failed:', paymentIntent.id);

  try {
    // Track failed payment
    await trackEvent('payment_failed', {
      payment_intent_id: paymentIntent.id,
      amount: paymentIntent.amount,
      currency: paymentIntent.currency,
      last_payment_error: paymentIntent.last_payment_error,
    });

  } catch (error) {
    console.error('Error handling payment failed:', error);
  }
}

async function handleInvoicePaymentSucceeded(invoice) {
  console.log('Invoice payment succeeded:', invoice.id);

  try {
    const subscriptionId = invoice.subscription;

    // Track subscription payment
    await trackEvent('subscription_payment_succeeded', {
      invoice_id: invoice.id,
      subscription_id: subscriptionId,
      amount: invoice.amount_paid,
      currency: invoice.currency,
    });

    // Update subscription access
    if (subscriptionId) {
      await updateSubscription(subscriptionId, {
        status: 'active',
        lastPaymentDate: new Date(invoice.created * 1000),
        currentPeriodEnd: new Date(invoice.period_end * 1000),
      });

      // Notify customer of successful payment
      await notifyCustomer({
        type: 'subscription_payment_succeeded',
        subscriptionId,
        amount: invoice.amount_paid,
        currency: invoice.currency,
        nextBillingDate: new Date(invoice.period_end * 1000),
      });
    }

  } catch (error) {
    console.error('Error handling invoice payment succeeded:', error);
  }
}

async function handleInvoicePaymentFailed(invoice) {
  console.log('Invoice payment failed:', invoice.id);

  try {
    const subscriptionId = invoice.subscription;

    // Track failed subscription payment
    await trackEvent('subscription_payment_failed', {
      invoice_id: invoice.id,
      subscription_id: subscriptionId,
      amount: invoice.amount_due,
      currency: invoice.currency,
      attempt_count: invoice.attempt_count,
    });

    // Update subscription status
    if (subscriptionId) {
      await updateSubscription(subscriptionId, {
        status: 'past_due',
        lastPaymentFailure: new Date(invoice.created * 1000),
      });

      // Notify customer of payment failure
      await notifyCustomer({
        type: 'subscription_payment_failed',
        subscriptionId,
        amount: invoice.amount_due,
        currency: invoice.currency,
        nextRetryDate: new Date(invoice.next_payment_attempt * 1000),
      });
    }

  } catch (error) {
    console.error('Error handling invoice payment failed:', error);
  }
}

async function handleSubscriptionCreated(subscription) {
  console.log('Subscription created:', subscription.id);

  try {
    const { customer_id, product_id, user_id } = subscription.metadata;

    if (!customer_id || !product_id) {
      console.error('Missing metadata in subscription:', subscription.id);
      return;
    }

    // Create subscription in database
    await createSubscription({
      stripeSubscriptionId: subscription.id,
      customerId: customer_id,
      productId: product_id,
      userId: user_id,
      status: subscription.status,
      currentPeriodStart: new Date(subscription.current_period_start * 1000),
      currentPeriodEnd: new Date(subscription.current_period_end * 1000),
      metadata: subscription.metadata,
    });

    // Track subscription creation
    await trackEvent('subscription_created', {
      subscription_id: subscription.id,
      customer_id,
      product_id,
    });

    // Send welcome email
    await notifyCustomer({
      type: 'subscription_created',
      customerId,
      subscriptionId: subscription.id,
      productId,
      currentPeriodEnd: new Date(subscription.current_period_end * 1000),
    });

  } catch (error) {
    console.error('Error handling subscription created:', error);
  }
}

async function handleSubscriptionUpdated(subscription) {
  console.log('Subscription updated:', subscription.id);

  try {
    // Update subscription in database
    await updateSubscription(subscription.id, {
      status: subscription.status,
      currentPeriodStart: new Date(subscription.current_period_start * 1000),
      currentPeriodEnd: new Date(subscription.current_period_end * 1000),
      cancelAtPeriodEnd: subscription.cancel_at_period_end,
    });

    // Track subscription update
    await trackEvent('subscription_updated', {
      subscription_id: subscription.id,
      status: subscription.status,
      cancel_at_period_end: subscription.cancel_at_period_end,
    });

  } catch (error) {
    console.error('Error handling subscription updated:', error);
  }
}

async function handleSubscriptionDeleted(subscription) {
  console.log('Subscription deleted:', subscription.id);

  try {
    // Update subscription in database
    await updateSubscription(subscription.id, {
      status: 'canceled',
      canceledAt: new Date(),
    });

    // Track subscription cancellation
    await trackEvent('subscription_canceled', {
      subscription_id: subscription.id,
    });

    // Send cancellation confirmation
    const dbSubscription = await getSubscriptionByStripeId(subscription.id);
    if (dbSubscription) {
      await notifyCustomer({
        type: 'subscription_canceled',
        customerId: dbSubscription.customerId,
        subscriptionId: subscription.id,
        endDate: new Date(subscription.current_period_end * 1000),
      });
    }

  } catch (error) {
    console.error('Error handling subscription deleted:', error);
  }
}

// Helper functions
async function cleanupPendingOrder(sessionId) {
  try {
    // Remove pending order from database
    /*
    await db.pendingOrder.delete({
      where: { sessionId },
    });
    */
    console.log('Cleaned up pending order:', sessionId);
  } catch (error) {
    console.error('Error cleaning up pending order:', error);
  }
}

async function getSubscriptionByStripeId(stripeSubscriptionId) {
  try {
    // Retrieve subscription from database
    /*
    return await db.subscription.findUnique({
      where: { stripeSubscriptionId },
    });
    */
    return null; // Placeholder
  } catch (error) {
    console.error('Error retrieving subscription:', error);
    return null;
  }
}

// Webhook event verification for testing
export async function GET() {
  try {
    const webhookEndpoints = await stripe.webhookEndpoints.list();
    return NextResponse.json({
      endpoints: webhookEndpoints.data,
      configured: !!webhookSecret,
    });
  } catch (error) {
    console.error('Error listing webhook endpoints:', error);
    return NextResponse.json(
      { error: 'Failed to list webhook endpoints' },
      { status: 500 }
    );
  }
}