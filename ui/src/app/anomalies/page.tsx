'use client';

import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api';
import { Anomaly, AnomalyFilters } from '@/lib/types';
import AnomalyTable from '@/components/AnomalyTable';

export default function AnomaliesPage() {
  const [anomalies, setAnomalies] = useState<Anomaly[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [total, setTotal] = useState(0);

  // Filter state
  const [filters, setFilters] = useState<AnomalyFilters>({
    streamType: 'all',
    status: 'all',
    pValueThreshold: 0.05,
    searchQuery: '',
  });

  const [searchInput, setSearchInput] = useState('');

  const loadAnomalies = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await apiClient.getAnomalies(filters, page, 50);
      setAnomalies(result.anomalies);
      setHasMore(result.hasMore);
      setTotal(result.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load anomalies');
      // Use mock data for development
      setAnomalies(getMockAnomalies());
      setTotal(5);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, [filters, page]);

  useEffect(() => {
    loadAnomalies();
  }, [loadAnomalies]);

  const handleSearch = () => {
    setFilters({ ...filters, searchQuery: searchInput });
    setPage(1);
  };

  const handleFilterChange = (key: keyof AnomalyFilters, value: any) => {
    setFilters({ ...filters, [key]: value });
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
            Anomalies
          </h2>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {total} total anomalies detected
          </p>
        </div>
        <button
          onClick={loadAnomalies}
          className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          Refresh
        </button>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Stream Type Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Stream Type
            </label>
            <select
              value={filters.streamType || 'all'}
              onChange={(e) => handleFilterChange('streamType', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            >
              <option value="all">All Streams</option>
              <option value="logs">Logs</option>
              <option value="metrics">Metrics</option>
              <option value="traces">Traces</option>
            </select>
          </div>

          {/* Status Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Status
            </label>
            <select
              value={filters.status || 'all'}
              onChange={(e) => handleFilterChange('status', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            >
              <option value="all">All Status</option>
              <option value="pending">Pending</option>
              <option value="acknowledged">Acknowledged</option>
              <option value="dismissed">Dismissed</option>
            </select>
          </div>

          {/* P-Value Threshold */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              P-Value Threshold: {filters.pValueThreshold?.toFixed(3)}
            </label>
            <input
              type="range"
              min="0.01"
              max="0.1"
              step="0.01"
              value={filters.pValueThreshold || 0.05}
              onChange={(e) =>
                handleFilterChange('pValueThreshold', parseFloat(e.target.value))
              }
              className="w-full"
            />
          </div>

          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Search
            </label>
            <div className="flex">
              <input
                type="text"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Search explanations..."
                className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-l-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              />
              <button
                onClick={handleSearch}
                className="px-4 py-2 bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-200 rounded-r-md hover:bg-gray-300 dark:hover:bg-gray-500"
              >
                Search
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md p-4">
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            Warning: {error}. Showing mock data for development.
          </p>
        </div>
      )}

      {/* Table */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
        {loading ? (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            Loading anomalies...
          </div>
        ) : (
          <AnomalyTable anomalies={anomalies} />
        )}
      </div>

      {/* Pagination */}
      {!loading && anomalies.length > 0 && (
        <div className="flex justify-between items-center">
          <button
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
            className="px-4 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700"
          >
            Previous
          </button>
          <span className="text-sm text-gray-600 dark:text-gray-400">
            Page {page}
          </span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={!hasMore}
            className="px-4 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}

// Mock data for development
function getMockAnomalies(): Anomaly[] {
  return [
    {
      id: '1',
      timestamp: new Date(Date.now() - 3600000).toISOString(),
      streamType: 'logs',
      ncdScore: 0.4523,
      pValue: 0.0087,
      status: 'pending',
      glassBoxExplanation:
        'High NCD score detected: window data differs significantly from baseline. Compression ratio dropped from 0.82 to 0.35, indicating structural change in log patterns.',
      compressionRatios: {
        baseline: 0.82,
        window: 0.35,
        combined: 0.51,
      },
      baselineData: '{"type":"baseline","samples":100}',
      windowData: '{"type":"window","samples":50}',
    },
    {
      id: '2',
      timestamp: new Date(Date.now() - 7200000).toISOString(),
      streamType: 'metrics',
      ncdScore: 0.3241,
      pValue: 0.0234,
      status: 'acknowledged',
      glassBoxExplanation:
        'Moderate anomaly: metrics compression efficiency decreased. P-value 0.0234 indicates statistical significance.',
      compressionRatios: {
        baseline: 0.75,
        window: 0.48,
        combined: 0.61,
      },
      baselineData: '{"type":"baseline","samples":100}',
      windowData: '{"type":"window","samples":50}',
    },
    {
      id: '3',
      timestamp: new Date(Date.now() - 10800000).toISOString(),
      streamType: 'traces',
      ncdScore: 0.5123,
      pValue: 0.0012,
      status: 'pending',
      glassBoxExplanation:
        'Critical anomaly detected: trace patterns show unprecedented deviation. NCD score 0.5123 exceeds normal range significantly.',
      compressionRatios: {
        baseline: 0.88,
        window: 0.29,
        combined: 0.47,
      },
      baselineData: '{"type":"baseline","samples":100}',
      windowData: '{"type":"window","samples":50}',
    },
    {
      id: '4',
      timestamp: new Date(Date.now() - 14400000).toISOString(),
      streamType: 'logs',
      ncdScore: 0.2876,
      pValue: 0.0421,
      status: 'dismissed',
      glassBoxExplanation:
        'Minor anomaly: slight deviation detected but within acceptable threshold. Marked as false positive.',
      compressionRatios: {
        baseline: 0.79,
        window: 0.62,
        combined: 0.71,
      },
      baselineData: '{"type":"baseline","samples":100}',
      windowData: '{"type":"window","samples":50}',
    },
    {
      id: '5',
      timestamp: new Date(Date.now() - 18000000).toISOString(),
      streamType: 'metrics',
      ncdScore: 0.4021,
      pValue: 0.0156,
      status: 'pending',
      glassBoxExplanation:
        'Anomaly in metric patterns: compression degradation suggests data quality issues or genuine anomaly.',
      compressionRatios: {
        baseline: 0.81,
        window: 0.42,
        combined: 0.58,
      },
      baselineData: '{"type":"baseline","samples":100}',
      windowData: '{"type":"window","samples":50}',
    },
  ];
}
