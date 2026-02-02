package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/database"
	"github.com/grainme/movie-api/internal/domain"
)

type PostgresReviewRepository struct {
	dbQueries *database.Queries
}

func NewPostgresReviewRepository(db database.DBTX) *PostgresReviewRepository {
	return &PostgresReviewRepository{
		dbQueries: database.New(db),
	}
}

func (r *PostgresReviewRepository) GetAllReviews(ctx context.Context) ([]domain.Review, error) {
	reviews, err := r.dbQueries.GetAllReviews(ctx)
	if err != nil {
		return nil, err
	}

	reviewsList := make([]domain.Review, len(reviews))
	for idx, r := range reviews {
		reviewsList[idx] = toDomainReview(r)
	}
	return reviewsList, nil
}

func (r *PostgresReviewRepository) AddReview(ctx context.Context, review *domain.Review) (domain.Review, error) {
	comment := sql.NullString{}
	if review.Comment != nil {
		comment = sql.NullString{
			String: *review.Comment,
			Valid:  true,
		}
	}

	dbReview, err := r.dbQueries.AddReview(ctx, database.AddReviewParams{
		ID:       review.ID,
		UserName: review.Username,
		Rating:   review.Rating,
		Comment:  comment,
		MovieID:  review.MovieID,
	})
	if err != nil {
		return domain.Review{}, err
	}

	return toDomainReview(dbReview), nil
}

func (r *PostgresReviewRepository) GetAllReviewsByMovieId(ctx context.Context, movieId uuid.UUID) ([]domain.Review, error) {
	reviews, err := r.dbQueries.GetAllReviewsByMovieId(ctx, movieId)
	if err != nil {
		return nil, err
	}

	reviewsList := make([]domain.Review, len(reviews))
	for idx, r := range reviews {
		reviewsList[idx] = toDomainReview(r)
	}
	return reviewsList, nil
}

// Helper(Mapper)
func toDomainReview(dbReview database.Review) domain.Review {
	review := domain.Review{
		ID:       dbReview.ID,
		Username: dbReview.UserName,
		Rating:   dbReview.Rating,
		Comment:  nil,
		MovieID:  dbReview.MovieID,
	}

	if dbReview.Comment.Valid {
		review.Comment = &dbReview.Comment.String
	}

	return review
}
