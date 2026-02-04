# Railway Configuration Instructions for OAuth Callback URLs

## Quick Summary
To fix the failing Playwright tests in CI, you need to configure Railway so that PR environments use the correct OAuth callback URL.

## What Was Changed in the Code

The Go backend (`cinema.log.server.golang/internal/auth/handler.go`) now checks for callback URLs in this priority order:

1. `CALLBACK_BASE_URL` environment variable (explicit override)
2. `RAILWAY_PUBLIC_DOMAIN` environment variable (Railway auto-detected)
3. `BACKEND_URL` environment variable (existing fallback)

## Required Railway Configuration

You need to configure the **PR environment** in Railway to set the callback URL. Here are your options:

### Option 1: Set CALLBACK_BASE_URL for Each PR (Recommended)

In Railway's PR environment template/settings, add this environment variable:

```
CALLBACK_BASE_URL=https://reverse-proxy-cinemalog-pr-${{PR_NUMBER}}.up.railway.app
```

**Note**: The exact syntax for injecting PR number may vary in Railway. Common patterns:
- `${{PR_NUMBER}}` 
- `${RAILWAY_PR_NUMBER}`
- Or it might need to be set via Railway's API/CLI per PR

### Option 2: Use Railway's Public Domain Variable

If Railway auto-provides the public domain in a variable, you can set:

```
RAILWAY_PUBLIC_DOMAIN=reverse-proxy-cinemalog-pr-${{PR_NUMBER}}.up.railway.app
```

The code will prepend `https://` automatically.

### Option 3: Script-based Configuration

If Railway doesn't support dynamic PR number injection, you may need to:
1. Use Railway's API to set `CALLBACK_BASE_URL` for each PR deployment
2. Or use a Railway deploy hook/script to set the variable dynamically

## Required GitHub OAuth App Configuration

⚠️ **IMPORTANT**: You must also update your GitHub OAuth App to allow callbacks from PR domains.

1. Go to: https://github.com/settings/developers
2. Select your OAuth App
3. In "Authorization callback URL", add:
   ```
   https://reverse-proxy-cinemalog-pr-*.up.railway.app/auth/github-callback
   ```
   
   Note: GitHub may require you to add each PR URL individually if wildcards aren't supported. In that case, you can:
   - Add multiple specific PR URLs as needed, OR
   - Create a separate OAuth App for testing with wildcard/flexible callback URLs

## Verification Steps

After configuring Railway:

1. **Trigger a new PR build** or redeploy the current PR
2. **Check the backend logs** in Railway to verify it's using the correct callback URL
3. **Run the Playwright tests** - they should now pass
4. **Verify the OAuth flow**:
   - Tests should authenticate successfully
   - No "access_denied" errors
   - Should redirect to the PR environment, not production

## Testing the Configuration

You can verify the environment variable is being used by:

1. SSH into your Railway deployment (if possible)
2. Run: `echo $CALLBACK_BASE_URL` or `echo $RAILWAY_PUBLIC_DOMAIN`
3. Check the backend logs when it starts - you may want to add a temporary log statement

## Troubleshooting

### Problem: Tests still redirect to production
**Solutions**:
- Verify `CALLBACK_BASE_URL` is set in Railway PR environment
- Check Railway logs to confirm the environment variable is available
- Ensure the PR has been redeployed after setting the variable

### Problem: GitHub OAuth returns "redirect_uri_mismatch"
**Solution**: 
- Add the PR environment callback URL to your GitHub OAuth App settings
- Pattern: `https://reverse-proxy-cinemalog-pr-{NUMBER}.up.railway.app/auth/github-callback`

### Problem: Railway doesn't support dynamic PR number
**Solution**: 
- Use Railway's CLI or API to set the variable programmatically
- Or create a GitHub Action that calls Railway's API to set the variable before tests run

## Production Deployment

✅ **Production is NOT affected**. The code maintains backward compatibility:
- Production will continue using `BACKEND_URL` 
- Only set `CALLBACK_BASE_URL` in PR environments, not production

## Questions?

If you encounter issues with Railway's environment variable configuration:
1. Check Railway's documentation for PR environment variables
2. Contact Railway support for dynamic PR environment configuration
3. Consider alternative deployment strategies if Railway limitations persist

---

For complete technical details, see: [OAUTH_CALLBACK_SETUP.md](./OAUTH_CALLBACK_SETUP.md)
