import { ArrowRight, Database, Gauge, FileCheck, Cloud } from "lucide-react";

const steps = [
  {
    number: "01",
    icon: Cloud,
    title: "Stream Ingest",
    description: "Send telemetry via REST API or OpenTelemetry Collector. Logs, metrics, traces, IoT data, LLM I/O — any structured data."
  },
  {
    number: "02",
    icon: Database,
    title: "Compression Engine",
    description: "OpenZL builds a baseline compression plan from your data. Uses Normalized Compression Distance (NCD) and statistical permutation tests for deterministic anomaly detection. No training required — format-aware compression learns your schema automatically."
  },
  {
    number: "03",
    icon: Gauge,
    title: "Anomaly Detection",
    description: "Live data that compresses poorly triggers alerts with field-level explanations: 'message field 5× larger' or 'unexpected schema element'."
  },
  {
    number: "04",
    icon: FileCheck,
    title: "Evidence Bundle",
    description: "Cryptographically signed PDF/JSON compliance pack generated automatically. Ready for DORA, NIS2, and AI Act audits."
  }
];

export const HowItWorks = () => {
  return (
    <section className="py-32 px-4 relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-subtle"></div>
      
      <div className="container mx-auto max-w-6xl relative z-10">
        <div className="text-center mb-20 space-y-4">
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            How <span className="text-gradient">Driftlock Works</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Edge deployment. Zero data retention. Explainable results.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {steps.map((step, idx) => (
            <div key={idx} className="relative">
              {/* Connector Line */}
              {idx < steps.length - 1 && (
                <div className="hidden lg:block absolute top-20 left-[calc(100%+1rem)] w-8 border-t-2 border-dashed border-primary/30">
                  <ArrowRight className="absolute -right-2 -top-3 w-5 h-5 text-primary" />
                </div>
              )}
              
              <div className="glass-card rounded-2xl p-6 hover-lift h-full">
                {/* Step Number */}
                <div className="text-5xl font-display font-bold text-primary/20 mb-4">
                  {step.number}
                </div>
                
                {/* Icon */}
                <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center mb-6">
                  <step.icon className="w-7 h-7 text-primary" />
                </div>
                
                {/* Content */}
                <h3 className="font-display font-bold text-xl mb-3">
                  {step.title}
                </h3>
                
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {step.description}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Key Differentiators */}
        <div className="mt-16 flex flex-wrap justify-center gap-6">
          <div className="glass-card rounded-full px-6 py-3 text-sm font-medium">
            ⚡ Edge Deployment
          </div>
          <div className="glass-card rounded-full px-6 py-3 text-sm font-medium">
            🔒 No Data Storage
          </div>
          <div className="glass-card rounded-full px-6 py-3 text-sm font-medium">
            📊 Field-Level Explanations
          </div>
          <div className="glass-card rounded-full px-6 py-3 text-sm font-medium">
            ⚖️ Compliance-Ready
          </div>
        </div>
      </div>
    </section>
  );
};
