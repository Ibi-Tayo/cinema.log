using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace cinema.log.server.Repositories.Migrations
{
    /// <inheritdoc />
    public partial class CachedSoundtrack : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "CachedSoundtracks",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmTitle = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    SpotifyAlbumId = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    SoundtrackName = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    Artists = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    AlbumArtUrl = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    LastUpdated = table.Column<DateTime>(type: "datetime2", nullable: false),
                    TracksJson = table.Column<string>(type: "nvarchar(max)", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_CachedSoundtracks", x => x.Id);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "CachedSoundtracks");
        }
    }
}
