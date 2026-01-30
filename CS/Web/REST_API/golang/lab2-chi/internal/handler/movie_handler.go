package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/service"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: service,
	}
}

func (h *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies := h.movieService.GetAllMovies(r.Context())
	respondJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	uuidFromId, err := extractIdAndParse(w, r)
	if err != nil {
		return
	}

	movie, err := h.movieService.GetMovieById(r.Context(), uuidFromId)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	var movieData domain.Movie
	if err := json.NewDecoder(r.Body).Decode(&movieData); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.AddMovie(r.Context(), &movieData)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, movie)
}

func (h *MovieHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	uuidFromId, err := extractIdAndParse(w, r)
	if err != nil {
		return
	}

	err = h.movieService.DeleteMovieById(r.Context(), uuidFromId)
	if err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
