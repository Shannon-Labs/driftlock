/**
 * Stripe Checkout API Route
 * Creates checkout sessions for compliance report purchases
 */

import { NextResponse } from 'next/server';
import Stripe from 'stripe';
import { auth } from '@/lib/auth';

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY, {
  apiVersion: '2023-10-16',
});

// Product catalog
const PRODUCTS = {
  dora_quarterly: {
    name: 'DORA Quarterly Compliance Report',
    description: 'Comprehensive DORA compliance report for regulatory submission',
    price: 29900, // $299.00 in cents
    currency: 'usd',
    metadata: {
      type: 'dora_quarterly',
      validity: '90d', // 90 days validity
    },
  },
  nis2_incident: {
    name: 'NIS2 Incident Report',
    description: 'NIS2 regulatory incident reporting documentation',
    price: 19900, // $199.00 in cents
    currency: 'usd',
    metadata: {
      type: 'nis2_incident',
      validity: '60d',
    },
  },
  eu_ai_audit: {
    name: 'EU AI Act Audit Trail',
    description: 'EU AI Act compliance audit trail and documentation',
    price: 14900, // $149.00 in cents
    currency: 'usd',
    metadata: {
      type: 'eu_ai_audit',
      validity: '60d',
    },
  },
  annual_package: {
    name: 'Annual Compliance Package',
    description: 'All compliance reports for 12 months + priority support',
    price: 149900, // $1,499.00 in cents
    currency: 'usd',
    metadata: {
      type: 'annual_package',
      validity: '365d',
    },
  },
};

export async function POST(request) {
  try {
    // Authenticate user
    const session = await auth();
    if (!session?.user) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const body = await request.json();
    const { productId, customerEmail, customerId, metadata = {} } = body;

    // Validate product
    const product = PRODUCTS[productId];
    if (!product) {
      return NextResponse.json({ error: 'Invalid product' }, { status: 400 });
    }

    // Check if customer already has an active subscription for this product
    if (productId === 'annual_package') {
      const existingSubscription = await checkExistingSubscription(customerId);
      if (existingSubscription) {
        return NextResponse.json(
          { error: 'Active subscription already exists' },
          { status: 409 }
        );
      }
    }

    // Create Stripe checkout session
    const checkoutSession = await stripe.checkout.sessions.create({
      payment_method_types: ['card'],
      mode: productId === 'annual_package' ? 'subscription' : 'payment',
      line_items: [
        {
          price_data: {
            currency: product.currency,
            product_data: {
              name: product.name,
              description: product.description,
              metadata: {
                ...product.metadata,
                customer_id: customerId,
                generated_by: 'shannon-labs-compliance',
              },
            },
            unit_amount: product.price,
            recurring: productId === 'annual_package' ? {
              interval: 'year',
              interval_count: 1,
            } : undefined,
          },
          quantity: 1,
        },
      ],
      customer_email: customerEmail,
      success_url: `${process.env.NEXT_PUBLIC_APP_URL}/success?session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${process.env.NEXT_PUBLIC_APP_URL}/pricing?cancelled=true`,
      metadata: {
        customer_id: customerId,
        product_id: productId,
        user_id: session.user.id,
        ...metadata,
      },
      subscription_data: productId === 'annual_package' ? {
        metadata: {
          customer_id: customerId,
          product_id: productId,
          user_id: session.user.id,
        },
        trial_period_days: 0,
      } : undefined,
      allow_promotion_codes: true,
      billing_address_collection: 'required',
      customer_creation: productId !== 'annual_package' ? 'always' : undefined,
    });

    // Store pending order in database
    if (productId !== 'annual_package') {
      await storePendingOrder({
        sessionId: checkoutSession.id,
        customerId,
        productId,
        amount: product.price,
        currency: product.currency,
        status: 'pending',
        metadata: checkoutSession.metadata,
      });
    }

    return NextResponse.json({
      sessionId: checkoutSession.id,
      url: checkoutSession.url,
    });

  } catch (error) {
    console.error('Checkout session creation failed:', error);
    return NextResponse.json(
      { error: 'Failed to create checkout session' },
      { status: 500 }
    );
  }
}

export async function GET(request) {
  try {
    const { searchParams } = new URL(request.url);
    const sessionId = searchParams.get('session_id');

    if (!sessionId) {
      return NextResponse.json({ error: 'Session ID required' }, { status: 400 });
    }

    // Retrieve checkout session
    const session = await stripe.checkout.sessions.retrieve(sessionId, {
      expand: ['line_items', 'customer'],
    });

    // Check session ownership
    const authSession = await auth();
    if (!authSession?.user || session.metadata.user_id !== authSession.user.id) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    return NextResponse.json({
      session: {
        id: session.id,
        status: session.status,
        payment_status: session.payment_status,
        customer_email: session.customer_email || session.customer?.email,
        amount_total: session.amount_total,
        currency: session.currency,
        created: session.created,
        metadata: session.metadata,
      },
    });

  } catch (error) {
    console.error('Checkout session retrieval failed:', error);
    return NextResponse.json(
      { error: 'Failed to retrieve checkout session' },
      { status: 500 }
    );
  }
}

// Helper functions
async function checkExistingSubscription(customerId) {
  // Check if customer has active annual subscription
  try {
    const subscriptions = await stripe.subscriptions.list({
      customer: customerId,
      status: 'active',
      limit: 1,
    });

    return subscriptions.data.length > 0;
  } catch (error) {
    console.error('Error checking existing subscription:', error);
    return false;
  }
}

async function storePendingOrder(orderData) {
  // Store pending order in database for fulfillment
  // This would typically use your ORM (Prisma, etc.)
  try {
    // Example implementation (adjust based on your database schema)
    /*
    await db.pendingOrder.create({
      data: {
        sessionId: orderData.sessionId,
        customerId: orderData.customerId,
        productId: orderData.productId,
        amount: orderData.amount,
        currency: orderData.currency,
        status: orderData.status,
        metadata: orderData.metadata,
        createdAt: new Date(),
        expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000), // 24 hours
      },
    });
    */
    console.log('Stored pending order:', orderData.sessionId);
  } catch (error) {
    console.error('Error storing pending order:', error);
    // Don't fail the request if storage fails
  }
}

// Price lookup for products
export async function getProductPrice(productId) {
  const product = PRODUCTS[productId];
  if (!product) {
    throw new Error('Invalid product');
  }
  return {
    amount: product.price,
    currency: product.currency,
    name: product.name,
    description: product.description,
  };
}

// Validate checkout session
export async function validateCheckoutSession(sessionId, customerId) {
  try {
    const session = await stripe.checkout.sessions.retrieve(sessionId);

    if (session.payment_status !== 'paid') {
      throw new Error('Payment not completed');
    }

    if (session.metadata.customer_id !== customerId) {
      throw new Error('Customer mismatch');
    }

    return session;
  } catch (error) {
    console.error('Checkout session validation failed:', error);
    throw error;
  }
}