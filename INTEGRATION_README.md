# DriftLock Integration Guide

This document explains how to run DriftLock with the integrated web-frontend.

## Architecture

DriftLock now consists of two main services:

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
   docker-compose up -d
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

- **Anomaly Synchronization**: Anomalies created in Go API are also stored in Supabase
- **Status Updates**: Anomaly status changes are synchronized between systems
- **Usage Tracking**: API usage is tracked in Supabase for billing
- **Webhook Notifications**: Go API can trigger Supabase Edge Functions

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SUPABASE_PROJECT_ID` | Supabase project ID | - |
| `SUPABASE_ANON_KEY` | Supabase anonymous key | - |
| `SUPABASE_SERVICE_ROLE_KEY` | Supabase service role key | - |
| `SUPABASE_BASE_URL` | Supabase project URL | - |
| `SUPABASE_WEBHOOK_URL` | Supabase webhook URL | - |

## Deployment

### Production Deployment

1. Build and push Docker images:
   ```bash
   docker-compose -f docker-compose.yml build
   docker-compose -f docker-compose.yml push
   ```

2. Configure production environment variables

3. Deploy with your preferred orchestration platform

### Monitoring

- **Go API**: Prometheus metrics at `http://localhost:9090/metrics`
- **Web Frontend**: Supabase dashboard and logs
- **Infrastructure**: Health checks at `http://localhost:8080/healthz`

## Troubleshooting

### Port Conflicts
- Web frontend runs on port 3000
- Go API runs on port 8080
- Change ports in docker-compose.yml if needed

### Database Issues
- Check PostgreSQL connection in Go API logs
- Verify Supabase project settings
- Ensure Supabase migrations are applied

### Authentication Issues
- Verify Supabase keys in `.env`
- Check RLS policies in Supabase
- Review API key configuration

## Development Notes

- The Go API server can run independently without Supabase
- Web frontend requires Supabase configuration
- Both services share anomaly data through Supabase
- Redis is optional but recommended for caching
