/*
 * Concrete implementation for in-memory storage
 */
package memory

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

type MemoryMovieRepository struct {
	mu     sync.RWMutex
	movies []*domain.Movie
}

func NewMemoryMovieRepository(movies []*domain.Movie) *MemoryMovieRepository {
	return &MemoryMovieRepository{
		mu:     sync.RWMutex{},
		movies: movies,
	}
}

func (r *MemoryMovieRepository) GetAllMovies(ctx context.Context) []*domain.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// this is a deep copy (instead of copy(?,?))
	// to avoid mutating internal props of the repo by the caller
	result := make([]*domain.Movie, len(r.movies))
	for idx, movie := range r.movies {
		movieCopy := *movie
		result[idx] = &movieCopy
	}

	return result
}

func (r *MemoryMovieRepository) GetMovieById(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, movie := range r.movies {
		if movie.Id == id {
			movieCopy := *movie
			return &movieCopy, nil
		}
	}

	return nil, domain.ErrMovieNotFound
}

func (r *MemoryMovieRepository) AddMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.movies = append(r.movies, movie)

	movieCopy := *movie
	return &movieCopy, nil
}

func (r *MemoryMovieRepository) UpdateMovieById(ctx context.Context, id uuid.UUID, newRating int) (*domain.Movie, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, mv := range r.movies {
		if mv.Id == id {
			mv.Rating = newRating
			// return copy to avoid data race
			movieCopy := *mv
			return &movieCopy, nil
		}
	}

	return nil, domain.ErrMovieNotFound
}

func (r *MemoryMovieRepository) DeleteMovieById(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, mv := range r.movies {
		if mv.Id == id {
			r.movies = slices.Delete(r.movies, idx, idx+1)
			return nil
		}
	}

	return domain.ErrMovieNotFound
}
