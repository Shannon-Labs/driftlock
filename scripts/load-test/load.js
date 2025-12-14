/**
 * Load Test - Production simulation
 *
 * Purpose: Simulate realistic production traffic patterns
 * Duration: 15 minutes (5 min ramp-up, 5 min sustained, 5 min cooldown)
 * VUs: 0 -> 25 -> 50 -> 100 -> 50 -> 0
 *
 * Usage: k6 run scripts/load-test/load.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend, Rate } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const detectLatency = new Trend('detect_latency', true);
const eventsProcessed = new Counter('events_processed');
const anomaliesDetected = new Counter('anomalies_detected');
const detectSuccess = new Rate('detect_success');

// Test configuration
export const options = {
  stages: [
    { duration: '1m', target: 25 },  // Ramp up to 25 VUs
    { duration: '2m', target: 50 },  // Ramp up to 50 VUs
    { duration: '2m', target: 100 }, // Ramp up to 100 VUs
    { duration: '5m', target: 100 }, // Sustain 100 VUs
    { duration: '3m', target: 50 },  // Ramp down to 50 VUs
    { duration: '2m', target: 0 },   // Cooldown
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000', 'p(99.9)<2000'],
    'http_req_failed': ['rate<0.01'],
    'http_reqs': ['rate>100'],
    'detect_latency': ['p(95)<400', 'p(99)<800'],
    'detect_success': ['rate>0.99'],
    'events_processed': ['count>10000'],
  },
  tags: {
    test_type: 'load',
  },
};

const BASE_URL = helpers.getBaseUrl();

export function setup() {
  console.log(`Starting load test against ${BASE_URL}`);
  console.log('Simulating production traffic patterns...');

  // Verify server is ready
  const health = http.get(`${BASE_URL}/healthz`);
  if (health.status !== 200) {
    throw new Error(`Server not healthy: ${health.status}`);
  }

  const ready = http.get(`${BASE_URL}/readyz`);
  if (ready.status !== 200) {
    throw new Error(`Server not ready: ${ready.status}`);
  }

  return {
    baseUrl: BASE_URL,
    startTime: Date.now(),
  };
}

export default function (data) {
  const baseUrl = data.baseUrl;
  const scenario = Math.random();

  // 70% - Normal detection requests
  if (scenario < 0.7) {
    const batchSize = Math.floor(Math.random() * 100) + 10; // 10-110 events
    const events = helpers.generateEventBatch(batchSize, 0.1);
    const streamId = `stream-${__VU}-${Math.floor(__ITER / 10)}`;
    const payload = helpers.createDetectPayload(events, streamId);

    const startTime = Date.now();
    const res = http.post(
      `${baseUrl}/v1/demo/detect`,
      payload,
      {
        headers: helpers.demoHeaders(),
        tags: { endpoint: 'detect', scenario: 'normal' },
      }
    );

    const success = check(res, {
      'detect status is 200': (r) => r.status === 200,
      'detect response is valid': (r) => {
        const result = helpers.validateDetectResponse(r);
        return result.valid;
      },
    });

    detectSuccess.add(success);

    if (res.status === 200) {
      try {
        const body = JSON.parse(res.body);

        if (body.detection_time_ms) {
          detectLatency.add(body.detection_time_ms);
        }

        if (body.processed) {
          eventsProcessed.add(body.processed);
        }

        if (body.anomalies) {
          anomaliesDetected.add(body.anomalies.length);
        }
      } catch (e) {
        console.error(`Failed to parse response: ${e.message}`);
      }
    }

    helpers.sleepWithJitter(1, 0.5);
  }
  // 20% - Large batch detection
  else if (scenario < 0.9) {
    const batchSize = Math.floor(Math.random() * 500) + 500; // 500-1000 events
    const events = helpers.generateMixedBatch(batchSize);
    const streamId = `large-batch-${__VU}`;
    const payload = helpers.createDetectPayload(events, streamId, 'strict');

    const res = http.post(
      `${baseUrl}/v1/demo/detect`,
      payload,
      {
        headers: helpers.demoHeaders(),
        tags: { endpoint: 'detect', scenario: 'large_batch' },
      }
    );

    const success = check(res, {
      'large batch status is 200': (r) => r.status === 200,
      'large batch processes all events': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.processed === batchSize;
        } catch {
          return false;
        }
      },
    });

    detectSuccess.add(success);

    if (res.status === 200) {
      try {
        const body = JSON.parse(res.body);
        if (body.detection_time_ms) {
          detectLatency.add(body.detection_time_ms);
        }
        if (body.processed) {
          eventsProcessed.add(body.processed);
        }
        if (body.anomalies) {
          anomaliesDetected.add(body.anomalies.length);
        }
      } catch (e) {
        console.error(`Failed to parse large batch response: ${e.message}`);
      }
    }

    helpers.sleepWithJitter(2, 1);
  }
  // 10% - Health checks
  else {
    const healthRes = http.get(`${baseUrl}/healthz`, {
      tags: { endpoint: 'health', scenario: 'monitoring' },
    });

    check(healthRes, {
      'health check is 200': (r) => r.status === 200,
    });

    sleep(0.5);
  }
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000;
  console.log(`Load test complete! Duration: ${duration.toFixed(2)}s`);
  console.log(`Tested against: ${data.baseUrl}`);
}
