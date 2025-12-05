// API Key interface - matches backend response
export interface APIKey {
  id: string
  name: string
  prefix: string
  key?: string  // Only present on creation
  role?: 'admin' | 'stream'  // Not returned by /v1/me/keys
  stream_id?: string
  created_at: string
  last_used_at?: string
  rate_limit_rps?: number
  status: 'active' | 'revoked'  // Returned by /v1/me/keys
}

// Firebase User subset - only fields we actually use
export interface AuthUser {
  uid: string
  email: string | null
  emailVerified: boolean
  displayName?: string | null
  getIdToken: (forceRefresh?: boolean) => Promise<string>
}
