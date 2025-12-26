package reviews

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

type Handler struct {
	ReviewService ReviewService
	RatingService RatingService
}

type ReviewService interface {
	GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error)
	CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error)
	UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error)
	DeleteReview(ctx context.Context, reviewId uuid.UUID) error
}

type RatingService interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	CreateRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID, initialRating float32) (*domain.UserFilmRating, error)
}

func NewHandler(reviewService ReviewService, ratingService RatingService) *Handler {
	return &Handler{
		ReviewService: reviewService,
		RatingService: ratingService,
	}
}

func (h *Handler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("userId")
	userId, err := utils.ParseUUID(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.ReviewService.GetAllReviewsByUserId(r.Context(), userId)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, reviews)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middleware.KeyUser).(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Content string    `json:"content"`
		Rating  float32   `json:"rating"`
		FilmId  uuid.UUID `json:"filmId"`
	}

	if err := utils.DecodeJSON(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review := domain.Review{
		ID:      uuid.New(),
		Content: req.Content,
		Date:    time.Now(),
		Rating:  req.Rating,
		FilmId:  req.FilmId,
		UserId:  user.ID,
	}

	createdReview, err := h.ReviewService.CreateReview(r.Context(), review)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	// Create initial ELO rating for the film if it doesn't exist
	_, err = h.RatingService.GetRating(r.Context(), user.ID, req.FilmId)
	if err != nil {
		// Rating doesn't exist, create it
		_, err = h.RatingService.CreateRating(r.Context(), user.ID, req.FilmId, req.Rating)
		if err != nil {
			// Log the error but don't fail the review creation
			// The rating can be created later
			fmt.Printf("Failed to create initial rating: %v\n", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	utils.SendJSON(w, createdReview)
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middleware.KeyUser).(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reviewIdStr := r.PathValue("id")
	reviewId, err := utils.ParseUUID(reviewIdStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string  `json:"content"`
		Rating  float32 `json:"rating"`
	}

	if err := utils.DecodeJSON(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review := domain.Review{
		ID:      reviewId,
		Content: req.Content,
		Date:    time.Now(),
		Rating:  req.Rating,
		UserId:  user.ID,
	}

	updatedReview, err := h.ReviewService.UpdateReview(r.Context(), review)
	if err != nil {
		if err == ErrReviewNotFound {
			http.Error(w, ErrReviewNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedReview)
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	// Get review ID from query parameter
	reviewIdStr := r.URL.Query().Get("id")
	if reviewIdStr == "" {
		http.Error(w, "Missing review ID", http.StatusBadRequest)
		return
	}

	reviewId, err := utils.ParseUUID(reviewIdStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	err = h.ReviewService.DeleteReview(r.Context(), reviewId)
	if err != nil {
		if err == ErrReviewNotFound {
			http.Error(w, ErrReviewNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
