import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";

const Terms = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-4xl">
          <h1 className="text-4xl font-bold mb-8">Terms of Service</h1>
          <Card className="p-8 bg-gradient-card border-primary/10 prose prose-invert max-w-none">
            <p className="text-muted-foreground">Last updated: January 15, 2025</p>
            
            <h2>Agreement to Terms</h2>
            <p>By accessing Driftlock, you agree to be bound by these Terms of Service.</p>
            
            <h2>Use License</h2>
            <p>Permission is granted to use Driftlock for commercial and non-commercial purposes subject to these terms.</p>
            
            <h2>Disclaimer</h2>
            <p>Driftlock is provided "as is" without warranties of any kind. We do not guarantee uninterrupted or error-free service.</p>
            
            <h2>Limitations</h2>
            <p>Shannon Labs shall not be liable for any damages arising from the use or inability to use Driftlock.</p>
            
            <h2>Contact</h2>
            <p>Questions about Terms of Service? Email <a href="mailto:hunter@shannonlabs.dev">hunter@shannonlabs.dev</a></p>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  );
};

export default Terms;
