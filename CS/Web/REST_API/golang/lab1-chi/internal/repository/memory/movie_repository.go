/*
 * Concrete implementation for in-memory storage
 */
package memory

import (
	"sync"

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

func (r *MemoryMovieRepository) GetAllMovies() []*domain.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.movies
}
