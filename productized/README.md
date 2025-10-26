# DriftLock - Anomaly Detection Platform

DriftLock is a comprehensive anomaly detection platform that monitors your systems, applications, and infrastructure for unusual behavior patterns. Our platform uses advanced algorithms to detect anomalies in logs, metrics, and traces, alerting you before issues impact your users.

## Features

- **Real-time Anomaly Detection**: Detects anomalies in logs, metrics, and traces as they happen
- **Multiple Detection Algorithms**: Statistical analysis, ML-based detection, and customizable rules
- **Multi-tenancy**: Secure separation of customer data and configurations
- **Comprehensive Alerting**: Email, Slack, webhook integrations with customizable thresholds
- **Customizable Dashboards**: User-configurable views with advanced visualizations
- **API Access**: Full-featured API for integration with your existing tools
- **Security First**: Enterprise-grade security with SOC2 compliance
- **Scalable Architecture**: Designed to handle millions of events per second

## Architecture

DriftLock follows a microservices architecture:

- **API Gateway**: Handles authentication, rate limiting, and request routing
- **Ingestion Service**: Processes incoming events from various sources
- **Anomaly Detection Engine**: Runs multiple detection algorithms in parallel
- **Alerting Service**: Manages delivery of alerts through multiple channels
- **Dashboard Service**: Real-time visualization and user interface
- **User Management**: Authentication, authorization, and user data management

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- PostgreSQL (or use the provided Docker setup)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/driftlock.git
   cd driftlock
   ```

2. Start the development environment:
   ```bash
   docker-compose up -d
   ```

3. Access the API at `http://localhost:8080`
4. Access the database admin at `http://localhost:8081`

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```bash
SERVER_ADDRESS=:8080
DATABASE_URL=postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
SESSION_SECRET=your-super-secret-session-key-change-in-production
KAFKA_BROKERS=localhost:9092
DEBUG=true
```

## API Documentation

See our [API Documentation](docs/api.md) for detailed information about endpoints, authentication, and usage examples.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please contact us at support@driftlock.com or open an issue in this repository.