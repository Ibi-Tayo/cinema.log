using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Repositories;

public class UserRepository(CinemaLogContext context, ILogger<UserRepository> logger) : IUserRepository
{
    private CinemaLogContext _context = context;
    private ILogger<UserRepository> _logger = logger;
    
    public async Task<User?> CreateUser(User user)
    {
        try
        {
            await _context.Users.AddAsync(user);
            await _context.SaveChangesAsync();
        }
        catch (Exception e)
        { 
            _logger.LogError(e, e.Message);
            return null;
        }
        return user;
    }

    public async Task<User?> GetUserById(Guid id)
    {
        return await _context.Users.FindAsync(id);
    }

    public async Task<User?> UpdateUser(User user)
    {
        try
        {
            _context.Users.Update(user);
            await _context.SaveChangesAsync();
        }
        catch (Exception e)
        {
            _logger.LogError(e, e.Message);
            return null;
        }
        return user;
    }

    public async Task<User?> DeleteUserById(Guid id)
    {
        var foundUser = await _context.Users.FindAsync(id);
        if (foundUser == null) return null;
        _context.Users.Remove(foundUser);
        await _context.SaveChangesAsync();
        return foundUser;
    }

    public async Task<List<Review>> GetUserReviews(Guid userId)
    {
        var reviews = await _context.Reviews
            .Where(review => review.User.UserId == userId)
            .ToListAsync();

        return reviews;
    }
}