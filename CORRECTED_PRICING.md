# DriftLock Corrected Pricing Configuration

## Actual Stripe Pricing (Website Configuration)

| Plan | Monthly Price | Product ID | Price ID | Included API Calls | Overage Rate |
|------|--------------|------------|----------|-------------------|--------------|
| **Pro** | $49/month | `prod_TJKXbWnB3ExnqJ` | `price_1SMhsZL4rhSbUSqA51lWvPlQ` | **50,000 calls** | **$0.001/call** |
| **Enterprise** | $249/month | `prod_TJKXEFXBjkcsAB` | `price_1SMhshL4rhSbUSqAyHfhWUSQ` | **500,000 calls** | **$0.0005/call** |

### What Was Incorrect

❌ **My Documentation Said:**
- Developer (Free): 1,000 detections
- Standard ($49): 50,000 detections
- Growth ($249): 500,000 detections

✅ **CORRECT Website Pricing:**
- Pro ($49): 50,000 API calls + $0.001/call overage
- Enterprise ($249): 500,000 API calls + $0.0005/call overage

### Key Differences

1. **Only 2 plans** (not 3)
2. **No free tier** (website doesn't have free option)
3. **Overage rates are critical** (much cheaper per call with higher tier)
4. **Metered usage** - only pay for actual API calls (not anomaly detections)

### Overage Examples

**Pro Plan ($49/month):**
- Included: 50,000 calls
- At 50,001 calls: $49 + (1 × $0.001) = $49.001
- At 75,000 calls: $49 + (25,000 × $0.001) = $74.00

**Enterprise Plan ($249/month):**
- Included: 500,000 calls
- At 600,000 calls: $249 + (100,000 × $0.0005) = $299.00
- At 750,000 calls: $249 + (250,000 × $0.0005) = $374.00

### Changes Needed

1. **Database**: Update plan_price_map table with correct Product/Price IDs
2. **Stripe Webhook**: Update edge function to use correct Price IDs
3. **Frontend**: Update billing components to reflect 2 tiers only
4. **Documentation**: Fix all pricing references

---

**Action Required:** Update all configuration files to match this corrected pricing.
