using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserFilmRatingService
{
    Task<UserFilmRatingDto> GetUserFilmRating(Guid userId, Guid filmId);
    Task<UserFilmRatingDto> AddUserFilmRating(UserFilmRatingDto filmRating);
    Task<UserFilmRatingDto> UpdateUserFilmRating(UserFilmRatingDto filmRating);
    Task<bool> DeleteUserFilmRating(Guid userId, Guid filmId);
}