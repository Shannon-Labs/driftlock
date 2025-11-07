// API client for DriftLock API
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface ApiError {
  message: string;
  status?: number;
}

export class ApiClient {
  private apiKey: string | null = null;
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
    // Load API key from localStorage if available
    const storedKey = localStorage.getItem('driftlock_api_key');
    if (storedKey) {
      this.apiKey = storedKey;
    }
  }

  setApiKey(key: string | null) {
    this.apiKey = key;
    if (key) {
      localStorage.setItem('driftlock_api_key', key);
    } else {
      localStorage.removeItem('driftlock_api_key');
    }
  }

  getApiKey(): string | null {
    return this.apiKey;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.apiKey) {
      headers['Authorization'] = `Bearer ${this.apiKey}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error: ApiError = {
        message: `API request failed: ${response.statusText}`,
        status: response.status,
      };
      try {
        const errorData = await response.json();
        error.message = errorData.error || error.message;
      } catch {
        // Ignore JSON parse errors
      }
      throw error;
    }

    // Handle empty responses
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return response.json();
    }
    return {} as T;
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async patch<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

