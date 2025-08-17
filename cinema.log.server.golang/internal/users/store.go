package users

import (
	"context"
	"database/sql"
	"errors"

	"cinema.log.server.golang/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{
		db: db,
	}
}

func (s *store) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User

	query := `SELECT id, github_id, name, username, profile_pic_url, created_at, updated_at FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.GithubId, &user.Name, &user.Username, &user.ProfilePicURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *store) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, github_id, name, username, profile_pic_url, created_at, updated_at 
	          FROM users WHERE id = $1`
	
	user := &domain.User{}
	row := s.db.QueryRowContext(ctx, query, id)
	
	err := row.Scan(&user.ID, &user.GithubId, &user.Name, &user.Username, &user.ProfilePicURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

func (s *store) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Generate a new UUID if not provided
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	
	query := `
		INSERT INTO users (id, github_id, name, username, profile_pic_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
		RETURNING created_at, updated_at`
	
	err := s.db.QueryRowContext(ctx, query, user.ID, user.GithubId, user.Name, user.Username, user.ProfilePicURL).
	Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *store) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users 
		SET github_id = $2, name = $3, username = $4, profile_pic_url = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`
	
	err := s.db.QueryRowContext(ctx, query, user.ID, user.GithubId, user.Name, user.Username, user.ProfilePicURL).Scan(&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

func (s *store) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

func (s *store) GetOrCreateUserByGithubId(ctx context.Context, githubID int64, 
	name string, username string, avatarUrl string) (*domain.User, error) {
	query := `SELECT id, github_id, name, username, profile_pic_url, created_at, updated_at 
			  FROM users WHERE github_id = $1`

	user := &domain.User{}
	row := s.db.QueryRowContext(ctx, query, githubID)

	err := row.Scan(&user.ID, &user.GithubId, &user.Name, &user.Username, &user.ProfilePicURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// If not found, create a new user
			user = &domain.User{
				GithubId: githubID,
				Name:     name,
				Username: username,
				ProfilePicURL: avatarUrl,
			}
			return s.CreateUser(ctx, user)
		}
		return nil, err
	}

	return user, nil
}