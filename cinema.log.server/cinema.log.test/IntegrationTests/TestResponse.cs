using System.Text.Json.Serialization;

namespace cinema.log.test.IntegrationTests;

public class TestResponse<T>
{
    [JsonPropertyName("statusCode")]
    public int? StatusCode { get; set; }
    
    [JsonPropertyName("statusMessage")]
    public string? StatusMessage { get; set; }
    
    [JsonPropertyName("data")]
    public T? Data { get; set; }
}