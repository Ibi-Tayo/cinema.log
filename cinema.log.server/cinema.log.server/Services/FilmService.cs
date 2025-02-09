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
    private readonly IFilmRepository _filmRepository;

    private static HttpClient _client = new()
    {
        BaseAddress = new Uri("https://api.themoviedb.org/3/")
    };

    public FilmService(IConfiguration config, IFilmRepository filmRepository)
    {
        _config = config;
        _filmRepository = filmRepository;
        // Ensure all outgoing requests have this header
        _client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
    }
    
    // Important method to ensure filmId is always present for every film that is selected, reviewed etc.
    public async Task<Response<Guid?>> GetFilmIdUsingDetails(string title, int externalId, string? director, int? releaseYear)
    {
        // Check if film exists in the film table
        var id = await _filmRepository.GetFilmId(title, director, releaseYear);
        if (id != null) return Response<Guid?>.BuildResponse(200, "Success, film exists", id);

        // If not, then add to the table using externalId to search for film details
        var film = await AddFilmToDb(externalId);
        if (film != null)
        {
            return Response<Guid?>.BuildResponse(200, "Success, new film added", film.FilmId);
        }
        return Response<Guid?>.BuildResponse(500, "Internal Server Error", null);
    }

    public async Task<Response<FilmDto>> GetFilmFromDb(Guid filmId)
    {
        var film = await _filmRepository.GetFilmById(filmId);
        if (film is null)
        {
            return Response<FilmDto>.BuildResponse(404, "Film not found", null);
        }

        var responseFilm = Mapper<Film, FilmDto>.Map(film);
        return Response<FilmDto>.BuildResponse(200, "Success", responseFilm);
    }

    public async Task<Film?> AddFilmToDb(int externalId)
    {
        var film = await GetFilmDetailsByExternalId(externalId);
        if (film is null) return null;
        var addedFilm = await _filmRepository.CreateFilm(film);
        return addedFilm;
    }

    public async Task<Response<FilmDto>> UpdateFilmPosterInDb(string posterUrl, Guid filmId)
    {
        var film = await _filmRepository.GetFilmById(filmId);
        if (film == null) return Response<FilmDto>.BuildResponse(404, "Film not found", null);
        film.PosterUrl = posterUrl;
        var updatedFilmDto = Mapper<Film, FilmDto>.Map(film);
        return Response<FilmDto>.BuildResponse(200, "Success", updatedFilmDto);
    }

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
        var key = _config["TmdbApiKey"];
        var reqUrl = $"movie/{externalId}/images?api_key={key}";

        using var response = await _client.GetAsync(reqUrl);
        response.EnsureSuccessStatusCode();
    
        var json = JsonSerializer.Deserialize<JsonElement>(await response.Content.ReadAsStringAsync());
        var images = new List<FilmImageDto>();

        if (json.TryGetProperty("backdrops", out var backdrops))
            images.AddRange(ProcessImageArray(backdrops, FilmImageType.Backdrop));
    
        if (json.TryGetProperty("posters", out var posters))
            images.AddRange(ProcessImageArray(posters, FilmImageType.Poster));
    
        if (json.TryGetProperty("logos", out var logos))
            images.AddRange(ProcessImageArray(logos, FilmImageType.Logo));

        return Response<List<FilmImageDto>>.BuildResponse(200, "Success", images);
    }
    
    private List<FilmImageDto> ProcessImageArray(JsonElement imagesArray, FilmImageType imageType)
    {
        return imagesArray.EnumerateArray()
            .Where(img => img.TryGetProperty("file_path", out _))
            .Select(img => new FilmImageDto
            {
                ImageType = imageType,
                AspectRatio = img.GetProperty("aspect_ratio").GetSingle(),
                Height = img.GetProperty("height").GetInt32(),
                Width = img.GetProperty("width").GetInt32(),
                Url = img.GetProperty("file_path").GetString()
            })
            .ToList();
    }
    
    private string GetApiKey()
    {
        // Try environment variable first (for CI/CD)
        var key = Environment.GetEnvironmentVariable("TmdbApiKey");
        Console.WriteLine($"Environment Variable Path: Key exists = {!string.IsNullOrEmpty(key)}");
    
        // Fall back to user secrets/configuration
        if (string.IsNullOrEmpty(key))
        {
            key = _config["TmdbApiKey"];
            Console.WriteLine($"Configuration Path: Key exists = {!string.IsNullOrEmpty(key)}");
        }

        if (string.IsNullOrEmpty(key))
        {
            throw new InvalidOperationException(
                "TMDB API key not found. Ensure it's set in user secrets for local development " +
                "or as an environment variable 'TmdbApiKey' for CI/CD.");
        }

        return key;
    }
    
    private async Task<Film?> GetFilmDetailsByExternalId(int externalId)
    {
        var key = GetApiKey();
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