import { createClient } from '@supabase/supabase-js';

// DriftLock API Gateway - Cloudflare Workers
export interface Env {
  ENV_SUPABASE_URL: string;
  ENV_SUPABASE_SERVICE_ROLE_KEY: string;
  ENV_GO_BACKEND_URL: string;
  ENV_JWT_SECRET: string;
}

// Initialize Supabase client
let supabase: any = null;

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    // Initialize Supabase client if not already done
    if (!supabase) {
      supabase = createClient(
        env.ENV_SUPABASE_URL,
        env.ENV_SUPABASE_SERVICE_ROLE_KEY,
        {
          auth: {
            persistSession: false,
            autoRefreshToken: false
          }
        }
      );
    }

    // Handle Stripe webhook requests
    if (request.method === 'POST' && (request.url.includes('/webhook') || request.url.endsWith('/'))) {
      // Check if this is a Stripe webhook
      const stripeSignature = request.headers.get('stripe-signature');
      if (stripeSignature) {
        return handleStripeWebhook(request, env, stripeSignature);
      }
    }

    // Handle regular API requests
    try {
      // Extract API key from Authorization header
      const authHeader = request.headers.get('Authorization');
      if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return new Response('Missing or invalid authorization header', { status: 401 });
      }

      const apiKey = authHeader.substring(7); // Remove "Bearer " prefix

      // Validate API key against Supabase
      const orgId = await validateApiKey(apiKey, env);
      if (!orgId) {
        return new Response('Invalid API key', { status: 401 });
      }

      // Check subscription status
      const subscriptionValid = await isSubscriptionValid(orgId, env);
      if (!subscriptionValid) {
        return new Response('Invalid subscription', { status: 402 });
      }

      // Determine if this request should be metered
      const shouldMeter = await isMeteredRequest(request);

      if (shouldMeter) {
        // Record usage before proxying
        const meterResult = await recordUsage(orgId, 1, env);
        if (!meterResult.success) {
          return new Response('Failed to record usage', { status: 500 });
        }
      }

      // Add organization context to request
      const newHeaders = new Headers(request.headers);
      newHeaders.set('X-Organization-ID', orgId);
      newHeaders.set('X-Forwarded-By', 'Cloudflare-Workers-Gateway');
      newHeaders.set('X-Real-IP', request.headers.get('CF-Connecting-IP') || '');
      
      // Get the path without the domain, just the path part
      const url = new URL(request.url);
      const path = url.pathname + url.search;

      // Proxy request to Go backend
      const backendRequest = new Request(`${env.ENV_GO_BACKEND_URL}${path}`, {
        method: request.method,
        headers: newHeaders,
        body: request.body,
      });

      // Fetch from the Go backend
      const backendResponse = await fetch(backendRequest);

      // Log audit event (do this in the background to not slow down the response)
      ctx.waitUntil(logAuditEvent(orgId, request, backendResponse, env));

      return backendResponse;

    } catch (error) {
      console.error('Gateway error:', error);
      return new Response('Internal server error', { status: 500 });
    }
  }
};

// Handle Stripe webhooks
async function handleStripeWebhook(request: Request, env: Env, signature: string): Promise<Response> {
  try {
    const payload = await request.text();
    
    // Verify the webhook signature (this would require access to the signing secret)
    // In a real implementation, you would use stripe.webhooks.constructEvent
    // For now, we'll just log the event and respond with success
    console.log('Received Stripe webhook:', request.headers.get('stripe-signature'));
    
    // Parse the event
    const event = JSON.parse(payload);
    
    // Log the webhook event for debugging
    console.log('Stripe event type:', event.type);
    console.log('Stripe event data:', event.data.object);
    
    // Process the event based on its type
    await processStripeEvent(event, env);
    
    // Response must be 200 to acknowledge receipt
    return new Response('OK', { status: 200 });
  } catch (err) {
    console.error('Webhook error:', err);
    return new Response('Webhook error', { status: 400 });
  }
}

