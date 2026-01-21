package domain

import (
	"github.com/google/uuid"
)

type Film struct {
	ID          uuid.UUID `json:"id"`
	ExternalID  int       `json:"externalId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PosterUrl   string    `json:"posterUrl"`
	ReleaseYear string    `json:"releaseYear"`
}

// Note this is a pure backend construct to track recommendations, not exposed via API
type FilmRecommendation struct {
	ID                     uuid.UUID `json:"id"`
	UserID                 uuid.UUID `json:"userId"`
	ExternalFilmID         int       `json:"externalFilmId"`
	HasSeen                bool      `json:"hasSeen"`
	HasBeenRecommended     bool      `json:"hasBeenRecommended"`
	RecommendationsGenerated bool     `json:"recommendationsGenerated"`
}
