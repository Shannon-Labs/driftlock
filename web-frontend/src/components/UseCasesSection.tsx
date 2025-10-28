import { CreditCard, Cpu, Bot, CloudRain } from "lucide-react";

const useCases = [
  {
    icon: CreditCard,
    industry: "Fintech",
    title: "Credit Card Transaction Monitoring",
    problem: "Detect anomalous transaction patterns without storing sensitive cardholder data.",
    solution: "Driftlock compresses transaction logs in real time, flagging unusual amounts, geographies, or merchant categories with field-level explanations."
  },
  {
    icon: Cpu,
    industry: "Manufacturing IoT",
    title: "Sensor Drift Detection",
    problem: "Identify failing sensors or unexpected environmental changes across thousands of devices.",
    solution: "Format-aware compression spots schema changes, out-of-range values, or message size spikes instantly."
  },
  {
    icon: Bot,
    industry: "AI/LLM Platforms",
    title: "Prompt/Response Anomalies",
    problem: "Monitor LLM I/O for unusual patterns, prompt injection attempts, or model drift.",
    solution: "Driftlock analyzes prompt length, token usage, and response structure changes — no PII storage required."
  },
  {
    icon: CloudRain,
    industry: "Data Feeds",
    title: "Weather/Market Data Integrity",
    problem: "Detect corrupt feeds, missing fields, or schema violations in high-frequency data streams.",
    solution: "Compression-based checks flag data quality issues before they propagate to downstream systems."
  }
];

export const UseCasesSection = () => {
  return (
    <section className="py-32 px-4 relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-subtle"></div>
      
      <div className="container mx-auto max-w-6xl relative z-10">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            Built for <span className="text-gradient">Regulated Industries</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            From finance to AI — wherever explainability and compliance matter
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {useCases.map((useCase, idx) => (
            <div 
              key={idx}
              className="glass-card rounded-2xl p-8 hover-lift relative overflow-hidden group"
            >
              <div className="absolute inset-0 bg-gradient-primary opacity-0 group-hover:opacity-5 transition-opacity duration-500"></div>
              
              <div className="relative">
                {/* Industry Badge */}
                <div className="inline-flex items-center gap-2 bg-primary/10 rounded-full px-3 py-1 text-xs font-medium text-primary mb-4">
                  <useCase.icon className="w-3 h-3" />
                  {useCase.industry}
                </div>
                
                <h3 className="font-display font-bold text-xl mb-4">
                  {useCase.title}
                </h3>
                
                <div className="space-y-4">
                  <div>
                    <div className="text-xs font-semibold text-destructive mb-1">Challenge</div>
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      {useCase.problem}
                    </p>
                  </div>
                  
                  <div>
                    <div className="text-xs font-semibold text-primary mb-1">Solution</div>
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      {useCase.solution}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
