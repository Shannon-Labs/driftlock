/**
 * Cloudflare Pages Function to proxy API requests
 * 
 * This function proxies requests to the actual Go API backend.
 * Set API_BACKEND_URL environment variable in Cloudflare Pages to point to your deployed API.
 */

interface EventContext {
  request: Request;
  params: { path?: string[] };
  env: { API_BACKEND_URL?: string };
}

const BACKEND_URL = (globalThis as any).API_BACKEND_URL || 'http://localhost:8080';

export async function onRequest(context: EventContext): Promise<Response> {
  const { request, params, env } = context;
  const backendUrl = env.API_BACKEND_URL || BACKEND_URL;
  const path = params.path as string[] || [];
  
  // Build the backend URL
  const backendPath = path.length > 0 ? path.join('/') : '';
  const searchParams = new URL(request.url).search;
  const url = new URL(`/${backendPath}${searchParams}`, backendUrl);
  
  // Forward the request to the backend
  const backendRequest = new Request(url.toString(), {
    method: request.method,
    headers: request.headers,
    body: request.body,
  });
  
  try {
    const response = await fetch(backendRequest);
    
    // Create a new response with CORS headers
    const corsHeaders = new Headers(response.headers);
    corsHeaders.set('Access-Control-Allow-Origin', '*');
    corsHeaders.set('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    corsHeaders.set('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With, X-Request-Id');
    
    return new Response(response.body, {
      status: response.status,
      statusText: response.statusText,
      headers: corsHeaders,
    });
  } catch (error) {
    return new Response(
      JSON.stringify({ 
        success: false, 
        error: 'Backend API unavailable',
        message: error instanceof Error ? error.message : 'Unknown error'
      }),
      {
        status: 503,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
      }
    );
  }
}

