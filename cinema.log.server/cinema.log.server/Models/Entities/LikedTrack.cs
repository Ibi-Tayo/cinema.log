namespace cinema.log.server.Models.Entities;

public class LikedTrack
{
    public Guid Id { get; set; }
    public Guid UserId { get; set; }
    public required string TrackTitle { get; set; }
    public UserFilmSoundtrackRating UserFilmSoundtrackRating { get; set; }
    public Guid UserFilmSoundtrackRatingId { get; set; }
    
}