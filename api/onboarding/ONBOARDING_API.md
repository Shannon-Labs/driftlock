# Tenant Onboarding API Specification

## Overview
Public API endpoints for user signup, trial activation, and self-service tenant management.

## Endpoints

### 1. POST /v1/onboard/signup
Create a new tenant with trial account.

**Request:**
```json
{
  "email": "user@example.com",
  "company_name": "Acme Corp",
  "plan": "trial",
  "source": "landing-page"  // optional: analytics tracking
}
```

**Response (201):**
```json
{
  "success": true,
  "tenant": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Acme Corp",
    "slug": "acme-corp",
    "plan": "trial",
    "status": "pending_verification",
    "created_at": "2025-01-18T12:00:00Z"
  },
  "api_key": "dlk_550e8400-e29b-41d4-a716-446655440000.sk_live_xxxxxxxx",
  "trial_end": "2025-02-01T12:00:00Z",
  "message": "Check your email to activate your account"
}
```

**Error Responses:**
- `429 Too Many Requests` - Rate limited
- `400 Bad Request` - Invalid email or company name
- `409 Conflict` - Email already registered

**Rate Limiting:** 5 requests per IP per hour

### 2. GET /v1/onboard/verify
Verify email address and activate tenant.

**Query Parameters:**
- `token` - Verification token from email

**Response (302 Redirect):**
Redirects to: `https://driftlock.net/verified?tenant={tenant_id}`

### 3. POST /v1/onboard/resend-verification
Resend verification email.

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Rate Limiting:** 3 requests per email per day

## Implementation Status
- [x] Database schema supports onboarding
- [x] CLI tenant creation exists (needs API wrapper)
- [ ] Public onboarding endpoint (to implement)
- [ ] Email verification flow (to implement)
- [ ] Rate limiting per IP (to implement)
- [ ] Welcome email templates (to create)

## Database Changes Needed
Add to tenants table:
```sql
ALTER TABLE tenants 
ADD COLUMN email TEXT,
ADD COLUMN verification_token TEXT,
ADD COLUMN verified_at TIMESTAMPTZ,
ADD COLUMN trial_ends_at TIMESTAMPTZ,
ADD COLUMN signup_ip INET;
```

## Security Considerations
- Use Argon2id for token hashing
- Log all signup attempts
- Implement CAPTCHA after 3 failed attempts
- Check disposable email domains
- Rate limit by IP and email
- Store IP for fraud detection (GDPR compliant)
