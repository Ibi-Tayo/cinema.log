package ratings

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"github.com/google/uuid"
)

type mockRatingService struct{}

func (m *mockRatingService) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	return &domain.UserFilmRating{ID: uuid.New(), UserId: userId, FilmId: filmId}, nil
}

func (m *mockRatingService) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRatingService) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, comparison domain.ComparisonHistory) (*domain.ComparisonPair, error) {
	return &ratings, nil
}

func (m *mockRatingService) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
	return &comparison, nil
}

func (m *mockRatingService) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockRatingService) GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error) {
	return nil, errors.New("not implemented")
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

func TestNewHandler(t *testing.T) {
	mockSvc := &mockRatingService{}
	handler := NewHandler(mockSvc)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.RatingService != mockSvc {
		t.Error("expected handler to contain the provided service")
	}
}

func TestHandler_GetRating_Success(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	userId := uuid.New()
	filmId := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/ratings?userId="+userId.String()+"&filmId="+filmId.String(), nil)
	w := httptest.NewRecorder()

	handler.GetRating(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_GetRating_MissingFilmId(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	userId := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/ratings?userId="+userId.String(), nil)
	w := httptest.NewRecorder()

	handler.GetRating(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetRating_InvalidFilmId(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	userId := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/ratings?userId="+userId.String()+"&filmId=invalid", nil)
	w := httptest.NewRecorder()

	handler.GetRating(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_CompareFilms_Success(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	userId := uuid.New()
	filmAId := uuid.New()
	filmBId := uuid.New()
	user := &domain.User{ID: userId, Name: "Test User", Username: "testuser"}

	comparison := `{
		"userId": "` + userId.String() + `",
		"filmAId": "` + filmAId.String() + `",
		"filmBId": "` + filmBId.String() + `",
		"winningFilmId": "` + filmAId.String() + `"
	}`

	req := httptest.NewRequest(http.MethodPost, "/ratings/compare-films", strings.NewReader(comparison))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CompareFilms(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestHandler_CompareFilms_InvalidJSON(t *testing.T) {
	handler := NewHandler(&mockRatingService{})
	userId := uuid.New()
	user := &domain.User{ID: userId, Name: "Test User", Username: "testuser"}

	req := httptest.NewRequest(http.MethodPost, "/ratings/compare-films", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CompareFilms(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
