package graph

import (
	"context"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGraphStore struct {
	mock.Mock
}

func (m *MockGraphStore) AddNode(ctx context.Context, node *domain.FilmGraphNode) error {
	args := m.Called(ctx, node)
	return args.Error(0)
}

func (m *MockGraphStore) NodeExists(ctx context.Context, userID uuid.UUID, externalFilmID int) (bool, error) {
	args := m.Called(ctx, userID, externalFilmID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGraphStore) AddEdge(ctx context.Context, edge *domain.FilmGraphEdge) error {
	args := m.Called(ctx, edge)
	return args.Error(0)
}

func (m *MockGraphStore) GetNodesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.FilmGraphNode), args.Error(1)
}

func (m *MockGraphStore) GetEdgesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphEdge, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.FilmGraphEdge), args.Error(1)
}

type MockFilmStore struct {
	mock.Mock
}

func (m *MockFilmStore) GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Film), args.Error(1)
}

func TestAddFilmToGraph_Success(t *testing.T) {
	mockGraphStore := new(MockGraphStore)
	mockFilmStore := new(MockFilmStore)
	service := NewService(mockGraphStore, mockFilmStore)

	ctx := context.Background()
	userID := uuid.New()
	filmID := uuid.New()
	recFilmID1 := uuid.New()
	recFilmID2 := uuid.New()

	film := domain.Film{
		ID:         filmID,
		ExternalID: 123,
		Title:      "Test Film",
	}

	recommendations := []domain.Film{
		{
			ID:         recFilmID1,
			ExternalID: 456,
			Title:      "Rec Film 1",
		},
		{
			ID:         recFilmID2,
			ExternalID: 789,
			Title:      "Rec Film 2",
		},
	}

	// Mock film lookups
	mockFilmStore.On("GetFilmByExternalId", ctx, 123).Return(&film, nil)
	mockFilmStore.On("GetFilmByExternalId", ctx, 456).Return(&recommendations[0], nil)
	mockFilmStore.On("GetFilmByExternalId", ctx, 789).Return(&recommendations[1], nil)

	// Mock graph operations
	mockGraphStore.On("AddNode", ctx, mock.AnythingOfType("*domain.FilmGraphNode")).Return(nil)
	mockGraphStore.On("NodeExists", ctx, userID, 456).Return(true, nil)
	mockGraphStore.On("NodeExists", ctx, userID, 789).Return(false, nil)
	mockGraphStore.On("AddEdge", ctx, mock.AnythingOfType("*domain.FilmGraphEdge")).Return(nil)

	err := service.AddFilmToGraph(ctx, userID, film, recommendations)

	assert.NoError(t, err)
	mockGraphStore.AssertExpectations(t)
	mockFilmStore.AssertExpectations(t)
}

func TestGetUserGraph_Success(t *testing.T) {
	mockGraphStore := new(MockGraphStore)
	mockFilmStore := new(MockFilmStore)
	service := NewService(mockGraphStore, mockFilmStore)

	ctx := context.Background()
	userID := uuid.New()

	expectedNodes := []domain.FilmGraphNode{
		{
			UserID:         userID,
			ExternalFilmID: 123,
			Title:          "Film 1",
		},
	}

	expectedEdges := []domain.FilmGraphEdge{
		{
			UserID:     userID,
			EdgeId:     uuid.New(),
			FromFilmID: 123,
			ToFilmID:   456,
		},
	}

	mockGraphStore.On("GetNodesByUser", ctx, userID).Return(expectedNodes, nil)
	mockGraphStore.On("GetEdgesByUser", ctx, userID).Return(expectedEdges, nil)

	nodes, edges, err := service.GetUserGraph(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedNodes, nodes)
	assert.Equal(t, expectedEdges, edges)
	mockGraphStore.AssertExpectations(t)
}
