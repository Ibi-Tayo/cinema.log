using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Utilities;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Distributed;
using Microsoft.IdentityModel.Tokens;
using Octokit;
using SecurityTokenValidationException = Microsoft.IdentityModel.Tokens.SecurityTokenValidationException;
using TokenProvider = cinema.log.server.Utilities.TokenProvider;
using User = cinema.log.server.Models.Entities.User;

namespace cinema.log.server.Controllers;

// TODO: Note that for now the auth controller is tightly coupled to the localhost implementation of the client
// TODO: Make sure to change base url depending on if we are running prod or dev
[ApiController]
[Route("auth")]
public class AuthController : ControllerBase
{
    private readonly string _clientId;
    private readonly string _clientSecret;
    private readonly GitHubClient _client;
    private readonly IDistributedCache _distributedCache;
    private readonly IUserRepository _userRepository;
    private readonly IConfiguration _config;

    public AuthController(IConfiguration config, IDistributedCache distributedCache, IUserRepository userRepository)
    {
        _config = config;
        _clientId = _config["GithubClientId"] ?? throw new InvalidOperationException();
        _clientSecret = _config["GithubClientSecret"] ?? throw new InvalidOperationException();
        _client = new GitHubClient(new ProductHeaderValue("cinema-log"));
        _distributedCache = distributedCache;
        _userRepository = userRepository;
    }


