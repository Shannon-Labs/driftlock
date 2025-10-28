import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { AlertCircle, TrendingUp, Calendar } from "lucide-react";
import { supabase } from "@/integrations/supabase/client";

interface UsageData {
  total_calls: number;
  included_calls: number;
  percent_used: number;
  overage_calls: number;
  estimated_overage_usd: number;
  days_remaining: number;
  dunning_state: string;
}

export const UsageOverview = ({ organizationId }: { organizationId: string }) => {
  const [usage, setUsage] = useState<UsageData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUsage();

    // Refresh every 30 seconds
    const interval = setInterval(fetchUsage, 30000);
    return () => clearInterval(interval);
  }, [organizationId]);

  const fetchUsage = async () => {
    try {
      const { data, error } = await supabase
        .from('v_current_period_usage')
        .select('*')
        .eq('organization_id', organizationId)
        .single();

      if (error) throw error;
      setUsage(data);
    } catch (error) {
      console.error('Error fetching usage:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-4 bg-muted rounded w-1/3"></div>
          <div className="h-8 bg-muted rounded"></div>
        </div>
      </Card>
    );
  }

  if (!usage) {
    return (
      <Card className="p-6">
        <p className="text-muted-foreground">No usage data available</p>
      </Card>
    );
  }

  const percentColor = 
    usage.percent_used >= 100 ? 'text-destructive' :
    usage.percent_used >= 90 ? 'text-orange-500' :
    usage.percent_used >= 70 ? 'text-yellow-500' :
    'text-primary';

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Usage Overview</h2>
          <p className="text-sm text-muted-foreground">Current billing period</p>
        </div>
        <Badge variant="secondary" className="flex items-center gap-2">
          <Calendar className="w-4 h-4" />
          {Math.round(usage.days_remaining)} days remaining
        </Badge>
      </div>

      {/* Main Usage Card */}
      <Card className="p-6 space-y-6">
        {/* Anomaly Detection Badge */}
        <Badge variant="outline" className="w-fit">
          💡 Only anomaly detections are billable
        </Badge>

        {/* Usage Progress */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">Quota Used</span>
            <span className={`text-2xl font-bold ${percentColor}`}>
              {usage.percent_used.toFixed(1)}%
            </span>
          </div>
          <Progress value={Math.min(usage.percent_used, 100)} className="h-3" />
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <span>{usage.total_calls.toLocaleString()} calls</span>
            <span>{usage.included_calls.toLocaleString()} included</span>
          </div>
        </div>

        {/* Overage Info */}
        {usage.overage_calls > 0 && (
          <div className="border-t pt-4 space-y-2">
            <div className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-orange-500" />
              <span className="font-semibold">Overage Usage</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">
                {usage.overage_calls.toLocaleString()} additional calls
              </span>
              <span className="text-lg font-bold">
                ${usage.estimated_overage_usd.toFixed(2)}
              </span>
            </div>
            <p className="text-xs text-muted-foreground">
              Estimated overage charges for this period
            </p>
          </div>
        )}

        {/* Warnings */}
        {usage.percent_used >= 70 && (
          <div className="border-t pt-4">
            <div className="flex items-start gap-3 p-3 bg-muted/50 rounded-lg">
              <AlertCircle className={`w-5 h-5 flex-shrink-0 ${
                usage.percent_used >= 100 ? 'text-destructive' :
                usage.percent_used >= 90 ? 'text-orange-500' :
                'text-yellow-500'
              }`} />
              <div className="space-y-1">
                <p className="text-sm font-medium">
                  {usage.percent_used >= 100 ? 'Quota Exceeded' :
                   usage.percent_used >= 90 ? 'High Usage Alert' :
                   'Usage Warning'}
                </p>
                <p className="text-xs text-muted-foreground">
                  {usage.percent_used >= 100 
                    ? 'You are being charged overage rates. Consider upgrading your plan.'
                    : usage.percent_used >= 90
                    ? 'You\'re approaching your quota limit. Upgrade to avoid overage charges.'
                    : 'You\'ve used over 70% of your quota this period.'}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Dunning State Warning */}
        {usage.dunning_state && usage.dunning_state !== 'ok' && (
          <div className="border-t pt-4">
            <div className="flex items-start gap-3 p-3 bg-destructive/10 rounded-lg">
              <AlertCircle className="w-5 h-5 flex-shrink-0 text-destructive" />
              <div className="space-y-1">
                <p className="text-sm font-medium text-destructive">Payment Issue</p>
                <p className="text-xs text-muted-foreground">
                  Please update your payment method to avoid service interruption.
                </p>
              </div>
            </div>
          </div>
        )}
      </Card>

      {/* Quick Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <p className="text-sm text-muted-foreground">Total Detections</p>
          <p className="text-2xl font-bold">{usage.total_calls.toLocaleString()}</p>
        </Card>
        <Card className="p-4">
          <p className="text-sm text-muted-foreground">Included Calls</p>
          <p className="text-2xl font-bold">{usage.included_calls.toLocaleString()}</p>
        </Card>
        <Card className="p-4">
          <p className="text-sm text-muted-foreground">Days Left</p>
          <p className="text-2xl font-bold">{Math.round(usage.days_remaining)}</p>
        </Card>
      </div>
    </div>
  );
};
