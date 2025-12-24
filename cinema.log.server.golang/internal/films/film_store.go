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

func (s *store) GetFilmsForRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) ([]domain.Film, error) {
	query := `
		SELECT f.film_id, f.external_id, f.title, f.description, f.poster_url, f.release_year
		FROM films f
		INNER JOIN user_film_ratings ufr ON f.film_id = ufr.film_id
		WHERE ufr.user_id = $1 AND f.film_id != $2
	`

	rows, err := s.db.QueryContext(ctx, query, userId, filmId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []domain.Film
	for rows.Next() {
		var film domain.Film
		err := rows.Scan(&film.ID, &film.ExternalID, &film.Title, &film.Description, &film.PosterUrl, &film.ReleaseYear)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return films, nil
}