// Process Stripe events
async function processStripeEvent(event: any, env: Env) {
  const eventType = event.type;
  const eventData = event.data.object;
  
  try {
    switch (eventType) {
      case 'checkout.session.completed':
        console.log('Handling checkout session completed:', eventData.id);
        // This is where you'd update your user's subscription in Supabase
        await handleCheckoutSessionCompleted(eventData, env);
        break;
        
      case 'customer.subscription.created':
      case 'customer.subscription.updated':
        console.log('Handling subscription update:', eventData.id);
        await handleSubscriptionUpdated(eventData, env);
        break;
        
      case 'customer.subscription.deleted':
        console.log('Handling subscription deletion:', eventData.id);
        await handleSubscriptionDeleted(eventData, env);
        break;
        
      case 'invoice.payment_succeeded':
        console.log('Handling payment succeeded:', eventData.id);
        await handleInvoicePaymentSucceeded(eventData, env);
        break;
        
      case 'invoice.payment_failed':
        console.log('Handling payment failed:', eventData.id);
        await handleInvoicePaymentFailed(eventData, env);
        break;
        
      default:
        console.log('Unhandled event type:', eventType);
    }
  } catch (err) {
    console.error(`Error processing ${eventType} event:`, err);
  }
}

// Handle checkout session completion
async function handleCheckoutSessionCompleted(session: any, env: Env) {
  // Extract customer and subscription info
  const customerId = session.customer;
  const subscriptionId = session.subscription;
  const organizationId = session.client_reference_id || null; // You can pass org ID via client_reference_id
  
  if (!organizationId) {
    console.error('No organization ID found in checkout session');
    return;
  }
  
  // Update the subscription in Supabase
  const { error } = await supabase
    .from('subscriptions')
    .update({
      stripe_subscription_id: subscriptionId,
      status: 'active',
      current_period_start: new Date(session.current_period_start * 1000).toISOString(),
      current_period_end: new Date(session.current_period_end * 1000).toISOString(),
    })
    .eq('organization_id', organizationId);
    
  if (error) {
    console.error('Error updating subscription:', error);
  }
}

// Handle subscription updated
async function handleSubscriptionUpdated(subscription: any, env: Env) {
  const organizationId = await getOrgIdFromStripeCustomer(subscription.customer);
  if (!organizationId) {
    console.error('No organization found for customer:', subscription.customer);
    return;
  }
  
  const { error } = await supabase
    .from('subscriptions')
    .update({
      status: subscription.status,
      current_period_start: new Date(subscription.current_period_start * 1000).toISOString(),
      current_period_end: new Date(subscription.current_period_end * 1000).toISOString(),
      plan: mapStripePriceToPlan(subscription.items.data[0]?.price.id) // Simplified mapping
    })
    .eq('stripe_subscription_id', subscription.id);
    
  if (error) {
    console.error('Error updating subscription:', error);
  }
}

// Handle subscription deleted
async function handleSubscriptionDeleted(subscription: any, env: Env) {
  const organizationId = await getOrgIdFromStripeCustomer(subscription.customer);
  if (!organizationId) {
    console.error('No organization found for customer:', subscription.customer);
    return;
  }
  
  const { error } = await supabase
    .from('subscriptions')
    .update({ status: 'canceled' })
    .eq('stripe_subscription_id', subscription.id);
    
  if (error) {
    console.error('Error updating subscription status:', error);
  }
}

// Handle invoice payment succeeded
async function handleInvoicePaymentSucceeded(invoice: any, env: Env) {
  // Update invoice status in Supabase
  const { error } = await supabase
    .from('invoices_mirror')
    .update({ 
      status: 'paid', 
      amount_paid_cents: invoice.amount_paid,
      paid_at: new Date().toISOString()
    })
    .eq('stripe_invoice_id', invoice.id);
    
  if (error) {
    console.error('Error updating invoice:', error);
  }
}

// Handle invoice payment failed
async function handleInvoicePaymentFailed(invoice: any, env: Env) {
  // Update invoice status in Supabase
  const { error } = await supabase
    .from('invoices_mirror')
    .update({ status: 'failed' })
    .eq('stripe_invoice_id', invoice.id);
    
  if (error) {
    console.error('Error updating invoice:', error);
  }
  
  // You might want to trigger dunning management here
  await handlePaymentFailure(invoice.customer, env);
}

// Helper to get organization ID from Stripe customer ID
async function getOrgIdFromStripeCustomer(stripeCustomerId: string): Promise<string | null> {
  const { data, error } = await supabase
    .from('billing_customers')
    .select('organization_id')
    .eq('stripe_customer_id', stripeCustomerId)
    .single();
    
  if (error) {
    console.error('Error getting org ID from customer:', error);
    return null;
  }
  
  return data?.organization_id || null;
}

