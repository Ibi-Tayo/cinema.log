using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserFilmRatingRepository
{
    Task<UserFilmRating?> CreateRating(UserFilmRating rating);
    Task<UserFilmRating?> GetRatingById(Guid id);
    Task<UserFilmRating?> UpdateRating(UserFilmRating rating);
    Task<UserFilmRating?> DeleteRatingById(Guid id);
}