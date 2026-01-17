package graph

import (
	"context"
	"database/sql"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrNodeNotFound = errors.New("film graph node not found")
	ErrEdgeNotFound = errors.New("film graph edge not found")
	ErrServer       = errors.New("internal server error")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// AddNode adds a film to a user's graph (idempotent - won't error if already exists)
func (s *Store) AddNode(ctx context.Context, node *domain.FilmGraphNode) error {
	query := `
		INSERT INTO film_graph_nodes (user_id, external_film_id, title)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, external_film_id) DO NOTHING
	`
	_, err := s.db.ExecContext(ctx, query, node.UserID, node.ExternalFilmID, node.Title)
	if err != nil {
		return err
	}
	return nil
}

// NodeExists checks if a film exists in a user's graph
func (s *Store) NodeExists(ctx context.Context, userID uuid.UUID, externalFilmID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM film_graph_nodes 
			WHERE user_id = $1 AND external_film_id = $2
		)
	`
	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, externalFilmID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AddEdge adds a bidirectional connection between two films in a user's graph
func (s *Store) AddEdge(ctx context.Context, edge *domain.FilmGraphEdge) error {
	query := `
		INSERT INTO film_graph_edges (user_id, edge_id, from_film_id, to_film_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (edge_id) DO NOTHING
	`

	// Add edge in both directions for bidirectional relationship
	_, err := s.db.ExecContext(ctx, query, edge.UserID, edge.EdgeId, edge.FromFilmID, edge.ToFilmID)
	if err != nil {
		return err
	}

	// Add reverse edge
	_, err = s.db.ExecContext(ctx, query, edge.UserID, uuid.New(), edge.ToFilmID, edge.FromFilmID)
	if err != nil {
		return err
	}

	return nil
}

// GetNodesByUser returns all film graph nodes for a specific user
func (s *Store) GetNodesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphNode, error) {
	query := `
		SELECT user_id, external_film_id, title
		FROM film_graph_nodes
		WHERE user_id = $1
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]domain.FilmGraphNode, 0)
	for rows.Next() {
		var node domain.FilmGraphNode
		if err := rows.Scan(&node.UserID, &node.ExternalFilmID, &node.Title); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

// GetEdgesByUser returns all edges for a specific user
// Returns only one direction of each bidirectional edge to avoid duplicates
func (s *Store) GetEdgesByUser(ctx context.Context, userID uuid.UUID) ([]domain.FilmGraphEdge, error) {
	query := `
		SELECT DISTINCT user_id, edge_id, from_film_id, to_film_id
		FROM film_graph_edges
		WHERE user_id = $1 AND from_film_id < to_film_id
		ORDER BY from_film_id, to_film_id
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	edges := make([]domain.FilmGraphEdge, 0)
	for rows.Next() {
		var edge domain.FilmGraphEdge
		if err := rows.Scan(&edge.UserID, &edge.EdgeId, &edge.FromFilmID, &edge.ToFilmID); err != nil {
			return nil, err
		}
		edges = append(edges, edge)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return edges, nil
}
