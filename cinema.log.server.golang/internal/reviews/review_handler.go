package reviews

import "net/http"

type Handler struct {
	ReviewService ReviewService
}

type ReviewService interface {
}

func NewHandler(reviewService ReviewService) *Handler {
	return &Handler{
		ReviewService: reviewService,
	}
}

func (h *Handler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
}