import { useState, useEffect } from "react";
import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { UsageOverview } from "@/components/dashboard/UsageOverview";
import { SensitivityControl } from "@/components/dashboard/SensitivityControl";
import { BillingOverview } from "@/components/dashboard/BillingOverview";
import { ApiKeyManagement } from "@/components/dashboard/ApiKeyManagement";
import { supabase } from "@/integrations/supabase/client";
import { useAuth } from "@/contexts/AuthContext";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { 
  Activity, 
  AlertTriangle, 
  TrendingUp, 
  Database,
  Clock,
  Zap,
  CheckCircle2,
  DollarSign,
  Settings,
  BarChart3,
  Key
} from "lucide-react";

interface Organization {
  id: string;
  name: string;
  slug: string;
}

const Dashboard = () => {
  const { user } = useAuth();
  const [organization, setOrganization] = useState<Organization | null>(null);
  const [anomalies, setAnomalies] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, [user]);

  const loadDashboardData = async () => {
    try {
      // Load user's primary organization
      const { data: orgData, error: orgError } = await supabase
        .from('organization_members')
        .select('organization_id, organizations(id, name, slug)')
        .eq('user_id', user?.id)
        .single();

      if (orgError) throw orgError;
      
      if (orgData?.organizations) {
        const org = Array.isArray(orgData.organizations) 
          ? orgData.organizations[0] 
          : orgData.organizations;
        setOrganization(org as Organization);

        // Load anomalies for this organization
        const { data: anomalyData, error: anomalyError } = await supabase
          .from('anomaly_events')
          .select('*')
          .eq('organization_id', org.id)
          .order('detection_timestamp', { ascending: false })
          .limit(10);

        if (!anomalyError && anomalyData) {
          setAnomalies(anomalyData);
        }
      }
    } catch (error: any) {
      console.error('Error loading dashboard:', error);
      toast.error('Failed to load dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const stats = [
    { label: "Events Processed", value: "0", change: "—", icon: Activity, color: "text-primary" },
    { label: "Anomalies Detected", value: anomalies.length.toString(), change: "—", icon: AlertTriangle, color: "text-secondary" },
    { label: "Active Streams", value: "0", change: "—", icon: Database, color: "text-primary" },
    { label: "API Keys", value: "1", change: "—", icon: Key, color: "text-accent" },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-7xl">
          {/* Header */}
          <div className="mb-8">
            <h1 className="text-4xl font-bold mb-2">
              <span className="text-gradient">Dashboard</span>
            </h1>
            <p className="text-muted-foreground">
              {organization ? `${organization.name} - Real-time monitoring` : 'Loading...'}
            </p>
          </div>

          {loading ? (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {[1, 2, 3, 4].map((i) => (
                  <Skeleton key={i} className="h-32" />
                ))}
              </div>
            </div>
          ) : (
            <>
              {/* Stats Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {stats.map((stat, idx) => (
                  <Card key={idx} className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">{stat.label}</p>
                        <p className="text-3xl font-bold mb-1">{stat.value}</p>
                        <p className="text-sm text-muted-foreground">{stat.change}</p>
                      </div>
                      <div className={`w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center ${stat.color}`}>
                        <stat.icon className="w-6 h-6" />
                      </div>
                    </div>
                  </Card>
                ))}
              </div>

              {/* Main Content */}
              <Tabs defaultValue="keys" className="space-y-6">
                <TabsList className="grid w-full grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-2">
                  <TabsTrigger value="keys">API Keys</TabsTrigger>
                  <TabsTrigger value="anomalies">Anomalies</TabsTrigger>
                  <TabsTrigger value="usage">Usage</TabsTrigger>
                  <TabsTrigger value="billing">Billing</TabsTrigger>
                  <TabsTrigger value="settings">Settings</TabsTrigger>
                </TabsList>

                {/* API Keys Tab */}
                <TabsContent value="keys" className="space-y-6">
                  <ApiKeyManagement />
                </TabsContent>

                {/* Anomalies Tab */}
                <TabsContent value="anomalies" className="space-y-6">
                  <Card className="p-6 bg-gradient-card border-primary/10">
                    <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
                      <AlertTriangle className="w-5 h-5 text-secondary" />
                      Recent Anomalies
                    </h3>
                    {anomalies.length === 0 ? (
                      <div className="text-center py-12">
                        <CheckCircle2 className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
                        <h4 className="font-semibold mb-2">No anomalies detected</h4>
                        <p className="text-sm text-muted-foreground">
                          Start sending events to detect anomalies in your data streams
                        </p>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {anomalies.map((anomaly) => (
                          <Card key={anomaly.id} className="p-4 bg-background/50 border-primary/10 hover:border-primary/30 transition-all">
                            <div className="flex items-start justify-between">
                              <div className="flex-1">
                                <div className="flex items-center gap-2 mb-2">
                                  <Badge 
                                    variant={
                                      anomaly.severity === "critical" || anomaly.severity === "high"
                                        ? "destructive" 
                                        : anomaly.severity === "medium" 
                                        ? "default" 
                                        : "secondary"
                                    }
                                  >
                                    {anomaly.severity}
                                  </Badge>
                                  <span className="text-sm text-muted-foreground flex items-center gap-1">
                                    <Clock className="w-3 h-3" />
                                    {new Date(anomaly.detection_timestamp).toLocaleString()}
                                  </span>
                                </div>
                                <p className="text-sm mb-2">{anomaly.description}</p>
                                <p className="text-xs text-muted-foreground">{anomaly.explanation}</p>
                              </div>
                              <Zap className="w-5 h-5 text-secondary" />
                            </div>
                          </Card>
                        ))}
                      </div>
                    )}
                  </Card>
                </TabsContent>

                {/* Usage Tab */}
                <TabsContent value="usage" className="space-y-6">
                  {organization && <UsageOverview organizationId={organization.id} />}
                </TabsContent>

                {/* Billing Tab */}
                <TabsContent value="billing" className="space-y-6">
                  {organization && <BillingOverview organizationId={organization.id} />}
                </TabsContent>

                {/* Settings Tab */}
                <TabsContent value="settings" className="space-y-6">
                  {organization && <SensitivityControl organizationId={organization.id} />}
                </TabsContent>
              </Tabs>
            </>
          )}
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Dashboard;
