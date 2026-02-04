package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/service"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

func (h *ReviewHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.reviewService.GetAllReviews(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) GetAllReviewsByMovieId(w http.ResponseWriter, r *http.Request) {
	movieId, err := extractIdAndParse(w, r)
	if err != nil {
		return
	}

	reviews, err := h.reviewService.GetAllReviewsByMovieId(r.Context(), movieId)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	var review domain.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	insertedReview, err := h.reviewService.AddReview(r.Context(), &review)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, insertedReview)
}
