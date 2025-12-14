/**
 * K6 Load Test Helpers for Driftlock API
 *
 * Provides event generators, auth utilities, and common test utilities.
 */

import { randomString, randomIntBetween, randomItem } from 'k6';

/**
 * Generate realistic log event
 */
export function generateLogEvent(streamId = null) {
  const levels = ['INFO', 'WARN', 'ERROR', 'DEBUG'];
  const services = ['api-server', 'worker', 'collector', 'billing'];
  const messages = [
    'Request processed successfully',
    'Database connection established',
    'Cache hit for key',
    'Background job completed',
    'Webhook delivered',
    'User authenticated',
    'Payment processed',
    'Anomaly detected',
  ];

  return {
    timestamp: new Date().toISOString(),
    level: randomItem(levels),
    service: randomItem(services),
    message: randomItem(messages),
    trace_id: randomString(16),
    span_id: randomString(8),
    stream_id: streamId || `stream-${randomIntBetween(1, 100)}`,
  };
}

/**
 * Generate anomalous log event (high entropy)
 */
export function generateAnomalousEvent(streamId = null) {
  return {
    timestamp: new Date().toISOString(),
    level: 'ERROR',
    service: 'unknown-service-' + randomString(10),
    message: 'CRITICAL: ' + randomString(50),
    trace_id: randomString(32),
    span_id: randomString(16),
    stream_id: streamId || `stream-${randomIntBetween(1, 100)}`,
    error: {
      type: 'UnexpectedError-' + randomString(10),
      stack: randomString(200),
      code: randomIntBetween(1000, 9999),
    },
  };
}

/**
 * Generate metric event
 */
export function generateMetricEvent(streamId = null) {
  const metricNames = [
    'cpu_usage_percent',
    'memory_bytes',
    'request_duration_ms',
    'http_requests_total',
    'error_rate',
  ];

  return {
    timestamp: new Date().toISOString(),
    name: randomItem(metricNames),
    value: Math.random() * 100,
    unit: 'percent',
    stream_id: streamId || `metrics-${randomIntBetween(1, 50)}`,
    tags: {
      host: `server-${randomIntBetween(1, 10)}`,
      region: randomItem(['us-east', 'us-west', 'eu-central']),
    },
  };
}

/**
 * Generate batch of events
 */
export function generateEventBatch(size, anomalyRate = 0.1, streamId = null) {
  const events = [];
  for (let i = 0; i < size; i++) {
    if (Math.random() < anomalyRate) {
      events.push(generateAnomalousEvent(streamId));
    } else {
      events.push(generateLogEvent(streamId));
    }
  }
  return events;
}

/**
 * Generate mixed event batch (logs, metrics, traces)
 */
export function generateMixedBatch(size, streamId = null) {
  const events = [];
  for (let i = 0; i < size; i++) {
    const type = Math.random();
    if (type < 0.6) {
      events.push(generateLogEvent(streamId));
    } else if (type < 0.9) {
      events.push(generateMetricEvent(streamId));
    } else {
      events.push(generateAnomalousEvent(streamId));
    }
  }
  return events;
}

/**
 * Create API key auth headers
 */
export function authHeaders(apiKey) {
  return {
    'Content-Type': 'application/json',
    'X-API-Key': apiKey,
  };
}

/**
 * Create demo request headers
 */
export function demoHeaders() {
  return {
    'Content-Type': 'application/json',
  };
}

/**
 * Get base URL from environment or default
 */
export function getBaseUrl() {
  return __ENV.BASE_URL || 'http://localhost:8080';
}

/**
 * Get API key from environment or use test key
 */
export function getApiKey() {
  return __ENV.API_KEY || 'test-api-key-12345';
}

/**
 * Create detect request payload
 */
export function createDetectPayload(events, streamId = null, profile = 'balanced') {
  return JSON.stringify({
    events: events,
    stream_id: streamId,
    profile: profile,
    return_anomalies: true,
  });
}

/**
 * Parse and validate detect response
 */
