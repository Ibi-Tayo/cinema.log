using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IReviewService
{
    Task<ReviewDto> GetReview(Guid reviewId);
    Task<ReviewDto> AddReview(ReviewDto review);
    Task<ReviewDto> UpdateReview(ReviewDto review);
    Task<bool> DeleteReview(Guid reviewId);
}