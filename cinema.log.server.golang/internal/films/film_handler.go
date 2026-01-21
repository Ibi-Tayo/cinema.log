package films

import (
	"context"
	"errors"
	"log"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrFilmNotFound               = errors.New("film not found")
	ErrFilmRecommendationNotFound = errors.New("film recommendation not found")
	ErrServer                     = errors.New("internal server error")
)

type Handler struct {
	FilmService   FilmService
	RatingService RatingService
}

type FilmService interface {
	CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error)
	GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	GetFilmsFromExternal(ctx context.Context, query string) ([]domain.Film, error) // ? pagination?
	GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
	GenerateFilmRecommendations(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error)
	GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error)
}

type RatingService interface {
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
}

func NewHandler(filmService FilmService, ratingService RatingService) *Handler {
	return &Handler{
		FilmService:   filmService,
		RatingService: ratingService,
	}
}

func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	var film domain.Film
	if err := utils.DecodeJSON(r, &film); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdFilm, err := h.FilmService.CreateFilm(r.Context(), &film)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, createdFilm)
}

func (h *Handler) GetFilmById(w http.ResponseWriter, r *http.Request) {
	reqId := r.PathValue("id")
	id, err := utils.ParseUUID(reqId)
	if err != nil {
		http.Error(w, "Error parsing UUID", http.StatusInternalServerError)
		return
	}

	film, err := h.FilmService.GetFilmById(r.Context(), id)
	if err != nil {
		if err == ErrFilmNotFound {
			http.Error(w, ErrFilmNotFound.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
			return
		}
	}
	utils.SendJSON(w, film)
}

func (h *Handler) GetFilmsFromExternal(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("f")
	if search == "" {
		http.Error(w, "Missing required query parameter 'f' for film search", http.StatusBadRequest)
		return
	}
	films, err := h.FilmService.GetFilmsFromExternal(r.Context(), search)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendJSON(w, films)
}

func (h *Handler) GetFilmsForComparison(w http.ResponseWriter, r *http.Request) {
	// Get userId and filmId from query parameters
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseUUID(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	filmIDStr := r.URL.Query().Get("filmId")
	if filmIDStr == "" {
		http.Error(w, "filmId is required", http.StatusBadRequest)
		return
	}

	filmID, err := utils.ParseUUID(filmIDStr)
	if err != nil {
		http.Error(w, "invalid filmId", http.StatusBadRequest)
		return
	}

	// Get films for rating that the user has already rated (excluding the current film)
	candidateFilms, err := h.FilmService.GetFilmsForRating(r.Context(), userID, filmID)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	// Filter out films that have already been compared
	var filmsForComparison []domain.Film
	for _, film := range candidateFilms {
		hasBeenCompared, err := h.RatingService.HasBeenCompared(r.Context(), userID, filmID, film.ID)
		if err != nil {
			log.Printf("error checking comparison history: %v", err)
			continue
		}
		if !hasBeenCompared {
			filmsForComparison = append(filmsForComparison, film)
		}
		// Limit to 10 films
		if len(filmsForComparison) >= 10 {
			break
		}
	}

	utils.SendJSON(w, filmsForComparison)
}

func (h *Handler) GenerateFilmRecommendations(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseUUID(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	var films []domain.Film
	if err := utils.DecodeJSON(r, &films); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	recommendations, err := h.FilmService.GenerateFilmRecommendations(r.Context(), userID, films)
	if err != nil {
		if errors.Is(err, ErrTooManyFilms) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, recommendations)
}

func (h *Handler) GetSeenUnratedFilms(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseUUID(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	films, err := h.FilmService.GetSeenUnratedFilms(r.Context(), userID)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, films)
}
