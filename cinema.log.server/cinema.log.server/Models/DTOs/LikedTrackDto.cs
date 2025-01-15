namespace cinema.log.server.Models.DTOs;

public class LikedTrackDto
{
    public Guid Id { get; set; }
    public Guid UserId { get; set; }
    public required string TrackTitle { get; set; }
    public Guid UserFilmSoundtrackRatingId { get; set; }
}