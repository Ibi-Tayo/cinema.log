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

    public async Task<UserFilmRating?> GetRatingFilmUserId(Guid userId, Guid filmId)
    {
        return await _context.UserFilmRatings.FirstOrDefaultAsync(ufr => ufr.User.UserId == userId &&
                                                                         ufr.Film.FilmId == filmId);
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

    public async Task<float?> GetUserFilmRating(Guid userId, Guid filmId)
    {
        var rating = await _context.UserFilmRatings
            .FirstOrDefaultAsync(r =>
                r.Film.FilmId == filmId &&
                r.User.UserId == userId);
           
        return rating?.EloRating;
    }

    public async Task<List<UserFilmRating>> GetAllUserFilmRatingsRanked(Guid userId)
    {
        var ratings = await _context.UserFilmRatings
            .Where(r => r.User.UserId == userId)
            .OrderByDescending(r => r.EloRating)
            .ToListAsync();

        return ratings;
    }
}