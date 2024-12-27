using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IFilmRepository
{
    Film CreateFilm(Film film);
    Film GetFilmById(Guid id);
    Film UpdateFilm(Film film);
    Film DeleteFilmById(Guid id);
}