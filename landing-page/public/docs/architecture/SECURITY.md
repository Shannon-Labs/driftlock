# Security Policy

## Supported Versions

| Version | Supported          |
|---------|-------------------|
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

The Shannon Labs security team takes security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

If you discover a security vulnerability in DriftLock, please report it to us privately before disclosing it publicly.

### How to Report

**Email**: hunter@shannonlabs.dev

Please include the following information in your report:

- Type of vulnerability
- Affected versions
- Steps to reproduce the vulnerability
- Potential impact of the vulnerability
- Any proof-of-concept or exploit code (if available)

### What to Expect

- **Confirmation**: We will acknowledge receipt of your report within 24 hours
- **Assessment**: We will investigate the vulnerability and assess its impact
- **Resolution**: We will work on a fix and provide a timeline for release
- **Coordination**: We will coordinate disclosure with you if desired
- **Credit**: With your permission, we will credit you in the security advisory

### Security Features

DriftLock includes several security features:

- **Input validation**: All inputs are validated before processing
- **Authentication**: API endpoints support JWT-based authentication
- **Authorization**: Role-based access control for sensitive operations
- **Encryption**: Data encryption in transit and at rest
- **Audit logging**: Comprehensive logging of security-relevant events

### Security Best Practices

When deploying DriftLock, we recommend:

1. **Network Security**
   - Use HTTPS/TLS for all communications
   - Deploy behind a firewall or VPN
   - Implement network segmentation

2. **Access Control**
   - Use strong, unique passwords
   - Implement multi-factor authentication where possible
   - Follow principle of least privilege

3. **Data Protection**
   - Encrypt sensitive data at rest
   - Regularly backup data and test restores
   - Implement data retention policies

4. **Monitoring**
   - Monitor security logs for suspicious activity
   - Set up alerts for security events
   - Regularly review access logs

## Security Advisories

We will publish security advisories for resolved vulnerabilities in our [GitHub Security Advisories](https://github.com/Shannon-Labs/driftlock/security/advisories).

## Security Team

The Shannon Labs security team can be reached at:
- **Security**: hunter@shannonlabs.dev
- **General Security Questions**: hunter@shannonlabs.dev

### Encryption Keys and Secrets

- Never commit API keys, passwords, or other secrets to the repository
- Use environment variables for configuration
- Rotate secrets regularly
- Use a secrets management system in production

### Frontend API Key Management

Client-side applications, like the `landing-page`, require special handling for API keys.

- **NEVER** embed API keys directly in the frontend code or commit them to version control in files like `.env`.
- The `VITE_FIREBASE_API_KEY` for the frontend is stored securely in **Google Secret Manager**.
- A **Firebase Cloud Function** (`getFirebaseConfig`) acts as a secure proxy. It is the only part of the system with permission to access the API key secret.
- The frontend application calls this Cloud Function at runtime to fetch the Firebase configuration.
- This pattern ensures that the API key is never exposed directly to the browser or public repositories.

## Dependency Security

We regularly scan our dependencies for known vulnerabilities using:
- GitHub Dependabot
- Snyk
- OWASP Dependency Check

If you discover a vulnerable dependency, please report it following the same process as other security vulnerabilities.

## Responsible Disclosure Policy

We follow a responsible disclosure policy:

- We will respond to security reports within 24 hours
- We will provide regular updates on our progress
- We will request a reasonable amount of time to fix the vulnerability
- We will credit you for your discovery (with your permission)
- We will not pursue legal action against researchers who follow this policy

## Security Recognition

We recognize and appreciate the work of security researchers who help us keep DriftLock secure. With your permission, we will:

- Add your name to our Security Hall of Fame
- Send you Shannon Labs swag
- Provide a financial reward for qualifying vulnerabilities (see our bug bounty program)

Thank you for helping keep DriftLock secure!