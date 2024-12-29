using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Services;

public class ReviewService : IReviewService
{
    public async Task<ReviewDto> GetReview(Guid reviewId)
    {
        throw new NotImplementedException();
    }

    public async Task<ReviewDto> AddReview(ReviewDto review)
    {
        throw new NotImplementedException();
    }

    public async Task<ReviewDto> UpdateReview(ReviewDto review)
    {
        throw new NotImplementedException();
    }

    public async Task<bool> DeleteReview(Guid reviewId)
    {
        throw new NotImplementedException();
    }
}