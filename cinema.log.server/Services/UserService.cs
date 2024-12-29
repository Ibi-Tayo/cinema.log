using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Services;

public class UserService : IUserService
{
    public async Task<UserDto> GetUser(Guid userId)
    {
        throw new NotImplementedException();
    }

    public async Task<UserDto> AddUser(UserDto user)
    {
        throw new NotImplementedException();
    }

    public async Task<UserDto> UpdateUser(UserDto user)
    {
        throw new NotImplementedException();
    }

    public async Task<bool> DeleteUser(Guid userId)
    {
        throw new NotImplementedException();
    }

    public async Task<List<ReviewDto>> GetUserReviews(Guid userId)
    {
        throw new NotImplementedException();
    }

    public async Task<List<FilmDto>> GetFilmsReviewedByUser(Guid userId)
    {
        throw new NotImplementedException();
    }
}