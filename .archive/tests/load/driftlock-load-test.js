import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// A custom metric to track error rate
const errorRate = new Rate('errors');

export let options = {
  // Key performance indicators (KPIs) for Driftlock
  thresholds: {
    // HTTP response time
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% of requests under 500ms, 99% under 1000ms
    
    // Success rate
    'http_req_failed': ['rate<0.01'], // Less than 1% of requests should fail
    
    // Error rate
    'errors': ['rate<0.01'], // Less than 1% of requests should fail
  },
  
  // Execution scenarios
  scenarios: {
    // Sustained load test: 1000 requests per second for 10 minutes
    sustained_load: {
      executor: 'constant-arrival-rate',
      rate: 1000, // requests per second
      timeUnit: '1s',
      duration: '10m',
      preAllocatedVUs: 50, // Pre-allocated virtual users
      maxVUs: 200, // Maximum virtual users
    },
    
    // Spike test: gradual increase to 5000 req/s
    spike_test: {
      executor: 'ramping-arrival-rate',
      startRate: 100,
      timeUnit: '1s',
      stages: [
        { target: 1000, duration: '2m' },   // Ramp up to 1000 req/s over 2 minutes
        { target: 1000, duration: '3m' },   // Stay at 1000 req/s for 3 minutes
        { target: 5000, duration: '1m' },   // Spike to 5000 req/s over 1 minute
        { target: 5000, duration: '2m' },   // Sustain 5000 req/s for 2 minutes
        { target: 100, duration: '2m' },    // Ramp down to 100 req/s over 2 minutes
      ],
      preAllocatedVUs: 100,
      maxVUs: 500,
    },
    
    // Stress test: gradually increase until failure
    stress_test: {
      executor: 'ramping-arrival-rate',
      startRate: 1000,
      timeUnit: '1s',
      stages: [
        { target: 2000, duration: '2m' },   // Increase to 2000 req/s
        { target: 4000, duration: '2m' },   // Increase to 4000 req/s
        { target: 8000, duration: '2m' },   // Increase to 8000 req/s
        { target: 10000, duration: '2m' },  // Increase to 10000 req/s
      ],
      preAllocatedVUs: 200,
      maxVUs: 1000,
    },
  },
};

// Environment variables for test configuration
const BASE_URL = __ENV.DRIFTLOCK_BASE_URL || 'http://localhost:8080';
const API_KEY = __ENV.DRIFTLOCK_API_KEY || 'test-api-key';

// Headers for API requests
const BASE_HEADERS = {
  'Content-Type': 'application/json',
  'X-API-Key': API_KEY,
};

export default function() {
  // Test the anomalies endpoint - this is where most load will occur
  let response = http.get(`${BASE_URL}/v1/anomalies`, {
    headers: BASE_HEADERS,
  });

  // Check response
  let success = check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
    'has anomalies in response': (r) => {
      try {
        let data = JSON.parse(r.body);
        return Array.isArray(data.anomalies);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!success);

  // Add slight delay to simulate more realistic usage
  sleep(0.1);
}

// Setup function - run once before the test
export function setup() {
  console.log('Starting Driftlock load test...');
  console.log(`Target URL: ${BASE_URL}`);
  
  // Verify API connectivity
  let healthCheck = http.get(`${BASE_URL}/healthz`);
  if (healthCheck.status !== 200) {
    console.error(`Health check failed: ${healthCheck.status}`);
    return null;
  }
  
  console.log('API health check passed');
  return { setupCompleted: true };
}

// Teardown function - run once after the test
export function teardown(data) {
  console.log('Load test completed');
  console.log('Generating performance report...');
}