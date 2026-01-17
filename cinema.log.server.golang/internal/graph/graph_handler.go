package graph

import (
	"context"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

type Handler struct {
	GraphService GraphService
}

type GraphService interface {
	GetUserGraph(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, []domain.FilmGraphEdge, error)
}

func NewHandler(graphService GraphService) *Handler {
	return &Handler{
		GraphService: graphService,
	}
}

type GetUserGraphResponse struct {
	Nodes []domain.FilmGraphNode `json:"nodes"`
	Edges []domain.FilmGraphEdge `json:"edges"`
}

// GetUserGraph returns the film graph for the authenticated user
func (h *Handler) GetUserGraph(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middleware.KeyUser).(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	nodes, edges, err := h.GraphService.GetUserGraph(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve graph", http.StatusInternalServerError)
		return
	}

	response := GetUserGraphResponse{
		Nodes: nodes,
		Edges: edges,
	}

	utils.SendJSON(w, response)
}