export function validateDetectResponse(response) {
  if (response.status !== 200) {
    return {
      valid: false,
      error: `Unexpected status: ${response.status}`,
    };
  }

  try {
    const body = JSON.parse(response.body);
    if (!body.hasOwnProperty('processed') || !body.hasOwnProperty('anomalies')) {
      return {
        valid: false,
        error: 'Missing required fields in response',
      };
    }

    return {
      valid: true,
      processed: body.processed,
      anomalies: body.anomalies,
      detectionTime: body.detection_time_ms || 0,
    };
  } catch (e) {
    return {
      valid: false,
      error: `Failed to parse response: ${e.message}`,
    };
  }
}

/**
 * Generate unique stream IDs for capacity testing
 */
export function generateUniqueStreamId(vuId, iteration) {
  return `load-test-stream-${vuId}-${iteration}-${Date.now()}`;
}

/**
 * Sleep with jitter to avoid thundering herd
 */
export function sleepWithJitter(baseSeconds, jitterSeconds = 1) {
  const sleep = require('k6').sleep;
  const jitter = Math.random() * jitterSeconds;
  sleep(baseSeconds + jitter);
}

/**
 * Custom metrics helpers
 */
export function recordDetectionMetrics(metrics, response, startTime) {
  const duration = Date.now() - startTime;

  if (response.status === 200) {
    const body = JSON.parse(response.body);
    if (body.detection_time_ms) {
      metrics.detectionLatency.add(body.detection_time_ms);
    }
    if (body.processed) {
      metrics.eventsProcessed.add(body.processed);
    }
    if (body.anomalies) {
      metrics.anomaliesDetected.add(body.anomalies.length);
    }
  }

  metrics.requestDuration.add(duration);
}

/**
 * Check if response indicates rate limiting
 */
export function isRateLimited(response) {
  return response.status === 429;
}

/**
 * Generate transaction data (for fraud detection scenarios)
 */
export function generateTransaction(isFraudulent = false) {
  if (isFraudulent) {
    return {
      timestamp: new Date().toISOString(),
      transaction_id: randomString(16),
      amount: randomIntBetween(5000, 50000), // High amount
      user_id: randomString(8),
      merchant: 'SUSPICIOUS-' + randomString(10),
      location: 'Unknown',
      device_id: randomString(32),
      ip_address: `${randomIntBetween(1, 255)}.${randomIntBetween(1, 255)}.${randomIntBetween(1, 255)}.${randomIntBetween(1, 255)}`,
      velocity: randomIntBetween(10, 50), // High velocity
    };
  }

  return {
    timestamp: new Date().toISOString(),
    transaction_id: randomString(16),
    amount: randomIntBetween(10, 500),
    user_id: `user-${randomIntBetween(1, 1000)}`,
    merchant: randomItem(['Amazon', 'Walmart', 'Target', 'Costco']),
    location: randomItem(['US', 'CA', 'UK', 'DE']),
    device_id: `device-${randomIntBetween(1, 100)}`,
    ip_address: `192.168.${randomIntBetween(1, 255)}.${randomIntBetween(1, 255)}`,
    velocity: randomIntBetween(1, 5),
  };
}

/**
 * Generate LLM I/O event
 */
export function generateLLMEvent(streamId = null) {
  const models = ['gpt-4', 'gpt-3.5-turbo', 'claude-3', 'gemini-pro'];
  const prompts = [
    'Write a function to',
    'Explain how',
    'What is the best way to',
    'Debug this code:',
  ];

  return {
    timestamp: new Date().toISOString(),
    model: randomItem(models),
    prompt: randomItem(prompts) + ' ' + randomString(20),
    completion: randomString(100),
    tokens_input: randomIntBetween(10, 500),
    tokens_output: randomIntBetween(50, 1000),
    latency_ms: randomIntBetween(100, 5000),
    stream_id: streamId || `llm-${randomIntBetween(1, 20)}`,
  };
}

export default {
  generateLogEvent,
  generateAnomalousEvent,
  generateMetricEvent,
  generateEventBatch,
  generateMixedBatch,
  generateTransaction,
  generateLLMEvent,
  authHeaders,
  demoHeaders,
  getBaseUrl,
  getApiKey,
  createDetectPayload,
  validateDetectResponse,
  generateUniqueStreamId,
  sleepWithJitter,
  recordDetectionMetrics,
  isRateLimited,
};
