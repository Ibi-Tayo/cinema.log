namespace cinema.log.server.Abstractions.Interfaces;

public interface ICalculationService
{
    double CalculateExpectedResult(double filmARating, double filmBRating);
    double RecalculateFilmRating(double expectedResult, double actualResult, 
        double currentRating, double filmKConstantValue);
}