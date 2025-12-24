package ratings

import (
	"context"
	"database/sql"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrRatingNotFound = errors.New("rating not found")
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) RatingStore {
	return &store{
		db: db,
	}
}

func (s *store) GetRating(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (*domain.UserFilmRating, error) {
	query := `
		SELECT user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
		FROM user_film_ratings
		WHERE user_id = $1 AND film_id = $2
	`

	rating := &domain.UserFilmRating{}
	row := s.db.QueryRowContext(ctx, query, userId, filmId)

	err := row.Scan(
		&rating.ID,
		&rating.UserId,
		&rating.FilmId,
		&rating.EloRating,
		&rating.NumberOfComparisons,
		&rating.LastUpdated,
		&rating.InitialRating,
		&rating.KConstantValue,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRatingNotFound
		}
		return nil, err
	}

	return rating, nil
}

func (s *store) GetAllRatings(ctx context.Context) ([]domain.UserFilmRating, error) {
	query := `
		SELECT user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
		FROM user_film_ratings
		ORDER BY last_updated DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []domain.UserFilmRating
	for rows.Next() {
		var rating domain.UserFilmRating
		err := rows.Scan(
			&rating.ID,
			&rating.UserId,
			&rating.FilmId,
			&rating.EloRating,
			&rating.NumberOfComparisons,
			&rating.LastUpdated,
			&rating.InitialRating,
			&rating.KConstantValue,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ratings, nil
}

func (s *store) GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRating, error) {
	query := `
		SELECT user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
		FROM user_film_ratings
		WHERE user_id = $1
		ORDER BY number_of_comparisons ASC, last_updated ASC
	`

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []domain.UserFilmRating
	for rows.Next() {
		var rating domain.UserFilmRating
		err := rows.Scan(
			&rating.ID,
			&rating.UserId,
			&rating.FilmId,
			&rating.EloRating,
			&rating.NumberOfComparisons,
			&rating.LastUpdated,
			&rating.InitialRating,
			&rating.KConstantValue,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ratings, nil
}

func (s *store) CreateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
	query := `
		INSERT INTO user_film_ratings (user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
	`

	createdRating := &domain.UserFilmRating{}
	row := s.db.QueryRowContext(ctx, query,
		rating.ID,
		rating.UserId,
		rating.FilmId,
		rating.EloRating,
		rating.NumberOfComparisons,
		rating.LastUpdated,
		rating.InitialRating,
		rating.KConstantValue,
	)

	err := row.Scan(
		&createdRating.ID,
		&createdRating.UserId,
		&createdRating.FilmId,
		&createdRating.EloRating,
		&createdRating.NumberOfComparisons,
		&createdRating.LastUpdated,
		&createdRating.InitialRating,
		&createdRating.KConstantValue,
	)

	if err != nil {
		return nil, err
	}

	return createdRating, nil
}

func (s *store) UpdateRating(ctx context.Context, rating domain.UserFilmRating) (*domain.UserFilmRating, error) {
	query := `
		UPDATE user_film_ratings
		SET elo_rating = $1,
		    number_of_comparisons = $2,
		    last_updated = $3,
		    k_constant_value = $4
		WHERE user_film_rating_id = $5
		RETURNING user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
	`

	updatedRating := &domain.UserFilmRating{}
	row := s.db.QueryRowContext(ctx, query,
		rating.EloRating,
		rating.NumberOfComparisons,
		rating.LastUpdated,
		rating.KConstantValue,
		rating.ID,
	)

	err := row.Scan(
		&updatedRating.ID,
		&updatedRating.UserId,
		&updatedRating.FilmId,
		&updatedRating.EloRating,
		&updatedRating.NumberOfComparisons,
		&updatedRating.LastUpdated,
		&updatedRating.InitialRating,
		&updatedRating.KConstantValue,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRatingNotFound
		}
		return nil, err
	}

	return updatedRating, nil
}

func (s *store) UpdateRatings(ctx context.Context, ratings domain.ComparisonPair) (*domain.ComparisonPair, error) {
	// Begin a transaction to ensure both updates succeed or fail together
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		UPDATE user_film_ratings
		SET elo_rating = $1,
		    number_of_comparisons = $2,
		    last_updated = $3,
		    k_constant_value = $4
		WHERE user_film_rating_id = $5
		RETURNING user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
	`

	// Update Film A
	updatedFilmA := &domain.UserFilmRating{}
	row := tx.QueryRowContext(ctx, query,
		ratings.FilmA.EloRating,
		ratings.FilmA.NumberOfComparisons,
		ratings.FilmA.LastUpdated,
		ratings.FilmA.KConstantValue,
		ratings.FilmA.ID,
	)

	err = row.Scan(
		&updatedFilmA.ID,
		&updatedFilmA.UserId,
		&updatedFilmA.FilmId,
		&updatedFilmA.EloRating,
		&updatedFilmA.NumberOfComparisons,
		&updatedFilmA.LastUpdated,
		&updatedFilmA.InitialRating,
		&updatedFilmA.KConstantValue,
	)

	if err != nil {
		return nil, err
	}

	// Update Film B
	updatedFilmB := &domain.UserFilmRating{}
	row = tx.QueryRowContext(ctx, query,
		ratings.FilmB.EloRating,
		ratings.FilmB.NumberOfComparisons,
		ratings.FilmB.LastUpdated,
		ratings.FilmB.KConstantValue,
		ratings.FilmB.ID,
	)

	err = row.Scan(
		&updatedFilmB.ID,
		&updatedFilmB.UserId,
		&updatedFilmB.FilmId,
		&updatedFilmB.EloRating,
		&updatedFilmB.NumberOfComparisons,
		&updatedFilmB.LastUpdated,
		&updatedFilmB.InitialRating,
		&updatedFilmB.KConstantValue,
	)

	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.ComparisonPair{
		FilmA: *updatedFilmA,
		FilmB: *updatedFilmB,
	}, nil
}
