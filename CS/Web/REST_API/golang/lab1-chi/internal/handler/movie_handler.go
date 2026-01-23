package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/grainme/movie-api/internal/service"
)

type MovieHandler struct {
	movieService service.MovieService
}

func NewMovieHandler(service service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: service,
	}
}

func (h *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies := h.movieService.GetAllMovies()
	var moviesData []byte

	if len(movies) == 0 {
		moviesData = []byte("No movies available.")
	} else {
		data, err := json.Marshal(movies)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		moviesData = data
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(moviesData)
	w.WriteHeader(http.StatusAccepted)
}
