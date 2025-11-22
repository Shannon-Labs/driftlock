# Migration Plan: Switch Everything to driftlock.net

## Overview
Migrate all services to use `driftlock.net` as the primary domain, replacing `driftlock.pages.dev` references.

## Current State
- **Frontend**: `https://driftlock.pages.dev` (Cloudflare Pages)
- **API**: `https://driftlock.net/api/v1` (Cloudflare Workers/Pages Functions)
- **CORS**: API currently allows `*` or needs to be configured for specific origins

## Target State
- **Frontend**: `https://driftlock.net` (Cloudflare Pages with custom domain)
- **API**: `https://driftlock.net/api/v1` (unchanged)
- **CORS**: API configured to allow `https://driftlock.net` and `https://www.driftlock.net`

---

## Implementation Steps

### Phase 1: Frontend Configuration

#### 1.1 Update Environment Variables
- [ ] **File**: `landing-page/.env.production` (create if doesn't exist)
  - Set `VITE_API_BASE_URL=https://driftlock.net/api/v1`
  - This ensures production builds use the correct API URL

#### 1.2 Update Default API URL in Code
- [ ] **File**: `landing-page/src/views/PlaygroundView.vue`
  - Current: `const apiBase = import.meta.env.VITE_API_BASE_URL || 'https://driftlock.net/api/v1'`
  - Status: ✅ Already correct, no change needed

#### 1.3 Update Wrangler Configuration
- [ ] **File**: `landing-page/wrangler.toml`
  - Add custom domain configuration (optional, can be done via dashboard)
  - Current configuration is fine, domain will be added via Cloudflare dashboard

#### 1.4 Update Documentation References
- [ ] **File**: `landing-page/README.md`
  - Update deployment URLs from `driftlock.pages.dev` to `driftlock.net`
  - Update any example URLs in documentation

- [ ] **File**: `landing-page/DEPLOYMENT.md`
  - Update all references to use `driftlock.net` instead of `driftlock.pages.dev`
  - Verify DNS configuration instructions are correct

- [ ] **File**: `landing-page/YC_APPLICATION_REVIEW.md`
  - Update Company URL from `https://driftlock.pages.dev` to `https://driftlock.net`

---

### Phase 2: API CORS Configuration

#### 2.1 Update API CORS Environment Variable
- [ ] **Location**: Cloudflare Dashboard → Workers & Pages → Your API Project → Settings → Variables
  - Variable: `CORS_ALLOW_ORIGINS`
  - Current value: Check current setting (may be `*` or empty)
  - New value: `https://driftlock.net,https://www.driftlock.net`
  - Environment: **Production** (and Preview if needed)
  - **Note**: See `docs/CLOUDFLARE_CORS_SETUP.md` for detailed step-by-step instructions

#### 2.2 Verify Changes Applied
- [ ] If using Cloudflare Workers: Changes take effect immediately
- [ ] If using Pages Functions: Ensure latest deployment is active
- [ ] Test API endpoint to verify CORS headers are correct

#### 2.3 Verify CORS Implementation
- [ ] **File**: `collector-processor/cmd/driftlock-http/main.go`
  - Verify CORS logic handles multiple origins correctly (lines 347-357, 558-583)
  - Status: ✅ Code already supports comma-separated origins via `parseAllowedOrigins()`
  - The `originAllowed()` function checks exact matches (case-insensitive)

#### 2.4 Test CORS Configuration
- [ ] Test from `https://driftlock.net/playground` that API calls work
- [ ] Verify preflight OPTIONS requests return correct headers
- [ ] Check browser console for CORS errors
- [ ] Test API health check shows "API Connected" status

---

### Phase 3: Cloudflare Pages Domain Configuration

#### 3.1 Add Custom Domain via Dashboard
- [ ] Go to Cloudflare Dashboard → Pages → `driftlock` project
- [ ] Navigate to "Custom domains" section
- [ ] Click "Set up a custom domain"
- [ ] Enter: `driftlock.net`
- [ ] Cloudflare will automatically configure DNS (if domain is managed by Cloudflare)

#### 3.2 Add www Subdomain (Optional)
- [ ] Add `www.driftlock.net` as additional custom domain
- [ ] Configure redirect: `www.driftlock.net` → `driftlock.net` (301 redirect)

#### 3.3 Verify SSL Certificate
- [ ] Wait for Cloudflare to provision SSL certificate (usually automatic, takes 1-5 minutes)
- [ ] Verify SSL is active: `curl -I https://driftlock.net`
- [ ] Check certificate validity in browser

#### 3.4 Update DNS Records (if domain not managed by Cloudflare)
- [ ] **If domain registrar is separate from Cloudflare:**
  - Add CNAME record: `driftlock.net` → `driftlock.pages.dev`
  - Or add A record pointing to Cloudflare Pages IP (if provided)
  - Add CNAME for www: `www.driftlock.net` → `driftlock.net`

---

### Phase 4: Testing & Verification

#### 4.1 Frontend Testing
- [ ] Verify `https://driftlock.net` loads correctly
- [ ] Test navigation: Home → Playground → Home
- [ ] Verify all internal links work
- [ ] Test dark mode toggle
- [ ] Verify responsive design on mobile

#### 4.2 Playground Testing
- [ ] Navigate to `https://driftlock.net/playground`
- [ ] Verify API health check shows "API Connected"
- [ ] Test file upload functionality
- [ ] Test sample data loading
- [ ] Test anomaly detection with sample data
- [ ] Verify results display correctly
- [ ] Check browser console for errors

#### 4.3 API Integration Testing
- [ ] Test API health endpoint: `curl https://driftlock.net/api/v1/healthz`
- [ ] Test CORS preflight: `curl -X OPTIONS https://driftlock.net/api/v1/detect -H "Origin: https://driftlock.net" -v`
- [ ] Verify `Access-Control-Allow-Origin` header includes `https://driftlock.net`
- [ ] Test actual API call from playground UI

#### 4.4 SEO & Meta Tags
- [ ] Verify canonical URLs use `driftlock.net`
- [ ] Check Open Graph tags reference correct domain
- [ ] Verify sitemap.xml (if exists) uses correct domain

---

### Phase 5: Cleanup & Redirects

#### 5.1 Set Up Redirect from Old Domain
- [ ] **Option A**: Keep `driftlock.pages.dev` active and add redirect rule
  - In Cloudflare Pages → Custom domains → Add redirect rule
  - `driftlock.pages.dev/*` → `https://driftlock.net$1` (301 redirect)

- [ ] **Option B**: Remove `driftlock.pages.dev` custom domain
  - Only if you want to completely deprecate the `.pages.dev` URL
  - **Recommendation**: Keep redirect active for at least 30 days

#### 5.2 Update External References
- [ ] Update GitHub repository description/URLs
- [ ] Update any marketing materials with new URL
- [ ] Update social media profiles
- [ ] Update email signatures
- [ ] Update any partner/integration documentation

---

## Environment Variables Summary

### Frontend (Cloudflare Pages)
```bash
# Production environment variables (set in Cloudflare Pages dashboard)
VITE_API_BASE_URL=https://driftlock.net/api/v1
```

### Backend API (Cloudflare Workers/Pages Functions)
```bash
# Production environment variables (set in Cloudflare dashboard or wrangler.toml)
CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net
# Other variables as needed for your API
```

---

## Rollback Plan

If issues occur during migration:

1. **Keep `driftlock.pages.dev` active** as fallback
2. **Revert CORS changes** in Vercel if API breaks
3. **Remove custom domain** from Cloudflare Pages if needed
4. **DNS changes** can be reverted at registrar

---

## Post-Migration Checklist

- [ ] All tests passing
- [ ] No CORS errors in browser console
- [ ] SSL certificate valid and active
- [ ] Redirect from old domain working (if configured)
- [ ] Documentation updated
- [ ] External references updated
- [ ] Monitor error logs for 24-48 hours

---

## Notes

- **DNS Propagation**: Changes may take up to 48 hours to propagate globally, though usually much faster
- **SSL Provisioning**: Cloudflare automatically provisions SSL certificates, usually within minutes
- **Zero Downtime**: Migration can be done with zero downtime by keeping both domains active during transition
- **CORS**: The API already supports multiple origins via comma-separated `CORS_ALLOW_ORIGINS` value

---

## Timeline Estimate

- **Phase 1** (Frontend Config): 15 minutes
- **Phase 2** (API CORS): 5 minutes
- **Phase 3** (Domain Setup): 10 minutes + wait for SSL (5-15 minutes)
- **Phase 4** (Testing): 30 minutes
- **Phase 5** (Cleanup): 15 minutes

**Total**: ~1.5 hours (excluding DNS propagation wait time)

---

## Questions to Resolve

1. **Do you want to keep `www.driftlock.net`?** (Recommended: Yes, with redirect to non-www)
2. **Do you want to redirect `driftlock.pages.dev`?** (Recommended: Yes, 301 redirect)
3. **Is the domain `driftlock.net` already managed by Cloudflare?** (Affects DNS setup steps)

