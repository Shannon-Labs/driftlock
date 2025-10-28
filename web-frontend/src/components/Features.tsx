import { Brain, Lock, Zap, LineChart, FileCheck, Gauge } from "lucide-react";

const features = [
  {
    icon: Brain,
    title: "Glass-Box Explainability",
    description: "Understand every anomaly with field-level compression metrics. No black boxes, no model training — just deterministic, mathematical explanations you can audit.",
    highlight: true
  },
  {
    icon: Lock,
    title: "Privacy by Design",
    description: "Driftlock never trains on or stores your data — perfect for regulated environments. All processing happens at the edge with zero data retention.",
  },
  {
    icon: FileCheck,
    title: "Compliance-Ready",
    description: "Generate signed evidence bundles for DORA, NIS2, and AI Act audits. Every anomaly includes cryptographically verifiable compliance artifacts.",
  },
];

export const Features = () => {
  return (
    <section className="py-32 px-4 relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-subtle"></div>
      
      <div className="container mx-auto relative z-10">
        <div className="max-w-3xl mx-auto text-center mb-20 space-y-6">
          <h2 className="text-4xl md:text-6xl font-display font-bold">
            Why <span className="text-gradient">Driftlock</span>?
          </h2>
          <p className="text-xl md:text-2xl text-muted-foreground font-light">
            Built for regulated industries that demand explainability, compliance, and performance
          </p>
        </div>

        {/* Value Proposition Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          {features.map((feature, idx) => (
            <div
              key={idx}
              className="group glass-card rounded-2xl p-8 hover-lift relative overflow-hidden"
            >
              {/* Gradient Overlay on Hover */}
              <div className="absolute inset-0 bg-gradient-primary opacity-0 group-hover:opacity-5 transition-opacity duration-500"></div>
              
              <div className="relative">
                {/* Icon */}
                <div className="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                  <feature.icon className="w-8 h-8 text-primary" />
                </div>

                {/* Content */}
                <h3 className="font-display font-bold text-2xl mb-4">
                  {feature.title}
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  {feature.description}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Bottom CTA */}
        <div className="text-center mt-20">
          <p className="text-muted-foreground mb-6">
            See how it works in production
          </p>
          <a 
            href="/dashboard"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-xl glass-card hover-lift text-sm font-medium"
          >
            Explore Dashboard
            <span className="text-primary">→</span>
          </a>
        </div>
      </div>
    </section>
  );
};
