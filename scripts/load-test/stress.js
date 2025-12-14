/**
 * Stress Test - Breaking point test
 *
 * Purpose: Push the system to its limits to find breaking points
 * Duration: 10 minutes
 * VUs: 0 -> 50 -> 100 -> 150 -> 200 -> 0
 *
 * Usage: k6 run scripts/load-test/stress.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend, Rate } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const detectLatency = new Trend('detect_latency', true);
const eventsProcessed = new Counter('events_processed');
const serverErrors = new Counter('server_errors');
const timeouts = new Counter('timeouts');
const errorRate = new Rate('error_rate');

// Test configuration
export const options = {
  stages: [
    { duration: '1m', target: 50 },   // Warm up
    { duration: '2m', target: 100 },  // Increase load
    { duration: '2m', target: 150 },  // Push harder
    { duration: '2m', target: 200 },  // Maximum stress
    { duration: '2m', target: 150 },  // Back off slightly
    { duration: '1m', target: 0 },    // Cooldown
  ],
  thresholds: {
    'http_req_duration': ['p(95)<1000', 'p(99)<2000'],
    'http_req_failed': ['rate<0.05'], // Allow higher error rate in stress test
    'detect_latency': ['p(95)<800', 'p(99)<1500'],
    'error_rate': ['rate<0.1'],
  },
  tags: {
    test_type: 'stress',
  },
};

const BASE_URL = helpers.getBaseUrl();

export function setup() {
  console.log(`Starting stress test against ${BASE_URL}`);
  console.log('WARNING: This test will push the system to its limits!');

  // Verify server is ready
  const health = http.get(`${BASE_URL}/healthz`);
  if (health.status !== 200) {
    throw new Error(`Server not healthy: ${health.status}`);
  }

  return {
    baseUrl: BASE_URL,
    startTime: Date.now(),
    errorThreshold: 0,
  };
}

export default function (data) {
  const baseUrl = data.baseUrl;

  // Generate aggressive load
  const batchSize = Math.floor(Math.random() * 800) + 200; // 200-1000 events
  const events = helpers.generateMixedBatch(batchSize);
  const streamId = helpers.generateUniqueStreamId(__VU, __ITER);
  const payload = helpers.createDetectPayload(events, streamId, 'strict');

  const startTime = Date.now();
  const res = http.post(
    `${baseUrl}/v1/demo/detect`,
    payload,
    {
      headers: helpers.demoHeaders(),
      tags: { endpoint: 'detect', stress_level: 'high' },
      timeout: '10s', // Increased timeout for stress conditions
    }
  );
  const duration = Date.now() - startTime;

  const hasError = !check(res, {
    'status is not 5xx': (r) => r.status < 500,
    'status is not timeout': (r) => r.status !== 0,
    'response time is acceptable': (r) => r.timings.duration < 5000,
  });

  errorRate.add(hasError);

  if (res.status === 0) {
    timeouts.add(1);
    console.warn(`Request timeout at VU=${__VU}, iter=${__ITER}`);
  } else if (res.status >= 500) {
    serverErrors.add(1);
    console.warn(`Server error ${res.status} at VU=${__VU}, iter=${__ITER}`);
  } else if (res.status === 200) {
    try {
      const body = JSON.parse(res.body);

      if (body.detection_time_ms) {
        detectLatency.add(body.detection_time_ms);
      }

      if (body.processed) {
        eventsProcessed.add(body.processed);
      }
    } catch (e) {
      console.error(`Failed to parse response: ${e.message}`);
    }
  }

  // Very short sleep to maintain pressure
  helpers.sleepWithJitter(0.5, 0.2);
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000;
  console.log(`Stress test complete! Duration: ${duration.toFixed(2)}s`);
  console.log(`Tested against: ${data.baseUrl}`);
  console.log('Check metrics to identify breaking points.');
}

export function handleSummary(data) {
  const summary = {
    'stress_test_results.json': JSON.stringify(data, null, 2),
  };

  // Log key insights
  console.log('\n=== Stress Test Summary ===');
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Failed requests: ${data.metrics.http_req_failed.values.rate * 100}%`);
  console.log(`P95 latency: ${data.metrics.http_req_duration.values['p(95)']}ms`);
  console.log(`P99 latency: ${data.metrics.http_req_duration.values['p(99)']}ms`);

  if (data.metrics.server_errors) {
    console.log(`Server errors: ${data.metrics.server_errors.values.count}`);
  }
  if (data.metrics.timeouts) {
    console.log(`Timeouts: ${data.metrics.timeouts.values.count}`);
  }

  return summary;
}
