import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Code, Rocket, Book } from "lucide-react";
import { Link } from "react-router-dom";

export const DeveloperOnboarding = () => {
  const codeExample = `// Stream telemetry to Driftlock
const response = await fetch('https://api.driftlock.dev/v1/stream', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    stream_id: 'prod-api-logs',
    data: {
      timestamp: Date.now(),
      level: 'info',
      message: 'User login successful',
      user_id: 'usr_123',
      ip: '192.168.1.1'
    }
  })
});

// Monitor for anomalies
const anomaly = await fetch('https://api.driftlock.dev/v1/monitor', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    stream_id: 'prod-api-logs'
  })
});

const result = await anomaly.json();
// result.is_anomaly: true
// result.explanation: "message field 5x larger than baseline"`;

  return (
    <section className="py-32 px-4 relative">
      <div className="container mx-auto max-w-6xl">
        <div className="text-center mb-16 space-y-4">
          <Badge variant="secondary" className="mb-2">
            <Code className="w-3 h-3 mr-1" />
            Developer-First
          </Badge>
          <h2 className="text-3xl md:text-5xl font-display font-bold">
            Get Started in <span className="text-gradient">5 Minutes</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Simple REST API. No SDKs required. OpenTelemetry Collector support included.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          {/* Code Example */}
          <Card className="p-6 bg-slate-950 border-primary/20">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-mono text-muted-foreground">JavaScript/TypeScript</span>
              <Badge variant="outline" className="text-xs">REST API</Badge>
            </div>
            <pre className="text-xs font-mono text-slate-300 overflow-x-auto">
              <code>{codeExample}</code>
            </pre>
          </Card>

          {/* Quick Start Steps */}
          <div className="space-y-6">
            <div className="glass-card rounded-xl p-6">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <span className="text-lg font-bold text-primary">1</span>
                </div>
                <div>
                  <h3 className="font-semibold mb-2">Sign up for free</h3>
                  <p className="text-sm text-muted-foreground">
                    Get 10,000 API calls/month on the Developer plan. No credit card required.
                  </p>
                </div>
              </div>
            </div>

            <div className="glass-card rounded-xl p-6">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <span className="text-lg font-bold text-primary">2</span>
                </div>
                <div>
                  <h3 className="font-semibold mb-2">Stream your data</h3>
                  <p className="text-sm text-muted-foreground">
                    Send telemetry via REST API or configure OpenTelemetry Collector. Works with any structured data.
                  </p>
                </div>
              </div>
            </div>

            <div className="glass-card rounded-xl p-6">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <span className="text-lg font-bold text-primary">3</span>
                </div>
                <div>
                  <h3 className="font-semibold mb-2">Get explainable alerts</h3>
                  <p className="text-sm text-muted-foreground">
                    Receive real-time anomalies with field-level explanations. Download compliance evidence bundles anytime.
                  </p>
                </div>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-4 pt-4">
              <Link to="/docs" className="flex-1">
                <Button size="lg" className="w-full bg-gradient-primary">
                  <Rocket className="w-4 h-4 mr-2" />
                  Start Free
                </Button>
              </Link>
              <Link to="/docs" className="flex-1">
                <Button size="lg" variant="outline" className="w-full">
                  <Book className="w-4 h-4 mr-2" />
                  API Docs
                </Button>
              </Link>
            </div>
          </div>
        </div>

        {/* Language Support */}
        <div className="mt-12 text-center">
          <p className="text-sm text-muted-foreground mb-4">SDKs coming soon</p>
          <div className="flex flex-wrap justify-center gap-3">
            <Badge variant="outline">Go</Badge>
            <Badge variant="outline">Rust</Badge>
            <Badge variant="outline">TypeScript</Badge>
            <Badge variant="outline">Python</Badge>
          </div>
        </div>
      </div>
    </section>
  );
};
