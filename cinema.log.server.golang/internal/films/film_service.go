package films

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	tmdbBaseUrl = "https://api.themoviedb.org/3/"
)

type Service struct {
	FilmStore FilmStore
}

type FilmStore interface {
	GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error)
	CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error)
	GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
}

type TMDBSearchResponse struct {
	Results []FilmSearchResult `json:"results"`
}

type FilmSearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"original_title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}

func NewService(f FilmStore) *Service {
	return &Service{
		FilmStore: f,
	}
}

func (s Service) CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error) {
	// Store layer handles UPSERT - if film with same external_id exists,
	// it will update and return existing film; otherwise creates new
	return s.FilmStore.CreateFilm(ctx, film)
}

func (s Service) GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
	return s.FilmStore.GetFilmById(ctx, id)
}

func (s Service) GetFilmsFromExternal(ctx context.Context, query string) ([]domain.Film, error) {
	if query == "" {
		return nil, errors.New("cannot obtain films with empty query string")
	}

	key := os.Getenv("TMDB_API_KEY")
	reqUrl := fmt.Sprintf("%ssearch/movie?query=%s&include_adult=false&language=en-US&page=1&api_key=%s", tmdbBaseUrl, query, key)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, ErrServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not process response from tmdb")
	}

	var tmdbResponse TMDBSearchResponse
	if err := json.Unmarshal(body, &tmdbResponse); err != nil {
		return nil, errors.New("could not parse tmdb response")
	}

	films := make([]domain.Film, 0, len(tmdbResponse.Results))

	for _, filmResult := range tmdbResponse.Results {
		films = append(films, domain.Film{
			ID:          uuid.New(),
			ExternalID:  filmResult.ID,
			Title:       filmResult.Title,
			Description: filmResult.Overview,
			PosterUrl:   filmResult.PosterPath,
			ReleaseYear: filmResult.ReleaseDate,
		})
	}
	return films, nil
}

func (s Service) GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error) {
	return s.FilmStore.GetFilmsForRating(ctx, userId, filmId)
}
