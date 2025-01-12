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

    public async Task<Response<FilmSoundtrackDto>> GetSoundtrackByFilmId(Guid filmId)
    {
        // for testing purposes
        var dto = new FilmSoundtrackDto() { SoundtrackName = await GetAccessToken() };
        return Response<FilmSoundtrackDto>.BuildResponse(200, "Success", dto);
        
        // Use Id to get film name
        // Use film name in search api call to spotify - returns array of results
        // Find a match in array where the title includes 1. Film name, 2. "Soundtrack" or "Score"
        // Get Spotify Album Id and use in api call to get album
        // Use data to populate film soundtrack dto and return
    }

    // use user id to get all liked tracks from all films
    public async Task<Response<List<LikedTrackDto>>> GetLikedTracksByUserId(Guid userId)
    {
        var res = await _context.LikedTracks.Where(lt => lt.UserId == userId).ToListAsync();
        var dto = res.Select(Mapper<LikedTrack, LikedTrackDto>.Map).ToList();
        return Response<List<LikedTrackDto>>.BuildResponse(200, "Success", dto);
    }
    
    // When looking at film rating, get all liked songs in specific film
    public async Task<Response<List<LikedTrackDto>>> GetLikedTracksFromFilmRatingId(Guid filmRatingId)
    {
        var res = await _context.LikedTracks
            .Where(lt => lt.UserFilmSoundtrackRatingId == filmRatingId).ToListAsync();
        var dto = res.Select(Mapper<LikedTrack, LikedTrackDto>.Map).ToList();
        return Response<List<LikedTrackDto>>.BuildResponse(200, "Success", dto);
    }

    // use dto to add new liked track to liked track table
    // (user needs to make sure they have the film soundtrack rating id)
    // (they'd get that by calling GetSoundtrackByFilmId)
    public async Task<Response<LikedTrackDto?>> SetLikedTrack(LikedTrackDto likedTrack)
    {
        var user = await _context.Users.FindAsync(likedTrack.UserId);
        if (user == null) 
            return Response<LikedTrackDto?>.BuildResponse(404, "User not found", null);
        
        var soundtrackRating = await _context.UserFilmSoundtrackRatings.FindAsync(likedTrack.UserFilmSoundtrackRatingId);
        if (soundtrackRating == null) 
            return Response<LikedTrackDto?>.BuildResponse(404, 
                "Soundtrack rating not found", null);
       
        likedTrack.Id = Guid.NewGuid(); // set liked track id here

        var entity = _context.LikedTracks.Add(Mapper<LikedTrackDto, LikedTrack>.Map(likedTrack)).Entity;
        await _context.SaveChangesAsync();
        return Response<LikedTrackDto?>.BuildResponse(200, "Success", 
            Mapper<LikedTrack, LikedTrackDto>.Map(entity));
    }

    public async Task<Response<bool>> DeleteLikedTrack(LikedTrackDto likedTrack)
    {
        var track = await _context.LikedTracks.FindAsync(likedTrack.Id);
        if (track == null)
            return Response<bool>.BuildResponse(404, "Track not found", false);
        
        var user = await _context.Users.FindAsync(likedTrack.UserId);
        if (user == null)
            return Response<bool>.BuildResponse(404, "User not found", false);
        
        var soundtrackRating = await _context.UserFilmSoundtrackRatings.FindAsync(likedTrack.UserFilmSoundtrackRatingId);
        if (soundtrackRating == null)
            return Response<bool>.BuildResponse(404, "Soundtrack rating not found", false);
        
        _context.LikedTracks.Remove(Mapper<LikedTrackDto, LikedTrack>.Map(likedTrack));
        await _context.SaveChangesAsync();
        return Response<bool>.BuildResponse(200, "Success", true);
        
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
        // expired token so request new one
        var resp = await RequestAccessToken();
        return resp.Item1;
    }
    #endregion

}