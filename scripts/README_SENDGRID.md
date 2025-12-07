# SendGrid CLI Tool

We've created a Python-based SendGrid CLI tool to help with email configuration and testing.

## Installation

The tool is already set up with the required Python SDK in `/tmp/sendgrid-env`.

## Usage

```bash
# List all verified senders
./scripts/sendgrid --list-senders

# Check if a specific sender is verified
./scripts/sendgrid --check-sender noreply@driftlock.net

# Show current configuration and solutions
./scripts/sendgrid --current-config

# Send a test email (replace with your email)
./scripts/sendgrid --test-email your-email@example.com
```

## Current Status

- ✅ **API Key**: Configured and working
- ✅ **Verified Sender**: `hunter@shannonlabs.dev` is verified
- ❌ **Problem**: Application tries to send from `noreply@driftlock.net` (not verified)

## Solutions

### Option 1: Quick Fix (Use existing sender)
Update the environment variable in Cloud Run:
```
EMAIL_FROM_ADDRESS=hunter@shannonlabs.dev
EMAIL_FROM_NAME=Driftlock
```

### Option 2: Verify noreply@driftlock.net
1. Go to: https://app.sendgrid.com/settings/sender_auth
2. Click "Verify a Single Sender"
3. Enter:
   - From Email: noreply@driftlock.net
   - From Name: Driftlock
   - Reply To: support@driftlock.net
4. Check the email inbox for verification link

### Option 3: Domain Verification (Recommended for Production)
1. Go to: https://app.sendgrid.com/settings/sender_auth
2. Click "Authenticate Your Domain"
3. Enter: driftlock.net
4. Add the DNS records provided by SendGrid

## Testing

After fixing the sender verification, test the email flow:

```bash
# Test the signup flow
./scripts/test-email.sh your-email@example.com

# Or send a direct test email
./scripts/sendgrid --test-email your-email@example.com
```

## API Commands

You can also use curl directly with the SendGrid API:

```bash
# Get API key from GCP Secret Manager
SG_KEY=$(gcloud secrets versions access latest --secret=sendgrid-api-key --project=driftlock)

# List verified senders
curl -s "https://api.sendgrid.com/v3/verified_senders" \
  -H "Authorization: Bearer $SG_KEY" \
  -H "Content-Type: application/json"

# Send a test email
curl -X POST https://api.sendgrid.com/v3/mail/send \
  -H "Authorization: Bearer $SG_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "personalizations": [{"to": [{"email": "your-email@example.com"}]}],
    "from": {"email": "hunter@shannonlabs.dev"},
    "subject": "Driftlock Test",
    "content": [{"type": "text/plain", "value": "Test email"}]
  }'
```