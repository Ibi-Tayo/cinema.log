package films

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	GetFilmById(ctx context.Context, id uuid.UUID) (domain.Film, error)
	GetFilmByExternalId(ctx context.Context, id int) (domain.Film, error)
	CreateFilm(ctx context.Context, film domain.Film) (domain.Film, error)
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

func (s Service) CreateFilm(ctx context.Context, film domain.Film) (domain.Film, error) {
	return s.FilmStore.CreateFilm(ctx, film)
}

func (s Service) GetFilmById(ctx context.Context, id uuid.UUID) (domain.Film, error) {
	film, err := s.FilmStore.GetFilmById(ctx, id)
	if err != nil {
		return domain.Film{}, err
	}
	return film, nil
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
		// check if film already exists in db
		film, err := s.FilmStore.GetFilmByExternalId(ctx, filmResult.ID)
		if (err != nil) {
			// create new film and put in database
			createdFilm, err := s.FilmStore.CreateFilm(ctx, domain.Film{
				ID: uuid.New(),
				ExternalID: filmResult.ID,
				Title: filmResult.Title,
				Description: filmResult.Overview,
				PosterUrl: filmResult.PosterPath,
				ReleaseYear: filmResult.ReleaseDate,
			})
			if (err != nil) {
				log.Printf("could not add film: %s to database:", filmResult.Title)
				continue
			}
			film = createdFilm
		} 
		// if film exists then add to list
		films = append(films, film)
	}

	return films, nil
}
