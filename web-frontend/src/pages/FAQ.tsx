import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const FAQ = () => {
  const faqs = [
    {
      category: "General",
      questions: [
        {
          q: "What is Driftlock?",
          a: "Driftlock is a compression-based anomaly detection (CBAD) platform for OpenTelemetry data. It uses Meta's OpenZL format-aware compression framework to provide explainable, deterministic anomaly detection for regulated industries."
        },
        {
          q: "How is Driftlock different from traditional ML-based anomaly detection?",
          a: "Unlike black-box ML approaches, Driftlock provides glass-box explanations based on compression theory. Every anomaly comes with mathematical proof of why it was flagged, making it ideal for regulated industries that need to explain decisions to auditors."
        },
        {
          q: "What industries is Driftlock designed for?",
          a: "Driftlock is built for regulated industries including financial services (DORA compliance), healthcare (HIPAA), cybersecurity (NIS2), and AI/ML systems (EU AI Act). Any organization that needs explainable anomaly detection can benefit."
        },
      ]
    },
    {
      category: "Technical",
      questions: [
        {
          q: "How does compression-based anomaly detection work?",
          a: "CBAD measures the compressibility of data using algorithms like OpenZL. Normal data compresses predictably, while anomalous data compresses poorly due to structural differences. This is grounded in Kolmogorov complexity theory and provides statistical significance via permutation testing."
        },
        {
          q: "What data formats does Driftlock support?",
          a: "Driftlock natively supports OpenTelemetry Protocol (OTLP) for logs, metrics, and traces. We also support LLM I/O monitoring (prompts, responses, tool calls) and can ingest custom JSON formats."
        },
        {
          q: "What are the performance characteristics?",
          a: "Driftlock processes 100k+ events/second with <100ms p95 API latency. The CBAD engine achieves sub-second detection times with configurable memory usage and can scale horizontally."
        },
        {
          q: "Can I run Driftlock on-premises?",
          a: "Yes! Enterprise plans include on-premises deployment with Docker images and Kubernetes Helm charts. You maintain full control of your data."
        },
      ]
    },
    {
      category: "Compliance",
      questions: [
        {
          q: "What compliance frameworks does Driftlock support?",
          a: "Driftlock supports DORA (Digital Operational Resilience Act), NIS2 (Network and Information Security), EU AI Act, HIPAA, and provides tools for SOC 2 compliance. Evidence bundles and audit trails are built-in."
        },
        {
          q: "How does Driftlock help with DORA compliance?",
          a: "Driftlock generates cryptographically signed evidence bundles, provides incident detection and reporting templates, and maintains deterministic audit trails required for DORA compliance in financial services."
        },
        {
          q: "Can Driftlock monitor AI/ML systems for AI Act compliance?",
          a: "Yes, Driftlock includes LLM I/O monitoring capabilities to track prompts, responses, and tool calls. This enables runtime AI monitoring and explainable anomaly detection for high-risk AI systems under the EU AI Act."
        },
      ]
    },
    {
      category: "Pricing & Support",
      questions: [
        {
          q: "Is there a free tier?",
          a: "Yes! Our Starter plan is free and includes 10k events/day, 7-day retention, and basic anomaly detection. Perfect for testing and small projects."
        },
        {
          q: "What support options are available?",
          a: "Starter plans include community support. Professional plans get priority email support. Enterprise plans include 24/7 dedicated support with custom SLAs and professional services."
        },
        {
          q: "Can I upgrade or downgrade my plan?",
          a: "Yes, you can change plans at any time. Upgrades take effect immediately. Downgrades take effect at the end of your current billing period."
        },
        {
          q: "Do you offer custom pricing for high-volume deployments?",
          a: "Yes, Enterprise plans include custom pricing based on your specific needs. Contact our sales team to discuss volume discounts and custom SLAs."
        },
      ]
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-4xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              <span className="text-gradient">Frequently Asked Questions</span>
            </h1>
            <p className="text-xl text-muted-foreground">
              Everything you need to know about Driftlock
            </p>
          </div>

          {/* FAQ Categories */}
          <div className="space-y-8">
            {faqs.map((category, idx) => (
              <div key={idx}>
                <h2 className="text-2xl font-bold mb-4">{category.category}</h2>
                <Card className="p-6 bg-gradient-card border-primary/10">
                  <Accordion type="single" collapsible className="w-full">
                    {category.questions.map((faq, faqIdx) => (
                      <AccordionItem key={faqIdx} value={`item-${idx}-${faqIdx}`}>
                        <AccordionTrigger className="text-left font-semibold">
                          {faq.q}
                        </AccordionTrigger>
                        <AccordionContent className="text-muted-foreground">
                          {faq.a}
                        </AccordionContent>
                      </AccordionItem>
                    ))}
                  </Accordion>
                </Card>
              </div>
            ))}
          </div>

          {/* Contact */}
          <Card className="mt-12 p-8 bg-gradient-card border-primary/10 text-center">
            <h3 className="text-xl font-bold mb-3">Still have questions?</h3>
            <p className="text-muted-foreground mb-6">
              Can't find the answer you're looking for? Contact our support team.
            </p>
            <a href="/contact">
              <button className="px-6 py-3 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                Contact Support
              </button>
            </a>
          </Card>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default FAQ;
