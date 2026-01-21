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
	"slices"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	tmdbBaseUrl = "https://api.themoviedb.org/3/"

	ErrEmptyQueryString    = errors.New("cannot obtain films with empty query string")
	ErrProcessTMDBResponse = errors.New("could not process response from tmdb")
	ErrParseTMDBResponse   = errors.New("could not parse tmdb response")
	ErrEmptyFilmList       = errors.New("cannot generate recommendations with empty film list")
	ErrTooManyFilms        = errors.New("cannot generate recommendations with more than 10 films")
)

type Service struct {
	FilmStore              FilmStore
	GraphService           GraphService
	tmdbRecommendationFunc func(film domain.Film) []domain.Film
}

type GraphService interface {
	AddFilmToGraph(ctx context.Context, userID uuid.UUID, film domain.Film, recommendations []domain.Film) error
}

type FilmStore interface {
	GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error)
	GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error)
	CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error)
	GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error)
	GetFilmRecommendation(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error)
	CreateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error)
	UpdateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error)
	GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error)
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

func NewService(f FilmStore, g GraphService) *Service {
	return &Service{
		FilmStore:              f,
		GraphService:           g,
		tmdbRecommendationFunc: getFilmRecommendationsFromTmdb,
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
		return nil, ErrEmptyQueryString
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
		return nil, ErrProcessTMDBResponse
	}

	var tmdbResponse TMDBSearchResponse
	if err := json.Unmarshal(body, &tmdbResponse); err != nil {
		return nil, ErrParseTMDBResponse
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

// Generates film recommendations using TMDB, assumption when using this is that films in the argument have been seen by the user
func (s Service) GenerateFilmRecommendations(ctx context.Context, userId uuid.UUID, films []domain.Film) ([]domain.Film, error) {
	if len(films) == 0 {
		return []domain.Film{}, ErrEmptyFilmList
	}

	if len(films) > 10 {
		return nil, ErrTooManyFilms
	}

	allRecommendations := make([]domain.Film, 0)

	for _, film := range films {
		// ensure film exists in films table
		_, err := s.FilmStore.CreateFilm(ctx, &film)
		if err != nil {
			return nil, err
		}
		// add/update the film_recommendation table: has_seen = true, recommendations_generated = true, all other stuff too
		existingRec, err := s.FilmStore.GetFilmRecommendation(ctx, userId, film.ExternalID)
		if err != nil {
			if err == ErrFilmRecommendationNotFound {
				// create new recommendation entry
				newRec := &domain.FilmRecommendation{
					ID:                       uuid.New(),
					UserID:                   userId,
					ExternalFilmID:           film.ExternalID,
					HasSeen:                  true,
					HasBeenRecommended:       false,
					RecommendationsGenerated: true,
				}
				_, err := s.FilmStore.CreateFilmRecommendation(ctx, newRec)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// update existing recommendation entry
			existingRec.HasSeen = true
			existingRec.RecommendationsGenerated = true
			_, err := s.FilmStore.UpdateFilmRecommendation(ctx, existingRec)
			if err != nil {
				return nil, err
			}
		}

		recommendations := s.tmdbRecommendationFunc(film)

		// Add film to user's graph with its recommendations
		if err := s.GraphService.AddFilmToGraph(ctx, userId, film, recommendations); err != nil {
			log.Printf("Failed to add film to graph: %v", err)
			// Don't fail the entire operation, just log and continue
		}

		allRecommendations = slices.Concat(allRecommendations, recommendations)
	}

	// Deduplicate recommendations by ExternalID
	// Multiple seed films can recommend the same film
	uniqueRecommendations := make(map[int]domain.Film)
	for _, recFilm := range allRecommendations {
		uniqueRecommendations[recFilm.ExternalID] = recFilm
	}

	allRecommendations = make([]domain.Film, 0, len(uniqueRecommendations))
	for _, film := range uniqueRecommendations {
		allRecommendations = append(allRecommendations, film)
	}

	// to prevent circular recommendations, we filter the all recommendations list by checking the film_recommendation_table
	// is there an entry? omit films where - has_seen = true (this means that recommended films could be re-recommended if they havent been seen, i'll have to see how circular this could get)
	// take allRecommendations and add/update the film_recommendation table: has_been_recommended = true

	filteredRecommendations := make([]domain.Film, 0)
	for _, recFilm := range allRecommendations {
		// ensure film exists in films table
		_, err := s.FilmStore.CreateFilm(ctx, &recFilm)
		if err != nil {
			return nil, err
		}
		// check film_recommendation table
		existingRec, err := s.FilmStore.GetFilmRecommendation(ctx, userId, recFilm.ExternalID)
		if err != nil {
			if err == ErrFilmRecommendationNotFound {
				// no existing recommendation, safe to add
				filteredRecommendations = append(filteredRecommendations, recFilm)
				_, err := s.FilmStore.CreateFilmRecommendation(ctx, &domain.FilmRecommendation{
					ID:                       uuid.New(),
					UserID:                   userId,
					ExternalFilmID:           recFilm.ExternalID,
					HasSeen:                  false,
					HasBeenRecommended:       true,
					RecommendationsGenerated: false,
				})
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// existing recommendation found, only add if has_seen is false
			if !existingRec.HasSeen {
				filteredRecommendations = append(filteredRecommendations, recFilm)
				existingRec.HasBeenRecommended = true
				_, err := s.FilmStore.UpdateFilmRecommendation(ctx, existingRec)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	allRecommendations = filteredRecommendations

	return allRecommendations, nil
}

// Gets seen but unrated films (should prompt user to rate these films)
func (s Service) GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error) {
	return s.FilmStore.GetSeenUnratedFilms(ctx, userId)
}

func getFilmRecommendationsFromTmdb(film domain.Film) []domain.Film {
	key := os.Getenv("TMDB_API_KEY")
	reqUrl := fmt.Sprintf("%smovie/%d/recommendations?api_key=%s", tmdbBaseUrl, film.ExternalID, key)

	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println("Error fetching recommendations from TMDB:", err)
		return []domain.Film{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("TMDB API returned status %d for film ID %d\n", resp.StatusCode, film.ExternalID)
		return []domain.Film{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading TMDB response body:", err)
		return []domain.Film{}
	}

	var tmdbResponse TMDBSearchResponse
	if err := json.Unmarshal(body, &tmdbResponse); err != nil {
		log.Println("Error parsing TMDB response:", err)
		return []domain.Film{}
	}

	// Limit to top 10 recommendations to avoid weak suggestions
	topRecommendations := tmdbResponse.Results
	if len(topRecommendations) > 10 {
		topRecommendations = topRecommendations[:10]
	}

	recommendedFilms := make([]domain.Film, 0, len(topRecommendations))

	for _, filmResult := range topRecommendations {
		recommendedFilms = append(recommendedFilms, domain.Film{
			ID:          uuid.New(),
			ExternalID:  filmResult.ID,
			Title:       filmResult.Title,
			Description: filmResult.Overview,
			PosterUrl:   filmResult.PosterPath,
			ReleaseYear: filmResult.ReleaseDate,
		})
	}

	return recommendedFilms
}
