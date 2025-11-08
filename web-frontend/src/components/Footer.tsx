import { Link } from "react-router-dom";
import { Github, Twitter, Linkedin, Mail } from "lucide-react";

export const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t border-border bg-background/95 backdrop-blur">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          {/* Brand */}
          <div>
            <Link to="/" className="flex items-center space-x-2 mb-4">
              <div className="w-8 h-8 rounded-lg bg-gradient-primary"></div>
              <span className="text-xl font-bold">Driftlock</span>
            </Link>
            <p className="text-sm text-muted-foreground mb-4">
              Explainable anomaly detection for regulated industries
            </p>
            <div className="flex items-center gap-4">
              <a 
                href="https://github.com/Shannon-Labs/driftlock" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                aria-label="GitHub"
              >
                <Github className="w-5 h-5" />
              </a>
              <a 
                href="https://twitter.com/shannonlabs" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                aria-label="Twitter"
              >
                <Twitter className="w-5 h-5" />
              </a>
              <a 
                href="https://linkedin.com/company/shannonlabs" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                aria-label="LinkedIn"
              >
                <Linkedin className="w-5 h-5" />
              </a>
            </div>
          </div>

          {/* Product */}
          <div>
            <h3 className="font-semibold mb-4">Product</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li><Link to="/docs" className="hover:text-primary transition-colors">Documentation</Link></li>
              <li><Link to="/dashboard" className="hover:text-primary transition-colors">Dashboard</Link></li>
              <li><Link to="/integrations" className="hover:text-primary transition-colors">Integrations</Link></li>
              <li><Link to="/pricing" className="hover:text-primary transition-colors">Pricing</Link></li>
              <li><Link to="/changelog" className="hover:text-primary transition-colors">Changelog</Link></li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h3 className="font-semibold mb-4">Company</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li><Link to="/about" className="hover:text-primary transition-colors">About</Link></li>
              <li><Link to="/use-cases" className="hover:text-primary transition-colors">Use Cases</Link></li>
              <li><Link to="/blog" className="hover:text-primary transition-colors">Blog</Link></li>
              <li><Link to="/contact" className="hover:text-primary transition-colors">Contact</Link></li>
              <li><Link to="/security" className="hover:text-primary transition-colors">Security</Link></li>
            </ul>
          </div>

          {/* Resources */}
          <div>
            <h3 className="font-semibold mb-4">Resources</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li><Link to="/faq" className="hover:text-primary transition-colors">FAQ</Link></li>
              <li><Link to="/comparison" className="hover:text-primary transition-colors">vs. ML Detection</Link></li>
              <li><a href="https://github.com/Shannon-Labs/driftlock" target="_blank" rel="noopener noreferrer" className="hover:text-primary transition-colors">GitHub</a></li>
              <li><a href="mailto:hunter@shannonlabs.dev" className="hover:text-primary transition-colors flex items-center gap-1">
                <Mail className="w-3 h-3" />
                Contact Sales
              </a></li>
            </ul>
          </div>
        </div>

        {/* Trust & Compliance */}
        <div className="py-8 border-t border-border">
          <p className="text-sm text-muted-foreground text-center max-w-3xl mx-auto mb-3">
            Driftlock is engineered to meet the world's toughest data-protection standards.
          </p>
          <div className="flex flex-wrap justify-center gap-4 text-xs text-muted-foreground/80">
            <span className="px-3 py-1 rounded-full bg-primary/5 border border-primary/10">EU AI Act</span>
            <span className="px-3 py-1 rounded-full bg-primary/5 border border-primary/10">DORA</span>
            <span className="px-3 py-1 rounded-full bg-primary/5 border border-primary/10">NIS2</span>
            <span className="px-3 py-1 rounded-full bg-primary/5 border border-primary/10">CPRA</span>
            <span className="px-3 py-1 rounded-full bg-primary/5 border border-primary/10">SOC 2 Type II</span>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="pt-4 border-t border-border flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-muted-foreground">
          <p>
            © {currentYear} Shannon Labs. All rights reserved.
          </p>
          <div className="flex gap-6">
            <Link to="/privacy" className="hover:text-primary transition-colors">Privacy Policy</Link>
            <Link to="/terms" className="hover:text-primary transition-colors">Terms of Service</Link>
          </div>
        </div>
      </div>
    </footer>
  );
};
