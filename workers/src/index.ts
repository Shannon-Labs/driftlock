/**
 * Driftlock Cloudflare Workers - SaaS Backend
 * Replaces Firebase Functions with Cloudflare Workers
 */

import { GoogleGenerativeAI } from '@google/generative-ai';

// Cloud Run API endpoint (our main backend)
const CLOUD_RUN_API = 'https://driftlock-api-o6kjgrsowq-uc.a.run.app';

// CORS headers
const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Api-Key, Stripe-Signature',
};

// Helper to handle CORS preflight
function handleCORS(request: Request): Response | null {
  if (request.method === 'OPTIONS') {
    return new Response(null, {
      status: 204,
      headers: corsHeaders,
    });
  }
  return null;
}

// Helper to add CORS headers to response
function addCORSHeaders(response: Response): Response {
  const newHeaders = new Headers(response.headers);
  Object.entries(corsHeaders).forEach(([key, value]) => {
    newHeaders.set(key, value);
  });
  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: newHeaders,
  });
}

// Proxy signup requests to Cloud Run backend
async function handleSignup(request: Request): Promise<Response> {
  if (request.method !== 'POST') {
    return new Response(JSON.stringify({ error: 'Method not allowed' }), {
      status: 405,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  try {
    const body = await request.json();
    const { email, company_name } = body;

    // Forward to Cloud Run backend
    const backendResponse = await fetch(`${CLOUD_RUN_API}/v1/onboard/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        company_name,
      }),
    });

    const result = await backendResponse.json();
    return new Response(JSON.stringify(result), {
      status: backendResponse.status,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    console.error('Signup error:', error);
    return new Response(JSON.stringify({ error: 'Signup failed' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Enhanced anomaly analysis with Gemini
async function handleAnalyzeAnomalies(request: Request, env: Env): Promise<Response> {
  if (request.method !== 'POST') {
    return new Response(JSON.stringify({ error: 'Method not allowed' }), {
      status: 405,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  try {
    const body = await request.json();
    const { anomalies, query, api_key } = body;

    if (!anomalies || !Array.isArray(anomalies)) {
      return new Response(JSON.stringify({ error: 'Invalid anomalies data' }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' },
      });
    }

    // Verify API key with backend if provided
    if (api_key) {
      try {
        const authResponse = await fetch(`${CLOUD_RUN_API}/v1/anomalies`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${api_key}`,
            'Content-Type': 'application/json',
          },
        });

        if (!authResponse.ok) {
          return new Response(JSON.stringify({ error: 'Invalid API key' }), {
            status: 401,
            headers: { 'Content-Type': 'application/json' },
          });
        }
      } catch (authError) {
        console.warn('API key validation failed', authError);
      }
    }

    const genAI = new GoogleGenerativeAI(env.GEMINI_API_KEY || '');
    const model = genAI.getGenerativeModel({ model: 'gemini-pro' });

    const prompt = `
      Analyze these anomalies detected by Driftlock's compression-based system:
      
      ${anomalies.map((anomaly: any) => `
      - ID: ${anomaly.id}
      - NCD Score: ${anomaly.ncd_score || 'N/A'} 
      - P-value: ${anomaly.p_value || 'N/A'}
      - Explanation: ${anomaly.explanation || 'No explanation available'}
      - Data Type: ${anomaly.stream_type || 'unknown'}
      - Detected: ${anomaly.detected_at || 'recent'}
      `).join('\n')}
      
      User Query: ${query || 'Provide insights about these anomalies'}
      
      Please provide:
      1. Executive summary of anomaly patterns
      2. Risk assessment (Critical/High/Medium/Low)
      3. Recommended immediate actions
      4. Compliance implications for financial services
      5. Business impact assessment
      
      Keep the response professional and actionable for DevOps and security teams.
    `;

    const result = await model.generateContent(prompt);
    const analysis = result.response.text();

    return new Response(JSON.stringify({
      success: true,
      analysis,
      processed_anomalies: anomalies.length,
      timestamp: new Date().toISOString(),
      confidence: 'high',
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    console.error('Error analyzing anomalies:', error);
    return new Response(JSON.stringify({ error: 'Analysis failed' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Generate compliance reports
async function handleComplianceReport(request: Request, env: Env): Promise<Response> {
  if (request.method !== 'POST') {
    return new Response(JSON.stringify({ error: 'Method not allowed' }), {
      status: 405,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  try {
    const body = await request.json();
    const { anomalies, regulation, tenant_info, api_key } = body;

    // Verify API key if provided
    if (api_key) {
      try {
        const authResponse = await fetch(`${CLOUD_RUN_API}/v1/anomalies`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${api_key}`,
            'Content-Type': 'application/json',
          },
        });

        if (!authResponse.ok) {
          return new Response(JSON.stringify({ error: 'Invalid API key' }), {
            status: 401,
            headers: { 'Content-Type': 'application/json' },
          });
        }
      } catch (authError) {
        console.warn('API key validation failed', authError);
      }
    }

    const genAI = new GoogleGenerativeAI(env.GEMINI_API_KEY || '');
    const model = genAI.getGenerativeModel({ model: 'gemini-pro' });

    const prompt = `
      Generate a ${regulation || 'DORA'} compliance report for these anomalies:
      
      Company: ${tenant_info?.company_name || 'Customer'}
      Report Date: ${new Date().toLocaleDateString()}
      
      Detected Anomalies:
      ${anomalies.map((anomaly: any) => `
      - Anomaly ID: ${anomaly.id}
      - Detection Time: ${anomaly.detected_at}
      - NCD Score: ${anomaly.ncd_score} (Mathematical Evidence)
      - Statistical Significance: P-value ${anomaly.p_value}
      - Technical Explanation: ${anomaly.explanation}
      - Risk Level: ${anomaly.ncd_score > 0.7 ? 'HIGH' : anomaly.ncd_score > 0.4 ? 'MEDIUM' : 'LOW'}
      `).join('\n')}
      
      Generate a formal compliance report including:
      1. Executive Summary with business impact
      2. Technical Analysis with mathematical evidence (NCD, p-values)
      3. Risk Assessment per ${regulation || 'DORA'} requirements  
      4. Audit Trail and Evidence Documentation
      5. Recommended Immediate Actions
      6. Long-term Monitoring Recommendations
      
      Format as a professional regulatory document suitable for auditor review.
      Include references to compression-based anomaly detection methodology.
    `;

    const result = await model.generateContent(prompt);
    const report = result.response.text();

    return new Response(JSON.stringify({
      success: true,
      report,
      regulation: regulation || 'DORA',
      anomaly_count: anomalies.length,
      generated_at: new Date().toISOString(),
      company: tenant_info?.company_name || 'Customer',
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    console.error('Error generating compliance report:', error);
    return new Response(JSON.stringify({ error: 'Report generation failed' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Proxy API requests to Cloud Run backend
async function handleApiProxy(request: Request): Promise<Response> {
  let apiPath = new URL(request.url).pathname;

  // Handle specific rewrites
  if (apiPath === '/webhooks/stripe') {
    apiPath = '/v1/billing/webhook';
  } else if (apiPath === '/api/v1/healthz' || apiPath === '/healthz') {
    apiPath = '/healthz';
  } else if (apiPath.startsWith('/api/proxy')) {
    apiPath = apiPath.replace('/api/proxy', '');
  } else if (apiPath.startsWith('/api/v1')) {
    apiPath = apiPath.replace('/api', '');
  }

  const backendUrl = `${CLOUD_RUN_API}${apiPath}`;

  // Determine headers to forward
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  // Forward Authorization header
  const authHeader = request.headers.get('authorization');
  if (authHeader) {
    headers['Authorization'] = authHeader;
  }

  // Forward Stripe headers for webhooks
  const stripeSignature = request.headers.get('stripe-signature');
  if (stripeSignature) {
    headers['Stripe-Signature'] = stripeSignature;
  }

  // Forward X-Api-Key if present
  const apiKey = request.headers.get('x-api-key');
  if (apiKey) {
    headers['X-Api-Key'] = apiKey;
  }

  // Get request body
  let body: BodyInit | null = null;
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    // For Stripe webhooks, we need to preserve raw body
    if (apiPath === '/v1/billing/webhook') {
      body = await request.arrayBuffer();
    } else {
      body = await request.text();
    }
  }

  try {
    const backendResponse = await fetch(backendUrl, {
      method: request.method,
      headers: headers,
      body: body,
    });

    // Check content type of response
    const contentType = backendResponse.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      const result = await backendResponse.json();
      return new Response(JSON.stringify(result), {
        status: backendResponse.status,
        headers: { 'Content-Type': 'application/json' },
      });
    } else {
      const text = await backendResponse.text();
      return new Response(text, {
        status: backendResponse.status,
        headers: { 'Content-Type': 'text/plain' },
      });
    }
  } catch (error) {
    console.error('API proxy error:', error);
    return new Response(JSON.stringify({ error: 'Backend service unavailable' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Health check for the entire stack
async function handleHealthCheck(): Promise<Response> {
  const health: { [key: string]: any } = {
    success: true,
    status: 'healthy',
    service: 'driftlock-saas-backend',
    timestamp: new Date().toISOString(),
    version: '2.0.0',
    features: [
      'user-signup',
      'anomaly-analysis',
      'compliance-reporting',
      'gemini-integration',
      'cloud-run-proxy',
    ],
    backend: { status: 'unknown' },
  };

  // Check Cloud Run backend health
  try {
    const candidates = [
      `${CLOUD_RUN_API}/healthz`,
      `${CLOUD_RUN_API}/v1/healthz`,
    ];

    let backendResponse: Response | null = null;

    for (const url of candidates) {
      const attempt = await fetch(url, { method: 'GET' });
      if (attempt.ok) {
        backendResponse = attempt;
        health.backend.checked = url;
        break;
      }
      backendResponse = attempt;
    }

    if (backendResponse && backendResponse.ok) {
      const backendHealth = await backendResponse.json();
      health.backend = {
        status: 'healthy',
        database: backendHealth.database || 'unknown',
        license: backendHealth.license ? 'valid' : 'unknown',
      };
      health.success = backendHealth.success !== false;
    } else {
      const code = backendResponse?.status;
      health.backend = { status: code ? `unhealthy (${code})` : 'unhealthy' };
      health.success = false;
    }
  } catch (error) {
    health.backend = { status: 'unreachable' };
    health.success = false;
  }

  return new Response(JSON.stringify(health), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  });
}

// Get Firebase config (for backward compatibility during migration)
async function handleGetFirebaseConfig(env: Env): Promise<Response> {
  try {
    const firebaseConfig = {
      apiKey: env.VITE_FIREBASE_API_KEY || '',
      authDomain: env.VITE_FIREBASE_AUTH_DOMAIN || 'driftlock.firebaseapp.com',
      projectId: env.VITE_FIREBASE_PROJECT_ID || 'driftlock',
      storageBucket: env.VITE_FIREBASE_STORAGE_BUCKET || 'driftlock.appspot.com',
      messagingSenderId: env.VITE_FIREBASE_MESSAGING_SENDER_ID || '131489574303',
      appId: env.VITE_FIREBASE_APP_ID || '1:131489574303:web:e83e3e433912d05a8d61aa',
      measurementId: env.VITE_FIREBASE_MEASUREMENT_ID || 'G-CXBMVS3G8H',
    };

    return new Response(JSON.stringify(firebaseConfig), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    console.error('Error getting Firebase config:', error);
    return new Response(JSON.stringify({ error: 'Could not retrieve Firebase configuration.' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Main request handler
export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    // Handle CORS preflight
    const corsResponse = handleCORS(request);
    if (corsResponse) {
      return addCORSHeaders(corsResponse);
    }

    const url = new URL(request.url);
    const path = url.pathname;

    // Route to appropriate handler
    if (path === '/api/v1/onboard/signup' || path === '/signup') {
      const response = await handleSignup(request);
      return addCORSHeaders(response);
    }

    if (path === '/api/analyze' || path === '/analyze') {
      const response = await handleAnalyzeAnomalies(request, env);
      return addCORSHeaders(response);
    }

    if (path === '/api/compliance' || path === '/compliance') {
      const response = await handleComplianceReport(request, env);
      return addCORSHeaders(response);
    }

    if (path === '/api/v1/healthz' || path === '/healthz' || path === '/v1/healthz') {
      const response = await handleHealthCheck();
      return addCORSHeaders(response);
    }

    if (path === '/getFirebaseConfig') {
      const response = await handleGetFirebaseConfig(env);
      return addCORSHeaders(response);
    }

    // Default: proxy to Cloud Run backend
    if (path.startsWith('/api/') || path.startsWith('/webhooks/')) {
      const response = await handleApiProxy(request);
      return addCORSHeaders(response);
    }

    // 404 for unknown routes
    return new Response(JSON.stringify({ error: 'Not found' }), {
      status: 404,
      headers: { 'Content-Type': 'application/json' },
    });
  },
};

// Environment variables interface
interface Env {
  GEMINI_API_KEY?: string;
  VITE_FIREBASE_API_KEY?: string;
  VITE_FIREBASE_AUTH_DOMAIN?: string;
  VITE_FIREBASE_PROJECT_ID?: string;
  VITE_FIREBASE_STORAGE_BUCKET?: string;
  VITE_FIREBASE_MESSAGING_SENDER_ID?: string;
  VITE_FIREBASE_APP_ID?: string;
  VITE_FIREBASE_MEASUREMENT_ID?: string;
  CLOUD_RUN_API_URL?: string;
}



