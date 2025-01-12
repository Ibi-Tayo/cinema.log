using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Models.Entities;

public class CinemaLogContext(DbContextOptions<CinemaLogContext> options) : DbContext(options)
{
    public DbSet<User> Users { get; set; }
    public DbSet<Film> Films { get; set; }
    public DbSet<Review> Reviews { get; set; }
    public DbSet<UserFilmRating> UserFilmRatings { get; set; }
    public DbSet<ComparisonHistory> ComparisonHistories { get; set; }
    public DbSet<UserFilmSoundtrackRating> UserFilmSoundtrackRatings { get; set; }
    public DbSet<LikedTrack> LikedTracks { get; set; }
    public DbSet<Spotify> SpotifyApi { get; set; }
    
}