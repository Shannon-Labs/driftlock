import { Navigation } from "@/components/Navigation";
import { Footer } from "@/components/Footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Shield, Lock, FileCheck, Eye, Server, CheckCircle2 } from "lucide-react";

const Security = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      
      <main className="pt-24 pb-16 px-4">
        <div className="container mx-auto max-w-6xl">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Security & <span className="text-gradient">Trust Center</span>
            </h1>
            <p className="text-xl text-muted-foreground">
              Enterprise-grade security built for regulated industries
            </p>
          </div>

          {/* Compliance Badges */}
          <div className="flex flex-wrap justify-center gap-4 mb-16">
            <Badge variant="secondary" className="px-4 py-2 text-sm bg-primary/20 border-primary/30">
              DORA Compliant
            </Badge>
            <Badge variant="secondary" className="px-4 py-2 text-sm bg-secondary/20 border-secondary/30">
              NIS2 Ready
            </Badge>
            <Badge variant="secondary" className="px-4 py-2 text-sm bg-accent/20 border-accent/30">
              AI Act Compatible
            </Badge>
            <Badge variant="secondary" className="px-4 py-2 text-sm bg-primary/20 border-primary/30">
              SOC 2 Type II (In Progress)
            </Badge>
          </div>

          {/* Security Features */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-16">
            <Card className="p-6 bg-gradient-card border-primary/10">
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Lock className="w-6 h-6 text-primary" />
              </div>
              <h3 className="text-xl font-bold mb-3">Data Encryption</h3>
              <p className="text-muted-foreground mb-4">
                End-to-end encryption for data in transit and at rest using industry-standard AES-256 and TLS 1.3.
              </p>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  TLS 1.3 for all API communications
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  AES-256 encryption at rest
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  Key rotation and management
                </li>
              </ul>
            </Card>

            <Card className="p-6 bg-gradient-card border-primary/10">
              <div className="w-12 h-12 rounded-lg bg-secondary/10 flex items-center justify-center mb-4">
                <Eye className="w-6 h-6 text-secondary" />
              </div>
              <h3 className="text-xl font-bold mb-3">Privacy Controls</h3>
              <p className="text-muted-foreground mb-4">
                Configurable data redaction and privacy-preserving anomaly detection for sensitive information.
              </p>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-secondary" />
                  PII redaction support
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-secondary" />
                  GDPR compliance tools
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-secondary" />
                  Data residency controls
                </li>
              </ul>
            </Card>

            <Card className="p-6 bg-gradient-card border-primary/10">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mb-4">
                <FileCheck className="w-6 h-6 text-accent" />
              </div>
              <h3 className="text-xl font-bold mb-3">Audit Trails</h3>
              <p className="text-muted-foreground mb-4">
                Comprehensive audit logging with cryptographic integrity for regulatory compliance.
              </p>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-accent" />
                  Cryptographically signed logs
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-accent" />
                  Immutable audit trails
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-accent" />
                  Evidence bundle export
                </li>
              </ul>
            </Card>

            <Card className="p-6 bg-gradient-card border-primary/10">
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Server className="w-6 h-6 text-primary" />
              </div>
              <h3 className="text-xl font-bold mb-3">Infrastructure Security</h3>
              <p className="text-muted-foreground mb-4">
                Secure deployment options including on-premises, private cloud, and multi-region support.
              </p>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  On-premises deployment
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  VPC/private network support
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  Multi-region redundancy
                </li>
              </ul>
            </Card>
          </div>

          {/* Compliance Section */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 mb-12">
            <h2 className="text-3xl font-bold mb-6">Regulatory Compliance</h2>
            <div className="space-y-6">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <Shield className="w-5 h-5 text-primary" />
                  <h3 className="text-xl font-semibold">DORA (Digital Operational Resilience Act)</h3>
                </div>
                <p className="text-muted-foreground pl-8">
                  Built-in evidence bundle generation, incident reporting templates, and operational resilience 
                  testing capabilities for financial institutions operating in the EU.
                </p>
              </div>

              <div>
                <div className="flex items-center gap-3 mb-2">
                  <Shield className="w-5 h-5 text-secondary" />
                  <h3 className="text-xl font-semibold">NIS2 (Network and Information Security Directive)</h3>
                </div>
                <p className="text-muted-foreground pl-8">
                  Automated incident detection, reporting workflows, and security event correlation to meet 
                  NIS2 requirements for critical infrastructure operators.
                </p>
              </div>

              <div>
                <div className="flex items-center gap-3 mb-2">
                  <Shield className="w-5 h-5 text-accent" />
                  <h3 className="text-xl font-semibold">EU AI Act</h3>
                </div>
                <p className="text-muted-foreground pl-8">
                  Runtime AI monitoring capabilities for LLM systems, including prompt/response tracking, 
                  model drift detection, and explainable anomaly reports.
                </p>
              </div>
            </div>
          </Card>

          {/* Security Practices */}
          <Card className="p-8 md:p-12 bg-gradient-card border-primary/10 mb-12">
            <h2 className="text-3xl font-bold mb-6">Security Practices</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div>
                <h3 className="font-semibold mb-3">Development Security</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Secure SDLC with automated security scanning</li>
                  <li>• Dependency vulnerability monitoring</li>
                  <li>• Code review and static analysis</li>
                  <li>• Regular penetration testing</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">Operational Security</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• 24/7 security monitoring</li>
                  <li>• Incident response procedures</li>
                  <li>• Regular security training</li>
                  <li>• Third-party security audits</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">Access Control</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Role-based access control (RBAC)</li>
                  <li>• Multi-factor authentication (MFA)</li>
                  <li>• Principle of least privilege</li>
                  <li>• Access logging and monitoring</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">Data Protection</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Automated backup procedures</li>
                  <li>• Disaster recovery planning</li>
                  <li>• Data retention policies</li>
                  <li>• Secure data disposal</li>
                </ul>
              </div>
            </div>
          </Card>

          {/* Contact Security Team */}
          <Card className="p-8 bg-gradient-card border-primary/10 text-center">
            <h2 className="text-2xl font-bold mb-4">Security Questions?</h2>
            <p className="text-muted-foreground mb-6">
              Our security team is here to answer your questions and provide additional documentation
            </p>
            <a href="mailto:hunter@shannonlabs.dev?subject=Security%20Inquiry">
              <button className="px-6 py-3 bg-gradient-primary text-primary-foreground rounded-lg font-medium hover:opacity-90 transition-opacity">
                Contact Security Team
              </button>
            </a>
          </Card>
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Security;
