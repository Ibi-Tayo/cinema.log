using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;

[ApiController]
[Route("[controller]")]
public class TestController : ControllerBase
{
    
    IUserService _userService;

    public TestController(IUserService userService)
    {
        _userService = userService;
    }

    [HttpGet]
    [Route("{userId}")]
    public async Task<ActionResult<Response<UserDto>>> GetUser(Guid userId)
    {
        var resp = await _userService.GetUser(userId);
        return resp.StatusCode switch
        {
            StatusCodes.Status200OK => Ok(resp),
            StatusCodes.Status404NotFound => NotFound(resp),
            _ => StatusCode(StatusCodes.Status500InternalServerError, resp)
        };
    }
    
    [HttpPost]
    [Route("AddUser")]
    public async Task<ActionResult<Response<UserDto>>> AddUser(UserDto newUser)
    {
        var resp = await _userService.AddUser(newUser);
        return resp.StatusCode switch
        {
            StatusCodes.Status201Created => CreatedAtAction("AddUser", resp),
            StatusCodes.Status400BadRequest => BadRequest(resp),
            _ => StatusCode(StatusCodes.Status500InternalServerError, resp)
        };
    }
    
    // this updated user needs an existing id, or else we wouldnt know what user to update
    [HttpPut]
    [Route("UpdateUser")]
    public async Task<ActionResult<Response<UserDto>>> UpdateUser(UserDto existingUser)
    {
        var resp = await _userService.UpdateUser(existingUser);
        return resp.StatusCode switch
        {
            StatusCodes.Status200OK => Ok(resp),
            StatusCodes.Status400BadRequest => BadRequest(resp),
            _ => StatusCode(StatusCodes.Status500InternalServerError, resp)
        };
    }

    [HttpDelete]
    [Route("DeleteUser")]
    public async Task<ActionResult<Response<UserDto>>> DeleteUser(Guid userId)
    {
        var response = await _userService.DeleteUser(userId);
        return response.StatusCode switch
        {
            StatusCodes.Status204NoContent => StatusCode(204, response),
            StatusCodes.Status404NotFound => NotFound(response),
            _ => StatusCode(StatusCodes.Status500InternalServerError, response)
        };
    }
}