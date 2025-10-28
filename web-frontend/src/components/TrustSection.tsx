import { Shield, Lock, FileCheck, Cloud } from "lucide-react";

const trustPoints = [
  {
    icon: Shield,
    title: "Zero Data Retention",
    description: "Your telemetry never leaves the edge. We analyze, not store."
  },
  {
    icon: Lock,
    title: "Cryptographically Signed",
    description: "Every evidence bundle includes tamper-proof audit trails."
  },
  {
    icon: FileCheck,
    title: "Audit-Ready Evidence",
    description: "PDF/JSON compliance packs for DORA, NIS2, AI Act inspections."
  },
  {
    icon: Cloud,
    title: "Edge Architecture",
    description: "Cloudflare Workers deployment. No centralized data collection."
  }
];

export const TrustSection = () => {
  return (
    <section className="py-32 px-4 relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-subtle"></div>
      
      <div className="container mx-auto max-w-6xl relative z-10">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            Built for <span className="text-gradient">Trust & Compliance</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Engineered to meet the world's toughest data-protection standards
          </p>
        </div>

        {/* Trust Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
          {trustPoints.map((point, idx) => (
            <div 
              key={idx}
              className="glass-card rounded-2xl p-6 text-center hover-lift"
            >
              <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-4">
                <point.icon className="w-7 h-7 text-primary" />
              </div>
              <h3 className="font-semibold mb-2">{point.title}</h3>
              <p className="text-sm text-muted-foreground">{point.description}</p>
            </div>
          ))}
        </div>

        {/* Compliance Badges */}
        <div className="glass-card rounded-2xl p-12 text-center">
          <h3 className="text-xl font-semibold mb-6">Compliance Frameworks Supported</h3>
          <div className="flex flex-wrap justify-center gap-4">
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              EU AI Act
            </div>
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              DORA
            </div>
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              NIS2
            </div>
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              CPRA
            </div>
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              SOC 2 Type II
            </div>
            <div className="px-6 py-3 rounded-full bg-primary/10 border border-primary/20 font-semibold">
              ISO 27001
            </div>
          </div>

          <p className="text-sm text-muted-foreground mt-8">
            Used by regulated organizations in finance, healthcare, manufacturing, and AI/ML platforms
          </p>
        </div>
      </div>
    </section>
  );
};
