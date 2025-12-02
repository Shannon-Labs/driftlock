# 🚀 Launch Driftlock TODAY - Quick Start Guide

Perfect! You want to launch today. Let's make this happen. Here's your step-by-step action plan.

## ⏰ Timeline: ~2-3 Hours to Launch

---

## Phase 1: Fix Authentication (15 minutes)

Your Firebase and Google Cloud credentials expired. Let's fix that first.

**Run this command:**
```bash
./scripts/launch-driftlock.sh
```

Choose option **"3) Re-authenticate Firebase & Google Cloud"**

Or manually:
```bash
# Firebase re-auth
firebase login --reauth

# Google Cloud re-auth
gcloud auth login
gcloud config set project driftlock
```

**Verify:**
```bash
firebase projects:list  # Should show driftlock project
gcloud auth list        # Should show your account
```

---

## Phase 2: Deploy Backend & Fix 404 (30 minutes)

The Cloud Run backend is returning 404. Let's redeploy it.

```bash
cd collector-processor/cmd/driftlock-http

# Redeploy with correct settings
gcloud run deploy driftlock-api \
  --source=. \
  --region=us-central1 \
  --allow-unauthenticated \
  --set-env-vars="DATABASE_URL=postgresql://user:pass@host:5432/dbname" \
  --set-env-vars="STRIPE_SECRET_KEY=sk_test_your_key" \
  --set-env-vars="STRIPE_WEBHOOK_SECRET=whsec_your_secret" \
  --set-env-vars="STRIPE_PRICE_ID_PRO=price_xxxxx"

# Verify it's working
curl https://driftlock-api-o6kjgrsowq-uc.a.run.app/healthz
```

**Expected response:**
```json
{"success":true,"database":"healthy","license":"valid","version":"2.0.0"}
```

If you need the database credentials:
```bash
gcloud secrets access database-url --project=driftlock
gcloud secrets access stripe-secret-key --project=driftlock
```

---

## Phase 3: Configure Stripe Webhooks (20 minutes)

### Option A: Quick Setup

```bash
./scripts/launch-driftlock.sh
```

Choose option **"2) Set up Stripe webhooks"**

### Option B: Manual Setup

1. Go to https://dashboard.stripe.com/webhooks
2. Click "Add endpoint"
3. URL: `https://driftlock.net/webhooks/stripe`
4. Select these events:
   - `checkout.session.completed`
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_succeeded`
   - `invoice.payment_failed`
5. Copy the **Signing secret** (starts with `whsec_`)

**Store the secret:**
```bash
cd functions
firebase functions:config:set stripe.webhook_secret="whsec_your_secret_here"
cd ..
```

### Test Webhooks Locally (Optional)

```bash
# Terminal 1: Start Firebase emulator
cd functions
npm run build
firebase emulators:start --only functions

# Terminal 2: Forward webhooks
./scripts/stripe-webhook-forward.sh

# Terminal 3: Trigger test events
stripe trigger customer.subscription.created
stripe trigger checkout.session.completed
```

---

## Phase 4: Deploy Everything (30 minutes)

```bash
./scripts/launch-driftlock.sh
```

Choose option **"6) Full launch sequence (tests → deploy → verify)"**

Or run manually:

```bash
# Deploy Firebase Functions
cd functions
npm install
npm run build
firebase deploy --only functions
cd ..

# Deploy Firebase Hosting
cd landing-page
npm install
npm run build
cd ..
firebase deploy --only hosting

# Verify deployment
./scripts/test-launch-readiness.sh
```

---

## Phase 5: Configure Cloudflare Domain (15 minutes)

### If using Cloudflare + Firebase:

1. **Get Firebase IP addresses:**
   ```bash
   firebase hosting:sites:info
   ```

2. **In Cloudflare Dashboard:**
   - Go to DNS > Records
   - Add these A records:
     ```
     Type: A
     Name: @
     Content: 199.36.158.100
     Proxy status: DNS only (gray cloud)
     
     Type: A
     Name: @
     Content: 199.36.158.101
     Proxy status: DNS only (gray cloud)
     ```

3. **In Firebase Console:**
   - Go to Hosting > Domains
   - Add custom domain: `driftlock.net`
   - Wait for SSL certificate (5-10 minutes)

4. **Verify:**
   ```bash
   curl -I https://driftlock.net
   # Should show: HTTP/2 200
   ```

---

## Phase 6: End-to-End Testing (30 minutes)

### Test 1: User Signup Flow

1. Open browser: https://driftlock.net
2. Click "Get Started" or "Sign Up"
3. Enter email: **test-`date +%s`@driftlock.dev** (unique email)
4. Submit form
5. **Expected:** Receive API key via email within 30 seconds

**Check logs:**
```bash
firebase functions:log --limit 20
gcloud run services logs tail driftlock-api --region=us-central1
```

### Test 2: Stripe Checkout Flow

1. Visit: https://driftlock.net/dashboard (sign in if needed)
2. Click "Upgrade to Pro" or "Billing"
3. Should redirect to Stripe Checkout
4. Use test card:
   ```
   Card: 4242 4242 4242 4242
   Expiry: Any future date
   CVC: Any 3 digits
   ```
5. Complete checkout
6. **Expected:** Redirect back to dashboard with success message

**Verify in Stripe Dashboard:**
- Payment succeeded
- Customer created
- Subscription active

**Verify in database:**
```bash
# Check tenant plan updated
SELECT email, plan, stripe_status FROM tenants WHERE email = 'test-xxx@driftlock.dev';
# Should show: plan = 'pro', stripe_status = 'active'
```

### Test 3: API Usage

```bash
# Get your API key from email or database

