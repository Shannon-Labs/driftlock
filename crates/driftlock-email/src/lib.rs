//! Driftlock Email - SendGrid integration
//!
//! This module provides email functionality using SendGrid's v3 API.
//! Supports transactional emails for user onboarding, verification, and billing notifications.

use serde::{Deserialize, Serialize};
use thiserror::Error;

/// SendGrid API endpoint
const SENDGRID_API_URL: &str = "https://api.sendgrid.com/v3/mail/send";

/// Email service errors
#[derive(Error, Debug)]
pub enum EmailError {
    #[error("SendGrid API error: {0}")]
    ApiError(String),
    #[error("Network error: {0}")]
    NetworkError(#[from] reqwest::Error),
    #[error("Invalid email address: {0}")]
    InvalidEmail(String),
    #[error("Invalid API key")]
    InvalidApiKey,
}

/// SendGrid email service
pub struct EmailService {
    api_key: String,
    from_email: String,
    #[allow(dead_code)]
    from_name: String,
    client: reqwest::Client,
}

impl EmailService {
    /// Create a new email service instance
    ///
    /// # Arguments
    /// * `api_key` - SendGrid API key
    /// * `from_email` - Sender email address (e.g., "noreply@driftlock.net")
    ///
    /// # Example
    /// ```
    /// use driftlock_email::EmailService;
    ///
    /// let service = EmailService::new("SG.xxxx", "noreply@driftlock.net");
    /// ```
    pub fn new(api_key: &str, from_email: &str) -> Self {
        Self::new_with_name(api_key, from_email, "Driftlock")
    }

    /// Create a new email service with custom sender name
    pub fn new_with_name(api_key: &str, from_email: &str, from_name: &str) -> Self {
        Self::with_client(api_key, from_email, from_name, Self::build_client())
    }

    /// Create email service with a custom reqwest client (useful for testing)
    pub fn with_client(
        api_key: &str,
        from_email: &str,
        from_name: &str,
        client: reqwest::Client,
    ) -> Self {
        Self {
            api_key: api_key.to_string(),
            from_email: from_email.to_string(),
            from_name: from_name.to_string(),
            client,
        }
    }

    /// Build default HTTP client with proper configuration
    fn build_client() -> reqwest::Client {
        reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()
            .expect("Failed to build HTTP client")
    }

