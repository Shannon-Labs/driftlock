import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Server, Layers, Database, LayoutDashboard, FileJson, Radio } from "lucide-react";

export const Architecture = () => {
  return (
    <section className="py-24 px-4 bg-gradient-subtle relative">
      <div className="container mx-auto">
        <div className="max-w-3xl mx-auto text-center mb-16 space-y-4">
          <h2 className="text-4xl md:text-5xl font-bold">
            Enterprise <span className="text-gradient">Architecture</span>
          </h2>
          <p className="text-xl text-muted-foreground">
            Production-ready components designed for scale and reliability
          </p>
        </div>

        {/* Architecture Flow */}
        <div className="max-w-6xl mx-auto mb-16">
          <Card className="p-8 bg-gradient-card border-primary/10">
            <h3 className="text-2xl font-bold mb-8 text-center">Data Flow Pipeline</h3>
            <div className="flex flex-col md:flex-row items-center justify-between gap-4 text-center">
              <div className="flex-1 space-y-2">
                <div className="w-16 h-16 rounded-lg bg-primary/20 flex items-center justify-center mx-auto">
                  <Radio className="w-8 h-8 text-primary" />
                </div>
                <div className="font-semibold">OTLP Sources</div>
                <div className="text-xs text-muted-foreground">Logs, Metrics, Traces</div>
              </div>
              <div className="text-muted-foreground">→</div>
              <div className="flex-1 space-y-2">
                <div className="w-16 h-16 rounded-lg bg-secondary/20 flex items-center justify-center mx-auto">
                  <Layers className="w-8 h-8 text-secondary" />
                </div>
                <div className="font-semibold">OTel Collector</div>
                <div className="text-xs text-muted-foreground">Process & Route</div>
              </div>
              <div className="text-muted-foreground">→</div>
              <div className="flex-1 space-y-2">
                <div className="w-16 h-16 rounded-lg bg-accent/20 flex items-center justify-center mx-auto">
                  <Server className="w-8 h-8 text-accent" />
                </div>
                <div className="font-semibold">CBAD Engine</div>
                <div className="text-xs text-muted-foreground">Detect Anomalies</div>
              </div>
              <div className="text-muted-foreground">→</div>
              <div className="flex-1 space-y-2">
                <div className="w-16 h-16 rounded-lg bg-primary/20 flex items-center justify-center mx-auto">
                  <Database className="w-8 h-8 text-primary" />
                </div>
                <div className="font-semibold">Storage</div>
                <div className="text-xs text-muted-foreground">PostgreSQL</div>
              </div>
              <div className="text-muted-foreground">→</div>
              <div className="flex-1 space-y-2">
                <div className="w-16 h-16 rounded-lg bg-secondary/20 flex items-center justify-center mx-auto">
                  <LayoutDashboard className="w-8 h-8 text-secondary" />
                </div>
                <div className="font-semibold">Dashboard</div>
                <div className="text-xs text-muted-foreground">Real-time UI</div>
              </div>
            </div>
          </Card>
        </div>

        {/* Component Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-7xl mx-auto">
          {/* CBAD Core */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center">
                <Layers className="w-6 h-6 text-primary" />
              </div>
              <Badge variant="secondary" className="bg-primary/20 text-primary">Rust</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">CBAD Core</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Compression-based algorithms with FFI bindings for Go and WASM target
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• OpenZL integration</li>
              <li>• C FFI bindings</li>
              <li>• Deterministic detection</li>
            </ul>
          </Card>

          {/* Collector Processor */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-secondary/10 flex items-center justify-center">
                <Radio className="w-6 h-6 text-secondary" />
              </div>
              <Badge variant="secondary" className="bg-secondary/20 text-secondary">Go</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">OTel Collector</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Native OpenTelemetry processor for logs, metrics, and traces
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• OTLP compatibility</li>
              <li>• Real-time streaming</li>
              <li>• LLM I/O receivers</li>
            </ul>
          </Card>

          {/* API Server */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center">
                <Server className="w-6 h-6 text-accent" />
              </div>
              <Badge variant="secondary" className="bg-accent/20 text-accent">Go</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">API Server</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Storage and retrieval with WebSocket/SSE streaming
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• PostgreSQL backend</li>
              <li>• GraphQL support</li>
              <li>• Real-time alerts</li>
            </ul>
          </Card>

          {/* Database */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center">
                <Database className="w-6 h-6 text-primary" />
              </div>
              <Badge variant="secondary" className="bg-primary/20 text-primary">Storage</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">Data Layer</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Time-series optimized storage with tiered archival
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• PostgreSQL hot storage</li>
              <li>• ClickHouse analytics</li>
              <li>• S3 cold archive</li>
            </ul>
          </Card>

          {/* Evidence Exporters */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-secondary/10 flex items-center justify-center">
                <FileJson className="w-6 h-6 text-secondary" />
              </div>
              <Badge variant="secondary" className="bg-secondary/20 text-secondary">Compliance</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">Evidence Bundles</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Cryptographically signed export packages for audit trails
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• JSON + PDF export</li>
              <li>• Digital signatures</li>
              <li>• DORA/NIS2 templates</li>
            </ul>
          </Card>

          {/* UI Dashboard */}
          <Card className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all">
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center">
                <LayoutDashboard className="w-6 h-6 text-accent" />
              </div>
              <Badge variant="secondary" className="bg-accent/20 text-accent">Next.js</Badge>
            </div>
            <h4 className="text-lg font-bold mb-2">Dashboard</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Real-time visualization with advanced analytics and filtering
            </p>
            <ul className="text-xs text-muted-foreground space-y-1">
              <li>• Interactive timelines</li>
              <li>• Root cause analysis</li>
              <li>• Custom dashboards</li>
            </ul>
          </Card>
        </div>

        {/* Performance Metrics */}
        <div className="mt-16 max-w-4xl mx-auto">
          <Card className="p-8 bg-gradient-card border-primary/10">
            <h3 className="text-2xl font-bold mb-6 text-center">Enterprise Performance</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              <div className="text-center">
                <div className="text-3xl font-bold text-primary mb-1">100k+</div>
                <div className="text-sm text-muted-foreground">Events/second</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-secondary mb-1">&lt;100ms</div>
                <div className="text-sm text-muted-foreground">API latency (p95)</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-accent mb-1">99.99%</div>
                <div className="text-sm text-muted-foreground">Uptime SLA</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-primary mb-1">100%</div>
                <div className="text-sm text-muted-foreground">Deterministic</div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </section>
  );
};
