package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	db, err := sql.Open("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	defer db.Close()

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

	r.Get("/movies", movieHandler.GetAllMovies)
	r.Get("/movies/{id}", movieHandler.GetMovieById)
	r.Post("/movies", movieHandler.AddMovie)
	r.Delete("/movies/{id}", movieHandler.DeleteById)
	r.Get("/reviews", reviewHandler.GetAllReviews)
	r.Get("/reviews/{id}", reviewHandler.GetAllReviewsByMovieId)
	r.Post("/reviews", reviewHandler.AddReview)

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
