//! Example: Verify a Firebase ID token
//!
//! This example demonstrates how to verify a Firebase ID token.
//!
//! Usage:
//!   cargo run --example verify_token -- <project_id> <token>
//!
//! Example:
//!   cargo run --example verify_token -- driftlock eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

use driftlock_auth::{FirebaseAuth, FirebaseError};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = std::env::args().collect();

    if args.len() != 3 {
        eprintln!("Usage: {} <project_id> <token>", args[0]);
        eprintln!("\nExample:");
        eprintln!(
            "  {} driftlock eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
            args[0]
        );
        std::process::exit(1);
    }

    let project_id = &args[1];
    let token = &args[2];

    println!("Verifying Firebase token...");
    println!("Project ID: {}", project_id);
    println!();

    // Create Firebase authenticator
    let firebase = FirebaseAuth::new(project_id);

    // Verify the token
    match firebase.verify_token(token).await {
        Ok(user) => {
            println!("Token verified successfully!");
            println!();
            println!("User ID: {}", user.uid);
            println!(
                "Email: {}",
                user.email.unwrap_or_else(|| "<none>".to_string())
            );
            println!("Email verified: {}", user.email_verified);
        }
        Err(FirebaseError::TokenExpired) => {
            eprintln!("Error: Token has expired");
            std::process::exit(1);
        }
        Err(FirebaseError::InvalidIssuer) => {
            eprintln!("Error: Token was not issued by Firebase");
            std::process::exit(1);
        }
        Err(FirebaseError::InvalidAudience) => {
            eprintln!("Error: Token was issued for a different project");
            std::process::exit(1);
        }
        Err(e) => {
            eprintln!("Error: Failed to verify token: {}", e);
            std::process::exit(1);
        }
    }

    Ok(())
}
