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
	ErrFilmNotFound = errors.New("film not found")
	ErrServer       = errors.New("internal server error")
)

type Handler struct {
	FilmService   FilmService
	RatingService RatingService
}

type FilmService interface {
	GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	GetFilmsFromExternal(ctx context.Context, query string) ([]domain.Film, error) // ? pagination?
	GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
}

type RatingService interface {
	GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error)
	FilterRatingsForComparison([]domain.UserFilmRating) []domain.UserFilmRating
}

func NewHandler(filmService FilmService) *Handler {
	return &Handler{
		FilmService: filmService,
	}
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

func (h *Handler) GetFilmsForRating(w http.ResponseWriter, r *http.Request) {
	// 1. Get all ratings using the rating service
	allRatings, err := h.RatingService.GetAllRatings(r.Context())
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Use rating service to priortise the top 5 or 10 based on no of comps and last updated
	filteredRatings := h.RatingService.FilterRatingsForComparison(allRatings)

	// 3. Use those film id's to get from db and send them as response
	var films []domain.Film
	for _, rating := range filteredRatings {
		if film, err := h.FilmService.GetFilmById(r.Context(), rating.FilmId); err != nil {
			log.Printf("could not find film with id %s %v", rating.FilmId, err)
		} else {
			films = append(films, *film)
		}
	}
	utils.SendJSON(w, films)
}
