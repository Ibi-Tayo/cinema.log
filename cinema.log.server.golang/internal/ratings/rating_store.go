package ratings

import (
	"context"
	"database/sql"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrRatingNotFound     = errors.New("rating not found")
	ErrComparisonNotFound = errors.New("comparison not found")
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
	query := /* sql */ `
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
	query := /* sql */ `
		SELECT user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
		FROM user_film_ratings
		ORDER BY last_updated DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := []domain.UserFilmRating{}
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

func (s *store) GetRatingsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserFilmRatingDetail, error) {
	// Fetch ratings along with film details for the given userId
	query := /* sql */ `
		SELECT r.user_film_rating_id, r.user_id, r.film_id, r.elo_rating, r.number_of_comparisons, r.last_updated, r.initial_rating, r.k_constant_value, f.title, f.release_year, f.poster_url
		FROM user_film_ratings r
		JOIN films f ON r.film_id = f.film_id
		WHERE r.user_id = $1
		ORDER BY r.elo_rating DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := []domain.UserFilmRatingDetail{}
	for rows.Next() {
		var rating domain.UserFilmRatingDetail
		err := rows.Scan(
			&rating.Rating.ID,
			&rating.Rating.UserId,
			&rating.Rating.FilmId,
			&rating.Rating.EloRating,
			&rating.Rating.NumberOfComparisons,
			&rating.Rating.LastUpdated,
			&rating.Rating.InitialRating,
			&rating.Rating.KConstantValue,
			&rating.FilmTitle,
			&rating.FilmReleaseYear,
			&rating.FilmPosterURL,
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
	query := /* sql */ `
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
	query := /* sql */ `
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

	query := /* sql */ `
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

func (s *store) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
	if comparison.ID == uuid.Nil {
		comparison.ID = uuid.New()
	}

	query := /* sql */ `
		INSERT INTO comparison_histories (comparison_history_id, user_id, film_a_film_id, film_b_film_id, winning_film_film_id, comparison_date, was_equal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING comparison_history_id, user_id, film_a_film_id, film_b_film_id, winning_film_film_id, comparison_date, was_equal
	`

	createdComparison := &domain.ComparisonHistory{}
	row := s.db.QueryRowContext(ctx, query,
		comparison.ID,
		comparison.UserId,
		comparison.FilmAId,
		comparison.FilmBId,
		comparison.WinningFilmId,
		comparison.ComparisonDate,
		comparison.WasEqual,
	)

	err := row.Scan(
		&createdComparison.ID,
		&createdComparison.UserId,
		&createdComparison.FilmAId,
		&createdComparison.FilmBId,
		&createdComparison.WinningFilmId,
		&createdComparison.ComparisonDate,
		&createdComparison.WasEqual,
	)

	if err != nil {
		return nil, err
	}

	return createdComparison, nil
}

func (s *store) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
	query := /* sql */ `
		SELECT COUNT(*) 
		FROM comparison_histories 
		WHERE user_id = $1 
		AND (
		(film_a_film_id = $2 AND film_b_film_id = $3) 
		OR 
		(film_a_film_id = $3 AND film_b_film_id = $2)
		)
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, userId, filmAId, filmBId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *store) GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error) {
	query := /* sql */ `
		SELECT comparison_history_id, user_id, film_a_film_id, film_b_film_id, winning_film_film_id, comparison_date, was_equal
		FROM comparison_histories
		WHERE user_id = $1
		ORDER BY comparison_date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comparisons []domain.ComparisonHistory
	for rows.Next() {
		var comparison domain.ComparisonHistory
		err := rows.Scan(
			&comparison.ID,
			&comparison.UserId,
			&comparison.FilmAId,
			&comparison.FilmBId,
			&comparison.WinningFilmId,
			&comparison.ComparisonDate,
			&comparison.WasEqual,
		)
		if err != nil {
			return nil, err
		}
		comparisons = append(comparisons, comparison)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comparisons, nil
}

// BulkGetRatings fetches multiple ratings for given film IDs in a single query
func (s *store) BulkGetRatings(ctx context.Context, userId uuid.UUID, filmIds []uuid.UUID) (map[uuid.UUID]*domain.UserFilmRating, error) {
	if len(filmIds) == 0 {
		return make(map[uuid.UUID]*domain.UserFilmRating), nil
	}

	// Build query with placeholders
	query := /* sql */ `
		SELECT user_film_rating_id, user_id, film_id, elo_rating, number_of_comparisons, last_updated, initial_rating, k_constant_value
		FROM user_film_ratings
		WHERE user_id = $1 AND film_id = ANY($2)
	`

	rows, err := s.db.QueryContext(ctx, query, userId, filmIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := make(map[uuid.UUID]*domain.UserFilmRating)
	for rows.Next() {
		rating := &domain.UserFilmRating{}
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
		ratings[rating.FilmId] = rating
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ratings, nil
}

// BulkHasBeenCompared checks if multiple film pairs have been compared
func (s *store) BulkHasBeenCompared(ctx context.Context, userId uuid.UUID, pairs []domain.ComparisonPair) (map[string]bool, error) {
	if len(pairs) == 0 {
		return make(map[string]bool), nil
	}

	// Build query to check all pairs at once
	query := /* sql */ `
		SELECT film_a_film_id, film_b_film_id 
		FROM comparison_histories 
		WHERE user_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build a set of existing comparisons
	existingComparisons := make(map[string]bool)
	for rows.Next() {
		var filmAId, filmBId uuid.UUID
		if err := rows.Scan(&filmAId, &filmBId); err != nil {
			return nil, err
		}
		// Store both directions
		key1 := filmAId.String() + "-" + filmBId.String()
		key2 := filmBId.String() + "-" + filmAId.String()
		existingComparisons[key1] = true
		existingComparisons[key2] = true
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Check each requested pair
	result := make(map[string]bool)
	for _, pair := range pairs {
		key := pair.FilmA.FilmId.String() + "-" + pair.FilmB.FilmId.String()
		result[key] = existingComparisons[key]
	}

	return result, nil
}

// BulkInsertComparisons inserts multiple comparison records in one query
func (s *store) BulkInsertComparisons(ctx context.Context, comparisons []domain.ComparisonHistory) error {
	if len(comparisons) == 0 {
		return nil
	}

	query := /* sql */ `
		INSERT INTO comparison_histories (comparison_history_id, user_id, film_a_film_id, film_b_film_id, winning_film_film_id, comparison_date, was_equal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Use a prepared statement for batch inserts
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, comparison := range comparisons {
		if comparison.ID == uuid.Nil {
			comparison.ID = uuid.New()
		}
		_, err := stmt.ExecContext(ctx,
			comparison.ID,
			comparison.UserId,
			comparison.FilmAId,
			comparison.FilmBId,
			comparison.WinningFilmId,
			comparison.ComparisonDate,
			comparison.WasEqual,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// BulkUpdateRatings updates multiple ratings using a CASE statement
func (s *store) BulkUpdateRatings(ctx context.Context, tx *sql.Tx, ratings []domain.UserFilmRating) error {
	if len(ratings) == 0 {
		return nil
	}

	// Use prepared statement for each update within the transaction
	query := /* sql */ `
		UPDATE user_film_ratings
		SET elo_rating = $1,
		    number_of_comparisons = $2,
		    last_updated = $3,
		    k_constant_value = $4
		WHERE user_film_rating_id = $5
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rating := range ratings {
		_, err := stmt.ExecContext(ctx,
			rating.EloRating,
			rating.NumberOfComparisons,
			rating.LastUpdated,
			rating.KConstantValue,
			rating.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// BeginTx starts a new transaction
func (s *store) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}
