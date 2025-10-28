import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Check, X } from "lucide-react";

const Comparison = () => {
  const features = [
    {
      category: "Detection Quality",
      items: [
        { feature: "Explainable anomalies", driftlock: true, ml: false },
        { feature: "Mathematical proof of detection", driftlock: true, ml: false },
        { feature: "Field-level attribution", driftlock: true, ml: false },
        { feature: "Statistical significance (p-values)", driftlock: true, ml: "partial" },
        { feature: "100% deterministic results", driftlock: true, ml: false },
        { feature: "No training required", driftlock: true, ml: false },
      ]
    },
    {
      category: "Compliance & Audit",
      items: [
        { feature: "DORA compliance ready", driftlock: true, ml: false },
        { feature: "NIS2 compliance ready", driftlock: true, ml: false },
        { feature: "AI Act compatible", driftlock: true, ml: false },
        { feature: "Cryptographic audit trails", driftlock: true, ml: "partial" },
        { feature: "Evidence bundle export", driftlock: true, ml: false },
        { feature: "Reproducible for auditors", driftlock: true, ml: false },
      ]
    },
    {
      category: "Performance",
      items: [
        { feature: "Real-time detection (<1s)", driftlock: true, ml: true },
        { feature: "100k+ events/second", driftlock: true, ml: "varies" },
        { feature: "Low false positive rate", driftlock: true, ml: "varies" },
        { feature: "Memory efficient", driftlock: true, ml: "varies" },
        { feature: "Horizontal scaling", driftlock: true, ml: true },
      ]
    },
    {
      category: "Operations",
      items: [
        { feature: "No model training", driftlock: true, ml: false },
        { feature: "No model drift issues", driftlock: true, ml: false },
        { feature: "Immediate deployment", driftlock: true, ml: false },
        { feature: "Works with small datasets", driftlock: true, ml: false },
        { feature: "Format-aware compression", driftlock: true, ml: false },
        { feature: "OpenTelemetry native", driftlock: true, ml: "partial" },
      ]
    },
  ];

  const renderCheckmark = (value: boolean | string) => {
    if (value === true) {
      return <Check className="w-5 h-5 text-green-500 mx-auto" />;
    } else if (value === false) {
      return <X className="w-5 h-5 text-red-500/50 mx-auto" />;
    } else {
      return <span className="text-xs text-yellow-500 mx-auto">{value}</span>;
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-6xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Driftlock vs. <span className="text-gradient">ML-Based Detection</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              Why compression-based anomaly detection is better for regulated industries
            </p>
          </div>

          {/* The Problem */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 mb-12">
            <h2 className="text-3xl font-bold mb-6">The Black Box Problem</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div>
                <h3 className="text-xl font-semibold mb-3 text-red-500/80">Traditional ML Approaches</h3>
                <ul className="space-y-2 text-muted-foreground">
                  <li>• "The neural network said so" is not an acceptable answer for regulators</li>
                  <li>• Model drift requires constant retraining and validation</li>
                  <li>• Non-deterministic results make auditing impossible</li>
                  <li>• Require large training datasets and expertise</li>
                  <li>• Can't explain why specific fields triggered detection</li>
                </ul>
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-3 text-green-500">Driftlock (CBAD)</h3>
                <ul className="space-y-2 text-muted-foreground">
                  <li>• Mathematical proof based on compression theory</li>
                  <li>• Deterministic: same input always produces same result</li>
                  <li>• Field-level explanations of what changed</li>
                  <li>• Works immediately with no training required</li>
                  <li>• Statistical significance with p-values</li>
                </ul>
              </div>
            </div>
          </Card>

          {/* Comparison Table */}
          <div className="space-y-8">
            {features.map((category, idx) => (
              <div key={idx}>
                <h2 className="text-2xl font-bold mb-4">{category.category}</h2>
                <Card className="overflow-hidden bg-gradient-card border-primary/10">
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead className="bg-muted/50">
                        <tr>
                          <th className="text-left p-4 font-semibold">Feature</th>
                          <th className="text-center p-4 font-semibold w-32">Driftlock</th>
                          <th className="text-center p-4 font-semibold w-32">ML-Based</th>
                        </tr>
                      </thead>
                      <tbody>
                        {category.items.map((item, itemIdx) => (
                          <tr key={itemIdx} className="border-t border-border">
                            <td className="p-4 text-sm">{item.feature}</td>
                            <td className="p-4 text-center">{renderCheckmark(item.driftlock)}</td>
                            <td className="p-4 text-center">{renderCheckmark(item.ml)}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </Card>
              </div>
            ))}
          </div>

          {/* Use Case Examples */}
          <Card className="mt-12 p-8 md:p-12 bg-gradient-card border-primary/10">
            <h2 className="text-3xl font-bold mb-6">Real-World Example</h2>
            <div className="space-y-6">
              <div>
                <h3 className="text-xl font-semibold mb-3">Scenario: Financial Transaction Anomaly</h3>
                <p className="text-muted-foreground mb-4">
                  A bank needs to explain to regulators why a specific transaction was flagged as anomalous.
                </p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card className="p-6 bg-background/50 border-red-500/20">
                  <h4 className="font-semibold mb-3 text-red-500/80">ML-Based System Response</h4>
                  <p className="text-sm text-muted-foreground mb-4">
                    "Our neural network detected an anomaly with 87% confidence. The model was trained on 
                    10M transactions and uses 200 hidden layers."
                  </p>
                  <p className="text-sm text-muted-foreground italic">
                    This explanation is not acceptable for DORA compliance and doesn't help investigators 
                    understand what actually went wrong.
                  </p>
                </Card>

                <Card className="p-6 bg-background/50 border-green-500/20">
                  <h4 className="font-semibold mb-3 text-green-500">Driftlock Response</h4>
                  <p className="text-sm text-muted-foreground mb-4">
                    "Transaction flagged due to: 1) Amount field compressed 43% worse than baseline 
                    (NCD=0.73, p=0.002), 2) New 'destination_country' field not in training schema, 
                    3) Timestamp pattern deviation."
                  </p>
                  <p className="text-sm text-muted-foreground italic">
                    Clear, actionable explanation with mathematical proof. Auditors can verify the 
                    detection logic and reproduce results.
                  </p>
                </Card>
              </div>
            </div>
          </Card>

          {/* CTA */}
          <Card className="mt-12 p-8 bg-gradient-card border-primary/10 text-center">
            <h3 className="text-2xl font-bold mb-4">Ready for Explainable Anomaly Detection?</h3>
            <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
              Join regulated organizations using Driftlock for compliance-ready anomaly detection
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a href="/docs">
                <button className="px-6 py-3 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                  Get Started Free
                </button>
              </a>
              <a href="/contact">
                <button className="px-6 py-3 bg-background border border-border rounded-lg font-medium hover:bg-muted transition-colors">
                  Schedule Demo
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

export default Comparison;
