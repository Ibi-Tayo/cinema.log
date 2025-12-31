package reviews

import (
	"context"
	"errors"
	"sort"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrReviewNotFound = errors.New("review not found")
	ErrServer         = errors.New("internal server error")
)

type Service struct {
	ReviewStore ReviewStore
}

type ReviewStore interface {
	GetReview(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error)
	GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error)
	CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error)
	UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error)
	DeleteReview(ctx context.Context, reviewId uuid.UUID) error
}

func NewService(reviewStore ReviewStore) *Service {
	return &Service{
		ReviewStore: reviewStore,
	}
}

func (s *Service) GetReview(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error) {
	return s.ReviewStore.GetReview(ctx, reviewId)
}

func (s *Service) GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
	reviews, err := s.ReviewStore.GetAllReviewsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	// sort reviews by most recent first
	sort.SliceStable(reviews, func(i, j int) bool {
		return reviews[i].Date.After(reviews[j].Date)
	})
	return reviews, nil
}

func (s *Service) CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return s.ReviewStore.CreateReview(ctx, review)
}

func (s *Service) UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return s.ReviewStore.UpdateReview(ctx, review)
}

func (s *Service) DeleteReview(ctx context.Context, reviewId uuid.UUID) error {
	return s.ReviewStore.DeleteReview(ctx, reviewId)
}
