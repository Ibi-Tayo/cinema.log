using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Repositories;

public class ReviewRepository(CinemaLogContext context) : IReviewRepository
{
    private CinemaLogContext _context = context;
    
    public async Task<Review> CreateReview(Review review)
    {
        throw new NotImplementedException();
    }

    public async Task<Review> GetReviewById(Guid id)
    {
        throw new NotImplementedException();
    }

    public async Task<Review> UpdateReview(Review review)
    {
        throw new NotImplementedException();
    }

    public async Task<Review> DeleteReviewById(Guid id)
    {
        throw new NotImplementedException();
    }
}