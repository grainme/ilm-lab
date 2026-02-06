package service

import (
	"context"
	"log"

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
	movie, err := cache.GetMovieById(ctx, s.rdb, id)
	if err != nil {
		log.Printf("failed to get movie from cache: %v", err)
	}

	// cache hit
	if movie != nil {
		// view counter (write-behind strategy)
		go func() {
			count, err := cache.IncrementViewCount(context.Background(), s.rdb, id)
			if err != nil {
				log.Printf("failed to increment view count for movie %s: %v", id, err)
				return
			}

			if count%100 == 0 {
				// Every 100 views, sync to DB
				// We don't have views track on the DB. but I got the idea of write-behind
				// s.repo.UpdateViewCount(context.Background(), id, count)
				log.Printf("syncing view count to DB: %d", count)
			}
		}()

		return movie, nil
	}

	log.Printf("CACHE MISS: %s\n", cache.MovieKey(id))
	movie, err = s.movieRepo.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}

	// save in cache
	if err := cache.SetMovie(ctx, s.rdb, id, *movie); err != nil {
		log.Printf("could not cache movie: %v\n", err)
	}

	return movie, err
}

func (s *MovieService) AddMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if movie == nil {
		return nil, domain.ErrInvalidMovie
	}

	if movie.Title == "" || len(movie.Title) > 40 {
		return nil, domain.ErrInvalidMovie
	}

	movie, err := s.movieRepo.AddMovie(ctx, movie)
	return movie, err
}

func (s *MovieService) UpdateMovieTitleById(ctx context.Context, id uuid.UUID, title string) error {
	if title == "" || len(title) > 40 {
		return domain.ErrInvalidMovie
	}

	err := s.movieRepo.UpdateMovieTitleById(ctx, id, title)
	if err != nil {
		return err
	}

	err = cache.DelMovie(ctx, s.rdb, id)
	return err
}

func (s *MovieService) DeleteMovieById(ctx context.Context, id uuid.UUID) error {
	err := s.movieRepo.DeleteMovieById(ctx, id)
	if err != nil {
		return err
	}

	err = cache.DelMovie(ctx, s.rdb, id)
	return err
}

func (s *MovieService) GetMovieWithReviews(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetMovieWithReviews(ctx, id)
	return movie, err
}
