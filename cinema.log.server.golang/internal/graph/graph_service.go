package graph

import (
	"context"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	GraphStore GraphStore
	FilmStore  FilmStore
}

type GraphStore interface {
	AddNode(ctx context.Context, node *domain.FilmGraphNode) error
	NodeExists(ctx context.Context, userID uuid.UUID, externalFilmID int) (bool, error)
	AddEdge(ctx context.Context, edge *domain.FilmGraphEdge) error
	GetNodesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, error)
	GetEdgesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphEdge, error)
}

type FilmStore interface {
	GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error)
}

func NewService(graphStore GraphStore, filmStore FilmStore) *Service {
	return &Service{
		GraphStore: graphStore,
		FilmStore:  filmStore,
	}
}

// AddFilmToGraph adds a film to the user's graph and creates connections to existing films
// based on recommendations. This should be called when a user confirms they've seen a film.
func (s *Service) AddFilmToGraph(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error {
	_, err := s.FilmStore.GetFilmByExternalId(ctx, film.ExternalID)
	if err != nil {
		return err
	}

	node := &domain.FilmGraphNode{
		UserID:         userID,
		ExternalFilmID: film.ExternalID,
		Title:          film.Title,
	}

	if err := s.GraphStore.AddNode(ctx, node); err != nil {
		return err
	}

	// Check each recommendation to see if it's already in the user's graph
	// If it is, create a bidirectional edge between the films
	for _, recFilm := range recommendations {
		_, err := s.FilmStore.GetFilmByExternalId(ctx, recFilm.ExternalID)
		if err != nil {
			// Skip this recommendation if we can't find it in the database
			continue
		}

		exists, err := s.GraphStore.NodeExists(ctx, userID, recFilm.ExternalID)
		if err != nil {
			continue
		}

		if exists {
			// Create bidirectional edge using external_id
			edge := &domain.FilmGraphEdge{
				UserID:     userID,
				EdgeId:     uuid.New(),
				FromFilmID: film.ExternalID,
				ToFilmID:   recFilm.ExternalID,
			}
			if err := s.GraphStore.AddEdge(ctx, edge); err != nil {
				continue
			}
		}
	}

	return nil
}

// GetUserGraph returns all nodes and edges for a user's film graph
func (s *Service) GetUserGraph(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, []domain.FilmGraphEdge, error) {
	nodes, err := s.GraphStore.GetNodesByUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	edges, err := s.GraphStore.GetEdgesByUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return nodes, edges, nil
}
