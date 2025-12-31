package reviews

import (
	"context"
	"errors"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type mockReviewStore struct{}

func (m *mockReviewStore) GetReview(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error) {
	return &domain.Review{ID: reviewId}, nil
}

func (m *mockReviewStore) GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
	return []domain.Review{{ID: uuid.New(), UserId: userId}}, nil
}

func (m *mockReviewStore) CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return &review, nil
}

func (m *mockReviewStore) UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return &review, nil
}

func (m *mockReviewStore) DeleteReview(ctx context.Context, reviewId uuid.UUID) error {
	return nil
}

func TestService_GetAllReviewsByUserId(t *testing.T) {
	service := NewService(&mockReviewStore{})
	result, err := service.GetAllReviewsByUserId(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) == 0 {
		t.Error("expected at least one review")
	}
}

func TestService_CreateReview(t *testing.T) {
	service := NewService(&mockReviewStore{})
	review := domain.Review{ID: uuid.New(), Content: "Test"}
	result, err := service.CreateReview(context.Background(), review)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != review.ID {
		t.Errorf("expected ID %v, got %v", review.ID, result.ID)
	}
}

func TestService_UpdateReview(t *testing.T) {
	service := NewService(&mockReviewStore{})
	review := domain.Review{ID: uuid.New(), Content: "Updated"}
	result, err := service.UpdateReview(context.Background(), review)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Content != review.Content {
		t.Errorf("expected content %s, got %s", review.Content, result.Content)
	}
}

func TestService_DeleteReview(t *testing.T) {
	service := NewService(&mockReviewStore{})
	err := service.DeleteReview(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

type errorStore struct{}

func (e *errorStore) GetReview(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error) {
	return nil, errors.New("database error")
}

func (e *errorStore) GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
	return nil, errors.New("database error")
}

func (e *errorStore) CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return nil, errors.New("database error")
}

func (e *errorStore) UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	return nil, errors.New("database error")
}

func (e *errorStore) DeleteReview(ctx context.Context, reviewId uuid.UUID) error {
	return errors.New("database error")
}

func TestService_Errors(t *testing.T) {
	service := NewService(&errorStore{})
	
	_, err := service.GetAllReviewsByUserId(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error from GetAllReviewsByUserId")
	}
	
	_, err = service.CreateReview(context.Background(), domain.Review{})
	if err == nil {
		t.Error("expected error from CreateReview")
	}
}
