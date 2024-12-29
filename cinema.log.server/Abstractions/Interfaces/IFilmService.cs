using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IFilmService
{
    Task<FilmDto> GetFilm(Guid filmId);
    Task<FilmDto> AddFilm(FilmDto film);
    Task<FilmDto> UpdateFilm(FilmDto film);
    Task<bool> DeleteFilm(Guid filmId);
}