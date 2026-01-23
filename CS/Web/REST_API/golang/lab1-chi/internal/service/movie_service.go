package service

import (
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
)

type MovieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) MovieService {
	return MovieService{
		movieRepo: repo,
	}
}

func (s *MovieService) GetAllMovies() []*domain.Movie {
	// should we have any other logic (validation...)?
	return s.movieRepo.GetAllMovies()
}
