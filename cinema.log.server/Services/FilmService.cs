using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

namespace cinema.log.server.Services;

public class FilmService : IFilmService
{
    private readonly IConfiguration _config;
    private static HttpClient _client = new()
    {
        BaseAddress = new Uri("https://api.themoviedb.org/3/")
    };

    public FilmService(IConfiguration config)
    {
        _config = config;
        // Ensure all outgoing requests have this header
        _client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
    }
    
    public async Task<Response<FilmDto>> GetFilmFromDb(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // this will have to get called server side - probably in the review service when a user leaves a review
    public async Task<bool> AddFilmToDb(int externalId)
    {
        // call GetFilmDetailsByExternalId() 
        // use the film we get back to add to db using repository
        throw new NotImplementedException();
    }

    public async Task<Response<FilmDto>> UpdateFilmInDb(FilmDto film)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<FilmDto>> DeleteFilmInDb(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // To remember:
    // 'External' basically means im sending an api call to tmdb
    // External id is the movie_id param that is an integer, obviously not to be confused with my own GUID
    
    public async Task<Response<List<FilmSearchResultDto>>> SearchFilmFromExternal(string searchTerm)
    {
        var films = new List<FilmSearchResultDto>();
        var key = _config["TmdbApiKey"];
        var reqUrl = $"search/movie?query={searchTerm}&include_adult=false&language=en-US&page=1&api_key={key}";
        
        using (var response = await _client.GetAsync(reqUrl))
        {
            response.EnsureSuccessStatusCode();
            var body = await response.Content.ReadAsStringAsync();
            var json = JsonSerializer.Deserialize<JsonElement>(body);
            json.TryGetProperty("results", out var results);
            
            foreach (var result in results.EnumerateArray())
            {
                result.TryGetProperty("id", out var id);
                // Not adding any films without a title
                if (!result.TryGetProperty("original_title", out var title)) continue;
                result.TryGetProperty("overview", out var description);
                result.TryGetProperty("release_date", out var releaseDate);
                result.TryGetProperty("poster_path", out var poster);
             
                films.Add(new FilmSearchResultDto()
                {
                    ExternalId = id.GetInt32(),
                    Title = title.GetString() ?? throw new NullReferenceException(),
                    Description = description.GetString(),
                    ReleaseYear = releaseDate.TryGetDateTime(out var date) ? date.Year : null,
                    PosterUrl = poster.GetString(),
                });
            }
        }
        return Response<List<FilmSearchResultDto>>.BuildResponse(200, "Success", films);
    }

    public async Task<Response<List<FilmImageDto>>> GetFilmImagesFromExternal(int externalId)
    {
        // call the tmdb api and hit images endpoint, get all the backdrops and all the posters
        // map these to the dtos and make a list
        
        throw new NotImplementedException();
    }
    
    // flow: when user has found a film, they should have the film details, when they leave a review, instead of sending all of this back
    // they just need to ping the server with the tmdb movie_id. this method will call tmdb and then map to our entity
    // then the entity will end up getting put in the db - see 'AddFilmToDb' above
    private async Task<Film?> GetFilmDetailsByExternalId(int externalId)
    {
        // this method will get called in AddFilmToDb
        
        var key = _config["TmdbApiKey"];
        var reqUrl = $"movie/{externalId}?api_key={key}";
        using var response = await _client.GetAsync(reqUrl);
        response.EnsureSuccessStatusCode();
        
        var body = await response.Content.ReadAsStringAsync();
        var json = JsonSerializer.Deserialize<JsonElement>(body);
            
        if (!json.TryGetProperty("original_title", out var title)) return null;
        json.TryGetProperty("overview", out var description);
        json.TryGetProperty("release_date", out var releaseDate);
        json.TryGetProperty("poster_path", out var poster);
        json.TryGetProperty("genres", out var genres);

        var genreString = new StringBuilder("");
        foreach (var genre in genres.EnumerateArray())
        {
            genre.TryGetProperty("name", out var name);
            genreString.Append(name.GetString());
            genreString.Append(',');
        }
                
        return new Film()
        {
            FilmId = Guid.NewGuid(),
            Title = title.GetString() ?? throw new NullReferenceException(),
            Description = description.GetString(),
            Genre = genreString.Length == 0 ? null : genreString.ToString().TrimEnd(','),
            Director = await GetFilmDirector(externalId),
            ReleaseYear = releaseDate.TryGetDateTime(out var date) ? date.Year : null,
            PosterUrl = poster.GetString(),
        };
    }

    private async Task<string?> GetFilmDirector(int externalId)
    {
        var key = _config["TmdbApiKey"];
        var reqUrl = $"movie/{externalId}/credits?api_key={key}";
        using var response = await _client.GetAsync(reqUrl);
        response.EnsureSuccessStatusCode();
        
        var body = await response.Content.ReadAsStringAsync();
        var json = JsonSerializer.Deserialize<JsonElement>(body);
        if (!json.TryGetProperty("crew", out var crew)) return null;
        foreach (var crewMember in crew.EnumerateArray())
        {
            if (crewMember.TryGetProperty("job", out var job) && 
                string.Equals(job.GetString(), "Director", StringComparison.OrdinalIgnoreCase))
            {
                return crewMember.TryGetProperty("name", out var name) ? name.GetString() : null;
            }
        }
        return null;
    }
}