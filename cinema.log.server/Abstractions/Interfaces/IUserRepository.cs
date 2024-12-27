using cinema.log.server.Models.Entities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface IUserRepository
{ 
    User CreateUser(User user);
    User GetUserById(Guid id);
    User UpdateUser(User user);
    User DeleteUserById(Guid id);
    ICollection<Review> GetUserReviews(Guid userId);
}