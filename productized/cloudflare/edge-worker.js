// DriftLock Edge Workers
// Handles rate limiting, caching, and edge computing for the application

// Rate limiting configuration
const RATE_LIMIT_KEY = 'rate_limit';
const RATE_LIMIT_WINDOW = 60; // 60 seconds
const RATE_LIMIT_MAX = 100; // 100 requests per window

// Cache configuration
const CACHE_TTL = 300; // 5 minutes

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const path = url.pathname;
    
    // Apply rate limiting to API routes
    if (path.startsWith('/api/')) {
      const ip = request.headers.get('CF-Connecting-IP');
      const key = `${RATE_LIMIT_KEY}:${ip}`;
      
      // Check rate limit
      const currentCount = await env.KV.get(key);
      const count = currentCount ? parseInt(currentCount) : 0;
      
      if (count >= RATE_LIMIT_MAX) {
        return new Response('Rate limit exceeded', { status: 429 });
      }
      
      // Increment rate limit counter
      ctx.waitUntil(env.KV.put(key, count + 1, { expirationTtl: RATE_LIMIT_WINDOW }));
    }
    
    // Add security headers
    const response = await fetch(request);
    const newResponse = new Response(response.body, response);
    newResponse.headers.set('X-Content-Type-Options', 'nosniff');
    newResponse.headers.set('X-Frame-Options', 'DENY');
    newResponse.headers.set('X-XSS-Protection', '1; mode=block');
    
    return newResponse;
  }
};