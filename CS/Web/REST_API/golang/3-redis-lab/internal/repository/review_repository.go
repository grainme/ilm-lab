package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

type ReviewRepository interface {
	AddReview(ctx context.Context, review *domain.Review) (domain.Review, error)
	GetAllReviews(ctx context.Context) ([]domain.Review, error)
	GetAllReviewsByMovieId(ctx context.Context, movieId uuid.UUID) ([]domain.Review, error)
}
