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
	createFilmFunc                  func(ctx context.Context, film *domain.Film) (*domain.Film, error)
	getFilmByIdFunc                 func(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	getFilmsFromExternalFunc        func(ctx context.Context, query string) ([]domain.Film, error)
	generateFilmRecommendationsFunc func(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error)
	getSeenUnratedFilmsFunc         func(ctx context.Context, userId uuid.UUID) ([]domain.Film, error)
}

func (m *mockFilmService) GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error) {
	if m.getFilmsFromExternalFunc != nil {
		return m.getFilmsFromExternalFunc(ctx, "")
	}
	return []domain.Film{{ID: uuid.New(), Title: "Film for Rating"}}, nil
}

func (m *mockFilmService) CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error) {
	if m.createFilmFunc != nil {
		return m.createFilmFunc(ctx, film)
	}
	return film, nil
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

func (m *mockFilmService) GenerateFilmRecommendations(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error) {
	if m.generateFilmRecommendationsFunc != nil {
		return m.generateFilmRecommendationsFunc(ctx, userId, films)
	}
	return films, nil
}

func (m *mockFilmService) GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error) {
	if m.getSeenUnratedFilmsFunc != nil {
		return m.getSeenUnratedFilmsFunc(ctx, userId)
	}
	return []domain.Film{{ID: uuid.New(), Title: "Seen Unrated Film"}}, nil
}

type mockRatingService struct {
	getAllRatingsFunc              func(ctx context.Context) ([]domain.UserFilmRating, error)
	filterRatingsForComparisonFunc func([]domain.UserFilmRating) []domain.UserFilmRating
	hasBeenComparedFunc            func(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
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

func (m *mockRatingService) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	if m.hasBeenComparedFunc != nil {
		return m.hasBeenComparedFunc(ctx, userId, filmAId, filmBId)
	}
	return false, nil
}

func TestNewHandler_Films(t *testing.T) {
	mockFilmSvc := &mockFilmService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

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
	handler := NewHandler(mockFilmSvc, mockRatingSvc)

	req := httptest.NewRequest(http.MethodGet, "/films/search?f=test", nil)
	w := httptest.NewRecorder()

	handler.GetFilmsFromExternal(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
