package films

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type mockFilmService struct {
	getFilmByIdFunc          func(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	getFilmsFromExternalFunc func(ctx context.Context, query string) ([]domain.Film, error)
}

func (m *mockFilmService) GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
	if m.getFilmByIdFunc != nil {
		return m.getFilmByIdFunc(ctx, id)
	}
	return &domain.Film{ID: id, Title: "Test Film"}, nil
}

func (m *mockFilmService) GetFilmsFromExternal(ctx context.Context, query string) ([]domain.Film, error) {
	if m.getFilmsFromExternalFunc != nil {
		return m.getFilmsFromExternalFunc(ctx, query)
	}
	return []domain.Film{{ID: uuid.New(), Title: "External Film"}}, nil
}

func (m *mockFilmService) GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error) {
	return nil, errors.New("not implemented")
}

type mockRatingService struct {
	getAllRatingsFunc              func(ctx context.Context) ([]domain.UserFilmRating, error)
	filterRatingsForComparisonFunc func([]domain.UserFilmRating) []domain.UserFilmRating
}

func (m *mockRatingService) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	if m.getAllRatingsFunc != nil {
		return m.getAllRatingsFunc(ctx)
	}
	filmId := uuid.New()
	return []domain.UserFilmRating{{ID: uuid.New(), FilmId: filmId}}, nil
}

func (m *mockRatingService) FilterRatingsForComparison(ratings []domain.UserFilmRating) []domain.UserFilmRating {
	if m.filterRatingsForComparisonFunc != nil {
		return m.filterRatingsForComparisonFunc(ratings)
	}
	return ratings
}

type mockComparisonService struct {
	hasBeenComparedFunc func(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
}

func (m *mockComparisonService) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	if m.hasBeenComparedFunc != nil {
		return m.hasBeenComparedFunc(ctx, userId, filmAId, filmBId)
	}
	return false, nil
}

func TestNewHandler_Films(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.FilmService != mockFilmSvc {
		t.Error("expected handler to contain the provided service")
	}
}

func TestHandler_GetFilmById_Success(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	filmId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/films/"+filmId.String(), nil)
	req.SetPathValue("id", filmId.String())
	w := httptest.NewRecorder()

	handler.GetFilmById(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_GetFilmById_InvalidUUID(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	req := httptest.NewRequest(http.MethodGet, "/films/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetFilmById(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_GetFilmById_NotFound(t *testing.T) {
	mockFilmSvc := &mockFilmService{
		getFilmByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
			return nil, ErrFilmNotFound
		},
	}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	filmId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/films/"+filmId.String(), nil)
	req.SetPathValue("id", filmId.String())
	w := httptest.NewRecorder()

	handler.GetFilmById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandler_GetFilmById_ServiceError(t *testing.T) {
	mockFilmSvc := &mockFilmService{
		getFilmByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
			return nil, errors.New("database error")
		},
	}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	filmId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/films/"+filmId.String(), nil)
	req.SetPathValue("id", filmId.String())
	w := httptest.NewRecorder()

	handler.GetFilmById(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_GetFilmsFromExternal_Success(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	req := httptest.NewRequest(http.MethodGet, "/films/search?f=inception", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsFromExternal(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_GetFilmsFromExternal_MissingQuery(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	req := httptest.NewRequest(http.MethodGet, "/films/search", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsFromExternal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetFilmsFromExternal_ServiceError(t *testing.T) {
	mockFilmSvc := &mockFilmService{
		getFilmsFromExternalFunc: func(ctx context.Context, query string) ([]domain.Film, error) {
			return nil, errors.New("external API error")
		},
	}
	mockRatingSvc := &mockRatingService{}
	mockComparisonSvc := &mockComparisonService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc, mockComparisonSvc)

	req := httptest.NewRequest(http.MethodGet, "/films/search?f=test", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsFromExternal(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_GetFilmsForRating_Success(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}

	handler := &Handler{
		FilmService:   mockFilmSvc,
		RatingService: mockRatingSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/films/candidates-for-comparison", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsForRating(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestHandler_GetFilmsForRating_RatingServiceError(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{
		getAllRatingsFunc: func(ctx context.Context) ([]domain.UserFilmRating, error) {
			return nil, errors.New("rating service error")
		},
	}

	handler := &Handler{
		FilmService:   mockFilmSvc,
		RatingService: mockRatingSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/films/candidates-for-comparison", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsForRating(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_GetFilmsForRating_FilmNotFound(t *testing.T) {
	// Test that handler gracefully handles films that can't be found
	mockFilmSvc := &mockFilmService{
		getFilmByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
			return nil, ErrFilmNotFound
		},
	}
	mockRatingSvc := &mockRatingService{}

	handler := &Handler{
		FilmService:   mockFilmSvc,
		RatingService: mockRatingSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/films/candidates-for-comparison", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsForRating(w, req)

	// Should still return 200 with empty array (or films that were found)
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
