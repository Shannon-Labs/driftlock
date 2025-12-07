#!/usr/bin/env python3
"""
SendGrid helper script for checking configuration and testing emails
Usage: ./sendgrid_helper.py [--test-email your@email.com] [--list-senders] [--verify-sender noreply@driftlock.net]
"""

import os
import sys
import argparse
import json
from sendgrid import SendGridAPIClient
from sendgrid.helpers.mail import Mail

def get_sendgrid_client():
    """Get SendGrid client with API key from GCP Secret Manager"""
    try:
        # Try to get from environment first
        api_key = os.getenv('SENDGRID_API_KEY')
        if not api_key:
            # Fallback: get from gcloud
            import subprocess
            result = subprocess.run(
                ['gcloud', 'secrets', 'versions', 'access', 'latest',
                 '--secret=sendgrid-api-key', '--project=driftlock'],
                capture_output=True, text=True
            )
            if result.returncode == 0:
                api_key = result.stdout.strip()
            else:
                print(f"❌ Error getting SendGrid API key: {result.stderr}")
                sys.exit(1)

        return SendGridAPIClient(api_key=api_key)
    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)

def list_verified_senders(client):
    """List all verified senders"""
    try:
        response = client.client.senders.get()
        senders = response.to_dict

        print("\n=== Verified Senders ===")
        if isinstance(senders, dict) and senders.get('results'):
            for sender in senders['results']:
                status = "✅ Verified" if sender.get('verified') else "❌ Not verified"
                print(f"• {sender.get('from_email')} ({sender.get('nickname', 'No nickname')}) - {status}")
                print(f"  From: {sender.get('from_name', 'N/A')}")
                print(f"  Reply-to: {sender.get('reply_to', 'N/A')}")
                print(f"  ID: {sender.get('id')}")
                print()
        elif isinstance(senders, list):
            for sender in senders:
                status = "✅ Verified" if sender.get('verified') else "❌ Not verified"
                print(f"• {sender.get('from_email')} ({sender.get('nickname', 'No nickname')}) - {status}")
                print(f"  From: {sender.get('from_name', 'N/A')}")
                print(f"  Reply-to: {sender.get('reply_to', 'N/A')}")
                print(f"  ID: {sender.get('id')}")
                print()
        else:
            print("No verified senders found")
            print(f"Response type: {type(senders)}")
            if hasattr(senders, '__dict__'):
                print(f"Response keys: {senders.__dict__.keys()}")

    except Exception as e:
        print(f"❌ Error listing senders: {e}")

def test_email_send(client, to_email):
    """Test sending an email"""
    try:
        # Use the currently verified sender
        from_email = "hunter@shannonlabs.dev"
        from_name = "Driftlock Test"

        message = Mail(
            from_email=(from_email, from_name),
            to_emails=to_email,
            subject="Driftlock Email Test ✅",
            html_content="""<strong>Driftlock Email Test</strong><br><br>
            This is a test email from the Driftlock system.<br><br>
            If you received this, SendGrid is configured correctly!<br><br>
            Best regards,<br>
            The Driftlock Team""",
            plain_text_content="""Driftlock Email Test

This is a test email from the Driftlock system.

If you received this, SendGrid is configured correctly!

Best regards,
The Driftlock Team"""
        )

        print(f"\n📧 Sending test email to {to_email}...")
        response = client.send(message)

        if response.status_code == 202:
            print(f"✅ Email sent successfully!")
            print(f"   Message ID: {response.headers.get('X-Message-Id', 'N/A')}")
        else:
            print(f"❌ Failed to send email: {response.status_code}")
            print(f"   Response: {response.body}")

    except Exception as e:
        print(f"❌ Error sending email: {e}")

def check_single_sender(client, email):
    """Check if a specific sender is verified"""
    try:
        response = client.client.senders.get()
        senders = response.to_dict

        # Handle different response formats
        if isinstance(senders, dict):
            senders_list = senders.get('results', [])
        elif isinstance(senders, list):
            senders_list = senders
        else:
            print(f"❌ Unexpected response format: {type(senders)}")
            return False

        for sender in senders_list:
            if sender.get('from_email', '').lower() == email.lower():
                if sender.get('verified'):
                    print(f"✅ {email} is verified!")
                    return True
                else:
                    print(f"❌ {email} exists but is not verified")
                    print(f"   Status: {sender.get('verification_status', 'Unknown')}")
                    return False

        print(f"❌ {email} not found in verified senders list")
        return False

    except Exception as e:
        print(f"❌ Error checking sender: {e}")
        return False

def main():
    parser = argparse.ArgumentParser(description='SendGrid helper tool')
    parser.add_argument('--test-email', help='Send a test email to this address')
    parser.add_argument('--list-senders', action='store_true', help='List all verified senders')
    parser.add_argument('--check-sender', help='Check if a specific sender is verified')
    parser.add_argument('--current-config', action='store_true', help='Show current configuration')

    args = parser.parse_args()

    print("=== SendGrid Helper Tool ===")

    if not any([args.test_email, args.list_senders, args.check_sender, args.current_config]):
        parser.print_help()
        sys.exit(1)

    # Initialize SendGrid client
    client = get_sendgrid_client()
    print("✅ Connected to SendGrid API")

    # No need to fix typo - just continue

    # Run requested actions
    if args.list_senders:
        list_verified_senders(client)

    if args.check_sender:
        check_single_sender(client, args.check_sender)

    if args.current_config:
        print("\n=== Current Configuration ===")
        print("Application sends from: noreply@driftlock.net")
        print("Issue: This sender is not verified in SendGrid")
        print("\nSolutions:")
        print("1. Update EMAIL_FROM_ADDRESS to use hunter@shannonlabs.dev")
        print("2. Verify noreply@driftlock.net in SendGrid dashboard")
        print("3. Verify driftlock.net domain in SendGrid dashboard")

    if args.test_email:
        test_email_send(client, args.test_email)

if __name__ == '__main__':
    main()