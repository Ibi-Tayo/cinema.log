package users

import (
	"context"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	// validation errors
	ErrUserNameInvalidLength = errors.New("name not between 5 and 20 characters")
	ErrNoId                  = errors.New("user ID is required")
	ErrInvalidId             = errors.New("invalid user ID format")
	ErrInvalidJson           = errors.New("invalid JSON format")
	//server errors
	ErrEncoding = errors.New("error encoding response")
	ErrServer   = errors.New("internal server error")
)

type service struct {
	store Store
}

type Store interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetOrCreateUserByGithubId(ctx context.Context, githubId int64, name string,
		username string, avatarUrl string) (*domain.User, error)
	GetOrCreateUserByGoogleId(ctx context.Context, googleId string, name string,
		username string, avatarUrl string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

func NewService(store Store) UserService {
	return &service{
		store: store,
	}
}

func (s *service) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.store.GetAllUsers(ctx)
}

func (s *service) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.store.GetUserById(ctx, id)
}

func (s *service) GetOrCreateUserByGithubId(ctx context.Context, githubId int64,
	name string, username string, avatarUrl string) (*domain.User, error) {
	return s.store.GetOrCreateUserByGithubId(ctx, githubId, name, username, avatarUrl)
}

func (s *service) GetOrCreateUserByGoogleId(ctx context.Context, googleId string,
	name string, username string, avatarUrl string) (*domain.User, error) {
	return s.store.GetOrCreateUserByGoogleId(ctx, googleId, name, username, avatarUrl)
}

func (s *service) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validation logic
	if len(user.Name) < 5 || len(user.Name) > 20 {
		return nil, ErrUserNameInvalidLength
	}

	if (len(user.Username) < 5 || len(user.Username) > 20) && user.Username != "" {
		return nil, ErrUserNameInvalidLength
	}
	return s.store.CreateUser(ctx, user)
}

func (s *service) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validation logic
	if len(user.Name) < 5 || len(user.Name) > 20 {
		return nil, ErrUserNameInvalidLength
	}

	if (len(user.Username) < 5 || len(user.Username) > 20) && user.Username != "" {
		return nil, ErrUserNameInvalidLength
	}
	return s.store.UpdateUser(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteUser(ctx, id)
}
