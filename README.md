Driftlock: Explainable AI anomaly detection for financial institutions facing €50M DORA fines

⚠️ STATUS: Research prototype → Enterprise pilot (3 banks in Q1 2025 evaluation)
🤖 AI CONTRIBUTION: Claude/Gemini assisted 40% of codebase; human-audited by founder (JD/MBA)
🎯 PROBLEM: EU regulators now require glass-box AI audits; Dynatrace/New Relic can't explain anomalies
💰 SOLUTION: Compression-based detection (CBAD) that runs on-prem with 0 API fees, full audit trails

git clone https://github.com/Shannon-Labs/driftlock.git
cd Driftlock
./start.sh  # Live demo: https://demo.driftlock.com (API key: yc-demo-2025)Driftlock provides compression-based anomaly detection (CBAD) with glass-box explainability, designed specifically for compliance with DORA, NIS2, and EU AI Act regulations.

**🎉 Now available as a pure open source release!** Driftlock runs without external dependencies - just Docker, PostgreSQL, and your API key.

**🌐 Live Demo:** https://9aac0d30.driftlock.pages.dev

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/shannon-labs/driftlock.git
cd Driftlock

# Copy the sample environment and configure it
cp .env.example .env
# Edit .env and set your API key and database password

# Launch everything with one command
./start.sh

# Access your dashboard
open http://localhost:3000
# Login with your API key from .env
```

**Or use Make commands:**
```bash
make quick-start  # Sets up .env and starts services
make dev          # Start development environment
make stop         # Stop all services
```

> Detailed installation options (Docker, bare metal, Kubernetes) live in [docs/installation.md](docs/installation.md).

### Wire Driftlock into OpenTelemetry

Add the Driftlock processor to your OpenTelemetry configuration:

```yaml
processors:
  Driftlock/anomaly:
    thresholds:
      compression_ratio: 0.7
      ncd_threshold: 0.3
    explanation:
      enabled: true
      detail_level: "detailed"

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [Driftlock/anomaly, batch]
      exporters: [otlp]
```

## ✨ Features

- **🔍 Glass-Box Anomaly Detection**: Every anomaly comes with human-readable explanations
- **📊 Regulatory Compliance**: Built-in audit trails for DORA/NIS2/EU AI Act
- **⚡ High Performance**: Rust core with Go API server
- **🔧 OpenTelemetry Native**: Seamless integration with existing observability stacks
- **🌐 Pure Open Source**: Apache 2.0 licensed, self-hosted, no external dependencies
- **📈 Real-time Monitoring**: Live dashboard with anomaly streaming
- **🔐 Simple Authentication**: API key-based access control
- **🐳 Docker Ready**: One-command deployment with Docker Compose

## 📊 Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Data Sources  │───▶│ OpenTelemetry    │───▶│ Driftlock Core  │
│ (Logs, Metrics, │    │ Collector        │    │ (Rust CBAD)     │
│  Traces)        │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Dashboard & API │◀───│   API Server     │◀───│  Explanations   │
│ Access          │    │    (Go)          │    │   & Audit Trail │
│ (API Key Auth)  │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🏢 Compliance Consulting Services

Need professional help with regulatory compliance? **[Shannon Labs](https://shannonlabs.ai)** offers expert consulting services:

- **📋 DORA Compliance Implementation**: Full compliance setup and documentation
- **🔍 NIS2 Incident Reporting**: Template-based reporting procedures
- **🤖 EU AI Act Audit Trails**: Complete transparency documentation
- **🎓 Regulatory Training**: Team education and compliance workshops
- **📊 Custom Compliance Frameworks**: Tailored solutions for your industry

**Why choose our consulting?**
- Former regulators and compliance officers on staff
- PhD-level technical expertise in AI systems
- Proven track record with financial institutions
- Fixed-price packages with clear deliverables

*Open source Driftlock provides the foundation. Our consulting helps you achieve full regulatory compliance with confidence.*

### Popular Consulting Packages

| Package | Duration | Investment | Best For |
|---------|----------|------------|----------|
| **DORA Quick Start** | 4-6 weeks | $5,000-10,000 | Financial institutions |
| **AI Act Readiness** | 2-4 weeks | $2,500-5,000 | High-risk AI systems |
| **Explainable AI Setup** | 6-10 weeks | $10,000-25,000 | Organizations needing transparency |
| **Compliance-as-a-Service** | Ongoing | $3,000-8,000/month | Continuous compliance support |

📧 **Contact us:** consulting@shannonlabs.ai for a free initial consultation

## 📖 Documentation

- [Installation Guide](docs/installation.md)
- [Architecture Overview](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Compliance Integration](docs/COMPLIANCE_DORA.md)
- [Examples](examples/)
- [Security Guidance](SECURITY.md)

## 🛠️ Development

### Prerequisites

- Go 1.24+
- Rust stable (1.75+ recommended)
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 14+

### Local Development

```bash
# Quick start (setup + start services)
make quick-start

# Or manually:
make setup      # Install dependencies
make test       # Run all tests
make dev        # Start development environment
make build      # Build production artifacts
```

### Project Structure

```
Driftlock/
├── api-server/        # Go HTTP API + ingest pipeline
├── cbad-core/         # Rust compression-based anomaly detector
├── web-frontend/      # React dashboard & marketing site
├── docs/              # Product + process documentation
├── deploy/ & k8s/     # Docker Compose, Helm, Kubernetes manifests
└── tools/, scripts/   # Developer utilities
```

## 🧪 Testing

```bash
# Quick sweep
make test

# Individual components
cd cbad-core && cargo test
cd api-server && go test ./...
cd web-frontend && npm run build

# Integration smoke test (requires Postgres)
make test-integration

# CI validation (lint + test + security scan)
make ci-check
```

## 🚀 Deployment

### Docker

```bash
# API server
docker build -t Driftlock/api-server .

# Dashboard (serve the static build)
pnpm --dir web-frontend build
docker build -f deploy/docker/Dockerfile.web -t Driftlock/web-frontend .
```

### Kubernetes / Helm

```
kubectl apply -f k8s/
# or
helm install Driftlock ./helm/Driftlock
```

## 📊 Compliance & Security

- **Explainable AI**: Every anomaly includes mathematical explanations
- **Audit Trails**: Complete logging for regulatory compliance
- **Data Privacy**: GDPR-compliant data handling
- **Security**: Built-in authentication and encryption
- **Transparency**: Open source with clear documentation

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines and local setup steps.

### Ways to Contribute

- 🐛 Report bugs
- 💡 Suggest features
- 📝 Improve documentation
- 🔧 Submit pull requests
- 🧪 Write tests
- 🌍 Translate documentation

## 📄 License

Licensed under **AGPL-3.0** with patent protections. See [LICENSE](LICENSE) and [PATENTS.md](PATENTS.md) for details.

### Commercial Licensing

For commercial use, proprietary licenses are available through Shannon Labs. This dual-licensing model allows you to:
- Use Driftlock under AGPL-3.0 for open source projects
- Purchase a commercial license for proprietary applications

Contact: licensing@shannonlabs.ai

## 🙏 Acknowledgments

- OpenTelemetry community for the observability framework
- Compression-based anomaly detection research community
- Regulatory compliance experts who provided insights
- All our amazing contributors

## 📞 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/Shannon-Labs/Driftlock/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Shannon-Labs/Driftlock/discussions)
- **Security**: security@shannonlabs.ai
- **Enterprise**: contact@shannonlabs.ai

---

**Built by [Shannon Labs](https://shannonlabs.ai)** - Making AI explainable and compliant.

If you find Driftlock useful, please give us a ⭐ on GitHub!
