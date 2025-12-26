package comparisons

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

type key int

const (
	keyUser key = iota
)

type Handler struct {
	ComparisonService ComparisonService
	RatingService     RatingService
}

type ComparisonService interface {
	CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
}

type RatingService interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair, winnerId uuid.UUID) (*domain.ComparisonPair, error)
}

func NewHandler(comparisonService ComparisonService, ratingService RatingService) *Handler {
	return &Handler{
		ComparisonService: comparisonService,
		RatingService:     ratingService,
	}
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
	hasBeenCompared, err := h.ComparisonService.HasBeenCompared(r.Context(), req.UserId, req.FilmAId, req.FilmBId)
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

	updatedPair, err := h.RatingService.UpdateRatings(r.Context(), pair, req.WinningFilmId)
	if err != nil {
		http.Error(w, "Failed to update ratings", http.StatusInternalServerError)
		return
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

	_, err = h.ComparisonService.CreateComparison(r.Context(), comparison)
	if err != nil {
		http.Error(w, "Failed to create comparison history", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedPair)
}
