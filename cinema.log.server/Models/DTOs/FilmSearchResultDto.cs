namespace cinema.log.server.Models.DTOs;

public class FilmSearchResultDto
{
    public int ExternalId { get; set; }
    public required string Title { get; set; }
    public string? Description { get; set; }
    public int? ReleaseYear { get; set; }
    public string? PosterUrl { get; set; }
}