# Shannon Labs Compliance Platform

**Private repository - Shannon Labs internal use**

A compliance reporting platform that generates audit-ready documentation for DORA, NIS2, and EU AI Act regulations based on anomaly data from the DriftLock open-source platform.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ DriftLock OSS   │───▶│ Compliance API   │───▶│ Report Generator│
│ Anomaly Data    │    │ (Go/Node.js)     │    │ (PDF + JSON)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Customer        │◀───│ Stripe Checkout  │◀───│ Report Storage  │
│ Portal          │    │ ($299-$1,499)    │    │ (S3/Local)      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Services

### 1. Report Generator (`src/reports/`)
- **DORA Quarterly Reports** ($299)
  - ICT risk assessment
  - Digital operational resilience testing
  - Incident reporting documentation
  - Backup and recovery validation

- **NIS2 Incident Reports** ($199)
  - Security incident notification
  - Risk management measures
  - Supply chain security
  - Business continuity planning

- **EU AI Act Audit Trails** ($149)
  - High-risk AI system classification
  - Transparency obligations
  - Human oversight documentation
  - Technical documentation

### 2. Customer Portal (`src/dashboard/`)
- Purchase history and report library
- Download management with expiring links
- Compliance status tracking
- Renewal prompts and notifications

### 3. Billing System (`src/billing/`)
- Stripe checkout integration
- Webhook handling for payment confirmation
- Subscription management
- Automated report fulfillment

## Getting Started

### Prerequisites
- Node.js 18+
- PostgreSQL
- Redis (for sessions)
- Stripe account with API keys

### Installation

```bash
# Install dependencies
npm install

# Set up environment variables
cp .env.example .env

# Set up database
npm run db:migrate

# Start development server
npm run dev
```

### Environment Variables

```bash
# Database
DATABASE_URL="postgresql://user:password@localhost:5432/compliance"

# Stripe
STRIPE_SECRET_KEY="sk_test_..."
STRIPE_PUBLISHABLE_KEY="pk_test_..."
STRIPE_WEBHOOK_SECRET="whsec_..."

# Storage
STORAGE_TYPE="s3" # or "local"
S3_BUCKET="Shannon-Labs-reports"
S3_REGION="us-east-1"
S3_ACCESS_KEY="..."
S3_SECRET_KEY="..."

# DriftLock OSS Integration
DRIFTLOCK_API_URL="https://api.driftlock.shannonlabs.ai"
DRIFTLOCK_API_KEY="..."

# Security
JWT_SECRET="your-jwt-secret"
NEXTAUTH_SECRET="your-nextauth-secret"
```

## Report Generation Process

### 1. Data Collection
```javascript
// Fetch anomaly data from DriftLock OSS
const anomalies = await fetchDriftLockData(customerId, dateRange);
```

### 2. Compliance Analysis
```javascript
// Analyze against regulatory frameworks
const doraAnalysis = analyzeDORACompliance(anomalies);
const nis2Analysis = analyzeNIS2Compliance(anomalies);
const euAIAnalysis = analyzeEUAICompliance(anomalies);
```

### 3. Report Generation
```javascript
// Generate PDF and JSON reports
const pdfReport = await generatePDFReport(doraAnalysis, template);
const jsonReport = await generateJSONReport(doraAnalysis);
```

### 4. Delivery
```javascript
// Store with expiring URLs
const downloadUrl = await storeReport(report, expires: '7d');
await notifyCustomer(customerId, downloadUrl);
```

## API Endpoints

### Public
- `GET /api/pricing` - Available report types and pricing
- `POST /api/checkout` - Create Stripe checkout session
- `GET /api/reports/:id/download` - Download purchased report

### Authenticated
- `GET /api/reports` - Customer's report library
- `POST /api/reports/:id/regenerate` - Regenerate expired report
- `GET /api/subscriptions` - Active subscriptions

### Webhooks
- `POST /api/webhooks/stripe` - Stripe payment webhooks

## Database Schema

### Core Tables
- `customers` - Customer information and billing
- `reports` - Generated reports and metadata
- `subscriptions` - Active subscriptions
- `downloads` - Download tracking and expiring links

## Security Considerations

1. **Data Isolation**: Customer data strictly separated
2. **Encryption**: All sensitive data encrypted at rest
3. **Access Control**: JWT-based authentication
4. **Audit Logging**: All actions logged for compliance
5. **Rate Limiting**: API endpoints rate-limited
6. **Secure Storage**: Reports stored with access controls

## Deployment

### Docker

```bash
# Build image
docker build -t Shannon-Labs/compliance-platform .

# Run with environment
docker run -p 3000:3000 \
  -e DATABASE_URL="..." \
  -e STRIPE_SECRET_KEY="..." \
  Shannon-Labs/compliance-platform
```

### Production Checklist
- [ ] SSL certificates configured
- [ ] Database backups enabled
- [ ] Stripe webhooks configured
- [ ] Rate limiting enabled
- [ ] Monitoring and alerting
- [ ] Security scanning completed
- [ ] Performance testing completed

## Monitoring

### Key Metrics
- Report generation success rate
- Payment conversion rate
- Customer satisfaction scores
- API response times
- Error rates by endpoint

### Alerts
- Payment failures
- Report generation failures
- High error rates
- Security events

## Support

- **Internal Documentation**: Confluence workspace
- **Issue Tracking**: Linear project
- **Customer Support**: Zendesk integration
- **Monitoring**: Datadog dashboard

## Legal & Compliance

This platform handles sensitive compliance data and must adhere to:
- GDPR data protection requirements
- SOC 2 Type II compliance
- ISO 27001 security standards
- Industry-specific regulations

## Development Workflow

1. Feature development in feature branches
2. Code review required for all changes
3. Automated testing in CI/CD
4. Staging deployment for review
5. Production deployment with monitoring

## Future Enhancements

- Multi-language support (EU regulations)
- Additional compliance frameworks (SOX, HIPAA)
- AI-powered report recommendations
- Integration with major compliance tools
- Mobile app for compliance managers