package ratings

import (
	"context"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	RatingStore RatingStore
}

type RatingStore interface {
	GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error)
}

func NewService(r RatingStore) *Service {
	return &Service{
		RatingStore: r,
	}
}

func (s Service) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	// TODO: get from rating store
	panic("not implemented")
}

func (s Service) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	// TODO: get from rating store
	panic("not implemented")
}

func (s Service) FilterRatingsForComparison([]domain.UserFilmRating) []domain.UserFilmRating {
	// TODO: Sort films based on 1. the ones that have had the least comparisons at the top 2. oldest comp date at the top
	// TODO: Then take up to the first 5 or 10
	panic("not implemented")
}

func (s Service) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error) {
	// TODO: update the elo ratings etc etc
	// TODO: update user_film_rating table using rating store, return updated pair of films
	panic("not implemented")
}
