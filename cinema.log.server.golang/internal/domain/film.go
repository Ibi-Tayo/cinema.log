package domain

import (
	"github.com/google/uuid"
)

type Film struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	Director    string    `json:"director"`
	PosterUrl   string    `json:"posterUrl"`
	ReleaseYear int32     `json:"releaseYear"`
}
