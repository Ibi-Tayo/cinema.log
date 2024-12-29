using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Services;

public class FilmService : IFilmService
{
    public async Task<FilmDto> GetFilm(Guid filmId)
    {
        throw new NotImplementedException();
    }

    public async Task<FilmDto> AddFilm(FilmDto film)
    {
        throw new NotImplementedException();
    }

    public async Task<FilmDto> UpdateFilm(FilmDto film)
    {
        throw new NotImplementedException();
    }

    public async Task<bool> DeleteFilm(Guid filmId)
    {
        throw new NotImplementedException();
    }
}