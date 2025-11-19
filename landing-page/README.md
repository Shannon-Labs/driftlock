# DriftLock Landing Page

Professional landing page for DriftLock - Explainable Anomaly Detection for EU Banks.

## Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Deploy to Cloudflare Pages
npm run deploy
```

## Features

- ✅ **Responsive Design** - Works perfectly on desktop, tablet, and mobile
- ✅ **Dark/Light Mode** - Automatic theme switching with system preference
- ✅ **High Performance** - Built with Vite for lightning-fast loading
- ✅ **SEO Optimized** - Complete meta tags and structured data
- ✅ **Lead Capture** - Forms ready for integration with your backend
- ✅ **Professional Animations** - Smooth, engaging user experience
- ✅ **DORA Compliance Focus** - Tailored messaging for EU banking regulations

## Technology Stack

- **Framework**: Vue.js 3 with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Icons**: Lucide Vue
- **Deployment**: Cloudflare Pages
- **Type Checking**: TypeScript + Vue TSC

## Project Structure

```
landing-page/
├── src/
│   ├── components/          # Vue components
│   │   ├── HeroSection.vue
│   │   ├── ProblemSection.vue
│   │   ├── SolutionSection.vue
│   │   ├── ProofSection.vue
│   │   ├── ComparisonSection.vue
│   │   └── CTASection.vue
│   ├── assets/             # Static assets
│   ├── main.ts            # App entry point
│   ├── style.css          # Global styles
│   └── App.vue            # Main app component
├── public/                # Public assets
├── dist/                  # Build output
├── package.json
├── vite.config.ts
├── tailwind.config.js
├── wrangler.toml          # Cloudflare Pages config
└── README.md
```

## Development

### Running Locally

```bash
# Install dependencies
npm install

