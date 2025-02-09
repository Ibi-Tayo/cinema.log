using System.Net;
using System.Text;
using System.Text.Json;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using DotNet.Testcontainers.Builders;
using DotNet.Testcontainers.Containers;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

namespace cinema.log.test.IntegrationTests;

[TestFixture]
public class UserControllerTests : IntegrationTest
{
    [Test]
    public async Task DeleteUserWithInvalidIdReturns404()
    {
        // Arrange
        var id = Guid.NewGuid();
        
        // Act
        var response = await Client.DeleteAsync($"/user/deleteuser?userId={id}");
        var responseContent = await response.Content.ReadAsStringAsync();
        var json = JsonSerializer.Deserialize<TestResponse<UserDto>>(responseContent);
        
        // Assert
        json.Should().NotBeNull();
        response.StatusCode.Should().Be(HttpStatusCode.NotFound);
        json.Data.Should().BeNull();
        json.StatusMessage.Should().Be("User not found");
    }
    
    [Test]
    public async Task AddUserWithValidDetailsReturns200()
    {
        // Arrange
        var user = new UserDto() { Username = "test", Name = "Tester" };
        var json = JsonSerializer.Serialize(user);
        var content = new StringContent(json, Encoding.UTF8, "application/json");
       
        // Act
        var response = await Client.PostAsync($"/user/adduser", content);
        var responseContent = await response.Content.ReadAsStringAsync();
        var jsonResponse = JsonSerializer.Deserialize<TestResponse<UserDto>>(responseContent);
        
        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.Created);
        jsonResponse.Should().NotBeNull();
        jsonResponse.StatusMessage.Should().Be("Success");
        
        jsonResponse.Data.Should().NotBeNull();
        jsonResponse.Data.Name.Should().Be(user.Name);
        jsonResponse.Data.Username.Should().Be(user.Username);
    }
}