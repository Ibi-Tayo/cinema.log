package films

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrFilmNotFound = errors.New("film not found")
	ErrServer       = errors.New("internal server error")
	ErrEncoding     = errors.New("error encoding response")
)

type Handler struct {
	FilmService FilmService
}

type FilmService interface {
	GetFilmById(ctx context.Context, id uuid.UUID) (domain.Film, error)
	GetFilmsFromExternal(ctx context.Context, query string) ([]domain.Film, error) // ? pagination?

	// All the film handler needs are the above methods - the review handler would need an interface that defines the method below
	// AddFilm(ctx context.Context, film domain.Film) (domain.Film, error)
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
	sendJSON(w, film)
}

func (h *Handler) GetFilmsFromExternal(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("f")
	films, err := h.FilmService.GetFilmsFromExternal(r.Context(), search)
	if err != nil {
		http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, films)
}

func sendJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
	}
}
