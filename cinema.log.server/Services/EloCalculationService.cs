using cinema.log.server.Abstractions.Interfaces;

namespace cinema.log.server.Services;

public class EloCalculationService: ICalculationService
{
    /*
    Calculate expected result
    ---
    Ea = 1 / (1 + 10^(Rb - Ra)/400)
    Where:
    Ea is expected score of film a
    Ra is current rating of film a
    Rb is current rating of film b
    */
    
    
    
    /*
    Recalculate elo rating
    ---
    R'a = Ra + K(Sa - Ea)
    Where: 
    R'a is new rating for film a
    Ra is current rating for film a
    K is K-factor (to be adjusted based on review date) (With the most recent having the highest K)
    Sa is actual result of match up (0 for loss, 0.5 for draw, 1 for win)
    Ea is expected result (Ea = 1 / (1 + 10^(Rb - Ra)/400)) 
    */

    public float CalculateExpectedResult(float filmARating, float filmBRating)
    {
        throw new NotImplementedException();
    }

    public float RecalculateFilmRating(float actualResult, float currentRating)
    {
        throw new NotImplementedException();
    }
}