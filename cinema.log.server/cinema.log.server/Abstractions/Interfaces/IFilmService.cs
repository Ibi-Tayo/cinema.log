using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IFilmService
{
    Task<Response<FilmDto>> GetFilmFromDb(Guid filmId);
    Task<Response<Guid?>> GetFilmIdUsingDetails(string title, int externalId, string? director = null, int? releaseYear = null);
    Task<Film?> AddFilmToDb(int id);
    Task<Response<FilmDto>> UpdateFilmPosterInDb(string posterUrl, Guid filmId);
    Task<Response<List<FilmSearchResultDto>>> SearchFilmFromExternal(string searchTerm);
    Task<Response<List<FilmImageDto>>> GetFilmImagesFromExternal(int externalId);
}