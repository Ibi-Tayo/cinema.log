namespace cinema.log.server.Models.DTOs;

public class UserFilmRatingDto
{
    public required Guid UserId { get; set; }
    public required Guid FilmId { get; set; }
    public float? EloRating { get; set; }
    public float InitialRating { get; set; }
}