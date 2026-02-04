package repository

import (
	"context"

	"github.com/grainme/movie-api/internal/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.CreateUserRequest) (domain.User, error)
	FindUserByName(ctx context.Context, username string) (domain.User, error)
}
