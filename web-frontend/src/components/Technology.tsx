import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Code2, Database, Boxes, Activity } from "lucide-react";

export const Technology = () => {
  return (
    <section className="py-24 px-4 relative">
      <div className="container mx-auto">
        <div className="max-w-3xl mx-auto text-center mb-16 space-y-4">
          <h2 className="text-4xl md:text-5xl font-bold">
            Powered by <span className="text-gradient">Advanced Technology</span>
          </h2>
          <p className="text-xl text-muted-foreground">
            Format-aware compression meets OpenTelemetry observability
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 max-w-6xl mx-auto">
          {/* OpenZL Card */}
          <Card className="p-8 bg-gradient-card border-primary/10 space-y-6">
            <div className="flex items-start justify-between">
              <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center">
                <Code2 className="w-7 h-7 text-primary" />
              </div>
              <Badge variant="secondary" className="bg-primary/20 text-primary border-primary/30">
                Core Technology
              </Badge>
            </div>

            <div className="space-y-4">
              <h3 className="text-2xl font-bold">Meta's OpenZL Framework</h3>
              <p className="text-muted-foreground leading-relaxed">
                Format-aware compression that understands your data structure (JSON logs, timeseries metrics, nested traces) 
                rather than treating it as opaque bytes.
              </p>

              <div className="space-y-3 pt-4">
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-primary mt-2"></div>
                  <div>
                    <div className="font-medium">1.5-2x Better Compression</div>
                    <div className="text-sm text-muted-foreground">Compared to zstd on structured data</div>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-primary mt-2"></div>
                  <div>
                    <div className="font-medium">20-40% Faster Speed</div>
                    <div className="text-sm text-muted-foreground">Compression and decompression performance</div>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-primary mt-2"></div>
                  <div>
                    <div className="font-medium">Deterministic Output</div>
                    <div className="text-sm text-muted-foreground">Fixed compression plans for audit trails</div>
                  </div>
                </div>
              </div>
            </div>
          </Card>

          {/* CBAD Card */}
          <Card className="p-8 bg-gradient-card border-primary/10 space-y-6">
            <div className="flex items-start justify-between">
              <div className="w-14 h-14 rounded-xl bg-secondary/10 flex items-center justify-center">
                <Activity className="w-7 h-7 text-secondary" />
              </div>
              <Badge variant="secondary" className="bg-secondary/20 text-secondary border-secondary/30">
                Detection Engine
              </Badge>
            </div>

            <div className="space-y-4">
              <h3 className="text-2xl font-bold">Compression-Based Anomaly Detection</h3>
              <p className="text-muted-foreground leading-relaxed">
                CBAD leverages compression theory (Kolmogorov complexity) to detect anomalies through 
                changes in data compressibility patterns.
              </p>

              <div className="space-y-3 pt-4">
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-secondary mt-2"></div>
                  <div>
                    <div className="font-medium">Explainable Detection</div>
                    <div className="text-sm text-muted-foreground">Field-level attribution of anomalies</div>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-secondary mt-2"></div>
                  <div>
                    <div className="font-medium">Statistical Significance</div>
                    <div className="text-sm text-muted-foreground">Permutation testing with p-values</div>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-secondary mt-2"></div>
                  <div>
                    <div className="font-medium">No Training Required</div>
                    <div className="text-sm text-muted-foreground">Mathematical foundation, not ML models</div>
                  </div>
                </div>
              </div>
            </div>
          </Card>

          {/* Architecture Card */}
          <Card className="p-8 bg-gradient-card border-primary/10 space-y-6 lg:col-span-2">
            <div className="flex items-center gap-4 mb-6">
              <div className="w-14 h-14 rounded-xl bg-accent/10 flex items-center justify-center">
                <Boxes className="w-7 h-7 text-accent" />
              </div>
              <div>
                <h3 className="text-2xl font-bold">Enterprise Architecture</h3>
                <p className="text-muted-foreground">Production-ready technology stack</p>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Database className="w-5 h-5 text-primary" />
                  <h4 className="font-semibold">Core Engine</h4>
                </div>
                <ul className="text-sm text-muted-foreground space-y-1.5 pl-7">
                  <li>Rust CBAD core</li>
                  <li>OpenZL integration</li>
                  <li>C FFI bindings</li>
                  <li>Go API server</li>
                </ul>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Activity className="w-5 h-5 text-secondary" />
                  <h4 className="font-semibold">Data Pipeline</h4>
                </div>
                <ul className="text-sm text-muted-foreground space-y-1.5 pl-7">
                  <li>OTel Collector</li>
                  <li>PostgreSQL storage</li>
                  <li>Real-time streaming</li>
                  <li>WebSocket/SSE</li>
                </ul>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Code2 className="w-5 h-5 text-accent" />
                  <h4 className="font-semibold">Deployment</h4>
                </div>
                <ul className="text-sm text-muted-foreground space-y-1.5 pl-7">
                  <li>Docker containers</li>
                  <li>Kubernetes ready</li>
                  <li>Multi-region support</li>
                  <li>99.99% uptime</li>
                </ul>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </section>
  );
};
