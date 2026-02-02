package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

// interface definition for movie repository (Data Access Layer)
// regardles the data source (cache, files, Db, in-memory...)
type MovieRepository interface {
	GetAllMovies(ctx context.Context) []*domain.Movie
	GetMovieById(ctx context.Context, id uuid.UUID) (*domain.Movie, error)
	AddMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	UpdateMovieTitleById(ctx context.Context, id uuid.UUID, title string) error
	DeleteMovieById(ctx context.Context, id uuid.UUID) error
	GetMovieWithReviews(ctx context.Context, id uuid.UUID) (*domain.Movie, error)
}
