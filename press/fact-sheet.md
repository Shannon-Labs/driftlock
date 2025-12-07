# Driftlock Fact Sheet

## Product Overview

**Name**: Driftlock
**Tagline**: Mathematical Observability
**Category**: Anomaly Detection / Observability / AIOps
**Launch Date**: [To be confirmed]

## Key Features

### Core Capabilities
- **Compression-Based Anomaly Detection (CBAD)**: Novel mathematical algorithm
- **Real-time Processing**: Detect anomalies as they occur
- **Zero Training Time**: No model training required
- **Explainable AI**: Every anomaly includes mathematical proof
- **Multi-Modal**: Logs, metrics, traces, and LLM I/O
- **OpenTelemetry Native**: Seamless integration

### Technical Highlights
- **O(1) Complexity**: Constant time per event
- **No ML Training**: Start detecting immediately
- **Mathematical Proofs**: NCD values, p-values, confidence scores
- **Privacy-First**: Optional zero-knowledge mode
- **Enterprise Ready**: SSO, audit trails, compliance reports

## Target Markets

### Primary Industries
- **Financial Services**: DORA compliance, transaction monitoring
- **Cybersecurity**: API security, intrusion detection
- **AI/ML Companies**: LLM monitoring, prompt injection detection
- **IoT/Manufacturing**: Equipment monitoring, predictive maintenance

### Compliance Frameworks
- **EU DORA**: Digital Operational Resilience Act
- **EU NIS2**: Network and Information Security Directive
- **EU AI Act**: AI transparency and explainability
- **GDPR**: Data protection and privacy

## Pricing Model

### Self-Service Tiers
- **Free**: $0/month (10k events/month)
- **Standard**: $15/month (500k events/month)
- **Pro**: $100/month (5M events/month)
- **Enterprise**: $299/month (25M events/month)

### Enterprise Features
- Custom volume pricing
- Dedicated support
- SLA guarantees
- On-premise deployment option

## Technical Specifications

### Infrastructure
- **Backend**: Go 1.22+
- **Core Algorithm**: Rust (FFI integration)
- **Frontend**: Vue 3 + TypeScript
- **Database**: PostgreSQL + Redis
- **Deployment**: GCP Cloud Run, GKE

### Performance
- **Latency**: <10ms per event
- **Throughput**: 10k+ events/second
- **Storage**: Compressed baselines (minimal overhead)
- **Availability**: 99.9% SLA

### Security
- **Authentication**: API Keys, JWT, SSO
- **Authorization**: Role-based access control
- **Encryption**: AES-256 at rest, TLS 1.3 in transit
- **Compliance**: SOC 2 Type II (in progress)

## Integrations

### Observability Stack
- **OpenTelemetry**: Native support
- **Prometheus**: Metrics integration
- **Grafana**: Visualization dashboards
- **ELK Stack**: Log analytics

### CI/CD & DevOps
- **GitHub Actions**: Workflow integration
- **GitLab CI**: Pipeline support
- **Jenkins**: Plugin available
- **Kubernetes**: Operator ready

## Company Information

**Name**: Shannon Labs, LLC
**Founded**: [To be confirmed]
**Headquarters**: Remote-first
**CEO**: Hunter Bown
**Team**: Distributed engineering team

### Contact
- **Website**: https://driftlock.net
- **Email**: hunter@shannonlabs.dev
- **GitHub**: https://github.com/Shannon-Labs/driftlock

## Resources

### Documentation
- Available at https://driftlock.net/docs
- API Reference: https://driftlock.net/docs/api
- User Guide: https://driftlock.net/docs/user-guide

### Community
- Discord: [Check website for invite]
- [Additional community links to be added]

---

**© 2025 Shannon Labs, LLC. All rights reserved.**