# CI/CD Authentication Approach

## Overview
This project uses a **dev login endpoint** for authentication in CI/CD environments instead of requiring GitHub OAuth credentials. This simplifies testing and eliminates the need to manage OAuth secrets in CI.

## How It Works

### Dev Login Endpoint
The backend provides a `/auth/dev/login` endpoint (in `cinema.log.server.golang/internal/auth/handler.go`) that:
- Is only available in non-production environments (controlled by `ENVIRONMENT` variable)
- Automatically creates or retrieves a test user (GitHub ID: 0, username: "devuser", name: "Dev User")
- Returns JWT authentication cookies without requiring OAuth flow

### Playwright Test Setup
The Playwright auth setup (`tests/auth.setup.ts`) now:
1. Calls the `/auth/dev/login` endpoint directly via API request
2. Receives authentication cookies automatically
3. Saves the authentication state for reuse in all tests
4. No GitHub OAuth interaction required

## Benefits
✅ **Simpler**: No need to manage GitHub OAuth credentials in CI
✅ **Faster**: Bypasses external OAuth flow, reducing test time significantly
✅ **Reliable**: No dependency on GitHub OAuth service availability
✅ **Secure**: Dev endpoint is disabled in production via environment check

## Environment Configuration

### For CI/PR Environments
**No special configuration needed!** The tests will work as long as:
- `ENVIRONMENT` is NOT set to "production" (or not set at all)
- The backend is running and accessible

### For Production
The dev login endpoint is automatically disabled when `ENVIRONMENT=production`, ensuring security.

## Local Testing
To test the Playwright suite locally:

```bash
# Start the dev environment (Angular + Go backend + Postgres)
./run-dev.zsh

# Run Playwright tests
npx playwright test
```

The tests will automatically use the dev login endpoint since `ENVIRONMENT` is not set to production in local development.

## Migration from OAuth Approach
This replaces the previous approach that required:
- ❌ GitHub OAuth app configuration with callback URLs for each PR environment
- ❌ GitHub test account credentials (email, password, TOTP secret)
- ❌ Complex OAuth flow navigation in Playwright tests
- ❌ Managing Railway environment variables for callback URLs

## Code Changes Summary

### Backend (`cinema.log.server.golang/internal/auth/handler.go`)
- Kept existing `DevLogin()` handler (already existed)
- Removed environment-based OAuth callback URL routing (not needed)

### Tests (`tests/auth.setup.ts`)
- Replaced GitHub OAuth flow with direct dev login API call
- Removed TOTP library dependency
- Simplified authentication to ~40 lines from ~70 lines

### CI Configuration (`.github/workflows/pipeline.yml`)
- Removed GitHub OAuth secret environment variables:
  - `TEST_GITHUB_EMAIL`
  - `TEST_GITHUB_PASSWORD`
  - `TEST_GITHUB_TOTP_SECRET`

### Dependencies (`package.json`)
- Removed `totp-generator` (no longer needed)

## Security Considerations

The dev login endpoint is **safe** because:
1. It's only enabled when `ENVIRONMENT != "production"`
2. Railway production deployments should have `ENVIRONMENT=production` set
3. The endpoint returns 404 in production environments
4. Test user has no special privileges beyond a normal user account

## Troubleshooting

### "Dev login failed with status 404"
**Cause**: The backend has `ENVIRONMENT=production` set
**Solution**: Ensure PR/test environments do NOT have `ENVIRONMENT=production`

### Tests can't reach the backend
**Cause**: Railway deployment not ready or BASE_URL incorrect
**Solution**: 
- Check the workflow's "Wait for Railway PR deployment" step
- Verify `BASE_URL` matches the actual Railway PR environment URL

### "navbar-profile-link not visible"
**Cause**: Authentication cookies not being set correctly
**Solution**:
- Check backend logs to ensure dev login endpoint was called
- Verify cookies are being saved in `.auth/user.json`
- Check CORS settings allow credentials from the frontend URL

## Future Improvements
Consider adding:
- Multiple test users with different roles/permissions
- Cleanup script to reset test data between runs
- Dev login endpoint with customizable user attributes for specific test scenarios
