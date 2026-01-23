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
	movies := h.movieService.GetAllMovies()

	data, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Marshalling failed", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "UUID parsing failed", http.StatusInternalServerError)
	}

	movie, err := h.movieService.GetMovieById(uuidFromId)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			w.Header().Add("Content-Type", "application/text")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Movie not found"))
			return
		} else {
			http.Error(w, "Service failed", http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Marshalling failed", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *MovieHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	var movieData domain.Movie
	if err := json.NewDecoder(r.Body).Decode(&movieData); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.AddMovie(&movieData)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidMovie) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Service failed", http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Marshalling failed", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *MovieHandler) UpdateMovieById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "UUID parsing failed", http.StatusInternalServerError)
	}

	var payload struct {
		Rating int `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.UpdateMovieById(uuidFromId, payload.Rating)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			w.Header().Add("Content-Type", "application/text")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Movie not found"))
			return
		} else {
			http.Error(w, "Service failed", http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Marshalling failed", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *MovieHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "UUID parsing failed", http.StatusInternalServerError)
	}

	err = h.movieService.DeleteMovieById(uuidFromId)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			w.Header().Add("Content-Type", "application/text")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Movie not found"))
			return
		} else {
			http.Error(w, "Service failed", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
