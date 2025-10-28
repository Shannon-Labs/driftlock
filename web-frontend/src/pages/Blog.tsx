import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Calendar, ArrowRight } from "lucide-react";

const Blog = () => {
  const posts = [
    {
      title: "Introducing Driftlock: Explainable Anomaly Detection",
      excerpt: "Why we built a compression-based anomaly detection platform and how it differs from traditional ML approaches.",
      date: "2025-01-15",
      category: "Product",
      readTime: "5 min read",
    },
    {
      title: "Understanding Compression-Based Anomaly Detection (CBAD)",
      excerpt: "A deep dive into the mathematical foundations of CBAD and why Kolmogorov complexity matters for explainability.",
      date: "2025-01-10",
      category: "Technical",
      readTime: "12 min read",
    },
    {
      title: "Format-Aware Compression with OpenZL",
      excerpt: "How Meta's OpenZL framework enables better anomaly detection through structural understanding of telemetry data.",
      date: "2025-01-05",
      category: "Technical",
      readTime: "8 min read",
    },
    {
      title: "DORA Compliance: What You Need to Know",
      excerpt: "Breaking down the Digital Operational Resilience Act and how Driftlock helps financial institutions meet requirements.",
      date: "2024-12-20",
      category: "Compliance",
      readTime: "6 min read",
    },
    {
      title: "Runtime AI Monitoring for the EU AI Act",
      excerpt: "How to monitor LLM systems for compliance with the EU AI Act using compression-based detection.",
      date: "2024-12-15",
      category: "AI/ML",
      readTime: "10 min read",
    },
    {
      title: "Deterministic vs. Probabilistic Anomaly Detection",
      excerpt: "Why determinism matters in regulated industries and how Driftlock achieves 100% reproducible results.",
      date: "2024-12-10",
      category: "Technical",
      readTime: "7 min read",
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-6xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              <span className="text-gradient">Blog</span>
            </h1>
            <p className="text-xl text-muted-foreground">
              Insights on anomaly detection, compliance, and observability
            </p>
          </div>

          {/* Featured Post */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all mb-12 cursor-pointer">
            <Badge variant="secondary" className="mb-4">Featured</Badge>
            <h2 className="text-3xl font-bold mb-4">{posts[0].title}</h2>
            <p className="text-muted-foreground mb-6 text-lg">{posts[0].excerpt}</p>
            <div className="flex items-center gap-6 text-sm text-muted-foreground mb-6">
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4" />
                {posts[0].date}
              </div>
              <Badge variant="secondary">{posts[0].category}</Badge>
              <span>{posts[0].readTime}</span>
            </div>
            <button className="text-primary font-medium flex items-center gap-2 hover:gap-3 transition-all">
              Read More <ArrowRight className="w-4 h-4" />
            </button>
          </Card>

          {/* Blog Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {posts.slice(1).map((post, idx) => (
              <Card key={idx} className="p-6 bg-gradient-card border-primary/10 hover:border-primary/30 transition-all cursor-pointer">
                <Badge variant="secondary" className="mb-3">{post.category}</Badge>
                <h3 className="text-xl font-bold mb-3">{post.title}</h3>
                <p className="text-sm text-muted-foreground mb-4 line-clamp-2">{post.excerpt}</p>
                <div className="flex items-center gap-4 text-xs text-muted-foreground mb-4">
                  <div className="flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    {post.date}
                  </div>
                  <span>{post.readTime}</span>
                </div>
                <button className="text-primary text-sm font-medium flex items-center gap-2 hover:gap-3 transition-all">
                  Read More <ArrowRight className="w-3 h-3" />
                </button>
              </Card>
            ))}
          </div>

          {/* Newsletter */}
          <Card className="mt-16 p-8 md:p-12 bg-gradient-card border-primary/10 text-center">
            <h2 className="text-2xl font-bold mb-4">Stay Updated</h2>
            <p className="text-muted-foreground mb-6">
              Get the latest on anomaly detection, compliance, and observability delivered to your inbox
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

export default Blog;
