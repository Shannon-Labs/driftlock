# DriftLock

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/shannon-labs/driftlock)](https://goreportcard.com/report/github.com/shannon-labs/driftlock)
[![Rust](https://img.shields.io/badge/rust-%23000000.svg?style=for-the-badge&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![CI](https://github.com/shannon-labs/driftlock/workflows/CI/badge.svg)](https://github.com/shannon-labs/driftlock/actions)

> **Explainable AI anomaly detection for regulated industries**

DriftLock provides compression-based anomaly detection (CBAD) with glass-box explainability, designed specifically for compliance with DORA, NIS2, and EU AI Act regulations.

**Now available as a standalone open source release!** DriftLock runs without external dependencies - just PostgreSQL and your API key.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/shannon-labs/driftlock.git
cd driftlock

# Copy the sample environment (never commit secrets)
cp .env.example .env

# Set your API key for dashboard access (optional, for development)
echo "DRIFTLOCK_DEV_API_KEY=your-secret-key-here" >> .env

# Launch Postgres + API + dashboard
docker compose up -d

# Web dashboard & API
open http://localhost:3000
curl  http://localhost:8080/healthz
```

> Detailed installation options (Docker, bare metal, Kubernetes) live in [docs/installation.md](docs/installation.md).

### Wire DriftLock into OpenTelemetry

Add the DriftLock processor to your OpenTelemetry configuration:

```yaml
processors:
  driftlock/anomaly:
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
      processors: [driftlock/anomaly, batch]
      exporters: [otlp]
```

## ✨ Features

- **🔍 Glass-Box Anomaly Detection**: Every anomaly comes with human-readable explanations
- **📊 Regulatory Compliance**: Built-in audit trails for DORA/NIS2/EU AI Act
- **⚡ High Performance**: Rust core with Go API server
- **🔧 OpenTelemetry Native**: Seamless integration with existing observability stacks
- **🌐 Open Source**: Apache 2.0 licensed, enterprise-friendly
- **📈 Real-time Monitoring**: Live dashboard with anomaly streaming

## 📊 Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Data Sources  │───▶│ OpenTelemetry    │───▶│ DriftLock Core  │
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

## 🆕 What's New in v0.1.0 (OSS Release)

### Standalone Operation
- ✅ **No External Dependencies**: Runs without Supabase or proprietary services
- ✅ **API Key Authentication**: Simple, secure access control
- ✅ **Simplified Deployment**: Just PostgreSQL and your API key
- ✅ **Full Functionality**: All features available in OSS version

### Migration from Previous Versions
If you're upgrading from a previous version:

1. **Configuration**: Update your `.env` file to use API keys instead of Supabase
2. **Dashboard**: Log in with your API key instead of email/password  
3. **API**: All endpoints remain compatible - no changes needed
4. **Data**: Your existing anomalies and configurations are preserved

### New Features
- **Self-Hosted**: Deploy anywhere - cloud, on-premise, or edge
- **Simplified Auth**: API key-based authentication for easy integration
- **Reduced Complexity**: Fewer moving parts, easier to maintain
- **Open Source**: Fully transparent, community-driven development

## 🏢 Enterprise Features (Optional)

## 🏢 Enterprise Features

Need compliance reports? Check out **[Shannon Labs Compliance Platform](https://compliance.shannonlabs.ai)**

- **DORA Quarterly Reports**: Audit-ready regulatory documentation ($299)
- **NIS2 Incident Reports**: Template-based incident reporting ($199)
- **EU AI Act Audit Trails**: Complete transparency documentation ($149)

## 📖 Documentation

- [Installation Guide](docs/installation.md)
- [Architecture Overview](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Compliance Integration](docs/COMPLIANCE_DORA.md)
- [Examples](examples/)
- [Security Guidance](SECURITY.md)

## 🛠️ Development

### Prerequisites

- Go 1.22+
- Rust stable (1.75+ recommended)
- Node.js 18+
- Docker & Docker Compose

### Local Development

```bash
# Install dependencies
make setup

# Run component tests
make test

# Start everything locally (Postgres, API, dashboard)
make dev

# Build production artifacts
make build
```

### Project Structure

```
driftlock/
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
```

## 🚀 Deployment

### Docker

```bash
# API server
docker build -t driftlock/api-server .

# Dashboard (serve the static build)
pnpm --dir web-frontend build
docker build -f deploy/docker/Dockerfile.web -t driftlock/web-frontend .
```

### Kubernetes / Helm

```
kubectl apply -f k8s/
# or
helm install driftlock ./helm/driftlock
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

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- OpenTelemetry community for the observability framework
- Compression-based anomaly detection research community
- Regulatory compliance experts who provided insights
- All our amazing contributors

## 📞 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/Shannon-Labs/driftlock/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)
- **Security**: security@shannonlabs.ai
- **Enterprise**: contact@shannonlabs.ai

---

**Built by [Shannon Labs](https://shannonlabs.ai)** - Making AI explainable and compliant.

If you find DriftLock useful, please give us a ⭐ on GitHub!
