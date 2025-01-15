using cinema.log.server.Models.Entities;
using FluentAssertions;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class DeleteUserTests : UserServiceTest
{
    [Test]
    public async Task DeleteUser_WhenUserExists_ReturnsNoContent()
    {
        // Arrange
        var userId = Guid.NewGuid();
        UserRepository.Setup(repo => repo.DeleteUserById(userId)).ReturnsAsync(TestUser);

        // Act
        var response = await Sut.DeleteUser(userId);

        // Assert
        response.StatusCode.Should().Be(204);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().BeNull();
    }

    [Test]
    public async Task DeleteUser_WhenUserDoesNotExist_ReturnsNotFound()
    {
        // Arrange
        var userId = Guid.NewGuid();
        UserRepository.Setup(repo => repo.DeleteUserById(userId)).ReturnsAsync((User)null);

        // Act
        var response = await Sut.DeleteUser(userId);

        // Assert
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User not found");
        response.Data.Should().BeNull();
    }
}