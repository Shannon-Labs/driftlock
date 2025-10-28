import { serve } from "https://deno.land/std@0.224.0/http/server.ts";
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2.45.0';
import { z } from 'https://esm.sh/zod@3.22.4';

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
};

const MeterUsageSchema = z.object({
  organization_id: z.string().uuid({ message: 'Invalid organization ID' }),
  count: z.number().int().positive().max(10000, { message: 'Count must be between 1 and 10000' }).default(1),
});

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response(null, { headers: corsHeaders });
  }

  try {
    // Verify service role authentication
    const authHeader = req.headers.get('Authorization');
    if (!authHeader?.startsWith('Bearer ')) {
      console.error('Unauthorized: Missing or invalid Authorization header');
      return new Response(JSON.stringify({ error: 'Unauthorized' }), {
        status: 401,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    const supabase = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
    );

    // Validate input
    const body = await req.json();
    const validationResult = MeterUsageSchema.safeParse(body);
    
    if (!validationResult.success) {
      console.error('Validation error:', validationResult.error.errors);
      return new Response(JSON.stringify({ 
        error: 'Invalid input',
        details: validationResult.error.errors.map(e => e.message),
      }), {
        status: 400,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    const { organization_id, count } = validationResult.data;

    console.log('Metering usage for org:', organization_id, 'count:', count);

    // Get current subscription
    const { data: subscription, error: subError } = await supabase
      .from('subscriptions')
      .select('*')
      .eq('organization_id', organization_id)
      .eq('status', 'active')
      .single();

    if (subError || !subscription) {
      console.error('No active subscription found:', subError);
      return new Response(JSON.stringify({ error: 'No active subscription' }), {
        status: 403,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    const periodStart = subscription.current_period_start;
    const periodEnd = subscription.current_period_end;

    // Get or create usage counter
    let { data: usageCounter, error: counterError } = await supabase
      .from('usage_counters')
      .select('*')
      .eq('organization_id', organization_id)
      .eq('period_start', periodStart)
      .single();

    if (counterError && counterError.code !== 'PGRST116') {
      console.error('Error fetching usage counter:', counterError);
      return new Response(JSON.stringify({ error: 'Database error' }), {
        status: 500,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    // Create counter if it doesn't exist
    if (!usageCounter) {
      const { error: insertError } = await supabase
        .from('usage_counters')
        .insert({
          organization_id,
          period_start: periodStart,
          period_end: periodEnd,
          total_calls: 0,
          included_calls_used: 0,
          overage_calls: 0,
          estimated_charges_cents: 0,
        });

      if (insertError) {
        console.error('Error creating usage counter:', insertError);
        return new Response(JSON.stringify({ error: 'Failed to create counter' }), {
          status: 500,
          headers: { ...corsHeaders, 'Content-Type': 'application/json' },
        });
      }

      const result = await supabase
        .from('usage_counters')
        .select('*')
        .eq('organization_id', organization_id)
        .eq('period_start', periodStart)
        .single();

      if (result.error) {
        return new Response(JSON.stringify({ error: 'Failed to get counter' }), {
          status: 500,
          headers: { ...corsHeaders, 'Content-Type': 'application/json' },
        });
      }

      usageCounter = result.data;
    }

    // Check soft cap (120% of included)
    const totalWithNewCount = usageCounter.total_calls + count;
    const softCap = subscription.included_calls * 1.2;

    if (totalWithNewCount > softCap) {
      // Get dunning behavior
      const { data: orgSettings } = await supabase
        .from('org_settings')
        .select('dunning_behavior')
        .eq('organization_id', organization_id)
        .single();

      const dunningBehavior = orgSettings?.dunning_behavior || 'soft_cap';

      if (dunningBehavior === 'block_immediately') {
        console.log('Usage limit exceeded, blocking request');
        return new Response(JSON.stringify({ error: 'Usage limit exceeded' }), {
          status: 429,
          headers: { ...corsHeaders, 'Content-Type': 'application/json' },
        });
      }
    }

    // Increment usage using RPC function
    const { error: incrementError } = await supabase.rpc('increment_usage', {
      p_org: organization_id,
      p_period_start: periodStart,
      p_count: count,
    });

    if (incrementError) {
      console.error('Failed to increment usage:', incrementError);
      return new Response(JSON.stringify({ error: 'Failed to record usage' }), {
        status: 500,
        headers: { ...corsHeaders, 'Content-Type': 'application/json' },
      });
    }

    // Calculate estimated charges
    const newTotal = usageCounter.total_calls + count;
    const newOverage = Math.max(0, newTotal - subscription.included_calls);
    const estimatedCharges = Math.round(newOverage * Number(subscription.overage_rate_per_call) * 100);

    await supabase
      .from('usage_counters')
      .update({ estimated_charges_cents: estimatedCharges })
      .eq('organization_id', organization_id)
      .eq('period_start', periodStart);

    // Check if we need to send alerts
    const percentUsed = (newTotal / subscription.included_calls) * 100;
    await checkUsageAlerts(organization_id, subscription, percentUsed, supabase);

    console.log('Usage recorded successfully:', { organization_id, newTotal, percentUsed });

    return new Response(JSON.stringify({ 
      success: true, 
      total_calls: newTotal,
      percent_used: Math.round(percentUsed * 100) / 100,
      overage_calls: newOverage,
      estimated_charges_cents: estimatedCharges,
    }), {
      status: 200,
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  } catch (error: any) {
    console.error('Error in meter-usage:', {
      message: error.message,
      stack: error.stack,
      code: error.code,
    });
    return new Response(JSON.stringify({ error: 'Internal server error' }), {
      status: 500,
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  }
});

async function checkUsageAlerts(
  organizationId: string,
  subscription: any,
  percentUsed: number,
  supabase: any
) {
  const { data: quotaPolicy } = await supabase
    .from('quota_policy')
    .select('*')
    .eq('organization_id', organizationId)
    .single();

  const alert70 = quotaPolicy?.alert_70 !== false;
  const alert90 = quotaPolicy?.alert_90 !== false;
  const alert100 = quotaPolicy?.alert_100 !== false;

  const lastAlertType = quotaPolicy?.last_alert_type;
  const lastAlertSent = quotaPolicy?.last_alert_sent_at 
    ? new Date(quotaPolicy.last_alert_sent_at).getTime()
    : 0;
  
  const now = Date.now();
  const oneHour = 3600000; // Don't spam alerts - at most once per hour

  let alertType: string | null = null;

  // Send alerts at configured thresholds (only send each once per period)
  if (alert70 && percentUsed >= 70 && percentUsed < 90 && lastAlertType !== 'usage_70' && (now - lastAlertSent > oneHour)) {
    alertType = 'usage_70';
  } else if (alert90 && percentUsed >= 90 && percentUsed < 100 && lastAlertType !== 'usage_90' && (now - lastAlertSent > oneHour)) {
    alertType = 'usage_90';
  } else if (alert100 && percentUsed >= 100 && lastAlertType !== 'usage_100' && (now - lastAlertSent > oneHour)) {
    alertType = 'usage_100';
  }

  if (alertType) {
    console.log(`Sending ${alertType} alert for org ${organizationId}`);

    const { data: usage } = await supabase
      .from('v_current_period_usage')
      .select('*')
      .eq('organization_id', organizationId)
      .single();

    // Call send-alert-email function
    try {
      await supabase.functions.invoke('send-alert-email', {
        body: {
          organization_id: organizationId,
          alert_type: alertType,
          data: usage,
        },
      });

      // Update last alert sent
      await supabase
        .from('quota_policy')
        .upsert({
          organization_id: organizationId,
          last_alert_type: alertType,
          last_alert_sent_at: new Date().toISOString(),
        });
    } catch (emailError: any) {
      console.error('Failed to send alert email:', emailError);
    }
  }
}
