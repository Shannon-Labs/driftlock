import { createClient } from '@supabase/supabase-js'
import { Hono } from 'hono'

// Environment variables interface
export interface Env {
  SUPABASE_URL: string
  SUPABASE_SERVICE_ROLE_KEY: string
  SUPABASE_ANON_KEY: string
  STRIPE_WEBHOOK_SECRET?: string
  CORS_ORIGIN?: string
}

// Initialize Supabase client
function createSupabaseClient(env: Env) {
  return createClient(
    env.SUPABASE_URL,
    env.SUPABASE_SERVICE_ROLE_KEY,
    {
      auth: {
        persistSession: false,
      },
    }
  )
}

// Initialize Hono app
const app = new Hono()

// CORS middleware
app.use('*', async (c, next) => {
  await next()
  const corsOrigin = (c.env as unknown as Env).CORS_ORIGIN || '*'
  c.header('Access-Control-Allow-Origin', corsOrigin)
  c.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
  c.header('Access-Control-Allow-Headers', 'Content-Type, Authorization, Stripe-Signature')
})

// Security headers
app.use('*', async (c, next) => {
  await next()
  c.header('X-Content-Type-Options', 'nosniff')
  c.header('X-Frame-Options', 'DENY')
  c.header('X-XSS-Protection', '1; mode=block')
  c.header('Strict-Transport-Security', 'max-age=31536000; includeSubDomains')
  // Adjust CSP as appropriate for your app frontends
  c.header('Content-Security-Policy', "default-src 'self'")
})

// Simple in-memory rate limit (best-effort per isolate)
const rlBuckets = new Map<string, {count: number; reset: number}>()
const RL_LIMIT = 100 // requests
const RL_PERIOD_MS = 60 * 1000 // 1 minute

app.use('*', async (c, next) => {
  const ip = c.req.header('cf-connecting-ip') || 'unknown'
  const now = Date.now()
  const b = rlBuckets.get(ip)
  if (!b || now > b.reset) {
    rlBuckets.set(ip, { count: 1, reset: now + RL_PERIOD_MS })
  } else {
    b.count++
    if (b.count > RL_LIMIT) {
      return c.json({ error: 'Rate limit exceeded' }, 429)
    }
  }
  await next()
})

// Handle OPTIONS for CORS preflight
app.options('*', (c) => {
  return c.text('', 204)
})

// Health check endpoint
app.get('/health', (c) => {
  return c.json({
    status: 'healthy',
    service: 'driftlock-api',
    timestamp: new Date().toISOString(),
    version: '1.0.0',
  })
})

