using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Repositories;

public class ReviewRepository(CinemaLogContext context, ILogger<ReviewRepository> logger) : IReviewRepository
{
    private CinemaLogContext _context = context;
    private ILogger<ReviewRepository> _logger = logger;
    
    public async Task<Review?> CreateReview(Review review)
    {
        try
        {
            await _context.Reviews.AddAsync(review);
            await _context.SaveChangesAsync();
            return review;
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
    }

    public async Task<Review?> GetReviewById(Guid id)
    {
        return await _context.Reviews.FindAsync(id);
    }

    public async Task<Review?> UpdateReview(Review review)
    {
        try
        {
            _context.Reviews.Update(review);
            await _context.SaveChangesAsync();
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }

        return review;
    }

    public async Task<bool> DeleteReviewById(Guid id)
    {
        var foundReview = await _context.Reviews.FindAsync(id);
        if (foundReview == null) return false;
        _context.Reviews.Remove(foundReview);
        await _context.SaveChangesAsync();
        return true;
    }

    public async Task<Review?> GetReviewByUserAndFilm(Guid userId, Guid filmId)
    {
        var foundReview = await _context.Reviews.FirstOrDefaultAsync(review => review.UserId == userId && 
                                                                               review.FilmId == filmId);
        return foundReview;
    }

    public async Task<List<Review>> GetReviewsByFilmId(Guid filmId)
    {
        return await _context.Reviews.Where(review => review.FilmId == filmId).ToListAsync();
    }
    
    public async Task<List<Review>> GetReviewsByUserId(Guid userId)
    {
        return await _context.Reviews.Where(review => review.UserId == userId).ToListAsync();
    }
}