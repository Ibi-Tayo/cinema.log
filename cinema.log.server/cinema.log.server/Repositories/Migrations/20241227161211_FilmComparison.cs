using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace cinema.log.server.Repositories.Migrations
{
    /// <inheritdoc />
    public partial class FilmComparison : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "ComparisonHistories",
                columns: table => new
                {
                    ComparisonHistoryId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    UserId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmAFilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmBFilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    WinningFilmFilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    ComparisonDate = table.Column<DateTime>(type: "datetime2", nullable: false),
                    WasEqual = table.Column<bool>(type: "bit", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ComparisonHistories", x => x.ComparisonHistoryId);
                    table.ForeignKey(
                        name: "FK_ComparisonHistories_Films_FilmAFilmId",
                        column: x => x.FilmAFilmId,
                        principalTable: "Films",
                        principalColumn: "FilmId");
                    table.ForeignKey(
                        name: "FK_ComparisonHistories_Films_FilmBFilmId",
                        column: x => x.FilmBFilmId,
                        principalTable: "Films",
                        principalColumn: "FilmId");
                    table.ForeignKey(
                        name: "FK_ComparisonHistories_Films_WinningFilmFilmId",
                        column: x => x.WinningFilmFilmId,
                        principalTable: "Films",
                        principalColumn: "FilmId");
                    table.ForeignKey(
                        name: "FK_ComparisonHistories_Users_UserId",
                        column: x => x.UserId,
                        principalTable: "Users",
                        principalColumn: "UserId",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "UserFilmRatings",
                columns: table => new
                {
                    UserFilmRatingId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    UserId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    FilmId = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    EloRating = table.Column<float>(type: "real", nullable: false),
                    NumberOfComparisons = table.Column<int>(type: "int", nullable: false),
                    LastUpdated = table.Column<DateTime>(type: "datetime2", nullable: false),
                    InitialRating = table.Column<float>(type: "real", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_UserFilmRatings", x => x.UserFilmRatingId);
                    table.ForeignKey(
                        name: "FK_UserFilmRatings_Films_FilmId",
                        column: x => x.FilmId,
                        principalTable: "Films",
                        principalColumn: "FilmId",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_UserFilmRatings_Users_UserId",
                        column: x => x.UserId,
                        principalTable: "Users",
                        principalColumn: "UserId",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_ComparisonHistories_FilmAFilmId",
                table: "ComparisonHistories",
                column: "FilmAFilmId");

            migrationBuilder.CreateIndex(
                name: "IX_ComparisonHistories_FilmBFilmId",
                table: "ComparisonHistories",
                column: "FilmBFilmId");

            migrationBuilder.CreateIndex(
                name: "IX_ComparisonHistories_UserId",
                table: "ComparisonHistories",
                column: "UserId");

            migrationBuilder.CreateIndex(
                name: "IX_ComparisonHistories_WinningFilmFilmId",
                table: "ComparisonHistories",
                column: "WinningFilmFilmId");

            migrationBuilder.CreateIndex(
                name: "IX_UserFilmRatings_FilmId",
                table: "UserFilmRatings",
                column: "FilmId");

            migrationBuilder.CreateIndex(
                name: "IX_UserFilmRatings_UserId",
                table: "UserFilmRatings",
                column: "UserId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "ComparisonHistories");

            migrationBuilder.DropTable(
                name: "UserFilmRatings");
        }
    }
}
