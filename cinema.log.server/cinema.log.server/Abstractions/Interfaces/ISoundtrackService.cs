using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface ISoundtrackService
{
    Task<Response<FilmSoundtrackDto>> GetSoundtrackByFilmId(Guid filmId, Guid userId);
    Task<Response<List<LikedTrackDto>>> GetLikedTracksByUserId(Guid userId);
    Task<Response<List<LikedTrackDto>>> GetLikedTracksFromFilmRatingId(Guid filmRatingId);
    Task<Response<LikedTrackDto?>> SetLikedTrack(LikedTrackDto likedTrack);
    Task<Response<bool>> DeleteLikedTrack(LikedTrackDto likedTrack);
}