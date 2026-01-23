package service

import (
	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
)

type MovieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: repo,
	}
}

func (s *MovieService) GetAllMovies() []*domain.Movie {
	// should we have any other logic (validation...)?
	movies := s.movieRepo.GetAllMovies()
	return movies
}

func (s *MovieService) GetMovieById(id uuid.UUID) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetMovieById(id)
	return movie, err
}

func (s *MovieService) AddMovie(movie *domain.Movie) (*domain.Movie, error) {
	if movie == nil {
		return nil, domain.ErrInvalidMovie
	}
	if movie.Rating < 0 || movie.Rating > 10 {
		return nil, domain.ErrInvalidRating
	}
	movie, err := s.movieRepo.AddMovie(movie)
	return movie, err
}

func (s *MovieService) UpdateMovieById(id uuid.UUID, newRating int) (*domain.Movie, error) {
	if newRating < 0 || newRating > 10 {
		return nil, domain.ErrInvalidRating
	}
	movie, err := s.movieRepo.UpdateMovieById(id, newRating)
	return movie, err
}

func (s *MovieService) DeleteMovieById(id uuid.UUID) error {
	err := s.movieRepo.DeleteMovieById(id)
	return err
}
