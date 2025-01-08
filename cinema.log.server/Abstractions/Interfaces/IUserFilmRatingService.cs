using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserFilmRatingService
{
    Task<Response<UserFilmRatingDto>> GetUserFilmRating(Guid userId, Guid filmId);
    Task<Response<UserFilmRatingDto>> AddUserFilmRating(UserFilmRatingDto filmRating);
    Task<Response<(UserFilmRatingDto?, UserFilmRatingDto?)>> FilmContest(Guid userId, Guid filmA, Guid filmB,
        Guid winnerId);
    Task<Response<List<FilmDto>>> GetFilmsForContest(Guid userId, Guid filmIdToContestAgainst);
    Task<Response<bool>> ResetAllRatings(Guid userId);
    Task<Response<bool>> DeleteRatingByUserAndFilmId(Guid userId, Guid filmId);
}