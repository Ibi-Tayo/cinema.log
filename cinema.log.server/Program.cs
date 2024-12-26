using cinema.log.server.Models;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddOpenApi();
builder.Services.AddSwaggerGen();
builder.Services.AddDbContext<CinemaLogContext>(opt 
    => opt.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseSwagger();
app.UseSwaggerUI(options =>
{
    options.SwaggerEndpoint("/swagger/v1/swagger.json", "v1");
    options.RoutePrefix = string.Empty;
});

// need to run migrations first before this can connect to database
using (var serviceScope = app.Services.CreateScope())
{
    var context = serviceScope.ServiceProvider.GetRequiredService<CinemaLogContext>();
    var connectionString = context.Database.GetDbConnection().ConnectionString;
    if (!context.Database.CanConnect())
    {
        throw new Exception($"Cannot connect to database: {connectionString}");
    }
    Console.WriteLine($"Database connection established: {connectionString}");
    context.Database.Migrate();
}

app.UseCors(b => b
    .AllowAnyOrigin()
    .AllowAnyMethod()
    .AllowAnyHeader());   

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.Run();