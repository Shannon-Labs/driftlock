import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ArrowRight, Sparkles } from "lucide-react";
import { Link } from "react-router-dom";
import { TerminalDemo } from "./TerminalDemo";

export const Hero = () => {
  return (
    <section className="relative min-h-screen flex items-center overflow-hidden">
      {/* Animated Mesh Background */}
      <div className="absolute inset-0 mesh-gradient"></div>
      <div className="absolute inset-0 bg-gradient-to-b from-background via-transparent to-background"></div>

      {/* Ambient background glow - removed floating orbs */}

      <div className="container mx-auto px-4 pt-32 pb-20 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-8 animate-fade-in">
            <Badge variant="secondary" className="glass-card px-4 py-2 w-fit">
              <Sparkles className="w-4 h-4 mr-2" />
              Powered by Meta's OpenZL Framework
            </Badge>

            <h1 className="text-5xl md:text-6xl lg:text-7xl font-display font-bold leading-tight">
              When your data stops compressing,{" "}
              <span className="text-gradient">we tell you why</span>
            </h1>

            <p className="text-xl md:text-2xl text-muted-foreground font-light leading-relaxed max-w-2xl">
              Format-aware compression anomaly detection — zero data retention, built for compliance. 
              No machine learning, no training, no black boxes.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 pt-4">
              <Link to="/docs">
                <Button size="lg" className="bg-gradient-primary hover:opacity-90 transition-opacity text-lg px-8 group">
                  Start Free (10k calls)
                  <ArrowRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
              <Link to="/docs">
                <Button size="lg" variant="outline" className="text-lg px-8 glass-card hover-lift border-primary/20">
                  See API Docs
                </Button>
              </Link>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 gap-8 pt-8">
              <div className="space-y-1">
                <div className="text-4xl font-display font-bold text-gradient">100k+</div>
                <div className="text-sm text-muted-foreground">Events/sec</div>
              </div>
              <div className="space-y-1">
                <div className="text-4xl font-display font-bold text-gradient">&lt;100ms</div>
                <div className="text-sm text-muted-foreground">Latency</div>
              </div>
            </div>
          </div>

          {/* Right Visual - Terminal Demo */}
          <div className="relative hidden lg:block animate-fade-in" style={{ animationDelay: '200ms' }}>
            <div className="relative">
              {/* Glow Effect */}
              <div className="absolute inset-0 bg-gradient-primary opacity-10 blur-[80px] rounded-3xl animate-glow pointer-events-none"></div>
              
              {/* Terminal Demo */}
              <TerminalDemo />
              
              {/* Single Floating Stat */}
              <div className="absolute -bottom-6 -right-6 glass-card rounded-2xl p-4 hover-lift">
                <div className="text-sm text-muted-foreground mb-1">Detection Accuracy</div>
                <div className="text-2xl font-display font-bold text-gradient">99.2%</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Scroll Indicator */}
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 animate-bounce">
        <div className="w-6 h-10 border-2 border-primary/30 rounded-full flex justify-center pt-2">
          <div className="w-1.5 h-2 bg-primary rounded-full animate-pulse"></div>
        </div>
      </div>
    </section>
  );
};
