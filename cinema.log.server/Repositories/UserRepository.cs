using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.Entities;

namespace cinema.log.server.Repositories;

public class UserRepository : IUserRepository
{
    public User CreateUser(User user)
    {
        throw new NotImplementedException();
    }

    public User GetUserById(Guid id)
    {
        throw new NotImplementedException();
    }

    public User UpdateUser(User user)
    {
        throw new NotImplementedException();
    }

    public User DeleteUserById(Guid id)
    {
        throw new NotImplementedException();
    }

    public ICollection<Review> GetUserReviews(Guid userId)
    {
        throw new NotImplementedException();
    }
}