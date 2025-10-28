import { Card } from "@/components/ui/card";
import { Shield, FileCheck, Bot, Lock } from "lucide-react";

const complianceFeatures = [
  {
    icon: Shield,
    title: "DORA Compliance",
    description: "Digital Operational Resilience Act evidence bundles with cryptographic integrity chains for financial services.",
    badge: "Financial Services",
  },
  {
    icon: FileCheck,
    title: "NIS2 Compliance",
    description: "EU cybersecurity incident reporting templates with automated compliance report generation.",
    badge: "EU Regulation",
  },
  {
    icon: Bot,
    title: "Runtime AI Monitoring",
    description: "AI Act compliance for LLM/ML systems with explainable anomaly detection and governance controls.",
    badge: "AI Systems",
  },
  {
    icon: Lock,
    title: "Cryptographic Audit Trails",
    description: "Tamper-evident evidence packages with digital signatures and automated integrity verification.",
    badge: "Enterprise Security",
  },
];

export const Compliance = () => {
  return (
    <section className="py-24 px-4 relative overflow-hidden">
      {/* Background Gradient */}
      <div className="absolute inset-0 bg-gradient-to-b from-card/20 via-background to-card/20"></div>
      
      <div className="container mx-auto relative z-10">
        <div className="max-w-3xl mx-auto text-center mb-16 space-y-4">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20">
            <Shield className="w-4 h-4 text-primary" />
            <span className="text-sm font-medium">Enterprise Compliance</span>
          </div>
          
          <h2 className="text-4xl md:text-5xl font-bold">
            Built for <span className="text-gradient">Regulated Industries</span>
          </h2>
          <p className="text-xl text-muted-foreground">
            Comprehensive compliance framework for financial services, healthcare, and critical infrastructure
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-5xl mx-auto">
          {complianceFeatures.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <Card 
                key={index}
                className="p-8 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all duration-300 hover:shadow-card group"
              >
                <div className="space-y-4">
                  <div className="flex items-start justify-between">
                    <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                      <Icon className="w-7 h-7 text-primary" />
                    </div>
                    <span className="text-xs font-medium px-3 py-1 rounded-full bg-secondary/20 text-secondary border border-secondary/30">
                      {feature.badge}
                    </span>
                  </div>
                  
                  <div className="space-y-2">
                    <h3 className="text-xl font-semibold">{feature.title}</h3>
                    <p className="text-muted-foreground leading-relaxed">
                      {feature.description}
                    </p>
                  </div>
                </div>
              </Card>
            );
          })}
        </div>

        {/* Additional Info */}
        <div className="mt-16 max-w-4xl mx-auto">
          <Card className="p-8 bg-gradient-card border-primary/20">
            <div className="flex flex-col md:flex-row gap-8 items-start">
              <div className="flex-shrink-0">
                <div className="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center">
                  <Lock className="w-8 h-8 text-primary" />
                </div>
              </div>
              
              <div className="space-y-4 flex-grow">
                <h3 className="text-2xl font-bold">Privacy-First Architecture</h3>
                <p className="text-muted-foreground leading-relaxed">
                  On-premises deployment with configurable data redaction. Your telemetry data never leaves your infrastructure. 
                  Full control over data retention, encryption, and access policies.
                </p>
                
                <div className="flex flex-wrap gap-3 pt-2">
                  <span className="text-sm px-3 py-1 rounded-full bg-primary/10 text-primary border border-primary/20">
                    On-Premises Deployment
                  </span>
                  <span className="text-sm px-3 py-1 rounded-full bg-primary/10 text-primary border border-primary/20">
                    Data Redaction
                  </span>
                  <span className="text-sm px-3 py-1 rounded-full bg-primary/10 text-primary border border-primary/20">
                    End-to-End Encryption
                  </span>
                  <span className="text-sm px-3 py-1 rounded-full bg-primary/10 text-primary border border-primary/20">
                    Zero-Trust Architecture
                  </span>
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </section>
  );
};
