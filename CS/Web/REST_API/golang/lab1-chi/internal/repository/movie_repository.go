package repository

import (
	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

// interface deftinio for movie repository (Data Access Layer)
// regardles the data source (cache, files, Db, in-memory...)
type MovieRepository interface {
	GetAllMovies() []*domain.Movie
	GetMovieById(id uuid.UUID) (*domain.Movie, error)
	AddMovie(movie *domain.Movie) (*domain.Movie, error)
	UpdateMovieById(id uuid.UUID, newRating int) (*domain.Movie, error)
	DeleteMovieById(id uuid.UUID) error
}

/*
 * These are the HTTP requests, i'm supporting:
 *
 * GET /movies - list all (200 OK)
 * GET /movies/:id - get one (200 or 404)
 * POST /movies - create (201 Created or 400 Bad Request)
 * PUT /movies/:id/rating - update (200 or 404)
 * DELETE /movies/:id - remove (204 No Content or 404)
 *
 */
