using System.ComponentModel.DataAnnotations;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Models.Entities;

public class ComparisonHistory
{
    public Guid ComparisonHistoryId { get; set; }

    public required User User { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    [Required]
    public required Film FilmA { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    [Required]
    public required Film FilmB { get; set; }
    
    [DeleteBehavior(DeleteBehavior.NoAction)]
    public Film? WinningFilm { get; set; }
    
    public DateTime ComparisonDate { get; set; } = DateTime.UtcNow;
    
    [Required]
    public bool WasEqual { get; set; }  // For when user says films are equal
}