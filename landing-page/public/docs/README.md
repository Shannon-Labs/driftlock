# Driftlock Documentation

Welcome to Driftlock's documentation! Learn how to use compression-based anomaly detection for your applications.

## 🚀 Getting Started

New to Driftlock? Start here:

- **[Quickstart Guide](./getting-started/quickstart.md)** - Get up and running in 5 minutes
- **[Core Concepts](./getting-started/concepts.md)** - Understand how Driftlock works
- **[Authentication](./getting-started/authentication.md)** - Set up API keys and Firebase Auth

## 📚 API Documentation

### REST API

- **[REST API Overview](./api/rest-api.md)** - Complete API reference
- **[POST /v1/detect](./api/endpoints/detect.md)** - Run anomaly detection
- **[GET /v1/anomalies](./api/endpoints/anomalies.md)** - Query anomalies
- **[GET /v1/anomalies/{id}](./api/endpoints/anomaly-detail.md)** - Get anomaly details

## 🛠️ SDKs & Libraries

Official client libraries to integrate Driftlock into your stack:

- **[Node.js SDK](./sdks/nodejs.md)** - Official Node.js/TypeScript client
- **[Python SDK](./sdks/python.md)** - Python client with async support
- **[Go SDK](./sdks/go.md)** - High-performance Go client
- **[JavaScript/Browser SDK](./sdks/javascript.md)** - Client-side integration

## 🔌 Integrations

Seamlessly integrate with your favorite frameworks:

- **[Express.js](./integrations/express.md)** - Middleware for Express applications
- **[Next.js](./integrations/nextjs.md)** - Full-stack Next.js integration
- **[Django](./integrations/django.md)** - Django middleware and views
- **[FastAPI](./integrations/fastapi.md)** - FastAPI dependency injection
- **[Spring Boot](./integrations/spring-boot.md)** - Java Spring Boot starter

## 💻 Developer Tools

Tools to speed up your development workflow:

- **[CLI Tool](./tools/cli.md)** - Command-line interface for testing
- **[Postman Collection](./tools/postman.md)** - API testing collection
- **[Webhooks](./tools/webhooks.md)** - Real-time notifications
- **[OpenAPI Spec](./tools/openapi.md)** - Swagger documentation
- **[Testing Strategies](./tools/testing.md)** - Unit and integration testing

## 🏭 Production Guide

Best practices for running Driftlock in production:

- **[Performance Tuning](./production/performance.md)** - Optimization guide
- **[Security](./production/security.md)** - API key management and security
- **[Monitoring](./production/monitoring.md)** - Metrics and dashboards
- **[Scaling](./production/scaling.md)** - Handling high throughput
- **[Troubleshooting](./production/troubleshooting.md)** - Common issues and solutions
- **Performance tiers:** Enterprise builds can enable our format-aware compression accelerator for sharper anomaly signals; default deployments stay on lightweight compressors with automatic fallback for portability.

## 💡 Real-World Examples

Complete implementation guides:

- **[Fraud Detection](./examples/fraud-detection.md)** - E-commerce fraud detection
- **[Log Analysis](./examples/log-analysis.md)** - DevOps log anomaly detection
- **[API Monitoring](./examples/api-monitoring.md)** - Microservices health monitoring
- **[IoT Sensors](./examples/iot-sensors.md)** - IoT device telemetry analysis

## 🔄 Migration Guides

Moving from another platform?

- **[From DataDog](./migration/from-datadog.md)** - Migration guide
- **[From New Relic](./migration/from-newrelic.md)** - Migration guide
- **[From Splunk](./migration/from-splunk.md)** - Migration guide
- **[Feature Comparison](./migration/comparison.md)** - Detailed comparison matrix

## 📘 Reference

- **[Code Examples](./api/examples/curl-examples.md)** - cURL and raw HTTP examples

## 🤝 Need Help?

- **Email**: support@driftlock.io
- **GitHub**: [Shannon-Labs/driftlock](https://github.com/Shannon-Labs/driftlock)
- **Community**: [GitHub Discussions](https://github.com/Shannon-Labs/driftlock/discussions)

---

**Documentation Version**: 1.1.0
**Last Updated**: November 19, 2025
