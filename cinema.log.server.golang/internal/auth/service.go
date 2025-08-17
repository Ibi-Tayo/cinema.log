package auth

import (
	"context"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/users"
	"github.com/google/go-github/v52/github"
)

type AuthService struct {
	userService users.UserService
}

func NewService(userService users.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) HandleGithubCallback(ctx context.Context, githubUser *github.User) (*domain.User, string, string, error) {
	// 1. Insert user into DB if doesn't exist
	// 2. Generate JWT

	// we return the user, the jwt and error
	return nil, "", "", nil
}

func (s *AuthService) GenerateJWT(user *domain.User) (string, error) {
	return "", nil
}