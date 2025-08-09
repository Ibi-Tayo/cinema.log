package users

import (
	"context"

	"cinema.log.server.golang/internal/domain"
)

type service struct {
	store Store
}

type Store interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	// Add other store methods as needed
}

func NewService(store Store) UserService {
	return &service{
		store: store,
	}
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.store.GetByID(ctx, id)
}