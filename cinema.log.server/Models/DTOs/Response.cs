namespace cinema.log.server.Models.DTOs;

public class Response<T>
{
    public int StatusCode { get; set; }
    public string? StatusMessage { get; set; }
    public T? Data { get; set; }
}