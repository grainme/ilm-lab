package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: repo,
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

	return insertedReview, nil
}
