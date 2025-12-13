# Driftlock Email Service

SendGrid v3 API integration for Driftlock transactional emails.

## Features

- Welcome emails with API keys
- Email verification links
- Trial ending notifications
- Payment failure alerts
- Full HTML email templates
- Proper error handling
- Async/await support

## Usage

```rust
use driftlock_email::{EmailService, EmailError};

#[tokio::main]
async fn main() -> Result<(), EmailError> {
    // Initialize service
    let email_service = EmailService::new(
        "SG.your_sendgrid_api_key",
        "noreply@driftlock.net"
    );

    // Send welcome email
    email_service
        .send_welcome(
            "user@example.com",
            "Acme Corp",
            "dk_live_abc123xyz"
        )
        .await?;

    // Send verification email
    email_service
        .send_verification(
            "user@example.com",
            "verification_token_xyz"
        )
        .await?;

    // Send trial ending warning
    email_service
        .send_trial_ending("user@example.com", 3)
        .await?;

    // Send payment failed notification
    email_service
        .send_payment_failed("user@example.com")
        .await?;

    Ok(())
}
```

## Environment Variables

```bash
SENDGRID_API_KEY=SG.your_api_key_here
FROM_EMAIL=noreply@driftlock.net
```

## Email Templates

### Welcome Email
- Includes API key
- Getting started instructions
- Dashboard link
- Usage example

### Verification Email
- Verification link (24-hour expiry)
- Security notice
- Manual link copy option

### Trial Ending
- Days remaining
- Pricing tiers
- Upgrade CTA
- Free tier fallback info

### Payment Failed
- Action required alert
- Common failure reasons
- Grace period (7 days)
- Update payment CTA

## SendGrid API

All emails use SendGrid's `/v3/mail/send` endpoint with:
- HTML content type
- Reply-to address
- Professional templates
- Responsive design

## Error Handling

```rust
use driftlock_email::EmailError;

match email_service.send_welcome("user@example.com", "Company", "key").await {
    Ok(()) => println!("Email sent successfully"),
    Err(EmailError::ApiError(msg)) => eprintln!("SendGrid error: {}", msg),
    Err(EmailError::NetworkError(e)) => eprintln!("Network error: {}", e),
    Err(EmailError::InvalidEmail(email)) => eprintln!("Invalid email: {}", email),
    Err(EmailError::InvalidApiKey) => eprintln!("Missing or invalid API key"),
}
```

## Testing

```bash
# Run tests
cargo test -p driftlock-email

# Build
cargo build -p driftlock-email
```

## Integration

From other crates:

```toml
[dependencies]
driftlock-email = { path = "../driftlock-email" }
```

```rust
use driftlock_email::EmailService;

let email = EmailService::new(&config.sendgrid_api_key, "noreply@driftlock.net");
email.send_welcome(&user_email, &company_name, &api_key).await?;
```
