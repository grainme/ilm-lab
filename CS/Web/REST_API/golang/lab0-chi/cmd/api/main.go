package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	handlers "github.com/grainme/movie-api/internal/handler"
	"github.com/grainme/movie-api/internal/repository/memory"
	"github.com/grainme/movie-api/internal/service"
)

func main() {
	movies := []*domain.Movie{
		{Id: uuid.New(), Title: "BadMovie1", Rating: 2},
		{Id: uuid.New(), Title: "BadMovie2", Rating: 5},
	}

	movieRepo := memory.NewMemoryMovieRepository(movies)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/movies", movieHandler.GetAllMovies)
	r.Get("/movies/{id}", movieHandler.GetMovieById)
	r.Post("/movies", movieHandler.AddMovie)
	r.Put("/movies/{id}", movieHandler.UpdateMovieById)
	r.Delete("/movies/{id}", movieHandler.DeleteById)

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
