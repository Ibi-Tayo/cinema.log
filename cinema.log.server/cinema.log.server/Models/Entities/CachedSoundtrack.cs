namespace cinema.log.server.Models.Entities;

public class CachedSoundtrack
{
    public Guid Id { get; set; }
    public Guid FilmId { get; set; } 
    public string FilmTitle { get; set; }
    public string SpotifyAlbumId { get; set; }
    public string SoundtrackName { get; set; }
    public string Artists { get; set; }
    public string AlbumArtUrl { get; set; }
    public DateTime LastUpdated { get; set; } 
    public string TracksJson { get; set; } // Serialized track data
}