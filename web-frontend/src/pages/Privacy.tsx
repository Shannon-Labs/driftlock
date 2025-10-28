import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";

const Privacy = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-4xl">
          <h1 className="text-4xl font-bold mb-8">Privacy Policy</h1>
          <Card className="p-8 bg-gradient-card border-primary/10 prose prose-invert max-w-none">
            <p className="text-muted-foreground">Last updated: January 15, 2025</p>
            <p>Shannon Labs ("we", "us", or "our") operates driftlock.net. This page informs you of our policies regarding the collection, use, and disclosure of personal data.</p>
            
            <h2>Information Collection</h2>
            <p>We collect information you provide directly to us, including name, email address, and company information when you contact us or sign up for our services.</p>
            
            <h2>Data Usage</h2>
            <p>We use collected data to provide and improve our services, respond to inquiries, and send updates about Driftlock.</p>
            
            <h2>Data Security</h2>
            <p>We implement appropriate security measures to protect your personal information. All data is encrypted in transit and at rest.</p>
            
            <h2>Contact</h2>
            <p>For privacy concerns, contact us at <a href="mailto:hunter@shannonlabs.dev">hunter@shannonlabs.dev</a></p>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  );
};

export default Privacy;
