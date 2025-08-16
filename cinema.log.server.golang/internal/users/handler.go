package users

import (
	"context"
	"encoding/json"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

type Handler struct {
	service UserService
}

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetOrCreateUserByGithubId(ctx context.Context, githubId uint64, name string,
	   						username string, avatarUrl string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

func NewHandler(s UserService) *Handler {
	return &Handler{
		service: s,
	}
}

// GetUserById handles GET requests for a specific user by ID
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from URL path parameter (you'll need to set this up in your router)
	userIDStr := r.PathValue("id") // Assuming Go 1.22+ with new ServeMux
	if userIDStr == "" {
		http.Error(w, ErrNoId.Error(), http.StatusBadRequest)
		return
	}

	// Parse UUID
	userID, err := utils.ParseUUID(userIDStr)
	if err != nil {
		http.Error(w, ErrInvalidId.Error(), http.StatusBadRequest)
		return
	}

	// Get user from service
	user, err := h.service.GetUserById(r.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	// Return user as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllUsers handles GET requests for all users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateUser handles POST requests to create a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, ErrInvalidJson.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), &user)
	if err != nil {
		if err == ErrUserExists {
			http.Error(w, ErrUserExists.Error(), http.StatusConflict)
			return
		}
		if err == ErrUserNameInvalidLength {
			http.Error(w, ErrUserNameInvalidLength.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteUser handles DELETE requests for a specific user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		http.Error(w, ErrNoId.Error(), http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseUUID(userIDStr)
	if err != nil {
		http.Error(w, ErrInvalidId.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		if err == ErrUserNotFound {
			http.Error(w, ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}