    [HttpGet]
    [Route("github-login")]
    public async Task<Response<string>> Login()
    {
        var csrf = Guid.NewGuid().ToString();
        await _distributedCache.SetStringAsync("github-oauth-state", csrf,
            new DistributedCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10) });

        var request = new OauthLoginRequest(_clientId)
        {
            State = csrf
        };

        // User gets navigated to GitHub for login
        var oauthLoginUrl = _client.Oauth.GetGitHubLoginUrl(request);
        
        return Response<string>.BuildResponse(200, "Redirect", oauthLoginUrl.AbsoluteUri);
    }

    [HttpGet]
    [Route("logout")]
    public async Task<Response<string>> Logout()
    {
        // This will clear by obtaining and setting as expired 
        foreach (var cookie in Request.Cookies)
        {
            Response.Cookies.Delete(cookie.Key);
        }

        return Response<string>.BuildResponse(200, "Redirect", "http://localhost:4200/login");
    }


    // This is the redirect uri specified with GitHub
    [HttpGet]
    [Route("github-callback")]
    public async Task<IActionResult> GithubCallback(string code, string state)
    {
        try
        {
            if (string.IsNullOrEmpty(code) || string.IsNullOrEmpty(state)) throw new AuthorizationException();
            // Interface with GitHub to get an access token from them
            var expectedState = await _distributedCache.GetStringAsync("github-oauth-state");
            if (!string.Equals(expectedState, state)) throw new AuthorizationException();
            var request = new OauthTokenRequest(_clientId, _clientSecret, code);
            var token = await _client.Oauth.CreateAccessToken(request);
            var accessToken = token.AccessToken;

            if (accessToken == null) throw new AuthorizationException();
            _client.Credentials = new Credentials(accessToken);

            var user = await _client.User.Current();
            if (user == null) throw new AuthorizationException();
            // Use the user above to find the user in our database
            var newOrCurrentUser = await _userRepository
                .GetOrCreateUserFromGithubId(user.Id, user.Name, user.Login, user.AvatarUrl);

            // Generate the JWT tokens 
            var jwtAccessToken = new TokenProvider(_config)
                .GenerateToken(newOrCurrentUser, DateTime.UtcNow.AddMinutes(10), null);

            var ip = Request.HttpContext.Connection.RemoteIpAddress ?? Request.HttpContext.Connection.LocalIpAddress;
            var client = Request.Headers.UserAgent.ToString();
            var jwtRefreshToken = new TokenProvider(_config)
                .GenerateToken(newOrCurrentUser, DateTime.UtcNow.AddDays(7), string.Join(',', ip, client));

            // Redirect the user to the homepage with the token 
            UpdateResponseCookies(jwtAccessToken, jwtRefreshToken, newOrCurrentUser);
            return new RedirectResult($"http://localhost:4200/profile/{newOrCurrentUser.Username}");
        }
        catch (AuthorizationException e)
        {
            // Any null or invalids throw this exception and redirects the user to the GitHub OAuth login page.
            var request = new OauthLoginRequest(_clientId);
            var oauthLoginUrl = _client.Oauth.GetGitHubLoginUrl(request);
            return new RedirectResult(oauthLoginUrl.AbsoluteUri);
            
        }
        catch (Exception e)
        {
            return new RedirectResult("http://localhost:4200/login");
        }
    }

    [HttpGet]
    [Route("refresh-token")]
    public async Task<IActionResult> RequestRefreshToken()
    {
        try
        {
            // Check the request cookies and get the refresh token.
            var requestRefreshToken = Request.Cookies["cinema-log-refresh-token"];
            if (string.IsNullOrEmpty(requestRefreshToken)) throw new AuthorizationException();
            // Validate the refresh token
            var handler = new JwtSecurityTokenHandler();
            var validationResult = await handler
                .ValidateTokenAsync(requestRefreshToken, ProvideTokenValidationParameters(_config));
            if (!validationResult.IsValid) throw new SecurityTokenValidationException();

            var userId = validationResult.Claims[ClaimTypes.UserData] as string ??
                         throw new SecurityTokenValidationException();
            var user = await _userRepository.GetUserById(new Guid(userId)) ?? throw new AuthorizationException();

            // Generate the JWT tokens 
            var jwtAccessToken = new TokenProvider(_config)
                .GenerateToken(user, DateTime.UtcNow.AddMinutes(10), null);
            
            var ip = Request.HttpContext.Connection.RemoteIpAddress ?? Request.HttpContext.Connection.LocalIpAddress;
            var client = Request.Headers.UserAgent.ToString();
            var jwtRefreshToken = new TokenProvider(_config)
                .GenerateToken(user, DateTime.UtcNow.AddDays(7), string.Join(',', ip, client));

            // Add cookies and just return an ok
            UpdateResponseCookies(jwtAccessToken, jwtRefreshToken, user);
            return Ok();
        }
        catch (AuthorizationException e)
        {
            return Unauthorized();
        }
        catch (SecurityTokenValidationException e)
        {
            return Unauthorized();
        }
    }

    private void UpdateResponseCookies(string jwtAccessToken, string jwtRefreshToken, User user)
    {
        Response.Cookies.Append("cinema-log-access-token", jwtAccessToken,
            ProvideCookieOptions(DateTime.UtcNow.AddMinutes(10)));
        Response.Cookies.Append("cinema-log-refresh-token", jwtRefreshToken,
            ProvideCookieOptions(DateTime.UtcNow.AddDays(7)));
        Response.Cookies.Append("userId", user.UserId.ToString(), 
            new CookieOptions() {Expires = DateTime.UtcNow.AddDays(7) });
    }
    
    private CookieOptions ProvideCookieOptions(DateTime expiry)
    {
        var cookieOptions = new CookieOptions()
        {
            IsEssential = true,
            Expires = expiry,
            Secure = true,
            HttpOnly = true,
            SameSite = SameSiteMode.None
        };
        return cookieOptions;
    }

    private TokenValidationParameters ProvideTokenValidationParameters(IConfiguration config)
    {
        return new TokenValidationParameters
        {
            IssuerSigningKey =
                new SymmetricSecurityKey(
                    Encoding.UTF8.GetBytes(config["TokenSecret"] ?? throw new InvalidOperationException())),
            ValidIssuer = "Cinema.Log.Server",
            ValidAudience = "Cinema.Log.Client",
            ClockSkew = TimeSpan.Zero
        };
    }
}