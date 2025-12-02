# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Linear integration for project management (28 issues imported)
- AI cost controls with per-plan limits and budget tracking
- Demo endpoint (`POST /v1/demo/detect`) with rate limiting
- Usage dashboard with daily charts
- Playground demo mode for unauthenticated users
- API documentation at `/docs/user-guide/api/`
- `.editorconfig` for consistent formatting
- `SECURITY.md` with vulnerability reporting process

### Changed
- Upgraded OpenTelemetry collector dependencies to v1.44.0
- Improved error messages for billing endpoints

### Fixed
- Stripe webhook event handling for trial expiration
- API key caching with proper TTL
- Grace period logic for failed payments

### Security
- Added TruffleHog secret scanning to CI
- Added Trivy container vulnerability scanning
- Configured CodeQL SAST analysis

## [0.1.0] - 2025-12-02

### Added
- Initial release of Driftlock CBAD platform
- Compression-based anomaly detection using OpenZL
- OpenTelemetry Collector processor integration
- REST API for event detection and anomaly retrieval
- Stripe billing integration with trials and subscriptions
- Firebase Authentication support
- SendGrid email verification
- Cloud Run deployment configuration
- PostgreSQL storage with goose migrations
- Prometheus metrics endpoint (`/metrics`)
- Health check endpoints (`/healthz`)

### Algorithms
- Compression ratio analysis
- Normalized Compression Distance (NCD)
- Shannon entropy calculation
- Permutation testing with deterministic seeds (ChaCha20)

### Compliance
- DORA evidence bundle generation
- NIS2 incident reporting support
- AI Act runtime monitoring capabilities

---

[Unreleased]: https://github.com/shannon-labs/driftlock/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/shannon-labs/driftlock/releases/tag/v0.1.0
