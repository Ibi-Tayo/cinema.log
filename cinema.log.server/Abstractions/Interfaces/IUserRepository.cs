using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserRepository
{ 
    Task<User?> CreateUser(User user);
    Task<User?> GetUserById(Guid id);
    Task<User?> UpdateUser(User user);
    Task<User?> DeleteUserById(Guid id);
}