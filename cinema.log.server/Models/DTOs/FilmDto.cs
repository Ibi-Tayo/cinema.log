namespace cinema.log.server.Models.DTOs;

public class FilmDto
{
    public Guid FilmId { get; set; }
    public required string Title { get; set; }
    public string? Description { get; set; }
    public string? Genre { get; set; }
    public string? Director { get; set; }
    public int? ReleaseYear { get; set; }
    public string? PosterUrl { get; set; }
}