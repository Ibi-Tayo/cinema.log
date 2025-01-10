using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models;
using cinema.log.server.Models.Entities;
using cinema.log.server.Repositories;
using cinema.log.server.Services;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddOpenApi();
builder.Services.AddSwaggerGen();
builder.Services.AddDbContext<CinemaLogContext>(opt 
    => opt.UseSqlServer(builder.Configuration["LocalDbConnectionString"]));

// Repositories for dependency injection
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IFilmRepository, FilmRepository>();
builder.Services.AddScoped<IReviewRepository, ReviewRepository>();
builder.Services.AddScoped<IUserFilmRatingRepository, UserFilmRatingRepository>();

// Services for dependency injection
builder.Services.AddTransient<IUserService, UserService>();
builder.Services.AddTransient<IFilmService, FilmService>();
builder.Services.AddTransient<IReviewService, ReviewService>();
builder.Services.AddTransient<IUserFilmRatingService, UserFilmRatingService>();
builder.Services.AddTransient<ICalculationService, EloCalculationService>();
builder.Services.AddTransient<ISoundtrackService, SoundtrackService>();

builder.Services.AddRouting(options => options.LowercaseUrls = true);

var app = builder.Build();

if (app.Environment.IsDevelopment()) { app.MapOpenApi(); }

app.UseSwagger();
app.UseSwaggerUI(options =>
{
    options.SwaggerEndpoint("/swagger/v1/swagger.json", "v1");
    options.RoutePrefix = string.Empty;
});

using (var serviceScope = app.Services.CreateScope())
{
    var context = serviceScope.ServiceProvider.GetRequiredService<CinemaLogContext>();
    var dbName = context.Database.GetDbConnection().Database;
    if (!context.Database.CanConnect())
    {
        Console.Error.WriteLine($"Cannot connect to database: {dbName}");
    }
    else
    {
        Console.WriteLine($"Database connection established: {dbName}");
    }
}

app.UseCors(b => b
    .AllowAnyOrigin()
    .AllowAnyMethod()
    .AllowAnyHeader());   

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();
app.Run();