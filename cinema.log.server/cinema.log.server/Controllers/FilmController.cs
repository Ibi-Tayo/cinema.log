using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;

[ApiController]
[Route("[controller]")]
public class FilmController(IFilmService filmService) : ControllerBase
{
    [HttpGet]
    [Route("/search")]
    public async Task<ActionResult<Response<List<FilmSearchResultDto>>>> SearchFilm(string query)
    {
        try
        {
            var resp = await filmService.SearchFilmFromExternal(query);
            return Ok(resp);
        }
        catch (HttpRequestException e)
        {
            return StatusCode(500, "Failed to send request to external API");
        }
    }
    
}