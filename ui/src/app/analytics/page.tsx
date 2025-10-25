'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import { PerformanceMetrics, StatisticalSummary } from '@/lib/types';

export default function AnalyticsPage() {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [summary, setSummary] = useState<StatisticalSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAnalytics();
  }, []);

  const loadAnalytics = async () => {
    try {
      setLoading(true);
      setError(null);
      const [metricsData, summaryData] = await Promise.all([
        apiClient.getPerformanceMetrics(),
        apiClient.getStatisticalSummary(),
      ]);
      setMetrics(metricsData);
      setSummary(summaryData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load analytics');
      // Use mock data for development
      setMetrics(getMockMetrics());
      setSummary(getMockSummary());
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500 dark:text-gray-400">
          Loading analytics...
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
            Analytics Dashboard
          </h2>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Statistical analysis and performance metrics
          </p>
        </div>
        <button
          onClick={loadAnalytics}
          className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
        >
          Refresh
        </button>
      </div>

      {error && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md p-4">
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            Warning: {error}. Showing mock data for development.
          </p>
        </div>
      )}

      {/* Performance Metrics */}
      {metrics && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <svg
                  className="h-8 w-8 text-indigo-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
              </div>
              <div className="ml-4">
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  Events/Second
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-white">
                  {metrics.eventsPerSecond.toFixed(0)}
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <svg
                  className="h-8 w-8 text-green-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <div className="ml-4">
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  Avg Latency
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-white">
                  {metrics.averageLatencyMs.toFixed(0)}ms
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <svg
                  className="h-8 w-8 text-red-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
              </div>
              <div className="ml-4">
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  Total Anomalies
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-white">
                  {metrics.totalAnomalies}
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <svg
                  className="h-8 w-8 text-yellow-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
              </div>
              <div className="ml-4">
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  Detection Rate
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-white">
                  {(metrics.detectionRate * 100).toFixed(2)}%
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Compression Efficiency */}
      {summary && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
            Compression Efficiency
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="text-center p-4 bg-gray-50 dark:bg-gray-900 rounded">
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Average Ratio
              </div>
              <div className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
                {summary.compressionEfficiency.averageRatio.toFixed(3)}
              </div>
            </div>
            <div className="text-center p-4 bg-gray-50 dark:bg-gray-900 rounded">
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Min Ratio
              </div>
              <div className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
                {summary.compressionEfficiency.minRatio.toFixed(3)}
              </div>
            </div>
            <div className="text-center p-4 bg-gray-50 dark:bg-gray-900 rounded">
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Max Ratio
              </div>
              <div className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
                {summary.compressionEfficiency.maxRatio.toFixed(3)}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Detection Rate Over Time */}
      {summary && summary.detectionRateOverTime.length > 0 && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
            Detection Rate Over Time
          </h3>
          <div className="space-y-2">
            {summary.detectionRateOverTime.slice(-10).map((point, idx) => (
              <div key={idx} className="flex items-center">
                <div className="w-32 text-sm text-gray-600 dark:text-gray-400">
                  {new Date(point.timestamp).toLocaleTimeString()}
                </div>
                <div className="flex-1 mx-4">
                  <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-6">
                    <div
                      className="bg-indigo-500 h-6 rounded-full flex items-center justify-end pr-2"
                      style={{ width: `${point.rate * 100}%` }}
                    >
                      <span className="text-xs text-white font-semibold">
                        {(point.rate * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* False Positive Tracking */}
      {summary &&
        summary.falsePositiveTracking &&
        summary.falsePositiveTracking.length > 0 && (
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              False Positive Rate
            </h3>
            <div className="space-y-2">
              {summary.falsePositiveTracking.slice(-10).map((point, idx) => (
                <div key={idx} className="flex items-center">
                  <div className="w-32 text-sm text-gray-600 dark:text-gray-400">
                    {new Date(point.timestamp).toLocaleTimeString()}
                  </div>
                  <div className="flex-1 mx-4">
                    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-6">
                      <div
                        className={`h-6 rounded-full flex items-center justify-end pr-2 ${
                          point.rate < 0.05
                            ? 'bg-green-500'
                            : point.rate < 0.1
                            ? 'bg-yellow-500'
                            : 'bg-red-500'
                        }`}
                        style={{ width: `${Math.min(point.rate * 100, 100)}%` }}
                      >
                        <span className="text-xs text-white font-semibold">
                          {(point.rate * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <div className="mt-4 text-sm text-gray-600 dark:text-gray-400">
              <p>
                Lower false positive rates indicate more accurate anomaly detection.
                Target: &lt;5%
              </p>
            </div>
          </div>
        )}

      {/* System Information */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          System Information
        </h3>
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <dt className="text-sm text-gray-500 dark:text-gray-400">
              Detection Algorithm
            </dt>
            <dd className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
              Compression-Based Anomaly Detection (CBAD)
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-500 dark:text-gray-400">
              Compression Method
            </dt>
            <dd className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
              OpenZL with NCD
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-500 dark:text-gray-400">
              Statistical Test
            </dt>
            <dd className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
              P-value significance testing
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-500 dark:text-gray-400">
              Explanation Method
            </dt>
            <dd className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
              Glass-box mathematical reasoning
            </dd>
          </div>
        </dl>
      </div>
    </div>
  );
}

function getMockMetrics(): PerformanceMetrics {
  return {
    eventsPerSecond: 1247,
    averageLatencyMs: 234,
    totalAnomalies: 127,
    detectionRate: 0.034,
    falsePositiveRate: 0.028,
  };
}

function getMockSummary(): StatisticalSummary {
  const now = Date.now();
  return {
    detectionRateOverTime: Array.from({ length: 24 }, (_, i) => ({
      timestamp: new Date(now - (23 - i) * 3600000).toISOString(),
      rate: 0.02 + Math.random() * 0.05,
    })),
    falsePositiveTracking: Array.from({ length: 24 }, (_, i) => ({
      timestamp: new Date(now - (23 - i) * 3600000).toISOString(),
      rate: 0.01 + Math.random() * 0.04,
    })),
    compressionEfficiency: {
      averageRatio: 0.742,
      minRatio: 0.287,
      maxRatio: 0.923,
    },
    performanceStats: getMockMetrics(),
  };
}
