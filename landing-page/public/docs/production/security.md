# Security Best Practices

Protect your data and your Driftlock account with these security best practices.

## API Key Management

Your API keys are the keys to your kingdom. Treat them with the same care as passwords.

### Secret vs. Public Keys

- **Secret Keys**: Full access, including creating streams and viewing data. **NEVER** expose these in client-side code (browsers, mobile apps).
- **Restricted Keys**: Scoped permissions (e.g., `write-only`). Use these for frontend or untrusted environments.

### Rotation

Periodically rotate your API keys, especially if you suspect a leak.
1. Generate a new key in the dashboard.
2. Update your applications to use the new key.
3. Revoke the old key.

## Data Privacy

Driftlock is designed to detect anomalies without needing to know *what* your data means.

### PII Redaction

**Do not send Personally Identifiable Information (PII)** like credit card numbers, social security numbers, or passwords in the event body.
- Hash sensitive fields before sending.
- Use internal IDs instead of names/emails.

```javascript
// BAD
driftlock.detect({
  body: { email: "user@example.com", credit_card: "4242..." }
});

// GOOD
driftlock.detect({
  body: { user_id: "u_123", payment_hash: "abc..." }
});
```

## Network Security

### TLS/SSL

All communication with the Driftlock API must be over HTTPS. Requests made over plain HTTP will be rejected.

### IP Whitelisting

(Enterprise Plan) You can restrict API access to specific IP addresses or CIDR blocks to ensure traffic only comes from your known servers.

## Compliance

### GDPR / CCPA

Driftlock acts as a Data Processor. We retain event data only for the duration of the retention window (default 30 days) to calculate baselines. You can request full data deletion at any time via support.

### SOC 2

Driftlock is currently undergoing SOC 2 Type I certification. Contact sales for more information.

## Vulnerability Reporting

If you discover a security vulnerability in Driftlock, please report it to [security@driftlock.io](mailto:security@driftlock.io). We offer a bug bounty program for responsible disclosure.
