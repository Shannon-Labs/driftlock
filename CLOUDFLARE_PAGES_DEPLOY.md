# Deploying Driftlock Marketing Site to Cloudflare Pages

The Driftlock marketing site (driftlock.net) is built with Vite + React and deployed to Cloudflare Pages.

## Quick Deploy

```bash
cd web-frontend
npm install
npm run build
```

Then deploy via Cloudflare Dashboard or Wrangler:

```bash
npm install -g wrangler
wrangler pages deploy dist --project-name=driftlock
```

## Configuration

### Environment Variables

Set these in Cloudflare Pages dashboard:

**Production:**
- `VITE_API_URL` - API endpoint for the dashboard (e.g., `https://api.driftlock.net`)
- `VITE_SUPABASE_PROJECT_ID` - Supabase project ID (optional)
- `VITE_SUPABASE_PUBLISHABLE_KEY` - Supabase anon key (optional)
- `VITE_STRIPE_PUBLISHABLE_KEY` - Stripe publishable key (optional)

**Staging:**
- Same variables with staging values

### Build Configuration

**Build command:** `npm run build`
**Build output directory:** `dist`
**Node.js version:** 18+

### Custom Domains

1. Add custom domain in Cloudflare Pages dashboard
2. Update DNS records to point to Cloudflare Pages
3. Configure SSL/TLS to Full (Strict)

## Local Development

```bash
cd web-frontend
npm install
npm run dev
```

Visit http://localhost:3000

## Features

The marketing site includes:

- **Landing Page** (`/`): Product overview, hero section, feature highlights
- **Documentation** (`/documentation`): API docs, quick start guide
- **Pricing** (`/pricing`): Tiered pricing with feature comparison
- **About** (`/about`): Company information, team
- **Security** (`/security`): Compliance highlights (DORA/NIS2/EU AI Act)
- **Contact** (`/contact`): Contact form, support information
- **Dashboard** (`/dashboard`): Main application (requires API key)

## Compliance Pages

Special pages for regulatory compliance:

- **DORA Compliance**: `/security` and `/documentation#dora`
- **NIS2 Incident Reporting**: `/security#incident-response`
- **EU AI Act Transparency**: `/security#ai-act`
- **Audit Trails**: Built into the dashboard application

## Performance

- Static site generation where possible
- Optimized images and assets
- Cloudflare CDN for global distribution
- Automatic HTTPS and HTTP/3

## Monitoring

- Cloudflare Web Analytics
- Custom event tracking (respects DNT)
- Error tracking via Sentry (optional)

## SEO

- Meta tags for social sharing
- Sitemap generation
- Structured data markup
- Open Graph and Twitter Card support

## Deployment Checklist

- [ ] Set environment variables in Cloudflare Dashboard
- [ ] Configure custom domain (driftlock.net)
- [ ] Enable automatic deployments from Git
- [ ] Set up branch deployments (main → production, develop → staging)
- [ ] Configure build notifications
- [ ] Test contact forms and CTAs
- [ ] Verify compliance page content
- [ ] Check analytics integration
- [ ] Validate mobile responsiveness
- [ ] Performance audit (Lighthouse)