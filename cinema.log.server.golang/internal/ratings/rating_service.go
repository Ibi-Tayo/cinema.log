package ratings

import (
	"context"
	"database/sql"
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
	GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error)
	CreateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	UpdateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error)
	CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
	BulkGetRatings(ctx context.Context, userId uuid.UUID, filmIds []uuid.UUID) (map[uuid.UUID]*domain.UserFilmRating, error)
	BulkHasBeenCompared(ctx context.Context, userId uuid.UUID, pairs []domain.ComparisonPair) (map[string]bool, error)
	BulkInsertComparisons(ctx context.Context, comparisons []domain.ComparisonHistory) error
	BulkUpdateRatings(ctx context.Context, tx *sql.Tx, ratings []domain.UserFilmRating) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
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

func (s Service) GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error) {
	return s.RatingStore.GetRatingsByUserId(ctx, userId)
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
	case numberOfComparisons >= 0 && numberOfComparisons < 20:
		return 40
	case numberOfComparisons >= 20 && numberOfComparisons < 40:
		return 20
	case numberOfComparisons >= 40:
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

// ProcessBatchComparisons processes multiple film comparisons in a single transaction
func (s Service) ProcessBatchComparisons(ctx context.Context, userId, targetFilmId uuid.UUID, comparisons []ComparisonItem) error {
	if len(comparisons) == 0 {
		return nil
	}

	// Cap at 50 comparisons
	if len(comparisons) > 50 {
		comparisons = comparisons[:50]
	}

	// Collect all film IDs (target + all challengers)
	filmIds := []uuid.UUID{targetFilmId}
	for _, comp := range comparisons {
		filmIds = append(filmIds, comp.ChallengerFilmId)
	}

	// Fetch all ratings in bulk
	ratingsMap, err := s.RatingStore.BulkGetRatings(ctx, userId, filmIds)
	if err != nil {
		return err
	}

	// Validate target film exists
	targetRating, ok := ratingsMap[targetFilmId]
	if !ok {
		return ErrRatingNotFound
	}

	// Filter out duplicates and validate all films exist
	var validComparisons []ComparisonItem
	for _, comp := range comparisons {
		challengerRating, ok := ratingsMap[comp.ChallengerFilmId]
		if !ok {
			continue // Skip missing films
		}

		// Check if already compared
		hasBeenCompared, err := s.RatingStore.HasBeenCompared(ctx, userId, targetFilmId, comp.ChallengerFilmId)
		if err != nil {
			return err
		}
		if hasBeenCompared {
			continue // Skip duplicates
		}

		validComparisons = append(validComparisons, comp)
		ratingsMap[comp.ChallengerFilmId] = challengerRating
	}

	if len(validComparisons) == 0 {
		return nil // Nothing to process
	}

	// Begin transaction
	tx, err := s.RatingStore.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Process comparisons sequentially to maintain K-factor progression
	var updatedRatings []domain.UserFilmRating
	var comparisonHistory []domain.ComparisonHistory
	currentTargetRating := *targetRating

	for _, comp := range validComparisons {
		challengerRating := ratingsMap[comp.ChallengerFilmId]

		// Determine winner
		var winningFilmId uuid.UUID
		var wasEqual bool
		var targetResult, challengerResult float64

		switch comp.Result {
		case "better":
			winningFilmId = targetFilmId
			targetResult = 1.0
			challengerResult = 0.0
		case "worse":
			winningFilmId = comp.ChallengerFilmId
			targetResult = 0.0
			challengerResult = 1.0
		case "same":
			winningFilmId = targetFilmId // Use target as placeholder
			wasEqual = true
			targetResult = 0.5
			challengerResult = 0.5
		}

		// Calculate expected results
		targetExpected := s.calculateExpectedResult(currentTargetRating.EloRating, challengerRating.EloRating)
		challengerExpected := s.calculateExpectedResult(challengerRating.EloRating, currentTargetRating.EloRating)

		// Update K constants
		currentTargetRating.KConstantValue = s.updateKConstantValue(currentTargetRating)
		challengerRating.KConstantValue = s.updateKConstantValue(*challengerRating)

		// Recalculate ratings
		targetNewRating := s.recalculateFilmRating(targetExpected, targetResult, currentTargetRating.EloRating, currentTargetRating.KConstantValue)
		challengerNewRating := s.recalculateFilmRating(challengerExpected, challengerResult, challengerRating.EloRating, challengerRating.KConstantValue)

		// Update ratings
		currentTargetRating.EloRating = targetNewRating
		currentTargetRating.LastUpdated = time.Now()
		currentTargetRating.NumberOfComparisons += 1

		challengerRating.EloRating = challengerNewRating
		challengerRating.LastUpdated = time.Now()
		challengerRating.NumberOfComparisons += 1

		// Add to update list
		updatedRatings = append(updatedRatings, currentTargetRating, *challengerRating)

		// Create comparison history
		comparisonHistory = append(comparisonHistory, domain.ComparisonHistory{
			ID:             uuid.New(),
			UserId:         userId,
			FilmAId:        targetFilmId,
			FilmBId:        comp.ChallengerFilmId,
			WinningFilmId:  winningFilmId,
			ComparisonDate: time.Now(),
			WasEqual:       wasEqual,
		})

		// Update ratingsMap for next iteration
		ratingsMap[comp.ChallengerFilmId] = challengerRating
	}

	// Bulk update all ratings
	if err := s.RatingStore.BulkUpdateRatings(ctx, tx, updatedRatings); err != nil {
		return err
	}

	// Bulk insert comparison history
	if err := s.RatingStore.BulkInsertComparisons(ctx, comparisonHistory); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
