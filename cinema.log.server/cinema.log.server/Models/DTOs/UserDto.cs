using System.Text.Json.Serialization;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Models.DTOs;

public class UserDto
{
    [JsonPropertyName("userId")]
    public Guid UserId { get; set; }
    [JsonPropertyName("name")]
    public required string Name { get; set; }
    [JsonPropertyName("username")]
    public required string Username { get; set; }
    [JsonPropertyName("profilePicUrl")]
    public string? ProfilePicUrl { get; set; }
    
}