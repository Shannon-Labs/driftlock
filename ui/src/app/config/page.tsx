'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import { DetectionConfig } from '@/lib/types';

export default function ConfigPage() {
  const [config, setConfig] = useState<DetectionConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiClient.getConfig();
      setConfig(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load config');
      // Use mock data for development
      setConfig(getMockConfig());
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!config) return;

    try {
      setSaving(true);
      setError(null);
      setSuccess(null);
      const updated = await apiClient.updateConfig(config);
      setConfig(updated);
      setSuccess('Configuration saved successfully!');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save config');
      // For development, just show success
      setSuccess('Configuration saved (mock mode)');
      setTimeout(() => setSuccess(null), 3000);
    } finally {
      setSaving(false);
    }
  };

  const handleReset = () => {
    setConfig(getMockConfig());
  };

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500 dark:text-gray-400">
          Loading configuration...
        </div>
      </div>
    );
  }

  if (!config) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500">Failed to load configuration</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
            Configuration
          </h2>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage detection thresholds and system settings
          </p>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={handleReset}
            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600"
          >
            Reset
          </button>
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md p-4">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {success && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-md p-4">
          <p className="text-sm text-green-800 dark:text-green-200">{success}</p>
        </div>
      )}

      {/* Detection Thresholds */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Detection Thresholds
        </h3>
        <div className="space-y-6">
          {/* P-Value Threshold */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              P-Value Threshold: {config.pValueThreshold.toFixed(3)}
            </label>
            <input
              type="range"
              min="0.001"
              max="0.1"
              step="0.001"
              value={config.pValueThreshold}
              onChange={(e) =>
                setConfig({
                  ...config,
                  pValueThreshold: parseFloat(e.target.value),
                })
              }
              className="w-full"
            />
            <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
              <span>0.001 (Very strict)</span>
              <span>0.05 (Standard)</span>
              <span>0.1 (Lenient)</span>
            </div>
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              Lower values reduce false positives but may miss subtle anomalies.
            </p>
          </div>

          {/* NCD Threshold */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              NCD Threshold: {config.ncdThreshold.toFixed(3)}
            </label>
            <input
              type="range"
              min="0.1"
              max="0.8"
              step="0.01"
              value={config.ncdThreshold}
              onChange={(e) =>
                setConfig({
                  ...config,
                  ncdThreshold: parseFloat(e.target.value),
                })
              }
              className="w-full"
            />
            <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
              <span>0.1 (Very sensitive)</span>
              <span>0.3 (Balanced)</span>
              <span>0.8 (Very conservative)</span>
            </div>
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              Minimum NCD score required to flag an anomaly.
            </p>
          </div>
        </div>
      </div>

      {/* Window Configuration */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Window Configuration
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Baseline Size
            </label>
            <input
              type="number"
              min="10"
              max="1000"
              value={config.baselineSize}
              onChange={(e) =>
                setConfig({
                  ...config,
                  baselineSize: parseInt(e.target.value),
                })
              }
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Number of samples in baseline window
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Window Size
            </label>
            <input
              type="number"
              min="5"
              max="500"
              value={config.windowSize}
              onChange={(e) =>
                setConfig({ ...config, windowSize: parseInt(e.target.value) })
              }
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Number of samples in detection window
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Hop Size
            </label>
            <input
              type="number"
              min="1"
              max="100"
              value={config.hopSize}
              onChange={(e) =>
                setConfig({ ...config, hopSize: parseInt(e.target.value) })
              }
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Window slide interval
            </p>
          </div>
        </div>
      </div>

      {/* Stream Management */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Stream Management
        </h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-900 rounded">
            <div>
              <div className="font-medium text-gray-900 dark:text-white">Logs</div>
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Monitor log streams for anomalies
              </div>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={config.enabledStreams.logs}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    enabledStreams: {
                      ...config.enabledStreams,
                      logs: e.target.checked,
                    },
                  })
                }
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 dark:peer-focus:ring-indigo-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-indigo-600"></div>
            </label>
          </div>

          <div className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-900 rounded">
            <div>
              <div className="font-medium text-gray-900 dark:text-white">
                Metrics
              </div>
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Monitor metric streams for anomalies
              </div>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={config.enabledStreams.metrics}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    enabledStreams: {
                      ...config.enabledStreams,
                      metrics: e.target.checked,
                    },
                  })
                }
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 dark:peer-focus:ring-indigo-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-indigo-600"></div>
            </label>
          </div>

          <div className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-900 rounded">
            <div>
              <div className="font-medium text-gray-900 dark:text-white">
                Traces
              </div>
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Monitor trace streams for anomalies
              </div>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={config.enabledStreams.traces}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    enabledStreams: {
                      ...config.enabledStreams,
                      traces: e.target.checked,
                    },
                  })
                }
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 dark:peer-focus:ring-indigo-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-indigo-600"></div>
            </label>
          </div>
        </div>
      </div>

      {/* Configuration Preview */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Configuration Preview
        </h3>
        <pre className="bg-gray-900 text-green-400 p-4 rounded overflow-x-auto text-sm">
          {JSON.stringify(config, null, 2)}
        </pre>
      </div>
    </div>
  );
}

function getMockConfig(): DetectionConfig {
  return {
    pValueThreshold: 0.05,
    ncdThreshold: 0.3,
    windowSize: 50,
    baselineSize: 100,
    hopSize: 10,
    enabledStreams: {
      logs: true,
      metrics: true,
      traces: true,
    },
  };
}
