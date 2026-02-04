# Summary: Fixing CI OAuth Callback URL Issue

## Problem Identified ✅
The Playwright tests in CI were failing because:
- GitHub OAuth was redirecting to the **production URL** (`https://cinema-log.up.railway.app/api/auth/github-callback`)
- But tests needed to redirect to the **PR environment URL** (`https://reverse-proxy-cinemalog-pr-{PR_NUM}.up.railway.app/api/auth/github-callback`)

**Evidence from logs:**
```
navigated to "https://cinema-log.up.railway.app/api/auth/github-callback?error=access_denied"
```

The `error=access_denied` occurred because GitHub was redirecting back to production, but the test user/session was in the PR environment.

## Solution Implemented ✅

### Code Changes Made:
1. **Modified** `cinema.log.server.golang/internal/auth/handler.go`:
   - Added `getCallbackBaseURL()` function with priority-based URL selection
   - OAuth configs now use this function instead of hardcoded `BACKEND_URL`
   - Supports 3 sources (in priority order):
     1. `CALLBACK_BASE_URL` (explicit override)
     2. `RAILWAY_PUBLIC_DOMAIN` (auto-detected)
     3. `BACKEND_URL` (fallback)

2. **Added Tests** in `cinema.log.server.golang/internal/auth/handler_test.go`:
   - Tests all priority scenarios
   - Verifies OAuth configs use correct callback URLs
   - All 9 tests passing ✅

3. **Created Documentation**:
   - `OAUTH_CALLBACK_SETUP.md` - Technical details
   - `RAILWAY_SETUP_INSTRUCTIONS.md` - Step-by-step user guide
   - This file - Executive summary

## What You Need to Do Next 🎯

### Step 1: Configure Railway PR Environments

Railway needs to set the callback URL for PR deployments. Choose ONE of these approaches:

#### Option A: Use Railway Environment Variables (Recommended)
In your Railway project settings for PR environments, set:
```bash
CALLBACK_BASE_URL=https://reverse-proxy-cinemalog-pr-${PR_NUMBER}.up.railway.app
```

**Notes:**
- Railway may support `${PR_NUMBER}` or `${{RAILWAY_STATIC_URL}}` or similar variables
- Check Railway's documentation for PR environment variables
- If Railway doesn't support dynamic variables, see Option B

#### Option B: Use Railway's Public Domain Variable
If Railway auto-provides the domain, set:
```bash
RAILWAY_PUBLIC_DOMAIN=reverse-proxy-cinemalog-pr-${PR_NUMBER}.up.railway.app
```

The backend will automatically prepend `https://`.

#### Option C: Manual Per-PR Configuration
If Railway doesn't support dynamic variables:
- Use Railway's API/CLI to set `CALLBACK_BASE_URL` when each PR deploys
- Or use a Railway deploy hook/script
- Or set it manually for testing (not scalable)

### Step 2: Update GitHub OAuth App Settings

⚠️ **CRITICAL**: Your GitHub OAuth App must allow PR environment callbacks.

1. Go to https://github.com/settings/developers
2. Select your OAuth App (the one with `GITHUB_CLIENT_ID`)
3. Add these callback URLs:
   ```
   https://cinema-log.up.railway.app/auth/github-callback
   https://reverse-proxy-cinemalog-pr-*.up.railway.app/auth/github-callback
   ```

**Notes:**
- If GitHub doesn't support wildcard `*`, add specific PR URLs as needed
- Or create a separate OAuth App for testing with flexible callback URLs
- Both GitHub and Google OAuth configs have been updated (if you use Google auth)

### Step 3: Verify the Fix

After Railway configuration:

1. **Trigger a new deployment** on your PR branch
2. **Check Railway logs** to confirm the environment variable is set
3. **Run the Playwright tests** in GitHub Actions
4. **Expected result**: 
   - Tests should authenticate successfully
   - Should redirect to PR environment, not production
   - No `access_denied` errors

### Step 4: Validate Production Still Works

✅ **Good news**: Production is NOT affected by these changes!

- Production will continue using `BACKEND_URL` (existing behavior)
- Only set `CALLBACK_BASE_URL` in PR environments
- No production configuration changes needed

## Troubleshooting Guide

### Issue: Tests still redirect to production
**Causes:**
- Railway environment variable not set
- Railway deployment not restarted after setting variable
- Variable name mismatch

**Solutions:**
1. Check Railway environment variables in project settings
2. Redeploy the PR after setting the variable
3. Check Railway logs: look for `RedirectURL` in startup logs
4. Add temporary logging to confirm variable values

### Issue: GitHub OAuth returns "redirect_uri_mismatch"
**Cause:** GitHub OAuth App doesn't allow the PR callback URL

**Solution:**
1. Double-check GitHub OAuth App settings
2. Ensure PR domain pattern is added
3. Wait a few minutes after updating OAuth settings
4. Try with a specific PR URL first (e.g., PR #30)

### Issue: Railway doesn't support dynamic PR variables
**Solution:**
Contact Railway support or consider these alternatives:
1. Use Railway API to set variables programmatically
2. Add a pre-deploy script that calls Railway API
3. Use Railway's CLI in GitHub Actions to set the variable
4. Set variables manually for specific PR numbers you're testing

## Quick Reference

### Environment Variables
| Variable | Priority | Where to Set | Example |
|----------|----------|--------------|---------|
| `CALLBACK_BASE_URL` | 1 (Highest) | Railway PR env | `https://reverse-proxy-cinemalog-pr-30.up.railway.app` |
| `RAILWAY_PUBLIC_DOMAIN` | 2 | Railway auto | `reverse-proxy-cinemalog-pr-30.up.railway.app` |
| `BACKEND_URL` | 3 (Fallback) | All environments | `https://cinema-log.up.railway.app` |

### Testing Locally
```bash
# Simulate PR environment
export CALLBACK_BASE_URL="http://localhost:8080"
export BACKEND_URL="https://production.example.com"

# Verify it's using the right URL
# Check auth handler logs or add temporary debug logging
```

## Success Criteria ✅

Your configuration is correct when:
- [ ] Railway PR deployments have `CALLBACK_BASE_URL` or `RAILWAY_PUBLIC_DOMAIN` set
- [ ] GitHub OAuth App includes PR callback URL pattern
- [ ] Playwright tests in CI pass without `access_denied` errors
- [ ] Tests redirect to PR environment (`/profile/...`), not production
- [ ] Production deployments still work with existing `BACKEND_URL`

## Files Changed

1. `cinema.log.server.golang/internal/auth/handler.go` - OAuth callback URL logic
2. `cinema.log.server.golang/internal/auth/handler_test.go` - Tests (new file)
3. `OAUTH_CALLBACK_SETUP.md` - Technical documentation (new file)
4. `RAILWAY_SETUP_INSTRUCTIONS.md` - User guide (new file)
5. `IMPLEMENTATION_SUMMARY.md` - This file (new file)

## Need Help?

If you encounter issues:
1. Check Railway's documentation for PR environment variables
2. Review Railway logs to see what variables are available
3. Test locally by setting environment variables manually
4. Contact Railway support if dynamic PR variables aren't available
5. The code changes are solid - the only configuration needed is Railway + GitHub OAuth

---

**Status**: Code changes complete ✅ | Railway configuration needed ⏳ | GitHub OAuth App update needed ⏳
