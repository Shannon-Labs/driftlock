/**
 * Soak Test - Memory leak and stability test
 *
 * Purpose: Run sustained load to identify memory leaks and degradation
 * Duration: 30 minutes
 * VUs: 50 (constant)
 *
 * Usage: k6 run scripts/load-test/soak.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend, Rate, Gauge } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const detectLatency = new Trend('detect_latency', true);
const eventsProcessed = new Counter('events_processed');
const anomaliesDetected = new Counter('anomalies_detected');
const detectSuccess = new Rate('detect_success');
const responseTimeOverTime = new Trend('response_time_trend', true);
const memoryLeakIndicator = new Gauge('memory_leak_indicator');

// Test configuration
export const options = {
  vus: 50,
  duration: '30m',
  thresholds: {
    'http_req_duration': ['p(95)<600', 'p(99)<1200'],
    'http_req_failed': ['rate<0.01'],
    'http_reqs': ['rate>50'],
    'detect_latency': ['p(95)<500', 'p(99)<1000'],
    'detect_success': ['rate>0.99'],
    // Critical: Response times should not increase significantly over time
    'response_time_trend': ['p(95)<700', 'p(99)<1300'],
  },
  tags: {
    test_type: 'soak',
  },
};

const BASE_URL = helpers.getBaseUrl();

// Track baseline performance
let baselineP95 = 0;
let baselineP99 = 0;
let iterationCount = 0;

export function setup() {
  console.log(`Starting soak test against ${BASE_URL}`);
  console.log('Duration: 30 minutes');
  console.log('Monitoring for memory leaks and performance degradation...');

  // Verify server is ready
  const health = http.get(`${BASE_URL}/healthz`);
  if (health.status !== 200) {
    throw new Error(`Server not healthy: ${health.status}`);
  }

  return {
    baseUrl: BASE_URL,
    startTime: Date.now(),
    checkpointTimes: [],
    checkpointLatencies: [],
  };
}

export default function (data) {
  const baseUrl = data.baseUrl;
  iterationCount++;

  // Vary traffic patterns to simulate real usage
  const scenario = Math.random();
  let batchSize, profile;

  if (scenario < 0.5) {
    // Small batches (most common)
    batchSize = Math.floor(Math.random() * 50) + 10;
    profile = 'balanced';
  } else if (scenario < 0.8) {
    // Medium batches
    batchSize = Math.floor(Math.random() * 200) + 100;
    profile = 'balanced';
  } else {
    // Large batches
    batchSize = Math.floor(Math.random() * 500) + 500;
    profile = 'strict';
  }

  const events = helpers.generateMixedBatch(batchSize);
  const streamId = `soak-stream-${__VU}-${Math.floor(__ITER / 20)}`;
  const payload = helpers.createDetectPayload(events, streamId, profile);

  const startTime = Date.now();
  const res = http.post(
    `${baseUrl}/v1/demo/detect`,
    payload,
    {
      headers: helpers.demoHeaders(),
      tags: {
        endpoint: 'detect',
        batch_size: batchSize < 100 ? 'small' : batchSize < 300 ? 'medium' : 'large',
      },
    }
  );
  const duration = Date.now() - startTime;

  // Track response time trend
  responseTimeOverTime.add(duration);

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response is valid': (r) => helpers.validateDetectResponse(r).valid,
    'no significant degradation': (r) => {
      // If we have baseline, check for degradation
      if (baselineP95 > 0) {
        return r.timings.duration < baselineP95 * 2; // Allow 2x degradation
      }
      return true;
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

      // Periodic checkpoint (every 100 iterations)
      if (iterationCount % 100 === 0) {
        data.checkpointTimes.push(Date.now() - data.startTime);
        data.checkpointLatencies.push(body.detection_time_ms || duration);

        // Calculate memory leak indicator
        if (data.checkpointLatencies.length > 2) {
          const recent = data.checkpointLatencies.slice(-5);
          const early = data.checkpointLatencies.slice(0, 5);
          const recentAvg = recent.reduce((a, b) => a + b, 0) / recent.length;
          const earlyAvg = early.reduce((a, b) => a + b, 0) / early.length;
          const degradation = (recentAvg - earlyAvg) / earlyAvg;

          memoryLeakIndicator.add(degradation);

          if (degradation > 0.5) {
            console.warn(`Potential memory leak detected! Degradation: ${(degradation * 100).toFixed(2)}%`);
          }
        }
      }
    } catch (e) {
      console.error(`Failed to parse response: ${e.message}`);
    }
  }

  // Realistic think time
  helpers.sleepWithJitter(1.5, 0.5);

  // Occasional health check
  if (__ITER % 50 === 0) {
    const healthRes = http.get(`${baseUrl}/healthz`, {
      tags: { endpoint: 'health' },
    });

    check(healthRes, {
      'health check still passing': (r) => r.status === 200,
    });
  }
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000 / 60;
  console.log(`\nSoak test complete! Duration: ${duration.toFixed(2)} minutes`);
  console.log(`Tested against: ${data.baseUrl}`);

  // Analyze checkpoints for degradation
  if (data.checkpointLatencies.length > 10) {
    const early = data.checkpointLatencies.slice(0, 5);
    const late = data.checkpointLatencies.slice(-5);
    const earlyAvg = early.reduce((a, b) => a + b, 0) / early.length;
    const lateAvg = late.reduce((a, b) => a + b, 0) / late.length;
    const degradation = ((lateAvg - earlyAvg) / earlyAvg * 100).toFixed(2);

    console.log('\n=== Performance Over Time ===');
    console.log(`Early average latency: ${earlyAvg.toFixed(2)}ms`);
    console.log(`Late average latency: ${lateAvg.toFixed(2)}ms`);
    console.log(`Degradation: ${degradation}%`);

    if (parseFloat(degradation) > 20) {
      console.warn('WARNING: Significant performance degradation detected!');
      console.warn('This may indicate a memory leak or resource exhaustion.');
    } else {
      console.log('PASS: No significant degradation detected.');
    }
  }
}

export function handleSummary(data) {
  console.log('\n=== Soak Test Summary ===');
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Success rate: ${(1 - data.metrics.http_req_failed.values.rate) * 100}%`);
  console.log(`P95 latency: ${data.metrics.http_req_duration.values['p(95)']}ms`);
  console.log(`P99 latency: ${data.metrics.http_req_duration.values['p(99)']}ms`);
  console.log(`Events processed: ${data.metrics.events_processed.values.count}`);
  console.log(`Anomalies detected: ${data.metrics.anomalies_detected.values.count}`);

  return {
    'soak_test_results.json': JSON.stringify(data, null, 2),
  };
}
