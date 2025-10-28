import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Code2, Cloud, Database, Activity, Box, Zap } from "lucide-react";

const Integrations = () => {
  const integrations = [
    {
      name: "OpenTelemetry",
      category: "Observability",
      icon: Activity,
      description: "Native OTLP support for logs, metrics, and traces",
      status: "Available",
      color: "text-primary",
      bgColor: "bg-primary/10"
    },
    {
      name: "Kubernetes",
      category: "Infrastructure",
      icon: Box,
      description: "Helm charts and native K8s resource monitoring",
      status: "Available",
      color: "text-secondary",
      bgColor: "bg-secondary/10"
    },
    {
      name: "PostgreSQL",
      category: "Database",
      icon: Database,
      description: "Time-series optimized storage backend",
      status: "Available",
      color: "text-accent",
      bgColor: "bg-accent/10"
    },
    {
      name: "Prometheus",
      category: "Monitoring",
      icon: Activity,
      description: "Metrics export and alerting integration",
      status: "Available",
      color: "text-primary",
      bgColor: "bg-primary/10"
    },
    {
      name: "Grafana",
      category: "Visualization",
      icon: Activity,
      description: "Pre-built dashboards and panels",
      status: "Available",
      color: "text-secondary",
      bgColor: "bg-secondary/10"
    },
    {
      name: "AWS",
      category: "Cloud Provider",
      icon: Cloud,
      description: "CloudWatch, ECS, EKS integration",
      status: "Available",
      color: "text-accent",
      bgColor: "bg-accent/10"
    },
    {
      name: "Google Cloud",
      category: "Cloud Provider",
      icon: Cloud,
      description: "Cloud Logging, GKE support",
      status: "Available",
      color: "text-primary",
      bgColor: "bg-primary/10"
    },
    {
      name: "Azure",
      category: "Cloud Provider",
      icon: Cloud,
      description: "Azure Monitor and AKS integration",
      status: "Available",
      color: "text-secondary",
      bgColor: "bg-secondary/10"
    },
    {
      name: "Datadog",
      category: "Observability",
      icon: Activity,
      description: "Log forwarding and custom metrics",
      status: "Coming Soon",
      color: "text-accent",
      bgColor: "bg-accent/10"
    },
    {
      name: "Splunk",
      category: "SIEM",
      icon: Activity,
      description: "Security event correlation",
      status: "Coming Soon",
      color: "text-primary",
      bgColor: "bg-primary/10"
    },
    {
      name: "Elasticsearch",
      category: "Search & Analytics",
      icon: Database,
      description: "Log aggregation and search",
      status: "Coming Soon",
      color: "text-secondary",
      bgColor: "bg-secondary/10"
    },
    {
      name: "Kafka",
      category: "Streaming",
      icon: Zap,
      description: "Event streaming and processing",
      status: "Coming Soon",
      color: "text-accent",
      bgColor: "bg-accent/10"
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-6xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              <span className="text-gradient">Integrations</span>
            </h1>
            <p className="text-xl text-muted-foreground">
              Connect Driftlock with your existing observability stack
            </p>
          </div>

          {/* Integration Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-16">
            {integrations.map((integration, idx) => (
              <Card key={idx} className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                <div className="flex items-start justify-between mb-4">
                  <div className={`w-12 h-12 rounded-lg ${integration.bgColor} flex items-center justify-center`}>
                    <integration.icon className={`w-6 h-6 ${integration.color}`} />
                  </div>
                  <Badge variant="secondary" className={
                    integration.status === "Available" 
                      ? "bg-green-500/20 text-green-500 border-green-500/30"
                      : "bg-yellow-500/20 text-yellow-500 border-yellow-500/30"
                  }>
                    {integration.status}
                  </Badge>
                </div>
                <h3 className="text-lg font-bold mb-1">{integration.name}</h3>
                <p className="text-sm text-muted-foreground mb-3">{integration.category}</p>
                <p className="text-sm text-muted-foreground">{integration.description}</p>
              </Card>
            ))}
          </div>

          {/* Custom Integration */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 text-center">
            <Code2 className="w-12 h-12 text-primary mx-auto mb-4" />
            <h2 className="text-2xl font-bold mb-4">Need a Custom Integration?</h2>
            <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
              Our REST and GraphQL APIs make it easy to integrate Driftlock with any system. 
              Enterprise customers get dedicated integration support.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a href="/docs">
                <button className="px-6 py-3 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                  View API Docs
                </button>
              </a>
              <a href="/contact">
                <button className="px-6 py-3 bg-background border border-border rounded-lg font-medium hover:bg-muted transition-colors">
                  Contact Sales
                </button>
              </a>
            </div>
          </Card>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Integrations;
