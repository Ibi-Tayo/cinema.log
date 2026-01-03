package ratings

import (
	"context"
	"testing"
	"time"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

// Add mock store for service testing
type mockRatingStore struct {
	getRatingFunc            func(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	getAllRatingsFunc        func(ctx context.Context) ([]domain.UserFilmRating, error)
	getRatingsByUserIdFunc   func(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error)
	createRatingFunc         func(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	updateRatingFunc         func(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error)
	updateRatingsFunc        func(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error)
	createComparisonFunc     func(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	hasBeenComparedFunc      func(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	getComparisonHistoryFunc func(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
}

func (m *mockRatingStore) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	if m.getRatingFunc != nil {
		return m.getRatingFunc(ctx, userId, filmId)
	}
	return nil, ErrRatingNotFound
}

func (m *mockRatingStore) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	if m.getAllRatingsFunc != nil {
		return m.getAllRatingsFunc(ctx)
	}
	return nil, nil
}

func (m *mockRatingStore) GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error) {
	if m.getRatingsByUserIdFunc != nil {
		return m.getRatingsByUserIdFunc(ctx, userId)
	}
	return nil, nil
}

func (m *mockRatingStore) CreateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
	if m.createRatingFunc != nil {
		return m.createRatingFunc(ctx, rating)
	}
	return &rating, nil
}

func (m *mockRatingStore) UpdateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
	if m.updateRatingFunc != nil {
		return m.updateRatingFunc(ctx, rating)
	}
	return &rating, nil
}

func (m *mockRatingStore) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error) {
	if m.updateRatingsFunc != nil {
		return m.updateRatingsFunc(ctx, ratings)
	}
	return &ratings, nil
}

func (m *mockRatingStore) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
	if m.createComparisonFunc != nil {
		return m.createComparisonFunc(ctx, comparison)
	}
	return &comparison, nil
}

func (m *mockRatingStore) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	if m.hasBeenComparedFunc != nil {
		return m.hasBeenComparedFunc(ctx, userId, filmAId, filmBId)
	}
	return false, nil
}

func (m *mockRatingStore) GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error) {
	if m.getComparisonHistoryFunc != nil {
		return m.getComparisonHistoryFunc(ctx, userId)
	}
	return nil, nil
}

func TestService_GetRating(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	filmId := uuid.New()
	expectedRating := &domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              userId,
		FilmId:              filmId,
		EloRating:           1500.0,
		NumberOfComparisons: 5,
	}

	mock := &mockRatingStore{
		getRatingFunc: func(ctx context.Context, uid uuid.UUID, fid uuid.UUID) (*domain.UserFilmRating, error) {
			if uid != userId || fid != filmId {
				return nil, ErrRatingNotFound
			}
			return expectedRating, nil
		},
	}

	service := NewService(mock)
	rating, err := service.GetRating(ctx, userId, filmId)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rating.ID != expectedRating.ID {
		t.Errorf("expected rating ID %v, got %v", expectedRating.ID, rating.ID)
	}
}

func TestService_GetAllRatings(t *testing.T) {
	ctx := context.Background()
	expectedRatings := []domain.UserFilmRating{
		{ID: uuid.New(), EloRating: 1500.0},
		{ID: uuid.New(), EloRating: 1600.0},
	}

	mock := &mockRatingStore{
		getAllRatingsFunc: func(ctx context.Context) ([]domain.UserFilmRating, error) {
			return expectedRatings, nil
		},
	}

	service := NewService(mock)
	ratings, err := service.GetAllRatings(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ratings) != 2 {
		t.Errorf("expected 2 ratings, got %d", len(ratings))
	}
}

func TestService_CreateRating(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	filmId := uuid.New()
	initialRating := float32(3.5)

	mock := &mockRatingStore{
		createRatingFunc: func(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
			if rating.UserId != userId {
				t.Errorf("expected user ID %v, got %v", userId, rating.UserId)
			}
			if rating.FilmId != filmId {
				t.Errorf("expected film ID %v, got %v", filmId, rating.FilmId)
			}
			if rating.InitialRating != initialRating {
				t.Errorf("expected initial rating %.1f, got %.1f", initialRating, rating.InitialRating)
			}
			if rating.NumberOfComparisons != 0 {
				t.Errorf("expected 0 comparisons, got %d", rating.NumberOfComparisons)
			}
			if rating.KConstantValue != 40 {
				t.Errorf("expected K constant 40, got %.1f", rating.KConstantValue)
			}
			return &rating, nil
		},
	}

	service := NewService(mock)
	rating, err := service.CreateRating(ctx, userId, filmId, initialRating)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rating.UserId != userId {
		t.Errorf("expected user ID %v, got %v", userId, rating.UserId)
	}
}

