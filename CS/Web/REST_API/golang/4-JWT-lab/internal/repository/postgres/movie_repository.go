package postgres

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/database"
	"github.com/grainme/movie-api/internal/domain"
)

type PostgresMovieRepository struct {
	dbQueries *database.Queries
}

func NewPostgresMovieRepository(db database.DBTX) *PostgresMovieRepository {
	return &PostgresMovieRepository{
		dbQueries: database.New(db),
	}
}

func (r *PostgresMovieRepository) GetAllMovies(ctx context.Context) []*domain.Movie {
	movies, err := r.dbQueries.GetMovies(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	moviesList := make([]*domain.Movie, len(movies))
	for idx, mv := range movies {
		moviesList[idx] = toDomainMovieFromDatabaseMovie(mv)
	}

	return moviesList
}

func (r *PostgresMovieRepository) GetMovieById(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie, err := r.dbQueries.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainMovieFromDatabaseMovie(movie), nil
}

func (r *PostgresMovieRepository) AddMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	dbMovie, err := r.dbQueries.AddMovie(ctx, database.AddMovieParams{
		ID:       uuid.New(),
		Title:    movie.Title,
		Director: movie.Director,
		Year:     movie.Year,
	})
	if err != nil {
		return nil, err
	}

	insertedMovie := toDomainMovieFromDatabaseMovie(dbMovie)
	return insertedMovie, nil
}

func (r *PostgresMovieRepository) UpdateMovieTitleById(ctx context.Context, id uuid.UUID, title string) error {
	err := r.dbQueries.UpdateMovieTitleById(ctx, database.UpdateMovieTitleByIdParams{
		ID:    id,
		Title: title,
	})

	return err
}

func (r *PostgresMovieRepository) DeleteMovieById(ctx context.Context, id uuid.UUID) error {
	err := r.dbQueries.DeleteMovieById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresMovieRepository) GetMovieWithReviews(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie, err := r.dbQueries.GetMovieWithReviews(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainMovieFromGetMovieWithReviewsRow(movie), nil
}

// Helper function (Mapper)
func toDomainMovieFromDatabaseMovie(movie database.Movie) *domain.Movie {
	domainMovie := domain.Movie{
		ID:       movie.ID,
		Title:    movie.Title,
		Director: movie.Director,
		Year:     movie.Year,
	}

	return &domainMovie
}

func toDomainMovieFromGetMovieWithReviewsRow(movie database.GetMovieWithReviewsRow) *domain.Movie {
	domainMovie := domain.Movie{
		ID:            movie.ID,
		Title:         movie.Title,
		Director:      movie.Director,
		Year:          movie.Year,
		AverageRating: movie.AvgRating,
		ReviewsCount:  movie.ReviewsCount,
	}

	return &domainMovie
}
