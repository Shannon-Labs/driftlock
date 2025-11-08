# Release Notes

## Version 0.1.0 - Standalone OSS Release

### Breaking Changes

- **Supabase Dependency Removed**: DriftLock now runs standalone without requiring Supabase. The API server and dashboard work independently using API key authentication.
- **Authentication Changed**: The web dashboard now uses API key-based authentication instead of Supabase auth. Set `DEFAULT_API_KEY` (and optionally `DRIFTLOCK_DEV_API_KEY`) to control dashboard access.

### New Features

- **API Key Authentication**: Simple API key-based authentication for OSS deployments
- **Standalone Dashboard**: Dashboard works without Supabase, fetching data directly from the API server
- **Simplified Setup**: No external dependencies required beyond PostgreSQL

### Improvements

- **Removed Hardcoded Credentials**: All Supabase credentials removed from source code
- **Optional Supabase Integration**: Supabase can still be configured optionally for compliance features, but is not required for core functionality
- **Updated Documentation**: Installation and setup guides updated to reflect standalone operation

### Migration Guide

If you were previously using Supabase:

1. **API Server**: No changes needed - Supabase integration is optional and defaults to disabled
2. **Dashboard**: Update authentication to use API keys:
   - Set `DEFAULT_API_KEY` (and optionally `DRIFTLOCK_DEV_API_KEY`) environment variables on the API server
   - Use the API key to log into the dashboard
3. **Environment Variables**: Remove Supabase-related env vars unless you need compliance features

### Removed Components

- **Cloudflare API Worker**: Removed as it was entirely Supabase-dependent. Use the main API server directly for edge deployments.

### Technical Details

- API server handlers gracefully handle nil Supabase client
- Dashboard components updated to use REST API endpoints
- All Supabase queries replaced with API calls
- Test suite updated to work without Supabase
