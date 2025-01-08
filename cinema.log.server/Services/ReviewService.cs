using System.Text;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class ReviewService : IReviewService
{
    private readonly IReviewRepository _reviewRepository;
    private readonly IFilmRepository _filmRepository;
    private readonly IUserRepository _userRepository;

    public ReviewService(IReviewRepository reviewRepository, IFilmRepository filmRepository,
        IUserRepository userRepository)
    {
        _reviewRepository = reviewRepository;
        _filmRepository = filmRepository;
        _userRepository = userRepository;
    }

    public async Task<Response<ReviewDto>> GetReview(Guid reviewId)
    {
        var review = await _reviewRepository.GetReviewById(reviewId);
        if (review == null)
        {
            return Response<ReviewDto>.BuildResponse(404, "Review not found", null);
        }

        var responseReview = Mapper<Review, ReviewDto>.Map(review);
        return Response<ReviewDto>.BuildResponse(200, "Success", responseReview);
    }

    public async Task<Response<List<ReviewDto>>> GetReviewsByUserId(Guid userId)
    {
        var reviews = await _reviewRepository.GetReviewsByUserId(userId);
     
        var reviewDtos = reviews
            .Select(Mapper<Review, ReviewDto>.Map)
            .ToList();

        if (reviewDtos.Count == 0)
        {
            return Response<List<ReviewDto>>.BuildResponse(
                StatusCodes.Status404NotFound,
                "No reviews found for this user",
                null);
        }

        return Response<List<ReviewDto>>.BuildResponse(
            StatusCodes.Status200OK,
            "Success",
            reviewDtos);
    }
    public async Task<Response<List<ReviewDto>>> GetReviewsByFilmId(Guid filmId)
    {
        var film = await _filmRepository.GetFilmById(filmId);
        if (film == null)
        {
            return Response<List<ReviewDto>>.BuildResponse(
                StatusCodes.Status404NotFound,
                "Film not found",
                null);
        }

        var reviews = await _reviewRepository.GetReviewsByFilmId(filmId);
        var reviewDtos = reviews
            .Select(Mapper<Review, ReviewDto>.Map)
            .ToList();

        if (reviewDtos.Count == 0)
        {
            return Response<List<ReviewDto>>.BuildResponse(
                StatusCodes.Status404NotFound,
                "No reviews found for this film",
                null);
        }

        return Response<List<ReviewDto>>.BuildResponse(
            StatusCodes.Status200OK,
            "Success",
            reviewDtos);
    }

    public async Task<Response<ReviewDto>> GetReviewByUserAndFilm(Guid userId, Guid filmId)
    {
        var review = await _reviewRepository.GetReviewByUserAndFilm(userId, filmId);
        if (review == null) return Response<ReviewDto>.BuildResponse(404, "Review not found", null);
        var responseReview = Mapper<Review, ReviewDto>.Map(review);
        return Response<ReviewDto>.BuildResponse(200, "Success", responseReview);
    }

    public async Task<Response<ReviewDto>> AddReview(ReviewDto review)
    {
        var response = await ValidateReview(review, true);
        if (response.StatusMessage != "Success")
        {
            return response;
        }

        var newReview = Mapper<ReviewDto, Review>.Map(review);
        var createdReview = await _reviewRepository.CreateReview(newReview);

        if (createdReview == null)
        {
            return Response<ReviewDto>.BuildResponse(
                StatusCodes.Status500InternalServerError,
                "Internal Server Error",
                null);
        }

        var responseDto = Mapper<Review, ReviewDto>.Map(createdReview);
        return Response<ReviewDto>.BuildResponse(
            StatusCodes.Status201Created,
            "Success",
            responseDto);
    }

    public async Task<Response<ReviewDto>> UpdateReview(ReviewDto review)
    {
        var response = await ValidateReview(review);
        if (response.StatusMessage != "Success")
        {
            return response;
        }

        var existingReview = await _reviewRepository.GetReviewById(review.ReviewId);
        if (existingReview == null)
        {
            return Response<ReviewDto>.BuildResponse(
                StatusCodes.Status404NotFound,
                "Review not found",
                null);
        }

        var updatedReview = Mapper<ReviewDto, Review>.Map(review);
        var result = await _reviewRepository.UpdateReview(updatedReview);

        if (result == null)
        {
            return Response<ReviewDto>.BuildResponse(
                StatusCodes.Status500InternalServerError,
                "Internal Server Error",
                null);
        }

        var responseDto = Mapper<Review, ReviewDto>.Map(result);
        return Response<ReviewDto>.BuildResponse(
            StatusCodes.Status200OK,
            "Success",
            responseDto);
    }

    public async Task<Response<bool>> DeleteReview(Guid reviewId)
    {
        var deletedReview = await _reviewRepository.DeleteReviewById(reviewId);
        if (!deletedReview)
        {
            return Response<bool>.BuildResponse(
                StatusCodes.Status404NotFound,
                "Review not found",
                false);
        }

        return Response<bool>.BuildResponse(
            StatusCodes.Status204NoContent,
            "Success",
            true);
    }

    #region Helper Methods
    private async Task<Response<ReviewDto>> ValidateReview(ReviewDto review, bool newReview = false)
    {
        var sb = new StringBuilder();

        // Check rating
        if (review.Rating < 0 || review.Rating > 5)
        {
            sb.AppendLine("Rating must be between 1 and 5,");
        }

        // Check content length
        if (!string.IsNullOrEmpty(review.Content) && review.Content.Length > 3000)
        {
            sb.AppendLine("Review text cannot exceed 1000 characters,");
        }

        // Check film exists
        var film = await _filmRepository.GetFilmById(review.FilmId);
        if (film == null)
        {
            sb.AppendLine("Film does not exist");
        }

        // Check user exists
        var user = await _userRepository.GetUserById(review.UserId);
        if (user == null)
        {
            sb.AppendLine("User does not exist");
        }

        // Check for duplicate review
        if (newReview)
        {
            var existingReview = await _reviewRepository.GetReviewByUserAndFilm(review.UserId, review.FilmId);
            if (existingReview != null)
            {
                sb.AppendLine("User has already reviewed this film");
            }
        }

        if (sb.Length > 0)
        {
            return Response<ReviewDto>.BuildResponse(
                StatusCodes.Status400BadRequest,
                sb.ToString(),
                null);
        }

        return Response<ReviewDto>.BuildResponse(
            StatusCodes.Status200OK,
            "Success",
            null);
    }
    #endregion
}