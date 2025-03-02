using System.Security.Claims;
using System.Text.Encodings.Web;
using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.Entities;
using DotNet.Testcontainers.Builders;
using DotNet.Testcontainers.Containers;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace cinema.log.test.IntegrationTests;

public abstract class IntegrationTest
{
    private const string Database = "master";
    private const string Username = "sa";
    private const string Password = "$trongPassword";
    private const ushort MsSqlPort = 1433;
    internal WebApplicationFactory<IApiMarker> Factory;
    internal HttpClient Client;
    private IContainer _container;

    [OneTimeSetUp]
    public async Task OneTimeSetUp()
    {
        _container = new ContainerBuilder()
            .WithImage("mcr.microsoft.com/mssql/server:2022-latest")
            .WithPortBinding(MsSqlPort, true)
            .WithEnvironment("ACCEPT_EULA", "Y")
            .WithEnvironment("SQLCMDUSER", Username)
            .WithEnvironment("SQLCMDPASSWORD", Password)
            .WithEnvironment("MSSQL_SA_PASSWORD", Password)
            .WithWaitStrategy(Wait.ForUnixContainer().UntilPortIsAvailable(MsSqlPort))
            .Build();

        await _container.StartAsync();

        var host = _container.Hostname;
        var port = _container.GetMappedPublicPort(MsSqlPort);

        // Replace connection string in DbContext
        var connectionString =
            $"Server={host},{port};Database={Database};User Id={Username};Password={Password};TrustServerCertificate=True";
        Factory = new WebApplicationFactory<IApiMarker>()
            .WithWebHostBuilder(builder =>
            {
                builder.ConfigureServices(services =>
                {
                    services.AddDbContext<CinemaLogContext>(options =>
                        options.UseSqlServer(connectionString));
                    // Remove the existing authentication and authorization services
                    var authenticationServiceDescriptor = services.SingleOrDefault(
                        d => d.ServiceType == typeof(IAuthenticationService));
                    if (authenticationServiceDescriptor != null)
                        services.Remove(authenticationServiceDescriptor);

                    var authenticationHandlerDescriptor = services.SingleOrDefault(
                        d => d.ServiceType == typeof(AuthenticationHandlerProvider));
                    if (authenticationHandlerDescriptor != null)
                        services.Remove(authenticationHandlerDescriptor);

                    // Replace authorization services with a mocked version that allows everything
                    services.AddSingleton<IAuthorizationService, AllowAllAuthorizationService>();
                    services.AddSingleton<IAuthorizationPolicyProvider, AllowAllAuthorizationPolicyProvider>();
                    services.AddSingleton<IAuthorizationHandler, AllowAllAuthorizationHandler>();

                    // Replace authentication with a simple no-op version
                    services.AddAuthentication("Test")
                        .AddScheme<AuthenticationSchemeOptions, TestAuthHandler>("Test", options => { });
                });
            });

        Client = Factory.CreateClient();
        // Initialize database
        var scope = Factory.Services.CreateScope();
        var dbContext = scope.ServiceProvider.GetRequiredService<CinemaLogContext>();
        await dbContext.Database.MigrateAsync();
    }

    [OneTimeTearDown]
    public async Task OneTimeTearDown()
    {
        await _container.StopAsync();
        await _container.DisposeAsync();
        Client.Dispose();
        await Factory.DisposeAsync();
    }
    
    // Simple authorization service that allows all requests
public class AllowAllAuthorizationService : IAuthorizationService
{
    public Task<AuthorizationResult> AuthorizeAsync(ClaimsPrincipal user, object? resource, IEnumerable<IAuthorizationRequirement> requirements)
    {
        return Task.FromResult(AuthorizationResult.Success());
    }

    public Task<AuthorizationResult> AuthorizeAsync(ClaimsPrincipal user, object? resource, string policyName)
    {
        return Task.FromResult(AuthorizationResult.Success());
    }
}

// Simple policy provider that allows all policies
public class AllowAllAuthorizationPolicyProvider : IAuthorizationPolicyProvider
{
    public Task<AuthorizationPolicy> GetDefaultPolicyAsync() => 
        Task.FromResult(new AuthorizationPolicyBuilder().RequireAssertion(_ => true).Build());

    public Task<AuthorizationPolicy?> GetFallbackPolicyAsync() =>
        Task.FromResult<AuthorizationPolicy?>(null);

    public Task<AuthorizationPolicy?> GetPolicyAsync(string policyName) =>
        Task.FromResult<AuthorizationPolicy?>(new AuthorizationPolicyBuilder().RequireAssertion(_ => true).Build());
}

// Simple authorization handler that succeeds for all requirements
public class AllowAllAuthorizationHandler : IAuthorizationHandler
{
    public Task HandleAsync(AuthorizationHandlerContext context)
    {
        foreach (var requirement in context.PendingRequirements.ToList())
        {
            context.Succeed(requirement);
        }
        return Task.CompletedTask;
    }
}

// Test authentication handler that creates a minimal identity
public class TestAuthHandler : AuthenticationHandler<AuthenticationSchemeOptions>
{
    public TestAuthHandler(IOptionsMonitor<AuthenticationSchemeOptions> options, ILoggerFactory logger, UrlEncoder encoder, ISystemClock clock) 
        : base(options, logger, encoder, clock)
    {
    }

    protected override Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        var claims = new[] { new Claim(ClaimTypes.Name, "Test User") };
        var identity = new ClaimsIdentity(claims, "Test");
        var principal = new ClaimsPrincipal(identity);
        var ticket = new AuthenticationTicket(principal, "Test");
        return Task.FromResult(AuthenticateResult.Success(ticket));
    }
}
}