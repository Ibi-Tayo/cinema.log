using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Abstractions.Interfaces;

public interface ISoundtrackService
{
    Task<Response<FilmSoundtrackDto>> GetSoundtrackByFilmId(Guid filmId);
    Response<List<LikedTrackDto>> GetLikedTracksByUserId(Guid userId);
    Response<LikedTrackDto> SetLikedTrack(LikedTrackDto likedTrack);
    Response<bool> DeleteLikedTrack(LikedTrackDto likedTrack);
}