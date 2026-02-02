package domain

import "github.com/google/uuid"

type Movie struct {
	Id     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Rating int       `json:"rating"`
}
