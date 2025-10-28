import { Check, X } from "lucide-react";
import { Card } from "@/components/ui/card";

const comparison = [
  {
    feature: "Detection Method",
    traditional: "Black-box ML models",
    driftlock: "Format-aware compression (OpenZL)"
  },
  {
    feature: "Explainability",
    traditional: "Generic anomaly scores",
    driftlock: "Field-level root cause analysis"
  },
  {
    feature: "Data Storage",
    traditional: "Stores data for training",
    driftlock: "Zero data retention"
  },
  {
    feature: "Compliance Evidence",
    traditional: "Manual audit preparation",
    driftlock: "Automatic signed evidence bundles"
  },
  {
    feature: "Domain Support",
    traditional: "Host/agent focused",
    driftlock: "Logs, metrics, traces, IoT, LLM I/O"
  },
  {
    feature: "Setup Time",
    traditional: "Weeks of training",
    driftlock: "Minutes (no training required)"
  }
];

export const WhyDriftlock = () => {
  return (
    <section className="py-32 px-4 relative">
      <div className="container mx-auto max-w-6xl">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            Why <span className="text-gradient">Driftlock</span>?
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Built for teams that need explainability, compliance, and performance — not just alerts.
          </p>
        </div>

        {/* Comparison Table */}
        <Card className="overflow-hidden bg-gradient-card border-primary/10">
          <div className="grid grid-cols-3 gap-px bg-border">
            {/* Header */}
            <div className="bg-background p-4"></div>
            <div className="bg-muted/50 p-4">
              <h3 className="font-semibold text-center">Traditional ML/Observability</h3>
            </div>
            <div className="bg-primary/5 p-4">
              <h3 className="font-semibold text-center text-primary">Driftlock</h3>
            </div>

            {/* Rows */}
            {comparison.map((row, idx) => (
              <>
                <div key={`feature-${idx}`} className="bg-background p-4 font-medium">
                  {row.feature}
                </div>
                <div key={`traditional-${idx}`} className="bg-muted/30 p-4">
                  <div className="flex items-start gap-2">
                    <X className="w-4 h-4 text-destructive flex-shrink-0 mt-0.5" />
                    <span className="text-sm text-muted-foreground">{row.traditional}</span>
                  </div>
                </div>
                <div key={`driftlock-${idx}`} className="bg-primary/5 p-4">
                  <div className="flex items-start gap-2">
                    <Check className="w-4 h-4 text-primary flex-shrink-0 mt-0.5" />
                    <span className="text-sm font-medium">{row.driftlock}</span>
                  </div>
                </div>
              </>
            ))}
          </div>
        </Card>

        {/* Bottom Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-16">
          <div className="text-center">
            <div className="text-4xl font-display font-bold text-gradient mb-2">1.5-2×</div>
            <div className="text-sm text-muted-foreground">Better compression vs generic tools</div>
          </div>
          <div className="text-center">
            <div className="text-4xl font-display font-bold text-gradient mb-2">100%</div>
            <div className="text-sm text-muted-foreground">Deterministic & reproducible</div>
          </div>
          <div className="text-center">
            <div className="text-4xl font-display font-bold text-gradient mb-2">&lt;400ms</div>
            <div className="text-sm text-muted-foreground">p95 latency at 10k+ events/sec</div>
          </div>
        </div>
      </div>
    </section>
  );
};
