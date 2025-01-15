using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using FluentAssertions;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class GetUser : UserServiceTest
{
    [Test]
    public async Task GetUser_WhenUserExists_ShouldReturnSuccessWithUser()
    {
        // Arrange
        var userDto = Mapper<User, UserDto>.Map(TestUser);
        UserRepository.Setup(repo => repo.GetUserById(TestUser.UserId)).ReturnsAsync(TestUser);
        
        // Act
        var response = await Sut.GetUser(TestUser.UserId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(200);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().BeEquivalentTo(userDto);
        UserRepository.Verify(repo => repo.GetUserById(TestUser.UserId), Times.Once);
    }
    
    [Test]
    public async Task GetUser_WhenUserDoesNotExist_ShouldReturnFail()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserById(TestUser.UserId)).ReturnsAsync(null as User);
        
        // Act
        var response = await Sut.GetUser(TestUser.UserId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User not found");
        response.Data.Should().BeNull();
        UserRepository.Verify(repo => repo.GetUserById(TestUser.UserId), Times.Once);
    }
    
    [Test]
    public async Task GetUser_WhenUserDoesExistButIdIsIncorrect_ShouldReturnFail()
    {
        // Arrange
        var incorrectId = Guid.NewGuid();
        UserRepository.Setup(repo => repo.GetUserById(TestUser.UserId)).ReturnsAsync(TestUser);
        
        // Act
        var response = await Sut.GetUser(incorrectId);
        
        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User not found");
        response.Data.Should().BeNull();
        UserRepository.Verify(repo => repo.GetUserById(TestUser.UserId), Times.Never);
    }
    
    [Test]
    public async Task GetUser_WhenRepositoryThrowsException_ShouldReturnInternalServerError()
    {
        // Arrange
        UserRepository
            .Setup(repo => repo.GetUserById(TestUser.UserId))
            .ThrowsAsync(new Exception("Database connection failed"));

        // Act
        var response = await Sut.GetUser(TestUser.UserId);

        // Assert
        response.Should().NotBeNull();
        response.Should().BeEquivalentTo(new
        {
            StatusCode = 500,
            StatusMessage = "An error occurred while processing your request",
            Data = null as UserDto
        });

        UserRepository.Verify(repo => repo.GetUserById(TestUser.UserId), Times.Once);
    }
}