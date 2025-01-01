namespace cinema.log.server.Models.Entities;

public class UserFilmRating
{
    public Guid UserFilmRatingId { get; set; }
    public User User { get; set; }
    public required Guid UserId { get; set; }
    public Film Film { get; set; }
    public required Guid FilmId { get; set; }
    public float? EloRating { get; set; }
    public int NumberOfComparisons { get; set; }
    public DateTime LastUpdated { get; set; } = DateTime.UtcNow;
    public float InitialRating { get; set; }
}