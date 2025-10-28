import { AlertCircle, DollarSign, FileQuestion, Lock } from "lucide-react";

const problems = [
  {
    icon: AlertCircle,
    title: "ML-based anomaly detection is opaque and hard to audit",
    description: "Black-box models can't explain why they flagged an anomaly, making compliance audits a nightmare."
  },
  {
    icon: DollarSign,
    title: "Traditional monitoring is expensive and inflexible",
    description: "Built for host/agent models, not IoT, finance feeds, or LLM I/O. Per-host pricing doesn't scale."
  },
  {
    icon: FileQuestion,
    title: "Regulators require explainability",
    description: "DORA, NIS2, and AI Act demand root-cause analysis, not just alerts. You need verifiable evidence."
  },
  {
    icon: Lock,
    title: "Data privacy concerns with training models",
    description: "Sending sensitive telemetry to train ML models creates compliance and security risks."
  }
];

export const ProblemSection = () => {
  return (
    <section className="py-24 px-4 relative">
      <div className="container mx-auto max-w-6xl">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            Traditional Monitoring <span className="text-gradient">Falls Short</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            The stakes are high: downtime costs $5,600/minute on average, and compliance fines can reach €20M+ under DORA.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {problems.map((problem, idx) => (
            <div 
              key={idx}
              className="glass-card rounded-2xl p-8 hover-lift relative overflow-hidden group"
            >
              <div className="absolute inset-0 bg-gradient-primary opacity-0 group-hover:opacity-5 transition-opacity duration-500"></div>
              
              <div className="relative">
                <div className="w-12 h-12 rounded-xl bg-destructive/10 flex items-center justify-center mb-6">
                  <problem.icon className="w-6 h-6 text-destructive" />
                </div>
                
                <h3 className="font-display font-bold text-xl mb-3">
                  {problem.title}
                </h3>
                
                <p className="text-muted-foreground leading-relaxed">
                  {problem.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
