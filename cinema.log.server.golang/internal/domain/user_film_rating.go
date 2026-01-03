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
	NumberOfComparisons int       `json:"numberOfComparisons"`
	LastUpdated         time.Time `json:"lastUpdated"`
	InitialRating       float32   `json:"initialRating"`
	KConstantValue      float64   `json:"kConstantValue"`
}

type UserFilmRatingDetail struct {
	Rating    UserFilmRating `json:"rating"`
	FilmTitle string         `json:"filmTitle"`
	FilmReleaseYear string   `json:"filmReleaseYear"`
	FilmPosterURL string     `json:"filmPosterUrl"`
}

type ComparisonPair struct {
    FilmA UserFilmRating
    FilmB UserFilmRating
}