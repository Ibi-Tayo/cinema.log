using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IReviewRepository
{
    Task<Review?> CreateReview(Review review);
    Task<Review?> GetReviewById(Guid id);
    Task<Review?> UpdateReview(Review review);
    Task<bool> DeleteReviewById(Guid id);
    Task<Review?> GetReviewByUserAndFilm(Guid userId, Guid filmId);
    Task<List<Review>> GetReviewsByFilmId(Guid filmId);
    Task<List<Review>> GetReviewsByUserId(Guid userId);
}