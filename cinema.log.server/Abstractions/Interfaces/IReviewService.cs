using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IReviewService
{
    Task<Response<ReviewDto>> GetReview(Guid reviewId);
    Task<Response<List<ReviewDto>>> GetReviewsByUserId(Guid userId);
    Task<Response<List<ReviewDto>>> GetReviewsByFilmId(Guid filmId);
    Task<Response<ReviewDto>> GetReviewByUserAndFilm(Guid userId, Guid filmId);
    Task<Response<ReviewDto>> AddReview(ReviewDto review);
    Task<Response<ReviewDto>> UpdateReview(ReviewDto review);
    Task<Response<bool>> DeleteReview(Guid reviewId);


    
}