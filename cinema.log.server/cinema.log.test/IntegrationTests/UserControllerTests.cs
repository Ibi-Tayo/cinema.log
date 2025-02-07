using System.Net;
using System.Text.Json;
using cinema.log.server.Abstractions.Interfaces;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
namespace cinema.log.test.IntegrationTests;

[TestFixture]
public class UserControllerTests : WebApplicationFactory<IApiMarker>
{
    private readonly HttpClient _httpClient;
    
    public UserControllerTests()
    {
        var appFactory = new WebApplicationFactory<IApiMarker>();
        _httpClient = appFactory.CreateClient();
    }

    [Test]
    public async Task DeleteUserWithInvalidIdReturns404()
    {
        var id = Guid.NewGuid();
        var response = await _httpClient.DeleteAsync($"/user/deleteuser?userId={id}");
        var responseContent = await response.Content.ReadAsStringAsync();
        var json = JsonSerializer.Deserialize<TestResponse>(responseContent);
        
        response.StatusCode.Should().Be(HttpStatusCode.NotFound);
        json.Data.Should().BeNull();
        json.StatusMessage.Should().Be("User not found");
    }
    
}