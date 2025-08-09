package users

import (
	"context"

	"cinema.log.server.golang/internal/database"
	"cinema.log.server.golang/internal/domain"
)

type store struct {
	db database.Service
}

func NewStore(db database.Service) Store {
	return &store{
		db: db,
	}
}

func (s *store) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	// For now, return a mock user. You'll implement the actual database query later
	user := &domain.User{
		ID:            id,
		GithubID:      "mock123",
		Name:          "Mock User",
		Username:      "mockuser",
		ProfilePicURL: "https://example.com/pic.jpg",
	}
	return user, nil
}