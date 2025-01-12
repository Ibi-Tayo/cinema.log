using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace cinema.log.server.Repositories.Migrations
{
    /// <inheritdoc />
    public partial class SoundtrackMigration : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AlterColumn<double>(
                name: "EloRating",
                table: "UserFilmRatings",
                type: "float",
                nullable: false,
                defaultValue: 0.0,
                oldClrType: typeof(float),
                oldType: "real",
                oldNullable: true);

            migrationBuilder.AddColumn<double>(
                name: "KConstantValue",
                table: "UserFilmRatings",
                type: "float",
                nullable: false,
                defaultValue: 0.0);

            migrationBuilder.CreateTable(
                name: "SpotifyApi",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    AccessToken = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    ExpiryDate = table.Column<DateTime>(type: "datetime2", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_SpotifyApi", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "UserFilmSoundtrackRatings",
                columns: table => new
                {
                    UserFilmSoundtrackRatingId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    UserId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    Rating = table.Column<int>(type: "int", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_UserFilmSoundtrackRatings", x => x.UserFilmSoundtrackRatingId);
                });

            migrationBuilder.CreateTable(
                name: "LikedTracks",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    UserId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    TrackTitle = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    UserFilmSoundtrackRatingId = table.Column<Guid>(type: "uniqueidentifier", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_LikedTracks", x => x.Id);
                    table.ForeignKey(
                        name: "FK_LikedTracks_UserFilmSoundtrackRatings_UserFilmSoundtrackRatingId",
                        column: x => x.UserFilmSoundtrackRatingId,
                        principalTable: "UserFilmSoundtrackRatings",
                        principalColumn: "UserFilmSoundtrackRatingId",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_LikedTracks_UserFilmSoundtrackRatingId",
                table: "LikedTracks",
                column: "UserFilmSoundtrackRatingId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "LikedTracks");

            migrationBuilder.DropTable(
                name: "SpotifyApi");

            migrationBuilder.DropTable(
                name: "UserFilmSoundtrackRatings");

            migrationBuilder.DropColumn(
                name: "KConstantValue",
                table: "UserFilmRatings");

            migrationBuilder.AlterColumn<float>(
                name: "EloRating",
                table: "UserFilmRatings",
                type: "real",
                nullable: true,
                oldClrType: typeof(double),
                oldType: "float");
        }
    }
}
