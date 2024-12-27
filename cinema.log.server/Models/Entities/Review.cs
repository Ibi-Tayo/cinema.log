using System.ComponentModel.DataAnnotations;

namespace cinema.log.server.Models.Entities;

public class Review
{
    public Guid ReviewId { get; set; }
    
    [MaxLength(3000)]
    public string? Content { get; set; }
    
    [Required]
    public DateTime Date { get; set; }
    
    [Required]
    public float Rating { get; set; }
    
    [Required]
    public Film Film { get; set; }
    
    [Required]
    public User User { get; set; }
}