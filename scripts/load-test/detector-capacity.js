/**
 * Detector Capacity Test - LRU eviction and multi-stream handling
 *
 * Purpose: Test detector capacity and LRU eviction with 1000+ unique streams
 * Duration: 10 minutes
 * VUs: 100
 *
 * Usage: k6 run scripts/load-test/detector-capacity.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend, Rate } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const uniqueStreamsCreated = new Counter('unique_streams_created');
const detectLatency = new Trend('detect_latency', true);
const eventsProcessed = new Counter('events_processed');
const detectSuccess = new Rate('detect_success');
const lruEvictions = new Counter('lru_evictions_estimated');

// Test configuration
export const options = {
  vus: 100,
  duration: '10m',
  thresholds: {
    'http_req_duration': ['p(95)<1000', 'p(99)<2000'],
    'http_req_failed': ['rate<0.05'],
    'detect_success': ['rate>0.95'],
    'unique_streams_created': ['value>1000'], // Must create 1000+ streams
  },
  tags: {
    test_type: 'capacity',
  },
};

const BASE_URL = helpers.getBaseUrl();

// Track unique stream IDs per VU
const streamRegistry = new Set();

export function setup() {
  console.log(`Starting capacity test against ${BASE_URL}`);
  console.log('Testing LRU eviction with 1000+ unique streams...');

  // Verify server is ready
  const health = http.get(`${BASE_URL}/healthz`);
  if (health.status !== 200) {
    throw new Error(`Server not healthy: ${health.status}`);
  }

  return {
    baseUrl: BASE_URL,
    startTime: Date.now(),
    targetStreams: 1500, // Target more than minimum to ensure we hit threshold
  };
}

export default function (data) {
  const baseUrl = data.baseUrl;

  // Create unique stream ID for each iteration
  // This ensures we create many unique detectors
  const streamId = helpers.generateUniqueStreamId(__VU, __ITER);
  streamRegistry.add(streamId);
  uniqueStreamsCreated.add(1);

  // Generate events with varying characteristics
  const batchSize = Math.floor(Math.random() * 100) + 50; // 50-150 events
  const events = helpers.generateMixedBatch(batchSize, streamId);

  // Vary detection profiles to test different configurations
  const profiles = ['sensitive', 'balanced', 'strict'];
  const profile = profiles[Math.floor(Math.random() * profiles.length)];

  const payload = helpers.createDetectPayload(events, streamId, profile);

  const res = http.post(
    `${baseUrl}/v1/demo/detect`,
    payload,
    {
      headers: helpers.demoHeaders(),
      tags: {
        endpoint: 'detect',
        profile: profile,
        test_phase: 'capacity',
      },
    }
  );

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response has stream_id': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.stream_id === streamId;
      } catch {
        return false;
      }
    },
    'detector handles new stream': (r) => {
      // Should not fail even with many streams
      return r.status !== 500;
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

      // If we're creating many streams, estimate evictions
      // (assuming LRU cache with ~1000 capacity)
      const estimatedEvictions = Math.max(0, streamRegistry.size - 1000);
      if (estimatedEvictions > 0 && __ITER % 100 === 0) {
        console.log(`Estimated LRU evictions: ${estimatedEvictions} (${streamRegistry.size} total streams)`);
      }
    } catch (e) {
      console.error(`Failed to parse response: ${e.message}`);
    }
  } else {
    console.warn(`Failed to create stream ${streamId}: ${res.status}`);
  }

  // Shorter sleep to create streams faster
  helpers.sleepWithJitter(0.5, 0.2);

  // Periodically report progress
  if (__ITER % 100 === 0) {
    console.log(`VU ${__VU}: Created ${streamRegistry.size} unique streams`);
  }
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000 / 60;
  const totalStreams = streamRegistry.size;

  console.log(`\nCapacity test complete! Duration: ${duration.toFixed(2)} minutes`);
  console.log(`Tested against: ${data.baseUrl}`);
  console.log(`Total unique streams created: ${totalStreams}`);

  if (totalStreams >= 1000) {
    console.log('PASS: Successfully created 1000+ unique streams');
    const evictions = Math.max(0, totalStreams - 1000);
    console.log(`Estimated LRU evictions: ${evictions}`);
  } else {
    console.warn(`WARNING: Only created ${totalStreams} streams (target: 1000+)`);
  }
}

export function handleSummary(data) {
  const totalStreams = data.metrics.unique_streams_created.values.count;

  console.log('\n=== Capacity Test Summary ===');
  console.log(`Unique streams created: ${totalStreams}`);
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Success rate: ${(1 - data.metrics.http_req_failed.values.rate) * 100}%`);
  console.log(`P95 latency: ${data.metrics.http_req_duration.values['p(95)']}ms`);
  console.log(`P99 latency: ${data.metrics.http_req_duration.values['p(99)']}ms`);

  if (data.metrics.detect_latency) {
    console.log(`Detect P95 latency: ${data.metrics.detect_latency.values['p(95)']}ms`);
    console.log(`Detect P99 latency: ${data.metrics.detect_latency.values['p(99)']}ms`);
  }

  return {
    'capacity_test_results.json': JSON.stringify(data, null, 2),
  };
}
