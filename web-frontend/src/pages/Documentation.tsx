import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card } from "@/components/ui/card";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";

export default function Documentation() {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16">
        <div className="container mx-auto px-4">
          <div className="max-w-6xl mx-auto">
            <div className="text-center mb-12">
              <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Driftlock Documentation
              </h1>
              <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
                Comprehensive guide to compression-based anomaly detection using Meta's OpenZL
              </p>
            </div>

            <Tabs defaultValue="overview" className="space-y-8">
              <TabsList className="grid w-full grid-cols-4 lg:grid-cols-8 gap-2">
                <TabsTrigger value="overview">Overview</TabsTrigger>
                <TabsTrigger value="concepts">Concepts</TabsTrigger>
                <TabsTrigger value="api">API</TabsTrigger>
                <TabsTrigger value="auth">Auth</TabsTrigger>
                <TabsTrigger value="integration">Integration</TabsTrigger>
                <TabsTrigger value="billing">Billing</TabsTrigger>
                <TabsTrigger value="security">Security</TabsTrigger>
                <TabsTrigger value="advanced">Advanced</TabsTrigger>
              </TabsList>

              {/* OVERVIEW */}
              <TabsContent value="overview" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Driftlock Platform Overview</h2>
                  <div className="space-y-6">
                    <div>
                      <p className="text-lg mb-4">
                        Driftlock is a cutting-edge <strong>compression-based anomaly detection platform</strong> that leverages 
                        Meta's OpenZL format-aware compression framework to detect anomalies with unparalleled precision and explainability.
                      </p>
                      
                      <div className="bg-primary/10 p-6 rounded-lg border-l-4 border-primary mt-6">
                        <p className="font-semibold mb-2">API Base URL</p>
                        <code className="text-sm bg-muted px-3 py-1 rounded">https://api.driftlock.net/api/v1/</code>
                        <p className="text-sm text-muted-foreground mt-2">
                          All requests require Bearer token authentication
                        </p>
                      </div>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Core Technology</h3>
                      <p className="mb-4">
                        Driftlock integrates with <strong>Meta's OpenZL</strong>, a format-aware compression framework that understands 
                        structured data and applies intelligent compression strategies. Unlike traditional compression (zstd, gzip, lz4), 
                        OpenZL parses structure and learns optimal patterns.
                      </p>
                      
                      <div className="grid md:grid-cols-3 gap-4 mt-4">
                        <Card className="p-4 bg-gradient-to-br from-primary/5 to-primary/10">
                          <h4 className="font-semibold mb-2">1.5-2x Better</h4>
                          <p className="text-sm text-muted-foreground">Superior compression ratios on structured data</p>
                        </Card>
                        <Card className="p-4 bg-gradient-to-br from-secondary/5 to-secondary/10">
                          <h4 className="font-semibold mb-2">20-40% Faster</h4>
                          <p className="text-sm text-muted-foreground">Optimized speed with intelligent caching</p>
                        </Card>
                        <Card className="p-4 bg-gradient-to-br from-accent/5 to-accent/10">
                          <h4 className="font-semibold mb-2">Glass-Box</h4>
                          <p className="text-sm text-muted-foreground">Exact explanations of anomalies</p>
                        </Card>
                      </div>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Use Cases</h3>
                      <div className="grid md:grid-cols-2 gap-4">
                        {[
                          { icon: "🔍", title: "Log Analysis", desc: "Detect unusual error patterns and new field types" },
                          { icon: "📊", title: "Metric Monitoring", desc: "Identify distribution changes in timeseries" },
                          { icon: "🔗", title: "Trace Analysis", desc: "Detect structural changes in distributed tracing" },
                          { icon: "🤖", title: "LLM I/O Monitoring", desc: "Identify unusual patterns and hallucinations" },
                          { icon: "💰", title: "Financial Data", desc: "Detect unusual transaction patterns" },
                          { icon: "📡", title: "IoT Data", desc: "Identify sensor data anomalies" }
                        ].map((useCase, idx) => (
                          <Card key={idx} className="p-4">
                            <h4 className="font-semibold mb-2">{useCase.icon} {useCase.title}</h4>
                            <p className="text-sm text-muted-foreground">{useCase.desc}</p>
                          </Card>
                        ))}
                      </div>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Architecture</h3>
                      <div className="bg-muted p-6 rounded-lg">
                        <div className="space-y-2 font-mono text-sm">
                          <div><strong>Edge Layer:</strong> Cloudflare Workers API Gateway</div>
                          <div><strong>Backend:</strong> Supabase (Auth + Postgres + Edge Functions)</div>
                          <div><strong>Processing:</strong> Go Microservices with OpenZL Integration</div>
                          <div><strong>Storage:</strong> PostgreSQL with RLS multi-tenant security</div>
                        </div>
                      </div>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Compliance Ready</h3>
                      <div className="flex flex-wrap gap-2">
                        {["DORA", "NIS2", "EU AI Act", "SOC2", "GDPR"].map(tag => (
                          <Badge key={tag} variant="outline">{tag}</Badge>
                        ))}
                      </div>
                    </div>
                  </div>
                </Card>
              </TabsContent>

              {/* CONCEPTS */}
              <TabsContent value="concepts" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Concepts & Theory</h2>
                  
                  <Accordion type="single" collapsible className="space-y-4">
                    <AccordionItem value="openzl">
                      <AccordionTrigger className="text-xl font-semibold">
                        OpenZL Integration
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div>
                          <h4 className="font-semibold mb-2">What is OpenZL?</h4>
                          <p className="text-muted-foreground">
                            OpenZL is Meta's <strong>format-aware compression framework</strong> that understands structured data like JSON, 
                            logs, and metrics. Unlike byte-stream compressors, OpenZL parses structure and applies intelligent transforms.
                          </p>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-3">Key Differences</h4>
                          <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                              <thead className="bg-muted">
                                <tr>
                                  <th className="p-2 text-left">Feature</th>
                                  <th className="p-2 text-left">Traditional (zstd)</th>
                                  <th className="p-2 text-left">OpenZL</th>
                                </tr>
                              </thead>
                              <tbody>
                                <tr className="border-b">
                                  <td className="p-2">Data Understanding</td>
                                  <td className="p-2 text-muted-foreground">Byte streams</td>
                                  <td className="p-2 text-primary font-semibold">Parses structure</td>
                                </tr>
                                <tr className="border-b">
                                  <td className="p-2">Strategy</td>
                                  <td className="p-2 text-muted-foreground">Generic patterns</td>
                                  <td className="p-2 text-primary font-semibold">Format-specific</td>
                                </tr>
                                <tr className="border-b">
                                  <td className="p-2">Learning</td>
                                  <td className="p-2 text-muted-foreground">Static dictionaries</td>
                                  <td className="p-2 text-primary font-semibold">Learns optimal strategies</td>
                                </tr>
                                <tr>
                                  <td className="p-2">Performance</td>
                                  <td className="p-2 text-muted-foreground">Baseline</td>
                                  <td className="p-2 text-primary font-semibold">1.5-2x better, 20-40% faster</td>
                                </tr>
                              </tbody>
                            </table>
                          </div>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-2">Benefits for Anomaly Detection</h4>
                          <ul className="space-y-2 list-disc list-inside text-muted-foreground">
                            <li><strong className="text-foreground">Structural Awareness:</strong> Detects when data structure deviates</li>
                            <li><strong className="text-foreground">Field-Level Metrics:</strong> Per-field compression ratios</li>
                            <li><strong className="text-foreground">Deterministic:</strong> Reproducible anomaly detection</li>
                            <li><strong className="text-foreground">Explainability:</strong> Shows exactly which fields caused issues</li>
                          </ul>
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="cbad">
                      <AccordionTrigger className="text-xl font-semibold">
                        Compression-Based Anomaly Detection (CBAD)
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div>
                          <h4 className="font-semibold mb-2">Mathematical Foundation</h4>
                          <p className="text-muted-foreground mb-3">
                            CBAD is based on <strong>Kolmogorov complexity theory</strong>: the complexity of data is the length 
                            of the shortest program producing it. Compression approximates this complexity.
                          </p>
                          <div className="bg-muted p-4 rounded-lg font-mono text-sm">
                            K(x) ≈ |compressed(x)|<br/>
                            Anomaly Score = |compressed(x)| - baseline_compression
                          </div>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-2">Detection Process</h4>
                          <ol className="space-y-2 list-decimal list-inside text-muted-foreground">
                            <li><strong className="text-foreground">Training:</strong> OpenZL learns optimal strategies from baseline</li>
                            <li><strong className="text-foreground">Detection:</strong> New data compressed using learned strategies</li>
                            <li><strong className="text-foreground">Comparison:</strong> Significant ratio drops indicate deviation</li>
                            <li><strong className="text-foreground">Analysis:</strong> Field-level metrics identify root cause</li>
                          </ol>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-3">Field-Level Insights Example</h4>
                          <pre className="bg-muted p-4 rounded-lg text-sm overflow-x-auto">
{`{
  "field_compression_metrics": {
    "user_id": { "ratio": 2.5, "anomaly_score": 0.1 },
    "message": { "ratio": 1.2, "anomaly_score": 0.85 },
    "timestamp": { "ratio": 3.0, "anomaly_score": 0.05 }
  }
}`}
                          </pre>
                          <p className="text-sm text-muted-foreground mt-2">
                            The "message" field has poor compression (1.2x) and high anomaly score (0.85), indicating unusual patterns.
                          </p>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-3">Glass-Box vs Black-Box ML</h4>
                          <div className="grid md:grid-cols-2 gap-4">
                            <Card className="p-4 border-destructive/20">
                              <h5 className="font-semibold mb-2 flex items-center gap-2">
                                <span className="text-destructive">❌</span> Black-Box ML
                              </h5>
                              <ul className="text-sm space-y-1 text-muted-foreground">
                                <li>• Heuristic pattern matching</li>
                                <li>• Unexplainable decisions</li>
                                <li>• Requires labeled data</li>
                                <li>• Prone to false positives</li>
                              </ul>
                            </Card>
                            <Card className="p-4 bg-primary/5 border-primary/20">
                              <h5 className="font-semibold mb-2 flex items-center gap-2">
                                <span className="text-primary">✅</span> CBAD (Glass-Box)
                              </h5>
                              <ul className="text-sm space-y-1 text-muted-foreground">
                                <li>• Mathematical foundation</li>
                                <li>• Fully explainable</li>
                                <li>• Unsupervised learning</li>
                                <li>• Precise field insights</li>
                              </ul>
                            </Card>
                          </div>
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="anomaly-types">
                      <AccordionTrigger className="text-xl font-semibold">
                        Data Anomaly Types Explained
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        {[
                          {
                            icon: "🔍",
                            title: "Log Anomalies",
                            what: "Unusual message patterns, field lengths, new fields",
                            example: "10KB error messages when baseline is 200 bytes",
                            signal: "Message field compression drops from 3.5x to 1.1x"
                          },
                          {
                            icon: "📊",
                            title: "Metric Anomalies",
                            what: "Compression failure in timeseries from distribution changes",
                            example: "CPU metrics spike from 40-60% to 95-99%",
                            signal: "Numeric encoding less efficient due to wider range"
                          },
                          {
                            icon: "🔗",
                            title: "Trace Anomalies",
                            what: "Structural changes in distributed tracing",
                            example: "New service in call chain or unusual span durations",
                            signal: "Span structure deviates from learned patterns"
                          },
                          {
                            icon: "🤖",
                            title: "LLM I/O Anomalies",
                            what: "Unusual prompt/response patterns, hallucinations",
                            example: "Model producing repetitive or nonsensical outputs",
                            signal: "Response text compresses poorly due to randomness"
                          },
                          {
                            icon: "💰",
                            title: "Financial Data Anomalies",
                            what: "Unusual transaction patterns, field combinations",
                            example: "Transactions with unusual amounts or new categories",
                            signal: "Field correlations break down, reducing compression"
                          },
                          {
                            icon: "📡",
                            title: "IoT Data Anomalies",
                            what: "Sensor data patterns deviating from normal",
                            example: "Temperature rapidly oscillating vs steady baseline",
                            signal: "Timeseries compression drops due to higher entropy"
                          }
                        ].map((type, idx) => (
                          <div key={idx} className="border-l-4 border-primary pl-4">
                            <h4 className="font-semibold mb-2">{type.icon} {type.title}</h4>
                            <div className="space-y-1 text-sm">
                              <p><strong>What:</strong> <span className="text-muted-foreground">{type.what}</span></p>
                              <p><strong>Example:</strong> <span className="text-muted-foreground">{type.example}</span></p>
                              <p><strong>Compression Signal:</strong> <span className="text-muted-foreground">{type.signal}</span></p>
                            </div>
                          </div>
                        ))}
                      </AccordionContent>
                    </AccordionItem>
                  </Accordion>
                </Card>
              </TabsContent>

              {/* API REFERENCE */}
              <TabsContent value="api" className="space-y-6">
                <Card className="p-8">
                  <div className="flex items-center justify-between mb-6">
                    <h2 className="text-3xl font-bold">API Reference</h2>
                    <Badge>v1</Badge>
                  </div>

                  <Accordion type="single" collapsible className="space-y-4">
                    <AccordionItem value="anomalies">
                      <AccordionTrigger className="text-xl font-semibold">
                        Anomaly Detection Endpoints
                      </AccordionTrigger>
                      <AccordionContent className="space-y-6 pt-4">
                        <div className="space-y-4">
                          <div>
                            <div className="flex items-center gap-2 mb-2">
                              <Badge variant="secondary">GET</Badge>
                              <code className="text-sm">/anomalies</code>
                            </div>
                            <p className="text-sm text-muted-foreground mb-3">List detected anomalies with filtering</p>
                            
                            <h5 className="font-semibold text-sm mb-2">Query Parameters</h5>
                            <pre className="bg-muted p-3 rounded text-xs overflow-x-auto">
{`page: integer (default: 1)
limit: integer (default: 20, max: 100)
status: active | resolved | all
severity: low | medium | high | critical
source: string
from_date: ISO 8601 date
to_date: ISO 8601 date`}
                            </pre>

                            <h5 className="font-semibold text-sm mt-3 mb-2">Response</h5>
                            <pre className="bg-muted p-3 rounded text-xs overflow-x-auto">
{`{
  "anomalies": [{
    "id": 123,
    "type": "log",
    "severity": "high",
    "title": "Unusual Error Spike",
    "compression_metrics": {
      "ratio_change": -2.3,
      "affected_fields": ["message"],
      "deviation_score": 0.85
    },
    "explanation": "Compression ratio dropped...",
    "detected_at": "2025-10-26T14:30:00Z"
  }],
  "pagination": {...}
}`}
                            </pre>
                          </div>

                          <div className="border-t pt-4">
                            <div className="flex items-center gap-2 mb-2">
                              <Badge variant="secondary">GET</Badge>
                              <code className="text-sm">/anomalies/:id</code>
                            </div>
                            <p className="text-sm text-muted-foreground">Get detailed anomaly information</p>
                          </div>

                          <div className="border-t pt-4">
                            <div className="flex items-center gap-2 mb-2">
                              <Badge>PUT</Badge>
                              <code className="text-sm">/anomalies/:id/resolve</code>
                            </div>
                            <p className="text-sm text-muted-foreground mb-3">Mark anomaly as resolved</p>
                            <pre className="bg-muted p-3 rounded text-xs">
{`{
  "resolution_notes": "Resolved by scaling service"
}`}
                            </pre>
                          </div>

                          <div className="border-t pt-4">
                            <div className="flex items-center gap-2 mb-2">
                              <Badge>POST</Badge>
                              <code className="text-sm">/anomalies/detect</code>
                            </div>
                            <p className="text-sm text-muted-foreground mb-3">Manual anomaly detection</p>
                            <pre className="bg-muted p-3 rounded text-xs">
{`{
  "data": {...},
  "algorithm": "compression_based",
  "sensitivity": 0.5
}`}
                            </pre>
                          </div>
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="ingestion">
                      <AccordionTrigger className="text-xl font-semibold">
                        Data Ingestion Endpoints
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div>
                          <div className="flex items-center gap-2 mb-2">
                            <Badge>POST</Badge>
                            <code className="text-sm">/events/ingest</code>
                          </div>
                          <Badge variant="outline" className="mb-3">⚠️ Metered - only anomaly detection counts</Badge>
                          
                          <h5 className="font-semibold text-sm mb-2">Request</h5>
                          <pre className="bg-muted p-3 rounded text-xs overflow-x-auto">
{`{
  "timestamp": "2025-10-26T14:30:00Z",
  "type": "log|metric|trace|llm_io",
  "source": "service-name",
  "data": {...},
  "metadata": {
    "environment": "production"
  }
}`}
                          </pre>

                          <h5 className="font-semibold text-sm mt-3 mb-2">Response</h5>
                          <pre className="bg-muted p-3 rounded text-xs overflow-x-auto">
{`{
  "success": true,
  "event_id": "uuid",
  "compression_metrics": {
    "baseline_ratio": 2.5,
    "current_ratio": 1.2,
    "anomaly_detected": true
  }
}`}
                          </pre>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  </Accordion>
                </Card>
              </TabsContent>

              {/* AUTH */}
              <TabsContent value="auth" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Authentication</h2>
                  
                  <div className="space-y-6">
                    <div>
                      <h3 className="text-xl font-semibold mb-3">Bearer Token Authentication</h3>
                      <p className="text-muted-foreground mb-3">All API requests require authentication:</p>
                      <pre className="bg-muted p-4 rounded">
{`Authorization: Bearer YOUR_API_KEY`}
                      </pre>
                    </div>

                    <div>
                      <h3 className="text-xl font-semibold mb-4">Authentication Endpoints</h3>
                      <div className="space-y-4">
                        <div>
                          <div className="flex items-center gap-2 mb-2">
                            <Badge>POST</Badge>
                            <code className="text-sm">/auth/login</code>
                          </div>
                          <pre className="bg-muted p-3 rounded text-xs">
{`{
  "email": "user@example.com",
  "password": "password"
}

Response: {
  "token": "jwt_token",
  "user": {...}
}`}
                          </pre>
                        </div>

                        <div>
                          <div className="flex items-center gap-2 mb-2">
                            <Badge>POST</Badge>
                            <code className="text-sm">/auth/register</code>
                          </div>
                          <p className="text-sm text-muted-foreground">Create new account</p>
                        </div>
                      </div>
                    </div>

                    <div>
                      <h3 className="text-xl font-semibold mb-3">Example Request</h3>
                      <pre className="bg-muted p-4 rounded text-sm overflow-x-auto">
{`curl -X GET https://api.driftlock.net/api/v1/anomalies \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json"`}
                      </pre>
                    </div>
                  </div>
                </Card>
              </TabsContent>

              {/* INTEGRATION */}
              <TabsContent value="integration" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Integration Guides</h2>

                  <Accordion type="single" collapsible className="space-y-4">
                    <AccordionItem value="javascript">
                      <AccordionTrigger className="text-xl font-semibold">
                        JavaScript/Node.js SDK
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div>
                          <h4 className="font-semibold mb-2">Installation</h4>
                          <pre className="bg-muted p-3 rounded text-sm">
{`npm install @driftlock/client`}
                          </pre>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-2">Initialization</h4>
                          <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`import { DriftlockClient } from '@driftlock/client';

const client = new DriftlockClient({
  apiKey: 'YOUR_API_KEY',
  baseURL: 'https://api.driftlock.net'
});`}
                          </pre>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-2">Sending Data</h4>
                          <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`const result = await client.ingestEvent({
  type: 'log',
  source: 'my-service',
  data: { level: 'error', message: 'Failed' }
});

if (result.anomaly_detected) {
  console.log('Anomaly!', result.compression_metrics);
}`}
                          </pre>
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="python">
                      <AccordionTrigger className="text-xl font-semibold">
                        Python SDK
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div>
                          <h4 className="font-semibold mb-2">Installation</h4>
                          <pre className="bg-muted p-3 rounded text-sm">
{`pip install driftlock`}
                          </pre>
                        </div>

                        <div>
                          <h4 className="font-semibold mb-2">Usage</h4>
                          <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`from driftlock import DriftlockClient

client = DriftlockClient(api_key='YOUR_KEY')

result = client.ingest_event(
    type='log',
    source='my-service',
    data={'level': 'error'}
)`}
                          </pre>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  </Accordion>
                </Card>
              </TabsContent>

              {/* BILLING */}
              <TabsContent value="billing" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Billing & Usage Model</h2>

                  <div className="space-y-6">
                    <div className="bg-primary/10 p-6 rounded-lg border-l-4 border-primary">
                      <h3 className="font-semibold mb-2">💡 Pay-for-Anomalies Model</h3>
                      <ul className="space-y-2 text-sm">
                        <li><strong>Data Ingestion:</strong> Unlimited and FREE</li>
                        <li><strong>Anomaly Detection:</strong> Metered (value-add computation)</li>
                        <li><strong>1 "call"</strong> = 1 compression analysis for anomalies</li>
                        <li><strong>Pooled:</strong> Both APIs draw from shared pool</li>
                      </ul>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Plans</h3>
                      <div className="grid md:grid-cols-2 gap-4">
                        <Card className="p-6">
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="text-xl font-bold">Developer</h4>
                            <Badge>FREE</Badge>
                          </div>
                          <p className="text-3xl font-bold mb-4">$0<span className="text-sm font-normal">/mo</span></p>
                          <ul className="space-y-2 text-sm">
                            <li>✓ 10K anomaly calls</li>
                            <li>✓ Unlimited ingestion</li>
                            <li>✓ 7-day retention</li>
                          </ul>
                        </Card>

                        <Card className="p-6 border-2 border-primary relative">
                          <Badge className="absolute -top-3 right-4">POPULAR</Badge>
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="text-xl font-bold">Standard</h4>
                            <Badge variant="secondary">LAUNCH50</Badge>
                          </div>
                          <p className="text-3xl font-bold mb-4">$49<span className="text-sm font-normal">/mo</span></p>
                          <ul className="space-y-2 text-sm">
                            <li>✓ 250K anomaly calls</li>
                            <li>✓ Unlimited ingestion</li>
                            <li>✓ $0.0035/call overage</li>
                            <li>✓ 30-day retention</li>
                            <li>✓ Advanced features</li>
                          </ul>
                        </Card>

                        <Card className="p-6">
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="text-xl font-bold">Growth</h4>
                            <Badge variant="secondary">LAUNCH50</Badge>
                          </div>
                          <p className="text-3xl font-bold mb-4">$249<span className="text-sm font-normal">/mo</span></p>
                          <ul className="space-y-2 text-sm">
                            <li>✓ 2M anomaly calls</li>
                            <li>✓ Unlimited ingestion</li>
                            <li>✓ $0.0018/call overage</li>
                            <li>✓ 90-day retention</li>
                          </ul>
                        </Card>

                        <Card className="p-6 bg-gradient-to-br from-primary/5 to-primary/10">
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="text-xl font-bold">Enterprise</h4>
                            <Badge variant="outline">CUSTOM</Badge>
                          </div>
                          <p className="text-3xl font-bold mb-4">Custom</p>
                          <ul className="space-y-2 text-sm">
                            <li>✓ Custom commit</li>
                            <li>✓ ~$0.001/call</li>
                            <li>✓ Dedicated support</li>
                            <li>✓ Custom SLA</li>
                          </ul>
                        </Card>
                      </div>
                    </div>

                    <div className="bg-accent/20 p-6 rounded-lg border-l-4 border-accent">
                      <h3 className="font-semibold mb-2">🎉 Launch Promotion</h3>
                      <p className="text-sm">
                        LAUNCH50: 50% off first 3 months for Standard/Growth plans
                      </p>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-3">Cost Examples</h3>
                      <div className="space-y-3">
                        <Card className="p-4">
                          <h4 className="font-semibold mb-2">Example 1: Standard Plan</h4>
                          <pre className="bg-muted p-3 rounded text-sm">
{`Included: 250K calls
Used: 185K calls
Cost: $49 + $0 = $49`}
                          </pre>
                        </Card>

                        <Card className="p-4">
                          <h4 className="font-semibold mb-2">Example 2: With Overage</h4>
                          <pre className="bg-muted p-3 rounded text-sm">
{`Included: 250K calls
Used: 275K calls
Overage: 25K × $0.0035 = $87.50
Cost: $49 + $87.50 = $136.50`}
                          </pre>
                        </Card>
                      </div>
                    </div>
                  </div>
                </Card>
              </TabsContent>

              {/* SECURITY */}
              <TabsContent value="security" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Security & Compliance</h2>

                  <div className="space-y-6">
                    <div>
                      <h3 className="text-2xl font-semibold mb-3">Data Security</h3>
                      <ul className="space-y-2">
                        <li><strong>End-to-End Encryption:</strong> TLS 1.3 in transit, AES-256 at rest</li>
                        <li><strong>API Key Security:</strong> Bcrypt hashed before storage</li>
                        <li><strong>Multi-Tenant Isolation:</strong> Row-level security (RLS)</li>
                        <li><strong>Data Retention:</strong> Configurable 7-90 days</li>
                      </ul>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-4">Compliance</h3>
                      <Accordion type="single" collapsible className="space-y-2">
                        <AccordionItem value="gdpr">
                          <AccordionTrigger>GDPR Compliance</AccordionTrigger>
                          <AccordionContent>
                            <ul className="space-y-2 text-sm">
                              <li>✓ Right to Erasure (30 days)</li>
                              <li>✓ Data Portability (JSON export)</li>
                              <li>✓ Explicit Consent Management</li>
                              <li>✓ DPA Available</li>
                            </ul>
                          </AccordionContent>
                        </AccordionItem>

                        <AccordionItem value="soc2">
                          <AccordionTrigger>SOC2 Compliance</AccordionTrigger>
                          <AccordionContent>
                            <ul className="space-y-2 text-sm">
                              <li>✓ Security Controls</li>
                              <li>✓ 99.9% Availability SLA</li>
                              <li>✓ Processing Integrity</li>
                              <li>✓ Confidentiality & Privacy</li>
                            </ul>
                          </AccordionContent>
                        </AccordionItem>

                        <AccordionItem value="dora">
                          <AccordionTrigger>DORA (Digital Operational Resilience Act)</AccordionTrigger>
                          <AccordionContent>
                            <ul className="space-y-2 text-sm">
                              <li>✓ ICT Risk Management</li>
                              <li>✓ Incident Reporting</li>
                              <li>✓ Resilience Testing</li>
                              <li>✓ Third-Party Risk</li>
                            </ul>
                          </AccordionContent>
                        </AccordionItem>

                        <AccordionItem value="euai">
                          <AccordionTrigger>EU AI Act Compliance</AccordionTrigger>
                          <AccordionContent>
                            <ul className="space-y-2 text-sm">
                              <li>✓ Transparency (glass-box detection)</li>
                              <li>✓ Full Explainability</li>
                              <li>✓ Human Oversight</li>
                              <li>✓ Technical Documentation</li>
                            </ul>
                          </AccordionContent>
                        </AccordionItem>
                      </Accordion>
                    </div>

                    <div>
                      <h3 className="text-2xl font-semibold mb-3">Best Practices</h3>
                      <ol className="space-y-2 list-decimal list-inside">
                        <li>Rotate API keys every 90 days</li>
                        <li>Use environment variables (never hardcode)</li>
                        <li>Implement rate limiting</li>
                        <li>Monitor unusual usage patterns</li>
                        <li>Always use HTTPS</li>
                      </ol>
                    </div>
                  </div>
                </Card>
              </TabsContent>

              {/* ADVANCED */}
              <TabsContent value="advanced" className="space-y-6">
                <Card className="p-8">
                  <h2 className="text-3xl font-bold mb-6">Advanced Features</h2>

                  <Accordion type="single" collapsible className="space-y-4">
                    <AccordionItem value="prediction">
                      <AccordionTrigger className="text-xl font-semibold">
                        Anomaly Prediction
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <p className="text-muted-foreground">
                          Predict potential anomalies by analyzing compression trends and identifying degrading patterns.
                        </p>

                        <div>
                          <h4 className="font-semibold mb-2">Sensitivity Configuration</h4>
                          <pre className="bg-muted p-3 rounded text-sm">
{`{
  "sensitivity": 0.7  // 0.0-1.0
}`}
                          </pre>
                          <ul className="mt-2 space-y-1 text-sm">
                            <li><strong>Low (0.0-0.3):</strong> Severe anomalies only</li>
                            <li><strong>Medium (0.4-0.6):</strong> Balanced (recommended)</li>
                            <li><strong>High (0.7-1.0):</strong> Detect subtle changes</li>
                          </ul>
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="alerts">
                      <AccordionTrigger className="text-xl font-semibold">
                        Alert Configuration
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div className="grid md:grid-cols-2 gap-4">
                          {[
                            { icon: "📧", title: "Email", desc: "Full details via email" },
                            { icon: "📱", title: "SMS", desc: "Critical alerts via SMS" },
                            { icon: "🔗", title: "Webhooks", desc: "Custom endpoints" },
                            { icon: "💬", title: "Slack/Teams", desc: "Channel notifications" }
                          ].map((channel, idx) => (
                            <Card key={idx} className="p-4">
                              <h5 className="font-semibold mb-2">{channel.icon} {channel.title}</h5>
                              <p className="text-sm text-muted-foreground">{channel.desc}</p>
                            </Card>
                          ))}
                        </div>
                      </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="troubleshooting">
                      <AccordionTrigger className="text-xl font-semibold">
                        Troubleshooting Guide
                      </AccordionTrigger>
                      <AccordionContent className="space-y-4 pt-4">
                        <div className="space-y-3">
                          <Card className="p-4 border-l-4 border-destructive">
                            <h5 className="font-semibold mb-2">401 Unauthorized</h5>
                            <ul className="text-sm space-y-1 text-muted-foreground">
                              <li>• Verify API key format</li>
                              <li>• Check Bearer token header</li>
                              <li>• Ensure key not revoked</li>
                            </ul>
                          </Card>

                          <Card className="p-4 border-l-4 border-yellow-500">
                            <h5 className="font-semibold mb-2">429 Rate Limited</h5>
                            <ul className="text-sm space-y-1 text-muted-foreground">
                              <li>• Implement exponential backoff</li>
                              <li>• Check Retry-After header</li>
                              <li>• Consider upgrading plan</li>
                            </ul>
                          </Card>

                          <Card className="p-4 border-l-4 border-yellow-500">
                            <h5 className="font-semibold mb-2">402 Payment Required</h5>
                            <ul className="text-sm space-y-1 text-muted-foreground">
                              <li>• Usage limit exceeded</li>
                              <li>• Upgrade plan or wait for reset</li>
                              <li>• Review usage dashboard</li>
                            </ul>
                          </Card>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  </Accordion>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
