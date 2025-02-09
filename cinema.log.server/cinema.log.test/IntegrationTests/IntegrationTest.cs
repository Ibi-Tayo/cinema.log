using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.Entities;
using DotNet.Testcontainers.Builders;
using DotNet.Testcontainers.Containers;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

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
}