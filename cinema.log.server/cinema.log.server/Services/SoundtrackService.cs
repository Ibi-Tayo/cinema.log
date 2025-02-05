using System.Net.Http.Headers;
using System.Text.Json;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.ApiResponse;
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
    public SoundtrackService(CinemaLogContext context, HttpClient httpClient, IConfiguration configuration)
    {
        _context = context;
        _httpClient = httpClient;
        _configuration = configuration;
    }

    public async Task<Response<FilmSoundtrackDto>> GetSoundtrackByFilmId(Guid filmId, Guid userId)
    {
        // Check the cache first
        var cachedSoundtrack = await _context.CachedSoundtracks
            .FirstOrDefaultAsync(cs => cs.FilmId == filmId);
        if (cachedSoundtrack != null)
        {
            var ufr = await GetOrDefaultUserFilmSoundtrackRating(filmId, userId);
            var fsDto = MapCachedToDto(cachedSoundtrack, ufr);
            return Response<FilmSoundtrackDto>.BuildResponse(200, "Success", fsDto);
        }

        // Use id to get film name
        var filmTitle = _context.Films.SingleOrDefault(film => film.FilmId == filmId)?.Title;
        if (string.IsNullOrEmpty(filmTitle))
            return Response<FilmSoundtrackDto>
                .BuildResponse(404, "Film not found", null);

        // Use film name in search api call to spotify - returns array of results
        var albums = await SearchSpotifyAlbums(filmTitle);

        // Find a match in array where the title includes 1. Film name, 2. "Soundtrack" or "Score" etc
        var soundtrack = albums?.FirstOrDefault(album =>
                             album.Name.Contains(filmTitle, StringComparison.OrdinalIgnoreCase) &&
                             (album.Name.Contains("Score") ||
                              album.Name.Contains("Soundtrack") ||
                              album.Name.Contains("Motion Picture")))
                         // try searching just film title
                         ?? albums?.FirstOrDefault(a => a.Name.Contains(filmTitle, StringComparison.OrdinalIgnoreCase));
        if (soundtrack == null)
            return Response<FilmSoundtrackDto>
                .BuildResponse(404, "Soundtrack not found in spotify search", null);

        // Get Spotify Album id and use in api call to get full album details with all tracks 
        var album = await GetSpotifyAlbum(soundtrack.Id);
        if (album == null)
            return Response<FilmSoundtrackDto>
                .BuildResponse(404, "Album Id not found in spotify", null);
        
        var userFilmRating = await GetOrDefaultUserFilmSoundtrackRating(filmId, userId);

        var dto = await CreateFilmSoundtrackDto(filmId, userFilmRating, filmTitle, album);
        return Response<FilmSoundtrackDto>.BuildResponse(200, "Success", dto);
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

        var soundtrackRating =
            await _context.UserFilmSoundtrackRatings.FindAsync(likedTrack.UserFilmSoundtrackRatingId);
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

        var soundtrackRating =
            await _context.UserFilmSoundtrackRatings.FindAsync(likedTrack.UserFilmSoundtrackRatingId);
        if (soundtrackRating == null)
            return Response<bool>.BuildResponse(404, "Soundtrack rating not found", false);

        _context.LikedTracks.Remove(Mapper<LikedTrackDto, LikedTrack>.Map(likedTrack));
        await _context.SaveChangesAsync();
        return Response<bool>.BuildResponse(200, "Success", true);
    }

    public async Task<Response<bool>> UpdateSoundtrackRating(Guid filmId, Guid userId, int rating)
    {
        var userSoundtrackRating = await GetOrDefaultUserFilmSoundtrackRating(filmId, userId);
        userSoundtrackRating.Rating = rating;
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
            ? accessToken.GetString()
            : null;
        var expiry = res.TryGetProperty("expires_in", out var expires) ? expires.GetInt32() : 0;
        var expiryTime = DateTime.UtcNow.AddSeconds(expiry);
        if (token == null || expiry == 0) throw new UnauthorizedAccessException();

        // save in database
        var existingToken = await _context.SpotifyApi.FirstOrDefaultAsync();
        if (existingToken == null)
        {
            await _context.SpotifyApi.AddAsync(new Spotify { AccessToken = token, ExpiryDate = expiryTime });
        }
        else
        {
            existingToken.AccessToken = token;
            existingToken.ExpiryDate = expiryTime;
        }

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

    #region Helpers

    private async Task<List<Album>?> SearchSpotifyAlbums(string filmTitle)
    {
        var token = await GetAccessToken();
        using var requestMessage =
            new HttpRequestMessage(HttpMethod.Get,
                $"https://api.spotify.com/v1/search?q={filmTitle}&type=album&limit=3");
        requestMessage.Headers.Authorization =
            new AuthenticationHeaderValue("Bearer", token);

        var responseMessage = await _httpClient.SendAsync(requestMessage);
        var responseString = await responseMessage.Content.ReadAsStringAsync();
        var responseObj = JsonSerializer.Deserialize<SpotifyAlbumSearchResponse>(responseString);
        return responseObj?.Albums.Items;
    }

    private async Task<Album?> GetSpotifyAlbum(string albumId)
    {
        var token = await GetAccessToken();
        using var requestMessage =
            new HttpRequestMessage(HttpMethod.Get, $"https://api.spotify.com/v1/albums/{albumId}");
        requestMessage.Headers.Authorization =
            new AuthenticationHeaderValue("Bearer", token);

        var responseMessage = await _httpClient.SendAsync(requestMessage);
        var responseString = await responseMessage.Content.ReadAsStringAsync();
        var responseObj = JsonSerializer.Deserialize<Album>(responseString);
        return responseObj;
    }

    private async Task<UserFilmSoundtrackRating> GetOrDefaultUserFilmSoundtrackRating(Guid filmId, Guid userId)
    {
        var userFilmRating = _context.UserFilmSoundtrackRatings
            .SingleOrDefault(filmSoundtrackRating => filmSoundtrackRating.FilmId == filmId &&
                                                     filmSoundtrackRating.UserId == userId);

        if (userFilmRating != null) return userFilmRating;
        // Create a new entry in the table without a rating
        userFilmRating = new UserFilmSoundtrackRating
        {
            UserFilmSoundtrackRatingId = Guid.NewGuid(),
            FilmId = filmId,
            UserId = userId,
        };
        await _context.UserFilmSoundtrackRatings.AddAsync(userFilmRating);
        await _context.SaveChangesAsync();

        return userFilmRating;
    }

    private async Task<FilmSoundtrackDto> CreateFilmSoundtrackDto(
        Guid filmId, UserFilmSoundtrackRating userFilmRating,
        string filmTitle, Album album)
    {
        // Add to cache
        var cachedSoundtrack = new CachedSoundtrack()
        {
            Id = Guid.NewGuid(), FilmId = filmId, FilmTitle = filmTitle, SpotifyAlbumId = album.Id,
            SoundtrackName = album.Name, Artists = string.Join(", ", album.Artists.Select(artist => artist.Name)),
            AlbumArtUrl = album.Images[0].Url, LastUpdated = DateTime.Now, 
            TracksJson = JsonSerializer.Serialize(album.Tracks)
        };
        await _context.CachedSoundtracks.AddAsync(cachedSoundtrack);
        await _context.SaveChangesAsync();

        return new FilmSoundtrackDto
        {
            UserFilmSoundtrackRatingId = userFilmRating.UserFilmSoundtrackRatingId, FilmId = filmId,
            FilmName = filmTitle, Artists = string.Join(", ", album.Artists.Select(artist => artist.Name)),
            SoundtrackName = album.Name, AlbumArtUrl = album.Images[0].Url, 
            Rating = userFilmRating.Rating,
            Tracks = album.Tracks.Items.Select(track => new TrackDto
            {
                ArtistName = string.Join(", ", track.Artists.Select(artist => artist.Name)),
                TrackTitle = track.Name, TrackUrl = track.ExternalUrls.Spotify
            }).ToList(),
            LikedTracks = _context.LikedTracks
                .Where(lt => lt.UserFilmSoundtrackRatingId == userFilmRating.UserFilmSoundtrackRatingId)
                .Select(likedTrack => new LikedTrackDto()
                {
                    Id = likedTrack.Id, UserId = likedTrack.UserId, TrackTitle = likedTrack.TrackTitle,
                    UserFilmSoundtrackRatingId = likedTrack.UserFilmSoundtrackRatingId
                }).ToList()
        };
    }

    private FilmSoundtrackDto MapCachedToDto(
        CachedSoundtrack cachedSoundtrack,
        UserFilmSoundtrackRating userFilmRating)
    {
        // Deserialize the tracks from JSON
        var tracks = JsonSerializer.Deserialize<List<Track>>(cachedSoundtrack.TracksJson);

        return new FilmSoundtrackDto
        {
            UserFilmSoundtrackRatingId = userFilmRating.UserFilmSoundtrackRatingId, FilmId = cachedSoundtrack.FilmId,
            FilmName = cachedSoundtrack.FilmTitle, Artists = cachedSoundtrack.Artists,
            SoundtrackName = cachedSoundtrack.SoundtrackName, AlbumArtUrl = cachedSoundtrack.AlbumArtUrl,
            Rating = userFilmRating.Rating,
            Tracks = tracks?.Select(track => new TrackDto
            {
                // Extract artist names from deserialized Track objects
                ArtistName = string.Join(", ", track.Artists.Select(a => a.Name)), TrackTitle = track.Name,
                TrackUrl = track.ExternalUrls.Spotify
            }).ToList() ?? [], // Handle null tracks
            LikedTracks = _context.LikedTracks
                .Where(lt => lt.UserFilmSoundtrackRatingId == userFilmRating.UserFilmSoundtrackRatingId)
                .Select(lt => new LikedTrackDto
                {
                    Id = lt.Id, UserId = lt.UserId, TrackTitle = lt.TrackTitle, 
                    UserFilmSoundtrackRatingId = lt.UserFilmSoundtrackRatingId
                }).ToList()
        };
    }

    #endregion
}