package films

import (
	"context"
	"database/sql"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) FilmStore {
	return &store{
		db: db,
	}
}

func (s *store) CreateFilm(ctx context.Context, film domain.Film) (*domain.Film, error) {
	// Generate a new UUID
	if film.ID == uuid.Nil {
		film.ID = uuid.New()
	}

	query := `
		INSERT INTO films (film_id, external_id, title, description, poster_url, release_year) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.ExecContext(ctx, query, film.ID, film.ExternalID, film.Title, film.Description, film.PosterUrl, film.ReleaseYear)

	if err != nil {
		return nil, err
	}

	return &film, nil
}

func (s *store) GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error) {
	query := `SELECT film_id, external_id, title, description, poster_url, release_year 
	          FROM films WHERE external_id = $1`

	film := &domain.Film{}
	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&film.ID, &film.ExternalID, &film.Title, &film.Description, &film.PosterUrl, &film.ReleaseYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFilmNotFound
		}
		return nil, err
	}

	return film, nil
}

func (s *store) GetFilmById(ctx context.Context, id uuid.UUID) (*domain.Film, error) {
	query := `SELECT film_id, external_id, title, description, poster_url, release_year 
	          FROM films WHERE film_id = $1`

	film := &domain.Film{}
	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&film.ID, &film.ExternalID, &film.Title, &film.Description, &film.PosterUrl, &film.ReleaseYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFilmNotFound
		}
		return nil, err
	}

	return film, nil
}
