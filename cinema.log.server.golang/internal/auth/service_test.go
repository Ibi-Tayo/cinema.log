package auth

import (
	"context"
	"errors"
	"os"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type mockUserService struct{}

func (m *mockUserService) GetOrCreateUserByGithubId(ctx context.Context, githubId int64, name, username, profilePicURL string) (*domain.User, error) {
	return &domain.User{ID: uuid.New(), Name: name, Username: username}, nil
}

func (m *mockUserService) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return &domain.User{ID: id, Name: "Test", Username: "test"}, nil
}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return errors.New("not implemented")
}

func TestMain(m *testing.M) {
	os.Setenv("TOKEN_SECRET", "test-secret-key")
	code := m.Run()
	os.Exit(code)
}

func TestAuthService_GenerateJWT(t *testing.T) {
	service := NewService(&mockUserService{})
	user := &domain.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Username: "testuser",
	}

	jwtToken, refreshToken, err := service.GenerateJWT(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if jwtToken == "" {
		t.Error("expected JWT token to be non-empty")
	}

	if refreshToken == "" {
		t.Error("expected refresh token to be non-empty")
	}
}

func TestAuthService_ValidateJWT_Success(t *testing.T) {
	userId := uuid.New()
	expectedUser := &domain.User{
		ID:       userId,
		Name:     "Test User",
		Username: "testuser",
	}

	service := NewService(&mockUserService{})

	jwtToken, _, err := service.GenerateJWT(expectedUser)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	user, err := service.ValidateJWT(jwtToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %v, got %v", expectedUser.ID, user.ID)
	}
}

func TestAuthService_ValidateJWT_InvalidToken(t *testing.T) {
	service := NewService(&mockUserService{})

	_, err := service.ValidateJWT("invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
