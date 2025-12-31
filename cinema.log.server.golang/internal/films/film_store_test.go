package films

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

var (
	testDB      *sql.DB
	testStore   FilmStore
	testDbSetup *utils.TestDatabase
)

func TestMain(m *testing.M) {
	var err error
	testDbSetup, err = utils.StartTestPostgres()
	if err != nil {
		log.Fatalf("could not start test database: %v", err)
	}

	testDB = testDbSetup.DB
	testStore = NewStore(testDB)

	code := m.Run()

	testDbSetup.Close()
	os.Exit(code)
}

func TestNewStore(t *testing.T) {
	store := NewStore(testDB)
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestCreateFilm(t *testing.T) {
	ctx := context.Background()

	film := domain.Film{
		ID:          uuid.New(),
		ExternalID:  123456,
		Title:       "Test Film",
		Description: "A test film description",
		PosterUrl:   "/test-poster.jpg",
		ReleaseYear: "2024",
	}

	createdFilm, err := testStore.CreateFilm(ctx, &film)
	if err != nil {
		t.Fatalf("failed to create film: %v", err)
	}

	if createdFilm.ID != film.ID {
		t.Errorf("expected film ID %v, got %v", film.ID, createdFilm.ID)
	}
	if createdFilm.Title != film.Title {
		t.Errorf("expected title %s, got %s", film.Title, createdFilm.Title)
	}
	if createdFilm.ExternalID != film.ExternalID {
		t.Errorf("expected external ID %d, got %d", film.ExternalID, createdFilm.ExternalID)
	}
}

func TestCreateFilm_GeneratesUUID(t *testing.T) {
	ctx := context.Background()

	film := domain.Film{
		ID:          uuid.Nil, // No UUID provided
		ExternalID:  789012,
		Title:       "Test Film Without UUID",
		Description: "Testing UUID generation",
		PosterUrl:   "/test-poster-2.jpg",
		ReleaseYear: "2024",
	}

	createdFilm, err := testStore.CreateFilm(ctx, &film)
	if err != nil {
		t.Fatalf("failed to create film: %v", err)
	}

	if createdFilm.ID == uuid.Nil {
		t.Error("expected non-nil UUID to be generated")
	}
}

func TestGetFilmById(t *testing.T) {
	ctx := context.Background()

	// First create a film
	film := domain.Film{
		ID:          uuid.New(),
		ExternalID:  345678,
		Title:       "Test Film Get By ID",
		Description: "Testing get by ID",
		PosterUrl:   "/test-poster-3.jpg",
		ReleaseYear: "2024",
	}

	createdFilm, err := testStore.CreateFilm(ctx, &film)
	if err != nil {
		t.Fatalf("failed to create film: %v", err)
	}

	// Now retrieve it
	retrievedFilm, err := testStore.GetFilmById(ctx, createdFilm.ID)
	if err != nil {
		t.Fatalf("failed to get film: %v", err)
	}

	if retrievedFilm.ID != createdFilm.ID {
		t.Errorf("expected film ID %v, got %v", createdFilm.ID, retrievedFilm.ID)
	}
	if retrievedFilm.Title != createdFilm.Title {
		t.Errorf("expected title %s, got %s", createdFilm.Title, retrievedFilm.Title)
	}
}

func TestGetFilmById_NotFound(t *testing.T) {
	ctx := context.Background()

	nonExistentID := uuid.New()
	_, err := testStore.GetFilmById(ctx, nonExistentID)

	if err == nil {
		t.Fatal("expected error for non-existent film")
	}
	if err != ErrFilmNotFound {
		t.Errorf("expected ErrFilmNotFound, got %v", err)
	}
}

func TestGetFilmByExternalId(t *testing.T) {
	ctx := context.Background()

	// First create a film
	film := domain.Film{
		ID:          uuid.New(),
		ExternalID:  567890,
		Title:       "Test Film Get By External ID",
		Description: "Testing get by external ID",
		PosterUrl:   "/test-poster-4.jpg",
		ReleaseYear: "2024",
	}

	createdFilm, err := testStore.CreateFilm(ctx, &film)
	if err != nil {
		t.Fatalf("failed to create film: %v", err)
	}

	// Now retrieve it by external ID
	retrievedFilm, err := testStore.GetFilmByExternalId(ctx, createdFilm.ExternalID)
	if err != nil {
		t.Fatalf("failed to get film by external ID: %v", err)
	}

	if retrievedFilm.ExternalID != createdFilm.ExternalID {
		t.Errorf("expected external ID %d, got %d", createdFilm.ExternalID, retrievedFilm.ExternalID)
	}
	if retrievedFilm.Title != createdFilm.Title {
		t.Errorf("expected title %s, got %s", createdFilm.Title, retrievedFilm.Title)
	}
}

func TestGetFilmByExternalId_NotFound(t *testing.T) {
	ctx := context.Background()

	nonExistentExternalID := 999999999
	_, err := testStore.GetFilmByExternalId(ctx, nonExistentExternalID)

	if err == nil {
		t.Fatal("expected error for non-existent film")
	}
	if err != ErrFilmNotFound {
		t.Errorf("expected ErrFilmNotFound, got %v", err)
	}
}

func TestCreateFilm_UpsertBehavior(t *testing.T) {
	ctx := context.Background()

	// Create a film with a specific external_id
	film1 := domain.Film{
		ID:          uuid.New(),
		ExternalID:  111222333,
		Title:       "Original Title",
		Description: "Original Description",
		PosterUrl:   "/original-poster.jpg",
		ReleaseYear: "2024",
	}

	createdFilm1, err := testStore.CreateFilm(ctx, &film1)
	if err != nil {
		t.Fatalf("failed to create first film: %v", err)
	}

	// Try to create another film with the same external_id but different UUID
	film2 := domain.Film{
		ID:          uuid.New(), // Different UUID
		ExternalID:  111222333,  // Same external_id
		Title:       "Updated Title",
		Description: "Updated Description",
		PosterUrl:   "/updated-poster.jpg",
		ReleaseYear: "2025",
	}

	createdFilm2, err := testStore.CreateFilm(ctx, &film2)
	if err != nil {
		t.Fatalf("failed to create second film (upsert): %v", err)
	}

	// The returned film should have the ORIGINAL UUID (from film1)
	// but the UPDATED fields (from film2)
	if createdFilm2.ID != createdFilm1.ID {
		t.Errorf("expected UUID %v (from first insert), got %v", createdFilm1.ID, createdFilm2.ID)
	}
	if createdFilm2.Title != film2.Title {
		t.Errorf("expected title to be updated to %s, got %s", film2.Title, createdFilm2.Title)
	}
	if createdFilm2.Description != film2.Description {
		t.Errorf("expected description to be updated to %s, got %s", film2.Description, createdFilm2.Description)
	}

	// Verify only one record exists for this external_id
	retrievedFilm, err := testStore.GetFilmByExternalId(ctx, 111222333)
	if err != nil {
		t.Fatalf("failed to retrieve film: %v", err)
	}
	if retrievedFilm.ID != createdFilm1.ID {
		t.Errorf("expected film to have original UUID %v, got %v", createdFilm1.ID, retrievedFilm.ID)
	}
	if retrievedFilm.Title != film2.Title {
		t.Errorf("expected film title to be updated to %s, got %s", film2.Title, retrievedFilm.Title)
	}
}
