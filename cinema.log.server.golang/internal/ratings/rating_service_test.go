package ratings

import (
	"testing"
	"time"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

func TestService_FilterRatingsForComparison(t *testing.T) {
	service := NewService(nil) // No store needed for this test

	tests := []struct {
		name     string
		ratings  []domain.UserFilmRating
		expected int
	}{
		{
			name: "returns all when less than 10",
			ratings: []domain.UserFilmRating{
				{NumberOfComparisons: 0, LastUpdated: time.Now().Add(-2 * time.Hour)},
				{NumberOfComparisons: 1, LastUpdated: time.Now().Add(-1 * time.Hour)},
				{NumberOfComparisons: 0, LastUpdated: time.Now()},
			},
			expected: 3,
		},
		{
			name: "returns first 10 when more than 10",
			ratings: func() []domain.UserFilmRating {
				ratings := make([]domain.UserFilmRating, 15)
				for i := 0; i < 15; i++ {
					ratings[i] = domain.UserFilmRating{
						NumberOfComparisons: i,
						LastUpdated:         time.Now().Add(-time.Duration(i) * time.Hour),
					}
				}
				return ratings
			}(),
			expected: 10,
		},
		{
			name:     "returns empty when no ratings",
			ratings:  []domain.UserFilmRating{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FilterRatingsForComparison(tt.ratings)
			if len(result) != tt.expected {
				t.Errorf("expected %d ratings, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestService_FilterRatingsForComparison_Sorting(t *testing.T) {
	service := NewService(nil)

	now := time.Now()
	ratings := []domain.UserFilmRating{
		{ID: uuid.New(), NumberOfComparisons: 5, LastUpdated: now.Add(-1 * time.Hour)},
		{ID: uuid.New(), NumberOfComparisons: 0, LastUpdated: now.Add(-3 * time.Hour)}, // Should be first (least comparisons, oldest)
		{ID: uuid.New(), NumberOfComparisons: 0, LastUpdated: now.Add(-2 * time.Hour)}, // Should be second
		{ID: uuid.New(), NumberOfComparisons: 3, LastUpdated: now},
	}

	result := service.FilterRatingsForComparison(ratings)

	if result[0].NumberOfComparisons != 0 {
		t.Errorf("expected first rating to have 0 comparisons, got %d", result[0].NumberOfComparisons)
	}

	if result[0].LastUpdated.After(result[1].LastUpdated) {
		t.Error("expected ratings with same comparison count to be sorted by oldest first")
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
		name      string
		winnerId  uuid.UUID
		expectedA float64
		expectedB float64
	}{
		{
			name:      "film A wins",
			winnerId:  filmA,
			expectedA: 1.0,
			expectedB: 0.0,
		},
		{
			name:      "film B wins",
			winnerId:  filmB,
			expectedA: 0.0,
			expectedB: 1.0,
		},
		{
			name:      "draw",
			winnerId:  uuid.New(), // Different ID
			expectedA: 0.5,
			expectedB: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultA, resultB := service.defineFilmContestResult(filmA, filmB, tt.winnerId)
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
			name:                "new rating (0-4 comparisons)",
			numberOfComparisons: 2,
			expectedK:           40,
		},
		{
			name:                "intermediate (5-9 comparisons)",
			numberOfComparisons: 7,
			expectedK:           20,
		},
		{
			name:                "established (10+ comparisons)",
			numberOfComparisons: 15,
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
			expectedElo:   1400,
		},
		{
			name:          "2.5 star rating",
			initialRating: 2.5,
			expectedElo:   1500,
		},
		{
			name:          "3.5 star rating",
			initialRating: 3.5,
			expectedElo:   1600,
		},
		{
			name:          "5 star rating",
			initialRating: 5.0,
			expectedElo:   1700,
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
