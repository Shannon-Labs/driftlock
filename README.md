# DriftLock by Shannon Labs

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shannon-Labs/driftlock)](https://goreportcard.com/report/github.com/Shannon-Labs/driftlock)
[![Rust](https://img.shields.io/badge/rust-%23000000.svg?style=for-the-badge&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![CI](https://github.com/Shannon-Labs/driftlock/workflows/CI/badge.svg)](https://github.com/Shannon-Labs/driftlock/actions)

> **Explainable AI anomaly detection for regulated industries**

DriftLock provides compression-based anomaly detection (CBAD) with glass-box explainability, designed specifically for compliance with DORA, NIS2, and EU AI Act regulations.

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Start with Docker Compose (recommended)
docker-compose up -d

# Access the dashboard
open http://localhost:3000
```

### Using the OpenTelemetry Collector

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
│ Compliance      │◀───│   API Server     │◀───│  Explanations   │
│ Reports         │    │    (Go)          │    │   & Audit Trail │
│ (Shannon Labs)  │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🏢 Enterprise Features

Need compliance reports? Check out **[Shannon Labs Compliance Platform](https://compliance.shannonlabs.ai)**

- **DORA Quarterly Reports**: Audit-ready regulatory documentation ($299)
- **NIS2 Incident Reports**: Template-based incident reporting ($199)
- **EU AI Act Audit Trails**: Complete transparency documentation ($149)

## 📖 Documentation

- [Installation Guide](docs/installation.md)
- [Architecture Overview](docs/architecture.md)
- [API Reference](docs/api-reference.md)
- [Compliance Integration](docs/compliance.md)
- [Examples](examples/)
- [Troubleshooting](docs/troubleshooting.md)

## 🛠️ Development

### Prerequisites

- Go 1.24+
- Rust 1.70+
- Node.js 18+
- Docker & Docker Compose

### Local Development

```bash
# Install dependencies
make setup

# Run tests
make test

# Run all services locally
make dev

# Build components
make build
```

### Project Structure

```
driftlock/
├── src/
│   ├── anomaly-detection/    # Rust CBAD core
│   ├── api-server/          # Go REST API
│   ├── dashboard/           # React web UI
│   └── otel-collector/      # OpenTelemetry processor
├── docs/                    # Documentation
├── examples/                # Usage examples
├── deployments/             # Docker, K8s, Helm
└── scripts/                 # Build & utility scripts
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific component tests
cd src/anomaly-detection && cargo test
cd src/api-server && go test ./...
cd src/dashboard && npm test

# Integration tests
make test-integration
```

## 🚀 Deployment

### Docker

```bash
docker build -t driftlock/api-server ./src/api-server
docker build -t driftlock/dashboard ./src/dashboard
```

### Kubernetes

```bash
kubectl apply -f deployments/kubernetes/
```

### Helm

```bash
helm install driftlock ./deployments/helm/driftlock
```

## 📊 Compliance & Security

- **Explainable AI**: Every anomaly includes mathematical explanations
- **Audit Trails**: Complete logging for regulatory compliance
- **Data Privacy**: GDPR-compliant data handling
- **Security**: Built-in authentication and encryption
- **Transparency**: Open source with clear documentation

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

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