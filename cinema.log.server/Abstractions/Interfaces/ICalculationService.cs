namespace cinema.log.server.Abstractions.Interfaces;

public interface ICalculationService
{
    float CalculateExpectedResult(float? filmARating, float? filmBRating);
    float RecalculateFilmRating(float expectedResult, float actualResult, float? currentRating);
}