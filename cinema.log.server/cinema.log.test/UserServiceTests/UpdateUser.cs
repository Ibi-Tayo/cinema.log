using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using FluentAssertions;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class UpdateUser : UserServiceTest
{

    [Test]
    public async Task UpdateUser_WhenUserExists_ReturnUpdatedUser()
    {
        // Arrange
        var updatedUser = new UserDto()
        {
            UserId = Guid.NewGuid(),
            Username = "TestUserNewName",
            Name = "New Name"
        };
        var updatedUserEntity = Mapper<UserDto, User>.Map(updatedUser);
        UserRepository.Setup(repo => repo.UpdateUser(It.IsAny<User>())).ReturnsAsync(updatedUserEntity);
        
        // Act
        var response = await Sut.UpdateUser(updatedUser);
        
        // Assert 
        response.StatusCode.Should().Be(200);
        response.Data.Should().BeEquivalentTo(updatedUser);
        response.StatusMessage.Should().Be("Success");
    }
    
}