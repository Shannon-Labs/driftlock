'use client';

import { Anomaly } from '@/lib/types';
import Link from 'next/link';

interface AnomalyTableProps {
  anomalies: Anomaly[];
}

export default function AnomalyTable({ anomalies }: AnomalyTableProps) {
  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  const getStatusBadge = (status: string) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
      acknowledged: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
      dismissed: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
    };
    return colors[status as keyof typeof colors] || colors.pending;
  };

  const getStreamBadge = (streamType: string) => {
    const colors = {
      logs: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
      metrics: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
      traces: 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200',
    };
    return colors[streamType as keyof typeof colors] || colors.logs;
  };

  const getSeverityColor = (pValue: number) => {
    if (pValue < 0.01) return 'text-red-600 dark:text-red-400 font-semibold';
    if (pValue < 0.05) return 'text-orange-600 dark:text-orange-400 font-semibold';
    return 'text-yellow-600 dark:text-yellow-400';
  };

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Timestamp
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Stream Type
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              NCD Score
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              P-Value
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Status
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
          {anomalies.map((anomaly) => (
            <tr key={anomaly.id} className="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                {formatTimestamp(anomaly.timestamp)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span
                  className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStreamBadge(
                    anomaly.streamType
                  )}`}
                >
                  {anomaly.streamType}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm">
                <span className="font-mono text-gray-900 dark:text-gray-100">
                  {anomaly.ncdScore.toFixed(4)}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm">
                <span className={`font-mono ${getSeverityColor(anomaly.pValue)}`}>
                  {anomaly.pValue.toFixed(4)}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span
                  className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusBadge(
                    anomaly.status
                  )}`}
                >
                  {anomaly.status}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <Link
                  href={`/anomalies/${anomaly.id}`}
                  className="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300"
                >
                  View Details
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {anomalies.length === 0 && (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No anomalies found
        </div>
      )}
    </div>
  );
}
