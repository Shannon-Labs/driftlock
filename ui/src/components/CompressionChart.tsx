'use client';

import { CompressionRatios } from '@/lib/types';

interface CompressionChartProps {
  compressionRatios: CompressionRatios;
}

export default function CompressionChart({
  compressionRatios,
}: CompressionChartProps) {
  const getBarWidth = (ratio: number) => `${ratio * 100}%`;

  const getBarColor = (ratio: number) => {
    if (ratio > 0.7) return 'bg-green-500';
    if (ratio > 0.5) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium text-gray-900 dark:text-white">
        Compression Ratios
      </h3>

      <div className="space-y-3">
        {/* Baseline */}
        <div>
          <div className="flex justify-between text-sm mb-1">
            <span className="text-gray-700 dark:text-gray-300">Baseline</span>
            <span className="font-mono text-gray-900 dark:text-white">
              {compressionRatios.baseline.toFixed(3)}
            </span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
            <div
              className={`h-4 rounded-full ${getBarColor(
                compressionRatios.baseline
              )}`}
              style={{ width: getBarWidth(compressionRatios.baseline) }}
            />
          </div>
        </div>

        {/* Window */}
        <div>
          <div className="flex justify-between text-sm mb-1">
            <span className="text-gray-700 dark:text-gray-300">Window</span>
            <span className="font-mono text-gray-900 dark:text-white">
              {compressionRatios.window.toFixed(3)}
            </span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
            <div
              className={`h-4 rounded-full ${getBarColor(
                compressionRatios.window
              )}`}
              style={{ width: getBarWidth(compressionRatios.window) }}
            />
          </div>
        </div>

        {/* Combined */}
        <div>
          <div className="flex justify-between text-sm mb-1">
            <span className="text-gray-700 dark:text-gray-300">Combined</span>
            <span className="font-mono text-gray-900 dark:text-white">
              {compressionRatios.combined.toFixed(3)}
            </span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
            <div
              className={`h-4 rounded-full ${getBarColor(
                compressionRatios.combined
              )}`}
              style={{ width: getBarWidth(compressionRatios.combined) }}
            />
          </div>
        </div>
      </div>

      <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
        <p className="text-sm text-gray-600 dark:text-gray-400">
          <strong>Interpretation:</strong> A significant drop in compression ratio
          from baseline to window indicates structural changes in the data pattern,
          suggesting an anomaly.
        </p>
        <div className="mt-2 text-xs text-gray-500 dark:text-gray-500">
          <div className="flex items-center space-x-4">
            <div className="flex items-center">
              <div className="w-3 h-3 bg-green-500 rounded mr-1"></div>
              <span>&gt; 0.7 (Good)</span>
            </div>
            <div className="flex items-center">
              <div className="w-3 h-3 bg-yellow-500 rounded mr-1"></div>
              <span>0.5-0.7 (Moderate)</span>
            </div>
            <div className="flex items-center">
              <div className="w-3 h-3 bg-red-500 rounded mr-1"></div>
              <span>&lt; 0.5 (Poor)</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
