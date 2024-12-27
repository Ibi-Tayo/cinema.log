using cinema.log.server.Models.Entities;
using Microsoft.EntityFrameworkCore;

namespace cinema.log.server.Models;

public class CinemaLogContext(DbContextOptions<CinemaLogContext> options) : DbContext(options)
{
    public DbSet<User> Users { get; set; }
    public DbSet<Film> Films { get; set; }
}