import type {
  Anomaly,
  PaginatedAnomalies,
  AnomalyFilters,
  DetectionConfig,
  PerformanceMetrics,
  CompressionDataPoint,
  NCDHeatmapData,
  StatisticalSummary,
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async fetch<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  // Anomaly endpoints
  async getAnomalies(
    filters?: AnomalyFilters,
    page: number = 1,
    pageSize: number = 50
  ): Promise<PaginatedAnomalies> {
    const params = new URLSearchParams({
      page: page.toString(),
      pageSize: pageSize.toString(),
    });

    if (filters?.startDate) {
      params.append('startDate', filters.startDate.toISOString());
    }
    if (filters?.endDate) {
      params.append('endDate', filters.endDate.toISOString());
    }
    if (filters?.streamType && filters.streamType !== 'all') {
      params.append('streamType', filters.streamType);
    }
    if (filters?.pValueThreshold) {
      params.append('pValueThreshold', filters.pValueThreshold.toString());
    }
    if (filters?.ncdThreshold) {
      params.append('ncdThreshold', filters.ncdThreshold.toString());
    }
    if (filters?.status && filters.status !== 'all') {
      params.append('status', filters.status);
    }
    if (filters?.searchQuery) {
      params.append('search', filters.searchQuery);
    }

    return this.fetch<PaginatedAnomalies>(`/v1/anomalies?${params.toString()}`);
  }

  async getAnomaly(id: string): Promise<Anomaly> {
    return this.fetch<Anomaly>(`/v1/anomalies/${id}`);
  }

  async updateAnomalyStatus(
    id: string,
    status: 'acknowledged' | 'dismissed'
  ): Promise<Anomaly> {
    return this.fetch<Anomaly>(`/v1/anomalies/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    });
  }

  async exportAnomaly(id: string): Promise<Blob> {
    const response = await fetch(`${this.baseUrl}/v1/anomalies/${id}/export`);
    return response.blob();
  }

  // Configuration endpoints
  async getConfig(): Promise<DetectionConfig> {
    return this.fetch<DetectionConfig>('/v1/config');
  }

  async updateConfig(config: Partial<DetectionConfig>): Promise<DetectionConfig> {
    return this.fetch<DetectionConfig>('/v1/config', {
      method: 'PATCH',
      body: JSON.stringify(config),
    });
  }

  // Analytics endpoints
  async getPerformanceMetrics(): Promise<PerformanceMetrics> {
    return this.fetch<PerformanceMetrics>('/v1/metrics/performance');
  }

  async getCompressionTimeline(
    streamType?: string,
    hours: number = 24
  ): Promise<CompressionDataPoint[]> {
    const params = new URLSearchParams({ hours: hours.toString() });
    if (streamType) {
      params.append('streamType', streamType);
    }
    return this.fetch<CompressionDataPoint[]>(
      `/v1/analytics/compression-timeline?${params.toString()}`
    );
  }

  async getNCDHeatmap(hours: number = 24): Promise<NCDHeatmapData[]> {
    return this.fetch<NCDHeatmapData[]>(
      `/v1/analytics/ncd-heatmap?hours=${hours}`
    );
  }

  async getStatisticalSummary(): Promise<StatisticalSummary> {
    return this.fetch<StatisticalSummary>('/v1/analytics/summary');
  }

  // Server-Sent Events for real-time updates
  createAnomalyStream(): EventSource {
    return new EventSource(`${this.baseUrl}/v1/stream/anomalies`);
  }
}

export const apiClient = new ApiClient();
