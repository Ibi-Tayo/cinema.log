using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserService
{
    Task<UserDto> GetUser(Guid userId);
    Task<UserDto> AddUser(UserDto user);
    Task<UserDto> UpdateUser(UserDto user);
    Task<bool> DeleteUser(Guid userId);
    Task<List<ReviewDto>> GetUserReviews(Guid userId);
    Task<List<FilmDto>> GetFilmsReviewedByUser(Guid userId);
}