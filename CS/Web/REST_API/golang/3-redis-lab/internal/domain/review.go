package domain

import "github.com/google/uuid"

type Review struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"user_name"`
	Rating   int32     `json:"rating"`
	Comment  *string   `json:"comment"`
	MovieID  uuid.UUID `json:"movie_id"`
}
