Driftlock

A compression-based anomaly detection (CBAD) platform for OpenTelemetry data, powered by Meta's OpenZL format-aware compression framework. Driftlock provides explainable anomaly detection for regulated industries through advanced compression analysis of logs, metrics, traces, and LLM I/O.

## Architecture

Driftlock consists of two main services:

1. **Go API Server** (`api-server/`): Core anomaly detection pipeline
   - Processes events and detects anomalies using CBAD algorithm
   - Exposes REST API for anomaly management
   - Integrates with Supabase for web-frontend data synchronization

2. **React Web Frontend** (`web-frontend/`): Customer dashboard and billing
   - Built with React/TypeScript and shadcn-ui
   - Uses Supabase as backend (PostgreSQL + Edge Functions)
   - Handles user authentication, billing, and anomaly visualization

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local development)
- Supabase account with project created

### Environment Setup

1. Copy environment template:
   ```bash
   cp .env.example .env
   ```

2. Update Supabase configuration in `.env`:
   ```bash
   SUPABASE_PROJECT_ID=your_actual_project_id
   SUPABASE_ANON_KEY=your_actual_anon_key
   SUPABASE_SERVICE_ROLE_KEY=your_actual_service_role_key
   SUPABASE_BASE_URL=https://your_project_id.supabase.co
   SUPABASE_WEBHOOK_URL=https://your_project_id.supabase.co/functions/v1/webhook
   ```

### Running with Docker Compose

1. Start all services:
   ```bash
   ./start.sh
   ```

2. Services will be available at:
   - Web Frontend: http://localhost:3000
   - Go API: http://localhost:8080
   - API Documentation: http://localhost:8080/healthz

### Local Development

#### Web Frontend
```bash
cd web-frontend
npm install
npm run dev
```

#### Go API Server
```bash
cd api-server
go run ./cmd/api-server
```

## Integration Details

### Data Flow

1. **Event Ingestion**: Events are sent to Go API `/v1/events` endpoint
2. **Anomaly Detection**: Go API processes events using CBAD algorithm
3. **Data Synchronization**: Anomalies are stored in both PostgreSQL and Supabase
4. **Real-time Updates**: Web frontend receives updates via Supabase subscriptions
5. **Billing**: Usage is tracked in Supabase for subscription management

### API Integration

The Go API server integrates with Supabase through:

- **Anomaly Synchronization**: Anomalies created in Go API are also stored in supabase
- **Status Updates**: Anomaly status changes are synchronized between systems
- **Usage Tracking**: API usage is tracked in supabase for billing
- **Webhook Notifications**: Go API can trigger Supabase Edge Functions

## Documentation

- [Integration Guide](INTEGRATION_README.md) - Detailed setup and integration instructions
- [API Documentation](docs/API.md) - REST API reference
- [Architecture](docs/ARCHITECTURE.md) - System design and components
