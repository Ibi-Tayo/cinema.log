using System.Text;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Services;

public class UserService : IUserService
{
    private IUserRepository _userRepository;
    private IFilmRepository _filmRepository;
    private ILogger _logger;
    public UserService(IUserRepository userRepository, IFilmRepository filmRepository, ILogger logger)
    {
        _userRepository = userRepository;
        _filmRepository = filmRepository;
        _logger = logger;
    }
    public async Task<Response<UserDto>> GetUser(Guid userId)
    {
        var user = await _userRepository.GetUserById(userId);
        if (user == null)
        {
            return Response<UserDto>.BuildResponse(404, "User not found", null);
        }

        var responseUser = Mapper<User, UserDto>.Map(user);
        return Response<UserDto>.BuildResponse(200, "Success", responseUser);
    }
    
    public async Task<Response<UserDto>> AddUser(UserDto user)
    {
        var response = ValidateUser(user);
        if (response.StatusCode == 200)
        {
            var newUser = Mapper<UserDto, User>.Map(user);
            var responseUser = await _userRepository.CreateUser(newUser);
            if (responseUser != null)
            {
                return Response<UserDto>.BuildResponse(200, "Success", user);
            }
            return Response<UserDto>.BuildResponse(500, "Internal Server Error", null);
        }
        return response;
    }

    public async Task<Response<UserDto>> UpdateUser(UserDto user)
    {
        var response = ValidateUser(user);
        if (response.StatusCode == 200)
        {
            var newUser = Mapper<UserDto, User>.Map(user);
            var responseUser = await _userRepository.UpdateUser(newUser);
            if (responseUser != null)
            {
                return Response<UserDto>.BuildResponse(200, "Success", user);
            }
            return Response<UserDto>.BuildResponse(500, "Internal Server Error", null);
        }
        return response;
    }

    public async Task<Response<UserDto>> DeleteUser(Guid userId)
    {
        var deletedUser = await _userRepository.DeleteUserById(userId);
        if (deletedUser == null)
        {
            return Response<UserDto>.BuildResponse(404, "User not found", null);
        }
        
        return Response<UserDto>.BuildResponse(200, "Success", null);
    }

    public async Task<Response<List<ReviewDto>>> GetUserReviews(Guid userId)
    {
        var userReviews = await _userRepository.GetUserReviews(userId);
        var responseReviews = userReviews
            .Select(Mapper<Review, ReviewDto>.Map)
            .ToList();

        if (responseReviews.Count == 0)
        {
            return Response<List<ReviewDto>>.BuildResponse(404, "User reviews not found", null);
        }
        
        return Response<List<ReviewDto>>.BuildResponse(200, "Success", responseReviews);
    }

    public async Task<Response<List<FilmDto>>> GetFilmsReviewedByUser(Guid userId)
    {
        var userReviews = await _userRepository.GetUserReviews(userId);
        var userFilmsWatchedResponse = new List<FilmDto>();
        foreach (var userReview in userReviews)
        {
            var film = await _filmRepository.GetFilmById(userReview.FilmId);
            if (film != null) userFilmsWatchedResponse.Add(Mapper<Film, FilmDto>.Map(film));
        }
        if (userFilmsWatchedResponse.Count == 0)
        {
            return Response<List<FilmDto>>.BuildResponse(404, "User list of films reviewed not found", null);
        }

        return Response<List<FilmDto>>.BuildResponse(200, "Success", userFilmsWatchedResponse);
    }

    #region Helper methods

    private Response<UserDto> ValidateUser(UserDto user)
    {
        var sb = new StringBuilder();

        ValidateString(user.Username, "Username", 3, 20);
        ValidateString(user.Name, "Name", 3, 40, true);

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

        void ValidateString(string value, string fieldName, int minLength, int maxLength, bool isName = false)
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                sb.Append($"{fieldName} is required");
            }
            if (isName && value.Any(char.IsDigit))
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