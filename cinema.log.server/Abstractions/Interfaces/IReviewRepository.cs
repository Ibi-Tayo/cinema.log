using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IReviewRepository
{
    Review CreateReview(Review review);
    Review GetReviewById(Guid id);
    Review UpdateReview(Review review);
    Review DeleteReviewById(Guid id);
}