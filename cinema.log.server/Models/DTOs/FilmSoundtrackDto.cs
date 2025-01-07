using cinema.log.server.Models.Entities;

namespace cinema.log.server.Models.DTOs;

public class FilmSoundtrackDto
{
    public Guid UserFilmSoundtrackRatingId { get; set; }
    public Guid FilmId { get; set; }
    public string FilmName { get; set; }
    public string SoundtrackName { get; set; }
    public string AlbumArtUrl { get; set; }
    public List<TrackDto> Tracks { get; set; }
    public List<LikedTrackDto> LikedTracks { get; set; }
}