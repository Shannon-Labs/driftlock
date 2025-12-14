/**
 * Rate Limit Validation Test
 *
 * Purpose: Verify rate limiting is working correctly on demo endpoint
 * Expected: 10 requests per minute per IP, then 429 responses
 * Duration: 3 minutes
 * VUs: 1 (to test from single IP)
 *
 * Usage: k6 run scripts/load-test/rate-limit-validation.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate } from 'k6/metrics';
import * as helpers from './helpers.js';

// Custom metrics
const rateLimitedResponses = new Counter('rate_limited_responses');
const successfulRequests = new Counter('successful_requests');
const rateLimitHitCorrectly = new Rate('rate_limit_hit_correctly');

// Test configuration
export const options = {
  vus: 1,
  duration: '3m',
  thresholds: {
    'http_req_duration': ['p(95)<200'],
    'rate_limited_responses': ['rate>0.5'], // Expect many 429s
  },
  tags: {
    test_type: 'rate_limit',
  },
};

const BASE_URL = helpers.getBaseUrl();
const REQUESTS_PER_MINUTE = 10;

export function setup() {
  console.log(`Starting rate limit validation test against ${BASE_URL}`);
  console.log(`Expected limit: ${REQUESTS_PER_MINUTE} requests per minute per IP`);
  console.log('Testing with single VU to simulate single IP...');

  return {
    baseUrl: BASE_URL,
    startTime: Date.now(),
    requestCount: 0,
    rateLimitCount: 0,
  };
}

export default function (data) {
  const baseUrl = data.baseUrl;
  data.requestCount++;

  // Generate small batch
  const events = helpers.generateEventBatch(10, 0.1);
  const payload = helpers.createDetectPayload(events);

  const res = http.post(
    `${baseUrl}/v1/demo/detect`,
    payload,
    {
      headers: helpers.demoHeaders(),
      tags: { endpoint: 'demo_detect' },
    }
  );

  const isRateLimited = helpers.isRateLimited(res);

  if (isRateLimited) {
    data.rateLimitCount++;
    rateLimitedResponses.add(1);

    check(res, {
      'rate limit returns 429': (r) => r.status === 429,
      'rate limit has retry-after header': (r) => r.headers['Retry-After'] !== undefined,
    });

    console.log(`Request ${data.requestCount}: RATE LIMITED (429) - Total rate limits: ${data.rateLimitCount}`);
  } else if (res.status === 200) {
    successfulRequests.add(1);
    console.log(`Request ${data.requestCount}: SUCCESS (200)`);
  } else {
    console.warn(`Request ${data.requestCount}: Unexpected status ${res.status}`);
  }

  // Check rate limit timing
  const elapsedMinutes = (Date.now() - data.startTime) / 1000 / 60;
  const expectedSuccess = Math.floor(elapsedMinutes * REQUESTS_PER_MINUTE);
  const actualSuccess = data.requestCount - data.rateLimitCount;

  const rateLimitWorking = check(res, {
    'rate limit is enforced correctly': () => {
      // Allow some tolerance
      return actualSuccess <= expectedSuccess + 2;
    },
  });

  rateLimitHitCorrectly.add(rateLimitWorking);

  // Send requests rapidly to trigger rate limit
  sleep(2); // 30 requests per minute = 1 every 2 seconds
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000 / 60;
  const totalRequests = data.requestCount;
  const rateLimited = data.rateLimitCount;
  const successful = totalRequests - rateLimited;
  const expectedSuccess = Math.floor(duration * REQUESTS_PER_MINUTE);

  console.log('\n=== Rate Limit Validation Summary ===');
  console.log(`Test duration: ${duration.toFixed(2)} minutes`);
  console.log(`Total requests sent: ${totalRequests}`);
  console.log(`Successful requests: ${successful}`);
  console.log(`Rate limited (429): ${rateLimited}`);
  console.log(`Expected max success: ~${expectedSuccess}`);

  if (successful <= expectedSuccess + 2 && rateLimited > 0) {
    console.log('\nPASS: Rate limiting is working correctly!');
  } else if (rateLimited === 0) {
    console.warn('\nWARNING: No rate limiting detected! Rate limiter may not be enabled.');
  } else {
    console.warn('\nWARNING: Rate limit behavior is unexpected.');
    console.warn(`Expected ~${expectedSuccess} successful requests, got ${successful}`);
  }
}

export function handleSummary(data) {
  return {
    'rate_limit_test_results.json': JSON.stringify(data, null, 2),
  };
}
