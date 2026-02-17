package films

import (
	"context"
	"errors"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

// Mock GraphService for testing
type mockGraphService struct {
	addFilmToGraphFunc func(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error
}

func (m *mockGraphService) AddFilmToGraph(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
	if m.addFilmToGraphFunc != nil {
		return m.addFilmToGraphFunc(ctx, userID, film, recommendations)
	}
	return nil // default: no error
}

// Mock FilmStore for testing
type mockFilmStore struct {
	getFilmByIdFunc                 func(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	getFilmByExternalIdFunc         func(ctx context.Context, id int) (*domain.Film, error)
	createFilmFunc                  func(ctx context.Context, film *domain.Film) (*domain.Film, error)
	getFilmsForRatingFunc           func(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
	getFilmRecommendation           func(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error)
	updateFilmRecommendationFunc    func(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error)
	createFilmRecommendationFunc    func(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error)
	getSeenUnratedFilmsFunc         func(ctx context.Context, userId uuid.UUID) ([]domain.Film, error)
	generateFilmRecommendationsFunc func(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error)
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

func (m *mockFilmStore) CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error) {
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

func (m *mockFilmStore) GetFilmRecommendation(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error) {
	if m.getFilmRecommendation != nil {
		return m.getFilmRecommendation(ctx, userId, externalFilmId)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) CreateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
	if m.createFilmRecommendationFunc != nil {
		return m.createFilmRecommendationFunc(ctx, recommendation)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) GenerateFilmRecommendations(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error) {
	if m.generateFilmRecommendationsFunc != nil {
		return m.generateFilmRecommendationsFunc(ctx, userId, films)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error) {
	if m.getSeenUnratedFilmsFunc != nil {
		return m.getSeenUnratedFilmsFunc(ctx, userId)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFilmStore) UpdateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
	if m.updateFilmRecommendationFunc != nil {
		return m.updateFilmRecommendationFunc(ctx, recommendation)
	}
	return nil, errors.New("not implemented")
}

func TestNewService(t *testing.T) {
	mockStore := &mockFilmStore{}
	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)

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
		createFilmFunc: func(ctx context.Context, film *domain.Film) (*domain.Film, error) {
			if film.ID != testFilm.ID {
				t.Errorf("expected film ID %v, got %v", testFilm.ID, film.ID)
			}
			return film, nil
		},
	}

	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)
	createdFilm, err := service.CreateFilm(ctx, &testFilm)

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
		createFilmFunc: func(ctx context.Context, film *domain.Film) (*domain.Film, error) {
			return nil, expectedError
		},
	}

	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)
	_, err := service.CreateFilm(ctx, &testFilm)

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

	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)
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

	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)
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
	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)

	_, err := service.GetFilmsFromExternal(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty query, got nil")
	}
	expectedErrMsg := "cannot obtain films with empty query string"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestService_GenerateFilmRecommendations_DeduplicatesResults(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	// Create test seed films
	seedFilm1 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  100,
		Title:       "The Matrix",
		Description: "A hacker discovers reality",
		PosterUrl:   "/matrix.jpg",
		ReleaseYear: "1999",
	}

	seedFilm2 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  200,
		Title:       "Inception",
		Description: "Dreams within dreams",
		PosterUrl:   "/inception.jpg",
		ReleaseYear: "2010",
	}

	// These films will be recommended by both seed films (duplicates)
	duplicateFilm1 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  300,
		Title:       "Interstellar",
		Description: "Space exploration",
		PosterUrl:   "/interstellar.jpg",
		ReleaseYear: "2014",
	}

	duplicateFilm2 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  400,
		Title:       "Tenet",
		Description: "Time inversion",
		PosterUrl:   "/tenet.jpg",
		ReleaseYear: "2020",
	}

	// Unique recommendations
	uniqueFilm1 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  500,
		Title:       "Blade Runner",
		Description: "Replicants and humanity",
		PosterUrl:   "/bladerunner.jpg",
		ReleaseYear: "1982",
	}

	uniqueFilm2 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  600,
		Title:       "The Prestige",
		Description: "Rival magicians",
		PosterUrl:   "/prestige.jpg",
		ReleaseYear: "2006",
	}

	mockStore := &mockFilmStore{
		createFilmFunc: func(ctx context.Context, film *domain.Film) (*domain.Film, error) {
			return film, nil
		},
		getFilmRecommendation: func(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error) {
			// Return not found for all (simulating fresh recommendations)
			return nil, ErrFilmRecommendationNotFound
		},
		createFilmRecommendationFunc: func(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
			return recommendation, nil
		},
	}

	mockGraph := &mockGraphService{
		addFilmToGraphFunc: func(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
			return nil
		},
	}

	service := NewService(mockStore, mockGraph)

	// Mock the TMDB recommendation function to return predictable duplicates
	service.tmdbRecommendationFunc = func(film domain.Film) []domain.Film {
		switch film.ExternalID {
		case 100: // The Matrix recommends duplicates + unique1
			return []domain.Film{duplicateFilm1, duplicateFilm2, uniqueFilm1}
		case 200: // Inception recommends same duplicates + unique2
			return []domain.Film{duplicateFilm1, duplicateFilm2, uniqueFilm2}
		default:
			return []domain.Film{}
		}
	}

	// Generate recommendations from both seed films
	results, err := service.GenerateFilmRecommendations(ctx, userID, []domain.Film{seedFilm1, seedFilm2})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify no duplicate external IDs in results
	seenExternalIDs := make(map[int]bool)
	for _, film := range results {
		if seenExternalIDs[film.ExternalID] {
			t.Errorf("found duplicate film with ExternalID %d (Title: %s)", film.ExternalID, film.Title)
		}
		seenExternalIDs[film.ExternalID] = true
	}

	// Verify we got exactly 4 unique films (2 duplicates + 2 uniques, not 6)
	expectedCount := 4
	if len(results) != expectedCount {
		t.Errorf("expected %d unique recommendations, got %d", expectedCount, len(results))
	}

	// Verify all expected films are present
	expectedFilmIDs := map[int]bool{
		300: true, // Interstellar (duplicate)
		400: true, // Tenet (duplicate)
		500: true, // Blade Runner (unique to Matrix)
		600: true, // The Prestige (unique to Inception)
	}

	for _, film := range results {
		if !expectedFilmIDs[film.ExternalID] {
			t.Errorf("unexpected film in results: %d - %s", film.ExternalID, film.Title)
		}
	}
}

