/**
 * Driftlock Firebase Functions - SaaS Backend
 * Integrated with Cloud Run API and Gemini AI
 * Updated: Force redeploy for invoker config (attempt 3)
 */

import { setGlobalOptions } from "firebase-functions";
import { onRequest } from "firebase-functions/v2/https";
import { logger } from "firebase-functions";
import { GoogleGenerativeAI } from "@google/generative-ai";
import { SecretManagerServiceClient } from "@google-cloud/secret-manager";

type InvokerConfig = "public" | "private" | string | string[];

interface ProjectInfo {
  projectId?: string;
  projectNumber?: string;
}

const projectInfo = resolveProjectInfo();
const invokerConfig = resolveInvokers(projectInfo);

// Set global options for cost control
setGlobalOptions({ maxInstances: 10, invoker: invokerConfig });
logger.info("Configured function invokers", {
  invokerConfig,
  projectInfo,
});

function parseInvokerList(value?: string): string[] {
  if (!value) return [];
  return value
    .split(/[,;\s]+/)
    .map((part) => part.trim())
    .filter((part) => part.length > 0);
}

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

  if (info.projectId) {
    derived.add(`${info.projectId}@appspot.gserviceaccount.com`);
  }

  const invokers = Array.from(derived).filter(Boolean);
  if (invokers.length === 0) {
    return "private";
  }

  return invokers;
}

// Initialize Gemini
const genAI = new GoogleGenerativeAI(process.env.GEMINI_API_KEY || "");

// Cloud Run API endpoint (our main backend)
const CLOUD_RUN_API = process.env.CLOUD_RUN_API_URL || "https://driftlock-api-run.a.app";

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
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        company_name,
        signup_source: 'firebase_frontend'
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
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${api_key}`,
            'Content-Type': 'application/json',
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
      - NCD Score: ${anomaly.ncd_score || 'N/A'} 
      - P-value: ${anomaly.p_value || 'N/A'}
      - Explanation: ${anomaly.explanation || 'No explanation available'}
      - Data Type: ${anomaly.stream_type || 'unknown'}
      - Detected: ${anomaly.detected_at || 'recent'}
      `).join('\n')}
      
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
      confidence: "high"
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
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${api_key}`,
            'Content-Type': 'application/json',
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

    response.json({
      success: true,
      report,
      regulation: regulation || 'DORA',
      anomaly_count: anomalies.length,
      generated_at: new Date().toISOString(),
      company: tenant_info?.company_name || 'Customer'
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
    const apiPath = request.path.replace('/api/proxy', '');
    const backendUrl = `${CLOUD_RUN_API}${apiPath}`;

    const backendResponse = await fetch(backendUrl, {
      method: request.method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': request.headers.authorization || '',
      },
      body: request.method !== 'GET' ? JSON.stringify(request.body) : undefined,
    });

    const result = await backendResponse.json();
    response.status(backendResponse.status).json(result);

  } catch (error) {
    logger.error("API proxy error", error);
    response.status(500).json({ error: "Backend service unavailable" });
  }
});

// Health check for the entire stack
export const healthCheck = onRequest({ cors: true }, async (request, response) => {
  const health: { [key: string]: any } = {
    status: "healthy",
    service: "driftlock-saas-backend",
    timestamp: new Date().toISOString(),
    version: "2.0.0",
    features: [
      "user-signup",
      "anomaly-analysis",
      "compliance-reporting",
      "gemini-integration",
      "cloud-run-proxy"
    ],
    backend: { status: "unknown" }
  };

  // Check Cloud Run backend health
  try {
    const backendResponse = await fetch(`${CLOUD_RUN_API}/healthz`, {
      method: 'GET',
    });

    if (backendResponse.ok) {
      const backendHealth = await backendResponse.json();
      health.backend = {
        status: "healthy",
        database: backendHealth.database || "unknown",
        license: backendHealth.license ? "valid" : "unknown"
      };
    } else {
      health.backend = { status: "unhealthy" };
    }
  } catch (error) {
    health.backend = { status: "unreachable" };
  }

  response.json(health);
});

export const getFirebaseConfig = onRequest({ cors: true }, async (request, response) => {
  logger.info("Firebase config request received", { structuredData: true });

  try {
    const secretManagerClient = new SecretManagerServiceClient();
    const secretName = "projects/driftlock/secrets/VITE_FIREBASE_API_KEY/versions/latest";

    const [version] = await secretManagerClient.accessSecretVersion({
      name: secretName,
    });

    const apiKey = version.payload?.data?.toString();

    if (!apiKey) {
      throw new Error("API key not found in Secret Manager.");
    }

    const firebaseConfig = {
      apiKey: apiKey,
      authDomain: "driftlock.firebaseapp.com",
      projectId: "driftlock",
      storageBucket: "driftlock.appspot.com",
      messagingSenderId: "131489574303",
      appId: "1:131489574303:web:e83e3e433912d05a8d61aa",
      measurementId: "G-CXBMVS3G8H",
    };

    response.json(firebaseConfig);
  } catch (error) {
    logger.error("Error getting Firebase config", error);
    response.status(500).json({ error: "Could not retrieve Firebase configuration." });
  }
});
