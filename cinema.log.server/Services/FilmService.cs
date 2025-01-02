using cinema.log.server.Abstractions.Interfaces;
using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;

namespace cinema.log.server.Services;

public class FilmService : IFilmService
{
    public async Task<Response<FilmDto>> GetFilmFromDb(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // this will have to get called server side - probably in the review service when a user leaves a review
    public async Task<bool> AddFilmToDb(int externalId)
    {
        // call GetFilmDetailsByExternalId() 
        // use the film we get back to add to db using repository
        throw new NotImplementedException();
    }

    public async Task<Response<FilmDto>> UpdateFilmInDb(FilmDto film)
    {
        throw new NotImplementedException();
    }

    public async Task<Response<FilmDto>> DeleteFilmInDb(Guid filmId)
    {
        throw new NotImplementedException();
    }

    // To remember:
    // 'External' basically means im sending an api call to tmdb
    // External id is the movie_id param that is an integer, obviously not to be confused with my own GUID
    
    // This method will be called from outside so need to return DTO
    public async Task<Response<List<FilmDto>>> SearchFilmFromExternal(string searchTerm)
    {
        // will need to get the search results and make a list out of it
        throw new NotImplementedException();
    }

    public async Task<Response<List<FilmImageDto>>> GetFilmImagesFromExternal(int externalId)
    {
        // call the tmdb api and hit images endpoint, get all the backdrops and all the posters
        // map these to the dtos and make a list
        
        throw new NotImplementedException();
    }
    
    // flow: when user has found a film, they should have the film details, when they leave a review, instead of sending all of this back
    // they just need to ping the server with the tmdb movie_id. this method will call tmdb and then map to our entity
    // then the entity will end up getting put in the db - see 'AddFilmToDb' above
    private async Task<Film> GetFilmDetailsByExternalId(int externalId)
    {
        // this method will get called in AddFilmToDb
        // send api call to tmdb
        // manually map it to my film entity and return
        throw new NotImplementedException();
    }
}