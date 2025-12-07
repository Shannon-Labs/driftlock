/**
 * Driftlock Firebase Functions - SaaS Backend
 * Integrated with Cloud Run API and Gemini AI
 * Updated: Force redeploy for invoker config (attempt 3)
 */

import { setGlobalOptions, logger } from "firebase-functions";
import { onRequest } from "firebase-functions/v2/https";
import { GoogleGenerativeAI } from "@google/generative-ai";
import { SecretManagerServiceClient } from "@google-cloud/secret-manager";

// type InvokerConfig = "public" | "private" | string | string[];

interface ProjectInfo {
  projectId?: string;
  projectNumber?: string;
}

const projectInfo = resolveProjectInfo();
// const invokerConfig = resolveInvokers(projectInfo);

// Set global options for cost control
setGlobalOptions({ maxInstances: 10 });
logger.info("Configured function invokers (Managed manually)", {
  // invokerConfig,
  projectInfo,
});

/*
function parseInvokerList(value?: string): string[] {
  if (!value) return [];
  return value
    .split(/[,;\s]+/)
    .map((part) => part.trim())
    .filter((part) => part.length > 0);
}
*/

function resolveProjectInfo(): ProjectInfo {
  let firebaseConfig: any;
  if (process.env.FIREBASE_CONFIG) {
    try {
      firebaseConfig = JSON.parse(process.env.FIREBASE_CONFIG);
    } catch (error) {
      logger.warn("Failed to parse FIREBASE_CONFIG", error as Error);
    }
  }

  return {
    projectId:
      process.env.GCLOUD_PROJECT ||
      process.env.GCP_PROJECT ||
      firebaseConfig?.projectId,
    projectNumber:
      process.env.GCLOUD_PROJECT_NUMBER ||
      process.env.GCP_PROJECT_NUMBER ||
      firebaseConfig?.projectNumber,
  };
}

/*
function resolveInvokers(info: ProjectInfo): InvokerConfig {
  const explicit = [
    ...parseInvokerList(process.env.FUNCTIONS_INVOKERS),
    ...parseInvokerList(process.env.FIREBASE_FUNCTIONS_INVOKERS),
    ...parseInvokerList(process.env.ALLOWED_FUNCTION_INVOKERS),
  ].filter(Boolean);

  if (explicit.length > 0) {
    return Array.from(new Set(explicit));
  }

  const derived = new Set<string>(["firebase-hosting@system.gserviceaccount.com"]);

  if (info.projectNumber) {
    derived.add(
      `service-${info.projectNumber}@gcp-sa-firebasehosting.iam.gserviceaccount.com`
    );
    derived.add(`${info.projectNumber}-compute@developer.gserviceaccount.com`);
  }

  // Removed non-existent appspot service account
  // if (info.projectId) {
  //   derived.add(`${info.projectId}@appspot.gserviceaccount.com`);
  // }

  const invokers = Array.from(derived).filter(Boolean);
  if (invokers.length === 0) {
    return "private";
  }

  return invokers;
}
*/

// Initialize Gemini
const genAI = new GoogleGenerativeAI(process.env.GEMINI_API_KEY || "");

// Cloud Run API endpoint (our main backend)
const CLOUD_RUN_API = process.env.CLOUD_RUN_API_URL || "https://driftlock-api-o6kjgrsowq-uc.a.run.app";

// Proxy signup requests to Cloud Run backend
export const signup = onRequest({ cors: true }, async (request, response) => {
  logger.info("Signup request received", { structuredData: true });

  if (request.method !== "POST") {
    response.status(405).json({ error: "Method not allowed" });
    return;
  }

  try {
    const { email, company_name } = request.body;

    // Forward to Cloud Run backend
    const backendResponse = await fetch(`${CLOUD_RUN_API}/v1/onboard/signup`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email,
        company_name,
        // signup_source: 'firebase_frontend' // Backend does not support this field
      }),
    });

    const result = await backendResponse.json();

    if (backendResponse.ok) {
      response.json(result);
    } else {
      response.status(backendResponse.status).json(result);
    }
  } catch (error) {
    logger.error("Signup error", error);
    response.status(500).json({ error: "Signup failed" });
  }
});

export * from "./genkit";

