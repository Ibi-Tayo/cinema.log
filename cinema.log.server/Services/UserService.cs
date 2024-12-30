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

    public async Task<UserDto> UpdateUser(UserDto user)
    {
        // validation
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
        var username = user.Username;
        var name = user.Name;

        if (string.IsNullOrWhiteSpace(username))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Username is required",
            };
        }
        if (username.Any(char.IsPunctuation))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Username cannot contain punctuation",
            };
        }
        if (username.Any(char.IsDigit))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Username cannot contain digits",
            };
        }
        if (username.Length is > 20 or < 3)
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Username must be between 3 and 20 characters",
            };
        }

        if (string.IsNullOrWhiteSpace(name))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Name is required",
            };
        }
        if (name.Any(char.IsPunctuation))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Name cannot contain punctuation",
            };
        }
        if (name.Any(char.IsDigit))
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Name cannot contain digits",
            };
        }
        if (name.Length is > 40 or < 3)
        {
            return new Response<UserDto>()
            {
                StatusCode = 400,
                StatusMessage = "Name must be between 3 and 40 characters",
            };
        }

        return new Response<UserDto>()
        {
            StatusCode = 200,
            StatusMessage = "Success",
        };
    }
    #endregion
    

}