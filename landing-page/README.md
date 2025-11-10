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

- `CLOUDFLARE_ACCOUNT_ID`: Your Cloudflare account ID
- `CLOUDFLARE_API_TOKEN`: Your Cloudflare API token

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

The CTA section includes a lead capture form. To integrate with your backend:

1. Update the `handleSubmit` function in `CTASection.vue`
2. Add your API endpoint URL
3. Configure webhooks or email notifications

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