'use client';

interface GlassBoxExplanationProps {
  explanation: string;
}

export default function GlassBoxExplanation({
  explanation,
}: GlassBoxExplanationProps) {
  return (
    <div className="bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 p-6 rounded-r-lg">
      <div className="flex items-start">
        <div className="flex-shrink-0">
          <svg
            className="h-6 w-6 text-blue-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>
        <div className="ml-3 flex-1">
          <h3 className="text-lg font-medium text-blue-900 dark:text-blue-100">
            Glass-Box Explanation
          </h3>
          <div className="mt-2 text-sm text-blue-800 dark:text-blue-200">
            <p>{explanation}</p>
          </div>
          <div className="mt-4">
            <p className="text-xs text-blue-700 dark:text-blue-300 italic">
              This explanation is generated from the mathematical properties of the
              compression-based anomaly detection algorithm, providing transparent
              reasoning for this detection.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
