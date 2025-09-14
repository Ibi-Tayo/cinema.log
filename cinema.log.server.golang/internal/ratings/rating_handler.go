package ratings

import (
	"context"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type Handler struct {
	RatingService RatingService
}

type RatingService interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error)
}

func NewHandler(ratingService RatingService) *Handler {
	return &Handler{
		RatingService: ratingService,
	}
}

func (h *Handler) GetRating(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CompareFilms(w http.ResponseWriter, r *http.Request) {
	// TODO: json payload in request should be a list of comparison history objects
	// 1. Use film id + user id in the object to get user_film_rating for both film a and b (use rating service)
	// 2. Update ratings based on result (again using rating service)
	// 3. Send new user_film_rating back for both films
}
