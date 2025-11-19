# Driftlock Agents — `landing-page/` (SaaS Frontend)

These instructions apply to all work under `landing-page/` - the **primary customer-facing SaaS platform**.

---

## 1. Purpose (Updated for SaaS Platform)

This is now the **main business interface**, not just marketing:
- **Customer acquisition**: Professional landing page with instant signup
- **User onboarding**: Seamless Firebase Auth + API key generation  
- **Product demonstration**: Interactive demo showcasing business value
- **Revenue generation**: Clear pricing tiers and upgrade paths

### Must Emphasize:
- Mathematical explainability for regulatory compliance
- Instant trial with no credit card required
- DORA, NIS2, EU AI Act compliance built-in
- AI insights as premium upgrade (not core requirement)

---

## 2. Tech Stack & Quality Bar

### Current Architecture:
- **Vue 3 + TypeScript** - Reactive, type-safe frontend
- **Vite** - Fast development and optimized production builds
- **Tailwind CSS** - Responsive, professional design system
- **Firebase Hosting** - Global CDN with custom domain support
- **Firebase Auth** - User authentication and management

### Quality Requirements:
- 🚀 **Performance**: Lighthouse score >90 (all categories)
- 📱 **Responsive**: Perfect on mobile, tablet, desktop
- 🌙 **Dark mode**: Support system preferences
- ♿ **Accessibility**: WCAG 2.1 AA compliance
- 🔒 **Security**: No secrets in frontend code, CSP headers

### Build Process:
```bash
# Development
npm run dev

# Production build (required before Firebase deploy)
npm run build

# Type checking
npm run type-check
```

---

## 3. Content & Messaging Strategy

### Business-First Approach:
- ✅ Lead with regulatory compliance value
- ✅ Show mathematical proof capability  
- ✅ Emphasize instant trial and ease of use
- ✅ Clear pricing with upgrade CTAs
- ❌ No technical implementation details
- ❌ No code examples or API endpoints
- ❌ No "developer-first" language

### Key Messages:
1. **"The Only Explainable Anomaly Detection Platform"**
2. **"DORA & NIS2 Compliant from Day One"**  
3. **"Get API Key in 30 Seconds - No Credit Card"**
4. **"Mathematical Proof for Every Decision"**

---

## 4. Component Architecture

### Core Components:
- `DemoComponent.vue` - Interactive anomaly detection showcase
- `SignupForm.vue` - Firebase Auth integration with instant API keys
- `PricingSection.vue` - Clear SaaS pricing tiers
- `ComplianceSection.vue` - Regulatory value proposition

### API Integration:
```javascript
// Use Firebase Functions endpoints only
/api/signup          // User registration
/api/analyze         // AI insights (premium)
/api/compliance      // Generate reports
/api/proxy/*         // Backend access
```

### Authentication Flow:
1. User signs up via Firebase Auth
2. Trigger Cloud Run backend to create tenant + API key
3. Display API key immediately with copy-to-clipboard
4. Guide user to make first API call

---

## 5. Cost Optimization Rules

### AI Usage (CRITICAL):
- ❌ **No AI calls in demo/default flows**
- ✅ **AI as premium upsell only**
- ✅ **Mathematical explanations are primary value**
- ✅ **"Upgrade for AI insights" messaging**

### Performance:
- Core demo must complete in <2 seconds
- No unnecessary API calls or external dependencies
- Optimize images and assets for fast loading
- Use lazy loading for non-critical components

---

## 6. Firebase Integration

### Hosting Configuration:
```json
{
  "rewrites": [
    {
      "source": "/api/**",
      "function": "corresponding Firebase Function"
    },
    {
      "source": "**",
      "destination": "/index.html"
    }
  ]
}
```

### Build & Deploy:
```bash
# Always build first
npm run build

# Deploy from root directory
cd .. && firebase deploy --only hosting
```

### Environment Variables:
- Use Firebase config for environment-specific settings
- No hard-coded API URLs or keys
- Support dev/staging/production environments

---

## 7. User Experience Requirements

### Signup Flow:
1. **Instant**: API key in <30 seconds
2. **Simple**: Email + company name only
3. **Clear**: Show what they get in trial
4. **Actionable**: Next steps after signup

### Demo Experience:  
1. **Fast**: Results in <2 seconds
2. **Clear**: Mathematical explanation prominent
3. **Upsell**: AI insights require upgrade
4. **Compliance**: Show audit-ready evidence

### Pricing Strategy:
- **Free Trial**: 14 days, 10K events, mathematical explanations
- **Pro**: $99/month, AI insights, advanced reports
- **Enterprise**: Custom pricing, white-glove support

---

## 8. Never Do This

❌ **Expose technical details** (cargo commands, API endpoints)
❌ **Make AI required** for core functionality  
❌ **Hard-code backend URLs** (use environment variables)
❌ **Break mobile responsiveness**
❌ **Add expensive external dependencies**
❌ **Slow down the signup or demo flows**

---

## 9. Success Metrics

### Technical:
- Build succeeds in <2 minutes
- Page load time <2 seconds
- Mobile Lighthouse score >90
- Zero console errors

### Business:
- Signup conversion >15%
- Demo completion >70%  
- API key usage >50% within 24h
- Upgrade consideration >10%

**Remember**: This is the primary customer interface for a real SaaS business. Every pixel should drive user acquisition and retention.
  - Strong Lighthouse scores (performance, accessibility, SEO).
  - Responsive design and dark mode support.
  - Clean, modular components in `src/`.
- Do not add heavy client‑side trackers or third‑party widgets without clear value and privacy review.

---

## 3. Content discipline

- When changing copy:
  - Align with the current roadmap and implementation reality as documented in `docs/ROADMAP_TO_LAUNCH.md` and `FINAL-STATUS.md`.
  - Make it clear which capabilities are available **now** vs. **planned**.
- When adding CTAs or forms:
  - Prefer simple, backend‑agnostic integrations (e.g., webhook endpoints provided separately).
  - Do not hardcode secrets or endpoints tied to a specific environment.

