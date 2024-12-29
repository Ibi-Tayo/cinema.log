using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace cinema.log.server.Repositories.Migrations
{
    /// <inheritdoc />
    public partial class FilmDetailsUpdate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<int>(
                name: "ReleaseYear",
                table: "Films",
                type: "int",
                nullable: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "ReleaseYear",
                table: "Films");
        }
    }
}
