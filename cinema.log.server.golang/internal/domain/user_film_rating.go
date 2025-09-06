package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserFilmRating struct {
	ID                  uuid.UUID `json:"id"`
	UserId              uuid.UUID `json:"userId"`
	FilmId              uuid.UUID `json:"filmId"`
	EloRating           float64   `json:"eloRating"`
	NumberOfComparisons string    `json:"numberOfComparisons"`
	LastUpdated         time.Time `json:"lastUpdated"`
	InitialRating       float32   `json:"initialRating"`
	KConstantValue      float64   `json:"kConstantValue"`
}
