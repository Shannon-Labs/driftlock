import { Button } from "@/components/ui/button";
import { ArrowRight, Shield, Lock, Zap, Sparkles } from "lucide-react";
import { Link } from "react-router-dom";

export const CallToAction = () => {
  return (
    <section className="py-32 px-4 relative overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 mesh-gradient"></div>
        <div className="absolute inset-0 bg-gradient-to-b from-background via-background/80 to-background"></div>
      </div>

      {/* Floating Elements */}
      <div className="absolute top-20 left-20 w-64 h-64 bg-primary/10 rounded-full blur-3xl animate-float"></div>
      <div className="absolute bottom-20 right-20 w-80 h-80 bg-secondary/10 rounded-full blur-3xl animate-float" style={{ animationDelay: '3s' }}></div>

      <div className="container mx-auto relative z-10">
        <div className="max-w-4xl mx-auto">
          {/* Main Card */}
          <div className="glass-card rounded-3xl p-12 md:p-16 text-center space-y-8 hover-lift">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20">
              <Sparkles className="w-4 h-4 text-primary" />
              <span className="text-sm font-medium">Start Your Free Trial</span>
            </div>

            {/* Heading */}
            <h2 className="text-4xl md:text-6xl font-display font-bold">
              Ready to deploy{" "}
              <span className="text-gradient">explainable</span>{" "}
              anomaly detection?
            </h2>

            {/* Description */}
            <p className="text-xl text-muted-foreground font-light max-w-2xl mx-auto">
              Join teams using Driftlock to detect anomalies in production with glass-box explanations 
              and compliance-ready audit trails.
            </p>

            {/* CTAs */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
              <Link to="/docs">
                <Button size="lg" className="bg-gradient-primary hover:opacity-90 transition-opacity text-lg px-10 group">
                  Get Started Free
                  <ArrowRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
              <Link to="/contact">
                <Button size="lg" variant="outline" className="text-lg px-10 glass-card hover-lift border-primary/20">
                  Contact Sales
                </Button>
              </Link>
            </div>

            {/* Trust Indicators */}
            <div className="flex flex-wrap justify-center gap-8 pt-8 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <Shield className="w-4 h-4 text-primary" />
                <span>DORA Compliant</span>
              </div>
              <div className="flex items-center gap-2">
                <Lock className="w-4 h-4 text-primary" />
                <span>NIS2 Ready</span>
              </div>
              <div className="flex items-center gap-2">
                <Zap className="w-4 h-4 text-primary" />
                <span>Production-Ready</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
