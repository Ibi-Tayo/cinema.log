using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Repositories;
using cinema.log.server.Utilities;


namespace cinema.log.server.Services;

public class UserFilmRatingService : IUserFilmRatingService
{
    private readonly IUserFilmRatingRepository _userFilmRatingRepository;
    private readonly IUserRepository _userRepository;
    private readonly IFilmRepository _filmRepository;
    private readonly ICalculationService _calculationService;
    private const int NumberOfFilmsForContest = 10;

    public UserFilmRatingService(IUserFilmRatingRepository userFilmRatingRepository,
        ICalculationService calculationService, IUserRepository userRepository, 
        IFilmRepository filmRepository)
    {
        _userFilmRatingRepository = userFilmRatingRepository;
        _calculationService = calculationService;
        _userRepository = userRepository;
        _filmRepository = filmRepository;
    }

    public async Task<Response<UserFilmRatingDto>> GetUserFilmRating(Guid userId, Guid filmId)
    {
        var rating = await _userFilmRatingRepository.GetRatingFilmUserId(userId, filmId);
        if (rating == null) return Response<UserFilmRatingDto>.BuildResponse(404, "Not Found", null);

        var ratingDto = Mapper<UserFilmRating, UserFilmRatingDto>.Map(rating);
        return Response<UserFilmRatingDto>.BuildResponse(200, "Success", ratingDto);
    }

    // When user leaves a review, two calls get made: AddReview (review service) and this (AddUserFilmRating)
    public async Task<Response<UserFilmRatingDto>> AddUserFilmRating(UserFilmRatingDto filmRating)
    {
        var validation = await ValidateRating(filmRating);
        if (validation.StatusCode != 200) return validation;

        var newFilmRating = new UserFilmRating()
        {
            UserFilmRatingId = Guid.NewGuid(),
            UserId = filmRating.UserId,
            FilmId = filmRating.FilmId,
            EloRating = GetInitialEloRating(filmRating.InitialRating),
            InitialRating = filmRating.InitialRating
        };
        var addedFilmRating = await _userFilmRatingRepository.CreateRating(newFilmRating);
        if (addedFilmRating == null)
            return Response<UserFilmRatingDto>.BuildResponse(500, "Internal server error", null);
        var ratingDto = Mapper<UserFilmRating, UserFilmRatingDto>.Map(addedFilmRating);
        return Response<UserFilmRatingDto>.BuildResponse(200, "Success", ratingDto);
    }
    
    // Client sends id of film that we want to match up against other films
    public async Task<Response<List<FilmDto>>> GetFilmsForContest(Guid userId, Guid contestFilmId)
    {
        var contestFilm = await _filmRepository.GetFilmById(contestFilmId);
        if (contestFilm == null) return Response<List<FilmDto>>.BuildResponse(404, "Film Id for contest not found", null);
        
        var allFilms = await _userFilmRatingRepository.GetAllFilmsRatedByUserId(userId);
        var allRatings = await _userFilmRatingRepository.GetAllRatings(userId);
        
        if (allFilms.Count == 0) return Response<List<FilmDto>>.BuildResponse(404, "No rated films found", null);
        
        // Here we implement some logic to try and prioritise by genre using the contest film information
        var contestFilmGenres = contestFilm.Genre?.Split(',');
        var sortedFilms = AssignFilmPriority(allFilms, allRatings, contestFilmGenres).Take(NumberOfFilmsForContest);
        
        return Response<List<FilmDto>>.BuildResponse(200, "Success", 
            sortedFilms.Select(Mapper<Film, FilmDto>.Map).ToList());
    }

    public async Task<Response<(UserFilmRatingDto?, UserFilmRatingDto?)>> FilmContest(Guid userId,
        Guid filmA, Guid filmB,
        Guid winnerId)
    {
        // Define result of contest based on winnerId 
        DefineFilmContestResult(filmA, filmB, winnerId, out var filmAResult, out var filmBResult);

        // Use UserId with Film A and B to get the UserFilmRating from table
        var filmARating = await _userFilmRatingRepository.GetRatingFilmUserId(userId, filmA);
        var filmBRating = await _userFilmRatingRepository.GetRatingFilmUserId(userId, filmB);

        if (filmARating == null || filmBRating == null)
            return Response<(UserFilmRatingDto, UserFilmRatingDto)>
                .BuildNullableResponse<UserFilmRatingDto, UserFilmRatingDto>(404,
                "Film Rating not found for one or both films", (null, null));

        // Calculate expected result for film A and film B
        var filmAExpectedResult = _calculationService
            .CalculateExpectedResult(filmARating.EloRating, filmBRating.EloRating);
        var filmBExpectedResult = _calculationService
            .CalculateExpectedResult(filmBRating.EloRating, filmARating.EloRating);

        // Check constant
        filmARating.KConstantValue = VerifyKConstantValue(filmARating);
        filmBRating.KConstantValue = VerifyKConstantValue(filmBRating);

        // Recalculate film rating for film A and film B
        var filmANewRating = _calculationService
            .RecalculateFilmRating(filmAExpectedResult, filmAResult, filmARating.EloRating, filmARating.KConstantValue);
        var filmBNewRating = _calculationService
            .RecalculateFilmRating(filmBExpectedResult, filmBResult, filmBRating.EloRating, filmBRating.KConstantValue);
        
        // Update film A
        filmARating.EloRating = filmANewRating;
        filmARating.LastUpdated = DateTime.Now;
        filmARating.NumberOfComparisons += 1;

        // Update film B
        filmBRating.EloRating = filmBNewRating;
        filmBRating.LastUpdated = DateTime.Now;
        filmBRating.NumberOfComparisons += 1;

        var resA = await _userFilmRatingRepository.UpdateRating(filmARating);
        var resB = await _userFilmRatingRepository.UpdateRating(filmBRating);

        if (resA == null || resB == null)
            return Response<(UserFilmRatingDto, UserFilmRatingDto)>
                .BuildNullableResponse<UserFilmRatingDto, UserFilmRatingDto>(500,
                "Internal server error, Film ratings couldn't get updated", (null, null));

        var resADto = Mapper<UserFilmRating, UserFilmRatingDto>.Map(resA);
        var resBDto = Mapper<UserFilmRating, UserFilmRatingDto>.Map(resB);

        return Response<(UserFilmRatingDto?, UserFilmRatingDto?)>
            .BuildResponse(200, "Success", (resADto, resBDto));
    }

