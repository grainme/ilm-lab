package domain

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Director string    `json:"director"`
	Year     int32     `json:"year"`
}
