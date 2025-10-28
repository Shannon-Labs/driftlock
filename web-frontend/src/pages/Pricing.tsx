import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Check, Zap, Building2, Sparkles } from "lucide-react";
import { Link } from "react-router-dom";

const Pricing = () => {
  const plans = [
    {
      name: "Developer",
      price: "Free",
      description: "Experimentation, testing, integration",
      icon: Zap,
      color: "text-primary",
      bgColor: "bg-primary/10",
      features: [
        "10k API calls/month",
        "Dual APIs (Stream + Monitor)",
        "Edge execution (Cloudflare)",
        "Explainable anomaly reports",
        "Community support",
        "7 day data retention",
      ],
    },
    {
      name: "Standard",
      price: "$49",
      period: "/month",
      description: "Startups and small teams",
      icon: Sparkles,
      color: "text-secondary",
      bgColor: "bg-secondary/10",
      popular: true,
      overage: "$0.004/call",
      features: [
        "100,000 total API calls",
        "Stream + Monitor/Explain APIs",
        "Ultra-low latency edge execution",
        "DORA/NIS2/AI Act evidence packs",
        "Priority support",
        "90 day data retention",
        "99.9% uptime SLA",
        "$0.004 per call overage",
      ],
    },
    {
      name: "Growth",
      price: "$249",
      period: "/month",
      description: "Expanding services, regulated workloads",
      icon: Building2,
      color: "text-accent",
      bgColor: "bg-accent/10",
      overage: "$0.002/call",
      features: [
        "1 million total API calls",
        "All Standard features",
        "Custom thresholds",
        "Dedicated region deployment",
        "Private Slack/Teams support",
        "Custom retention periods",
        "99.9% uptime SLA",
        "$0.002 per call overage",
      ],
    },
    {
      name: "Enterprise",
      price: "Custom",
      description: "Large-scale telemetry, critical infrastructure",
      icon: Building2,
      color: "text-accent",
      bgColor: "bg-accent/10",
      features: [
        "Unlimited API calls",
        "Volume contract pricing ($0.001/call)",
        "Multi-region deployment",
        "24/7 dedicated support",
        "Custom SLA (99.99%+)",
        "On-premises deployment",
        "Signed audit trail archive",
        "White-label dashboard for MSPs",
        "Professional services",
      ],
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-7xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Transparent, predictable, and <span className="text-gradient">compliance-ready pricing</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              Driftlock runs at the edge, charges only for what you analyze, and never stores your data. 
              Every tier includes full OpenZL-powered anomaly detection, explainable insights, and built-in 
              compliance evidence bundles for frameworks like <strong>DORA</strong>, <strong>NIS2</strong>, and the <strong>EU AI Act</strong>.
            </p>
          </div>

          {/* Pricing Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-16">
            {plans.map((plan, idx) => (
              <Card 
                key={idx} 
                className={`p-8 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all relative ${
                  plan.popular ? "ring-2 ring-primary" : ""
                }`}
              >
                {plan.popular && (
                  <Badge className="absolute -top-3 left-1/2 -translate-x-1/2 bg-gradient-primary">
                    Most Popular
                  </Badge>
                )}

                <div className={`w-16 h-16 rounded-xl ${plan.bgColor} flex items-center justify-center mb-6`}>
                  <plan.icon className={`w-8 h-8 ${plan.color}`} />
                </div>

                <h3 className="text-2xl font-bold mb-2">{plan.name}</h3>
                <p className="text-sm text-muted-foreground mb-6">{plan.description}</p>

                <div className="mb-6">
                  <span className="text-4xl font-bold">{plan.price}</span>
                  {plan.period && <span className="text-muted-foreground">{plan.period}</span>}
                </div>

                <Link to="/docs">
                  <Button 
                    className={`w-full mb-6 ${
                      plan.popular 
                        ? "bg-gradient-primary" 
                        : "bg-background hover:bg-muted"
                    }`}
                  >
                    {plan.name === "Enterprise" ? "Contact Sales" : "Get Started"}
                  </Button>
                </Link>

                <div className="space-y-3">
                  {plan.features.map((feature, featureIdx) => (
                    <div key={featureIdx} className="flex items-start gap-3">
                      <Check className={`w-5 h-5 ${plan.color} flex-shrink-0 mt-0.5`} />
                      <span className="text-sm">{feature}</span>
                    </div>
                  ))}
                </div>
              </Card>
            ))}
          </div>

          {/* All Plans Include Section */}
          <div className="max-w-4xl mx-auto mb-16">
            <Card className="p-8 bg-gradient-card border-primary/10">
              <h2 className="text-2xl font-bold text-center mb-6">All plans include</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">Dual APIs (Stream + Monitor/Explain)</span>
                </div>
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">Edge execution via Cloudflare Workers</span>
                </div>
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">Ultra-low latency, no data exfiltration</span>
                </div>
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">Deterministic, explainable anomaly reports</span>
                </div>
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">DORA/NIS2/AI Act evidence bundles</span>
                </div>
                <div className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                  <span className="text-sm">99.9% uptime SLA (paid plans)</span>
                </div>
              </div>
            </Card>
          </div>

          {/* Add-ons Section */}
          <div className="max-w-4xl mx-auto mb-16">
            <h2 className="text-3xl font-bold text-center mb-8">Add-ons</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">Dedicated Region Deployment</h3>
                <p className="text-sm text-muted-foreground">
                  EU, US, or APAC region-specific deployment for data sovereignty compliance.
                </p>
              </Card>
              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">Private Slack/Teams Support</h3>
                <p className="text-sm text-muted-foreground">
                  Direct channel access to our engineering team for real-time support.
                </p>
              </Card>
              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">Signed Audit Trail Archive</h3>
                <p className="text-sm text-muted-foreground">
                  Cryptographically signed evidence storage for long-term compliance audits.
                </p>
              </Card>
              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">White-Label Dashboard</h3>
                <p className="text-sm text-muted-foreground">
                  Custom-branded monitoring interface for MSPs and resellers.
                </p>
              </Card>
            </div>
          </div>

          {/* FAQ Section */}
          <div className="max-w-3xl mx-auto">
            <h2 className="text-3xl font-bold text-center mb-8">
              Frequently Asked Questions
            </h2>
            
            <div className="space-y-6">
              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">How does call-based pricing work?</h3>
                <p className="text-sm text-muted-foreground">
                  API calls include any request to our Stream or Monitor/Explain endpoints. 
                  You're only charged for compression analysis requests — no hidden fees for data transfer or storage.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">What happens if I exceed my plan limits?</h3>
                <p className="text-sm text-muted-foreground">
                  We'll notify you when approaching your limit. Overage charges are transparent and predictable. 
                  You can upgrade anytime to get better per-call rates and higher baseline limits.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">Can I self-host Driftlock?</h3>
                <p className="text-sm text-muted-foreground">
                  Yes! Enterprise plans include on-premises deployment options. We provide Docker images, 
                  Kubernetes Helm charts, and full support for your infrastructure.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">What compliance frameworks are supported?</h3>
                <p className="text-sm text-muted-foreground">
                  All paid plans include built-in DORA (Digital Operational Resilience Act), NIS2 (Network and Information Security), 
                  and Runtime AI monitoring for EU AI Act compliance. Evidence bundles are cryptographically signed for audit trails.
                </p>
              </Card>

              <Card className="p-6 bg-gradient-card border-primary/10">
                <h3 className="font-semibold mb-2">How do I get started with Driftlock?</h3>
                <p className="text-sm text-muted-foreground">
                  Sign up for a free Developer account with 10k included calls per month. Upgrade to paid plans anytime 
                  to access higher quotas and advanced features. Enterprise customers can schedule a custom proof-of-concept 
                  deployment with our solutions team.
                </p>
              </Card>
            </div>
          </div>

          {/* CTA Section */}
          <div className="mt-16 text-center">
            <Card className="p-12 bg-gradient-card border-primary/10 max-w-3xl mx-auto">
              <h2 className="text-3xl font-bold mb-4">
                Start detecting explainable anomalies today
              </h2>
              <p className="text-muted-foreground mb-8">
                Deploy Driftlock in minutes — no model training, no data retention, full compliance visibility.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Link to="/docs">
                  <Button size="lg" className="bg-gradient-primary">
                    Get Started
                  </Button>
                </Link>
                <Link to="/docs">
                  <Button size="lg" variant="outline">
                    View API Docs
                  </Button>
                </Link>
              </div>
            </Card>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Pricing;
