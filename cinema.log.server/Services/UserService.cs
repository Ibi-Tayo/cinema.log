using System.Text;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Services;

public class UserService : IUserService
{
    private IUserRepository _userRepository;
    private ILogger _logger;
    public UserService(IUserRepository userRepository, ILogger logger)
    {
        _userRepository = userRepository;
        _logger = logger;
    }
    public async Task<UserDto?> GetUser(Guid userId)
    {
        var user = await _userRepository.GetUserById(userId);
        if (user == null) return null;
        var responseUser = MapUserToDto(user);
        return responseUser;
    }
    
    public async Task<Response<UserDto>> AddUser(UserDto user)
    {
        var response = ValidateUser(user);
        if (response.StatusCode == 200)
        {
            var newUser = MapDtoToUser(user);
            var responseUser = await _userRepository.CreateUser(newUser);
            if (responseUser != null)
            {
                response.Data = user;
            }
            else
            {
                response.StatusCode = 500;
                response.StatusMessage = "Internal Server Error";
            }
        }
        return response;
    }

    public async Task<Response<UserDto>> UpdateUser(UserDto user)
    {
        var response = ValidateUser(user);
        if (response.StatusCode == 200)
        {
            var newUser = MapDtoToUser(user);
            var responseUser = await _userRepository.UpdateUser(newUser);
            if (responseUser != null)
            {
                response.Data = user;
            }
            else
            {
                response.StatusCode = 500;
                response.StatusMessage = "Internal Server Error";
            }
        }
        return response;
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

    #region Helper methods
    private UserDto MapUserToDto(User user)
    {
        var responseUser = new UserDto()
        {
            UserId = user.UserId,
            Name = user.Name,
            Username = user.Username,
            ProfilePicUrl = user.ProfilePicUrl,
        };
        return responseUser;
    }
    private User MapDtoToUser(UserDto user)
    {
        var newUser = new User()
        {
            UserId = user.UserId,
            Name = user.Name,
            Username = user.Username,
            ProfilePicUrl = user.ProfilePicUrl,
        };
        return newUser;
    }

    private Response<UserDto> ValidateUser(UserDto user)
    {
        var sb = new StringBuilder();

        ValidateString(user.Username, "Username", 3, 20);
        ValidateString(user.Name, "Name", 3, 40);

        if (sb.Length > 0)
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = sb.ToString(),
            };
        }

        return new Response<UserDto>()
        {
            StatusCode = 200,
            StatusMessage = "Success",
        };

        void ValidateString(string value, string fieldName, int minLength, int maxLength)
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                sb.Append($"{fieldName} is required");
            }
            if (value.Any(char.IsPunctuation))
            {
                sb.Append($" {fieldName} cannot contain punctuation");
            }
            if (value.Any(char.IsDigit))
            {
                sb.Append($" {fieldName} cannot contain digits");
            }
            if (value.Length < minLength || value.Length > maxLength)
            {
                sb.Append($" {fieldName} must be between {minLength} and {maxLength} characters");
            }
        }
    }

    #endregion
    

}