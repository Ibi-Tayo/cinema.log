using System.ComponentModel.DataAnnotations;

namespace cinema.log.server.Models;

public class User
{
    [Required]
    public Guid UserId { get; set; }
    
    [Required]
    [MaxLength(40, ErrorMessage = "Name is too long")]
    public required string Name { get; set; }
    
    [Required]
    [MaxLength(20, ErrorMessage = "Username is too long")]
    public required string Username { get; set; }
    
    [Required]
    [MaxLength(500, ErrorMessage = "Url is too long")]
    public string? ProfilePicUrl { get; set; }
    
    public ICollection<Review>? Reviews { get; set; }
}