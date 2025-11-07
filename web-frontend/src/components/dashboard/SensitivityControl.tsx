import { useState, useEffect } from "react";
import { Card } from "@/components/ui/card";
import { Slider } from "@/components/ui/slider";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Info } from "lucide-react";
import { apiClient } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";

export const SensitivityControl = ({ organizationId }: { organizationId: string }) => {
  const [sensitivity, setSensitivity] = useState(0.5);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    fetchSensitivity();
  }, [organizationId]);

  const fetchSensitivity = async () => {
    try {
      // Fetch current config from API
      const config = await apiClient.get<{ ncd_threshold: number; p_value_threshold: number }>('/v1/config');
      // Use p_value_threshold as sensitivity indicator (inverted: lower = more sensitive)
      setSensitivity(1 - (config.p_value_threshold || 0.05));
    } catch (error) {
      console.error('Error fetching sensitivity:', error);
      // Default to medium sensitivity if API call fails
      setSensitivity(0.5);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      // Update config via API (p_value_threshold = 1 - sensitivity)
      const pValueThreshold = 1 - sensitivity;
      await apiClient.patch('/v1/config', {
        p_value_threshold: pValueThreshold,
      });

      toast({
        title: "Settings saved",
        description: "Anomaly sensitivity updated successfully",
      });
    } catch (error) {
      console.error('Error saving sensitivity:', error);
      toast({
        title: "Error",
        description: "Failed to save settings",
        variant: "destructive",
      });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Card className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-4 bg-muted rounded w-1/2"></div>
          <div className="h-12 bg-muted rounded"></div>
        </div>
      </Card>
    );
  }

  const sensitivityLabel = 
    sensitivity < 0.3 ? 'Low' :
    sensitivity < 0.7 ? 'Medium' :
    'High';

  const sensitivityColor = 
    sensitivity < 0.3 ? 'text-green-500' :
    sensitivity < 0.7 ? 'text-yellow-500' :
    'text-orange-500';

  return (
    <Card className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">Anomaly Sensitivity</h3>
          <p className="text-sm text-muted-foreground">
            Control detection frequency and costs
          </p>
        </div>
        <Badge variant="outline" className={sensitivityColor}>
          {sensitivityLabel} Sensitivity
        </Badge>
      </div>

      <div className="space-y-4">
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground w-12">Low</span>
          <Slider
            value={[sensitivity]}
            onValueChange={(value) => setSensitivity(value[0])}
            min={0}
            max={1}
            step={0.1}
            className="flex-1"
          />
          <span className="text-sm text-muted-foreground w-12 text-right">High</span>
        </div>

        <div className="text-center">
          <p className="text-3xl font-bold text-gradient">{(sensitivity * 100).toFixed(0)}%</p>
        </div>

        <div className="flex items-start gap-2 p-3 bg-muted/50 rounded-lg">
          <Info className="w-5 h-5 text-muted-foreground flex-shrink-0 mt-0.5" />
          <p className="text-xs text-muted-foreground">
            <strong>Lower sensitivity</strong> = fewer detections = lower costs. 
            <strong> Higher sensitivity</strong> = more detections = better coverage but higher costs.
            Only anomaly detections consume your quota.
          </p>
        </div>

        <Button 
          onClick={handleSave} 
          className="w-full" 
          disabled={saving}
        >
          {saving ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Saving...
            </>
          ) : (
            'Save Settings'
          )}
        </Button>
      </div>
    </Card>
  );
};
