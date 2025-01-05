namespace cinema.log.server.Utilities;

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
    
    public static Response<(T1?, T2?)> BuildNullableResponse<T1, T2>(int statusCode, string message, (T1?, T2?) data)
    {
        return new Response<(T1?, T2?)>
        {
            StatusCode = statusCode, 
            StatusMessage= message,
            Data = data
        };
    }
}