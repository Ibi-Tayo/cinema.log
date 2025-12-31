package reviews

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/middleware"
	"github.com/google/uuid"
)

type mockReviewService struct {
	getReview 				  func(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error)
	getAllReviewsByUserIdFunc func(ctx context.Context, userId uuid.UUID) ([]domain.Review, error)
	createReviewFunc          func(ctx context.Context, review domain.Review) (*domain.Review, error)
	updateReviewFunc          func(ctx context.Context, review domain.Review) (*domain.Review, error)
	deleteReviewFunc          func(ctx context.Context, reviewId uuid.UUID) error
}

func (m *mockReviewService) GetReview(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error) {
	if m.getReview != nil {
		return m.getReview(ctx, reviewId)
	}
	return &domain.Review{ID: reviewId}, nil
}

func (m *mockReviewService) GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
	if m.getAllReviewsByUserIdFunc != nil {
		return m.getAllReviewsByUserIdFunc(ctx, userId)
	}
	return []domain.Review{{ID: uuid.New(), UserId: userId}}, nil
}

func (m *mockReviewService) CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	if m.createReviewFunc != nil {
		return m.createReviewFunc(ctx, review)
	}
	return &review, nil
}

func (m *mockReviewService) UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	if m.updateReviewFunc != nil {
		return m.updateReviewFunc(ctx, review)
	}
	return &review, nil
}

func (m *mockReviewService) DeleteReview(ctx context.Context, reviewId uuid.UUID) error {
	if m.deleteReviewFunc != nil {
		return m.deleteReviewFunc(ctx, reviewId)
	}
	return nil
}

type mockRatingService struct {
	getRatingFunc    func(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error)
	createRatingFunc func(ctx context.Context, userId uuid.UUID, filmId uuid.UUID, initialRating float32) (*domain.UserFilmRating, error)
}

func (m *mockRatingService) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	if m.getRatingFunc != nil {
		return m.getRatingFunc(ctx, userId, filmId)
	}
	return &domain.UserFilmRating{ID: uuid.New(), UserId: userId, FilmId: filmId}, nil
}

func (m *mockRatingService) CreateRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID, initialRating float32) (*domain.UserFilmRating, error) {
	if m.createRatingFunc != nil {
		return m.createRatingFunc(ctx, userId, filmId, initialRating)
	}
	return &domain.UserFilmRating{ID: uuid.New(), UserId: userId, FilmId: filmId, InitialRating: initialRating}, nil
}

func TestNewHandler(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.ReviewService != mockReviewSvc {
		t.Error("expected handler to contain the provided service")
	}
}

func TestHandler_GetAllReviews_Success(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	userId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/reviews/"+userId.String(), nil)
	req.SetPathValue("userId", userId.String())
	w := httptest.NewRecorder()

	handler.GetAllReviews(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_GetAllReviews_InvalidUserId(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	req := httptest.NewRequest(http.MethodGet, "/reviews/invalid", nil)
	req.SetPathValue("userId", "invalid")
	w := httptest.NewRecorder()

	handler.GetAllReviews(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetAllReviews_ServiceError(t *testing.T) {
	mockReviewSvc := &mockReviewService{
		getAllReviewsByUserIdFunc: func(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
			return nil, errors.New("database error")
		},
	}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	userId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/reviews/"+userId.String(), nil)
	req.SetPathValue("userId", userId.String())
	w := httptest.NewRecorder()

	handler.GetAllReviews(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_CreateReview_Success(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	userId := uuid.New()
	filmId := uuid.New()
	user := &domain.User{ID: userId, Name: "Test User", Username: "testuser"}

	reviewReq := map[string]interface{}{
		"content": "Great movie!",
		"rating":  4.5,
		"filmId":  filmId.String(),
	}
	body, _ := json.Marshal(reviewReq)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateReview(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d, body: %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestHandler_CreateReview_Unauthorized(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	reviewReq := map[string]interface{}{
		"content": "Great movie!",
		"rating":  4.5,
		"filmId":  uuid.New().String(),
	}
	body, _ := json.Marshal(reviewReq)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateReview(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_CreateReview_InvalidJSON(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	userId := uuid.New()
	user := &domain.User{ID: userId, Name: "Test User", Username: "testuser"}

	req := httptest.NewRequest(http.MethodPost, "/reviews", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateReview(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_UpdateReview_Success(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	userId := uuid.New()
	reviewId := uuid.New()
	user := &domain.User{ID: userId, Name: "Test User", Username: "testuser"}

	// Put a dummy review in the mock service to be updated
	mockReviewSvc.getReview = func(ctx context.Context, reviewId uuid.UUID) (*domain.Review, error) {
		return &domain.Review{ID: reviewId, UserId: userId, Content: "Old review"}, nil
	}

	updateReq := map[string]interface{}{
		"content": "Updated review!",
		"reviewId": reviewId.String(),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/reviews/"+ reviewId.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", reviewId.String())
	ctx := context.WithValue(req.Context(), middleware.KeyUser, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.UpdateReview(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestHandler_DeleteReview_Success(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	reviewId := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/reviews?id="+reviewId.String(), nil)
	w := httptest.NewRecorder()

	handler.DeleteReview(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestHandler_DeleteReview_MissingReviewId(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	req := httptest.NewRequest(http.MethodDelete, "/reviews", nil)
	w := httptest.NewRecorder()

	handler.DeleteReview(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_DeleteReview_InvalidReviewId(t *testing.T) {
	mockReviewSvc := &mockReviewService{}
	mockRatingSvc := &mockRatingService{}
	handler := NewHandler(mockReviewSvc, mockRatingSvc)

	req := httptest.NewRequest(http.MethodDelete, "/reviews?id=invalid", nil)
	w := httptest.NewRecorder()

	handler.DeleteReview(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
