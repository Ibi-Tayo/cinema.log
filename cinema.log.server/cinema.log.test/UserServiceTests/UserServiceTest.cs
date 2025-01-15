using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.Entities;
using cinema.log.server.Services;
using Moq;

namespace cinema.log.test.UserServiceTests;

public class UserServiceTest
{
    internal Mock<IUserRepository> UserRepository;
    internal Mock<IFilmRepository> FilmRepository;
    internal User TestUser;
    internal UserService Sut;
    
    [SetUp]
    public void SetUp()
    {
        UserRepository = new Mock<IUserRepository>();
        FilmRepository = new Mock<IFilmRepository>();
        Sut = new UserService(UserRepository.Object, FilmRepository.Object);
        TestUser = new User()
        {
            UserId = Guid.NewGuid(),
            Name = "TestUser",
            Username = "TestUser",
        };
    }

}