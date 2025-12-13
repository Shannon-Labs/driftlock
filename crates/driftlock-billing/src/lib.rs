//! Driftlock Billing - Stripe integration

use stripe::{
    BillingPortalSession, CheckoutSession, CheckoutSessionMode, Client, CreateBillingPortalSession,
    CreateCheckoutSession, CreateCheckoutSessionLineItems, Customer, CustomerId, Webhook,
};
use thiserror::Error;
use uuid::Uuid;

#[derive(Error, Debug)]
pub enum BillingError {
    #[error("Stripe error: {0}")]
    StripeError(#[from] stripe::StripeError),
    #[error("Invalid webhook signature")]
    InvalidSignature,
    #[error("Configuration error: {0}")]
    ConfigError(String),
}

pub struct StripeClient {
    client: Client,
    webhook_secret: Option<String>,
}

impl StripeClient {
    pub fn new(secret_key: &str, webhook_secret: Option<&str>) -> Self {
        Self {
            client: Client::new(secret_key),
            webhook_secret: webhook_secret.map(String::from),
        }
    }

    /// Create a checkout session for subscription
    pub async fn create_checkout_session(
        &self,
        customer_email: &str,
        price_id: &str,
        success_url: &str,
        cancel_url: &str,
        tenant_id: Uuid,
    ) -> Result<CheckoutSession, BillingError> {
        let tenant_id_str = tenant_id.to_string();
        let mut params = CreateCheckoutSession::new();
        params.customer_email = Some(customer_email);
        params.mode = Some(CheckoutSessionMode::Subscription);
        params.success_url = Some(success_url);
        params.cancel_url = Some(cancel_url);
        params.line_items = Some(vec![CreateCheckoutSessionLineItems {
            price: Some(price_id.to_string()),
            quantity: Some(1),
            ..Default::default()
        }]);
        params.client_reference_id = Some(&tenant_id_str);

        let session = CheckoutSession::create(&self.client, params).await?;
        Ok(session)
    }

    /// Create a customer portal session
    pub async fn create_portal_session(
        &self,
        customer_id: &str,
        return_url: &str,
    ) -> Result<BillingPortalSession, BillingError> {
        let customer_id: CustomerId = customer_id
            .parse()
            .map_err(|_| BillingError::ConfigError("Invalid customer ID".to_string()))?;

        let mut params = CreateBillingPortalSession::new(customer_id);
        params.return_url = Some(return_url);

        let session = BillingPortalSession::create(&self.client, params).await?;
        Ok(session)
    }

    /// Verify webhook signature and parse event
    pub fn verify_webhook(
        &self,
        payload: &str,
        signature: &str,
    ) -> Result<stripe::Event, BillingError> {
        let secret = self.webhook_secret.as_ref().ok_or_else(|| {
            BillingError::ConfigError("Webhook secret not configured".to_string())
        })?;

        let event = Webhook::construct_event(payload, signature, secret)
            .map_err(|_| BillingError::InvalidSignature)?;

        Ok(event)
    }

    /// Get customer by ID
    pub async fn get_customer(&self, customer_id: &str) -> Result<Customer, BillingError> {
        let id: CustomerId = customer_id
            .parse()
            .map_err(|_| BillingError::ConfigError("Invalid customer ID".to_string()))?;

        let customer = Customer::retrieve(&self.client, &id, &[]).await?;
        Ok(customer)
    }
}

// =============================================================================
// API Safeguards - Cost Control Levers
// =============================================================================

/// Maximum size of a single event payload (bytes). Events larger than this are rejected.
/// This prevents customers from sending huge JSON blobs that consume excessive CPU/storage.
pub const MAX_EVENT_SIZE_BYTES: usize = 8 * 1024; // 8 KB

/// Maximum total request body size (bytes). Requests larger than this are rejected.
pub const MAX_REQUEST_SIZE_BYTES: usize = 5 * 1024 * 1024; // 5 MB

/// Minimum events per batch request. Single-event requests are inefficient.
/// Requests with fewer events than this get a warning header.
pub const MIN_BATCH_SIZE: usize = 10;

/// Maximum events per single request. Prevents memory exhaustion.
pub const MAX_BATCH_SIZE: usize = 10_000;

/// Rate limit for demo endpoint (requests per hour per IP)
pub const DEMO_RATE_LIMIT_PER_HOUR: u32 = 10;

// =============================================================================
// Plan Definitions - Free/Pro/Team/Enterprise
// =============================================================================

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Plan {
    Free,       // Free tier - 10k events/mo, 5 streams, 14 days retention
    Pro,        // $99/mo - 500k events, 20 streams, 90 days retention
    Team,       // $199/mo - 5M events, 100 streams, 365 days retention
    Enterprise, // Custom - committed volume + overages, 500+ streams
}

impl Plan {
    pub fn event_limit(&self) -> u64 {
        match self {
            Plan::Free => 10_000,
            Plan::Pro => 500_000,
            Plan::Team => 5_000_000,
            Plan::Enterprise => u64::MAX,
        }
    }

    pub fn stream_limit(&self) -> u64 {
        match self {
            Plan::Free => 5,
            Plan::Pro => 20,
            Plan::Team => 100,
            Plan::Enterprise => 500,
        }
    }

    /// Data retention period in days
    pub fn retention_days(&self) -> u32 {
        match self {
            Plan::Free => 14,
            Plan::Pro => 90,
            Plan::Team => 365,
            Plan::Enterprise => 730, // 2 years default, negotiable
        }
    }

    /// Rate limit (requests per minute)
    pub fn rate_limit_rpm(&self) -> u32 {
        match self {
            Plan::Free => 60,       // 1/sec
            Plan::Pro => 300,       // 5/sec
            Plan::Team => 600,      // 10/sec
            Plan::Enterprise => 3000, // 50/sec
        }
    }

