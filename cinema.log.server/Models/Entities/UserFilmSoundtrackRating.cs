namespace cinema.log.server.Models.Entities;

public class UserFilmSoundtrackRating
{
    public Guid UserFilmSoundtrackRatingId { get; set; }
    public Guid FilmId { get; set; }
    public Guid UserId { get; set; }
    public int Rating { get; set; }
    public ICollection<LikedTrack> LikedTracks { get; set; }
}