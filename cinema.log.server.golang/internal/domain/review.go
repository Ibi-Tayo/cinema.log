package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"title"`
	Date    time.Time `json:"date"`
	Rating  float32   `json:"rating"`
	FilmId  uuid.UUID `json:"filmId"`
	UserId  uuid.UUID `json:"userId"`
}
