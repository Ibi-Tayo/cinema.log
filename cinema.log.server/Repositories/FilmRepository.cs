using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Repositories;

public class FilmRepository: IFilmRepository
{
    CinemaLogContext _context;
    ILogger<FilmRepository> _logger;

    public FilmRepository(CinemaLogContext context, ILogger<FilmRepository> logger)
    {
        _context = context;
        _logger = logger;
    }
    public async Task<Film?> CreateFilm(Film film)
    {
        try
        {
            await _context.Films.AddAsync(film);
            await _context.SaveChangesAsync();
            return film;
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
    }

    public async Task<Film?> GetFilmById(Guid id)
    {
        return await _context.Films.FindAsync(id);
    }

    public async Task<List<Film>> GetFilmByTitle(string title)
    {
        return await _context.Films
            .Where(film => film.Title.ToLower() == title.ToLower()).ToListAsync();
    }

    public async Task<Film?> UpdateFilm(Film film)
    {
        try
        {
            _context.Films.Update(film);
            await _context.SaveChangesAsync();
            return film;
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
    }

    public async Task<Film?> DeleteFilmById(Guid id)
    {
        var foundFilm = await _context.Films.FindAsync(id);
        if (foundFilm == null) return null;
        _context.Films.Remove(foundFilm);
        await _context.SaveChangesAsync();
        return foundFilm;
    }
}