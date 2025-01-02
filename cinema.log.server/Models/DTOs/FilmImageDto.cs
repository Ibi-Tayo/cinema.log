namespace cinema.log.server.Models.DTOs;

public class FilmImageDto
{
    public FilmImageType ImageType { get; set; }
    public float AspectRatio { get; set; }
    public int Height { get; set; }
    public int Width { get; set; }
    public string Url { get; set; }
}


public enum FilmImageType
{
    Backdrop,
    Logo,
    Poster
}