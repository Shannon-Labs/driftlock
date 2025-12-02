# Driftlock Launch Plan 🚀

## Status Overview

**Current Issues:**
- Firebase authentication expired (needs re-auth)
- Cloud Run backend returning 404 (need to verify deployment)
- Stripe webhook secret not configured
- Cloudflare domain not connected to Firebase Hosting

**Goal:** Full launch today

---

## 1. Stripe Webhook Setup (Priority 1)

### Stripe CLI Installation & Setup

**If you don't have Stripe CLI installed:**
```bash
# macOS (using Homebrew)
brew install stripe/stripe-cli/stripe

# Or download directly from https://stripe.com/docs/stripe-cli
```

**Configure Stripe CLI for local testing:**
```bash
# Login to Stripe (opens browser)
stripe login

# Or use API key directly
stripe login --api-key sk_test_your_key_here
```

### Create Stripe Webhook Listener

**For local development testing:**
```bash
# Terminal 1: Start Firebase Functions emulator
cd functions
npm run build
firebase emulators:start --only functions

# Terminal 2: Forward Stripe webhooks to local
cd ..  # Back to project root
./scripts/stripe-webhook-forward.sh
```

**For production testing:**
```bash
# Listen to Stripe events and forward to production
stripe listen --forward-to https://driftlock.net/webhooks/stripe

# In another terminal, trigger test events
stripe trigger customer.subscription.created
stripe trigger customer.subscription.deleted
stripe trigger invoice.payment_succeeded
stripe trigger invoice.payment_failed
```

### Configure Webhook in Stripe Dashboard

1. Go to https://dashboard.stripe.com/webhooks
2. Click "Add endpoint"
3. URL: `https://driftlock.net/webhooks/stripe`
4. Events to listen for:
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_succeeded`
   - `invoice.payment_failed`
   - `checkout.session.completed`
5. Copy the **Signing secret** (starts with `whsec_`)

### Store Webhook Secret in Firebase/Google Cloud

```bash
# Store in Firebase Functions environment
cd functions
firebase functions:config:set stripe.webhook_secret="whsec_your_secret_here"

# Or set in Google Cloud Secret Manager
echo -n "whsec_your_secret_here" | gcloud secrets create stripe-webhook-secret --data-file=- --project=driftlock

# Or update .env for local testing
echo "STRIPE_WEBHOOK_SECRET=whsec_your_secret_here" >> functions/.env
```

---

## 2. Cloudflare + Firebase Domain Configuration

### Option A: Move to Google Domains (Recommended)

**Pros:**
- Simpler management (all Google ecosystem)
- Better Firebase integration
- Automatic SSL/CDN

**Steps:**
1. Transfer driftlock.net to Google Domains
2. In Firebase Console > Hosting > Connect Domain:
   - Enter: `driftlock.net`
   - Follow verification steps
3. Update DNS to point to Firebase:
   ```
   A @ <firebase_ip_1>
   A @ <firebase_ip_2>
   ```

### Option B: Keep Cloudflare + Point to Firebase

**Current state**: driftlock.net is on Cloudflare

**Steps:**
1. In Firebase Console > Hosting:
   - Add custom domain: `driftlock.net`
   - Note the 2 A records provided

2. In Cloudflare Dashboard:
   - Go to DNS > Records
   - Add the 2 A records from Firebase
   - Set Proxy status to "DNS only" (gray cloud, NOT orange)

3. Wait for SSL certificate (Firebase handles automatically)

4. Test:
   ```bash
   curl -I https://driftlock.net
   # Should return 200 OK
   ```

---

## 3. Re-authenticate & Deploy

### Re-authenticate Firebase

```bash
firebase login --reauth

# Or generate CI token
firebase login:ci
# Save token securely
```

### Re-authenticate Google Cloud

```bash
gcloud auth login
gcloud config set project driftlock
gcloud auth configure-docker
```

### Verify Cloud Run Backend Status

```bash
gcloud run services describe driftlock-api --region=us-central1 --format="yaml(status.address.url)"

# Should show URL: https://driftlock-api-o6kjgrsowq-uc.a.run.app
```

### Deploy Backend if Needed

```bash
# If backend is down or needs update
cd collector-processor/cmd/driftlock-http
gcloud run deploy driftlock-api \
  --source=. \
  --region=us-central1 \
  --allow-unauthenticated \
  --set-env-vars="STRIPE_SECRET_KEY=sk_test_xxx,STRIPE_WEBHOOK_SECRET=whsec_yyy"
```

### Deploy Firebase Functions

```bash
cd functions
npm install
npm run build
firebase deploy --only functions

# Check function logs
firebase functions:log
```

### Deploy Firebase Hosting

```bash
# Build landing page
cd landing-page
npm install
npm run build
cd ..

# Deploy hosting
firebase deploy --only hosting
```

---

## 4. Comprehensive Testing Checklist

### Backend API Tests

```bash
# Test health endpoint
curl -s https://driftlock-api-o6kjgrsowq-uc.a.run.app/healthz | jq

# Expected: {"success":true,"database":"healthy","license":"valid","version":"2.0.0"}

# Test signup endpoint
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@driftlock.dev","company_name":"Test Co"}' | jq

# Test anomaly detection
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{"stream_id":"test","events":[{"data":"test event"}]}' | jq
```

### Firebase Functions Tests

```bash
# Test health check via Firebase Functions
curl -s https://us-central1-driftlock.cloudfunctions.net/healthCheck | jq

# Test proxy to backend
curl -s https://us-central1-driftlock.cloudfunctions.net/apiProxy/api/v1/healthz | jq

