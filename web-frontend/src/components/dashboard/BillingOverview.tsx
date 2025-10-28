import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CreditCard, Download, ExternalLink, TrendingUp } from "lucide-react";
import { supabase } from "@/integrations/supabase/client";

interface Subscription {
  plan: string;
  status: string;
  included_calls: number;
  overage_rate_per_call: number;
  current_period_end: string;
  stripe_subscription_id: string;
}

interface Invoice {
  id: string;
  stripe_invoice_id: string;
  status: string;
  amount_due_cents: number;
  amount_paid_cents: number;
  hosted_invoice_url: string;
  invoice_pdf_url: string;
  created_at: string;
  paid_at: string;
}

export const BillingOverview = ({ organizationId }: { organizationId: string }) => {
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBillingData();
  }, [organizationId]);

  const fetchBillingData = async () => {
    try {
      // Fetch subscription
      const { data: subData, error: subError } = await supabase
        .from('subscriptions')
        .select('*')
        .eq('organization_id', organizationId)
        .single();

      if (subError && subError.code !== 'PGRST116') throw subError;
      setSubscription(subData);

      // Fetch invoices
      const { data: invoiceData, error: invoiceError } = await supabase
        .from('invoices_mirror')
        .select('*')
        .eq('organization_id', organizationId)
        .order('created_at', { ascending: false })
        .limit(10);

      if (invoiceError) throw invoiceError;
      setInvoices(invoiceData || []);
    } catch (error) {
      console.error('Error fetching billing data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpgrade = () => {
    // Navigate to pricing/checkout
    window.location.href = '/pricing';
  };

  const handleManageBilling = () => {
    // Open Stripe Customer Portal
    // In production, call an edge function to create a portal session
    console.log('Open customer portal');
  };

  if (loading) {
    return (
      <Card className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-4 bg-muted rounded w-1/3"></div>
          <div className="h-20 bg-muted rounded"></div>
        </div>
      </Card>
    );
  }

  const planNames: Record<string, string> = {
    developer: 'Developer',
    standard: 'Standard',
    growth: 'Growth',
    enterprise: 'Enterprise',
  };

  const planName = subscription ? planNames[subscription.plan] || subscription.plan : 'No Plan';

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Billing & Invoices</h2>

      {/* Current Plan Card */}
      <Card className="p-6 space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold">{planName} Plan</h3>
            <p className="text-sm text-muted-foreground">
              {subscription?.included_calls.toLocaleString()} included calls/month
            </p>
          </div>
          <Badge variant={subscription?.status === 'active' ? 'default' : 'secondary'}>
            {subscription?.status || 'Inactive'}
          </Badge>
        </div>

        {subscription && (
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Overage Rate</span>
              <span className="font-medium">
                ${(subscription.overage_rate_per_call * 100).toFixed(4)} per call
              </span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Period Ends</span>
              <span className="font-medium">
                {new Date(subscription.current_period_end).toLocaleDateString()}
              </span>
            </div>
          </div>
        )}

        <div className="flex gap-3 pt-4 border-t">
          <Button onClick={handleUpgrade} className="flex-1">
            <TrendingUp className="w-4 h-4 mr-2" />
            Upgrade Plan
          </Button>
          <Button onClick={handleManageBilling} variant="outline" className="flex-1">
            <CreditCard className="w-4 h-4 mr-2" />
            Manage Billing
          </Button>
        </div>
      </Card>

      {/* Invoices List */}
      <Card className="p-6 space-y-4">
        <h3 className="text-lg font-semibold">Recent Invoices</h3>
        
        {invoices.length === 0 ? (
          <p className="text-sm text-muted-foreground py-4 text-center">
            No invoices yet
          </p>
        ) : (
          <div className="space-y-3">
            {invoices.map((invoice) => (
              <div
                key={invoice.id}
                className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium">
                      ${(invoice.amount_due_cents / 100).toFixed(2)}
                    </p>
                    <Badge
                      variant={invoice.status === 'paid' ? 'default' : 'secondary'}
                      className="text-xs"
                    >
                      {invoice.status}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {new Date(invoice.created_at).toLocaleDateString()} • {invoice.stripe_invoice_id}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  {invoice.hosted_invoice_url && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => window.open(invoice.hosted_invoice_url, '_blank')}
                    >
                      <ExternalLink className="w-4 h-4" />
                    </Button>
                  )}
                  {invoice.invoice_pdf_url && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => window.open(invoice.invoice_pdf_url, '_blank')}
                    >
                      <Download className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
};
