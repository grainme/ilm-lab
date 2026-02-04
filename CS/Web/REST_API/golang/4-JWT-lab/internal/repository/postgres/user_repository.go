package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/auth"
	"github.com/grainme/movie-api/internal/database"
	"github.com/grainme/movie-api/internal/domain"
)

type PostgresUserRepository struct {
	dbQueries *database.Queries
}

func NewPostgresUserRepository(db database.DBTX) *PostgresUserRepository {
	return &PostgresUserRepository{
		dbQueries: database.New(db),
	}
}

func (r *PostgresUserRepository) AddUser(ctx context.Context, user domain.CreateUserRequest) (domain.User, error) {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	createdUser, err := r.dbQueries.AddUser(ctx, database.AddUserParams{
		ID:           uuid.New(),
		Username:     user.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return domain.User{}, err
	}

	toDomainUser, err := DatabaseUserToDomainUser(createdUser)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser, nil
}

func (r *PostgresUserRepository) FindUserByName(ctx context.Context, username string) (domain.User, error) {
	dbUser, err := r.dbQueries.FindUserByName(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	toDomainUser, err := DatabaseUserToDomainUser(dbUser)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser, nil
}

// -------- helpers (mappers)
func DatabaseUserToDomainUser(du database.User) (domain.User, error) {
	var role domain.Role
	if du.Role.Valid {
		role = domain.Role(du.Role.UserRole)
	} else {
		return domain.User{}, errors.New("invalid user role")
	}

	return domain.User{
		ID:           du.ID,
		Username:     du.Username,
		PasswordHash: du.PasswordHash,
		Role:         role,
		CreatedAt:    du.CreatedAt.Time,
		UpdatedAt:    du.UpdatedAt.Time,
	}, nil
}
