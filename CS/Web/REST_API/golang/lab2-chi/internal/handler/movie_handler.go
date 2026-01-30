package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/service"
)

type MovieResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Director string    `json:"director"`
	Year     int32     `json:"year"`
}

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

	moviesResponse := make([]MovieResponse, len(movies))
	for idx, movie := range movies {
		movieResponse := MovieResponse{
			ID:       movie.ID,
			Title:    movie.Title,
			Director: movie.Director,
			Year:     movie.Year,
		}
		moviesResponse[idx] = movieResponse
	}

	respondJSON(w, http.StatusOK, moviesResponse)
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

	movieResponse := MovieResponse{
		ID:       movie.ID,
		Title:    movie.Title,
		Director: movie.Director,
		Year:     movie.Year,
	}

	respondJSON(w, http.StatusOK, movieResponse)
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

func (h *MovieHandler) GetMovieWithReviews(w http.ResponseWriter, r *http.Request) {
	uuidFromId, err := extractIdAndParse(w, r)
	if err != nil {
		return
	}

	movie, err := h.movieService.GetMovieWithReviews(r.Context(), uuidFromId)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, movie)
}