# Test signup via Firebase
curl -X POST https://us-central1-driftlock.cloudfunctions.net/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@driftlock.dev","company_name":"Test Co"}' | jq
```

### Stripe Webhook Tests

```bash
# In one terminal: start Stripe listener
stripe listen --forward-to https://driftlock.net/webhooks/stripe

# In another terminal: trigger test events
stripe trigger checkout.session.completed
stripe trigger customer.subscription.created

# Verify webhook in Stripe Dashboard shows 200 OK
```

### Frontend Tests

```bash
# Test landing page
curl -I https://driftlock.net

# Test dashboard (if authenticated)
curl -s https://driftlock.net/dashboard | head -20

# Test API endpoints via frontend
# Open browser console and check network requests
```

### End-to-End User Flow Test

1. **User Signup:**
   - Visit https://driftlock.net
   - Click "Get Started"
   - Enter email: test@driftlock.dev
   - Should receive API key in ~30 seconds

2. **Checkout Flow:**
   - Visit billing page
   - Click "Upgrade to Pro"
   - Should redirect to Stripe Checkout
   - Use test card: 4242 4242 4242 4242
   - Should redirect back to dashboard

3. **API Usage:**
   - Use API key to call `/v1/anomalies`
   - Should return 200 OK with results

4. **Webhook Processing:**
   - Check Stripe Dashboard for successful webhooks
   - Check Firebase Functions logs
   - Verify tenant plan updated in database

---

## 5. Monitoring & Verification

### Check Logs

```bash
# Firebase Functions logs
firebase functions:log --limit 50

# Cloud Run logs
gcloud run services logs tail driftlock-api --region=us-central1

# Check for errors
gcloud logging read "resource.type=cloud_run_revision AND severity>=ERROR" --limit=10
```

### Verify Database

```bash
# Connect to database (if using Cloud SQL)
gcloud sql connect driftlock-db --user=postgres --project=driftlock

# Check tenant table
SELECT id, email, plan, stripe_status FROM tenants ORDER BY created_at DESC LIMIT 10;
```

### Monitor Stripe Events

```bash
# List recent events
stripe events list --limit 10

# Check specific event
stripe events evt_xxx
```

---

## 6. Launch Day Checklist

### Morning (Before Launch)
- [ ] Re-authenticate Firebase and Google Cloud
- [ ] Verify Cloud Run backend is deployed and responding
- [ ] Deploy latest Firebase Functions
- [ ] Deploy latest Firebase Hosting
- [ ] Test all API endpoints

### Midday (During Launch)
- [ ] Configure Stripe webhook secret in Functions
- [ ] Set up Cloudflare DNS records for Firebase
- [ ] Test Stripe checkout flow with test card
- [ ] Verify webhook processing
- [ ] Monitor error rates

### Afternoon (Post-Launch)
- [ ] Monitor logs for errors
- [ ] Test signup flow end-to-end
- [ ] Verify email delivery
- [ ] Check Firebase Functions invocations
- [ ] Confirm billing works for real customers

### Evening (Verification)
- [ ] Check all monitoring dashboards
- [ ] Verify SSL certificates
- [ ] Test mobile responsiveness
- [ ] Review Firebase usage/costs
- [ ] Celebrate! 🎉

---

## 7. Troubleshooting

### Cloud Run 404 Errors
**Problem:** Backend returning 404
**Solution:**
```bash
# Redeploy backend
cd collector-processor/cmd/driftlock-http
gcloud run deploy driftlock-api \
  --source=. \
  --region=us-central1 \
  --allow-unauthenticated \
  --set-env-vars="DATABASE_URL=...,STRIPE_SECRET_KEY=sk_test_..."
```

### Firebase Functions Timeout
**Problem:** Functions returning "unavailable"
**Solution:**
```bash
# Check invoker permissions
gcloud run services add-iam-policy-binding driftlock-api \
  --member="serviceAccount:driftlock@appspot.gserviceaccount.com" \
  --role="roles/run.invoker" \
  --region=us-central1
```

### Stripe Webhook Errors
**Problem:** Stripe receiving 5xx errors
**Solution:**
```bash
# Check Functions logs
firebase functions:log --only apiProxy

# Verify webhook secret matches
echo $STRIPE_WEBHOOK_SECRET
```

### Domain Not Resolving
**Problem:** driftlock.net not loading
**Solution:**
```bash
# Check DNS propagation
dig +short driftlock.net
dig +short driftlock.net @8.8.8.8

# Verify Firebase hosting config
firebase hosting Channel:list
```

---

## 8. Emergency Contacts & Resources

**Firebase Console:** https://console.firebase.google.com/project/driftlock
**Stripe Dashboard:** https://dashboard.stripe.com/
**Google Cloud Console:** https://console.cloud.google.com/home/dashboard?project=driftlock
**Cloudflare Dashboard:** https://dash.cloudflare.com/

**Support Tickets:**
- Firebase: Priority support via console
- Stripe: 24/7 support in dashboard
- Google Cloud: Billing support

---

## 9. Post-Launch Monitoring

**Week 1:**
- Monitor error rates hourly
- Check signup conversion rates
- Verify billing process works
- Review Firebase costs daily

**Month 1:**
- Set up alerting for errors > 1%
- Monitor API response times
- Track customer acquisition cost
- Review and optimize costs

**Ongoing:**
- Weekly security scans
- Monthly dependency updates
- Quarterly access reviews
- Continuous monitoring

---

**Ready to launch? Let's do this!** 🚀

Next steps:
1. Run the Stripe webhook setup script
2. Re-authenticate Firebase & Google Cloud
3. Follow the testing checklist
4. Deploy and monitor
