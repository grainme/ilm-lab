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

type ReviewService struct {
	reviewRepo repository.ReviewRepository
	rdb        *redis.Client
}

func NewReviewService(repo repository.ReviewRepository, rdb *redis.Client) *ReviewService {
	return &ReviewService{
		reviewRepo: repo,
		rdb:        rdb,
	}
}

func (s *ReviewService) GetAllReviews(ctx context.Context) ([]domain.Review, error) {
	reviews, err := s.reviewRepo.GetAllReviews(ctx)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *ReviewService) GetAllReviewsByMovieId(ctx context.Context, movieId uuid.UUID) ([]domain.Review, error) {
	reviews, err := s.reviewRepo.GetAllReviewsByMovieId(ctx, movieId)
	if err != nil {
		return []domain.Review{}, err
	}

	return reviews, nil
}

func (s *ReviewService) AddReview(ctx context.Context, review *domain.Review) (domain.Review, error) {
	insertedReview, err := s.reviewRepo.AddReview(ctx, review)
	if err != nil {
		return domain.Review{}, err
	}

	// invalidate the movie cache since its data (avg_rating) has changed
	err = cache.DelMovie(ctx, s.rdb, review.MovieID)
	if err != nil {
		// Log the error but don't crash. The cache will expire on its own.
		log.Printf("failed to invalidate movie cache for movie %s: %v", review.MovieID, err)
	}

	return insertedReview, nil
}
