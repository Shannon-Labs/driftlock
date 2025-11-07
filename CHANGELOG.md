# Changelog

All notable changes to DriftLock will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial open source release of DriftLock
- Compression-based anomaly detection (CBAD) engine in Rust
- Go API server with OpenTelemetry integration
- React dashboard with real-time anomaly monitoring
- API key authentication for OSS deployments
- Docker Compose setup for easy local development
- Comprehensive test suite for all components

### Changed
- Removed Supabase dependency - now runs standalone
- Dashboard authentication changed from Supabase to API key-based
- Simplified deployment configuration
- Updated documentation for OSS usage

### Removed
- Cloudflare API Worker (Supabase-dependent)
- Hardcoded Supabase credentials
- Proprietary authentication flows

## [0.1.0] - 2024-11-06

### Added
- Standalone OSS release
- API key authentication system
- REST API for all dashboard functionality
- Simplified configuration management
- Comprehensive documentation

### Breaking Changes
- **Supabase Integration Removed**: DriftLock no longer requires Supabase for core functionality
- **Authentication Changed**: Dashboard now uses API keys instead of Supabase auth
- **Configuration Changes**: Environment variables updated to reflect standalone operation

### Migration Guide
For users upgrading from previous versions:

1. **API Key Setup**: Configure `DEFAULT_API_KEY` in your environment for dashboard access
2. **Environment Variables**: Remove Supabase-related variables unless using compliance features
3. **Dashboard Access**: Use API key to log in instead of email/password
4. **API Integration**: Update any direct Supabase queries to use REST API endpoints

### Technical Details
- API server handlers gracefully handle nil Supabase client
- All Supabase queries replaced with REST API calls
- Dashboard components updated to use API endpoints
- Test suite updated to work without Supabase
- Docker configuration simplified

[Unreleased]: https://github.com/shannon-labs/driftlock/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/shannon-labs/driftlock/releases/tag/v0.1.0
