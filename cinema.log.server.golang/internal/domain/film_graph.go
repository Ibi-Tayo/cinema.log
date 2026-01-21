package domain

import "github.com/google/uuid"

type FilmGraphNode struct {
	UserID         uuid.UUID `json:"userId"`
	ExternalFilmID int       `json:"externalFilmId"`
	Title          string    `json:"title"`
}

type FilmGraphEdge struct {
	UserID     uuid.UUID `json:"userId"`
	EdgeId     uuid.UUID `json:"edgeId"`
	FromFilmID int       `json:"fromFilmId"`
	ToFilmID   int       `json:"toFilmId"`
}