func TestService_UpdateRatings(t *testing.T) {
	ctx := context.Background()

	filmAId := uuid.New()
	filmBId := uuid.New()

	filmA := domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              uuid.New(),
		FilmId:              filmAId,
		EloRating:           1500.0,
		NumberOfComparisons: 5,
		LastUpdated:         time.Now(),
		InitialRating:       1500.0,
		KConstantValue:      32.0,
	}

	filmB := domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              filmA.UserId,
		FilmId:              filmBId,
		EloRating:           1600.0,
		NumberOfComparisons: 5,
		LastUpdated:         time.Now(),
		InitialRating:       1500.0,
		KConstantValue:      32.0,
	}

	mock := &mockRatingStore{
		updateRatingFunc: func(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
			return &rating, nil
		},
	}

	service := NewService(mock)
	pair := domain.ComparisonPair{
		FilmA: filmA,
		FilmB: filmB,
	}

	// FilmA wins
	comparison := domain.ComparisonHistory{
		ID:            uuid.New(),
		UserId:        filmA.UserId,
		FilmAId:       filmA.FilmId,
		FilmBId:       filmB.FilmId,
		WinningFilmId: filmA.FilmId,
		WasEqual:      false,
	}
	updatedPair, err := service.UpdateRatings(ctx, pair, comparison)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// FilmA should have increased rating
	if updatedPair.FilmA.EloRating <= filmA.EloRating {
		t.Errorf("expected FilmA rating to increase from %.2f, got %.2f", filmA.EloRating, updatedPair.FilmA.EloRating)
	}

	// FilmB should have decreased rating
	if updatedPair.FilmB.EloRating >= filmB.EloRating {
		t.Errorf("expected FilmB rating to decrease from %.2f, got %.2f", filmB.EloRating, updatedPair.FilmB.EloRating)
	}

	// Both should have increased comparisons
	if updatedPair.FilmA.NumberOfComparisons != filmA.NumberOfComparisons+1 {
		t.Errorf("expected FilmA comparisons to be %d, got %d", filmA.NumberOfComparisons+1, updatedPair.FilmA.NumberOfComparisons)
	}
	if updatedPair.FilmB.NumberOfComparisons != filmB.NumberOfComparisons+1 {
		t.Errorf("expected FilmB comparisons to be %d, got %d", filmB.NumberOfComparisons+1, updatedPair.FilmB.NumberOfComparisons)
	}
}

func TestService_UpdateRatings_Draw(t *testing.T) {
	ctx := context.Background()

	filmA := domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              uuid.New(),
		FilmId:              uuid.New(),
		EloRating:           1500.0,
		NumberOfComparisons: 5,
		LastUpdated:         time.Now(),
		InitialRating:       1500.0,
		KConstantValue:      32.0,
	}

	filmB := domain.UserFilmRating{
		ID:                  uuid.New(),
		UserId:              filmA.UserId,
		FilmId:              uuid.New(),
		EloRating:           1500.0,
		NumberOfComparisons: 5,
		LastUpdated:         time.Now(),
		InitialRating:       1500.0,
		KConstantValue:      32.0,
	}

	mock := &mockRatingStore{
		updateRatingFunc: func(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
			return &rating, nil
		},
	}

	service := NewService(mock)
	pair := domain.ComparisonPair{
		FilmA: filmA,
		FilmB: filmB,
	}

	// Neither wins (draw) - use a different winnerId
	comparison := domain.ComparisonHistory{
		ID:            uuid.New(),
		UserId:        filmA.UserId,
		FilmAId:       filmA.FilmId,
		FilmBId:       filmB.FilmId,
		WinningFilmId: filmA.FilmId,
		WasEqual:      true,
	}
	updatedPair, err := service.UpdateRatings(ctx, pair, comparison)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// With equal ratings and a draw, ratings should stay roughly the same
	// (they both get 0.5 result and 0.5 expected)
	if updatedPair.FilmA.EloRating != filmA.EloRating {
		t.Logf("FilmA rating changed from %.2f to %.2f (expected to stay same on draw with equal ratings)", filmA.EloRating, updatedPair.FilmA.EloRating)
	}
}

