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
    
    public double CalculateExpectedResult(double filmARating, double filmBRating)
    {
        return (1 / (1 + Math.Pow(10, (filmBRating - filmARating) / 400)));
    }
    
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
    
    public double RecalculateFilmRating(double expectedResult, double actualResult,
        double currentRating, double filmKConstantValue)
    {
        return currentRating + filmKConstantValue * (actualResult - expectedResult);
    }
}