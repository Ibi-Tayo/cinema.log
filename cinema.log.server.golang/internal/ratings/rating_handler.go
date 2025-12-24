package ratings

import (
	"context"
	"encoding/json"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

type Handler struct {
	RatingService RatingService
}

type RatingService interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	GetRatingsForComparison(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, winnerId uuid.UUID) (*domain.ComparisonPair, error)
}

func NewHandler(ratingService RatingService) *Handler {
	return &Handler{
		RatingService: ratingService,
	}
}

func (h *Handler) GetRating(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	filmIDStr := r.URL.Query().Get("filmId")
	if filmIDStr == "" {
		http.Error(w, "filmId is required", http.StatusBadRequest)
		return
	}

	filmID, err := uuid.Parse(filmIDStr)
	if err != nil {
		http.Error(w, "invalid filmId", http.StatusBadRequest)
		return
	}

	rating, err := h.RatingService.GetRating(r.Context(), userID, filmID)
	if err != nil {
		http.Error(w, "failed to get rating", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, rating)
}

func (h *Handler) GetRatingsForComparison(w http.ResponseWriter, r *http.Request) {
	// TODO: Get userId from request context/auth
	// For now, expecting it as a query parameter
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	ratings, err := h.RatingService.GetRatingsForComparison(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get ratings", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, ratings)
}

func (h *Handler) CompareFilms(w http.ResponseWriter, r *http.Request) {
	var comparison domain.ComparisonHistory
	if err := json.NewDecoder(r.Body).Decode(&comparison); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	filmARating, err := h.RatingService.GetRating(r.Context(), comparison.UserId, comparison.FilmAId)
	if err != nil {
		http.Error(w, "failed to get rating for film A", http.StatusInternalServerError)
		return
	}

	filmBRating, err := h.RatingService.GetRating(r.Context(), comparison.UserId, comparison.FilmBId)
	if err != nil {
		http.Error(w, "failed to get rating for film B", http.StatusInternalServerError)
		return
	}

	pair := domain.ComparisonPair{
		FilmA: *filmARating,
		FilmB: *filmBRating,
	}

	updatedPair, err := h.RatingService.UpdateRatings(r.Context(), pair, comparison.WinningFilmId)
	if err != nil {
		http.Error(w, "failed to update ratings", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedPair)
}
