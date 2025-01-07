using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class SoundtrackService : ISoundtrackService
{
    // use film Id to get data from film table,
    // send request to spotify to get details
    public Response<FilmSoundtrackDto> GetSoundtrackByFilmId(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // use user id to get all liked tracks from liked track table
    public Response<List<LikedTrackDto>> GetLikedTracksByUserId(Guid userId)
    {
        throw new NotImplementedException();
    }

    // use dto to add new liked track to liked track table
    // (user needs to make sure they have the film soundtrack rating id)
    // (they'd get that by calling GetSoundtrackByFilmId)
    public Response<LikedTrackDto> SetLikedTrack(LikedTrackDto likedTrack)
    {
        throw new NotImplementedException();
    }

    // Self-explanatory
    public Response<bool> DeleteLikedTrack(LikedTrackDto likedTrack)
    {
        throw new NotImplementedException();
    }
}