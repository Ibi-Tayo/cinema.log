using cinema.log.server.Models.DTOs;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserService
{
    Task<Response<UserDto>> GetUser(Guid userId);
    Task<Response<UserDto>> AddUser(UserDto user);
    Task<Response<UserDto>> UpdateUser(UserDto user);
    Task<Response<UserDto>> DeleteUser(Guid userId);
    Task<Response<List<ReviewDto>>> GetUserReviews(Guid userId);
    Task<Response<List<FilmDto>>> GetFilmsReviewedByUser(Guid userId);
}