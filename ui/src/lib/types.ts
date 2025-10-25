// Core anomaly detection types matching the backend Go API

export interface Anomaly {
  id: string;
  timestamp: string;
  streamType: 'logs' | 'metrics' | 'traces';
  ncdScore: number;
  pValue: number;
  status: 'pending' | 'acknowledged' | 'dismissed';
  glassBoxExplanation: string;
  compressionRatios: CompressionRatios;
  baselineData: string;
  windowData: string;
  metadata?: Record<string, any>;
}

export interface CompressionRatios {
  baseline: number;
  window: number;
  combined: number;
}

export interface AnomalyFilters {
  startDate?: Date;
  endDate?: Date;
  streamType?: 'logs' | 'metrics' | 'traces' | 'all';
  pValueThreshold?: number;
  ncdThreshold?: number;
  status?: 'pending' | 'acknowledged' | 'dismissed' | 'all';
  searchQuery?: string;
}

export interface PaginatedAnomalies {
  anomalies: Anomaly[];
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
}

export interface DetectionConfig {
  pValueThreshold: number;
  ncdThreshold: number;
  windowSize: number;
  baselineSize: number;
  hopSize: number;
  enabledStreams: {
    logs: boolean;
    metrics: boolean;
    traces: boolean;
  };
}

export interface PerformanceMetrics {
  eventsPerSecond: number;
  averageLatencyMs: number;
  totalAnomalies: number;
  detectionRate: number;
  falsePositiveRate?: number;
}

export interface CompressionDataPoint {
  timestamp: string;
  compressionRatio: number;
  isAnomaly: boolean;
  ncdScore?: number;
  streamType: string;
}

export interface NCDHeatmapData {
  streamType: string;
  timeSlot: string;
  ncdScore: number;
  count: number;
}

export interface StatisticalSummary {
  detectionRateOverTime: Array<{
    timestamp: string;
    rate: number;
  }>;
  falsePositiveTracking: Array<{
    timestamp: string;
    rate: number;
  }>;
  compressionEfficiency: {
    averageRatio: number;
    minRatio: number;
    maxRatio: number;
  };
  performanceStats: PerformanceMetrics;
}
