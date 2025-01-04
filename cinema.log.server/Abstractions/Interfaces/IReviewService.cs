using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IReviewService
{
    Task<Response<ReviewDto>> GetReview(Guid reviewId);
    Task<Response<ReviewDto>> AddReview(ReviewDto review);
    Task<Response<ReviewDto>> UpdateReview(ReviewDto review);
    Task<Response<ReviewDto>> DeleteReview(Guid reviewId);
    Task<Response<List<ReviewDto>>> GetReviewsByFilmId(Guid filmId);
}