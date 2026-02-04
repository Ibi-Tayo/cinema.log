# OAuth Callback URL Configuration for CI/PR Environments

## Problem
GitHub OAuth callbacks need to redirect to the correct environment URL. In production, this is the main domain, but in PR environments (Railway PR deployments), this needs to be the PR-specific URL.

## Solution
The backend now supports environment-based callback URL routing with the following priority:

1. **CALLBACK_BASE_URL** - Explicit override (highest priority)
2. **RAILWAY_PUBLIC_DOMAIN** - Auto-detected from Railway
3. **BACKEND_URL** - Default fallback

## Railway Configuration Required

### For PR Environments
Railway needs to set one of the following environment variables for PR deployments:

**Option 1: Use Railway's Public Domain (Recommended)**
- Railway should automatically provide `RAILWAY_PUBLIC_DOMAIN` 
- The code will auto-detect this and use `https://${RAILWAY_PUBLIC_DOMAIN}`
- **No manual configuration needed if Railway sets this variable**

**Option 2: Manual Override**
Set `CALLBACK_BASE_URL` environment variable in Railway PR environment settings:
```
CALLBACK_BASE_URL=https://reverse-proxy-cinemalog-pr-{PR_NUMBER}.up.railway.app
```

### For Production
- Keep existing `BACKEND_URL` environment variable
- Do NOT set `CALLBACK_BASE_URL` or `RAILWAY_PUBLIC_DOMAIN` (unless Railway auto-sets it)
- Production will use `BACKEND_URL` as before

## GitHub OAuth App Configuration

Ensure your GitHub OAuth App has the PR environment callback URLs added to the allowed callback URLs:

1. Go to GitHub Settings > Developer settings > OAuth Apps > Your App
2. Add the following callback URLs:
   - Production: `https://cinema-log.up.railway.app/auth/github-callback`
   - PR Template: `https://reverse-proxy-cinemalog-pr-*.up.railway.app/auth/github-callback`
   
Note: GitHub OAuth allows wildcards or you can add specific PR URLs as needed.

## Testing Locally

To test the callback URL logic locally:

```bash
# Test with explicit override
export CALLBACK_BASE_URL="http://localhost:8080"
export BACKEND_URL="https://production.example.com"
# OAuth will use: http://localhost:8080/auth/github-callback

# Test with Railway domain simulation
export RAILWAY_PUBLIC_DOMAIN="my-app-pr-30.up.railway.app"
export BACKEND_URL="https://production.example.com"
# OAuth will use: https://my-app-pr-30.up.railway.app/auth/github-callback

# Test with fallback to BACKEND_URL
unset CALLBACK_BASE_URL
unset RAILWAY_PUBLIC_DOMAIN
export BACKEND_URL="https://production.example.com"
# OAuth will use: https://production.example.com/auth/github-callback
```

## Implementation Details

The callback URL logic is in `cinema.log.server.golang/internal/auth/handler.go`:

```go
func getCallbackBaseURL() string {
    // Priority 1: Explicit override
    if callbackBaseURL := os.Getenv("CALLBACK_BASE_URL"); callbackBaseURL != "" {
        return callbackBaseURL
    }
    
    // Priority 2: Railway auto-provided domain
    if railwayDomain := os.Getenv("RAILWAY_PUBLIC_DOMAIN"); railwayDomain != "" {
        return "https://" + railwayDomain
    }
    
    // Priority 3: Fallback to original BACKEND_URL
    return BackendURL
}
```

This function is used in both GitHub and Google OAuth configurations.

## Troubleshooting

### Playwright Tests Failing with "access_denied"
**Symptom**: Tests redirect to production URL instead of PR environment
**Solution**: Ensure Railway PR environment has either:
- `RAILWAY_PUBLIC_DOMAIN` set automatically, OR
- `CALLBACK_BASE_URL` set manually to the PR environment URL

### OAuth Redirects to Wrong Domain
**Check**:
1. Verify the environment variables are set correctly in Railway
2. Check GitHub OAuth App callback URLs include the PR domain pattern
3. Review backend logs to see which callback URL is being used

### Local Development Issues
**Solution**: 
- For local dev, ensure `BACKEND_URL=http://localhost:8080`
- Clear `CALLBACK_BASE_URL` and `RAILWAY_PUBLIC_DOMAIN` if testing locally
- Make sure GitHub OAuth app has `http://localhost:8080/auth/github-callback` in allowed URLs

## Railway Deployment Checklist

- [ ] Verify `RAILWAY_PUBLIC_DOMAIN` is available in Railway PR environments
- [ ] If not, set `CALLBACK_BASE_URL` in Railway PR environment template
- [ ] Update GitHub OAuth App callback URLs to include PR domain pattern
- [ ] Test Playwright E2E tests pass after deployment
- [ ] Verify production environment still uses `BACKEND_URL` correctly