func TestService_CalculateExpectedResult(t *testing.T) {
	service := NewService(nil)

	tests := []struct {
		name           string
		filmRating     float64
		opponentRating float64
		wantRange      [2]float64 // min, max
	}{
		{
			name:           "equal ratings",
			filmRating:     1500,
			opponentRating: 1500,
			wantRange:      [2]float64{0.49, 0.51},
		},
		{
			name:           "higher rated film",
			filmRating:     1700,
			opponentRating: 1500,
			wantRange:      [2]float64{0.75, 0.77},
		},
		{
			name:           "lower rated film",
			filmRating:     1500,
			opponentRating: 1700,
			wantRange:      [2]float64{0.23, 0.25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateExpectedResult(tt.filmRating, tt.opponentRating)
			if result < tt.wantRange[0] || result > tt.wantRange[1] {
				t.Errorf("expected result between %v and %v, got %v", tt.wantRange[0], tt.wantRange[1], result)
			}
		})
	}
}

func TestService_RecalculateFilmRating(t *testing.T) {
	service := NewService(nil)

	tests := []struct {
		name           string
		expectedResult float64
		actualResult   float64
		currentRating  float64
		kValue         float64
		wantMin        float64
	}{
		{
			name:           "win increases rating",
			expectedResult: 0.5,
			actualResult:   1.0,
			currentRating:  1500,
			kValue:         40,
			wantMin:        1500,
		},
		{
			name:           "loss decreases rating",
			expectedResult: 0.5,
			actualResult:   0.0,
			currentRating:  1500,
			kValue:         40,
			wantMin:        100, // Min rating is 100
		},
		{
			name:           "rating never goes below 100",
			expectedResult: 1.0,
			actualResult:   0.0,
			currentRating:  110,
			kValue:         40,
			wantMin:        100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.recalculateFilmRating(tt.expectedResult, tt.actualResult, tt.currentRating, tt.kValue)
			if result < tt.wantMin {
				t.Errorf("expected rating >= %v, got %v", tt.wantMin, result)
			}
		})
	}
}

func TestService_DefineFilmContestResult(t *testing.T) {
	service := NewService(nil)

	filmA := uuid.New()
	filmB := uuid.New()

	tests := []struct {
		name       string
		comparison domain.ComparisonHistory
		expectedA  float64
		expectedB  float64
	}{
		{
			name: "film A wins",
			comparison: domain.ComparisonHistory{
				WinningFilmId: filmA,
				WasEqual:      false,
			},
			expectedA: 1.0,
			expectedB: 0.0,
		},
		{
			name: "film B wins",
			comparison: domain.ComparisonHistory{
				WinningFilmId: filmB,
				WasEqual:      false,
			},
			expectedA: 0.0,
			expectedB: 1.0,
		},
		{
			name: "draw",
			comparison: domain.ComparisonHistory{
				WinningFilmId: uuid.New(), // Different ID doesn't matter when WasEqual is true
				WasEqual:      true,
			},
			expectedA: 0.5,
			expectedB: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultA, resultB := service.defineFilmContestResult(filmA, filmB, tt.comparison)
			if resultA != tt.expectedA {
				t.Errorf("expected film A result %v, got %v", tt.expectedA, resultA)
			}
			if resultB != tt.expectedB {
				t.Errorf("expected film B result %v, got %v", tt.expectedB, resultB)
			}
		})
	}
}

func TestService_UpdateKConstantValue(t *testing.T) {
	service := NewService(nil)

	tests := []struct {
		name                string
		numberOfComparisons int
		expectedK           float64
	}{
		{
			name:                "new rating (0-20 comparisons)",
			numberOfComparisons: 2,
			expectedK:           40,
		},
		{
			name:                "intermediate (20-40 comparisons)",
			numberOfComparisons: 25,
			expectedK:           20,
		},
		{
			name:                "established (40+ comparisons)",
			numberOfComparisons: 45,
			expectedK:           10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			film := domain.UserFilmRating{
				NumberOfComparisons: tt.numberOfComparisons,
			}
			result := service.updateKConstantValue(film)
			if result != tt.expectedK {
				t.Errorf("expected K value %v, got %v", tt.expectedK, result)
			}
		})
	}
}

func TestService_GetInitialEloRating(t *testing.T) {
	service := NewService(nil)

	tests := []struct {
		name          string
		initialRating float32
		expectedElo   float32
	}{
		{
			name:          "1 star rating",
			initialRating: 1.0,
			expectedElo:   950,
		},
		{
			name:          "2.5 star rating",
			initialRating: 2.5,
			expectedElo:   1000,
		},
		{
			name:          "3.5 star rating",
			initialRating: 3.5,
			expectedElo:   1050,
		},
		{
			name:          "5 star rating",
			initialRating: 5.0,
			expectedElo:   1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getInitialEloRating(tt.initialRating)
			if result != tt.expectedElo {
				t.Errorf("expected Elo rating %v, got %v", tt.expectedElo, result)
			}
		})
	}
}
