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
    public async Task<Response<UserDto>> GetUser(Guid userId) => await userService.GetUser(userId);

    [HttpPost]
    [Route("AddUser")]
    public async Task<Response<UserDto>> AddUser(UserDto newUser) => await userService.AddUser(newUser);

    [HttpPut]
    [Route("UpdateUser")]
    public async Task<ActionResult<Response<UserDto>>> UpdateUser(UserDto existingUser) 
        => await userService.UpdateUser(existingUser);


    [HttpDelete]
    [Route("DeleteUser")]
    public async Task<ActionResult<Response<UserDto>>> DeleteUser(Guid userId) 
        => await userService.DeleteUser(userId);
}