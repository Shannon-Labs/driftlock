# DriftLock - Anomaly Detection Platform

DriftLock is a comprehensive anomaly detection platform that monitors your systems, applications, and infrastructure for unusual behavior patterns. Our platform uses advanced algorithms to detect anomalies in logs, metrics, and traces, alerting you before issues impact your users.

## Features

- **Real-time Anomaly Detection**: Detects anomalies in logs, metrics, and traces as they happen
- **Multiple Detection Algorithms**: Statistical analysis, ML-based detection, and customizable rules
- **Multi-tenancy**: Secure separation of customer data and configurations
- **Subscription Billing**: Integrated Stripe-based subscription management
- **Comprehensive Alerting**: Email, Slack, webhook integrations with customizable thresholds
- **Customizable Dashboards**: User-configurable views with advanced visualizations
- **API Access**: Full-featured API for integration with your existing tools
- **Security First**: Enterprise-grade security with SOC2 compliance
- **Scalable Architecture**: Designed to handle millions of events per second
- **Analytics & Audit Logging**: Built-in analytics and compliance logging
- **Onboarding Flow**: Guided setup process for new users
- **Cloudflare Integration**: Edge computing, CDN, and security powered by Cloudflare

## Architecture

DriftLock follows a microservices architecture:

- **API Gateway**: Handles authentication, rate limiting, and request routing
- **Ingestion Service**: Processes incoming events from various sources
- **Anomaly Detection Engine**: Runs multiple detection algorithms in parallel
- **Billing Service**: Integrated subscription and payment processing with Stripe
- **Alerting Service**: Manages delivery of alerts through multiple channels
- **Email Service**: Notification system with multiple providers
- **Analytics Service**: GA4 and custom analytics integration
- **Audit Service**: Comprehensive compliance and security logging
- **Dashboard Service**: Real-time visualization and user interface
- **User Management**: Authentication, authorization, and user data management
- **Onboarding Service**: Guided user setup and configuration

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 18+ (for frontend)
- PostgreSQL (or use the provided Docker setup)
- Stripe account for billing
- SendGrid account for emails (or SMTP server)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/driftlock.git
   cd driftlock
   ```

2. Copy the environment example:
   ```bash
   cd productized
   cp .env.example .env
   # Update the .env file with your actual values
   ```

3. Start the development environment:
   ```bash
   cd productized
   docker-compose up -d
   ```

4. Install Go dependencies:
   ```bash
   cd productized
   go mod tidy
   ```

5. Run the backend API:
   ```bash
   cd productized
   go run cmd/server/main.go
   ```

6. In a new terminal, run the frontend:
   ```bash
   cd productized/frontend
   npm install
   npm start
   ```

7. Access the API at `http://localhost:8080`
8. Access the frontend at `http://localhost:3000`
9. Access the database admin at `http://localhost:8081`

### Environment Variables

See the [.env.example](.env.example) file for all configuration options.

## Frontend

The frontend is built with React and can be found in the `frontend/` directory:

```bash
cd productized/frontend
npm install
npm run dev
```

## API Documentation

See our [API Documentation](docs/api.md) for detailed information about endpoints, authentication, and usage examples.

## Billing Integration

The billing system integrates with Stripe for subscription management. Key features include:
- Subscription plans (Free, Pro, Business)
- Usage-based billing
- Automatic invoicing
- Webhook handling for events
- Customer portal for managing subscriptions

## Email Service

Email notifications are handled through multiple providers:
- SMTP support for custom email servers
- SendGrid integration for transactional emails
- Predefined templates for common notifications
- Welcome emails for new users
- Anomaly alerts and billing notifications

## Analytics & Audit Logging

Analytics tracking includes:
- Google Analytics 4 integration
- Custom event tracking
- User behavior analytics
- API usage metrics
- Comprehensive audit logging for compliance
- GDPR-ready data retention policies

## Onboarding Flow

New users are guided through a step-by-step onboarding process:
- Profile setup
- Data source connection
- Alert configuration
- Dashboard overview
- Anomaly resolution guidance

## Cloudflare Integration

The platform leverages Cloudflare services:
- Cloudflare Pages for frontend hosting
- Cloudflare Workers for edge computing
- DDoS protection and WAF
- Rate limiting at the edge
- Custom domains and SSL
- CDN for static assets

## Production Deployment

For production deployment, see the [deployment documentation](deployment/README.md) for detailed instructions on setting up the infrastructure using Kubernetes and Cloudflare.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please contact us at support@driftlock.com or open an issue in this repository.