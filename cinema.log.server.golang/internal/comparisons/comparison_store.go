package comparisons

import (
	"context"
	"database/sql"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrComparisonNotFound = errors.New("comparison not found")
)

type ComparisonStore interface {
	CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error)
	HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error)
	GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error)
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) ComparisonStore {
	return &store{
		db: db,
	}
}

func (s *store) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
	if comparison.ID == uuid.Nil {
		comparison.ID = uuid.New()
	}

	query := `
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
	query := `
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
	query := `
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
