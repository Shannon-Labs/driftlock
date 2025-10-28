import { serve } from "https://deno.land/std@0.224.0/http/server.ts";
import { Resend } from "https://esm.sh/resend@2.0.0";
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2.45.0';
import { z } from 'https://esm.sh/zod@3.22.4';

const resend = new Resend(Deno.env.get("RESEND_API_KEY"));

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
};

const AlertEmailSchema = z.object({
  organization_id: z.string().uuid({ message: 'Invalid organization ID' }),
  alert_type: z.enum(['usage_70', 'usage_90', 'usage_100', 'payment_failed', 'invoice_threshold'], {
    errorMap: () => ({ message: 'Invalid alert type' }),
  }),
  data: z.record(z.any()),
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
    const validationResult = AlertEmailSchema.safeParse(body);
    
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

    const { organization_id, alert_type, data } = validationResult.data;

    console.log('Sending alert email:', { organization_id, alert_type });

    // Get organization and billing customer info
    const { data: org } = await supabase
      .from('organizations')
      .select('name')
      .eq('id', organization_id)
      .single();

    const { data: billing } = await supabase
      .from('billing_customers')
      .select('billing_email')
      .eq('organization_id', organization_id)
      .single();

    const recipientEmail = billing?.billing_email || 'support@driftlock.io';

    let subject = '';
    let html = '';

    switch (alert_type) {
      case 'usage_70':
        subject = `Driftlock Usage Alert: 70% Quota Used`;
        html = `
          <h2>Usage Alert for ${org?.name || 'Your Organization'}</h2>
          <p>You've used <strong>70%</strong> of your monthly anomaly detection quota.</p>
          <p><strong>Current usage:</strong> ${data.total_calls?.toLocaleString()} / ${data.included_calls?.toLocaleString()} calls</p>
          <p><strong>Days remaining:</strong> ${Math.round(data.days_remaining)} days</p>
          <p>Consider upgrading your plan if you need more capacity.</p>
          <p><a href="https://driftlock.io/dashboard/billing">Manage Billing</a></p>
        `;
        break;

      case 'usage_90':
        subject = `Driftlock Usage Alert: 90% Quota Used`;
        html = `
          <h2>⚠️ High Usage Alert for ${org?.name || 'Your Organization'}</h2>
          <p>You've used <strong>90%</strong> of your monthly anomaly detection quota.</p>
          <p><strong>Current usage:</strong> ${data.total_calls?.toLocaleString()} / ${data.included_calls?.toLocaleString()} calls</p>
          <p><strong>Days remaining:</strong> ${Math.round(data.days_remaining)} days</p>
          <p><strong>Estimated overage:</strong> $${((data.estimated_charges_cents || 0) / 100).toFixed(2)}</p>
          <p>You may incur overage charges. Consider upgrading to avoid interruptions.</p>
          <p><a href="https://driftlock.io/dashboard/billing">Upgrade Plan</a></p>
        `;
        break;

      case 'usage_100':
        subject = `Driftlock Usage Alert: Quota Exceeded`;
        html = `
          <h2>🚨 Quota Exceeded for ${org?.name || 'Your Organization'}</h2>
          <p>You've exceeded your monthly anomaly detection quota.</p>
          <p><strong>Current usage:</strong> ${data.total_calls?.toLocaleString()} / ${data.included_calls?.toLocaleString()} calls</p>
          <p><strong>Overage calls:</strong> ${data.overage_calls?.toLocaleString()}</p>
          <p><strong>Estimated overage charges:</strong> $${((data.estimated_charges_cents || 0) / 100).toFixed(2)}</p>
          <p>Overage charges apply at your plan rate. To avoid additional costs, consider upgrading.</p>
          <p><a href="https://driftlock.io/dashboard/billing">Manage Plan</a></p>
        `;
        break;

      case 'payment_failed':
        subject = `Driftlock Payment Failed`;
        html = `
          <h2>Payment Failed for ${org?.name || 'Your Organization'}</h2>
          <p>We were unable to process your payment for the current billing period.</p>
          <p><strong>Amount due:</strong> $${((data.amount_due_cents || 0) / 100).toFixed(2)}</p>
          <p>Please update your payment method to avoid service interruption.</p>
          <p><a href="${data.hosted_invoice_url}">View Invoice</a></p>
          <p><a href="https://driftlock.io/dashboard/billing">Update Payment Method</a></p>
        `;
        break;

      case 'invoice_threshold':
        subject = `Driftlock Mid-Cycle Invoice`;
        html = `
          <h2>Mid-Cycle Billing for ${org?.name || 'Your Organization'}</h2>
          <p>Your overage usage has exceeded $${(data.threshold_cents / 100).toFixed(2)}.</p>
          <p><strong>Current overage:</strong> $${((data.estimated_charges_cents || 0) / 100).toFixed(2)}</p>
          <p>We'll issue an invoice now to keep your billing manageable.</p>
          <p><a href="https://driftlock.io/dashboard/billing">View Billing</a></p>
        `;
        break;

      default:
        console.error('Unknown alert type:', alert_type);
        return new Response(JSON.stringify({ error: 'Unknown alert type' }), {
          status: 400,
          headers: { ...corsHeaders, 'Content-Type': 'application/json' },
        });
    }

    const emailResponse = await resend.emails.send({
      from: 'Driftlock <alerts@driftlock.io>',
      to: [recipientEmail],
      subject,
      html,
    });

    console.log('Email sent successfully:', emailResponse);

    return new Response(JSON.stringify({ success: true, data: emailResponse }), {
      status: 200,
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  } catch (error: any) {
    console.error('Error sending alert email:', {
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
