/**
 * Smoke Test - Quick sanity check
 *
 * Purpose: Verify all endpoints return expected status codes
 * Duration: 1 minute
 * VUs: 5
 *
 * Usage: k6 run scripts/load-test/smoke.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const healthCheckSuccess = new Rate('health_check_success');
const demoDetectSuccess = new Rate('demo_detect_success');

// Test configuration
export const options = {
  vus: 5,
  duration: '1m',
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'],
    'http_req_failed': ['rate<0.01'],
    'health_check_success': ['rate>0.99'],
    'demo_detect_success': ['rate>0.95'],
  },
  tags: {
    test_type: 'smoke',
  },
};

const BASE_URL = helpers.getBaseUrl();

export function setup() {
  console.log(`Starting smoke test against ${BASE_URL}`);
  console.log('Testing basic endpoint functionality...');

  // Verify server is running
  const health = http.get(`${BASE_URL}/healthz`);
  if (health.status !== 200) {
    throw new Error(`Server not healthy: ${health.status}`);
  }

  return { baseUrl: BASE_URL };
}

export default function (data) {
  const baseUrl = data.baseUrl;

  // Test 1: Health check
  const healthRes = http.get(`${baseUrl}/healthz`);
  const healthOk = check(healthRes, {
    'health check returns 200': (r) => r.status === 200,
    'health check has ok status': (r) => r.body.includes('ok') || r.status === 200,
  });
  healthCheckSuccess.add(healthOk);

  sleep(0.5);

  // Test 2: Readiness check
  const readyRes = http.get(`${baseUrl}/readyz`);
  check(readyRes, {
    'readiness check returns 200': (r) => r.status === 200,
  });

  sleep(0.5);

  // Test 3: Demo detect endpoint (small batch)
  const events = helpers.generateEventBatch(10, 0.2);
  const detectPayload = helpers.createDetectPayload(events);

  const detectRes = http.post(
    `${baseUrl}/v1/demo/detect`,
    detectPayload,
    { headers: helpers.demoHeaders() }
  );

  const detectOk = check(detectRes, {
    'demo detect returns 200': (r) => r.status === 200,
    'demo detect has processed count': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.hasOwnProperty('processed');
      } catch {
        return false;
      }
    },
    'demo detect has anomalies array': (r) => {
      try {
        const body = JSON.parse(r.body);
        return Array.isArray(body.anomalies);
      } catch {
        return false;
      }
    },
  });
  demoDetectSuccess.add(detectOk);

  sleep(1);

  // Test 4: Invalid endpoint returns 404
  const notFoundRes = http.get(`${baseUrl}/v1/nonexistent`);
  check(notFoundRes, {
    'invalid endpoint returns 404': (r) => r.status === 404,
  });

  sleep(1);
}

export function teardown(data) {
  console.log('Smoke test complete!');
  console.log(`Tested against: ${data.baseUrl}`);
}
