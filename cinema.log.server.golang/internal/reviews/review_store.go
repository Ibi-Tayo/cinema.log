package reviews

import (
	"context"
	"database/sql"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) ReviewStore {
	return &store{
		db: db,
	}
}

func (s *store) GetAllReviewsByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Review, error) {
	query := `SELECT review_id, content, date, rating, film_id, user_id 
	          FROM reviews WHERE user_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var review domain.Review
		err := rows.Scan(&review.ID, &review.Content, &review.Date, &review.Rating, &review.FilmId, &review.UserId)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *store) CreateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	if review.ID == uuid.Nil {
		review.ID = uuid.New()
	}

	query := `
		INSERT INTO reviews (review_id, content, date, rating, film_id, user_id) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.ExecContext(ctx, query, review.ID, review.Content, review.Date, review.Rating, review.FilmId, review.UserId)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (s *store) UpdateReview(ctx context.Context, review domain.Review) (*domain.Review, error) {
	query := `
		UPDATE reviews 
		SET content = $1, date = $2, rating = $3, film_id = $4, user_id = $5
		WHERE review_id = $6`

	result, err := s.db.ExecContext(ctx, query, review.Content, review.Date, review.Rating, review.FilmId, review.UserId, review.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrReviewNotFound
	}

	return &review, nil
}

func (s *store) DeleteReview(ctx context.Context, reviewId uuid.UUID) error {
	query := `DELETE FROM reviews WHERE review_id = $1`

	result, err := s.db.ExecContext(ctx, query, reviewId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrReviewNotFound
	}

	return nil
}