# Test anomaly detection
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "test-launch",
    "events": [
      {"data": "normal event 1"},
      {"data": "normal event 2"},
      {"data": "ANOMALOUS EVENT WITH UNUSUAL PATTERN"},
      {"data": "normal event 3"}
    ]
  }' | jq

# Should return detected anomalies
```

### Test 4: Webhook Processing

**In Stripe Dashboard:**
1. Go to Developers > Webhooks
2. Find your endpoint: `https://driftlock.net/webhooks/stripe`
3. Click "Send test webhook"
4. Choose event: `customer.subscription.created`
5. Send test
6. **Expected:** 200 OK response

**Check Functions logs:**
```bash
firebase functions:log --only apiProxy --limit=10
```

---

## 🎯 You're Live! (5 minutes)

If all tests pass, you're ready to launch!

**Final verification:**
```bash
# Run complete test suite
./scripts/test-launch-readiness.sh

# Should show all green ✅
```

**Announce launch:**
- Post on social media
- Email early access list
- Update PH/YC if applicable
- Monitor closely for first hour

**Monitor these:**
```bash
# Real-time logs
firebase functions:log --tail
gcloud run services logs tail driftlock-api --region=us-central1

# Error rates
firebase monitoring:console  # Check Functions usage
```

---

## 🆘 Quick Troubleshooting

### Problem: Domain not working
```bash
# Check DNS
dig +short driftlock.net
dig +short driftlock.net @8.8.8.8

# Check Firebase hosting
firebase hosting:sites:versions:list

# If issues, redeploy
firebase deploy --only hosting
```

### Problem: Webhooks failing
```bash
# Check Functions logs
firebase functions:log --only apiProxy --limit=20

# Verify secret matches
firebase functions:config:get stripe.webhook_secret

# Test manually
curl -X POST https://driftlock.net/webhooks/stripe \
  -H "Stripe-Signature: test" \
  -d '{"test": true}'
```

### Problem: API returning 401
```bash
# Check API key is valid
curl https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/onboard/verify \
  -H "X-API-Key: YOUR_KEY"

# If fails, regenerate key
```

### Problem: Functions timeout
```bash
# Check invoker permissions
gcloud run services add-iam-policy-binding driftlock-api \
  --member="serviceAccount:driftlock@appspot.gserviceaccount.com" \
  --role="roles/run.invoker" \
  --region=us-central1
```

---

## 📞 Support Resources

**Documentation:**
- `LAUNCH_PLAN.md` - Full detailed plan
- `CLOUDFLARE_DEPLOYMENT.md` - Domain setup
- `AGENTS.md` - Project architecture

**Firebase Console:** https://console.firebase.google.com/project/driftlock
**Stripe Dashboard:** https://dashboard.stripe.com/
**Google Cloud Console:** https://console.cloud.google.com/project/driftlock

---

## ✅ Post-Launch Checklist

**First hour:**
- [ ] Monitor error rates (< 1%)
- [ ] Check signup flow works
- [ ] Verify Stripe test payment
- [ ] Monitor Firebase Functions invocations

**First day:**
- [ ] Check cost metrics
- [ ] Review logs for warnings
- [ ] Test all user flows again
- [ ] Set up alerts

**First week:**
- [ ] Daily cost monitoring
- [ ] Track conversion rates
- [ ] Review support requests
- [ ] Optimize based on usage

---

## 🎉 You're Ready!

You have everything you need to launch Driftlock today:

✅ **Complete launch scripts** in `scripts/`  
✅ **Detailed testing plan** in `LAUNCH_PLAN.md`  
✅ **Automated readiness tests**  
✅ **Stripe webhook forwarding** for local testing  
✅ **Troubleshooting guides**  

**Time to launch:**
```bash
# Quick start
./scripts/launch-driftlock.sh

# Or run full sequence
./scripts/launch-driftlock.sh
# Choose option 6: Full launch sequence
```

**You've got this!** 🚀

Questions? Check the docs or run:
```bash
./scripts/test-launch-readiness.sh
```

Good luck with the launch! 🎉
