export interface DriftlockClientOptions {
  apiKey: string;
  baseUrl?: string;
  timeoutMs?: number;
  fetchImpl?: typeof fetch;
}

export interface DetectRequest {
  stream_id: string;
  events: unknown[];
  request_id?: string;
  config_override?: Record<string, unknown>;
}

export interface DetectResponse {
  success?: boolean;
  anomaly_count?: number;
  anomalies?: Anomaly[];
  metrics?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface Anomaly {
  id: string;
  stream_id?: string;
  ncd?: number;
  p_value?: number;
  compression_ratio?: number;
  entropy_change?: number;
  confidence?: number;
  explanation?: string;
  status?: string;
  detected_at?: string;
  metrics?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface HealthResponse {
  license?: Record<string, unknown>;
  database?: string;
  queue?: Record<string, unknown>;
  cbad?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface ListAnomaliesParams {
  limit?: number;
  page_token?: string;
  stream_id?: string;
  min_ncd?: number;
  max_p_value?: number;
  status?: string;
  since?: string;
  until?: string;
  has_evidence?: boolean;
}

export class DriftlockError extends Error {
  public status?: number;
  public details?: unknown;

  constructor(message: string, status?: number, details?: unknown) {
    super(message);
    this.name = "DriftlockError";
    this.status = status;
    this.details = details;
  }
}

export class DriftlockClient {
  private readonly apiKey: string;
  private readonly baseUrl: string;
  private readonly timeoutMs: number;
  private readonly fetchImpl: typeof fetch;

  constructor(options: DriftlockClientOptions) {
    if (!options.apiKey) {
      throw new DriftlockError("apiKey is required");
    }
    this.apiKey = options.apiKey;
    this.baseUrl = (options.baseUrl ?? "http://localhost:8080").replace(/\/+$/, "");
    this.timeoutMs = options.timeoutMs ?? 10_000;
    const impl = options.fetchImpl ?? globalThis.fetch;
    if (!impl) {
      throw new DriftlockError("No fetch implementation found. Provide fetchImpl or use Node 18+ / browser.");
    }
    this.fetchImpl = impl;
  }

  async health(): Promise<HealthResponse> {
    return this.request<HealthResponse>("/healthz", { method: "GET" });
  }

  async detect(payload: DetectRequest): Promise<DetectResponse> {
    return this.request<DetectResponse>("/v1/detect", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
  }

  async listAnomalies(params: ListAnomaliesParams = {}): Promise<Record<string, unknown>> {
    const qs = new URLSearchParams(
      Object.entries(params)
        .filter(([, value]) => value !== undefined && value !== null)
        .map(([key, value]) => [key, String(value)])
    );
    const suffix = qs.toString() ? `?${qs.toString()}` : "";
    return this.request<Record<string, unknown>>(`/v1/anomalies${suffix}`, { method: "GET" });
  }

  async getAnomaly(id: string): Promise<Anomaly> {
    if (!id) {
      throw new DriftlockError("anomaly id is required");
    }
    return this.request<Anomaly>(`/v1/anomalies/${encodeURIComponent(id)}`, { method: "GET" });
  }

  private async request<T>(path: string, init: RequestInit): Promise<T> {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);
    try {
      const response = await this.fetchImpl(`${this.baseUrl}${path}`, {
        ...init,
        headers: {
          ...(init.headers ?? {}),
          "X-Api-Key": this.apiKey
        },
        signal: controller.signal
      });

      const contentType = response.headers.get("content-type") || "";
      const isJson = contentType.includes("application/json");
      const payload = isJson ? await response.json() : await response.text();

      if (!response.ok) {
        throw new DriftlockError(
          `Request failed with status ${response.status}`,
          response.status,
          payload
        );
      }

      return payload as T;
    } catch (err) {
      if (err instanceof DriftlockError) {
        throw err;
      }
      if ((err as Error).name === "AbortError") {
        throw new DriftlockError("Request timed out", undefined, { timeoutMs: this.timeoutMs });
      }
      throw new DriftlockError("Network or parsing error", undefined, { cause: err });
    } finally {
      clearTimeout(timer);
    }
  }
}
