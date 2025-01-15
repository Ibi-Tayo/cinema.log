using cinema.log.server.Models.Entities;

namespace cinema.log.server.Models.DTOs;

public class UserDto
{
    public Guid UserId { get; set; }
    public required string Name { get; set; }
    public required string Username { get; set; }
    public string? ProfilePicUrl { get; set; }
    
}