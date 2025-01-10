namespace cinema.log.server.Models.Entities;

public class Spotify
{
    public Guid Id { get; set; }
    public string AccessToken { get; set; }
    public DateTime ExpiryDate { get; set; }
}