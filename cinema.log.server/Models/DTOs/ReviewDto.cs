namespace cinema.log.server.Models.DTOs;

public class ReviewDto
{
    public Guid ReviewId { get; set; }
    public string? Content { get; set; }
    public DateTime Date { get; set; }
    public float Rating { get; set; }
    public Guid FilmId { get; set; }
    public Guid UserId { get; set; }
}