# Start development server (http://localhost:5173)
npm run dev
```

### Type Checking

```bash
# Run TypeScript type checking
npm run type-check
```

## Deployment

### Cloudflare Pages

The project is configured for Cloudflare Pages deployment:

```bash
# Build and deploy to Cloudflare Pages
npm run deploy
```

#### Manual Deployment

1. Build the project:
   ```bash
   npm run build
   ```

2. Deploy using Wrangler CLI:
   ```bash
   npx wrangler pages publish dist
   ```

#### Environment Variables

The landing page requires Firebase environment variables for authentication. Set these in Cloudflare Pages:

**Required Firebase Variables:**
- `VITE_FIREBASE_API_KEY`
- `VITE_FIREBASE_AUTH_DOMAIN`
- `VITE_FIREBASE_PROJECT_ID`
- `VITE_FIREBASE_STORAGE_BUCKET`
- `VITE_FIREBASE_MESSAGING_SENDER_ID`
- `VITE_FIREBASE_APP_ID`
- `VITE_FIREBASE_MEASUREMENT_ID` (optional, for Analytics)

**Setup Options:**

1. **Automatic Setup (Recommended):**
   ```bash
   ./scripts/setup-cloudflare-env.sh
   ```

2. **Manual Setup via Wrangler CLI:**
   ```bash
   wrangler pages secret put VITE_FIREBASE_API_KEY --project-name="driftlock"
   wrangler pages secret put VITE_FIREBASE_AUTH_DOMAIN --project-name="driftlock"
   # ... repeat for each variable
   ```

3. **Manual Setup via Cloudflare Dashboard:**
   - Go to Dashboard → Pages → driftlock → Settings → Environment Variables
   - Add each `VITE_FIREBASE_*` variable with values from your Firebase Console

**Getting Firebase Values:**
1. Go to Firebase Console → Project Settings → General → Your apps
2. Copy the config values from the "SDK setup and configuration" section

After setting environment variables, redeploy the site to apply changes.

### DNS Configuration

After deployment, configure your DNS:

1. **A Record** (for driftlock.net):
   ```
   Type: A
   Name: @
   Value: Cloudflare Pages IP (provided after deployment)
   TTL: 300
   ```

2. **CNAME Record** (for www):
   ```
   Type: CNAME
   Name: www
   Value: driftlock.net
   TTL: 300
   ```

## Customization

### Colors and Branding

Edit `tailwind.config.js` to customize:

```javascript
theme: {
  extend: {
    colors: {
      driftlock: {
        blue: '#1a73e8',
        teal: '#6c757d',
        orange: '#ff9800',
      }
    }
  }
}
```

### Content

Main content is organized in components:

- **HeroSection.vue**: Main value proposition and CTAs
- **ProblemSection.vue**: DORA compliance challenges
- **SolutionSection.vue**: CBAD technology explanation
- **ProofSection.vue**: Demo integration and evidence
- **ComparisonSection.vue**: Competitive differentiation
- **CTASection.vue**: Lead capture and contact forms

### Forms

The contact CTA now submits to `/api/v1/contact`, handled by a Cloudflare Pages Function (`functions/api/v1/contact.ts`). To configure it:

- **Step 1:** Deploy the Pages Function and set Cloudflare env variables such as `CRM_WEBHOOK_URL` (CRM/webhook endpoint) and, optionally, `CONTACT_LOG_KV` (KV namespace for logging submissions when no webhook exists).
- **Step 2:** Update `wrangler.toml` or the Cloudflare dashboard so the KV binding (if configured) is available to the function.
- **Step 3:** Customize the UI copy or fields inside `src/views/HomeView.vue` if you need additional data points.

Frontend validation lives in `handleContactSubmit` within `HomeView.vue`, mirroring the server-side checks so visitors get instant feedback. When no webhook is configured, the function logs submissions but still returns a friendly success state.

### Cloudflare Contact Proxy

- **File**: `functions/api/v1/contact.ts`
- **Behavior**: Validates payloads, forwards them to `CRM_WEBHOOK_URL` when set, otherwise logs to KV/console so the UX never fails.
- **Extensibility**: You can extend this file to send emails (e.g., via SendGrid) or enrich payloads with geolocation data provided by Cloudflare.

### QA Workflow (Manual)

1. **Type checking/build** – `npm run type-check && npm run build` before opening a PR.
2. **Lighthouse** – Use Chrome DevTools → Lighthouse → Run (Mobile + Desktop) to capture Performance/Accessibility/Best Practices/SEO scores. Attach JSON to the PR if scores change.
3. **DevTools Recorder** – Record a flow that tabs through nav, opens dark mode, and submits the contact form (happy + error paths). Export the flow for future Playwright automation.
4. **Smoke test** – Use the Cloudflare preview deployment and confirm `/api/v1/contact` returns `{"ok": true}` when `CRM_WEBHOOK_URL` is unset.

## Performance

### Lighthouse Scores

The site is optimized for high Lighthouse scores:

- **Performance**: 95+
- **Accessibility**: 100
- **Best Practices**: 100
- **SEO**: 100

### Optimization Features

- Lazy loading for images
- Code splitting for vendor libraries
- Optimized bundle sizes
- CSS purging with Tailwind
- Gzip compression enabled

## Analytics and Tracking

### Cloudflare Analytics

Cloudflare Pages includes built-in analytics. Track:

- Page views and unique visitors
- Geographic distribution
- Device and browser breakdown
- Performance metrics

### Custom Analytics

To add custom analytics (Google Analytics, etc.):

1. Add tracking script to `index.html`
2. Update `App.vue` to emit custom events
3. Configure conversion tracking for the CTA form

## Security

### Content Security Policy

The site includes basic CSP headers. Customize as needed:

```html
<meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';">
```

### Form Security

- Input validation on frontend and backend
- Rate limiting for form submissions
- CSRF protection for form endpoints
- SSL/HTTPS enforced

## Browser Support

- Chrome 88+
- Firefox 85+
- Safari 14+
- Edge 88+

## Contributing

1. Follow the existing code style
2. Use TypeScript for all new code
3. Test on multiple screen sizes
4. Run type checking before committing
5. Check Lighthouse scores after changes

## License

Licensed under Apache 2.0 - see the main project license for details.

## Support

For technical issues or questions:

- Create an issue in the main DriftLock repository
- Contact the development team
- Check the documentation at `/docs`
