using cinema.log.server.Services;
using FluentAssertions;

namespace cinema.log.test.EloCalculationServiceTests;

[TestFixture]
public class EloCalculationTests
{
    [TestCase(2400, 2000, 0.91)]
    [TestCase(2000, 2400, 0.09)]
    public void CalculateExpectedResult_GivesCorrectResult(
        double expectedResultFilmRating, 
        double filmBRating,
        double actualExpectedResult)
    {
        // Arrange
        var sut = new EloCalculationService();
        
        // Act
        var result = sut.CalculateExpectedResult(expectedResultFilmRating, filmBRating);
        
        // Assert
        result.Should().Be(actualExpectedResult);
    }

    [TestCase(0.91, 0, 2400, 32, 2371)]
    [TestCase(0.91, 1, 2400, 32, 2403)]
    [TestCase(0.09, 0, 2000, 32, 1997)]
    [TestCase(0.09, 1, 2000, 32, 2029)]
    [TestCase(0.02, 0, 101, 32, 100)]
    public void RecalculateFilmRating_GivesCorrectRatingUpdate(
        double expectedResult, double actualResult,
        double currentRating, double filmKConstantValue,
        double actualExpectedUpdatedFilmRating)
    {
        // Arrange
        var sut = new EloCalculationService();
        
        // Act
        var result = sut.RecalculateFilmRating(expectedResult, actualResult, currentRating, filmKConstantValue);
        
        // Assert
        result.Should().Be(actualExpectedUpdatedFilmRating);
    }
}