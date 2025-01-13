using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Services;
using cinema.log.server.Utilities;
using FluentAssertions;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class GetUser
{
    private Mock<IUserRepository> _userRepository;
    private Mock<IFilmRepository> _filmRepository;
    private User _testUser;

    [SetUp]
    public void SetUp()
    {
        _userRepository = new Mock<IUserRepository>();
        _filmRepository = new Mock<IFilmRepository>();
        _testUser = new User()
        {
            UserId = Guid.NewGuid(),
            Name = "TestUser",
            Username = "TestUser",
        };
    }
    
    [Test]
    public async Task GetUser_WhenUserExists_ShouldReturnSuccessWithUser()
    {
        // Arrange
        var userDto = Mapper<User, UserDto>.Map(_testUser);
        _userRepository.Setup(repo => repo.GetUserById(_testUser.UserId)).ReturnsAsync(_testUser);
        var sut = new UserService(_userRepository.Object, _filmRepository.Object);
        
        // Act
        var response = await sut.GetUser(_testUser.UserId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(200);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().BeEquivalentTo(userDto);
        _userRepository.Verify(repo => repo.GetUserById(_testUser.UserId), Times.Once);
    }
    
    [Test]
    public async Task GetUser_WhenUserDoesNotExist_ShouldReturnFail()
    {
        // Arrange
        _userRepository.Setup(repo => repo.GetUserById(_testUser.UserId)).ReturnsAsync(null as User);
        var sut = new UserService(_userRepository.Object, _filmRepository.Object);
        
        // Act
        var response = await sut.GetUser(_testUser.UserId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User not found");
        response.Data.Should().BeNull();
        _userRepository.Verify(repo => repo.GetUserById(_testUser.UserId), Times.Once);
    }
    
    [Test]
    public async Task GetUser_WhenUserDoesExistButIdIsIncorrect_ShouldReturnFail()
    {
        // Arrange
        var incorrectId = Guid.NewGuid();
        _userRepository.Setup(repo => repo.GetUserById(_testUser.UserId)).ReturnsAsync(_testUser);
        var sut = new UserService(_userRepository.Object, _filmRepository.Object);
        
        // Act
        var response = await sut.GetUser(incorrectId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User not found");
        response.Data.Should().BeNull();
        _userRepository.Verify(repo => repo.GetUserById(_testUser.UserId), Times.Never);
    }
    
    [Test]
    public async Task GetUser_WhenRepositoryThrowsException_ShouldReturnInternalServerError()
    {
        // Arrange
        var sut = new UserService(_userRepository.Object, _filmRepository.Object);

        _userRepository
            .Setup(repo => repo.GetUserById(_testUser.UserId))
            .ThrowsAsync(new Exception("Database connection failed"));

        // Act
        var response = await sut.GetUser(_testUser.UserId);

        // Assert
        response.Should().NotBeNull();
        response.Should().BeEquivalentTo(new
        {
            StatusCode = 500,
            StatusMessage = "An error occurred while processing your request",
            Data = null as UserDto
        });

        _userRepository.Verify(repo => repo.GetUserById(_testUser.UserId), Times.Once);
    }
}