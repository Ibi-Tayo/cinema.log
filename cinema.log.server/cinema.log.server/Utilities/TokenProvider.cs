using System.Security.Claims;
using System.Text;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Utilities;

public class TokenProvider(IConfiguration config)
{
    public string GenerateToken(User user, DateTime expiry, string? identity)
    {
        var data = Encoding.UTF8.GetBytes(config["TokenSecret"] ?? throw new InvalidOperationException());
        var securityKey = new Microsoft.IdentityModel.Tokens.SymmetricSecurityKey(data);

        var claims = new Dictionary<string, object>
        {
            [ClaimTypes.Name] = user.Name,
            [ClaimTypes.UserData] = user.UserId,
            ["IdentityPartial"] = identity ?? string.Empty,
        };
        var descriptor = new Microsoft.IdentityModel.Tokens.SecurityTokenDescriptor
        {
            Issuer = "Cinema.Log.Server",
            Audience = "Cinema.Log.Client",
            Claims = claims,
            NotBefore = DateTime.UtcNow,
            Expires = expiry,
            SigningCredentials = new Microsoft.IdentityModel.Tokens.SigningCredentials(securityKey, 
                Microsoft.IdentityModel.Tokens.SecurityAlgorithms.HmacSha256Signature)
        };

        var handler = new Microsoft.IdentityModel.JsonWebTokens.JsonWebTokenHandler
            {
                SetDefaultTimesOnTokenCreation = false
            };
        var tokenString = handler.CreateToken(descriptor);
        return tokenString;
    }
}