using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;


[ApiController]
[Route("[controller]")]
public class SoundtrackController(ISoundtrackService soundtrackService, 
    IFilmService filmService, IUserService userService) : ControllerBase
{
    [HttpGet]
    [Route("/soundtrack-integration-test")]
    public async Task<Response<FilmSoundtrackDto>> CanGetSoundtrackFromFilmId()
    {
        // make film for test (by searching to increase test coverage)
        var results = await filmService.SearchFilmFromExternal("inception");
        var extId = results.Data
            .FirstOrDefault(q => q.Title.Contains("inception", StringComparison.OrdinalIgnoreCase))
            .ExternalId;
        var film = await filmService.AddFilmToDb(extId);
        // make user
        var user = await userService.AddUser(new UserDto()
        {
            Name = "ibitayotest",
            Username = "user",
        });

        // test
        var resp = await soundtrackService.GetSoundtrackByFilmId(film.FilmId, user.Data.UserId);
        return resp;
    }
    
}