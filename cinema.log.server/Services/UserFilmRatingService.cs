using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Services;

public class UserFilmRatingService : IUserFilmRatingService
{
    public async Task<UserFilmRatingDto> GetUserFilmRating(Guid userId, Guid filmId)
    {
        throw new NotImplementedException();
    }

    public async Task<UserFilmRatingDto> AddUserFilmRating(UserFilmRatingDto filmRating)
    {
        throw new NotImplementedException();
    }

    public async Task<UserFilmRatingDto> UpdateUserFilmRating(UserFilmRatingDto filmRating)
    {
        throw new NotImplementedException();
    }

    public async Task<bool> DeleteUserFilmRating(Guid userId, Guid filmId)
    {
        throw new NotImplementedException();
    }
}