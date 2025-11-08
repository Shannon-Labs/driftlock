# Driftlock

Driftlock is an open-source anomaly detection system that uses compression-based algorithms to provide explainable results for observability data. Designed for organizations requiring transparent AI systems, it processes telemetry data through OpenTelemetry and outputs detailed explanations for detected anomalies.

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

Or use Make commands:

```bash
make quick-start  # Sets up .env and starts services
make dev          # Start development environment
make stop         # Stop all services
```

Detailed installation options (Docker, bare metal, Kubernetes) live in docs/installation.md.

## Wire Driftlock into OpenTelemetry

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

## Features

- **Compression-Based Anomaly Detection**: Uses normalized compression distance (NCD) algorithms
- **Explainable Results**: Each anomaly includes mathematical explanations and context
- **OpenTelemetry Integration**: Processes logs, metrics, and traces via standard collectors
- **Self-Hosted**: Runs on-premise with Docker and PostgreSQL
- **API Access**: RESTful API for programmatic access and integration
- **Web Dashboard**: Browser-based interface for monitoring and configuration
- **Audit Logging**: Complete audit trail for compliance and debugging

## Architecture

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

## Testing

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

## Deployment

### Docker

```bash
# API server
docker build -t driftlock-api .

# Dashboard (serve the static build)
pnpm --dir web-frontend build
docker build -f deploy/docker/Dockerfile.web -t driftlock-web-frontend .
```

### Kubernetes / Helm

```bash
kubectl apply -f k8s/
# or
helm install driftlock ./helm/driftlock
```

## Security & Compliance

- **Explainable AI**: Every anomaly includes mathematical explanations
- **Audit Trails**: Complete logging for regulatory compliance
- **Data Privacy**: GDPR-compliant data handling
- **Security**: Built-in authentication and encryption
- **Transparency**: Open source with clear documentation

## Contributing

We welcome contributions! Please see CONTRIBUTING.md for guidelines and local setup steps.

## License

Licensed under **Apache 2.0** with patent protections. See [LICENSE](LICENSE) and [PATENTS.md](PATENTS.md) for details.

### Commercial Licensing

For commercial use, proprietary licenses are available through Shannon Labs. This dual-licensing model allows you to:
- Use Driftlock under Apache 2.0 for open source and commercial projects
- Purchase a commercial license for proprietary applications

Contact: hunter@shannonlabs.dev

## Acknowledgments

- OpenTelemetry community for the observability framework
- Compression-based anomaly detection research community  
- Regulatory compliance experts who provided insights
- All our amazing contributors

## Support

- **Documentation**: docs/
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Security**: hunter@shannonlabs.dev
- **Enterprise**: hunter@shannonlabs.dev

Developed by Shannon Labs.
