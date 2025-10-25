# Driftlock UI

Production-ready Next.js 14 web interface for Driftlock anomaly detection platform.

## Features

- **Anomaly Dashboard**: Browse and filter detected anomalies with advanced search
- **Detail View**: Investigate individual anomalies with glass-box explanations
- **Live Feed**: Real-time anomaly stream via Server-Sent Events (SSE)
- **Analytics**: Performance metrics, detection rates, and compression efficiency
- **Configuration**: Manage detection thresholds and stream settings
- **Mobile-Responsive**: Optimized for desktop, tablet, and mobile devices

## Tech Stack

- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Charts**: Recharts
- **State**: React Hooks + React Query
- **Real-time**: Server-Sent Events (EventSource API)

## Quick Start

```bash
# Install dependencies
pnpm install

# Run development server
pnpm dev

# Build for production
pnpm build

# Start production server
pnpm start
```

The UI will be available at http://localhost:3000

## Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Pages

- `/` - Dashboard with system overview
- `/anomalies` - Anomaly list with filters and search
- `/anomalies/[id]` - Detailed anomaly investigation
- `/live` - Real-time anomaly feed
- `/analytics` - Statistical analysis and performance metrics
- `/config` - System configuration and tuning

## API Integration

The UI connects to the Driftlock API server (default: http://localhost:8080) and expects these endpoints:

- `GET /v1/anomalies` - List anomalies with pagination and filters
- `GET /v1/anomalies/:id` - Get anomaly details
- `PATCH /v1/anomalies/:id/status` - Update anomaly status
- `GET /v1/stream/anomalies` - SSE stream for real-time updates
- `GET /v1/config` - Get detection configuration
- `PATCH /v1/config` - Update configuration
- `GET /v1/metrics/performance` - Performance metrics
- `GET /v1/analytics/*` - Analytics endpoints

## Development Features

- **Mock Data**: The UI includes mock data for development when the API is unavailable
- **Error Handling**: Graceful fallbacks and user-friendly error messages
- **Type Safety**: Full TypeScript coverage for API responses
- **Dark Mode**: Automatic dark mode support based on system preferences

## Phase 3 Completion Status

✅ Next.js 14 App Router setup with TypeScript
✅ Anomaly list view with filters, search, and pagination
✅ Anomaly detail view with glass-box explanations
✅ Real-time live feed with SSE
✅ Compression ratio visualizations
✅ Analytics dashboard with performance metrics
✅ Configuration management UI
✅ Mobile-responsive design
✅ Dark mode support

## Next Steps

- Connect to production API endpoints
- Add authentication/authorization
- Implement advanced chart visualizations
- Add export functionality for compliance reports
- Optimize bundle size and performance
