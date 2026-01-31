package service

import (
	"context"

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

func (s *MovieService) GetAllMovies(ctx context.Context) []*domain.Movie {
	// should we have any other logic (validation...)?
	movies := s.movieRepo.GetAllMovies(ctx)
	return movies
}

func (s *MovieService) GetMovieById(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetMovieById(ctx, id)
	return movie, err
}

func (s *MovieService) AddMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if movie == nil {
		return nil, domain.ErrInvalidMovie
	}
	// this does not make sense, we should not allow to add a movie with a rating
	// this done from the rating service
	// if *movie.AvgRating < 0 || *movie.AvgRating > 10 {
	// 	return nil, domain.ErrInvalidRating
	// }
	if movie.Title == "" || len(movie.Title) > 40 {
		return nil, domain.ErrInvalidMovie
	}

	movie, err := s.movieRepo.AddMovie(ctx, movie)
	return movie, err
}

func (s *MovieService) DeleteMovieById(ctx context.Context, id uuid.UUID) error {
	err := s.movieRepo.DeleteMovieById(ctx, id)
	return err
}

func (s *MovieService) GetMovieWithReviews(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetMovieWithReviews(ctx, id)
	return movie, err
}
