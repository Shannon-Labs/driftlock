# Next.js Integration

Integrate Driftlock into your Next.js application (App Router or Pages Router).

## Installation

```bash
npm install @driftlock/next
```

## App Router (Next.js 13+)

### Middleware Integration

Create or update `middleware.ts` in your project root to monitor all requests.

```typescript
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { driftlockMiddleware } from '@driftlock/next/middleware';

export function middleware(request: NextRequest) {
  // Run Driftlock anomaly detection
  return driftlockMiddleware(request, {
    apiKey: process.env.DRIFTLOCK_API_KEY!,
    streamId: 'nextjs-app',
    // Optional: Block anomalous requests
    blockAnomalies: false 
  });
}

export const config = {
  matcher: '/api/:path*', // Monitor API routes
};
```

### Route Handlers

You can also use Driftlock directly in specific Route Handlers.

```typescript
// app/api/hello/route.ts
import { NextResponse } from 'next/server';
import { Driftlock } from '@driftlock/next';

const driftlock = new Driftlock({ apiKey: process.env.DRIFTLOCK_API_KEY });

export async function POST(request: Request) {
  const body = await request.json();

  // Detect anomalies in the request body
  const detection = await driftlock.detect({
    streamId: 'hello-api',
    events: [{ type: 'request', body }]
  });

  if (detection.anomalies.length > 0) {
    return NextResponse.json({ error: 'Anomaly detected' }, { status: 403 });
  }

  return NextResponse.json({ message: 'Hello World' });
}
```

## Pages Router

### API Routes

Wrap your API handlers with `withDriftlock`.

```typescript
// pages/api/users.ts
import type { NextApiRequest, NextApiResponse } from 'next';
import { withDriftlock } from '@driftlock/next';

async function handler(req: NextApiRequest, res: NextApiResponse) {
  res.status(200).json({ name: 'John Doe' });
}

export default withDriftlock(handler, {
  apiKey: process.env.DRIFTLOCK_API_KEY,
  streamId: 'users-api'
});
```

### getServerSideProps

Monitor server-side rendering performance and anomalies.

```typescript
import { withDriftlockSSR } from '@driftlock/next';

export const getServerSideProps = withDriftlockSSR(async (context) => {
  // Your logic here
  return {
    props: {},
  };
}, {
  apiKey: process.env.DRIFTLOCK_API_KEY,
  streamId: 'ssr-pages'
});
```

## Vercel Edge Functions

Driftlock is fully compatible with Vercel Edge Functions.

```typescript
import { DriftlockClient } from '@driftlock/client'; // Use the core client

export const config = {
  runtime: 'edge',
};

export default async function handler(req: Request) {
  const client = new DriftlockClient({ apiKey: process.env.DRIFTLOCK_API_KEY });
  
  // Async detection (fire and forget)
  context.waitUntil(
    client.detect({
      streamId: 'edge-functions',
      events: [{ type: 'edge-req', body: { url: req.url } }]
    })
  );

  return new Response('Hello from Edge!');
}
```

## Configuration

The `@driftlock/next` package accepts the same configuration options as the Node.js SDK, plus Next.js specific options:

- `excludePaths`: Array of paths to ignore (middleware only)
- `captureBody`: Boolean, whether to include request body in analysis (default: `true`)

## Troubleshooting

### "API Key Missing"
Ensure `DRIFTLOCK_API_KEY` is set in your `.env.local` file and added to your Vercel project environment variables.

### "Cold Start Latency"
Driftlock's HTTP requests are optimized, but if you are on a serverless cold start, the first request might be slower. Use `blockAnomalies: false` (default) to ensure detection happens asynchronously and doesn't impact user latency.
