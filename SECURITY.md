# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

**Please do NOT create public GitHub issues for security vulnerabilities.**

To report a security issue, please email **security@driftlock.io** with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Affected versions
4. Any potential mitigations you've identified

### What to Expect

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 5 business days
- **Resolution Timeline**: Depends on severity
  - Critical: 24-48 hours
  - High: 7 days
  - Medium: 30 days
  - Low: 90 days

### Disclosure Policy

We follow coordinated disclosure:
- We will work with you to understand and resolve the issue
- We will credit you in our security advisory (unless you prefer anonymity)
- We ask that you give us reasonable time to address the issue before public disclosure

## Security Measures

### Infrastructure
- All secrets stored in Google Cloud Secret Manager
- TLS 1.3 enforced for all connections
- Cloud Run with VPC isolation
- Regular dependency audits via Dependabot

### Application
- Firebase Authentication for user sessions
- API key hashing with bcrypt
- Rate limiting on all public endpoints
- Input validation on all user data
- SQL injection prevention via parameterized queries

### Monitoring
- CodeQL SAST scanning on every PR
- Trivy container vulnerability scanning
- TruffleHog secret scanning
- Automated security workflow weekly

## Security Contacts

- **Primary**: security@driftlock.io
- **PGP Key**: Available upon request
