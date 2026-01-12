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

func (s *store) CreateFilm(ctx context.Context, film *domain.Film) (*domain.Film, error) {
	// Generate a new UUID if not provided
	if film.ID == uuid.Nil {
		film.ID = uuid.New()
	}

	// Use UPSERT to prevent race conditions
	// If film with same external_id exists, return it; otherwise create new
	query := /* sql */ `
		INSERT INTO films (film_id, external_id, title, description, poster_url, release_year) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (external_id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			poster_url = EXCLUDED.poster_url,
			release_year = EXCLUDED.release_year
		RETURNING film_id, external_id, title, description, poster_url, release_year`

	var createdFilm domain.Film
	err := s.db.QueryRowContext(ctx, query, film.ID, film.ExternalID, film.Title, film.Description, film.PosterUrl, film.ReleaseYear).Scan(
		&createdFilm.ID,
		&createdFilm.ExternalID,
		&createdFilm.Title,
		&createdFilm.Description,
		&createdFilm.PosterUrl,
		&createdFilm.ReleaseYear,
	)

	if err != nil {
		return nil, err
	}

	return &createdFilm, nil
}

func (s *store) GetFilmByExternalId(ctx context.Context, id int) (*domain.Film, error) {
	query := /* sql */ `SELECT film_id, external_id, title, description, poster_url, release_year 
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
	query := /* sql */ `SELECT film_id, external_id, title, description, poster_url, release_year 
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
	// Select films rated by the user, ordered by closeness of ELO rating to the specified film (maintaining competitive balance)
	// and then by the number of comparisons (ascending)
	query := /* sql */ `
		SELECT
			f.film_id,
			f.external_id,
			f.title,
			f.description,
			f.poster_url,
			f.release_year
		FROM films f
		INNER JOIN user_film_ratings ufr
			ON f.film_id = ufr.film_id
		WHERE
			ufr.user_id = $1
			AND f.film_id != $2
		ORDER BY ABS(
			ufr.elo_rating - (
				SELECT elo_rating
				FROM user_film_ratings
				WHERE user_id = $1
				AND film_id = $2
			)
		),
		ufr.number_of_comparisons ASC;
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

func (s *store) CreateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
	query := /* sql */ `
		INSERT INTO film_recommendation (film_recommendation_id, user_id, external_film_id, has_seen, has_been_recommended, recommendations_generated) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		recommendation.ID,
		recommendation.UserID,
		recommendation.ExternalFilmID,
		recommendation.HasSeen,
		recommendation.HasBeenRecommended,
		recommendation.RecommendationsGenerated,
	)

	if err != nil {
		return nil, err
	}

	return recommendation, nil
}

func (s *store) GetFilmRecommendation(ctx context.Context, userId uuid.UUID, externalFilmId int) (*domain.FilmRecommendation, error) {
	query := /* sql */ `
		SELECT film_recommendation_id, user_id, external_film_id, has_seen, has_been_recommended, recommendations_generated
		FROM film_recommendation
		WHERE user_id = $1 AND external_film_id = $2
	`

	recommendation := &domain.FilmRecommendation{}
	row := s.db.QueryRowContext(ctx, query, userId, externalFilmId)

	err := row.Scan(&recommendation.ID, &recommendation.UserID, &recommendation.ExternalFilmID, &recommendation.HasSeen, &recommendation.HasBeenRecommended, &recommendation.RecommendationsGenerated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFilmRecommendationNotFound
		}
		return nil, err
	}

	return recommendation, nil
}

func (s *store) UpdateFilmRecommendation(ctx context.Context, recommendation *domain.FilmRecommendation) (*domain.FilmRecommendation, error) {
	query := /* sql */ `
		UPDATE film_recommendation
		SET has_seen = $1, has_been_recommended = $2, recommendations_generated = $3
		WHERE film_recommendation_id = $4
	`

	_, err := s.db.ExecContext(ctx, query,
		recommendation.HasSeen,
		recommendation.HasBeenRecommended,
		recommendation.RecommendationsGenerated,
		recommendation.ID,
	)

	if err != nil {
		return nil, err
	}

	return recommendation, nil
}

// return list of films that have been seen (film_recommendation table) AND have not been rated (user_id and film_id on user_film_ratings)
// might need to link up via external_id -> film table -> film_id -> user_film_ratings
func (s *store) GetSeenUnratedFilms(ctx context.Context, userId uuid.UUID) ([]domain.Film, error) {
	query := /* sql */ `
		SELECT
			f.film_id,
			f.external_id,
			f.title,
			f.description,
			f.poster_url,
			f.release_year
		FROM films f
		INNER JOIN film_recommendation fr
			ON f.external_id = fr.external_film_id
		LEFT JOIN user_film_ratings ufr
			ON f.film_id = ufr.film_id AND ufr.user_id = $1
		WHERE
			fr.user_id = $1
			AND fr.has_seen = TRUE
			AND ufr.film_id IS NULL
	`

	rows, err := s.db.QueryContext(ctx, query, userId)
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
