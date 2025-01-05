using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IFilmRepository
{
    Task<Film?> CreateFilm(Film film);
    Task<Film?> GetFilmById(Guid id);
    Task<List<Film>> GetFilmsByIds(List<Guid> ids);
    Task<Film?> UpdateFilm(Film film);
    Task<Film?> DeleteFilmById(Guid id);
    Task<Guid?> GetFilmId(string title, string? director, int? releaseYear);
}