// Enhanced anomaly analysis with Gemini
export const analyzeAnomalies = onRequest({ cors: true }, async (request, response) => {
  logger.info("Anomaly analysis request", { structuredData: true });

  if (request.method !== "POST") {
    response.status(405).json({ error: "Method not allowed" });
    return;
  }

  try {
    const { anomalies, query, api_key } = request.body;

    if (!anomalies || !Array.isArray(anomalies)) {
      response.status(400).json({ error: "Invalid anomalies data" });
      return;
    }

    // Verify API key with backend if provided
    if (api_key) {
      try {
        const authResponse = await fetch(`${CLOUD_RUN_API}/v1/anomalies`, {
          method: "GET",
          headers: {
            "Authorization": `Bearer ${api_key}`,
            "Content-Type": "application/json",
          },
        });

        if (!authResponse.ok) {
          response.status(401).json({ error: "Invalid API key" });
          return;
        }
      } catch (authError) {
        logger.warn("API key validation failed", authError);
      }
    }

    const model = genAI.getGenerativeModel({ model: "gemini-pro" });

    const prompt = `
      Analyze these anomalies detected by Driftlock's compression-based system:
      
      ${anomalies.map((anomaly: any) => `
      - ID: ${anomaly.id}
      - NCD Score: ${anomaly.ncd_score || "N/A"} 
      - P-value: ${anomaly.p_value || "N/A"}
      - Explanation: ${anomaly.explanation || "No explanation available"}
      - Data Type: ${anomaly.stream_type || "unknown"}
      - Detected: ${anomaly.detected_at || "recent"}
      `).join("\n")}
      
      User Query: ${query || "Provide insights about these anomalies"}
      
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

    response.json({
      success: true,
      analysis,
      processed_anomalies: anomalies.length,
      timestamp: new Date().toISOString(),
      confidence: "high",
    });
  } catch (error) {
    logger.error("Error analyzing anomalies", error);
    response.status(500).json({ error: "Analysis failed" });
  }
});

// Generate compliance reports
export const generateComplianceReport = onRequest({ cors: true }, async (request, response) => {
  logger.info("Compliance report generation", { structuredData: true });

  if (request.method !== "POST") {
    response.status(405).json({ error: "Method not allowed" });
    return;
  }

  try {
    const { anomalies, regulation, tenant_info, api_key } = request.body;

    // Verify API key if provided
    if (api_key) {
      try {
        const authResponse = await fetch(`${CLOUD_RUN_API}/v1/anomalies`, {
          method: "GET",
          headers: {
            "Authorization": `Bearer ${api_key}`,
            "Content-Type": "application/json",
          },
        });

        if (!authResponse.ok) {
          response.status(401).json({ error: "Invalid API key" });
          return;
        }
      } catch (authError) {
        logger.warn("API key validation failed", authError);
      }
    }

    const model = genAI.getGenerativeModel({ model: "gemini-pro" });

    const prompt = `
      Generate a ${regulation || "DORA"} compliance report for these anomalies:
      
      Company: ${tenant_info?.company_name || "Customer"}
      Report Date: ${new Date().toLocaleDateString()}
      
      Detected Anomalies:
      ${anomalies.map((anomaly: any) => `
      - Anomaly ID: ${anomaly.id}
      - Detection Time: ${anomaly.detected_at}
      - NCD Score: ${anomaly.ncd_score} (Mathematical Evidence)
      - Statistical Significance: P-value ${anomaly.p_value}
      - Technical Explanation: ${anomaly.explanation}
      - Risk Level: ${anomaly.ncd_score > 0.7 ? "HIGH" : anomaly.ncd_score > 0.4 ? "MEDIUM" : "LOW"}
      `).join("\n")}
      
      Generate a formal compliance report including:
      1. Executive Summary with business impact
      2. Technical Analysis with mathematical evidence (NCD, p-values)
      3. Risk Assessment per ${regulation || "DORA"} requirements  
      4. Audit Trail and Evidence Documentation
      5. Recommended Immediate Actions
      6. Long-term Monitoring Recommendations
      
      Format as a professional regulatory document suitable for auditor review.
      Include references to compression-based anomaly detection methodology.
    `;

    const result = await model.generateContent(prompt);
    const report = result.response.text();

    response.json({
      success: true,
      report,
      regulation: regulation || "DORA",
      anomaly_count: anomalies.length,
      generated_at: new Date().toISOString(),
      company: tenant_info?.company_name || "Customer",
    });
  } catch (error) {
    logger.error("Error generating compliance report", error);
    response.status(500).json({ error: "Report generation failed" });
  }
});

// Proxy API requests to Cloud Run backend
export const apiProxy = onRequest({ cors: true }, async (request, response) => {
  logger.info("API proxy request", { path: request.path, method: request.method });

  try {
    let apiPath = request.path;

    // Handle specific rewrites first
    if (request.path === "/webhooks/stripe") {
      apiPath = "/v1/billing/webhook";
    } else if (request.path === "/api/v1/healthz" || request.path === "/healthz") {
      // Health check should go directly to backend /healthz (no auth required)
      apiPath = "/healthz";
    } else if (request.path.startsWith("/api/proxy")) {
      apiPath = request.path.replace("/api/proxy", "");
    } else if (request.path.startsWith("/api/v1")) {
      apiPath = request.path.replace("/api", "");
    }

    const backendUrl = `${CLOUD_RUN_API}${apiPath}`;

    // Determine headers to forward
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    // Forward Authorization header
    if (request.headers.authorization) {
      headers["Authorization"] = request.headers.authorization;
    }

    // Forward Stripe headers for webhooks
    if (request.headers["stripe-signature"]) {
      headers["Stripe-Signature"] = request.headers["stripe-signature"] as string;
    }

    // Forward X-Api-Key if present
    if (request.headers["x-api-key"]) {
      headers["X-Api-Key"] = request.headers["x-api-key"] as string;
    }

    // For webhooks, we need raw body sometimes, but fetch takes body as string/buffer
    // request.body in Firebase Functions is already parsed if JSON
    // BUT Stripe webhooks need raw body for signature verification.
    // Firebase Functions v2 usually provides rawBody buffer if rawBody is enabled?
    // Actually onRequest by default parses body.
    // For webhooks verification we might need the raw buffer.
    // However, passing JSON.stringify(request.body) might break signature verification
    // if the parsing/stringifying changes the byte order.
    // Since we are acting as a proxy, we should ideally forward the raw bytes.

    let body: any;
    if (request.method !== "GET" && request.method !== "HEAD") {
      // CRITICAL: For Stripe webhooks, we MUST use the raw request body buffer.
      // JSON.stringify() re-serializes the parsed body, potentially changing key order
      // or formatting, which invalidates the cryptographic signature.
      if (request.path === "/webhooks/stripe" && (request as any).rawBody) {
        body = (request as any).rawBody;
      } else {
        // For other endpoints, re-serializing is generally fine (or we could use rawBody if available)
        // usage of rawBody is safer for proxying in general to avoid any parsing/re-parsing issues
        body = (request as any).rawBody ? (request as any).rawBody : JSON.stringify(request.body);
      }
    }

    const backendResponse = await fetch(backendUrl, {
      method: request.method,
      headers: headers,
      body: body,
    });

    // Check content type of response
    const contentType = backendResponse.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      const result = await backendResponse.json();
      response.status(backendResponse.status).json(result);
    } else {
      const text = await backendResponse.text();
      response.status(backendResponse.status).send(text);
    }
  } catch (error) {
    logger.error("API proxy error", error);
    response.status(500).json({ error: "Backend service unavailable" });
  }
});

// Health check for the entire stack
export const healthCheck = onRequest({ cors: true }, async (request, response) => {
  const health: { [key: string]: any } = {
    success: true,
    status: "healthy",
    service: "driftlock-saas-backend",
    timestamp: new Date().toISOString(),
    version: "2.0.0",
    features: [
      "user-signup",
      "anomaly-analysis",
      "compliance-reporting",
      "gemini-integration",
      "cloud-run-proxy",
    ],
    backend: { status: "unknown" },
  };

  // Check Cloud Run backend health
  try {
    const candidates = [
      `${CLOUD_RUN_API}/healthz`,
      `${CLOUD_RUN_API}/v1/healthz`,
    ];

    let backendResponse: any = null;

    for (const url of candidates) {
      const attempt = await fetch(url, { method: "GET" });
      if (attempt.ok) {
        backendResponse = attempt;
        health.backend.checked = url;
        break;
      }

      // Keep last response for context if all fail
      backendResponse = attempt;
    }

    if (backendResponse && backendResponse.ok) {
      const backendHealth = await backendResponse.json();
      health.backend = {
        status: "healthy",
        database: backendHealth.database || "unknown",
        license: backendHealth.license ? "valid" : "unknown",
      };
      // Ensure success is true if backend is healthy
      health.success = backendHealth.success !== false;
    } else {
      const code = backendResponse?.status;
      health.backend = { status: code ? `unhealthy (${code})` : "unhealthy" };
      health.success = false;
    }
  } catch (error) {
    health.backend = { status: "unreachable" };
    health.success = false;
  }

  response.json(health);
});

export const getFirebaseConfig = onRequest({ cors: true }, async (request, response) => {
  logger.info("Firebase config request received", { structuredData: true });

  try {
    const secretManagerClient = new SecretManagerServiceClient();
    const secretName = "projects/driftlock/secrets/VITE_FIREBASE_API_KEY/versions/latest";

    let apiKey: string | undefined = process.env.VITE_FIREBASE_API_KEY;

    if (!apiKey) {
      const [version] = await secretManagerClient.accessSecretVersion({
        name: secretName,
      });

      apiKey = version.payload?.data?.toString();
    }

    if (!apiKey) {
      throw new Error("API key not found in Secret Manager or env.");
    }

    const firebaseConfig = {
      apiKey: apiKey,
      authDomain: process.env.VITE_FIREBASE_AUTH_DOMAIN || "driftlock.firebaseapp.com",
      projectId: process.env.VITE_FIREBASE_PROJECT_ID || "driftlock",
      storageBucket: process.env.VITE_FIREBASE_STORAGE_BUCKET || "driftlock.appspot.com",
      messagingSenderId: process.env.VITE_FIREBASE_MESSAGING_SENDER_ID || "131489574303",
      appId: process.env.VITE_FIREBASE_APP_ID || "1:131489574303:web:e83e3e433912d05a8d61aa",
      measurementId: process.env.VITE_FIREBASE_MEASUREMENT_ID || "G-CXBMVS3G8H",
    };

    response.json(firebaseConfig);
  } catch (error) {
    logger.error("Error getting Firebase config", error);
    response.status(500).json({ error: "Could not retrieve Firebase configuration." });
  }
});