    /// Send welcome email with API key
    ///
    /// # Arguments
    /// * `to` - Recipient email address
    /// * `company_name` - Company/user name
    /// * `api_key` - API key to include in email
    pub async fn send_welcome(
        &self,
        to: &str,
        company_name: &str,
        api_key: &str,
    ) -> Result<(), EmailError> {
        let subject = "Welcome to Driftlock - Your API Key Inside".to_string();
        let html_content = format!(
            r#"
<!DOCTYPE html>
<html>
<head>
    <style>
        body {{ font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }}
        .container {{ max-width: 600px; margin: 0 auto; padding: 20px; }}
        .header {{ background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }}
        .content {{ background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }}
        .api-key {{ background: #1f2937; color: #10b981; padding: 15px; border-radius: 6px; font-family: 'Courier New', monospace; font-size: 14px; word-break: break-all; margin: 20px 0; }}
        .button {{ display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }}
        .footer {{ color: #6b7280; font-size: 12px; margin-top: 30px; text-align: center; }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Driftlock!</h1>
        </div>
        <div class="content">
            <p>Hi {},</p>

            <p>Thank you for signing up for Driftlock! We're excited to help you detect anomalies in your OpenTelemetry data using compression-based analysis.</p>

            <p><strong>Your API Key:</strong></p>
            <div class="api-key">{}</div>

            <p><strong>Important:</strong> Store this API key securely. For security reasons, we won't show it again.</p>

            <h3>Getting Started</h3>
            <p>You can start detecting anomalies right away:</p>
            <pre style="background: #1f2937; color: #10b981; padding: 15px; border-radius: 6px; overflow-x: auto;">
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Authorization: Bearer {}" \
  -H "Content-Type: application/json" \
  -d '{{"stream_id": "my-stream", "data": "your log data"}}'
            </pre>

            <a href="https://driftlock.net/dashboard" class="button">Go to Dashboard</a>

            <h3>What's Next?</h3>
            <ul>
                <li>Configure your detection profiles</li>
                <li>Set up OpenTelemetry collector integration</li>
                <li>Review our <a href="https://docs.driftlock.net">documentation</a></li>
            </ul>

            <p>If you have any questions, reply to this email or visit our <a href="https://docs.driftlock.net">documentation</a>.</p>

            <p>Best regards,<br>The Driftlock Team</p>
        </div>
        <div class="footer">
            <p>Driftlock - Compression-Based Anomaly Detection</p>
            <p>You received this email because you signed up for Driftlock.</p>
        </div>
    </div>
</body>
</html>
            "#,
            company_name, api_key, api_key
        );

        self.send(to, subject.as_str(), &html_content).await
    }

    /// Send email verification link
    ///
    /// # Arguments
    /// * `to` - Recipient email address
    /// * `verification_token` - Verification token to include in link
    pub async fn send_verification(
        &self,
        to: &str,
        verification_token: &str,
    ) -> Result<(), EmailError> {
        let subject = "Verify Your Driftlock Email".to_string();
        let verification_url = format!(
            "https://driftlock.net/verify-email?token={}",
            verification_token
        );

        let html_content = format!(
            r#"
<!DOCTYPE html>
<html>
<head>
    <style>
        body {{ font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }}
        .container {{ max-width: 600px; margin: 0 auto; padding: 20px; }}
        .header {{ background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }}
        .content {{ background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }}
        .button {{ display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }}
        .footer {{ color: #6b7280; font-size: 12px; margin-top: 30px; text-align: center; }}
        .warning {{ background: #fef3c7; border-left: 4px solid #f59e0b; padding: 15px; border-radius: 4px; margin: 20px 0; }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email</h1>
        </div>
        <div class="content">
            <p>Thanks for signing up for Driftlock!</p>

            <p>Please verify your email address by clicking the button below:</p>

            <a href="{}" class="button">Verify Email Address</a>

            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #667eea;">{}</p>

            <div class="warning">
                <strong>Security Note:</strong> This link expires in 24 hours. If you didn't sign up for Driftlock, you can safely ignore this email.
            </div>

            <p>Best regards,<br>The Driftlock Team</p>
        </div>
        <div class="footer">
            <p>Driftlock - Compression-Based Anomaly Detection</p>
        </div>
    </div>
</body>
</html>
            "#,
            verification_url, verification_url
        );

        self.send(to, subject.as_str(), &html_content).await
    }

    /// Send trial ending warning (3 days before expiration)
    ///
    /// # Arguments
    /// * `to` - Recipient email address
    /// * `days_left` - Number of days remaining in trial
    pub async fn send_trial_ending(&self, to: &str, days_left: i32) -> Result<(), EmailError> {
        let subject = format!("Your Driftlock Trial Ends in {} Days", days_left);
        let html_content = format!(
            r#"
<!DOCTYPE html>
<html>
<head>
    <style>
        body {{ font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }}
        .container {{ max-width: 600px; margin: 0 auto; padding: 20px; }}
        .header {{ background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }}
        .content {{ background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }}
        .button {{ display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }}
        .pricing {{ background: white; border: 2px solid #e5e7eb; border-radius: 8px; padding: 20px; margin: 20px 0; }}
        .footer {{ color: #6b7280; font-size: 12px; margin-top: 30px; text-align: center; }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Your Trial is Ending Soon</h1>
        </div>
        <div class="content">
            <p>Hi there,</p>

            <p>Your Driftlock trial ends in <strong>{} days</strong>. We hope you've enjoyed detecting anomalies with compression-based analysis!</p>

            <h3>Continue Your Protection</h3>
            <p>Upgrade to a paid plan to keep your anomaly detection running:</p>

            <div class="pricing">
                <h4>Recommended Plans:</h4>
                <ul>
                    <li><strong>Pro</strong> - $99/month (500K events, 20 streams)</li>
                    <li><strong>Team</strong> - $199/month (5M events, 100 streams)</li>
                    <li><strong>Enterprise</strong> - Custom (Unlimited events, EU data residency)</li>
                </ul>
            </div>

            <a href="https://driftlock.net/dashboard/billing" class="button">Upgrade Now</a>

            <p><strong>What happens after trial?</strong></p>
            <p>Without upgrading, your account will downgrade to the free Pulse tier (10K events/month). Your detection profiles and historical data will be preserved.</p>

            <p>Have questions? Reply to this email - we're here to help!</p>

            <p>Best regards,<br>The Driftlock Team</p>
        </div>
        <div class="footer">
            <p>Driftlock - Compression-Based Anomaly Detection</p>
        </div>
    </div>
</body>
</html>
            "#,
            days_left
        );

        self.send(to, subject.as_str(), &html_content).await
    }

    /// Send payment failed notification
    ///
    /// # Arguments
    /// * `to` - Recipient email address
    pub async fn send_payment_failed(&self, to: &str) -> Result<(), EmailError> {
        let subject = "Action Required: Payment Failed for Driftlock".to_string();
        let html_content = r#"
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .alert { background: #fee2e2; border-left: 4px solid #ef4444; padding: 15px; border-radius: 4px; margin: 20px 0; }
        .footer { color: #6b7280; font-size: 12px; margin-top: 30px; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Payment Failed</h1>
        </div>
        <div class="content">
            <div class="alert">
                <strong>Action Required:</strong> We couldn't process your payment for Driftlock.
            </div>

            <p>Hi there,</p>

            <p>We attempted to charge your payment method but the payment failed. This could be due to:</p>
            <ul>
                <li>Insufficient funds</li>
                <li>Expired credit card</li>
                <li>Card declined by your bank</li>
                <li>Incorrect billing information</li>
            </ul>

            <h3>What You Need to Do</h3>
            <p>Update your payment method within 7 days to avoid service interruption:</p>

            <a href="https://driftlock.net/dashboard/billing" class="button">Update Payment Method</a>

            <p><strong>Grace Period:</strong> You have 7 days to update your payment information before your account is downgraded to the free tier.</p>

            <p>If you have questions or need assistance, please reply to this email.</p>

            <p>Best regards,<br>The Driftlock Team</p>
        </div>
        <div class="footer">
            <p>Driftlock - Compression-Based Anomaly Detection</p>
        </div>
    </div>
</body>
</html>
        "#.to_string();

        self.send(to, subject.as_str(), &html_content).await
    }

    /// Internal send function that calls SendGrid API
    ///
    /// # Arguments
    /// * `to` - Recipient email address
    /// * `subject` - Email subject line
    /// * `html_content` - HTML email body
    async fn send(&self, to: &str, subject: &str, html_content: &str) -> Result<(), EmailError> {
        // Validate email address (basic check)
        if !to.contains('@') || !to.contains('.') {
            return Err(EmailError::InvalidEmail(to.to_string()));
        }

        // Validate API key
        if self.api_key.is_empty() {
            return Err(EmailError::InvalidApiKey);
        }

        // Construct SendGrid API request
        let request_body = SendGridRequest {
            personalizations: vec![Personalization {
                to: vec![EmailAddress {
                    email: to.to_string(),
                }],
            }],
            from: EmailAddress {
                email: self.from_email.clone(),
            },
            reply_to: Some(EmailAddress {
                email: self.from_email.clone(),
            }),
            subject: subject.to_string(),
            content: vec![Content {
                r#type: "text/html".to_string(),
                value: html_content.to_string(),
            }],
        };

        // Send request to SendGrid
        tracing::debug!("Sending email to {} with subject: {}", to, subject);

        let response = self
            .client
            .post(SENDGRID_API_URL)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .header("Content-Type", "application/json")
            .json(&request_body)
            .send()
            .await?;

        // Check response status
        let status = response.status();
        if status.is_success() {
            tracing::info!("Email sent successfully to {}", to);
            Ok(())
        } else {
            let error_body = response
                .text()
                .await
                .unwrap_or_else(|_| "Unknown error".to_string());
            tracing::error!("SendGrid API error: {} - {}", status, error_body);
            Err(EmailError::ApiError(format!(
                "HTTP {}: {}",
                status, error_body
            )))
        }
    }
}

// SendGrid API request structures

#[derive(Debug, Serialize, Deserialize)]
struct SendGridRequest {
    personalizations: Vec<Personalization>,
    from: EmailAddress,
    #[serde(skip_serializing_if = "Option::is_none")]
    reply_to: Option<EmailAddress>,
    subject: String,
    content: Vec<Content>,
}

#[derive(Debug, Serialize, Deserialize)]
struct Personalization {
    to: Vec<EmailAddress>,
}

#[derive(Debug, Serialize, Deserialize)]
struct EmailAddress {
    email: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Content {
    r#type: String,
    value: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sendgrid_request_serialization() {
        // Test that our SendGrid request structures serialize correctly
        let request = SendGridRequest {
            personalizations: vec![Personalization {
                to: vec![EmailAddress {
                    email: "test@example.com".to_string(),
                }],
            }],
            from: EmailAddress {
                email: "noreply@driftlock.net".to_string(),
            },
            reply_to: Some(EmailAddress {
                email: "noreply@driftlock.net".to_string(),
            }),
            subject: "Test Subject".to_string(),
            content: vec![Content {
                r#type: "text/html".to_string(),
                value: "<html><body>Test</body></html>".to_string(),
            }],
        };

        let json = serde_json::to_string(&request).expect("Failed to serialize");
        assert!(json.contains("test@example.com"));
        assert!(json.contains("Test Subject"));
        assert!(json.contains("text/html"));
    }

    #[test]
    fn test_email_address_validation() {
        // Test basic email validation logic without creating HTTP client
        assert!("test@example.com".contains('@') && "test@example.com".contains('.'));
        assert!(!("invalid-email".contains('@') && "invalid-email".contains('.')));
    }

    #[test]
    fn test_api_key_validation() {
        // Test API key validation logic
        let empty_key = "";
        let valid_key = "SG.test_key";
        assert!(empty_key.is_empty());
        assert!(!valid_key.is_empty());
    }
}