func TestService_GenerateFilmRecommendations_EmptyFilmList(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockStore := &mockFilmStore{}
	mockGraph := &mockGraphService{}
	service := NewService(mockStore, mockGraph)

	_, err := service.GenerateFilmRecommendations(ctx, userID, []domain.Film{})

	if err == nil {
		t.Fatal("expected error for empty film list, got nil")
	}
	expectedErrMsg := "cannot generate recommendations with empty film list"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestService_GenerateFilmRecommendations_FiltersCircularRecommendations(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	seedFilm := domain.Film{
		ID:          uuid.New(),
		ExternalID:  100,
		Title:       "The Matrix",
		Description: "A hacker discovers reality",
		PosterUrl:   "/matrix.jpg",
		ReleaseYear: "1999",
	}

	// This film will be recommended but already marked as seen
	alreadySeenFilm := domain.Film{
		ID:          uuid.New(),
		ExternalID:  200,
		Title:       "Inception",
		Description: "Dreams within dreams",
		PosterUrl:   "/inception.jpg",
		ReleaseYear: "2010",
	}

	// This film should be included in results
	newRecommendation := domain.Film{
		ID:          uuid.New(),
		ExternalID:  300,
		Title:       "Interstellar",
		Description: "Space exploration",
		PosterUrl:   "/interstellar.jpg",
		ReleaseYear: "2014",
	}

	mockStore := &mockFilmStore{
		createFilmFunc: func(ctx context.Context, film *domain.Film) (*domain.Film, error) {
			return film, nil
		},
		getFilmRecommendation: func(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error) {
			// Simulate that alreadySeenFilm was already marked as seen
			if externalFilmId == 200 {
				return &domain.FilmRecommendation{
					ID:                       uuid.New(),
					UserID:                   userId,
					ExternalFilmID:           externalFilmId,
					HasSeen:                  true,
					HasBeenRecommended:       false,
					RecommendationsGenerated: false,
				}, nil
			}
			// All others are new
			return nil, ErrFilmRecommendationNotFound
		},
		createFilmRecommendationFunc: func(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
			return recommendation, nil
		},
	}

	mockGraph := &mockGraphService{
		addFilmToGraphFunc: func(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
			return nil
		},
	}

	service := NewService(mockStore, mockGraph)

	// Mock TMDB to return both films
	service.tmdbRecommendationFunc = func(film domain.Film) []domain.Film {
		if film.ExternalID == 100 {
			return []domain.Film{alreadySeenFilm, newRecommendation}
		}
		return []domain.Film{}
	}

	results, err := service.GenerateFilmRecommendations(ctx, userID, []domain.Film{seedFilm})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should only return newRecommendation, not alreadySeenFilm
	if len(results) != 1 {
		t.Errorf("expected 1 recommendation (filtered), got %d", len(results))
	}

	if len(results) > 0 && results[0].ExternalID != 300 {
		t.Errorf("expected Interstellar (300), got %d - %s", results[0].ExternalID, results[0].Title)
	}
}

func TestService_GenerateFilmRecommendations_ReturnsEmptyWhenAllFiltered(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	seedFilm := domain.Film{
		ID:          uuid.New(),
		ExternalID:  100,
		Title:       "The Matrix",
		Description: "A hacker discovers reality",
		PosterUrl:   "/matrix.jpg",
		ReleaseYear: "1999",
	}

	// Both recommendations will be filtered out (already seen)
	seenFilm1 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  200,
		Title:       "Inception",
		Description: "Dreams within dreams",
		PosterUrl:   "/inception.jpg",
		ReleaseYear: "2010",
	}

	seenFilm2 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  300,
		Title:       "Interstellar",
		Description: "Space exploration",
		PosterUrl:   "/interstellar.jpg",
		ReleaseYear: "2014",
	}

	mockStore := &mockFilmStore{
		createFilmFunc: func(ctx context.Context, film *domain.Film) (*domain.Film, error) {
			return film, nil
		},
		getFilmRecommendation: func(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error) {
			// Mark all recommended films as already seen
			if externalFilmId == 200 || externalFilmId == 300 {
				return &domain.FilmRecommendation{
					ID:                       uuid.New(),
					UserID:                   userId,
					ExternalFilmID:           externalFilmId,
					HasSeen:                  true,
					HasBeenRecommended:       false,
					RecommendationsGenerated: false,
				}, nil
			}
			return nil, ErrFilmRecommendationNotFound
		},
		createFilmRecommendationFunc: func(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
			return recommendation, nil
		},
	}

	mockGraph := &mockGraphService{
		addFilmToGraphFunc: func(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
			return nil
		},
	}

	service := NewService(mockStore, mockGraph)

	// Mock TMDB to return films that user has already seen
	service.tmdbRecommendationFunc = func(film domain.Film) []domain.Film {
		if film.ExternalID == 100 {
			return []domain.Film{seenFilm1, seenFilm2}
		}
		return []domain.Film{}
	}

	results, err := service.GenerateFilmRecommendations(ctx, userID, []domain.Film{seedFilm})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should return empty array when all recommendations are filtered
	if len(results) != 0 {
		t.Errorf("expected 0 recommendations (all filtered), got %d", len(results))
	}
}

