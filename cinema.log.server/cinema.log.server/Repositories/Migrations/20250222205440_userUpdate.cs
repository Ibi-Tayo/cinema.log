using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace cinema.log.server.Repositories.Migrations
{
    /// <inheritdoc />
    public partial class userUpdate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<long>(
                name: "GithubId",
                table: "Users",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "GithubId",
                table: "Users");
        }
    }
}
