using System.Text;
using System.Text.Json;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Services;

public class SoundtrackService : ISoundtrackService
{
    private readonly CinemaLogContext _context;
    private readonly HttpClient _httpClient;

    private readonly IConfiguration _configuration;

    // use film Id to get data from film table,
    // send request to spotify to get details
    public SoundtrackService(CinemaLogContext context, HttpClient httpClient, IConfiguration configuration)
    {
        _context = context;
        _httpClient = httpClient;
        _configuration = configuration;
    }

    public Response<FilmSoundtrackDto> GetSoundtrackByFilmId(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // use user id to get all liked tracks from liked track table
    public Response<List<LikedTrackDto>> GetLikedTracksByUserId(Guid userId)
    {
        throw new NotImplementedException();
    }

    // use dto to add new liked track to liked track table
    // (user needs to make sure they have the film soundtrack rating id)
    // (they'd get that by calling GetSoundtrackByFilmId)
    public Response<LikedTrackDto> SetLikedTrack(LikedTrackDto likedTrack)
    {
        throw new NotImplementedException();
    }

    // Self-explanatory
    public Response<bool> DeleteLikedTrack(LikedTrackDto likedTrack)
    {
        throw new NotImplementedException();
    }

    #region AccessToken Methods
    private async Task<(string?, DateTime)> RequestAccessToken()
    {
        var parameters = new Dictionary<string, string?>
        {
            { "grant_type", "client_credentials" },
            { "client_id", _configuration["SpotifyClientId"] },
            { "client_secret", _configuration["SpotifyClientSecret"] }
        };
        using var jsonContent = new FormUrlEncodedContent(parameters);

        using var response = await _httpClient.PostAsync("https://accounts.spotify.com/api/token",
            jsonContent);
        response.EnsureSuccessStatusCode();
        // process response
        var res = await response.Content.ReadFromJsonAsync<JsonElement>();
        var token = res.TryGetProperty("access_token", out var accessToken) 
            ? accessToken.GetString() : null;
        var expiry = res.TryGetProperty("expires_in", out var expires) ? expires.GetInt32() : 0;
        var expiryTime = DateTime.UtcNow.AddSeconds(expiry);
        if (token == null || expiry == 0) throw new UnauthorizedAccessException();
        
        // save in database
        var existingToken = await _context.SpotifyApi.FirstOrDefaultAsync();
        if (existingToken == null)
        {
            await _context.SpotifyApi.AddAsync(new Spotify { AccessToken = token, ExpiryDate = expiryTime });
        }
        else { existingToken.AccessToken = token; existingToken.ExpiryDate = expiryTime; }
        await _context.SaveChangesAsync();
        return (token, expiryTime);
    }

    private async Task<string?> GetAccessToken()
    {
        var existingToken = await _context.SpotifyApi.FirstOrDefaultAsync();
        if (existingToken == null)
        {
            var res = await RequestAccessToken();
            return res.Item1;
        }

        if (DateTime.UtcNow < existingToken.ExpiryDate) return existingToken.AccessToken;
        
        var resp = await RequestAccessToken();
        return resp.Item1;
    }
    

    #endregion

}