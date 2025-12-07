# Full End-to-End Testing Guide

## Prerequisites
- Email access (Gmail, Outlook, etc.)
- Web browser
- Command line access

## Option 1: Quick Manual Test

### User Signup Flow
1. **Go to**: https://driftlock.net
2. **Click**: Sign Up / Get Started
3. **Enter**:
   - Email: `your-email@example.com`
   - Company: `Test Company`
   - Password: `any-password`
   - Plan: Select "Pilot" (free trial)
4. **Submit**: Create account
5. **Check Email**: Look for "Welcome to Driftlock!" email
6. **Click Verification Link**: Will open in browser
7. **Copy API Key**: Shown in success message
8. **Test API**:
   ```bash
   curl -X POST https://driftlock.net/api/v1/detect \
     -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"events": ["test event 1", "test event 2"]}'
   ```

### Stripe Checkout Flow (after signup)
1. **Login**: Go to https://driftlock.net
2. **Click**: Upgrade button
3. **Select Plan**: Radar ($100/mo) or Tensor ($299/mo)
4. **Pay with Stripe Test Card**: `4242 4242 4242 4242`
   - Expiry: Any future date
   - CVC: Any 3 digits
   - ZIP: Any 5 digits
5. **Success**: Redirected back to dashboard

## Option 2: Automated Testing Scripts

### Test Signup Only
```bash
# Replace with your email
./scripts/test-full-signup.sh your-email@example.com
```

### Test Complete Flow (Signup → Checkout)
```bash
# Replace with your email and desired plan
./scripts/test-full-checkout.sh your-email@example.com radar
```

## What to Verify

### ✅ Successful Signup Flow
- [ ] Account created successfully
- [ ] Verification email received within 1 minute
- [ ] Verification link works (no 404 errors)
- [ ] API key generated and displayed
- [ ] API key works with detect endpoint
- [ ] Dashboard loads with tenant data

### ✅ Successful Checkout Flow
- [ ] Signup flow complete
- [ ] Checkout session created without errors
- [ ] Stripe checkout page loads
- [ ] Payment with test card succeeds
- [ ] Redirect back to application
- [ ] Subscription updated in database
- - Paid plan API access working
- [ ] Confirmation email received

## Troubleshooting

### Common Issues
1. **Email not received**
   - Check spam/junk folder
   - Try `hunter@shannonlabs.dev` (already verified)

2. **Verification link expires**
   - Links expire after 24 hours
   - Request a new verification email via API

3. **API key doesn't work**
   - Ensure email is verified first
   - Check API key format: `dlk_<uuid>.<secret>`

4. **Stripe checkout fails**
   - Ensure email is verified
   - Check Stripe is in test mode
   - Use test card details

### Checking Status
```bash
# Check tenant status
curl -s "https://driftlock.net/api/v1/onboard/status?tenant_id=<ID>" \
  -H "Content-Type: application/json"

# Check billing status
curl -s "https://driftlock.net/api/v1/me/billing" \
  -H "Authorization: Bearer <API_KEY>" \
  -H "Content-Type: application/json"
```

## Test Data

### For Anomaly Detection
```json
{
  "events": [
    {"timestamp": "2025-01-06T10:00:00Z", "level": "INFO", "message": "User login"},
    {"timestamp": "2025-01-06T10:01:00Z", "level": "INFO", "message": "Normal operation"},
    {"timestamp": "2025-01-06T10:02:00Z", "level": "ERROR", "message": "Database connection failed!"},
    {"timestamp": "2025-01-06T10:03:00Z", "level": "ERROR", "message": "Database connection failed!"}
  ]
}
```

## Pricing Tiers
- **Pilot**: Free trial (10K events/day)
- **Radar**: $15/month (250K events/month) ← Good for testing
- **Tensor**: $100/month (2.5M events/month)
- **Orbit**: $299/month (10M events/month)