import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Calendar } from "lucide-react";

const Changelog = () => {
  const releases = [
    {
      version: "0.2.0",
      date: "2025-01-15",
      type: "Feature Release",
      changes: [
        { type: "feature", text: "Enhanced Go FFI bridge with lifecycle management" },
        { type: "feature", text: "Streaming anomaly detection interface" },
        { type: "feature", text: "Configurable privacy redaction support" },
        { type: "improvement", text: "Performance optimizations: 1000+ events/second" },
        { type: "improvement", text: "Thread-safe operations with mutex protection" },
      ]
    },
    {
      version: "0.1.5",
      date: "2025-01-10",
      type: "Bug Fix",
      changes: [
        { type: "fix", text: "Fixed memory leak in compression adapter" },
        { type: "fix", text: "Resolved race condition in event processing" },
        { type: "improvement", text: "Improved error messages for invalid configurations" },
      ]
    },
    {
      version: "0.1.0",
      date: "2025-01-05",
      type: "Initial Release",
      changes: [
        { type: "feature", text: "Core CBAD engine with Rust implementation" },
        { type: "feature", text: "OpenZL format-aware compression integration" },
        { type: "feature", text: "OpenTelemetry Collector processor" },
        { type: "feature", text: "Go API server with PostgreSQL backend" },
        { type: "feature", text: "Real-time anomaly detection" },
        { type: "feature", text: "Statistical significance testing" },
        { type: "feature", text: "Evidence bundle export (JSON)" },
      ]
    },
  ];

  const getTypeColor = (type: string) => {
    switch (type) {
      case "feature":
        return "bg-green-500/20 text-green-500 border-green-500/30";
      case "improvement":
        return "bg-blue-500/20 text-blue-500 border-blue-500/30";
      case "fix":
        return "bg-yellow-500/20 text-yellow-500 border-yellow-500/30";
      default:
        return "bg-primary/20 text-primary border-primary/30";
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-4xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              <span className="text-gradient">Changelog</span>
            </h1>
            <p className="text-xl text-muted-foreground">
              Product updates and release notes
            </p>
          </div>

          {/* Releases */}
          <div className="space-y-8">
            {releases.map((release, idx) => (
              <Card key={idx} className="p-8 bg-gradient-card border-primary/10">
                <div className="flex items-start justify-between mb-6">
                  <div>
                    <h2 className="text-2xl font-bold mb-2">Version {release.version}</h2>
                    <Badge variant="secondary" className="bg-primary/20 border-primary/30">
                      {release.type}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Calendar className="w-4 h-4" />
                    {release.date}
                  </div>
                </div>

                <div className="space-y-3">
                  {release.changes.map((change, changeIdx) => (
                    <div key={changeIdx} className="flex items-start gap-3">
                      <Badge 
                        variant="secondary" 
                        className={`${getTypeColor(change.type)} flex-shrink-0 mt-0.5`}
                      >
                        {change.type}
                      </Badge>
                      <p className="text-muted-foreground">{change.text}</p>
                    </div>
                  ))}
                </div>
              </Card>
            ))}
          </div>

          {/* Subscribe */}
          <Card className="mt-12 p-8 bg-gradient-card border-primary/10 text-center">
            <h3 className="text-xl font-bold mb-3">Stay Updated</h3>
            <p className="text-muted-foreground mb-6">
              Subscribe to release notifications and get updates delivered to your inbox
            </p>
            <div className="flex flex-col sm:flex-row gap-4 max-w-md mx-auto">
              <input
                type="email"
                placeholder="your.email@company.com"
                className="flex-1 px-4 py-2 rounded-lg bg-background border border-border focus:outline-none focus:ring-2 focus:ring-primary"
              />
              <button className="px-6 py-2 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                Subscribe
              </button>
            </div>
          </Card>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Changelog;
