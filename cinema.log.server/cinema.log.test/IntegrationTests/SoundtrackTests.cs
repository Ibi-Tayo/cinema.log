using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.Extensions.DependencyInjection;

namespace cinema.log.test.IntegrationTests;

[TestFixture]
public class SoundtrackTests : WebApplicationFactory<IApiMarker>
{
    private readonly HttpClient _httpClient;
    private readonly IFilmService filmService;
    private readonly IUserService userService;
    private readonly ISoundtrackService soundtrackService;
    
    public SoundtrackTests()
    {
        var appFactory = new WebApplicationFactory<IApiMarker>();
        _httpClient = appFactory.CreateClient();
        filmService = appFactory.Services.GetRequiredService<IFilmService>();
        userService = appFactory.Services.GetRequiredService<IUserService>();
        soundtrackService = appFactory.Services.GetRequiredService<ISoundtrackService>();
    }
    
    [Test]
    public async Task CanGetSoundtrackFromFilmId()
    {
        // Arrange: Make film for test (by searching to increase test coverage)
        var results = await filmService.SearchFilmFromExternal("inception");
        var extId = results.Data
            .FirstOrDefault(q => q.Title.Contains("inception", StringComparison.OrdinalIgnoreCase))
            .ExternalId;
        var film = await filmService.AddFilmToDb(extId);
        // Arrange: Make user
        var user = await userService.AddUser(new UserDto()
        {
            Name = "ibitayotest",
            Username = "user",
        });

        // Act
        var resp = await soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
        
        // Assert
        resp.StatusCode.Should().Be(200);
        resp.StatusMessage.Should().Be("Success");
        resp.Data.Should().NotBeNull();
    }
    
}