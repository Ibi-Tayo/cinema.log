using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserFilmRatingRepository
{
    Task<UserFilmRating?> CreateRating(UserFilmRating rating);
    Task<UserFilmRating?> GetRatingById(Guid id);
    Task<UserFilmRating?> GetRatingFilmUserId(Guid userId, Guid filmId);
    Task<List<UserFilmRating>> GetAllRatings(Guid userId);
    Task<UserFilmRating?> UpdateRating(UserFilmRating rating);
    Task<UserFilmRating?> DeleteRatingById(Guid id);
    Task<bool> DeleteRatingByUserFilmId(Guid userId, Guid filmId);
}