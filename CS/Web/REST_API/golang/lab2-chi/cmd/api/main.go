package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	handlers "github.com/grainme/movie-api/internal/handler"
	"github.com/grainme/movie-api/internal/repository/postgres"
	"github.com/grainme/movie-api/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// movies := []*domain.Movie{
	// 	{Id: uuid.New(), Title: "BadMovie1", Rating: 2},
	// 	{Id: uuid.New(), Title: "BadMovie2", Rating: 5},
	// }

	dbDSN := os.Getenv("DB_DSN")

	log.Println("Running database migrations...")
	m, err := migrate.New("file://db/migrations", dbDSN)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// ErrNoChanges means DB is already up-to-date
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
		os.Exit(1)
	}

	movieRepo := postgres.NewPostgresMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	reviewRepo := postgres.NewPostgresReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// movie routes
	r.Get("/movies", movieHandler.GetAllMovies)
	r.Get("/movies/{id}", movieHandler.GetMovieById)
	r.Post("/movies", movieHandler.AddMovie)
	r.Delete("/movies/{id}", movieHandler.DeleteById)
	r.Get("/movies/{id}/reviews", movieHandler.GetMovieWithReviews)
	// reviews routes
	r.Get("/reviews", reviewHandler.GetAllReviews)
	r.Post("/reviews", reviewHandler.AddReview)
	r.Get("/reviews/{id}", reviewHandler.GetAllReviewsByMovieId)

	port := os.Getenv("PORT")
	if port == "" {
		// default value
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
