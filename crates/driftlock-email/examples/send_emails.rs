//! Example: Sending emails with Driftlock Email Service
//!
//! Usage:
//!   cargo run --example send_emails
//!
//! Environment variables:
//!   SENDGRID_API_KEY - Your SendGrid API key
//!   TEST_EMAIL       - Email address to send test emails to

use driftlock_email::{EmailError, EmailService};
use std::env;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load environment variables
    let sendgrid_api_key =
        env::var("SENDGRID_API_KEY").expect("SENDGRID_API_KEY environment variable must be set");

    let test_email = env::var("TEST_EMAIL").unwrap_or_else(|_| "test@example.com".to_string());

    println!("Initializing Driftlock Email Service...");

    // Create email service
    let email_service = EmailService::new(&sendgrid_api_key, "noreply@driftlock.net");

    println!("\nSending test emails to: {}\n", test_email);

    // Test 1: Send welcome email
    println!("1. Sending welcome email...");
    match email_service
        .send_welcome(&test_email, "Acme Corporation", "dk_live_1234567890abcdef")
        .await
    {
        Ok(()) => println!("   ✓ Welcome email sent successfully"),
        Err(e) => eprintln!("   ✗ Error sending welcome email: {}", e),
    }

    // Test 2: Send verification email
    println!("\n2. Sending verification email...");
    match email_service
        .send_verification(&test_email, "verify_abc123def456")
        .await
    {
        Ok(()) => println!("   ✓ Verification email sent successfully"),
        Err(e) => eprintln!("   ✗ Error sending verification email: {}", e),
    }

    // Test 3: Send trial ending email
    println!("\n3. Sending trial ending notification...");
    match email_service.send_trial_ending(&test_email, 3).await {
        Ok(()) => println!("   ✓ Trial ending email sent successfully"),
        Err(e) => eprintln!("   ✗ Error sending trial ending email: {}", e),
    }

    // Test 4: Send payment failed email
    println!("\n4. Sending payment failed notification...");
    match email_service.send_payment_failed(&test_email).await {
        Ok(()) => println!("   ✓ Payment failed email sent successfully"),
        Err(e) => eprintln!("   ✗ Error sending payment failed email: {}", e),
    }

    // Test 5: Test email validation
    println!("\n5. Testing email validation...");
    match email_service
        .send_welcome("invalid-email-address", "Test", "test_key")
        .await
    {
        Err(EmailError::InvalidEmail(email)) => {
            println!("   ✓ Correctly rejected invalid email: {}", email);
        }
        Ok(()) => eprintln!("   ✗ Should have rejected invalid email"),
        Err(e) => eprintln!("   ✗ Unexpected error: {}", e),
    }

    // Test 6: Test API key validation
    println!("\n6. Testing API key validation...");
    let empty_service = EmailService::new("", "noreply@driftlock.net");
    match empty_service.send_welcome(&test_email, "Test", "key").await {
        Err(EmailError::InvalidApiKey) => {
            println!("   ✓ Correctly rejected empty API key");
        }
        Ok(()) => eprintln!("   ✗ Should have rejected empty API key"),
        Err(e) => eprintln!("   ✗ Unexpected error: {}", e),
    }

    println!("\n✓ All tests completed!");
    println!("\nNote: Check your inbox at {} for test emails", test_email);

    Ok(())
}
