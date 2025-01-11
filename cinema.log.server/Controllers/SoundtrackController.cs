using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;


[ApiController]
[Route("[controller]")]
public class SoundtrackController(ISoundtrackService soundtrackService) : ControllerBase
{
    [HttpGet]
    [Route("/test")]
    public async Task<ActionResult> Test()
    {
        var res = await soundtrackService.GetSoundtrackByFilmId(Guid.Empty);
        return Ok(res);
    }
    
}