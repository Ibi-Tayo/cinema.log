namespace cinema.log.server.Models.DTOs;

public class CreateReviewDto
{
    public string? Content { get; set; }
    public float Rating { get; set; }
    public Guid FilmId { get; set; }
    public Guid UserId { get; set; }
}