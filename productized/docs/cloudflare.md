# DriftLock Cloudflare Integration Documentation

This document details how DriftLock leverages Cloudflare services for edge computing, security, and performance.

## Overview

DriftLock integrates with Cloudflare to provide:
- Global CDN for static assets
- DDoS protection and WAF
- Edge computing with Cloudflare Workers
- Secure custom domains with SSL
- Rate limiting at the edge
- Improved security and performance

## Components

### Cloudflare Pages (Frontend Hosting)

The React frontend is deployed to Cloudflare Pages with:
- Automatic builds from GitHub
- Global CDN distribution
- Custom domain support
- Free SSL certificates
- Preview deployments for pull requests

#### Configuration
- Build command: `npm run build`
- Build output directory: `build`
- Environment variables configured per environment

### Cloudflare Workers (Edge Computing)

Custom workers handle:
- Rate limiting for API endpoints
- Request filtering and caching
- Custom authentication logic at the edge
- Request/response transformations

#### Worker Configuration
- Rate limit: 100 requests per minute per IP for API endpoints
- Cache static assets for 5 minutes
- Block requests with malicious patterns

### DNS Configuration

DNS records are configured to:
- Point `app.driftlock.com` to Cloudflare Pages
- Point `api.driftlock.com` to backend API servers
- Enable Cloudflare proxy for security and performance

### Security Features

#### Web Application Firewall (WAF)
- Block common attack patterns (SQL injection, XSS, etc.)
- Custom rules for application-specific threats
- Regular updates to threat signatures

#### DDoS Protection
- Automatic mitigation of DDoS attacks
- Rate limiting at network level
- Challenge pages for suspicious traffic

### SSL/TLS Configuration
- SSL mode: "Full" (encrypts traffic between Cloudflare and origin)
- Minimum TLS version: 1.2
- Automatic HTTPS rewrites
- HSTS enabled

## Environment Configuration

### Production Environment
- Domain: `app.driftlock.com`
- API: `api.driftlock.com`
- SSL: Strict mode
- Security: High
- Caching: Aggressive

### Staging Environment
- Domain: `staging.driftlock.com`
- API: `staging-api.driftlock.com`
- SSL: Full mode
- Security: Medium
- Caching: Standard

## Performance Optimizations

### Caching
- Static assets cached for 1 year (with proper versioning)
- API responses cache where appropriate
- Browser caching headers for CSS and JS files

### Compression
- Brotli compression enabled
- Gzip as fallback
- Image optimization

### Edge Locations
- Leverage Cloudflare's 300+ data centers globally
- Route users to nearest edge location
- Minimize latency for all regions

## Monitoring

### Analytics
- Page load times
- Error rates
- API response times
- Request volumes

### Logs
- Access logs for debugging
- Security event logs
- Performance metrics

## Troubleshooting

### Common Issues

1. **SSL Certificate Issues**
   - Ensure Cloudflare proxy is enabled (orange cloud icon)
   - Check DNS settings are correct

2. **Rate Limiting Too Aggressive**
   - Adjust rate limit rules in Cloudflare dashboard
   - Monitor actual usage patterns

3. **Caching Issues**
   - Clear Cloudflare cache when needed
   - Check cache rules and headers

4. **Worker Errors**
   - Monitor worker logs in Cloudflare dashboard
   - Test with staging worker first

## Security Best Practices

1. **API Protection**
   - Always proxy API through Cloudflare
   - Use rate limiting to prevent abuse
   - Block known malicious IPs

2. **Origin Protection**
   - Never expose origin server IP directly
   - Only allow requests through Cloudflare
   - Use Cloudflare-only DNS records

3. **Certificate Management**
   - Use Cloudflare's Origin CA certificates
   - Enable Always Use HTTPS
   - Monitor certificate expiration

## Deployment Process

1. Deploy frontend to Cloudflare Pages
2. Update DNS to point to new version (if needed)
3. Test all functionality through Cloudflare
4. Monitor for any issues in Cloudflare dashboard
5. Update security rules as needed based on traffic patterns