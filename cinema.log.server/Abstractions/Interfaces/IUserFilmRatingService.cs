using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserFilmRatingService
{
    Task<Response<UserFilmRatingDto>> GetUserFilmRating(Guid userId, Guid filmId);
    Task<Response<UserFilmRatingDto>> AddUserFilmRating(UserFilmRatingDto filmRating);
    Task<Response<UserFilmRatingDto>> UpdateUserFilmRating(UserFilmRatingDto filmRating);
    Task<Response<UserFilmRatingDto>> DeleteUserFilmRating(Guid userId, Guid filmId);
}