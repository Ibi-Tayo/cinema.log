using cinema.log.server.Models.DTOs;
using cinema.log.server.Models.Entities;
using cinema.log.server.Utilities;
using FluentAssertions;
using Moq;

namespace cinema.log.test.UserServiceTests;

[TestFixture]
public class GetUserReviews: UserServiceTest
{
    private List<Review> _testReviews;

    [SetUp]
    public new void SetUp()
    {
        base.SetUp();
        _testReviews = new List<Review>
        {
            new() { ReviewId = Guid.NewGuid(), UserId = TestUser.UserId, FilmId = Guid.NewGuid(), Rating = 5 },
            new() { ReviewId = Guid.NewGuid(), UserId = TestUser.UserId, FilmId = Guid.NewGuid(), Rating = 4 }
        };
    }

    [Test]
    public async Task GetUserReviews_WhenUserHasReviews_ReturnsReviewsList()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserReviews(TestUser.UserId)).ReturnsAsync(_testReviews);

        // Act
        var response = await Sut.GetUserReviews(TestUser.UserId);

        // Assert
        response.StatusCode.Should().Be(200);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().HaveCount(2);
        response.Data.Should().BeEquivalentTo(_testReviews.Select(Mapper<Review, ReviewDto>.Map));
    }

    [Test]
    public async Task GetUserReviews_WhenUserHasNoReviews_ReturnsNotFound()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserReviews(TestUser.UserId)).ReturnsAsync(new List<Review>());

        // Act
        var response = await Sut.GetUserReviews(TestUser.UserId);

        // Assert
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User reviews not found");
        response.Data.Should().BeNull();
    }
}

[TestFixture]
public class GetFilmsReviewedByUser : UserServiceTest
{
    private List<Review> _testReviews;
    private List<Film> _testFilms;

    [SetUp]
    public new void SetUp()
    {
        base.SetUp();
        _testReviews = new List<Review>
        {
            new() { ReviewId = Guid.NewGuid(), UserId = TestUser.UserId, FilmId = Guid.NewGuid(), Rating = 5 },
            new() { ReviewId = Guid.NewGuid(), UserId = TestUser.UserId, FilmId = Guid.NewGuid(), Rating = 4 }
        };

        _testFilms = new List<Film>
        {
            new() { FilmId = _testReviews[0].FilmId, Title = "Test Film 1" },
            new() { FilmId = _testReviews[1].FilmId, Title = "Test Film 2" }
        };
    }

    [Test]
    public async Task GetFilmsReviewedByUser_WhenUserHasReviewedFilms_ReturnsFilmsList()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserReviews(TestUser.UserId)).ReturnsAsync(_testReviews);

        foreach (var film in _testFilms)
        {
            FilmRepository.Setup(repo => repo.GetFilmById(film.FilmId)).ReturnsAsync(film);
        }

        // Act
        var response = await Sut.GetFilmsReviewedByUser(TestUser.UserId);

        // Assert
        response.StatusCode.Should().Be(200);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().HaveCount(2);
        response.Data.Should().BeEquivalentTo(_testFilms.Select(Mapper<Film, FilmDto>.Map));
    }

    [Test]
    public async Task GetFilmsReviewedByUser_WhenUserHasNoReviews_ReturnsNotFound()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserReviews(TestUser.UserId)).ReturnsAsync(new List<Review>());

        // Act
        var response = await Sut.GetFilmsReviewedByUser(TestUser.UserId);

        // Assert
        response.StatusCode.Should().Be(404);
        response.StatusMessage.Should().Be("User list of films reviewed not found");
        response.Data.Should().BeNull();
    }

    [Test]
    public async Task GetFilmsReviewedByUser_WhenSomeFilmsAreNotFound_ReturnsOnlyFoundFilms()
    {
        // Arrange
        UserRepository.Setup(repo => repo.GetUserReviews(TestUser.UserId)).ReturnsAsync(_testReviews);

        // Only set up one film to be found
        FilmRepository.Setup(repo => repo.GetFilmById(_testFilms[0].FilmId)).ReturnsAsync(_testFilms[0]);
        FilmRepository.Setup(repo => repo.GetFilmById(_testFilms[1].FilmId)).ReturnsAsync((Film)null);

        // Act
        var response = await Sut.GetFilmsReviewedByUser(TestUser.UserId);

        // Assert
        response.StatusCode.Should().Be(200);
        response.StatusMessage.Should().Be("Success");
        response.Data.Should().HaveCount(1);
        response.Data.Should().BeEquivalentTo(new List<FilmDto> { Mapper<Film, FilmDto>.Map(_testFilms[0]) });
    }
}