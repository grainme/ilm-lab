package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/cache"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
	"github.com/redis/go-redis/v9"
)

type MovieService struct {
	movieRepo repository.MovieRepository
	rdb       *redis.Client
}

func NewMovieService(repo repository.MovieRepository, rdb *redis.Client) *MovieService {
	return &MovieService{
		movieRepo: repo,
		rdb:       rdb,
	}
}

func (s *MovieService) GetAllMovies(ctx context.Context) []*domain.Movie {
	// should we have any other logic (validation...)?
	movies := s.movieRepo.GetAllMovies(ctx)
	return movies
}

func (s *MovieService) GetMovieById(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	startTime := time.Now()
	movie, err := cache.GetMovieById(ctx, *s.rdb, id)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, fmt.Errorf("movie is nil")
	}

	// cache hit
	if movie != nil {
		log.Printf("CACHE HIT %s: %dms\n", cache.MovieKey(movie.ID), time.Since(startTime).Milliseconds())
		return movie, nil
	}

	log.Printf("CACHE MISS: %s\n", cache.MovieKey(movie.ID))
	movie, err = s.movieRepo.GetMovieById(ctx, id)

	// save in cache
	if err := cache.SetMovie(ctx, s.rdb, id, *movie); err != nil {
		log.Printf("could not cache movie: %v\n", err)
	}

	log.Printf("(Postgres) duration: %dms\n", time.Since(startTime).Milliseconds())
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
