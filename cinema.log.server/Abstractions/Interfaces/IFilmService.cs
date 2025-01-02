using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IFilmService
{
    Task<Response<FilmDto>> GetFilmFromDb(Guid filmId);
    Task<bool> AddFilmToDb(int id);
    Task<Response<FilmDto>> UpdateFilmInDb(FilmDto film);
    Task<Response<FilmDto>>  DeleteFilmInDb(Guid filmId);
    Task<Response<List<FilmDto>>> SearchFilmFromExternal(string searchTerm);
    Task<Response<List<FilmImageDto>>> GetFilmImagesFromExternal(int externalId);
}