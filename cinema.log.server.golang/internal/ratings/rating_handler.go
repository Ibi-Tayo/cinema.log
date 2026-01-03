package ratings

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

type Handler struct {
	RatingService RatingService
}

type RatingService interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error)
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, comparison domain.ComparisonHistory) (*domain.ComparisonPair, error)
	CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
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
		http.Error(w, "rating not found", http.StatusNotFound)
		return
	}

	utils.SendJSON(w, rating)
}

// Fetch all ratings in order (ranked by elo rating)
func (h *Handler) GetRatingsByUserId(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	ratings, err := h.RatingService.GetRatingsByUserId(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get ratings", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, ratings)
}

type CompareFilmsRequest struct {
	UserId        uuid.UUID `json:"userId"`
	FilmAId       uuid.UUID `json:"filmAId"`
	FilmBId       uuid.UUID `json:"filmBId"`
	WinningFilmId uuid.UUID `json:"winningFilmId"`
	WasEqual      bool      `json:"wasEqual"`
}

func (h *Handler) CompareFilms(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middleware.KeyUser).(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CompareFilmsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Ensure the userId in the request matches the authenticated user
	if req.UserId != user.ID {
		http.Error(w, "Unauthorized: user ID mismatch", http.StatusUnauthorized)
		return
	}

	// Check if films have already been compared
	hasBeenCompared, err := h.RatingService.HasBeenCompared(r.Context(), req.UserId, req.FilmAId, req.FilmBId)
	if err != nil {
		http.Error(w, "Failed to check comparison history", http.StatusInternalServerError)
		return
	}
	if hasBeenCompared {
		http.Error(w, "Films have already been compared", http.StatusBadRequest)
		return
	}

	// Get ratings for both films
	filmARating, err := h.RatingService.GetRating(r.Context(), req.UserId, req.FilmAId)
	if err != nil {
		http.Error(w, "Failed to get rating for film A", http.StatusInternalServerError)
		return
	}

	filmBRating, err := h.RatingService.GetRating(r.Context(), req.UserId, req.FilmBId)
	if err != nil {
		http.Error(w, "Failed to get rating for film B", http.StatusInternalServerError)
		return
	}

	// Update ELO ratings
	pair := domain.ComparisonPair{
		FilmA: *filmARating,
		FilmB: *filmBRating,
	}

	// Create comparison history
	comparison := domain.ComparisonHistory{
		ID:             uuid.New(),
		UserId:         req.UserId,
		FilmAId:        req.FilmAId,
		FilmBId:        req.FilmBId,
		WinningFilmId:  req.WinningFilmId,
		ComparisonDate: time.Now(),
		WasEqual:       req.WasEqual,
	}

	updatedPair, err := h.RatingService.UpdateRatings(r.Context(), pair, comparison)
	if err != nil {
		http.Error(w, "Failed to update ratings", http.StatusInternalServerError)
		return
	}

	_, err = h.RatingService.CreateComparison(r.Context(), comparison)
	if err != nil {
		http.Error(w, "Failed to create comparison history", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedPair)
}
