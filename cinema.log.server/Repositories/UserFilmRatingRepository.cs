using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Repositories;

public class UserFilmRatingRepository(CinemaLogContext context, ILogger<UserFilmRatingRepository> logger) : IUserFilmRatingRepository
{
    private CinemaLogContext _context = context;
    private ILogger<UserFilmRatingRepository> _logger = logger;
    
    public async Task<UserFilmRating?> CreateRating(UserFilmRating rating)
    {
        try
        {
            await _context.UserFilmRatings.AddAsync(rating);
            await _context.SaveChangesAsync();
            return rating;
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
    }

    public async Task<UserFilmRating?> GetRatingById(Guid id)
    {
        return await _context.UserFilmRatings.FindAsync(id);
    }

    public async Task<List<UserFilmRating>> GetAllRatings(Guid userId)
    {
        return await _context.UserFilmRatings.Where(x => x.UserId == userId).ToListAsync();
    }

    public async Task<List<Guid>> GetAllFilmIds(Guid userId)
    {
        return await _context.UserFilmRatings
            .Where(ufr => ufr.UserId == userId)
            .Select(rating => rating.FilmId).ToListAsync();
    }

    public async Task<List<Film>> GetAllFilmsRatedByUserId(Guid userId)
    {
        return await _context.UserFilmRatings
            .Where(ufr => ufr.UserId == userId)
            .Join(
                _context.Films,
                ufr => ufr.FilmId,
                film => film.FilmId,
                (ufr, film) => film
            )
            .ToListAsync();
    }

    public async Task<UserFilmRating?> GetRatingFilmUserId(Guid userId, Guid filmId)
    {
        return await _context.UserFilmRatings.FirstOrDefaultAsync(ufr => ufr.UserId == userId &&
                                                                         ufr.FilmId == filmId);
    }

    public async Task<UserFilmRating?> UpdateRating(UserFilmRating rating)
    {
        try
        {
            _context.UserFilmRatings.Update(rating);
            await _context.SaveChangesAsync();
            return rating;
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
    }

    public async Task<UserFilmRating?> DeleteRatingById(Guid id)
    {
        var foundRating = await _context.UserFilmRatings.FindAsync(id);
        if (foundRating == null) return null;
        _context.UserFilmRatings.Remove(foundRating);
        await _context.SaveChangesAsync();
        return foundRating;
    }

    public async Task<bool> DeleteRatingByUserFilmId(Guid userId, Guid filmId)
    {
        var rating = await _context.UserFilmRatings
            .FirstOrDefaultAsync(ufr => ufr.UserId == userId && ufr.FilmId == filmId);
        if (rating == null) return false;
        _context.UserFilmRatings.Remove(rating);
        await _context.SaveChangesAsync();
        return true;
    }
    
    public async Task<List<UserFilmRating>> GetAllUserFilmRatingsRanked(Guid userId)
    {
        var ratings = await _context.UserFilmRatings
            .Where(r => r.UserId == userId)
            .OrderByDescending(r => r.EloRating)
            .ToListAsync();

        return ratings;
    }
}