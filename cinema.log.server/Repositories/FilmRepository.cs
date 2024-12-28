using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Repositories;

public class FilmRepository(CinemaLogContext context) : IFilmRepository
{
    private CinemaLogContext _context = context;

    public async Task<Film> CreateFilm(Film film)
    {
        throw new NotImplementedException();
    }

    public async Task<Film> GetFilmById(Guid id)
    {
        throw new NotImplementedException();
    }

    public async Task<Film> UpdateFilm(Film film)
    {
        throw new NotImplementedException();
    }

    public async Task<Film> DeleteFilmById(Guid id)
    {
        throw new NotImplementedException();
    }
}