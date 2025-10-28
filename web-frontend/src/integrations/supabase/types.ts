export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json | undefined }
  | Json[]

export type Database = {
  // Allows to automatically instantiate createClient with right options
  // instead of createClient<Database, { PostgrestVersion: 'XX' }>(URL, KEY)
  __InternalSupabase: {
    PostgrestVersion: "13.0.5"
  }
  public: {
    Tables: {
      anomaly_events: {
        Row: {
          anomaly_type: string
          created_at: string | null
          description: string | null
          detection_timestamp: string | null
          explanation: string | null
          id: string
          organization_id: string
          raw_data: Json | null
          resolved_at: string | null
          resolved_by: string | null
          severity: string
        }
        Insert: {
          anomaly_type: string
          created_at?: string | null
          description?: string | null
          detection_timestamp?: string | null
          explanation?: string | null
          id?: string
          organization_id: string
          raw_data?: Json | null
          resolved_at?: string | null
          resolved_by?: string | null
          severity: string
        }
        Update: {
          anomaly_type?: string
          created_at?: string | null
          description?: string | null
          detection_timestamp?: string | null
          explanation?: string | null
          id?: string
          organization_id?: string
          raw_data?: Json | null
          resolved_at?: string | null
          resolved_by?: string | null
          severity?: string
        }
        Relationships: [
          {
            foreignKeyName: "anomaly_events_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      api_keys: {
        Row: {
          created_at: string
          created_by: string | null
          id: string
          is_active: boolean
          key_hash: string
          key_name: string | null
          key_prefix: string
          last_used_at: string | null
          name: string
          organization_id: string | null
          permissions: Json | null
          tier: string
          updated_at: string
          user_id: string
        }
        Insert: {
          created_at?: string
          created_by?: string | null
          id?: string
          is_active?: boolean
          key_hash: string
          key_name?: string | null
          key_prefix: string
          last_used_at?: string | null
          name: string
          organization_id?: string | null
          permissions?: Json | null
          tier?: string
          updated_at?: string
          user_id: string
        }
        Update: {
          created_at?: string
          created_by?: string | null
          id?: string
          is_active?: boolean
          key_hash?: string
          key_name?: string | null
          key_prefix?: string
          last_used_at?: string | null
          name?: string
          organization_id?: string | null
          permissions?: Json | null
          tier?: string
          updated_at?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "api_keys_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "api_keys_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      api_usage: {
        Row: {
          api_key_id: string | null
          created_at: string
          endpoint: string
          id: string
          ip_address: string | null
          metadata: Json | null
          method: string
          request_size_bytes: number | null
          response_size_bytes: number | null
          response_time_ms: number | null
          status_code: number | null
          user_agent: string | null
          user_id: string | null
        }
        Insert: {
          api_key_id?: string | null
          created_at?: string
          endpoint: string
          id?: string
          ip_address?: string | null
          metadata?: Json | null
          method: string
          request_size_bytes?: number | null
          response_size_bytes?: number | null
          response_time_ms?: number | null
          status_code?: number | null
          user_agent?: string | null
          user_id?: string | null
        }
        Update: {
          api_key_id?: string | null
          created_at?: string
          endpoint?: string
          id?: string
          ip_address?: string | null
          metadata?: Json | null
          method?: string
          request_size_bytes?: number | null
          response_size_bytes?: number | null
          response_time_ms?: number | null
          status_code?: number | null
          user_agent?: string | null
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "api_usage_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "profiles"
            referencedColumns: ["id"]
          },
        ]
      }
      api_usage_logs: {
        Row: {
          api_key_id: string | null
          created_at: string
          endpoint: string | null
          id: string
          ip_address: string | null
          method: string | null
          request_size_bytes: number | null
          response_size_bytes: number | null
          response_time_ms: number | null
          status_code: number | null
          user_agent: string | null
          user_id: string | null
        }
        Insert: {
          api_key_id?: string | null
          created_at?: string
          endpoint?: string | null
          id?: string
          ip_address?: string | null
          method?: string | null
          request_size_bytes?: number | null
          response_size_bytes?: number | null
          response_time_ms?: number | null
          status_code?: number | null
          user_agent?: string | null
          user_id?: string | null
        }
        Update: {
          api_key_id?: string | null
          created_at?: string
          endpoint?: string | null
          id?: string
          ip_address?: string | null
          method?: string | null
          request_size_bytes?: number | null
          response_size_bytes?: number | null
          response_time_ms?: number | null
          status_code?: number | null
          user_agent?: string | null
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "api_usage_logs_api_key_id_fkey"
            columns: ["api_key_id"]
            isOneToOne: false
            referencedRelation: "api_keys"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "api_usage_logs_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      articles: {
        Row: {
          created_at: string | null
          hash: string
          id: string
          lang: string | null
          published_at: string | null
          source: string | null
          summary: string | null
          title: string
          updated_at: string | null
          url: string
        }
        Insert: {
          created_at?: string | null
          hash: string
          id?: string
          lang?: string | null
          published_at?: string | null
          source?: string | null
          summary?: string | null
          title: string
          updated_at?: string | null
          url: string
        }
        Update: {
          created_at?: string | null
          hash?: string
          id?: string
          lang?: string | null
          published_at?: string | null
          source?: string | null
          summary?: string | null
          title?: string
          updated_at?: string | null
          url?: string
        }
        Relationships: []
      }
      audit_logs: {
        Row: {
          action: string
          created_at: string | null
          details: Json | null
          id: number
          ip_address: unknown
          organization_id: string | null
          resource_id: string | null
          resource_type: string
          user_agent: string | null
          user_id: string | null
        }
        Insert: {
          action: string
          created_at?: string | null
          details?: Json | null
          id?: number
          ip_address?: unknown
          organization_id?: string | null
          resource_id?: string | null
          resource_type: string
          user_agent?: string | null
          user_id?: string | null
        }
        Update: {
          action?: string
          created_at?: string | null
          details?: Json | null
          id?: number
          ip_address?: unknown
          organization_id?: string | null
          resource_id?: string | null
          resource_type?: string
          user_agent?: string | null
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "audit_logs_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      billing_actions: {
        Row: {
          action: string
          created_at: string | null
          details: Json | null
          id: string
          organization_id: string | null
          result: string
          stripe_event_id: string | null
        }
        Insert: {
          action: string
          created_at?: string | null
          details?: Json | null
          id?: string
          organization_id?: string | null
          result: string
          stripe_event_id?: string | null
        }
        Update: {
          action?: string
          created_at?: string | null
          details?: Json | null
          id?: string
          organization_id?: string | null
          result?: string
          stripe_event_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "billing_actions_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "billing_actions_stripe_event_id_fkey"
            columns: ["stripe_event_id"]
            isOneToOne: false
            referencedRelation: "stripe_events"
            referencedColumns: ["event_id"]
          },
        ]
      }
      billing_customers: {
        Row: {
          billing_email: string | null
          company_name: string | null
          created_at: string | null
          default_payment_method_brand: string | null
          default_payment_method_last4: string | null
          organization_id: string
          stripe_customer_id: string
          tax_country: string | null
          tax_id: string | null
          tax_postal_code: string | null
          updated_at: string | null
        }
        Insert: {
          billing_email?: string | null
          company_name?: string | null
          created_at?: string | null
          default_payment_method_brand?: string | null
          default_payment_method_last4?: string | null
          organization_id: string
          stripe_customer_id: string
          tax_country?: string | null
          tax_id?: string | null
          tax_postal_code?: string | null
          updated_at?: string | null
        }
        Update: {
          billing_email?: string | null
          company_name?: string | null
          created_at?: string | null
          default_payment_method_brand?: string | null
          default_payment_method_last4?: string | null
          organization_id?: string
          stripe_customer_id?: string
          tax_country?: string | null
          tax_id?: string | null
          tax_postal_code?: string | null
          updated_at?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "billing_customers_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: true
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      billing_events: {
        Row: {
          created_at: string | null
          event_type: string
          id: string
          organization_id: string | null
          payload: Json
          processed_at: string | null
          stripe_event_id: string | null
        }
        Insert: {
          created_at?: string | null
          event_type: string
          id?: string
          organization_id?: string | null
          payload: Json
          processed_at?: string | null
          stripe_event_id?: string | null
        }
        Update: {
          created_at?: string | null
          event_type?: string
          id?: string
          organization_id?: string | null
          payload?: Json
          processed_at?: string | null
          stripe_event_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "billing_events_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      credit_purchases: {
        Row: {
          amount_cents: number
          created_at: string
          credits_purchased: number
          credits_remaining: number
          id: string
          status: string | null
          stripe_payment_intent_id: string
          user_id: string
        }
        Insert: {
          amount_cents: number
          created_at?: string
          credits_purchased: number
          credits_remaining: number
          id?: string
          status?: string | null
          stripe_payment_intent_id: string
          user_id: string
        }
        Update: {
          amount_cents?: number
          created_at?: string
          credits_purchased?: number
          credits_remaining?: number
          id?: string
          status?: string | null
          stripe_payment_intent_id?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "credit_purchases_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "profiles"
            referencedColumns: ["id"]
          },
        ]
      }
      credit_transactions: {
        Row: {
          amount: number
          created_at: string
          description: string
          id: string
          metadata: Json | null
          reference_id: string | null
          type: string
          user_id: string
        }
        Insert: {
          amount: number
          created_at?: string
          description: string
          id?: string
          metadata?: Json | null
          reference_id?: string | null
          type: string
          user_id: string
        }
        Update: {
          amount?: number
          created_at?: string
          description?: string
          id?: string
          metadata?: Json | null
          reference_id?: string | null
          type?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "credit_transactions_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      detections: {
        Row: {
          algorithm_used: string | null
          anomalies_found: number
          api_key_id: string | null
          confidence_score: number | null
          created_at: string
          data_points: number
          error_message: string | null
          id: string
          organization_id: string | null
          processing_time_ms: number | null
          result: Json | null
          status: string | null
          updated_at: string
          user_id: string
        }
        Insert: {
          algorithm_used?: string | null
          anomalies_found?: number
          api_key_id?: string | null
          confidence_score?: number | null
          created_at?: string
          data_points: number
          error_message?: string | null
          id?: string
          organization_id?: string | null
          processing_time_ms?: number | null
          result?: Json | null
          status?: string | null
          updated_at?: string
          user_id: string
        }
        Update: {
          algorithm_used?: string | null
          anomalies_found?: number
          api_key_id?: string | null
          confidence_score?: number | null
          created_at?: string
          data_points?: number
          error_message?: string | null
          id?: string
          organization_id?: string | null
          processing_time_ms?: number | null
          result?: Json | null
          status?: string | null
          updated_at?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "detections_api_key_id_fkey"
            columns: ["api_key_id"]
            isOneToOne: false
            referencedRelation: "api_keys"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "detections_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "detections_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      dunning_states: {
        Row: {
          notes: string | null
          organization_id: string
          since: string | null
          state: string | null
          updated_at: string | null
        }
        Insert: {
          notes?: string | null
          organization_id: string
          since?: string | null
          state?: string | null
          updated_at?: string | null
        }
        Update: {
          notes?: string | null
          organization_id?: string
          since?: string | null
          state?: string | null
          updated_at?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "dunning_states_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: true
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      invoices_mirror: {
        Row: {
          amount_due_cents: number
          amount_paid_cents: number
          created_at: string | null
          finalized_at: string | null
          hosted_invoice_url: string | null
          id: string
          invoice_pdf_url: string | null
          organization_id: string | null
          paid_at: string | null
          period_end: string | null
          period_start: string | null
          status: string
          stripe_invoice_id: string
        }
        Insert: {
          amount_due_cents: number
          amount_paid_cents: number
          created_at?: string | null
          finalized_at?: string | null
          hosted_invoice_url?: string | null
          id?: string
          invoice_pdf_url?: string | null
          organization_id?: string | null
          paid_at?: string | null
          period_end?: string | null
          period_start?: string | null
          status: string
          stripe_invoice_id: string
        }
        Update: {
          amount_due_cents?: number
          amount_paid_cents?: number
          created_at?: string | null
          finalized_at?: string | null
          hosted_invoice_url?: string | null
          id?: string
          invoice_pdf_url?: string | null
          organization_id?: string | null
          paid_at?: string | null
          period_end?: string | null
          period_start?: string | null
          status?: string
          stripe_invoice_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "invoices_mirror_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      org_settings: {
        Row: {
          anomaly_sensitivity: number | null
          created_at: string | null
          dunning_behavior: string | null
          organization_id: string
          updated_at: string | null
          usage_alert_thresholds: Json | null
        }
        Insert: {
          anomaly_sensitivity?: number | null
          created_at?: string | null
          dunning_behavior?: string | null
          organization_id: string
          updated_at?: string | null
          usage_alert_thresholds?: Json | null
        }
        Update: {
          anomaly_sensitivity?: number | null
          created_at?: string | null
          dunning_behavior?: string | null
          organization_id?: string
          updated_at?: string | null
          usage_alert_thresholds?: Json | null
        }
        Relationships: [
          {
            foreignKeyName: "org_settings_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: true
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      organization_members: {
        Row: {
          joined_at: string | null
          organization_id: string
          role: string
          user_id: string
        }
        Insert: {
          joined_at?: string | null
          organization_id: string
          role: string
          user_id: string
        }
        Update: {
          joined_at?: string | null
          organization_id?: string
          role?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "organization_members_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      organizations: {
        Row: {
          created_at: string | null
          id: string
          name: string
          settings: Json | null
          slug: string
          updated_at: string | null
        }
        Insert: {
          created_at?: string | null
          id?: string
          name: string
          settings?: Json | null
          slug: string
          updated_at?: string | null
        }
        Update: {
          created_at?: string | null
          id?: string
          name?: string
          settings?: Json | null
          slug?: string
          updated_at?: string | null
        }
        Relationships: []
      }
      plan_price_map: {
        Row: {
          created_at: string | null
          currency: string
          id: string
          plan_code: string
          stripe_price_id: string
          stripe_product_id: string | null
        }
        Insert: {
          created_at?: string | null
          currency?: string
          id?: string
          plan_code: string
          stripe_price_id: string
          stripe_product_id?: string | null
        }
        Update: {
          created_at?: string | null
          currency?: string
          id?: string
          plan_code?: string
          stripe_price_id?: string
          stripe_product_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "plan_price_map_plan_code_fkey"
            columns: ["plan_code"]
            isOneToOne: false
            referencedRelation: "plans"
            referencedColumns: ["code"]
          },
        ]
      }
      plans: {
        Row: {
          base_price_cents: number
          code: string
          created_at: string | null
          display_name: string
          features: Json | null
          included_calls: number
          is_active: boolean | null
          overage_rate_cents: number
        }
        Insert: {
          base_price_cents: number
          code: string
          created_at?: string | null
          display_name: string
          features?: Json | null
          included_calls: number
          is_active?: boolean | null
          overage_rate_cents: number
        }
        Update: {
          base_price_cents?: number
          code?: string
          created_at?: string | null
          display_name?: string
          features?: Json | null
          included_calls?: number
          is_active?: boolean | null
          overage_rate_cents?: number
        }
        Relationships: []
      }
      profiles: {
        Row: {
          billing_period_ends: string | null
          created_at: string
          email: string
          grace_period_end: string | null
          id: string
          name: string | null
          plan_tier: Database["public"]["Enums"]["subscription_tier"] | null
          role: string | null
          stripe_customer_id: string | null
          stripe_subscription_id: string | null
          subscription_status: string | null
          updated_at: string
        }
        Insert: {
          billing_period_ends?: string | null
          created_at?: string
          email: string
          grace_period_end?: string | null
          id: string
          name?: string | null
          plan_tier?: Database["public"]["Enums"]["subscription_tier"] | null
          role?: string | null
          stripe_customer_id?: string | null
          stripe_subscription_id?: string | null
          subscription_status?: string | null
          updated_at?: string
        }
        Update: {
          billing_period_ends?: string | null
          created_at?: string
          email?: string
          grace_period_end?: string | null
          id?: string
          name?: string | null
          plan_tier?: Database["public"]["Enums"]["subscription_tier"] | null
          role?: string | null
          stripe_customer_id?: string | null
          stripe_subscription_id?: string | null
          subscription_status?: string | null
          updated_at?: string
        }
        Relationships: []
      }
      promotions: {
        Row: {
          applies_to_plans: string[] | null
          code: string
          created_at: string | null
          ends_at: string | null
          id: string
          is_active: boolean | null
          max_redemptions: number | null
          percent_off: number
          starts_at: string | null
          stripe_promotion_code: string | null
          times_redeemed: number | null
        }
        Insert: {
          applies_to_plans?: string[] | null
          code: string
          created_at?: string | null
          ends_at?: string | null
          id?: string
          is_active?: boolean | null
          max_redemptions?: number | null
          percent_off: number
          starts_at?: string | null
          stripe_promotion_code?: string | null
          times_redeemed?: number | null
        }
        Update: {
          applies_to_plans?: string[] | null
          code?: string
          created_at?: string | null
          ends_at?: string | null
          id?: string
          is_active?: boolean | null
          max_redemptions?: number | null
          percent_off?: number
          starts_at?: string | null
          stripe_promotion_code?: string | null
          times_redeemed?: number | null
        }
        Relationships: []
      }
      quota_policy: {
        Row: {
          alert_100: boolean | null
          alert_70: boolean | null
          alert_90: boolean | null
          behavior_on_exceed: string | null
          cap_percent: number | null
          created_at: string | null
          invoice_threshold_cents: number | null
          last_alert_sent_at: string | null
          last_alert_type: string | null
          organization_id: string
          updated_at: string | null
        }
        Insert: {
          alert_100?: boolean | null
          alert_70?: boolean | null
          alert_90?: boolean | null
          behavior_on_exceed?: string | null
          cap_percent?: number | null
          created_at?: string | null
          invoice_threshold_cents?: number | null
          last_alert_sent_at?: string | null
          last_alert_type?: string | null
          organization_id: string
          updated_at?: string | null
        }
        Update: {
          alert_100?: boolean | null
          alert_70?: boolean | null
          alert_90?: boolean | null
          behavior_on_exceed?: string | null
          cap_percent?: number | null
          created_at?: string | null
          invoice_threshold_cents?: number | null
          last_alert_sent_at?: string | null
          last_alert_type?: string | null
          organization_id?: string
          updated_at?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "quota_policy_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: true
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      stripe_events: {
        Row: {
          event_id: string
          id: string
          payload: Json
          received_at: string | null
          type: string
        }
        Insert: {
          event_id: string
          id?: string
          payload: Json
          received_at?: string | null
          type: string
        }
        Update: {
          event_id?: string
          id?: string
          payload?: Json
          received_at?: string | null
          type?: string
        }
        Relationships: []
      }
      subscription_tiers: {
        Row: {
          display_name: string
          features: Json | null
          monthly_requests: number
          price_cents: number | null
          rate_limit_per_minute: number | null
          stripe_price_id: string | null
          stripe_product_id: string | null
          tier: Database["public"]["Enums"]["subscription_tier"]
        }
        Insert: {
          display_name: string
          features?: Json | null
          monthly_requests?: number
          price_cents?: number | null
          rate_limit_per_minute?: number | null
          stripe_price_id?: string | null
          stripe_product_id?: string | null
          tier: Database["public"]["Enums"]["subscription_tier"]
        }
        Update: {
          display_name?: string
          features?: Json | null
          monthly_requests?: number
          price_cents?: number | null
          rate_limit_per_minute?: number | null
          stripe_price_id?: string | null
          stripe_product_id?: string | null
          tier?: Database["public"]["Enums"]["subscription_tier"]
        }
        Relationships: []
      }
      subscriptions: {
        Row: {
          billing_period_ends: string | null
          created_at: string
          current_period_end: string | null
          current_period_start: string | null
          has_launch50_promo: boolean | null
          id: string
          included_calls: number | null
          organization_id: string | null
          overage_rate_per_call: number | null
          plan: string | null
          price_monitor_id: string | null
          price_stream_id: string | null
          quota: number | null
          status: string | null
          stripe_customer_id: string | null
          stripe_subscription_id: string | null
          tier: Database["public"]["Enums"]["subscription_tier"]
          updated_at: string
          usage_count: number
          user_id: string
        }
        Insert: {
          billing_period_ends?: string | null
          created_at?: string
          current_period_end?: string | null
          current_period_start?: string | null
          has_launch50_promo?: boolean | null
          id?: string
          included_calls?: number | null
          organization_id?: string | null
          overage_rate_per_call?: number | null
          plan?: string | null
          price_monitor_id?: string | null
          price_stream_id?: string | null
          quota?: number | null
          status?: string | null
          stripe_customer_id?: string | null
          stripe_subscription_id?: string | null
          tier?: Database["public"]["Enums"]["subscription_tier"]
          updated_at?: string
          usage_count?: number
          user_id: string
        }
        Update: {
          billing_period_ends?: string | null
          created_at?: string
          current_period_end?: string | null
          current_period_start?: string | null
          has_launch50_promo?: boolean | null
          id?: string
          included_calls?: number | null
          organization_id?: string | null
          overage_rate_per_call?: number | null
          plan?: string | null
          price_monitor_id?: string | null
          price_stream_id?: string | null
          quota?: number | null
          status?: string | null
          stripe_customer_id?: string | null
          stripe_subscription_id?: string | null
          tier?: Database["public"]["Enums"]["subscription_tier"]
          updated_at?: string
          usage_count?: number
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "subscriptions_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
          {
            foreignKeyName: "subscriptions_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: true
            referencedRelation: "profiles"
            referencedColumns: ["id"]
          },
        ]
      }
      usage_counters: {
        Row: {
          estimated_charges_cents: number
          included_calls_used: number
          organization_id: string
          overage_calls: number
          period_end: string
          period_start: string
          total_calls: number
          updated_at: string | null
        }
        Insert: {
          estimated_charges_cents?: number
          included_calls_used?: number
          organization_id: string
          overage_calls?: number
          period_end: string
          period_start: string
          total_calls?: number
          updated_at?: string | null
        }
        Update: {
          estimated_charges_cents?: number
          included_calls_used?: number
          organization_id?: string
          overage_calls?: number
          period_end?: string
          period_start?: string
          total_calls?: number
          updated_at?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "usage_counters_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      user_credits: {
        Row: {
          balance: number
          total_used: number
          updated_at: string
          user_id: string
        }
        Insert: {
          balance?: number
          total_used?: number
          updated_at?: string
          user_id: string
        }
        Update: {
          balance?: number
          total_used?: number
          updated_at?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "user_credits_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: true
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      user_subscriptions: {
        Row: {
          cancel_at_period_end: boolean | null
          created_at: string
          current_period_end: string | null
          current_period_start: string | null
          id: string
          price_id: string | null
          status: string
          stripe_subscription_id: string | null
          updated_at: string
          user_id: string
        }
        Insert: {
          cancel_at_period_end?: boolean | null
          created_at?: string
          current_period_end?: string | null
          current_period_start?: string | null
          id?: string
          price_id?: string | null
          status: string
          stripe_subscription_id?: string | null
          updated_at?: string
          user_id: string
        }
        Update: {
          cancel_at_period_end?: boolean | null
          created_at?: string
          current_period_end?: string | null
          current_period_start?: string | null
          id?: string
          price_id?: string | null
          status?: string
          stripe_subscription_id?: string | null
          updated_at?: string
          user_id?: string
        }
        Relationships: [
          {
            foreignKeyName: "user_subscriptions_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["id"]
          },
        ]
      }
      users: {
        Row: {
          created_at: string
          credits: number
          email: string
          firebase_uid: string | null
          id: string
          is_active: boolean
          is_email_verified: boolean
          last_login_at: string | null
          name: string | null
          stripe_customer_id: string | null
          tier: string
          updated_at: string
        }
        Insert: {
          created_at?: string
          credits?: number
          email: string
          firebase_uid?: string | null
          id?: string
          is_active?: boolean
          is_email_verified?: boolean
          last_login_at?: string | null
          name?: string | null
          stripe_customer_id?: string | null
          tier?: string
          updated_at?: string
        }
        Update: {
          created_at?: string
          credits?: number
          email?: string
          firebase_uid?: string | null
          id?: string
          is_active?: boolean
          is_email_verified?: boolean
          last_login_at?: string | null
          name?: string | null
          stripe_customer_id?: string | null
          tier?: string
          updated_at?: string
        }
        Relationships: []
      }
    }
    Views: {
      v_current_period_usage: {
        Row: {
          behavior_on_exceed: string | null
          cap_percent: number | null
          current_period_end: string | null
          current_period_start: string | null
          days_remaining: number | null
          dunning_state: string | null
          estimated_charges_cents: number | null
          estimated_overage_usd: number | null
          included_calls: number | null
          included_calls_used: number | null
          needs_100_alert: boolean | null
          needs_70_alert: boolean | null
          needs_90_alert: boolean | null
          organization_id: string | null
          overage_calls: number | null
          percent_used: number | null
          plan: string | null
          status: string | null
          total_calls: number | null
        }
        Relationships: [
          {
            foreignKeyName: "subscriptions_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
      v_org_anomaly_summary: {
        Row: {
          critical_anomalies: number | null
          current_period_end: string | null
          current_period_start: string | null
          first_detection_in_period: string | null
          high_anomalies: number | null
          last_detection: string | null
          low_anomalies: number | null
          medium_anomalies: number | null
          organization_id: string | null
          total_anomalies: number | null
        }
        Relationships: [
          {
            foreignKeyName: "anomaly_events_organization_id_fkey"
            columns: ["organization_id"]
            isOneToOne: false
            referencedRelation: "organizations"
            referencedColumns: ["id"]
          },
        ]
      }
    }
    Functions: {
      add_user_credits: {
        Args: { p_amount: number; p_user_id: string }
        Returns: undefined
      }
      check_user_credits: { Args: { p_user_id: string }; Returns: boolean }
      deduct_user_credit: { Args: { p_user_id: string }; Returns: undefined }
      increment_usage: {
        Args: { p_count?: number; p_org: string; p_period_start: string }
        Returns: undefined
      }
      is_org_member: { Args: { org_id: string }; Returns: boolean }
      is_org_member_or_service: { Args: { org_id: string }; Returns: boolean }
    }
    Enums: {
      subscription_tier:
        | "free"
        | "starter"
        | "professional"
        | "enterprise"
        | "scale"
    }
    CompositeTypes: {
      [_ in never]: never
    }
  }
}

type DatabaseWithoutInternals = Omit<Database, "__InternalSupabase">

type DefaultSchema = DatabaseWithoutInternals[Extract<keyof Database, "public">]

export type Tables<
  DefaultSchemaTableNameOrOptions extends
    | keyof (DefaultSchema["Tables"] & DefaultSchema["Views"])
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof (DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
        DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Views"])
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? (DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
      DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Views"])[TableName] extends {
      Row: infer R
    }
    ? R
    : never
  : DefaultSchemaTableNameOrOptions extends keyof (DefaultSchema["Tables"] &
        DefaultSchema["Views"])
    ? (DefaultSchema["Tables"] &
        DefaultSchema["Views"])[DefaultSchemaTableNameOrOptions] extends {
        Row: infer R
      }
      ? R
      : never
    : never

export type TablesInsert<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Insert: infer I
    }
    ? I
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Insert: infer I
      }
      ? I
      : never
    : never

export type TablesUpdate<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Update: infer U
    }
    ? U
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Update: infer U
      }
      ? U
      : never
    : never

export type Enums<
  DefaultSchemaEnumNameOrOptions extends
    | keyof DefaultSchema["Enums"]
    | { schema: keyof DatabaseWithoutInternals },
  EnumName extends DefaultSchemaEnumNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"]
    : never = never,
> = DefaultSchemaEnumNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"][EnumName]
  : DefaultSchemaEnumNameOrOptions extends keyof DefaultSchema["Enums"]
    ? DefaultSchema["Enums"][DefaultSchemaEnumNameOrOptions]
    : never

export type CompositeTypes<
  PublicCompositeTypeNameOrOptions extends
    | keyof DefaultSchema["CompositeTypes"]
    | { schema: keyof DatabaseWithoutInternals },
  CompositeTypeName extends PublicCompositeTypeNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"]
    : never = never,
> = PublicCompositeTypeNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"][CompositeTypeName]
  : PublicCompositeTypeNameOrOptions extends keyof DefaultSchema["CompositeTypes"]
    ? DefaultSchema["CompositeTypes"][PublicCompositeTypeNameOrOptions]
    : never

export const Constants = {
  public: {
    Enums: {
      subscription_tier: [
        "free",
        "starter",
        "professional",
        "enterprise",
        "scale",
      ],
    },
  },
} as const