    public async Task<Response<bool>> ResetAllRatings(Guid userId)
    {
        var allRatings = await _userFilmRatingRepository.GetAllRatings(userId);
        if (allRatings.Count == 0) return Response<bool>.BuildResponse(404, "No ratings found", false);

        foreach (var rating in allRatings)
        {
            rating.EloRating = GetInitialEloRating(rating.InitialRating);
            rating.NumberOfComparisons = 0;
            rating.LastUpdated = DateTime.Now;
        }
        return Response<bool>.BuildResponse(200, "Success", true);
    }

    public async Task<Response<bool>> DeleteRatingByUserAndFilmId(Guid userId, Guid filmId)
    {
        var res = await _userFilmRatingRepository.DeleteRatingByUserFilmId(userId, filmId);
        if (!res) return Response<bool>.BuildResponse(404, "No rating found", false);
        return Response<bool>.BuildResponse(200, "Success", true);
    }
    
    #region Helper Methods

    private async Task<Response<UserFilmRatingDto>> ValidateRating(UserFilmRatingDto filmRating)
    {
        // Check if user exists
        var user = await _userRepository.GetUserById(filmRating.UserId);
        if (user == null) return Response<UserFilmRatingDto>.BuildResponse(404, "User not found", null);
        // Check if film exists
        var film = await _filmRepository.GetFilmById(filmRating.FilmId);
        if (film == null) return Response<UserFilmRatingDto>.BuildResponse(404, "Film not found", null);
        // Check if InitialRating is between 0-5
        if (filmRating.InitialRating is < 0 or > 5)
            return Response<UserFilmRatingDto>.BuildResponse(400, "Initial rating must be between 0 and 5", null);
        return Response<UserFilmRatingDto>.BuildResponse(200, "Success", filmRating);
    }

    private float GetInitialEloRating(float rating)
    {
        return rating switch
        {
            >= 0 and < 2 => 1400,
            >= 2 and < 3 => 1500,
            >= 3 and < 4 => 1600,
            >= 4 and <= 5 => 1700,
            _ => throw new ArgumentOutOfRangeException(nameof(rating), "Rating must be between 0 and 5"),
        };
    }

    private void DefineFilmContestResult(Guid filmA, Guid filmB, Guid winnerId,
        out float filmAResult,
        out float filmBResult)
    {
        if (filmA == winnerId)
        {
            filmAResult = 1;
            filmBResult = 0;
        }
        else if (filmB == winnerId)
        {
            filmAResult = 0;
            filmBResult = 1;
        }
        else
        {
            filmAResult = 0.5f;
            filmBResult = 0.5f;
        }
    }

    private List<Film> AssignFilmPriority(IEnumerable<Film> films, IEnumerable<UserFilmRating> ratings, string[]? contestFilmGenres)
    {
        var r = new Random();
        var filmsArray = films.ToArray();
        r.Shuffle(filmsArray); 
        
        if (contestFilmGenres == null || contestFilmGenres.Length == 0)
        {
            // No priority so return the list with the least compared films at the start of the list
            return filmsArray
                .OrderBy(film => ratings
                    .Where(ufr => ufr.FilmId == film.FilmId)
                    .Sum(u => u.NumberOfComparisons))
                .ToList();
        }
        
        var normalizedContestFilmGenres = contestFilmGenres
            .Select(g => g.Trim().ToLowerInvariant())
            .ToHashSet();
        
        var filmPriority = new List<(Film, int)>();
        
        // Go through each film and get the array of genres
        foreach (var film in filmsArray)
        {
            if (film.Genre == null)
            {
                filmPriority.Add((film, 0));
                continue;
            }
            
            // Count how many genres match with the contestFilmGenres Array
            var count = film.Genre
                .Split(',')
                .Select(g => g.Trim().ToLowerInvariant())
                .Count(genre => normalizedContestFilmGenres.Contains(genre));
            
            filmPriority.Add((film, count));
        }
        
        // Sort by int (descending order) to prioritise genre
        // But then also sort to prioritise the least compared films
        var sort = filmPriority
            .OrderByDescending(x => x.Item2)
            .ThenBy(film => ratings
                .Where(ufr => ufr.FilmId == film.Item1.FilmId)
                .Sum(u => u.NumberOfComparisons))
            .Select(x => x.Item1)
            .ToList();

        return sort;
    }

    private double VerifyKConstantValue(UserFilmRating filmRating)
    {
        var numberOfComparisons = filmRating.NumberOfComparisons;
        return numberOfComparisons switch
        {
            >= 0 and < 5 => 40,
            >= 5 and < 10  => 20,
            >= 10 => 10,
            _ => 40,
        };
        
    }
    #endregion
}