// Map Stripe price IDs to your plan names
function mapStripePriceToPlan(priceId: string): string {
  // This mapping should match your Stripe price IDs
  // You'll need to customize this based on your actual price IDs
  const planMap: Record<string, string> = {
    'price_standard': 'standard',
    'price_growth': 'growth',
    'price_enterprise': 'enterprise',
    // Add your actual price IDs here
  };
  
  return planMap[priceId] || 'unknown';
}

// Handle payment failure (update dunning status, etc.)
async function handlePaymentFailure(stripeCustomerId: string, env: Env) {
  const orgId = await getOrgIdFromStripeCustomer(stripeCustomerId);
  if (!orgId) return;
  
  // Update dunning status in Supabase (you'll need to implement this table)
  // This is where you'd implement your dunning management logic
  console.log(`Payment failed for organization: ${orgId}`);
}

// Validate API key against Supabase
async function validateApiKey(apiKey: string, env: Env): Promise<string | null> {
  try {
    // In a real implementation, you would hash the API key before comparing
    // For now, we'll use a simplified version - in production you should hash the key
    const { data, error } = await supabase
      .from('api_keys')
      .select('organization_id')
      .eq('key_hash', await hashApiKey(apiKey)) // Assuming keys are stored hashed
      .single();

    if (error) {
      console.error('API key validation error:', error);
      return null;
    }

    return data?.organization_id || null;
  } catch (error) {
    console.error('Error validating API key:', error);
    return null;
  }
}

// Check if organization has valid subscription
async function isSubscriptionValid(organizationId: string, env: Env): Promise<boolean> {
  try {
    const { data, error } = await supabase
      .from('subscriptions')
      .select('status')
      .eq('organization_id', organizationId)
      .in('status', ['active', 'trialing'])
      .single();

    if (error || !data) {
      return false;
    }

    return ['active', 'trialing'].includes(data.status);
  } catch (error) {
    console.error('Error checking subscription status:', error);
    return false;
  }
}

// Determine if request should be metered
async function isMeteredRequest(request: Request): Promise<boolean> {
  const url = new URL(request.url);
  const path = url.pathname;

  // Only meter anomaly detection endpoints
  // The Stream API is free, Monitor/Anomalies are billable
  const billableEndpoints = [
    '/api/v1/anomalies',           // Anomaly detection
    '/api/v1/anomalies/detect',    // Direct detection endpoint
    '/api/v1/analyze',             // Analysis endpoints
    '/api/v1/monitor',             // Monitoring endpoints
  ];

  return billableEndpoints.some(endpoint => path.startsWith(endpoint));
}

// Record usage metering
async function recordUsage(organizationId: string, count: number, env: Env): Promise<{ success: boolean }> {
  try {
    // Call the metering endpoint in Supabase
    // This assumes you'll have the metering-ingest function set up in Supabase
    const response = await fetch(`${env.ENV_SUPABASE_URL}/functions/v1/metering-ingest`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${env.ENV_SUPABASE_SERVICE_ROLE_KEY}`,
      },
      body: JSON.stringify({
        event: 'anomaly_detection',
        count: count,
        organization_id: organizationId
      })
    });

    return { success: response.ok };
  } catch (error) {
    console.error('Error recording usage:', error);
    return { success: false };
  }
}

// Log audit event
async function logAuditEvent(organizationId: string, request: Request, response: Response, env: Env): Promise<void> {
  try {
    const { error } = await supabase
      .from('audit_logs')
      .insert({
        organization_id: organizationId,
        action: `${request.method} ${new URL(request.url).pathname}`,
        resource_type: 'api_call',
        ip_address: request.headers.get('CF-Connecting-IP') || '',
        user_agent: request.headers.get('User-Agent') || '',
        details: {
          method: request.method,
          path: new URL(request.url).pathname,
          status_code: response.status,
          content_length: response.headers.get('Content-Length'),
          cf_ray: request.headers.get('CF-Ray')
        }
      });

    if (error) {
      console.error('Audit log error:', error);
    }
  } catch (error) {
    console.error('Failed to log audit event:', error);
  }
}

// Simple API key hashing (using Web Crypto API available in Cloudflare Workers)
async function hashApiKey(apiKey: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(apiKey);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}