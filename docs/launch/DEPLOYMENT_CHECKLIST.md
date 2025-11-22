# Driftlock Deployment Checklist

## ✅ Completed

- [x] Frontend deployed to Firebase Hosting (`https://driftlock.web.app`)
- [x] Backend API deployed to Cloud Run (`driftlock-api`)
- [x] API proxy function deployed and public
- [x] Organization Policy "Domain Restricted Sharing" disabled
- [x] Cloud Run services made public (`driftlock-api`, `apiproxy`)
- [x] Firebase Hosting rewrite rules configured (`/api/v1/**` → `apiproxy`)
- [x] Signup form updated to use Firebase Auth
- [x] API crypto test script created (`scripts/api_crypto_test.py`)

## 🔧 Required Before Launch

### 1. Stripe Configuration

**Required Secret:** `stripe-price-id-pro` (already referenced in `cloudbuild.yaml`)

**Action Items:**
1. Create Stripe products/prices in your Stripe dashboard:
   - **Radar Plan**: $20/month (or your pricing)
   - **Lock Plan**: $200/month (or your pricing)
2. Get the Price IDs (format: `price_xxxxx`)
3. Update Secret Manager:
   ```bash
   echo -n "price_xxxxx" | gcloud secrets create stripe-price-id-pro --data-file=-
   ```
4. Or update existing secret:
   ```bash
   echo -n "price_xxxxx" | gcloud secrets versions add stripe-price-id-pro --data-file=-
   ```

**Current Status:** Check `cloudbuild.yaml` for `STRIPE_PRICE_ID_PRO` secret reference.

### 2. Database Migration (Firebase UID)

**Action:** Add `firebase_uid` column to `tenants` table if not already present.

```sql
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS firebase_uid TEXT UNIQUE;
CREATE INDEX IF NOT EXISTS idx_tenants_firebase_uid ON tenants(firebase_uid);
```

**Status:** Backend code updated to accept Firebase tokens, but database schema needs migration.

### 3. Test Signup Flow

**Action:** Test the complete signup flow:
1. Visit `https://driftlock.web.app`
2. Fill out signup form
3. Verify Firebase Auth user is created
4. Verify tenant is created in database
5. Verify API key is returned
6. Test API key works: `curl -X POST https://driftlock.web.app/api/v1/detect -H "X-Api-Key: <key>" ...`

### 4. Run Live Crypto Test

**Action:** Start the API-based crypto test to demonstrate the platform:

```bash
# Get an API key first (sign up or use existing)
export DRIFTLOCK_API_KEY="dlk_..."
export DRIFTLOCK_API_URL="https://driftlock.web.app/api/v1"

# Run the test
python3 scripts/api_crypto_test.py
```

This will:
- Stream live Binance crypto trade data
- Send batches to your API
- Display detected anomalies in real-time
- Perfect for demos and validation

### 5. Environment Variables

**Verify these are set in Cloud Run:**
- `FIREBASE_SERVICE_ACCOUNT_KEY` - For Firebase Auth verification
- `DATABASE_URL` - PostgreSQL connection string
- `DRIFTLOCK_LICENSE_KEY` - If required
- `SENDGRID_API_KEY` - For welcome emails
- `STRIPE_PRICE_ID_PRO` - Stripe price ID (see above)
- `GEMINI_API_KEY` - For AI analysis features

## 📋 Post-Launch Monitoring

1. **Monitor API Usage:**
   - Check Cloud Run logs for errors
   - Monitor rate limiting
   - Track signup success rate

2. **Monitor Costs:**
   - Cloud Run instance hours
   - Firebase Functions invocations
   - Gemini API calls (if enabled)
   - Database connections

3. **User Feedback:**
   - Monitor signup completion rate
   - Track API key usage
   - Collect error reports

## 🚀 Quick Start Commands

### Deploy Everything
```bash
# Frontend
cd landing-page && npm run build && cd .. && firebase deploy --only hosting

# Backend (via Cloud Build)
gcloud builds submit --config cloudbuild.yaml .
```

### Test API
```bash
# Health check
curl https://driftlock.web.app/api/v1/healthz

# Test detection (requires API key)
curl -X POST https://driftlock.web.app/api/v1/detect \
  -H "X-Api-Key: <your-key>" \
  -H "Content-Type: application/json" \
  -d '{"events": [{"message": "test"}], "window_size": 50}'
```

### Run Crypto Test
```bash
export DRIFTLOCK_API_KEY="dlk_..."
python3 scripts/api_crypto_test.py
```

