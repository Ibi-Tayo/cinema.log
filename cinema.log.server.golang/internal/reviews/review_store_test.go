package reviews

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
	testStore   ReviewStore
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

func TestReviewStore_NewStore(t *testing.T) {
	store := NewStore(testDB)
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestReviewStore_CreateReview(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	review := domain.Review{
		ID:      uuid.New(),
		Content: "Great movie!",
		Date:    time.Now(),
		Rating:  5.0,
		FilmId:  filmId,
		UserId:  userId,
	}

	createdReview, err := testStore.CreateReview(ctx, review)
	if err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	if createdReview.ID != review.ID {
		t.Errorf("expected review ID %v, got %v", review.ID, createdReview.ID)
	}
	if createdReview.Content != review.Content {
		t.Errorf("expected content %s, got %s", review.Content, createdReview.Content)
	}
	if createdReview.Rating != review.Rating {
		t.Errorf("expected rating %.1f, got %.1f", review.Rating, createdReview.Rating)
	}
}

func TestReviewStore_CreateReview_GeneratesUUID(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	review := domain.Review{
		ID:      uuid.Nil, // No UUID provided
		Content: "Another great movie!",
		Date:    time.Now(),
		Rating:  4.0,
		FilmId:  filmId,
		UserId:  userId,
	}

	createdReview, err := testStore.CreateReview(ctx, review)
	if err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	if createdReview.ID == uuid.Nil {
		t.Error("expected non-nil UUID to be generated")
	}
}

func TestReviewStore_GetAllReviewsByUserId(t *testing.T) {
	ctx := context.Background()
	userId := createTestUser(ctx, t)
	filmId1 := createTestFilm(ctx, t)
	filmId2 := createTestFilm(ctx, t)
	
	// Create multiple reviews for the same user
	review1 := domain.Review{
		ID:      uuid.New(),
		Content: "First review",
		Date:    time.Now(),
		Rating:  5.0,
		FilmId:  filmId1,
		UserId:  userId,
	}
	
	review2 := domain.Review{
		ID:      uuid.New(),
		Content: "Second review",
		Date:    time.Now(),
		Rating:  4.0,
		FilmId:  filmId2,
		UserId:  userId,
	}

	_, err := testStore.CreateReview(ctx, review1)
	if err != nil {
		t.Fatalf("failed to create review1: %v", err)
	}

	_, err = testStore.CreateReview(ctx, review2)
	if err != nil {
		t.Fatalf("failed to create review2: %v", err)
	}

	// Get all reviews by user ID
	reviews, err := testStore.GetAllReviewsByUserId(ctx, userId)
	if err != nil {
		t.Fatalf("failed to get reviews: %v", err)
	}

	if len(reviews) < 2 {
		t.Errorf("expected at least 2 reviews, got %d", len(reviews))
	}

	// Verify all reviews belong to the user
	for _, review := range reviews {
		if review.UserId != userId {
			t.Errorf("expected user ID %v, got %v", userId, review.UserId)
		}
	}
}

func TestReviewStore_UpdateReview(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	// Create a review first
	review := domain.Review{
		ID:      uuid.New(),
		Content: "Initial review",
		Date:    time.Now(),
		Rating:  3.0,
		FilmId:  filmId,
		UserId:  userId,
	}

	createdReview, err := testStore.CreateReview(ctx, review)
	if err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	// Update the review
	createdReview.Content = "Updated review"
	createdReview.Rating = 5.0
	createdReview.Date = time.Now()

	updatedReview, err := testStore.UpdateReview(ctx, *createdReview)
	if err != nil {
		t.Fatalf("failed to update review: %v", err)
	}

	if updatedReview.Content != "Updated review" {
		t.Errorf("expected content 'Updated review', got %s", updatedReview.Content)
	}
	if updatedReview.Rating != 5.0 {
		t.Errorf("expected rating 5.0, got %.1f", updatedReview.Rating)
	}
}

func TestReviewStore_UpdateReview_NotFound(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	nonExistentReview := domain.Review{
		ID:      uuid.New(),
		Content: "Non-existent review",
		Date:    time.Now(),
		Rating:  3.0,
		FilmId:  filmId,
		UserId:  userId,
	}

	_, err := testStore.UpdateReview(ctx, nonExistentReview)
	
	if err == nil {
		t.Fatal("expected error for non-existent review")
	}
	if err != ErrReviewNotFound {
		t.Errorf("expected ErrReviewNotFound, got %v", err)
	}
}

func TestReviewStore_DeleteReview(t *testing.T) {
	ctx := context.Background()
	
	userId := createTestUser(ctx, t)
	filmId := createTestFilm(ctx, t)
	
	// Create a review first
	review := domain.Review{
		ID:      uuid.New(),
		Content: "Review to delete",
		Date:    time.Now(),
		Rating:  4.0,
		FilmId:  filmId,
		UserId:  userId,
	}

	createdReview, err := testStore.CreateReview(ctx, review)
	if err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	// Delete the review
	err = testStore.DeleteReview(ctx, createdReview.ID)
	if err != nil {
		t.Fatalf("failed to delete review: %v", err)
	}

	// Verify the review is deleted by trying to get it
	reviews, err := testStore.GetAllReviewsByUserId(ctx, createdReview.UserId)
	if err != nil {
		t.Fatalf("failed to get reviews: %v", err)
	}

	for _, r := range reviews {
		if r.ID == createdReview.ID {
			t.Error("expected review to be deleted, but it still exists")
		}
	}
}

func TestReviewStore_DeleteReview_NotFound(t *testing.T) {
	ctx := context.Background()
	
	nonExistentID := uuid.New()
	err := testStore.DeleteReview(ctx, nonExistentID)
	
	if err == nil {
		t.Fatal("expected error for non-existent review")
	}
	if err != ErrReviewNotFound {
		t.Errorf("expected ErrReviewNotFound, got %v", err)
	}
}
