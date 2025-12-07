#!/bin/bash
# Check SendGrid configuration and help with setup

set -e

echo "=== SendGrid Configuration Check ==="
echo ""

# Get SendGrid API key
echo "1. Getting SendGrid API key..."
SG_KEY=$(gcloud secrets versions access latest --secret=sendgrid-api-key --project=driftlock 2>/dev/null)

if [ -z "$SG_KEY" ]; then
    echo "❌ SendGrid API key not found in Secret Manager"
    echo "   Expected secret: sendgrid-api-key"
    exit 1
fi

echo "✅ SendGrid API key found (length: ${#SG_KEY})"
echo ""

# Check verified senders
echo "2. Checking verified senders..."
echo ""

RESPONSE=$(curl -s "https://api.sendgrid.com/v3/verified_senders" \
  -H "Authorization: Bearer $SG_KEY" \
  -H "Content-Type: application/json")

echo "$RESPONSE" | jq -r '.results[] | "• \(.from_email) (\(.nickname)) - Verified: \(.verified)"' 2>/dev/null || {
    echo "Could not parse response. Raw output:"
    echo "$RESPONSE"
}

echo ""
echo "3. Current email configuration in code..."
echo "   From: collector-processor/cmd/driftlock-http/email.go:27"
echo "   Address: noreply@driftlock.net"
echo ""

# Check if noreply@driftlock.net is verified
if echo "$RESPONSE" | grep -q "noreply@driftlock.net"; then
    echo "✅ noreply@driftlock.net is already verified!"
else
    echo "❌ noreply@driftlock.net is NOT verified"
    echo ""
    echo "=== SOLUTIONS ==="
    echo ""
    echo "Option A: Use existing verified sender"
    echo "  Update EMAIL_FROM_ADDRESS environment variable to use: hunter@shannonlabs.dev"
    echo ""
    echo "Option B: Verify noreply@driftlock.net"
    echo "  1. Go to: https://app.sendgrid.com/settings/sender_auth"
    echo "  2. Click 'Verify a Single Sender'"
    echo "  3. Enter:"
    echo "     - From Email: noreply@driftlock.net"
    echo "     - From Name: Driftlock"
    echo "     - Reply To: support@driftlock.net (or noreply@driftlock.net)"
    echo "  4. Complete verification process"
    echo ""
    echo "Option C: Verify the entire domain (recommended for production)"
    echo "  1. Go to: https://app.sendgrid.com/settings/sender_auth"
    echo "  2. Click 'Authenticate Your Domain'"
    echo "  3. Enter: driftlock.net"
    echo "  4. Add DNS records as provided by SendGrid"
    echo ""
fi

echo "4. Testing email delivery..."
echo "   After fixing the sender verification, run:"
echo "   ./scripts/test-email.sh your-email@example.com"
echo ""

echo "=== Quick Test ==="
echo "To test if SendGrid is working with the current sender:"
echo ""

# Create test with existing verified sender
cat > /tmp/test_email_config.json <<EOF
{
  "from": "hunter@shannonlabs.dev",
  "to": "test@example.com",
  "subject": "Driftlock Test",
  "content": "This is a test email from Driftlock."
}
EOF

echo "Test command (replace with your email):"
echo 'curl -X POST https://api.sendgrid.com/v3/mail/send \\'
echo '  -H "Authorization: Bearer $SG_KEY" \\'
echo '  -H "Content-Type: application/json" \\'
echo '  -d "{'
echo '    \"personalizations\": [{\"to\": [{\"email\": \"your-email@example.com\"}]}],'
echo '    \"from\": {\"email\": \"hunter@shannonlabs.dev\"},'
echo '    \"subject\": \"Driftlock Test\",'
echo '    \"content\": [{\"type\": \"text/plain\", \"value\": \"Test email\"}]'
echo '  }"'