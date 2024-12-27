using System.ComponentModel.DataAnnotations;

namespace cinema.log.server.Models.Entities;

public class Film
{
    [Required]
    public Guid FilmId { get; set; }
    
    [Required]
    public string Title { get; set; }
    
    public string? Description { get; set; }
    
    public string? Genre { get; set; }
    
    public string? Director { get; set; }
    
    public string? PosterUrl { get; set; }
    
    public ICollection<Review>? Reviews { get; set; }
}