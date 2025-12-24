package films

import (
	"context"
	"errors"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

// Mock FilmStore for testing
type mockFilmStore struct {
	getFilmByIdFunc         func(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	getFilmByExternalIdFunc func(ctx context.Context, id int) (*domain.Film, error)
	createFilmFunc          func(ctx context.Context, film domain.Film) (*domain.Film, error)
	getFilmsForRatingFunc   func(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
}

func (m *mockFilmStore) GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
	if m.getFilmByIdFunc != nil {
		return m.getFilmByIdFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error) {
	if m.getFilmByExternalIdFunc != nil {
		return m.getFilmByExternalIdFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) CreateFilm(ctx context.Context, film domain.Film) (*domain.Film, error) {
	if m.createFilmFunc != nil {
		return m.createFilmFunc(ctx, film)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error) {
	if m.getFilmsForRatingFunc != nil {
		return m.getFilmsForRatingFunc(ctx, userId, filmId)
	}
	return nil, errors.New("not implemented")
}

func TestNewService(t *testing.T) {
	mockStore := &mockFilmStore{}
	service := NewService(mockStore)

	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.FilmStore != mockStore {
		t.Error("expected service to contain the provided store")
	}
}

func TestService_CreateFilm(t *testing.T) {
	ctx := context.Background()
	testFilm := domain.Film{
		ID:          uuid.New(),
		ExternalID:  123,
		Title:       "Test Film",
		Description: "Test Description",
	}

	mockStore := &mockFilmStore{
		createFilmFunc: func(ctx context.Context, film domain.Film) (*domain.Film, error) {
			if film.ID != testFilm.ID {
				t.Errorf("expected film ID %v, got %v", testFilm.ID, film.ID)
			}
			return &film, nil
		},
	}

	service := NewService(mockStore)
	createdFilm, err := service.CreateFilm(ctx, testFilm)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdFilm.Title != testFilm.Title {
		t.Errorf("expected title %s, got %s", testFilm.Title, createdFilm.Title)
	}
}

func TestService_CreateFilm_Error(t *testing.T) {
	ctx := context.Background()
	testFilm := domain.Film{
		ID:    uuid.New(),
		Title: "Test Film",
	}

	expectedError := errors.New("database error")
	mockStore := &mockFilmStore{
		createFilmFunc: func(ctx context.Context, film domain.Film) (*domain.Film, error) {
			return nil, expectedError
		},
	}

	service := NewService(mockStore)
	_, err := service.CreateFilm(ctx, testFilm)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestService_GetFilmById(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()
	testFilm := &domain.Film{
		ID:          testID,
		ExternalID:  456,
		Title:       "Test Film",
		Description: "Test Description",
	}

	mockStore := &mockFilmStore{
		getFilmByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
			if id != testID {
				t.Errorf("expected ID %v, got %v", testID, id)
			}
			return testFilm, nil
		},
	}

	service := NewService(mockStore)
	film, err := service.GetFilmById(ctx, testID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if film.ID != testID {
		t.Errorf("expected film ID %v, got %v", testID, film.ID)
	}
	if film.Title != testFilm.Title {
		t.Errorf("expected title %s, got %s", testFilm.Title, film.Title)
	}
}

func TestService_GetFilmById_NotFound(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	mockStore := &mockFilmStore{
		getFilmByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
			return nil, ErrFilmNotFound
		},
	}

	service := NewService(mockStore)
	_, err := service.GetFilmById(ctx, testID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrFilmNotFound {
		t.Errorf("expected ErrFilmNotFound, got %v", err)
	}
}

func TestService_GetFilmsFromExternal_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	mockStore := &mockFilmStore{}
	service := NewService(mockStore)

	_, err := service.GetFilmsFromExternal(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty query, got nil")
	}
	expectedErrMsg := "cannot obtain films with empty query string"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}
