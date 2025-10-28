import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Target, Users, Lightbulb, Award } from "lucide-react";

const About = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-6xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              About <span className="text-gradient">Shannon Labs</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              Building the future of explainable anomaly detection for regulated industries
            </p>
          </div>

          {/* Mission */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 mb-12">
            <div className="flex items-start gap-4 mb-6">
              <div className="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                <Target className="w-8 h-8 text-primary" />
              </div>
              <div>
                <h2 className="text-3xl font-bold mb-4">Our Mission</h2>
                <p className="text-lg text-muted-foreground leading-relaxed">
                  We believe that anomaly detection systems should be explainable, deterministic, and built for 
                  the demands of regulated industries. Traditional black-box ML approaches fail when you need to 
                  explain <em>why</em> something was flagged as anomalous to regulators, auditors, or stakeholders.
                </p>
                <p className="text-lg text-muted-foreground leading-relaxed mt-4">
                  Driftlock combines compression-based anomaly detection (CBAD) with Meta's OpenZL format-aware 
                  compression framework to deliver glass-box explanations that meet DORA, NIS2, and AI Act compliance requirements.
                </p>
              </div>
            </div>
          </Card>

          {/* Values */}
          <div className="mb-16">
            <h2 className="text-3xl font-bold text-center mb-8">Our Values</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                  <Lightbulb className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-xl font-bold mb-3">Transparency First</h3>
                <p className="text-muted-foreground">
                  Every anomaly detection comes with mathematical explanations. No black boxes, 
                  no hidden algorithms—just pure compression theory and statistical significance.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                <div className="w-12 h-12 rounded-lg bg-secondary/10 flex items-center justify-center mb-4">
                  <Award className="w-6 h-6 text-secondary" />
                </div>
                <h3 className="text-xl font-bold mb-3">Compliance Ready</h3>
                <p className="text-muted-foreground">
                  Built from day one for DORA, NIS2, and Runtime AI compliance. Evidence bundles, 
                  audit trails, and cryptographic integrity come standard.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
                <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mb-4">
                  <Users className="w-6 h-6 text-accent" />
                </div>
                <h3 className="text-xl font-bold mb-3">Developer Focused</h3>
                <p className="text-muted-foreground">
                  Clean APIs, comprehensive docs, and native OpenTelemetry integration. 
                  Deploy in minutes, not months.
                </p>
              </Card>
            </div>
          </div>

          {/* Story */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 mb-12">
            <h2 className="text-3xl font-bold mb-6">The Driftlock Story</h2>
            <div className="space-y-4 text-muted-foreground">
              <p>
                Driftlock was born from frustration with existing anomaly detection systems that couldn't 
                explain their decisions. When a financial services company needed to justify anomaly alerts 
                to regulators, they found that "the neural network said so" wasn't an acceptable answer.
              </p>
              <p>
                We realized that compression-based anomaly detection (CBAD), grounded in Kolmogorov complexity 
                theory, could provide the mathematical rigor and explainability that regulated industries demand. 
                By leveraging Meta's OpenZL format-aware compression framework, we could achieve better detection 
                rates than traditional ML while maintaining complete transparency.
              </p>
              <p>
                Today, Driftlock powers anomaly detection for teams that need to prove their systems work—not 
                just trust that they do.
              </p>
            </div>
          </Card>

          {/* Team */}
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-4">Built by Shannon Labs</h2>
            <p className="text-muted-foreground mb-6">
              A research-driven team focused on bringing academic rigor to production systems
            </p>
            <div className="flex justify-center gap-4">
              <a 
                href="mailto:hunter@shannonlabs.dev"
                className="text-primary hover:text-primary/80 transition-colors font-medium"
              >
                Get in Touch →
              </a>
            </div>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default About;