    /// Whether this plan is metered (has hard limits vs committed volume)
    /// Enterprise uses committed volume + overages, not hard limits
    pub fn is_metered(&self) -> bool {
        !matches!(self, Plan::Enterprise)
    }

    pub fn from_price_id(price_id: &str) -> Option<Self> {
        // Map Stripe price IDs to plans (supports legacy + new IDs)
        match price_id {
            // New canonical price IDs
            "price_pro" | "price_1Pro" => Some(Plan::Pro),
            "price_team" | "price_1Team" => Some(Plan::Team),
            "price_enterprise" | "price_1Enterprise" => Some(Plan::Enterprise),
            // Legacy price IDs (backwards compatibility)
            "price_radar" | "price_1Radar" => Some(Plan::Pro),
            "price_tensor" | "price_1Tensor" => Some(Plan::Team),
            "price_orbit" | "price_1Orbit" => Some(Plan::Enterprise),
            _ => None,
        }
    }

    /// Parse plan from string name (supports legacy aliases)
    pub fn from_name(name: &str) -> Option<Self> {
        match name.to_lowercase().as_str() {
            // Canonical names
            "free" => Some(Plan::Free),
            "pro" => Some(Plan::Pro),
            "team" => Some(Plan::Team),
            "enterprise" => Some(Plan::Enterprise),
            // Legacy aliases
            "trial" | "pilot" | "pulse" => Some(Plan::Free),
            "radar" | "signal" | "starter" => Some(Plan::Pro),
            "tensor" | "growth" | "lock" => Some(Plan::Team),
            "orbit" | "horizon" => Some(Plan::Enterprise),
            _ => None,
        }
    }

    pub fn name(&self) -> &'static str {
        match self {
            Plan::Free => "free",
            Plan::Pro => "pro",
            Plan::Team => "team",
            Plan::Enterprise => "enterprise",
        }
    }

    pub fn display_name(&self) -> &'static str {
        match self {
            Plan::Free => "Free",
            Plan::Pro => "Pro",
            Plan::Team => "Team",
            Plan::Enterprise => "Enterprise",
        }
    }

    pub fn monthly_price_cents(&self) -> u64 {
        match self {
            Plan::Free => 0,
            Plan::Pro => 9900,       // $99
            Plan::Team => 19900,     // $199
            Plan::Enterprise => 0,   // Custom pricing
        }
    }
}

/// Webhook event types we handle
pub enum WebhookEvent {
    CheckoutCompleted {
        customer_id: String,
        subscription_id: String,
        tenant_id: Option<String>,
    },
    SubscriptionUpdated {
        customer_id: String,
        subscription_id: String,
        status: String,
        price_id: Option<String>,
    },
    SubscriptionDeleted {
        customer_id: String,
        subscription_id: String,
    },
    InvoicePaid {
        customer_id: String,
        subscription_id: Option<String>,
        amount_paid: i64,
    },
    InvoicePaymentFailed {
        customer_id: String,
        subscription_id: Option<String>,
    },
    Unknown,
}

impl WebhookEvent {
    /// Parse a Stripe event into our internal representation
    pub fn from_stripe_event(event: &stripe::Event) -> Self {
        match event.type_ {
            stripe::EventType::CheckoutSessionCompleted => {
                if let stripe::EventObject::CheckoutSession(session) = &event.data.object {
                    return WebhookEvent::CheckoutCompleted {
                        customer_id: session
                            .customer
                            .as_ref()
                            .map(|c| c.id().to_string())
                            .unwrap_or_default(),
                        subscription_id: session
                            .subscription
                            .as_ref()
                            .map(|s| s.id().to_string())
                            .unwrap_or_default(),
                        tenant_id: session.client_reference_id.clone(),
                    };
                }
            }
            stripe::EventType::CustomerSubscriptionUpdated => {
                if let stripe::EventObject::Subscription(sub) = &event.data.object {
                    return WebhookEvent::SubscriptionUpdated {
                        customer_id: sub.customer.id().to_string(),
                        subscription_id: sub.id.to_string(),
                        status: sub.status.to_string(),
                        price_id: sub
                            .items
                            .data
                            .first()
                            .and_then(|item| item.price.as_ref())
                            .map(|p| p.id.to_string()),
                    };
                }
            }
            stripe::EventType::CustomerSubscriptionDeleted => {
                if let stripe::EventObject::Subscription(sub) = &event.data.object {
                    return WebhookEvent::SubscriptionDeleted {
                        customer_id: sub.customer.id().to_string(),
                        subscription_id: sub.id.to_string(),
                    };
                }
            }
            stripe::EventType::InvoicePaid => {
                if let stripe::EventObject::Invoice(invoice) = &event.data.object {
                    return WebhookEvent::InvoicePaid {
                        customer_id: invoice
                            .customer
                            .as_ref()
                            .map(|c| c.id().to_string())
                            .unwrap_or_default(),
                        subscription_id: invoice.subscription.as_ref().map(|s| s.id().to_string()),
                        amount_paid: invoice.amount_paid.unwrap_or(0),
                    };
                }
            }
            stripe::EventType::InvoicePaymentFailed => {
                if let stripe::EventObject::Invoice(invoice) = &event.data.object {
                    return WebhookEvent::InvoicePaymentFailed {
                        customer_id: invoice
                            .customer
                            .as_ref()
                            .map(|c| c.id().to_string())
                            .unwrap_or_default(),
                        subscription_id: invoice.subscription.as_ref().map(|s| s.id().to_string()),
                    };
                }
            }
            _ => {}
        }
        WebhookEvent::Unknown
    }
}
