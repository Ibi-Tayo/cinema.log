using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.Extensions.DependencyInjection;

namespace cinema.log.test.IntegrationTests;

[TestFixture]
public class SoundtrackTests : IntegrationTest
{
    private IFilmService _filmService;
    private IUserService _userService;
    private ISoundtrackService _soundtrackService;
    
    [OneTimeSetUp]
    public void OneTimeSetup()
    {
        var scope = Factory.Services.CreateScope();
        _filmService = scope.ServiceProvider.GetRequiredService<IFilmService>();
        _userService = scope.ServiceProvider.GetRequiredService<IUserService>();
        _soundtrackService = scope.ServiceProvider.GetRequiredService<ISoundtrackService>();
    }
    
    [Test]
    public async Task SoundtrackShouldBeReturned()
    {
        // Arrange: Make film for test (by searching to increase test coverage)
        var results = await _filmService.SearchFilmFromExternal("inception");
        var extId = results.Data
            .FirstOrDefault(q => q.Title.Contains("inception", StringComparison.OrdinalIgnoreCase))
            .ExternalId;
        var film = await _filmService.AddFilmToDb(extId);
        // Arrange: Make user
        var user = await _userService.AddUser(new UserDto()
        {
            Name = "ibitayotest",
            Username = "user",
        });

        // Act
        var resp = 
            await _soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
        
        // Assert
        resp.StatusCode.Should().Be(200);
        resp.StatusMessage.Should().Be("Success");
        resp.Data.Should().NotBeNull();
        resp.Data.Tracks.Count.Should().Be(12);
        resp.Data.SoundtrackName.Should().Be("Inception (Music from the Motion Picture)");
        resp.Data.Artists.Should().Be("Hans Zimmer");
    }

    [Test]
    public async Task SoundTrackReturnedAndIsCorrect()
    {
       // Arrange 
       var film = await _filmService.AddFilmToDb(693134); // Dune part two
       var user = await _userService.AddUser(new UserDto() { Name = "ibitayotest", Username = "user" });
       
       // Act
       var resp = 
           await _soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
       
       // Assert
       resp.StatusCode.Should().Be(200);
       resp.StatusMessage.Should().Be("Success");
       resp.Data.Should().NotBeNull();
       resp.Data.Tracks.Count.Should().Be(25);
       resp.Data.SoundtrackName.Should().Be("Dune: Part Two (Original Motion Picture Soundtrack)");
       resp.Data.Artists.Should().Be("Hans Zimmer");
    }
    
    [Test]
    public async Task SoundTrackNotReturnedBecauseOfficialSoundtrackDoesNotExist()
    {
        // Arrange 
        var film = await _filmService.AddFilmToDb(667); // You only live twice
        var user = await _userService.AddUser(new UserDto() { Name = "ibitayotest", Username = "user" });
       
        // Act
        var resp = 
            await _soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
       
        // Assert
        resp.StatusCode.Should().Be(404);
        resp.StatusMessage.Should().Be("Soundtrack not found in spotify search");
        resp.Data.Should().BeNull();
    }
    
    /// <summary>
    /// The point of this test is to verify that we can return the correct soundtrack when there are -
    /// multiple matches. "The Lion King" matches the 1994 original, 2019 remake and 2024 prequel
    /// </summary>
    [Test]
    public async Task SoundTrackWithMultipleVersionsReturnedIsCorrect()
    {
        // Arrange 
        var film = await _filmService.AddFilmToDb(8587); // lion king 1994
        var film2 = await _filmService.AddFilmToDb(420818); // lion king 2019
        var film3 = await _filmService.AddFilmToDb(762509); // lion king 2024
        var user = await _userService.AddUser(new UserDto() { Name = "ibitayotest", Username = "user" });
       
        // Act
        var resp = 
            await _soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
        var resp2 = 
            await _soundtrackService.GetSoundtrackByFilmId(film2.FilmId, user.Data.UserId);
        var resp3 = 
            await _soundtrackService.GetSoundtrackByFilmId(film3.FilmId, user.Data.UserId);
       
        // Assert
        resp.StatusCode.Should().Be(200);
        resp.Data.Should().NotBeNull();
        resp.Data.Tracks.Count.Should().Be(12);
        resp.Data.SoundtrackName.Should().Be("The Lion King");
        
        resp2.StatusCode.Should().Be(200);
        resp2.Data.Should().NotBeNull();
        resp2.Data.Tracks.Count.Should().Be(19);
        resp2.Data.SoundtrackName.Should().Be("The Lion King (Original Motion Picture Soundtrack)");
        
        resp3.StatusCode.Should().Be(200);
        resp3.Data.Should().NotBeNull();
        resp3.Data.Tracks.Count.Should().Be(7);
        resp3.Data.SoundtrackName.Should().Be("Mufasa: The Lion King (Original Motion Picture Soundtrack)");
        
    }
    
}