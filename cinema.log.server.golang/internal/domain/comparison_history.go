package domain

import (
	"time"

	"github.com/google/uuid"
)

type ComparisonHistory struct {
	ID             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"userId"`
	FilmAId        uuid.UUID `json:"filmAId"`
	FilmBId        uuid.UUID `json:"filmBId"`
	WinningFilmId  uuid.UUID `json:"winningFilmId"`
	ComparisonDate time.Time `json:"comparisonDate"`
	WasEqual       bool      `json:"wasEqual"`
}
