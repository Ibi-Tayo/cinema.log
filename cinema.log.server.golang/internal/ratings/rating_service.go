package ratings

import (
	"context"
	"math"
	"time"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	RatingStore RatingStore
}

type RatingStore interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRating, error)
	CreateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	UpdateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error)
	CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
}

func NewService(r RatingStore) *Service {
	return &Service{
		RatingStore: r,
	}
}

func (s Service) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	return s.RatingStore.GetRating(ctx, userId, filmId)
}

func (s Service) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	return s.RatingStore.GetAllRatings(ctx)
}

func (s Service) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
	return s.RatingStore.CreateComparison(ctx, comparison)
}

func (s Service) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	return s.RatingStore.HasBeenCompared(ctx, userId, filmAId, filmBId)
}

func (s Service) GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error) {
	return s.RatingStore.GetComparisonHistory(ctx, userId)
}

func (s Service) CreateRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID, initialRating float32) (*domain.UserFilmRating, error) {
	// Create a new rating with initial values
	rating := domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              userId,
		FilmId:              filmId,
		EloRating:           float64(s.getInitialEloRating(initialRating)),
		NumberOfComparisons: 0,
		LastUpdated:         time.Now(),
		InitialRating:       initialRating,
		KConstantValue:      40, // Start with highest K value for new ratings
	}

	return s.RatingStore.CreateRating(ctx, rating)
}

func (s Service) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, comparison domain.ComparisonHistory) (*domain.ComparisonPair, error) {
	filmA := ratings.FilmA
	filmB := ratings.FilmB
	// Set the results from the film head to head
	filmAResult, filmBResult := s.defineFilmContestResult(filmA.FilmId, filmB.FilmId, comparison)
	// Calculate expected results for film A and film B
	filmAExpectedResult := s.calculateExpectedResult(filmA.EloRating, filmB.EloRating)
	filmBExpectedResult := s.calculateExpectedResult(filmB.EloRating, filmA.EloRating)

	// Update K Constants
	filmA.KConstantValue = s.updateKConstantValue(filmA)
	filmB.KConstantValue = s.updateKConstantValue(filmB)

	// Recalculate film rating for film A and film B
	filmANewRating := s.recalculateFilmRating(filmAExpectedResult, filmAResult, filmA.EloRating, filmA.KConstantValue)
	filmBNewRating := s.recalculateFilmRating(filmBExpectedResult, filmBResult, filmB.EloRating, filmB.KConstantValue)

	// Update films
	filmA.EloRating = filmANewRating
	filmA.LastUpdated = time.Now()
	filmA.NumberOfComparisons += 1

	filmB.EloRating = filmBNewRating
	filmB.LastUpdated = time.Now()
	filmB.NumberOfComparisons += 1

	// Update both films in the store
	updatedFilmA, err := s.RatingStore.UpdateRating(ctx, filmA)
	if err != nil {
		return nil, err
	}

	updatedFilmB, err := s.RatingStore.UpdateRating(ctx, filmB)
	if err != nil {
		return nil, err
	}

	// Return the updated comparison pair
	return &domain.ComparisonPair{
		FilmA: *updatedFilmA,
		FilmB: *updatedFilmB,
	}, nil
}

/*
   Calculate expected result
   ---
   Ea = 1 / (1 + 10^(Rb - Ra)/400)
   Where:
   Ea is expected score of film a
   Ra is current rating of film a
   Rb is current rating of film b
*/

func (s Service) calculateExpectedResult(filmUnderReviewEloRating float64, challengerFilmEloRating float64) float64 {
	rawCalc := (1 / (1 + math.Pow(10, (challengerFilmEloRating-filmUnderReviewEloRating)/400)))

	ratio := math.Pow(10, float64(2)) // round to 2 dp
	return math.Round(rawCalc*ratio) / ratio
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

func (s Service) recalculateFilmRating(expectedResult float64, actualResult float64,
	currentRating float64, filmKConstantValue float64) float64 {

	rawCalc := currentRating + filmKConstantValue*(actualResult-expectedResult)
	if rawCalc <= 100 {
		return 100 // 100 is the lowest rating value you can get, so no further decreases past this point
	}

	return math.Round(rawCalc)
}

func (s Service) defineFilmContestResult(filmA uuid.UUID, filmB uuid.UUID, comparison domain.ComparisonHistory) (float64, float64) {
	var filmAResult, filmBResult float64
	if comparison.WasEqual {
		filmAResult = 0.5
		filmBResult = 0.5
		return filmAResult, filmBResult
	}

	if filmA == comparison.WinningFilmId {
		filmAResult = 1
		filmBResult = 0
	} else if filmB == comparison.WinningFilmId {
		filmAResult = 0
		filmBResult = 1
	}
	return filmAResult, filmBResult
}

func (s Service) updateKConstantValue(film domain.UserFilmRating) float64 {
	numberOfComparisons := film.NumberOfComparisons
	switch {
	case numberOfComparisons >= 0 && numberOfComparisons < 10:
		return 40
	case numberOfComparisons >= 10 && numberOfComparisons < 20:
		return 20
	case numberOfComparisons >= 20:
		return 10
	default:
		return 40
	}
}

// this is generated from the users initial 5 star rating 'InitialRating'
func (s Service) getInitialEloRating(rating float32) float32 {
	switch {
	case rating >= 0 && rating < 2:
		return 950
	case rating >= 2 && rating < 3:
		return 1000
	case rating >= 3 && rating < 4:
		return 1050
	case rating >= 4 && rating <= 5:
		return 1100
	default:
		return 1000
	}
}
