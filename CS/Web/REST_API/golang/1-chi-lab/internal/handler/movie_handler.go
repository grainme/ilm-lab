package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *MovieHandler) UpdateMovieById(w http.ResponseWriter, r *http.Request) {
	uuidFromId, err := extractIdAndParse(w, r)
	if err != nil {
		return
	}

	var payload struct {
		Rating int `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.UpdateMovieById(r.Context(), uuidFromId, payload.Rating)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, movie)
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

// ------ private helpers ------

func respondJSON(w http.ResponseWriter, status int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal json response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func respondError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrMovieNotFound) {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
		return
	} else {
		http.Error(w, "movie service failed", http.StatusInternalServerError)
		return
	}
}

func extractIdAndParse(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	id := chi.URLParam(r, "id")
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "UUID parsing failed", http.StatusBadRequest)
		return uuid.UUID{}, err
	}

	return uuidFromId, nil
}
