using System.Text.Json.Serialization;

namespace cinema.log.test.IntegrationTests;

public class TestResponse
{
    [JsonPropertyName("statusCode")]
    public int? StatusCode { get; set; }
    
    [JsonPropertyName("statusMessage")]
    public string? StatusMessage { get; set; }
    
    [JsonPropertyName("data")]
    public object? Data { get; set; }
}