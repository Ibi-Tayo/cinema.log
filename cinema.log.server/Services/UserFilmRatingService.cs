using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class UserFilmRatingService : IUserFilmRatingService
{
    // remember to inject calculation service
    public async Task<Response<UserFilmRatingDto>> GetUserFilmRating(Guid userId, Guid filmId)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<UserFilmRatingDto>> AddUserFilmRating(UserFilmRatingDto filmRating)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<UserFilmRatingDto>> UpdateUserFilmRating(UserFilmRatingDto filmRating)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<UserFilmRatingDto>> DeleteUserFilmRating(Guid userId, Guid filmId)
    {
        throw new NotImplementedException();
    }
}