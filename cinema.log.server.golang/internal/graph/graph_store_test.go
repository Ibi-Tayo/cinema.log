package graph

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB      *sql.DB
	testStore   *Store
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

// Helper function to create test user
func createTestUser(ctx context.Context, t *testing.T, userID uuid.UUID) {
	query := `INSERT INTO users (user_id, name, username, github_id, profile_pic_url) 
	          VALUES ($1, $2, $3, $4, $5)`
	githubID := int(uuid.New().ID() % 2147483647) // Use UUID for uniqueness
	_, err := testDB.ExecContext(ctx, query, userID, "Test User", "testuser"+userID.String()[:8], githubID, "http://example.com/pic.jpg")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}

// Helper function to create test film
func createTestFilm(ctx context.Context, t *testing.T, filmID uuid.UUID) int {
	query := `INSERT INTO films (film_id, external_id, title, description, poster_url, release_year) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	externalID := int(uuid.New().ID() % 2147483647) // Use UUID for uniqueness
	_, err := testDB.ExecContext(ctx, query, filmID, externalID, "Test Film "+filmID.String()[:8], "Description", "/poster.jpg", "2024")
	if err != nil {
		t.Fatalf("failed to create test film: %v", err)
	}
	return externalID
}

func TestStore_AddNode(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	userID := uuid.New()
	filmID := uuid.New()

	// Create test user and film first
	createTestUser(ctx, t, userID)
	externalFilmID := createTestFilm(ctx, t, filmID)

	node := &domain.FilmGraphNode{
		UserID:         userID,
		ExternalFilmID: externalFilmID,
		Title:          "Test Film",
	}

	err := store.AddNode(ctx, node)
	require.NoError(t, err)

	// Verify node was added
	exists, err := store.NodeExists(ctx, userID, externalFilmID)
	require.NoError(t, err)
	assert.True(t, exists)

	// Test idempotency - adding same node again should not error
	err = store.AddNode(ctx, node)
	require.NoError(t, err)
}

func TestStore_NodeExists(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	userID := uuid.New()
	filmID := uuid.New()

	// Create test user and film first
	createTestUser(ctx, t, userID)
	externalFilmID := createTestFilm(ctx, t, filmID)

	// Node should not exist initially
	exists, err := store.NodeExists(ctx, userID, externalFilmID)
	require.NoError(t, err)
	assert.False(t, exists)

	// Add node
	node := &domain.FilmGraphNode{
		UserID:         userID,
		ExternalFilmID: externalFilmID,
		Title:          "Test Film",
	}
	err = store.AddNode(ctx, node)
	require.NoError(t, err)

	// Node should now exist
	exists, err = store.NodeExists(ctx, userID, externalFilmID)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestStore_AddEdge(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	userID := uuid.New()
	filmID1 := uuid.New()
	filmID2 := uuid.New()

	// Create test user and films first
	createTestUser(ctx, t, userID)
	fromFilmExtID := createTestFilm(ctx, t, filmID1)
	toFilmExtID := createTestFilm(ctx, t, filmID2)

	// Add nodes first
	node1 := &domain.FilmGraphNode{
		UserID:         userID,
		ExternalFilmID: fromFilmExtID,
		Title:          "Film 1",
	}
	node2 := &domain.FilmGraphNode{
		UserID:         userID,
		ExternalFilmID: toFilmExtID,
		Title:          "Film 2",
	}
	err := store.AddNode(ctx, node1)
	require.NoError(t, err)
	err = store.AddNode(ctx, node2)
	require.NoError(t, err)

	// Add edge
	edge := &domain.FilmGraphEdge{
		UserID:     userID,
		EdgeId:     uuid.New(),
		FromFilmID: fromFilmExtID,
		ToFilmID:   toFilmExtID,
	}
	err = store.AddEdge(ctx, edge)
	require.NoError(t, err)

	// Verify edge exists (only one direction returned)
	edges, err := store.GetEdgesByUser(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, edges, 1) // Should return only one direction
}

func TestStore_GetNodesByUser(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()
	filmID1 := uuid.New()
	filmID2 := uuid.New()

	// Create test users and films first
	createTestUser(ctx, t, userID1)
	createTestUser(ctx, t, userID2)
	externalFilmID1 := createTestFilm(ctx, t, filmID1)
	externalFilmID2 := createTestFilm(ctx, t, filmID2)

	// Add nodes for user 1
	node1 := &domain.FilmGraphNode{
		UserID:         userID1,
		ExternalFilmID: externalFilmID1,
		Title:          "Film 1",
	}
	node2 := &domain.FilmGraphNode{
		UserID:         userID1,
		ExternalFilmID: externalFilmID2,
		Title:          "Film 2",
	}
	err := store.AddNode(ctx, node1)
	require.NoError(t, err)
	err = store.AddNode(ctx, node2)
	require.NoError(t, err)

	// Add node for user 2
	node3 := &domain.FilmGraphNode{
		UserID:         userID2,
		ExternalFilmID: externalFilmID1,
		Title:          "Film 1",
	}
	err = store.AddNode(ctx, node3)
	require.NoError(t, err)

	// Get nodes for user 1
	nodes, err := store.GetNodesByUser(ctx, userID1)
	require.NoError(t, err)
	assert.Len(t, nodes, 2)

	// Get nodes for user 2
	nodes, err = store.GetNodesByUser(ctx, userID2)
	require.NoError(t, err)
	assert.Len(t, nodes, 1)
}

func TestStore_GetEdgesByUser(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	userID := uuid.New()
	filmUUID1 := uuid.New()
	filmUUID2 := uuid.New()
	filmUUID3 := uuid.New()

	// Create test user and films first
	createTestUser(ctx, t, userID)
	filmID1 := createTestFilm(ctx, t, filmUUID1)
	filmID2 := createTestFilm(ctx, t, filmUUID2)
	filmID3 := createTestFilm(ctx, t, filmUUID3)

	// Add nodes
	for _, filmID := range []int{filmID1, filmID2, filmID3} {
		node := &domain.FilmGraphNode{
			UserID:         userID,
			ExternalFilmID: filmID,
			Title:          "Film",
		}
		err := store.AddNode(ctx, node)
		require.NoError(t, err)
	}

	// Add edges
	edge1 := &domain.FilmGraphEdge{
		UserID:     userID,
		EdgeId:     uuid.New(),
		FromFilmID: filmID1,
		ToFilmID:   filmID2,
	}
	edge2 := &domain.FilmGraphEdge{
		UserID:     userID,
		EdgeId:     uuid.New(),
		FromFilmID: filmID2,
		ToFilmID:   filmID3,
	}
	err := store.AddEdge(ctx, edge1)
	require.NoError(t, err)
	err = store.AddEdge(ctx, edge2)
	require.NoError(t, err)

	// Get edges (only one direction of each edge returned)
	edges, err := store.GetEdgesByUser(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, edges, 2) // 2 edges, one direction each
}