// Create anomaly endpoint
app.post('/anomalies', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const body = await c.req.json()

    const { organization_id, event_type, severity, metadata, anomaly_score } = body

    // Validate required fields
    if (!organization_id || !event_type) {
      return c.json({ error: 'Missing required fields: organization_id, event_type' }, 400)
    }

    // Insert anomaly into Supabase
    const { data, error } = await supabase
      .from('anomalies')
      .insert({
        organization_id,
        event_type,
        severity: severity || 'medium',
        metadata: metadata || {},
        anomaly_score: anomaly_score || 0.5,
        created_at: new Date().toISOString(),
      })
      .select()
      .single()

    if (error) {
      console.error('Error creating anomaly:', error)
      return c.json({ error: 'Failed to create anomaly', details: error.message }, 500)
    }

    // Call meter-usage edge function if this is a true anomaly
    try {
      await supabase.functions.invoke('meter-usage', {
        body: {
          organization_id,
          count: 1,
          anomaly: true,
        },
      })
    } catch (usageError) {
      console.error('Error metering usage:', usageError)
      // Don't fail the anomaly creation if metering fails
    }

    return c.json(data, 201)
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Get anomalies endpoint
app.get('/anomalies', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const orgId = c.req.query('organization_id')
    const limit = parseInt(c.req.query('limit') || '50')
    const offset = parseInt(c.req.query('offset') || '0')

    if (!orgId) {
      return c.json({ error: 'organization_id is required' }, 400)
    }

    // Query anomalies for the organization
    const { data, error } = await supabase
      .from('anomalies')
      .select('*')
      .eq('organization_id', orgId)
      .order('created_at', { ascending: false })
      .range(offset, offset + limit - 1)

    if (error) {
      console.error('Error fetching anomalies:', error)
      return c.json({ error: 'Failed to fetch anomalies', details: error.message }, 500)
    }

    return c.json({
      data,
      pagination: {
        limit,
        offset,
        count: data?.length || 0,
      },
    })
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Track usage endpoint
app.post('/usage', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const body = await c.req.json()

    const { organization_id, count, is_anomaly } = body

    if (!organization_id || count === undefined) {
      return c.json({ error: 'Missing required fields: organization_id, count' }, 400)
    }

    // Only meter if anomaly is true (pay-for-anomalies model)
    if (is_anomaly) {
      const { data, error } = await supabase.functions.invoke('meter-usage', {
        body: {
          organization_id,
          count,
          anomaly: true,
        },
      })

      if (error) {
        console.error('Error metering usage:', error)
        return c.json({ error: 'Failed to meter usage', details: error.message }, 500)
      }

      return c.json({ success: true, data })
    }

    return c.json({ success: true, message: 'No billing required (not an anomaly)' })
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Get usage statistics endpoint
app.get('/usage', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const orgId = c.req.query('organization_id')

    if (!orgId) {
      return c.json({ error: 'organization_id is required' }, 400)
    }

    // Get current period usage
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const period_start = startOfMonth.toISOString()

    const { data, error } = await supabase
      .from('usage_counters')
      .select('*')
      .eq('organization_id', orgId)
      .eq('period_start', period_start)
      .single()

    if (error && error.code !== 'PGRST116') { // PGRST116 is "no rows returned"
      console.error('Error fetching usage:', error)
      return c.json({ error: 'Failed to fetch usage', details: error.message }, 500)
    }

    // Get subscription info
    const { data: subscription } = await supabase
      .from('subscriptions')
      .select('*, plan_price_map(*)')
      .eq('organization_id', orgId)
      .eq('status', 'active')
      .single()

    return c.json({
      usage: data || {
        organization_id: orgId,
        period_start,
        anomaly_count: 0,
        estimated_charges_cents: 0,
      },
      subscription,
    })
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Stripe webhook endpoint
app.post('/stripe-webhook', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const body = await c.req.text()
    const signature = c.req.header('stripe-signature')
    const env = c.env as unknown as Env

    if (!signature || !env.STRIPE_WEBHOOK_SECRET) {
      return c.json({ error: 'Missing Stripe signature or webhook secret' }, 400)
    }

    // Verify webhook with Supabase edge function
    const { data, error } = await supabase.functions.invoke('stripe-webhook', {
      body: {
        payload: body,
        signature,
      },
    })

    if (error) {
      console.error('Error processing Stripe webhook:', error)
      return c.json({ error: 'Failed to process webhook', details: error.message }, 500)
    }

    return c.json({ received: true, data })
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Get subscription info endpoint
app.get('/subscription', async (c) => {
  try {
    const supabase = createSupabaseClient(c.env as unknown as Env)
    const orgId = c.req.query('organization_id')

    if (!orgId) {
      return c.json({ error: 'organization_id is required' }, 400)
    }

    const { data, error } = await supabase
      .from('subscriptions')
      .select('*, plan_price_map(*)')
      .eq('organization_id', orgId)
      .eq('status', 'active')
      .single()

    if (error && error.code !== 'PGRST116') {
      console.error('Error fetching subscription:', error)
      return c.json({ error: 'Failed to fetch subscription', details: error.message }, 500)
    }

    return c.json(data || null)
  } catch (err) {
    console.error('Unexpected error:', err)
    return c.json({ error: 'Internal server error' }, 500)
  }
})

// Root endpoint
app.get('/', (c) => {
  return c.json({
    name: 'DriftLock API',
    version: '1.0.0',
    status: 'running',
    endpoints: [
      'GET /health - Health check',
      'POST /anomalies - Create anomaly',
      'GET /anomalies - Get anomalies',
      'POST /usage - Track usage',
      'GET /usage - Get usage stats',
      'GET /subscription - Get subscription',
      'POST /stripe-webhook - Stripe webhook handler',
    ],
  })
})

export default app
