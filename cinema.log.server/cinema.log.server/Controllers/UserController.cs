using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;
using Microsoft.AspNetCore.Mvc;

namespace cinema.log.server.Controllers;

[ApiController]
[Route("[controller]")]
public class UserController(IUserService userService) : ControllerBase
{
    [HttpGet]
    [Route("{userId}")]
    public async Task<ActionResult<Response<UserDto>>> GetUser(Guid userId)
    {
        var resp = await userService.GetUser(userId);
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
        var resp = await userService.AddUser(newUser);
        return resp.StatusCode switch
        {
            StatusCodes.Status201Created => CreatedAtAction("AddUser", resp),
            StatusCodes.Status400BadRequest => BadRequest(resp),
            _ => StatusCode(StatusCodes.Status500InternalServerError, resp)
        };
    }

    [HttpPut]
    [Route("UpdateUser")]
    public async Task<ActionResult<Response<UserDto>>> UpdateUser(UserDto existingUser)
    {
        var resp = await userService.UpdateUser(existingUser);
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
        var response = await userService.DeleteUser(userId);
        return response.StatusCode switch
        {
            StatusCodes.Status204NoContent => StatusCode(204, response),
            StatusCodes.Status404NotFound => NotFound(response),
            _ => StatusCode(StatusCodes.Status500InternalServerError, response)
        };
    }
}