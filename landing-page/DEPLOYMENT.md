# DriftLock Landing Page Deployment Guide

## Prerequisites

- Cloudflare account with Pages access
- Domain name (driftlock.net) registered
- Wrangler CLI installed (`npm install -g wrangler`)
- Node.js 18+ and npm installed

## Step 1: Cloudflare Authentication

```bash
# Login to Cloudflare
npx wrangler auth login

# Verify authentication
npx wrangler whoami
```

## Step 2: Deploy to Cloudflare Pages

```bash
# Navigate to landing page directory
cd landing-page

# Install dependencies
npm install

# Deploy to Cloudflare Pages
npm run deploy
```

The deployment will output a URL like `https://driftlock-landing.pages.dev`

## Step 3: Domain Configuration

### Option A: Using Cloudflare Dashboard

1. Go to Cloudflare Dashboard → Pages → DriftLock Landing
2. Click "Custom domains"
3. Add `driftlock.net`
4. Add `www.driftlock.net` (optional)

### Option B: Using Wrangler CLI

```bash
# Add custom domain
npx wrangler pages project create driftlock-landing --production-branch main
npx wrangler pages domain add driftlock.net
```

## Step 4: DNS Configuration

### At Your Domain Registrar

Update your DNS settings to point to Cloudflare:

```
Type: A
Name: @ (or driftlock.net)
Value: 192.0.2.1  # Use actual Cloudflare Pages IP
TTL: 300 (Automatic)

Type: CNAME
Name: www
Value: driftlock.net
TTL: 300 (Automatic)
```

### Cloudflare Nameservers

If using Cloudflare for DNS management:

```
Nameserver 1: dina.ns.cloudflare.com
Nameserver 2: walt.ns.cloudflare.com
```

## Step 5: SSL Certificate

Cloudflare automatically provisions SSL certificates:

- **Universal SSL**: Free, auto-renewing certificate
- **Advanced SSL**: Available for paid plans
- **HSTS**: Enable additional security headers

## Step 6: Performance Optimization

### Build Optimization

The project includes several optimizations:

```bash
# Analyze bundle size
npm run build

# Check for unused dependencies
npm audit

# Preview build locally
npm run preview
```

### Cloudflare Features

Enable these Cloudflare features:

1. **Brotli compression**: Automatic
2. **HTTP/2**: Automatic
3. **HTTP/3**: Enable in settings
4. **WebP conversion**: Enable in speed optimization
5. **Minify HTML/CSS/JS**: Enable in speed optimization

## Step 7: Analytics and Monitoring

### Cloudflare Analytics

1. Go to Cloudflare Dashboard → Analytics & Logs
2. Enable detailed analytics
3. Set up custom dashboards

### External Analytics

Add to `index.html`:

```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_MEASUREMENT_ID');
</script>
```

### Form Submission Setup

Configure the CTA form in `CTASection.vue`:

```javascript
const handleSubmit = async () => {
  // Add your API endpoint here
  const response = await fetch('https://api.driftlock.com/leads', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(form)
  });

  // Handle response
}
```

## Step 8: Testing and Verification

### Pre-Launch Checklist

- [ ] All forms work correctly
- [ ] Mobile responsive design confirmed
- [ ] Dark/light mode switching works
- [ ] Page load time < 3 seconds
- [ ] SSL certificate active
- [ ] DNS propagation complete
- [ ] Analytics tracking installed
- [ ] SEO meta tags correct

### Performance Testing

```bash
# Lighthouse CLI testing
npm install -g lighthouse
lighthouse https://driftlock.net --output html --output-path ./lighthouse-report.html

# WebPageTest API
curl "https://www.webpagetest.org/runtest.php?url=https://driftlock.net&f=json&k=YOUR_API_KEY"
```

### Cross-Browser Testing

Test in:
- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers

## Step 9: Launch

### DNS Propagation

Check DNS propagation:

```bash
# Check A record
dig driftlock.net A

# Check propagation status
for server in 8.8.8.8 1.1.1.1 208.67.222.222; do
  echo "Querying $server:"
  dig @$server driftlock.net A +short
done
```

### Final Verification

1. **HTTPS**: Ensure SSL is working
2. **Redirects**: Test www to non-www redirect
3. **Performance**: Verify load times
4. **Forms**: Test lead capture
5. **Mobile**: Confirm responsive design

## Troubleshooting

### Common Issues

**Build Failures:**
```bash
# Clear cache
rm -rf node_modules dist
npm install
npm run build
```

**DNS Issues:**
```bash
# Flush DNS cache
sudo dscacheutil -flushcache  # macOS
sudo systemctl flush-dns      # Linux
ipconfig /flushdns           # Windows
```

**SSL Issues:**
- Check certificate status in Cloudflare dashboard
- Verify DNS CNAME/A records
- Ensure SSL mode is "Full (strict)"

**Performance Issues:**
- Enable Brotli compression
- Check Cloudflare caching settings
- Analyze bundle size with `npm run build`

## Maintenance

### Regular Updates

```bash
# Update dependencies
npm update

# Security audit
npm audit

# Re-deploy after updates
npm run deploy
```

### Monitoring

Monitor:
- Page load times
- Form submission rates
- Error rates
- SSL certificate expiration

### Backup

Keep backups of:
- Source code in git
- Configuration files
- DNS records
- Analytics data

## Support

For deployment issues:

1. Check Cloudflare status page
2. Review error logs in dashboard
3. Contact Cloudflare support
4. File GitHub issues for project-specific problems

---

**Note**: This deployment guide assumes you have admin access to both Cloudflare and your domain registrar. Adjust steps according to your specific setup.