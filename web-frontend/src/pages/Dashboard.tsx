import { useState, useEffect } from "react";
import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { SensitivityControl } from "@/components/dashboard/SensitivityControl";
import { useAuth } from "@/contexts/AuthContext";
import { apiClient } from "@/lib/api";
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

interface Anomaly {
  id: string;
  timestamp: string;
  stream_type: string;
  ncd_score: number;
  p_value: number;
  status: string;
  glass_box_explanation: string;
  severity?: string;
}

const Dashboard = () => {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const [anomalies, setAnomalies] = useState<Anomaly[]>([]);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalAnomalies: 0,
    eventsProcessed: 0,
    activeStreams: 0,
    apiKeys: 1,
  });

  useEffect(() => {
    if (isAuthenticated && !authLoading) {
      loadDashboardData();
    }
  }, [isAuthenticated, authLoading]);

  const loadDashboardData = async () => {
    try {
      // Load anomalies from API
      const response = await apiClient.get<{ anomalies: Anomaly[]; total: number }>('/v1/anomalies?limit=10');
      setAnomalies(response.anomalies || []);
      setStats(prev => ({ ...prev, totalAnomalies: response.total || 0 }));
    } catch (error: any) {
      console.error('Error loading dashboard:', error);
      toast.error('Failed to load dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const displayStats = [
    { label: "Events Processed", value: stats.eventsProcessed.toString(), change: "—", icon: Activity, color: "text-primary" },
    { label: "Anomalies Detected", value: stats.totalAnomalies.toString(), change: "—", icon: AlertTriangle, color: "text-secondary" },
    { label: "Active Streams", value: stats.activeStreams.toString(), change: "—", icon: Database, color: "text-primary" },
    { label: "API Keys", value: stats.apiKeys.toString(), change: "—", icon: Key, color: "text-accent" },
  ];

  if (authLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Navigation />
        <main className="pt-24 pb-16 px-4">
          <div className="container mx-auto max-w-7xl">
            <Skeleton className="h-32 w-full" />
          </div>
        </main>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-background">
        <Navigation />
        <main className="pt-24 pb-16 px-4">
          <div className="container mx-auto max-w-7xl text-center">
            <Card className="p-8">
              <h2 className="text-2xl font-bold mb-4">Authentication Required</h2>
              <p className="text-muted-foreground mb-4">
                Please sign in with your API key to access the dashboard.
              </p>
            </Card>
          </div>
        </main>
      </div>
    );
  }

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
              Real-time anomaly monitoring
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
                {displayStats.map((stat, idx) => (
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
              <Tabs defaultValue="anomalies" className="space-y-6">
                <TabsList className="grid w-full grid-cols-2 md:grid-cols-3 gap-2">
                  <TabsTrigger value="anomalies">Anomalies</TabsTrigger>
                  <TabsTrigger value="settings">Settings</TabsTrigger>
                </TabsList>

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
                        {anomalies.map((anomaly) => {
                          const severity = anomaly.severity || (anomaly.p_value < 0.001 ? 'critical' : anomaly.p_value < 0.01 ? 'high' : anomaly.p_value < 0.05 ? 'medium' : 'low');
                          return (
                            <Card key={anomaly.id} className="p-4 bg-background/50 border-primary/10 hover:border-primary/30 transition-all">
                              <div className="flex items-start justify-between">
                                <div className="flex-1">
                                  <div className="flex items-center gap-2 mb-2">
                                    <Badge 
                                      variant={
                                        severity === "critical" || severity === "high"
                                          ? "destructive" 
                                          : severity === "medium" 
                                          ? "default" 
                                          : "secondary"
                                      }
                                    >
                                      {severity}
                                    </Badge>
                                    <span className="text-sm text-muted-foreground flex items-center gap-1">
                                      <Clock className="w-3 h-3" />
                                      {new Date(anomaly.timestamp).toLocaleString()}
                                    </span>
                                  </div>
                                  <p className="text-sm mb-2">{anomaly.glass_box_explanation}</p>
                                  <p className="text-xs text-muted-foreground">
                                    Stream: {anomaly.stream_type} | NCD: {anomaly.ncd_score.toFixed(3)} | p-value: {anomaly.p_value.toFixed(4)}
                                  </p>
                                </div>
                                <Zap className="w-5 h-5 text-secondary" />
                              </div>
                            </Card>
                          );
                        })}
                      </div>
                    )}
                  </Card>
                </TabsContent>

                {/* Settings Tab */}
                <TabsContent value="settings" className="space-y-6">
                  <SensitivityControl organizationId="default" />
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
