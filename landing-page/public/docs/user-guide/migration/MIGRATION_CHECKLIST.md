# Quick Migration Checklist: driftlock.net

## Pre-Migration (Before Adding Domain)

### Frontend Code Updates
- [x] API URL already defaults to `https://driftlock.net/api/v1` ✅
- [ ] Create `.env.production` with `VITE_API_BASE_URL=https://driftlock.net/api/v1`
- [ ] Update documentation URLs

### API Configuration
- [ ] Set `CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net` in Vercel

---

## Domain Setup (After You Add Domain)

### Cloudflare Pages
- [ ] Add custom domain `driftlock.net` in Cloudflare Dashboard
- [ ] (Optional) Add `www.driftlock.net` with redirect to `driftlock.net`
- [ ] Wait for SSL certificate provisioning (~5 minutes)
- [ ] Verify `https://driftlock.net` loads

### Testing
- [ ] Test homepage loads: `https://driftlock.net`
- [ ] Test playground: `https://driftlock.net/playground`
- [ ] Verify API connection works (no CORS errors)
- [ ] Test navigation between routes
- [ ] Check browser console for errors

### Redirects (Optional)
- [ ] Set up redirect: `driftlock.pages.dev` → `driftlock.net` (301)

---

## Post-Migration

- [ ] Update GitHub repo URLs
- [ ] Update marketing materials
- [ ] Monitor for 24-48 hours

---

**Full details**: See `docs/MIGRATION_TO_DRIFTLOCK_NET.md`

