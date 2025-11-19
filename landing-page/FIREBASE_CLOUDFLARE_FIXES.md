# Firebase & Cloudflare Pages Fixes

This document outlines the fixes applied to resolve two key issues in the Driftlock landing page deployment.

## Issues Fixed

### 1. `auth/invalid-api-key` Error

**Problem**: Cloudflare Pages didn't have Firebase environment variables, causing authentication failures.

**Solution**: 
- Added all required Firebase environment variables to `.env.production`
- Updated `src/firebase.ts` to gracefully handle missing environment variables
- Updated `src/stores/auth.ts` to handle null auth instances
- Created automated setup script at `scripts/setup-cloudflare-env.sh`

**Required Environment Variables**:
```
VITE_FIREBASE_API_KEY
VITE_FIREBASE_AUTH_DOMAIN
VITE_FIREBASE_PROJECT_ID
VITE_FIREBASE_STORAGE_BUCKET
VITE_FIREBASE_MESSAGING_SENDER_ID
VITE_FIREBASE_APP_ID
VITE_FIREBASE_MEASUREMENT_ID (optional)
```

### 2. `mce-autosize-textarea already defined` Error

**Problem**: Cloudflare Pages DevTools overlay was causing custom element conflicts.

**Solution**: 
- Removed the overlay import logic from `src/main.ts`
- Added documentation explaining that Cloudflare handles the overlay automatically
- The overlay is only used for preview builds and development, not production

## Setup Instructions

### Automatic Setup (Recommended)

Run the setup script to configure Firebase environment variables:

```bash
cd landing-page
./scripts/setup-cloudflare-env.sh
```

### Manual Setup

1. **Via Wrangler CLI:**
   ```bash
   wrangler pages secret put VITE_FIREBASE_API_KEY --project-name="driftlock"
   # Repeat for each variable
   ```

2. **Via Cloudflare Dashboard:**
   - Go to Dashboard → Pages → driftlock → Settings → Environment Variables
   - Add each `VITE_FIREBASE_*` variable

3. **Get Firebase Values:**
   - Firebase Console → Project Settings → General → Your apps
   - Copy values from "SDK setup and configuration"

### Deploy After Setup

After setting environment variables, redeploy:

```bash
npm run build
npm run deploy:cloudflare
```

Or trigger a new deployment by pushing to your connected git branch.

## Files Modified

1. **`.env.production`** - Added Firebase environment variable placeholders
2. **`src/firebase.ts`** - Added graceful handling of missing environment variables
3. **`src/stores/auth.ts`** - Added null checks for auth instance
4. **`src/main.ts`** - Removed problematic overlay import logic
5. **`scripts/setup-cloudflare-env.sh`** - New automated setup script
6. **`README.md`** - Updated with Firebase setup instructions

## Verification

After applying fixes:

1. ✅ Build completes successfully (`npm run build`)
2. ✅ No TypeScript errors
3. ✅ Firebase auth gracefully handles missing config
4. ✅ No custom element conflicts in Cloudflare Pages
5. ✅ Automated setup script available for easy configuration

## Notes

- Firebase config validation warns in console if variables are missing but doesn't break the app
- Auth functions return appropriate errors when Firebase is not configured
- The landing page continues to function even without Firebase auth configured
- Cloudflare Pages DevTools overlay is handled automatically in production