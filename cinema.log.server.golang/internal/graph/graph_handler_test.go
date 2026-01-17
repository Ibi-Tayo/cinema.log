package graph

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGraphService struct {
	mock.Mock
}

func (m *MockGraphService) GetUserGraph(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, []domain.FilmGraphEdge, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.FilmGraphNode), args.Get(1).([]domain.FilmGraphEdge), args.Error(2)
}

func (m *MockGraphService) AddFilmToGraph(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
	args := m.Called(ctx, userID, film, recommendations)
	return args.Error(0)
}

func TestNewHandler(t *testing.T) {
	mockSvc := &MockGraphService{}
	handler := NewHandler(mockSvc)

	assert.NotNil(t, handler)
	assert.Equal(t, mockSvc, handler.GraphService)
}

func TestHandler_GetUserGraph_Success(t *testing.T) {
	mockSvc := new(MockGraphService)
	handler := NewHandler(mockSvc)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Name:     "Test User",
		Username: "testuser",
	}

	expectedNodes := []domain.FilmGraphNode{
		{
			UserID:         userID,
			ExternalFilmID: 123,
			Title:          "Film 1",
		},
		{
			UserID:         userID,
			ExternalFilmID: 456,
			Title:          "Film 2",
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

	mockSvc.On("GetUserGraph", mock.Anything, userID).Return(expectedNodes, expectedEdges, nil)

	req := httptest.NewRequest(http.MethodGet, "/graph", nil)
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetUserGraph(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetUserGraphResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Nodes, 2)
	assert.Len(t, response.Edges, 1)
	assert.Equal(t, expectedNodes[0].Title, response.Nodes[0].Title)
	assert.Equal(t, expectedEdges[0].FromFilmID, response.Edges[0].FromFilmID)

	mockSvc.AssertExpectations(t)
}

func TestHandler_GetUserGraph_Unauthorized(t *testing.T) {
	mockSvc := new(MockGraphService)
	handler := NewHandler(mockSvc)

	// Request without user in context
	req := httptest.NewRequest(http.MethodGet, "/graph", nil)
	w := httptest.NewRecorder()

	handler.GetUserGraph(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertNotCalled(t, "GetUserGraph")
}

func TestHandler_GetUserGraph_ServiceError(t *testing.T) {
	mockSvc := new(MockGraphService)
	handler := NewHandler(mockSvc)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Name:     "Test User",
		Username: "testuser",
	}

	mockSvc.On("GetUserGraph", mock.Anything, userID).Return(
		[]domain.FilmGraphNode{},
		[]domain.FilmGraphEdge{},
		assert.AnError,
	)

	req := httptest.NewRequest(http.MethodGet, "/graph", nil)
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetUserGraph(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_GetUserGraph_EmptyGraph(t *testing.T) {
	mockSvc := new(MockGraphService)
	handler := NewHandler(mockSvc)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Name:     "Test User",
		Username: "testuser",
	}

	emptyNodes := []domain.FilmGraphNode{}
	emptyEdges := []domain.FilmGraphEdge{}

	mockSvc.On("GetUserGraph", mock.Anything, userID).Return(emptyNodes, emptyEdges, nil)

	req := httptest.NewRequest(http.MethodGet, "/graph", nil)
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetUserGraph(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetUserGraphResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Nodes, 0)
	assert.Len(t, response.Edges, 0)

	mockSvc.AssertExpectations(t)
}
