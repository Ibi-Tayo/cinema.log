package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cinema.log.server.golang/internal/auth"
	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

func TestIsAuthExempt(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/auth/github-login", true},
		{"/auth/github-callback", true},
		{"/auth/refresh-token", true},
		{"/users", false},
		{"/films", false},
		{"/ratings", false},
		{"/auth/logout", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isAuthExempt(tt.path)
			if result != tt.expected {
				t.Errorf("isAuthExempt(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	server := &Server{}
	handler := server.corsMiddleware(nextHandler)

	t.Run("adds CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Error("expected Access-Control-Allow-Origin header to be set")
		}
		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("expected Access-Control-Allow-Methods header to be set")
		}
		if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
			t.Error("expected Access-Control-Allow-Credentials header to be 'true'")
		}
	})

	t.Run("uses default frontend URL when not set", func(t *testing.T) {
		oldURL := os.Getenv("FRONTEND_URL")
		os.Unsetenv("FRONTEND_URL")
		defer os.Setenv("FRONTEND_URL", oldURL)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:4200" {
			t.Errorf("expected default origin 'http://localhost:4200', got %q", origin)
		}
	})

	t.Run("uses custom frontend URL when set", func(t *testing.T) {
		oldURL := os.Getenv("FRONTEND_URL")
		os.Setenv("FRONTEND_URL", "https://example.com")
		defer os.Setenv("FRONTEND_URL", oldURL)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "https://example.com" {
			t.Errorf("expected origin 'https://example.com', got %q", origin)
		}
	})

	t.Run("handles OPTIONS preflight requests", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d for OPTIONS, got %d", http.StatusNoContent, w.Code)
		}
	})
}

func TestAuthMiddleware_ExemptPaths(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	server := &Server{}
	handler := server.authMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/auth/github-login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_NoCookie(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	server := &Server{}
	handler := server.authMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Create a mock user service
	mockUserService := &mockUserServiceForAuth{}

	// Set up the test environment with a token secret
	oldSecret := os.Getenv("TOKEN_SECRET")
	os.Setenv("TOKEN_SECRET", "test-secret")
	defer os.Setenv("TOKEN_SECRET", oldSecret)

	authService := auth.NewService(mockUserService)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	server := &Server{authService: authService}
	handler := server.authMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.AddCookie(&http.Cookie{
		Name:  "cinema-log-access-token",
		Value: "invalid-token",
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Mock user service for auth middleware tests
type mockUserServiceForAuth struct{}

func (m *mockUserServiceForAuth) GetOrCreateUserByGithubId(ctx context.Context, githubId int64, name, username, profilePicURL string) (*domain.User, error) {
	return &domain.User{ID: uuid.New(), Name: name, Username: username}, nil
}

func (m *mockUserServiceForAuth) GetOrCreateUserByGoogleId(ctx context.Context, googleId string, name, username, profilePicURL string) (*domain.User, error) {
	return &domain.User{ID: uuid.New(), Name: name, Username: username}, nil
}

func (m *mockUserServiceForAuth) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return &domain.User{ID: id, Name: "Test", Username: "test"}, nil
}

func (m *mockUserServiceForAuth) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}

func (m *mockUserServiceForAuth) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserServiceForAuth) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserServiceForAuth) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return nil
}
