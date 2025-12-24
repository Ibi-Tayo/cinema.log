package ratings

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

var (
	testDB      *sql.DB
	testStore   RatingStore
	testDbSetup *utils.TestDatabase
)

// Helper function to create test user
func createTestUser(ctx context.Context, t *testing.T) uuid.UUID {
	userID := uuid.New()
	query := `INSERT INTO users (user_id, name, username, github_id, profile_pic_url) 
	          VALUES ($1, $2, $3, $4, $5)`
	githubID := int(time.Now().UnixNano() % 2147483647) // Use nanoseconds for uniqueness
	_, err := testDB.ExecContext(ctx, query, userID, "Test User", "testuser"+userID.String()[:8], githubID, "http://example.com/pic.jpg")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return userID
}

// Helper function to create test film
func createTestFilm(ctx context.Context, t *testing.T) uuid.UUID {
	filmID := uuid.New()
	query := `INSERT INTO films (film_id, external_id, title, description, poster_url, release_year) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	externalID := int(time.Now().UnixNano() % 2147483647) // Use nanoseconds for uniqueness
	_, err := testDB.ExecContext(ctx, query, filmID, externalID, "Test Film "+filmID.String()[:8], "Description", "/poster.jpg", "2024")
	if err != nil {
		t.Fatalf("failed to create test film: %v", err)
	}
	return filmID
}

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

func TestRatingStore_NewStore(t *testing.T) {
	store := NewStore(testDB)
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestRatingStore_CreateRating(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	rating := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId,
		EloRating:            1500.0,
		NumberOfComparisons:  0,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	createdRating, err := testStore.CreateRating(ctx, rating)
	if err != nil {
		t.Fatalf("failed to create rating: %v", err)
	}

	if createdRating.ID != rating.ID {
		t.Errorf("expected rating ID %v, got %v", rating.ID, createdRating.ID)
	}
	if createdRating.EloRating != rating.EloRating {
		t.Errorf("expected elo rating %.2f, got %.2f", rating.EloRating, createdRating.EloRating)
	}
	if createdRating.UserId != rating.UserId {
		t.Errorf("expected user ID %v, got %v", rating.UserId, createdRating.UserId)
	}
	if createdRating.FilmId != rating.FilmId {
		t.Errorf("expected film ID %v, got %v", rating.FilmId, createdRating.FilmId)
	}
}

func TestRatingStore_GetRating(t *testing.T) {
	ctx := context.Background()
	
	// Create a rating first
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	rating := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId,
		EloRating:            1600.0,
		NumberOfComparisons:  5,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	_, err := testStore.CreateRating(ctx, rating)
	if err != nil {
		t.Fatalf("failed to create rating: %v", err)
	}

	// Now retrieve it
	retrievedRating, err := testStore.GetRating(ctx, userId, filmId)
	if err != nil {
		t.Fatalf("failed to get rating: %v", err)
	}

	if retrievedRating.UserId != userId {
		t.Errorf("expected user ID %v, got %v", userId, retrievedRating.UserId)
	}
	if retrievedRating.FilmId != filmId {
		t.Errorf("expected film ID %v, got %v", filmId, retrievedRating.FilmId)
	}
	if retrievedRating.EloRating != rating.EloRating {
		t.Errorf("expected elo rating %.2f, got %.2f", rating.EloRating, retrievedRating.EloRating)
	}
}

func TestRatingStore_GetRating_NotFound(t *testing.T) {
	ctx := context.Background()
	
	nonExistentUserId := uuid.New()
	nonExistentFilmId := uuid.New()
	
	_, err := testStore.GetRating(ctx, nonExistentUserId, nonExistentFilmId)
	
	if err == nil {
		t.Fatal("expected error for non-existent rating")
	}
	if err != ErrRatingNotFound {
		t.Errorf("expected ErrRatingNotFound, got %v", err)
	}
}

func TestRatingStore_GetAllRatings(t *testing.T) {
	ctx := context.Background()
	
	// Create multiple ratings
	userId1 := createTestUser(ctx, t)
	filmId1 := createTestFilm(ctx, t)
	userId2 := createTestUser(ctx, t)
	filmId2 := createTestFilm(ctx, t)
	
	rating1 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId1,
		FilmId:               filmId1,
		EloRating:            1500.0,
		NumberOfComparisons:  0,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}
	
	rating2 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId2,
		FilmId:               filmId2,
		EloRating:            1700.0,
		NumberOfComparisons:  10,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	_, err := testStore.CreateRating(ctx, rating1)
	if err != nil {
		t.Fatalf("failed to create rating1: %v", err)
	}

	_, err = testStore.CreateRating(ctx, rating2)
	if err != nil {
		t.Fatalf("failed to create rating2: %v", err)
	}

	// Get all ratings
	ratings, err := testStore.GetAllRatings(ctx)
	if err != nil {
		t.Fatalf("failed to get all ratings: %v", err)
	}

	if len(ratings) < 2 {
		t.Errorf("expected at least 2 ratings, got %d", len(ratings))
	}
}

func TestRatingStore_GetRatingsByUserId(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId1 := createTestFilm(ctx, t)
	filmId2 := createTestFilm(ctx, t)
	
	// Create ratings for specific user
	rating1 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId1,
		EloRating:            1500.0,
		NumberOfComparisons:  0,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}
	
	rating2 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId2,
		EloRating:            1600.0,
		NumberOfComparisons:  5,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	_, err := testStore.CreateRating(ctx, rating1)
	if err != nil {
		t.Fatalf("failed to create rating1: %v", err)
	}

	_, err = testStore.CreateRating(ctx, rating2)
	if err != nil {
		t.Fatalf("failed to create rating2: %v", err)
	}

	// Get ratings by user ID
	ratings, err := testStore.GetRatingsByUserId(ctx, userId)
	if err != nil {
		t.Fatalf("failed to get ratings by user ID: %v", err)
	}

	if len(ratings) < 2 {
		t.Errorf("expected at least 2 ratings for user, got %d", len(ratings))
	}

	// Verify all ratings belong to the user
	for _, rating := range ratings {
		if rating.UserId != userId {
			t.Errorf("expected user ID %v, got %v", userId, rating.UserId)
		}
	}
}

func TestRatingStore_UpdateRating(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	// Create a rating first
	rating := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId,
		EloRating:            1500.0,
		NumberOfComparisons:  0,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	createdRating, err := testStore.CreateRating(ctx, rating)
	if err != nil {
		t.Fatalf("failed to create rating: %v", err)
	}

	// Update the rating
	createdRating.EloRating = 1650.0
	createdRating.NumberOfComparisons = 10
	createdRating.LastUpdated = time.Now()
	createdRating.KConstantValue = 24.0

	updatedRating, err := testStore.UpdateRating(ctx, *createdRating)
	if err != nil {
		t.Fatalf("failed to update rating: %v", err)
	}

	if updatedRating.EloRating != 1650.0 {
		t.Errorf("expected elo rating 1650.0, got %.2f", updatedRating.EloRating)
	}
	if updatedRating.NumberOfComparisons != 10 {
		t.Errorf("expected 10 comparisons, got %d", updatedRating.NumberOfComparisons)
	}
	if updatedRating.KConstantValue != 24.0 {
		t.Errorf("expected k constant 24.0, got %.2f", updatedRating.KConstantValue)
	}
}

func TestRatingStore_UpdateRating_NotFound(t *testing.T) {
	ctx := context.Background()
	
	nonExistentRating := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               uuid.New(),
		FilmId:               uuid.New(),
		EloRating:            1500.0,
		NumberOfComparisons:  0,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	_, err := testStore.UpdateRating(ctx, nonExistentRating)
	
	if err == nil {
		t.Fatal("expected error for non-existent rating")
	}
	if err != ErrRatingNotFound {
		t.Errorf("expected ErrRatingNotFound, got %v", err)
	}
}

func TestRatingStore_UpdateRatings(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId1 := createTestFilm(ctx, t)
	filmId2 := createTestFilm(ctx, t)
	
	// Create two ratings
	rating1 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId1,
		EloRating:            1500.0,
		NumberOfComparisons:  5,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}
	
	rating2 := domain.UserFilmRating{
		ID:                   uuid.New(),
		UserId:               userId,
		FilmId:               filmId2,
		EloRating:            1600.0,
		NumberOfComparisons:  5,
		LastUpdated:          time.Now(),
		InitialRating:        1500.0,
		KConstantValue:       32.0,
	}

	createdRating1, err := testStore.CreateRating(ctx, rating1)
	if err != nil {
		t.Fatalf("failed to create rating1: %v", err)
	}

	createdRating2, err := testStore.CreateRating(ctx, rating2)
	if err != nil {
		t.Fatalf("failed to create rating2: %v", err)
	}

	// Update both ratings
	createdRating1.EloRating = 1550.0
	createdRating1.NumberOfComparisons = 6
	createdRating2.EloRating = 1650.0
	createdRating2.NumberOfComparisons = 6

	comparisonPair := domain.ComparisonPair{
		FilmA: *createdRating1,
		FilmB: *createdRating2,
	}

	updatedPair, err := testStore.UpdateRatings(ctx, comparisonPair)
	if err != nil {
		t.Fatalf("failed to update ratings: %v", err)
	}

	if updatedPair.FilmA.EloRating != 1550.0 {
		t.Errorf("expected FilmA elo rating 1550.0, got %.2f", updatedPair.FilmA.EloRating)
	}
	if updatedPair.FilmB.EloRating != 1650.0 {
		t.Errorf("expected FilmB elo rating 1650.0, got %.2f", updatedPair.FilmB.EloRating)
	}
}
