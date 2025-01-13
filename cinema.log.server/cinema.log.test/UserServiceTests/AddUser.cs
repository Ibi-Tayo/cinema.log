using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using FluentAssertions;
using Microsoft.AspNetCore.Http;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class AddUser : UserServiceTest
{
    [Test]
    public async Task AddUser_WhenUserIsValid_ShouldCreateUserAndReturn201()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "TestUser",
            Name = "Test User"
        };
        var createdUser = new User 
        { 
            UserId = Guid.NewGuid(),
            Username = userDto.Username,
            Name = userDto.Name
        };

        UserRepository
            .Setup(repo => repo.CreateUser(It.IsAny<User>()))
            .ReturnsAsync(createdUser);

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.Should().BeEquivalentTo(new
        {
            StatusCode = StatusCodes.Status201Created,
            StatusMessage = "Success",
            Data = Mapper<User, UserDto>.Map(createdUser)
        });

        UserRepository.Verify(repo => repo.CreateUser(It.Is<User>(u => 
            u.Username == userDto.Username && 
            u.Name == userDto.Name)), Times.Once);
    }

    [Test]
    public async Task AddUser_WhenUsernameIsTooShort_ShouldReturn400()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "ab",  // Less than 3 characters
            Name = "Test User"
        };

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(400);
        response.StatusMessage.Should().Contain("Username must be between 3 and 20 characters");
        
        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Never);
    }

    [Test]
    public async Task AddUser_WhenUsernameIsTooLong_ShouldReturn400()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = new string('a', 21),  // 21 characters
            Name = "Test User"
        };

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(400);
        response.StatusMessage.Should().Contain("Username must be between 3 and 20 characters");
        
        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Never);
    }

    [Test]
    public async Task AddUser_WhenNameContainsDigits_ShouldReturn400()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "TestUser",
            Name = "Test User123"
        };

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(400);
        response.StatusMessage.Should().Contain("Name cannot contain digits");
        
        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Never);
    }

    [Test]
    public async Task AddUser_WhenNameIsMissing_ShouldReturn400()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "TestUser",
            Name = ""
        };

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(400);
        response.StatusMessage.Should().Contain("Name is required");
        
        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Never);
    }

    [Test]
    public async Task AddUser_WhenRepositoryCreateFails_ShouldReturn500()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "TestUser",
            Name = "Test User"
        };

        UserRepository
            .Setup(repo => repo.CreateUser(It.IsAny<User>()))
            .ReturnsAsync(null as User);

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.Should().BeEquivalentTo(new
        {
            StatusCode = StatusCodes.Status500InternalServerError,
            StatusMessage = "Internal Server Error",
            Data = null as UserDto
        });

        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Once);
    }

    [Test]
    public async Task AddUser_WhenRepositoryThrowsException_ShouldReturn500()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "TestUser",
            Name = "Test User"
        };

        UserRepository
            .Setup(repo => repo.CreateUser(It.IsAny<User>()))
            .ThrowsAsync(new Exception("Database connection failed"));

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.Should().BeEquivalentTo(new
        {
            StatusCode = StatusCodes.Status500InternalServerError,
            StatusMessage = "Internal Server Error",
            Data = null as UserDto
        });

        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Once);
    }

    [Test]
    public async Task AddUser_WhenMultipleValidationErrors_ShouldReturnAllErrors()
    {
        // Arrange
        var userDto = new UserDto 
        { 
            Username = "a",  // Too short
            Name = "Test123"  // Contains digits
        };

        // Act
        var response = await Sut.AddUser(userDto);

        // Assert
        response.Should().NotBeNull();
        response.StatusCode.Should().Be(400);
        response.StatusMessage.Should().Contain("Username must be between 3 and 20 characters");
        response.StatusMessage.Should().Contain("Name cannot contain digits");
        
        UserRepository.Verify(repo => repo.CreateUser(It.IsAny<User>()), Times.Never);
    }
}