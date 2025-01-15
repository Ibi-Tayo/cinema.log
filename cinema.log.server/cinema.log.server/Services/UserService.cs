using System.Text;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class UserService : IUserService
{
    private IUserRepository _userRepository;
    private IFilmRepository _filmRepository;

    public UserService(IUserRepository userRepository, IFilmRepository filmRepository)
    {
        _userRepository = userRepository;
        _filmRepository = filmRepository;
    }
    public async Task<Response<UserDto>> GetUser(Guid userId)
    {
        try
        {
            var user = await _userRepository.GetUserById(userId);
            if (user == null)
            {
                return Response<UserDto>.BuildResponse(404, "User not found", null);
            }

            var responseUser = Mapper<User, UserDto>.Map(user);
            return Response<UserDto>.BuildResponse(200, "Success", responseUser);
        }
        catch (Exception e)
        {
            return Response<UserDto>.BuildResponse(500, 
                "An error occurred while processing your request", null);
        }
    }
    
    public async Task<Response<UserDto>> AddUser(UserDto user)
    {
        try
        {
            var response = ValidateUser(user);
            if (response.StatusMessage == "Success")
            {
                var newUser = Mapper<UserDto, User>.Map(user);
                var responseUser = await _userRepository.CreateUser(newUser);
                if (responseUser != null)
                {
                    var responseDto = Mapper<User, UserDto>.Map(responseUser);
                    return Response<UserDto>.BuildResponse(StatusCodes.Status201Created, "Success", responseDto);
                }
                return Response<UserDto>.BuildResponse(StatusCodes.Status500InternalServerError, "Internal Server Error", null);
            }
            return response;
        }
        catch (Exception e)
        {
            return Response<UserDto>.BuildResponse(StatusCodes.Status500InternalServerError, "Internal Server Error", null);
        }
    }

    public async Task<Response<UserDto>> UpdateUser(UserDto user)
    {
        var response = ValidateUser(user);
        if (response.StatusMessage == "Success")
        {
            var newUser = Mapper<UserDto, User>.Map(user);
            var responseUser = await _userRepository.UpdateUser(newUser);
            if (responseUser != null)
            {
                var responseDto = Mapper<User, UserDto>.Map(responseUser);
                return Response<UserDto>.BuildResponse(StatusCodes.Status200OK, "Success", responseDto);
            }
            return Response<UserDto>.BuildResponse(StatusCodes.Status500InternalServerError, "Internal Server Error", null);
        }
        return response;
    }

    public async Task<Response<UserDto>> DeleteUser(Guid userId)
    {
        var deletedUser = await _userRepository.DeleteUserById(userId);
        if (deletedUser == null)
        {
            return Response<UserDto>.BuildResponse(StatusCodes.Status404NotFound, "User not found", null);
        }
        
        return Response<UserDto>.BuildResponse(StatusCodes.Status204NoContent, "Success", null);
    }

    public async Task<Response<List<ReviewDto>>> GetUserReviews(Guid userId)
    {
        var userReviews = await _userRepository.GetUserReviews(userId);
        var responseReviews = userReviews
            .Select(Mapper<Review, ReviewDto>.Map)
            .ToList();

        if (responseReviews.Count == 0)
        {
            return Response<List<ReviewDto>>.BuildResponse(StatusCodes.Status404NotFound, "User reviews not found", null);
        }
        
        return Response<List<ReviewDto>>.BuildResponse(StatusCodes.Status200OK, "Success", responseReviews);
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
            return Response<List<FilmDto>>.BuildResponse(StatusCodes.Status404NotFound, "User list of films reviewed not found", null);
        }

        return Response<List<FilmDto>>.BuildResponse(StatusCodes.Status200OK, "Success", userFilmsWatchedResponse);
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

        return new Response<UserDto>() { StatusMessage = "Success" };

        void ValidateString(string value, string fieldName, int minLength, int maxLength, bool isName = false)
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                sb.AppendLine($"{fieldName} is required");
            }
            if (isName && value.Any(char.IsDigit))
            {
                sb.AppendLine($"{fieldName} cannot contain digits");
            }
            if (value.Length < minLength || value.Length > maxLength)
            {
                sb.AppendLine($" {fieldName} must be between {minLength} and {maxLength} characters");
            }
        }
    }

    #endregion
}