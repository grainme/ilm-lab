package main

import (
	"net/http"

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
		&domain.Movie{
			Id:     uuid.New(),
			Title:  "BadMovie1",
			Rating: 2,
		},
		&domain.Movie{
			Id:     uuid.New(),
			Title:  "BadMovie2",
			Rating: 5,
		},
	}

	movieRepo := memory.NewMemoryMovieRepository(movies)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", movieHandler.GetAllMovies)

	http.ListenAndServe(":3000", r)
}
