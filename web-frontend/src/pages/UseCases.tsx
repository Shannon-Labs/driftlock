import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Building2, Heart, Bot, Shield, TrendingUp, Database } from "lucide-react";

const UseCases = () => {
  const useCases = [
    {
      icon: Building2,
      title: "Financial Services",
      subtitle: "DORA Compliance & Fraud Detection",
      color: "text-primary",
      bgColor: "bg-primary/10",
      description: "Meet DORA (Digital Operational Resilience Act) requirements with explainable anomaly detection for transaction monitoring, trading systems, and operational resilience.",
      features: [
        "Transaction anomaly detection with field-level attribution",
        "Real-time fraud detection with statistical significance",
        "DORA-compliant evidence bundles and audit trails",
        "Deterministic results for regulatory inquiries"
      ],
      metrics: {
        detection: "99.2% detection rate",
        falsePositive: "<0.5% false positives",
        latency: "<50ms p95 latency"
      }
    },
    {
      icon: Heart,
      title: "Healthcare",
      subtitle: "HIPAA Compliance & Patient Safety",
      color: "text-secondary",
      bgColor: "bg-secondary/10",
      description: "Detect anomalies in EHR access patterns, medication orders, and clinical workflows while maintaining HIPAA compliance and patient privacy.",
      features: [
        "Privacy-preserving anomaly detection with data redaction",
        "Clinical workflow deviation detection",
        "Medication order validation with compression analysis",
        "HIPAA audit logs with cryptographic integrity"
      ],
      metrics: {
        detection: "Early warning systems",
        privacy: "Zero PHI exposure",
        compliance: "HIPAA compliant"
      }
    },
    {
      icon: Bot,
      title: "AI/LLM Systems",
      subtitle: "Runtime AI Act Compliance",
      color: "text-accent",
      bgColor: "bg-accent/10",
      description: "Monitor LLM prompts, responses, and tool calls for compliance with the EU AI Act. Detect prompt injection, data leakage, and model drift.",
      features: [
        "LLM I/O monitoring with compression-based detection",
        "Prompt injection and jailbreak detection",
        "Model drift detection through output analysis",
        "Runtime AI compliance evidence generation"
      ],
      metrics: {
        coverage: "100% LLM I/O tracked",
        detection: "Real-time alerts",
        compliance: "AI Act ready"
      }
    },
    {
      icon: Shield,
      title: "Cybersecurity",
      subtitle: "NIS2 Compliance & Threat Detection",
      color: "text-primary",
      bgColor: "bg-primary/10",
      description: "Meet NIS2 (Network and Information Security) requirements with explainable threat detection, incident reporting, and security monitoring.",
      features: [
        "Network traffic anomaly detection",
        "Authentication pattern analysis",
        "NIS2-compliant incident reporting",
        "Security event correlation with CBAD"
      ],
      metrics: {
        detection: "Sub-second detection",
        reporting: "Automated reports",
        compliance: "NIS2 compliant"
      }
    },
    {
      icon: TrendingUp,
      title: "SaaS & E-commerce",
      subtitle: "Performance & Revenue Protection",
      color: "text-secondary",
      bgColor: "bg-secondary/10",
      description: "Detect performance degradation, pricing anomalies, and revenue leakage in production systems before they impact customers.",
      features: [
        "API performance anomaly detection",
        "Pricing and billing validation",
        "User behavior analysis for fraud prevention",
        "Revenue impact quantification"
      ],
      metrics: {
        uptime: "99.99% SLA",
        detection: "Proactive alerts",
        impact: "Revenue protection"
      }
    },
    {
      icon: Database,
      title: "Data Engineering",
      subtitle: "Pipeline Monitoring & Data Quality",
      color: "text-accent",
      bgColor: "bg-accent/10",
      description: "Monitor data pipelines for schema drift, data quality issues, and processing anomalies with format-aware compression analysis.",
      features: [
        "Schema drift detection with structural analysis",
        "Data quality validation through compression",
        "Pipeline performance monitoring",
        "Automated data quality reports"
      ],
      metrics: {
        coverage: "End-to-end pipelines",
        detection: "Real-time validation",
        quality: "Automated checks"
      }
    }
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-7xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              <span className="text-gradient">Use Cases</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              Driftlock powers explainable anomaly detection across regulated industries
            </p>
          </div>

          {/* Use Cases Grid */}
          <div className="space-y-8">
            {useCases.map((useCase, idx) => (
              <Card key={idx} className="p-8 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                  {/* Header */}
                  <div>
                    <div className={`w-16 h-16 rounded-xl ${useCase.bgColor} flex items-center justify-center mb-4`}>
                      <useCase.icon className={`w-8 h-8 ${useCase.color}`} />
                    </div>
                    <h2 className="text-2xl font-bold mb-2">{useCase.title}</h2>
                    <p className="text-sm text-muted-foreground mb-4">{useCase.subtitle}</p>
                    <div className="space-y-2">
                      {Object.entries(useCase.metrics).map(([key, value]) => (
                        <Badge key={key} variant="secondary" className="mr-2">
                          {value}
                        </Badge>
                      ))}
                    </div>
                  </div>

                  {/* Description & Features */}
                  <div className="lg:col-span-2">
                    <p className="text-muted-foreground mb-6 leading-relaxed">
                      {useCase.description}
                    </p>

                    <h3 className="font-semibold mb-3">Key Capabilities</h3>
                    <ul className="space-y-2">
                      {useCase.features.map((feature, featureIdx) => (
                        <li key={featureIdx} className="flex items-start gap-2 text-sm text-muted-foreground">
                          <div className={`w-1.5 h-1.5 rounded-full ${useCase.bgColor} mt-2 flex-shrink-0`}></div>
                          {feature}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </Card>
            ))}
          </div>

          {/* CTA */}
          <Card className="mt-16 p-12 bg-gradient-card border-primary/10 text-center">
            <h2 className="text-3xl font-bold mb-4">Ready to get started?</h2>
            <p className="text-muted-foreground mb-8 max-w-2xl mx-auto">
              See how Driftlock can provide explainable anomaly detection for your specific use case
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a href="/contact">
                <button className="px-6 py-3 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                  Schedule a Demo
                </button>
              </a>
              <a href="/docs">
                <button className="px-6 py-3 bg-background border border-border rounded-lg font-medium hover:bg-muted transition-colors">
                  Read Documentation
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

export default UseCases;
