namespace cinema.log.server.Models.DTOs;

public class Response<T>
{
    public int? StatusCode { get; set; }
    public string? StatusMessage { get; set; }
    public T? Data { get; set; }

    public static Response<T> BuildResponse(int? statusCode, string statusMessage, T? data)
    {
        return new Response<T>()
        {
            StatusCode = statusCode,
            StatusMessage = statusMessage,
            Data = data
        };
    }
}