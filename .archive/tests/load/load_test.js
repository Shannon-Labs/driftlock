import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Custom metrics
const apiResponseTime = new Trend('api_response_time');

export let options = {
  // Key configurations for performance testing
  stages: [
    // Ramp up to 100 users over 2 minutes
    { duration: '2m', target: 100 },
    // Stay at 100 users for 5 minutes
    { duration: '5m', target: 100 },
    // Spike to 1000 users over 2 minutes
    { duration: '2m', target: 1000 },
    // Stay at 1000 users for 5 minutes
    { duration: '5m', target: 1000 },
    // Ramp down to 0 users over 2 minutes
    { duration: '2m', target: 0 },
  ],
  
  // Thresholds for performance requirements
  thresholds: {
    // HTTP response time thresholds (ms)
    'api_response_time': ['p(95)<500', 'p(99)<1000'], // 95% of requests < 500ms, 99% < 1s
    'http_req_duration': ['p(95)<500', 'p(99)<1000'],
    
    // Error rate thresholds
    'http_req_failed': ['rate<0.01'], // Less than 1% of requests should fail
    
    // Throughput - aim for at least 500 requests per second
    'http_reqs{method:GET}': ['rate>250'],
    'http_reqs{method:POST}': ['rate>50'],
  },
  
  // Configuration
  noVUConnectionReuse: false,
  userAgent: 'Driftlock-Load-Test/1.0',
};

// Helper function to generate random anomaly data
function generateAnomalyData() {
  const streamTypes = ['logs', 'metrics', 'traces', 'llm'];
  const statuses = ['pending', 'acknowledged', 'dismissed', 'investigating'];
  
  return {
    timestamp: new Date().toISOString(),
    stream_type: streamTypes[Math.floor(Math.random() * streamTypes.length)],
    ncd_score: Math.random() * 0.5, // 0.0 to 0.5
    p_value: Math.random() * 0.1, // 0 to 0.1
    glass_box_explanation: `Compression ratio increased from ${Math.random() * 0.5} to ${Math.random() * 0.5 + 0.3}`,
    compression_baseline: Math.random() * 0.7,
    compression_window: Math.random() * 0.7,
    compression_combined: Math.random() * 0.7,
    confidence_level: 0.9 + Math.random() * 0.09, // 0.9 to 0.99
    baseline_data: { sample: 'data', value: Math.random() * 100 },
    window_data: { sample: 'data', value: Math.random() * 100 },
    metadata: { source: 'load-test', version: '1.0' },
    tags: ['load-test', 'performance', 'k6']
  };
}

export default function() {
  const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
  
  // Test GET /v1/anomalies endpoint (list anomalies)
  let response = http.get(`${BASE_URL}/v1/anomalies`);
  check(response, {
    'list anomalies status is 200': (r) => r.status === 200,
    'list anomalies response time < 500ms': (r) => r.timings.duration < 500,
  });
  apiResponseTime.add(response.timings.duration);
  
  if (response.status === 200) {
    // Extract total count from response if available
    const responseJson = response.json();
    console.log(`Retrieved ${responseJson.anomalies?.length || 0} anomalies, total: ${responseJson.total || 0}`);
  }

  // Test GET /v1/config endpoint
  response = http.get(`${BASE_URL}/v1/config`);
  check(response, {
    'get config status is 200': (r) => r.status === 200,
    'get config response time < 200ms': (r) => r.timings.duration < 200,
  });
  apiResponseTime.add(response.timings.duration);

  // Test GET /v1/analytics/summary endpoint
  response = http.get(`${BASE_URL}/v1/analytics/summary`);
  check(response, {
    'analytics summary status is 200': (r) => r.status === 200,
    'analytics summary response time < 1000ms': (r) => r.timings.duration < 1000,
  });
  apiResponseTime.add(response.timings.duration);

  // Test POST request - simulate anomaly creation (if endpoint exists)
  // Note: The actual API might not have a direct anomaly creation endpoint
  // but we'll test with a potential endpoint for anomaly reporting
  const anomalyData = generateAnomalyData();
  response = http.post(
    `${BASE_URL}/v1/anomalies`,
    JSON.stringify(anomalyData),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  
  // Check if the endpoint exists (could return 404 or 405 if not implemented)
  if (response.status !== 404 && response.status !== 405) {
    check(response, {
      'create anomaly response time < 1000ms': (r) => r.timings.duration < 1000,
    });
    apiResponseTime.add(response.timings.duration);
  }

  // Add a small random sleep to simulate real user behavior
  sleep(Math.random() * 0.5 + 0.1); // Sleep between 0.1 and 0.6 seconds
}

// Setup function - run once before all VUs start
export function setup() {
  console.log('Starting Driftlock Performance Test');
  
  // Verify API is accessible
  const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
  const response = http.get(`${BASE_URL}/healthz`);
  
  if (response.status !== 200) {
    console.error('API is not accessible, health check failed');
    return { error: 'API not accessible' };
  }
  
  console.log('API health check passed');
  return { apiAccessible: true };
}

// Teardown function - run once after all VUs finish
export function teardown(data) {
  console.log('Performance test completed');
  console.log('Final results:');
  console.log(`- API was accessible: ${data.apiAccessible}`);
}