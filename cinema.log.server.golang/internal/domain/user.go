package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	GithubId      int64     `json:"githubId"`
	Name          string    `json:"name"`
	Username      string    `json:"username"`
	ProfilePicURL string    `json:"profilePicUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
