// Driftlock Load Testing Script
// Run with: k6 run scripts/load-test.js
// Or: k6 run --vus 10 --duration 30s scripts/load-test.js

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const healthLatency = new Trend('health_latency');
const detectLatency = new Trend('detect_latency');

// Configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 10 },   // Stay at 10 users
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests should be < 500ms
    errors: ['rate<0.01'],              // Error rate should be < 1%
    health_latency: ['p(95)<200'],      // Health check p95 < 200ms
    detect_latency: ['p(95)<5000'],     // Detection p95 < 5s
  },
};

// Get configuration from environment
const BASE_URL = __ENV.API_URL || 'http://localhost:8080';
const API_KEY = __ENV.API_KEY || 'dlk_test-key';


// Generate baseline events (normal pattern)
function generateBaselineEvents(count) {
  const events = [];
  for (let i = 0; i < count; i++) {
    events.push({
      timestamp: new Date(Date.now() - (count - i) * 1000).toISOString(),
      level: 'info',
      message: `Normal operation ${i}`,
      user_id: `user_${Math.floor(Math.random() * 100)}`,
      duration_ms: Math.floor(Math.random() * 100) + 10,
    });
  }
  return events;
}

// Generate events with potential anomaly
function generateEventsWithAnomaly(count) {
  const events = generateBaselineEvents(count - 5);

  // Add anomalous events at the end
  for (let i = 0; i < 5; i++) {
    events.push({
      timestamp: new Date().toISOString(),
      level: 'error',
      message: `CRITICAL: Unusual activity detected - ${i}`,
      user_id: 'unknown',
      duration_ms: Math.floor(Math.random() * 10000) + 5000,
      anomaly_score: 0.95,
    });
  }
  return events;
}

export default function () {
  // Test 1: Health check
  const healthRes = http.get(`${BASE_URL}/healthz`);
  healthLatency.add(healthRes.timings.duration);

  const healthCheck = check(healthRes, {
    'health status is 200': (r) => r.status === 200,
    'health response is success': (r) => {
      try {
        return JSON.parse(r.body).success === true;
      } catch {
        return false;
      }
    },
  });

  if (!healthCheck) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }

  sleep(0.1);

  // Test 2: Anomaly detection (if API key provided)
  if (API_KEY && API_KEY !== 'dlk_test-key') {
    const events = Math.random() > 0.8
      ? generateEventsWithAnomaly(450)
      : generateBaselineEvents(450);

    const detectRes = http.post(
      `${BASE_URL}/v1/detect`,
      JSON.stringify({
        stream_id: 'default',
        events: events,
        config_override: {
          baseline_size: 400,
          window_size: 50,
        },
      }),
      {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${API_KEY}`,
        },
        timeout: '30s',
      }
    );

    detectLatency.add(detectRes.timings.duration);

    const detectCheck = check(detectRes, {
      'detect status is 200': (r) => r.status === 200,
      'detect response has success': (r) => {
        try {
          return JSON.parse(r.body).success === true;
        } catch {
          return false;
        }
      },
      'detect response has batch_id': (r) => {
        try {
          return JSON.parse(r.body).batch_id !== undefined;
        } catch {
          return false;
        }
      },
    });

    if (!detectCheck) {
      errorRate.add(1);
      console.log(`Detect failed: ${detectRes.status} - ${detectRes.body}`);
    } else {
      errorRate.add(0);
    }
  }

  sleep(0.5);
}

// Summary handler
export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'load-test-results.json': JSON.stringify(data, null, 2),
  };
}

// Text summary helper
function textSummary(data, options) {
  const metrics = data.metrics;

  let output = '\n========== LOAD TEST SUMMARY ==========\n\n';

  output += 'HTTP Requests:\n';
  output += `  Total: ${metrics.http_reqs?.values?.count || 0}\n`;
  output += `  Failed: ${metrics.http_req_failed?.values?.passes || 0}\n`;
  output += `  Duration (p95): ${(metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2)}ms\n\n`;

  output += 'Custom Metrics:\n';
  output += `  Error Rate: ${((metrics.errors?.values?.rate || 0) * 100).toFixed(2)}%\n`;
  output += `  Health Latency (p95): ${(metrics.health_latency?.values?.['p(95)'] || 0).toFixed(2)}ms\n`;
  output += `  Detect Latency (p95): ${(metrics.detect_latency?.values?.['p(95)'] || 0).toFixed(2)}ms\n\n`;

  output += 'Thresholds:\n';
  for (const [name, threshold] of Object.entries(data.thresholds || {})) {
    const status = threshold.ok ? '✓ PASS' : '✗ FAIL';
    output += `  ${name}: ${status}\n`;
  }

  output += '\n========================================\n';

  return output;
}
