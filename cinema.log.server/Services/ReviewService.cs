using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class ReviewService : IReviewService
{
    public async Task<Response<ReviewDto>> GetReview(Guid reviewId)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<ReviewDto>> AddReview(ReviewDto review)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<ReviewDto>> UpdateReview(ReviewDto review)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<ReviewDto>> DeleteReview(Guid reviewId)
    {
        throw new NotImplementedException();
    }
}