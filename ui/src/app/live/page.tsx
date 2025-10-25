'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { apiClient } from '@/lib/api';
import { Anomaly } from '@/lib/types';
import Link from 'next/link';

export default function LiveFeedPage() {
  const [anomalies, setAnomalies] = useState<Anomaly[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [todayCount, setTodayCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);

  const showToast = useCallback((anomaly: Anomaly) => {
    // In a real implementation, this would use a toast library
    console.log('New anomaly detected:', anomaly);
  }, []);

  const simulateLiveEvents = useCallback(() => {
    // Simulate new anomalies every 10 seconds for development
    const interval = setInterval(() => {
      const mockAnomaly = generateMockAnomaly();
      setAnomalies((prev) => [mockAnomaly, ...prev.slice(0, 49)]);
      setTodayCount((prev) => prev + 1);
      setLastUpdate(new Date());
      showToast(mockAnomaly);
    }, 10000);

    return () => clearInterval(interval);
  }, [showToast]);

  const connectToStream = useCallback(() => {
    try {
      const eventSource = apiClient.createAnomalyStream();
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        setIsConnected(true);
        setError(null);
      };

      eventSource.onmessage = (event) => {
        try {
          const anomaly = JSON.parse(event.data) as Anomaly;
          setAnomalies((prev) => [anomaly, ...prev.slice(0, 49)]); // Keep last 50
          setTodayCount((prev) => prev + 1);
          setLastUpdate(new Date());
          showToast(anomaly);
        } catch (err) {
          console.error('Failed to parse anomaly:', err);
        }
      };

      eventSource.onerror = () => {
        setIsConnected(false);
        setError('Connection lost. Retrying...');
        eventSource.close();
        // Retry after 5 seconds
        setTimeout(connectToStream, 5000);
        // For development, simulate events
        simulateLiveEvents();
      };
    } catch (err) {
      setError('Failed to connect to live stream');
      // For development, simulate events
      simulateLiveEvents();
    }
  }, [simulateLiveEvents, showToast]);

  const loadTodayCount = useCallback(async () => {
    try {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const result = await apiClient.getAnomalies(
        { startDate: today },
        1,
        1
      );
      setTodayCount(result.total);
    } catch (err) {
      console.error('Failed to load count:', err);
      setTodayCount(0);
    }
  }, []);

  useEffect(() => {
    // Connect to SSE stream
    connectToStream();

    // Load today's count
    loadTodayCount();

    // Cleanup on unmount
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, [connectToStream, loadTodayCount]);

  const formatTime = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  const getSeverityColor = (pValue: number) => {
    if (pValue < 0.01) return 'border-l-4 border-red-500';
    if (pValue < 0.05) return 'border-l-4 border-orange-500';
    return 'border-l-4 border-yellow-500';
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
          Live Anomaly Feed
        </h2>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Real-time stream of detected anomalies
        </p>
      </div>

      {/* Status Banner */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="flex items-center">
              <div
                className={`w-3 h-3 rounded-full mr-2 ${
                  isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
                }`}
              />
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {isConnected ? 'Connected' : 'Disconnected'}
              </span>
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              {lastUpdate && (
                <>Last update: {lastUpdate.toLocaleTimeString()}</>
              )}
            </div>
          </div>
          <div className="flex items-center space-x-6">
            <div>
              <div className="text-2xl font-bold text-gray-900 dark:text-white">
                {todayCount}
              </div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                Today&apos;s Anomalies
              </div>
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900 dark:text-white">
                {anomalies.length}
              </div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                In Feed
              </div>
            </div>
          </div>
        </div>
      </div>

      {error && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md p-4">
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            {error}
          </p>
        </div>
      )}

      {/* Live Feed */}
      <div className="space-y-3">
        {anomalies.length === 0 && (
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-12 text-center">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
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
            <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
              Waiting for anomalies...
            </h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              The feed will update in real-time as anomalies are detected.
            </p>
          </div>
        )}

        {anomalies.map((anomaly) => (
          <div
            key={`${anomaly.id}-${anomaly.timestamp}`}
            className={`bg-white dark:bg-gray-800 rounded-lg shadow p-4 ${getSeverityColor(
              anomaly.pValue
            )} hover:shadow-lg transition-shadow`}
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-3">
                  <span className="text-sm font-mono text-gray-500 dark:text-gray-400">
                    {formatTime(anomaly.timestamp)}
                  </span>
                  <span
                    className={`px-2 py-1 text-xs font-semibold rounded ${
                      anomaly.streamType === 'logs'
                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                        : anomaly.streamType === 'metrics'
                        ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                        : 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200'
                    }`}
                  >
                    {anomaly.streamType}
                  </span>
                  <span
                    className={`px-2 py-1 text-xs font-semibold rounded ${
                      anomaly.pValue < 0.01
                        ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                        : anomaly.pValue < 0.05
                        ? 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
                        : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                    }`}
                  >
                    p={anomaly.pValue.toFixed(4)}
                  </span>
                </div>
                <p className="mt-2 text-sm text-gray-900 dark:text-white line-clamp-2">
                  {anomaly.glassBoxExplanation}
                </p>
                <div className="mt-3 flex items-center space-x-4 text-xs text-gray-500 dark:text-gray-400">
                  <span>NCD: {anomaly.ncdScore.toFixed(4)}</span>
                  <span>
                    Compression: {anomaly.compressionRatios.baseline.toFixed(3)} →{' '}
                    {anomaly.compressionRatios.window.toFixed(3)}
                  </span>
                </div>
              </div>
              <Link
                href={`/anomalies/${anomaly.id}`}
                className="ml-4 px-4 py-2 bg-indigo-600 text-white text-sm rounded hover:bg-indigo-700"
              >
                Investigate
              </Link>
            </div>
          </div>
        ))}
      </div>

      {/* Auto-refresh info */}
      <div className="text-center text-sm text-gray-500 dark:text-gray-400">
        <p>This feed updates in real-time via Server-Sent Events</p>
        <p className="mt-1">
          Showing the last {anomalies.length} anomalies
        </p>
      </div>
    </div>
  );
}

// Generate mock anomaly for development
function generateMockAnomaly(): Anomaly {
  const streamTypes: Array<'logs' | 'metrics' | 'traces'> = [
    'logs',
    'metrics',
    'traces',
  ];
  const streamType = streamTypes[Math.floor(Math.random() * streamTypes.length)];
  const ncdScore = 0.2 + Math.random() * 0.4;
  const pValue = 0.001 + Math.random() * 0.09;

  return {
    id: `mock-${Date.now()}-${Math.random()}`,
    timestamp: new Date().toISOString(),
    streamType,
    ncdScore,
    pValue,
    status: 'pending',
    glassBoxExplanation: `New anomaly detected in ${streamType} stream. Compression ratio degradation observed with NCD score ${ncdScore.toFixed(
      4
    )}. Statistical significance confirmed with p-value ${pValue.toFixed(4)}.`,
    compressionRatios: {
      baseline: 0.7 + Math.random() * 0.2,
      window: 0.3 + Math.random() * 0.3,
      combined: 0.5 + Math.random() * 0.2,
    },
    baselineData: '{"type":"baseline","samples":100}',
    windowData: '{"type":"window","samples":50}',
  };
}
