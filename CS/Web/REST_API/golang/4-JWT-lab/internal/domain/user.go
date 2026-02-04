package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	Regular Role = "regular"
	Admin   Role = "admin"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     Role      `json:"role"`
}

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string // this is Argon2id hash not plain password
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
