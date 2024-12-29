namespace cinema.log.server.Abstractions.Interfaces;

public interface ICalculationService
{
    Task<float> CalculateExpectedResult(float filmARating, float filmBRating);
    Task<float> RecalculateFilmRating(float actualResult, float currentRating);
}