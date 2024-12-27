using System.ComponentModel.DataAnnotations;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Models.Entities;

public class ComparisonHistory
{
    public Guid ComparisonHistoryId { get; set; }

    public User User { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    [Required]
    public Film FilmA { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    [Required]
    public Film FilmB { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    [Required]
    public Film WinningFilm { get; set; }
    
    public DateTime ComparisonDate { get; set; } = DateTime.UtcNow;
    
    [Required]
    public bool WasEqual { get; set; }  // For when user says films are equal
}