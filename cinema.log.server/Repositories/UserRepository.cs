using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Repositories;

public class UserRepository(CinemaLogContext context) : IUserRepository
{
    private CinemaLogContext _context = context;
    
    public async Task<User> CreateUser(User user)
    {
        throw new NotImplementedException();
    }

    public async Task<User> GetUserById(Guid id)
    {
        throw new NotImplementedException();
    }

    public async Task<User> UpdateUser(User user)
    {
        throw new NotImplementedException();
    }

    public async Task<User> DeleteUserById(Guid id)
    {
        throw new NotImplementedException();
    }

    public ICollection<Review> GetUserReviews(Guid userId)
    {
        throw new NotImplementedException();
    }
}