package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

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
