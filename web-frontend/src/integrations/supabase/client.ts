// This file is kept for backward compatibility but Supabase is no longer required for OSS deployments.
// The dashboard now uses API key-based authentication instead.
// 
// If you need Supabase integration for compliance features, configure it via environment variables:
// - VITE_SUPABASE_URL
// - VITE_SUPABASE_ANON_KEY

import { createClient } from '@supabase/supabase-js';
import type { Database } from './types';

const SUPABASE_URL = import.meta.env.VITE_SUPABASE_URL || "";
const SUPABASE_PUBLISHABLE_KEY = import.meta.env.VITE_SUPABASE_ANON_KEY || "";

// Only create client if both URL and key are provided
let supabase: ReturnType<typeof createClient<Database>> | null = null;

if (SUPABASE_URL && SUPABASE_PUBLISHABLE_KEY) {
  supabase = createClient<Database>(SUPABASE_URL, SUPABASE_PUBLISHABLE_KEY, {
    auth: {
      storage: localStorage,
      persistSession: true,
      autoRefreshToken: true,
    }
  });
} else {
  // Create a no-op client that will fail gracefully
  supabase = createClient<Database>("https://placeholder.supabase.co", "placeholder-key", {
    auth: {
      storage: localStorage,
      persistSession: false,
      autoRefreshToken: false,
    }
  });
}

export { supabase };