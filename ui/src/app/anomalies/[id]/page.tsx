'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api';
import { Anomaly } from '@/lib/types';
import GlassBoxExplanation from '@/components/GlassBoxExplanation';
import CompressionChart from '@/components/CompressionChart';

export default function AnomalyDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [anomaly, setAnomaly] = useState<Anomaly | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showBaselineData, setShowBaselineData] = useState(false);
  const [showWindowData, setShowWindowData] = useState(false);

  useEffect(() => {
    loadAnomaly();
  }, [id]);

  const loadAnomaly = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiClient.getAnomaly(id);
      setAnomaly(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load anomaly');
      // Use mock data for development
      setAnomaly(getMockAnomaly(id));
    } finally {
      setLoading(false);
    }
  };

  const handleStatusUpdate = async (status: 'acknowledged' | 'dismissed') => {
    if (!anomaly) return;

    try {
      const updated = await apiClient.updateAnomalyStatus(anomaly.id, status);
      setAnomaly(updated);
    } catch (err) {
      console.error('Failed to update status:', err);
      // Update locally for mock
      setAnomaly({ ...anomaly, status });
    }
  };

  const handleExport = async () => {
    if (!anomaly) return;

    try {
      const blob = await apiClient.exportAnomaly(anomaly.id);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `anomaly-${anomaly.id}-evidence.json`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      console.error('Failed to export:', err);
      // Fallback: create JSON export from current data
      const dataStr = JSON.stringify(anomaly, null, 2);
      const blob = new Blob([dataStr], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `anomaly-${anomaly.id}-evidence.json`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    }
  };

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500 dark:text-gray-400">Loading anomaly...</div>
      </div>
    );
  }

  if (error && !anomaly) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500">Error: {error}</div>
        <button
          onClick={() => router.push('/anomalies')}
          className="mt-4 text-indigo-600 hover:text-indigo-800"
        >
          Back to Anomalies
        </button>
      </div>
    );
  }

  if (!anomaly) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500 dark:text-gray-400">Anomaly not found</div>
      </div>
    );
  }

  const getSeverityBadge = (pValue: number) => {
    if (pValue < 0.01)
      return (
        <span className="px-3 py-1 bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200 rounded-full text-sm font-semibold">
          Critical (p &lt; 0.01)
        </span>
      );
    if (pValue < 0.05)
      return (
        <span className="px-3 py-1 bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200 rounded-full text-sm font-semibold">
          High (p &lt; 0.05)
        </span>
      );
    return (
      <span className="px-3 py-1 bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200 rounded-full text-sm font-semibold">
        Moderate
      </span>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <button
            onClick={() => router.push('/anomalies')}
            className="text-indigo-600 hover:text-indigo-800 dark:text-indigo-400 mb-2 text-sm"
          >
            ← Back to Anomalies
          </button>
          <h2 className="text-3xl font-bold text-gray-900 dark:text-white">
            Anomaly Details
          </h2>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            ID: {anomaly.id} | {new Date(anomaly.timestamp).toLocaleString()}
          </p>
        </div>
        <div className="flex space-x-2">
          {anomaly.status === 'pending' && (
            <>
              <button
                onClick={() => handleStatusUpdate('acknowledged')}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Acknowledge
              </button>
              <button
                onClick={() => handleStatusUpdate('dismissed')}
                className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
              >
                Dismiss
              </button>
            </>
          )}
          <button
            onClick={handleExport}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
          >
            Export Evidence
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md p-4">
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            Warning: {error}. Showing mock data for development.
          </p>
        </div>
      )}

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Stream Type
          </div>
          <div className="mt-1 text-2xl font-semibold text-gray-900 dark:text-white capitalize">
            {anomaly.streamType}
          </div>
        </div>
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            NCD Score
          </div>
          <div className="mt-1 text-2xl font-mono font-semibold text-gray-900 dark:text-white">
            {anomaly.ncdScore.toFixed(4)}
          </div>
        </div>
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div className="text-sm text-gray-500 dark:text-gray-400">P-Value</div>
          <div className="mt-1 text-2xl font-mono font-semibold text-gray-900 dark:text-white">
            {anomaly.pValue.toFixed(4)}
          </div>
        </div>
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div className="text-sm text-gray-500 dark:text-gray-400">Severity</div>
          <div className="mt-2">{getSeverityBadge(anomaly.pValue)}</div>
        </div>
      </div>

      {/* Glass-Box Explanation */}
      <GlassBoxExplanation explanation={anomaly.glassBoxExplanation} />

      {/* Compression Analysis */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <CompressionChart compressionRatios={anomaly.compressionRatios} />
      </div>

      {/* NCD Visualization */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Normalized Compression Distance Analysis
        </h3>
        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
              NCD Score: {anomaly.ncdScore.toFixed(4)}
            </p>
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-6 relative">
              <div
                className="h-6 bg-gradient-to-r from-green-500 via-yellow-500 to-red-500 rounded-full"
                style={{ width: `${anomaly.ncdScore * 100}%` }}
              />
              <div className="absolute inset-0 flex items-center justify-center text-xs font-semibold text-white mix-blend-difference">
                {(anomaly.ncdScore * 100).toFixed(1)}%
              </div>
            </div>
            <div className="flex justify-between text-xs text-gray-500 dark:text-gray-500 mt-1">
              <span>0.0 (Identical)</span>
              <span>0.5 (Moderate)</span>
              <span>1.0 (Completely Different)</span>
            </div>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              <strong>Formula:</strong> NCD(x,y) = (C(xy) - min(C(x),C(y))) /
              max(C(x),C(y))
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-500 mt-2">
              Where C(x) is the compressed size. Higher NCD indicates greater
              dissimilarity between baseline and window data.
            </p>
          </div>
        </div>
      </div>

      {/* Statistical Significance */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Statistical Significance
        </h3>
        <div className="space-y-3">
          <div className="flex justify-between items-center">
            <span className="text-gray-700 dark:text-gray-300">P-Value:</span>
            <span className="font-mono text-lg text-gray-900 dark:text-white">
              {anomaly.pValue.toFixed(6)}
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-gray-700 dark:text-gray-300">
              Confidence Level:
            </span>
            <span className="font-semibold text-gray-900 dark:text-white">
              {((1 - anomaly.pValue) * 100).toFixed(2)}%
            </span>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              {anomaly.pValue < 0.01 ? (
                <>
                  <strong>Highly Significant:</strong> This anomaly is statistically
                  significant at the 99% confidence level. The probability of this
                  pattern occurring by chance is less than 1%.
                </>
              ) : anomaly.pValue < 0.05 ? (
                <>
                  <strong>Significant:</strong> This anomaly is statistically
                  significant at the 95% confidence level. The probability of this
                  pattern occurring by chance is less than 5%.
                </>
              ) : (
                <>
                  <strong>Moderate Significance:</strong> This anomaly shows some
                  statistical significance but may require additional investigation.
                </>
              )}
            </p>
          </div>
        </div>
      </div>

      {/* Raw Data */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow space-y-4">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white">
          Raw Data
        </h3>

        {/* Baseline Data */}
        <div>
          <button
            onClick={() => setShowBaselineData(!showBaselineData)}
            className="flex items-center justify-between w-full p-3 bg-gray-50 dark:bg-gray-900 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            <span className="font-medium text-gray-900 dark:text-white">
              Baseline Data
            </span>
            <svg
              className={`w-5 h-5 transform transition-transform ${
                showBaselineData ? 'rotate-180' : ''
              }`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>
          {showBaselineData && (
            <pre className="mt-2 p-4 bg-gray-900 text-green-400 rounded overflow-x-auto text-xs">
              {JSON.stringify(JSON.parse(anomaly.baselineData), null, 2)}
            </pre>
          )}
        </div>

        {/* Window Data */}
        <div>
          <button
            onClick={() => setShowWindowData(!showWindowData)}
            className="flex items-center justify-between w-full p-3 bg-gray-50 dark:bg-gray-900 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            <span className="font-medium text-gray-900 dark:text-white">
              Window Data
            </span>
            <svg
              className={`w-5 h-5 transform transition-transform ${
                showWindowData ? 'rotate-180' : ''
              }`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>
          {showWindowData && (
            <pre className="mt-2 p-4 bg-gray-900 text-green-400 rounded overflow-x-auto text-xs">
              {JSON.stringify(JSON.parse(anomaly.windowData), null, 2)}
            </pre>
          )}
        </div>
      </div>
    </div>
  );
}

// Mock data for development
function getMockAnomaly(id: string): Anomaly {
  return {
    id,
    timestamp: new Date(Date.now() - 3600000).toISOString(),
    streamType: 'logs',
    ncdScore: 0.4523,
    pValue: 0.0087,
    status: 'pending',
    glassBoxExplanation:
      'High NCD score detected: window data differs significantly from baseline. Compression ratio dropped from 0.82 to 0.35, indicating structural change in log patterns. The algorithm detected a 57% decrease in compression efficiency, suggesting that the recent data contains novel patterns not present in the historical baseline. This dramatic shift in compressibility is a strong signal of anomalous behavior, as normal operations typically maintain consistent compression characteristics.',
    compressionRatios: {
      baseline: 0.82,
      window: 0.35,
      combined: 0.51,
    },
    baselineData: JSON.stringify({
      type: 'baseline',
      samples: 100,
      timestamp_start: new Date(Date.now() - 86400000).toISOString(),
      timestamp_end: new Date(Date.now() - 7200000).toISOString(),
      pattern_summary: 'Regular log patterns observed',
    }),
    windowData: JSON.stringify({
      type: 'window',
      samples: 50,
      timestamp_start: new Date(Date.now() - 7200000).toISOString(),
      timestamp_end: new Date(Date.now() - 3600000).toISOString(),
      pattern_summary: 'Anomalous patterns detected - unusual error codes',
    }),
  };
}
