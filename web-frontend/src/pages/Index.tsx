import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Hero } from "@/components/Hero";
import { ProblemSection } from "@/components/ProblemSection";
import { HowItWorks } from "@/components/HowItWorks";
import { Features } from "@/components/Features";
import { WhyDriftlock } from "@/components/WhyDriftlock";
import { UseCasesSection } from "@/components/UseCasesSection";
import { DeveloperOnboarding } from "@/components/DeveloperOnboarding";
import { TrustSection } from "@/components/TrustSection";
import { CallToAction } from "@/components/CallToAction";

const Index = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <Hero />
      <ProblemSection />
      <HowItWorks />
      <Features />
      <WhyDriftlock />
      <UseCasesSection />
      <DeveloperOnboarding />
      <TrustSection />
      <CallToAction />
      <Footer />
    </div>
  );
};

export default Index;
