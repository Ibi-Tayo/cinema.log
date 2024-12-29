using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Repositories;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;

[ApiController]
[Route("[controller]")]
public class TestController : ControllerBase
{
    
    IFilmRepository _filmRepository;
    CinemaLogContext _context;

    public TestController(IFilmRepository repo, CinemaLogContext context)
    {
        _filmRepository = repo;
        _context = context;
    }
    
    [HttpPost]
    [Route("AddFilm")]
    public async Task<ActionResult<Film>> AddFilm()
    {
        var film = new Film()
        {
            Title = "New Film",
            Description = "a new film coming out in 2026",
            Director = "Ibitayo",
        };
        var newFilm = await _filmRepository.CreateFilm(film);
        return Ok(newFilm);

    }


}