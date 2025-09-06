package domain

import (
	"github.com/google/uuid"
)

type Film struct {
	ID          uuid.UUID `json:"id"`
	ExternalID  int     `json:"externalId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PosterUrl   string    `json:"posterUrl"`
	ReleaseYear string    `json:"releaseYear"`
}
