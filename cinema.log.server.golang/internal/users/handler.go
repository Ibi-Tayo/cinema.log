package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cinema.log.server.golang/internal/domain"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	// add other methods like Create(ctx, user), Update(ctx, user), etc.
}

type Handler struct {
	service UserService
}

func NewHandler(s UserService) *Handler {
	return &Handler{
		service: s,
	}
}

// EXAMPLE
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 1. **Decode and Validate the Request**
	// We get the user ID from the URL path. In a real app, you'd use a router
	// like chi or gin to extract this, but here we'll do it manually for clarity.
	idStr := r.PathValue("id") // Available in Go 1.22+ with http.ServeMux
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusBadRequest)
		return
	}

	// 2. **Call the Business Logic (Service Layer)**
	// We pass the request context, which can handle timeouts and cancellations.
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 3. **Encode and Send the Response**
	// Set the content type header and write the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set the status code to 200 OK

	// Encode the user object into JSON and send it.
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}