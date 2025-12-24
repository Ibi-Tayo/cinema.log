package ratings

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type mockRatingService struct{}

func (m *mockRatingService) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	return &domain.UserFilmRating{ID: uuid.New(), UserId: userId, FilmId: filmId}, nil
}

func (m *mockRatingService) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRatingService) GetRatingsForComparison(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRating, error) {
	return []domain.UserFilmRating{{ID: uuid.New()}}, nil
}

func (m *mockRatingService) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, winnerId uuid.UUID) (*domain.ComparisonPair, error) {
	return &ratings, nil
}

func TestHandler_GetRating_MissingUserId(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	req := httptest.NewRequest(http.MethodGet, "/ratings?filmId="+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	handler.GetRating(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetRating_InvalidUserId(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	req := httptest.NewRequest(http.MethodGet, "/ratings?userId=invalid&filmId="+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	handler.GetRating(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetRatingsForComparison_MissingUserId(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	req := httptest.NewRequest(http.MethodGet, "/ratings/for-comparison", nil)
	w := httptest.NewRecorder()
	handler.GetRatingsForComparison(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
