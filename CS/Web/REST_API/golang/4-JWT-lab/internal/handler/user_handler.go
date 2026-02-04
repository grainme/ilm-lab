package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// the request should contain a body (username, password)
	var userRequest domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Printf("Decoding failed: %v", err)
		http.Error(w, "Decoding failed", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	user, err := h.userService.Login(r.Context(), userRequest.Username, userRequest.Password)
	if err != nil {
		log.Printf("user could not login: %v", err)
		http.Error(w, "Login failed", http.StatusBadRequest)
		return
	}
	// bad validation, I think!
	if user.Username == "" {
		http.Error(w, "Login failed", http.StatusBadRequest)
		return
	}

	log.Printf("User connected\n")
	// we should not return the user as payload (just for testing)
	respondJSON(w, http.StatusFound, user)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// the request should contain a body (username, password)
	var userRequest domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Printf("Decoding failed: %v", err)
		http.Error(w, "Decoding failed", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	user, err := h.userService.Register(r.Context(), userRequest)
	if err != nil {
		log.Printf("user could not sign up: %v", err)
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	log.Printf("User registered\n")
	// we should not return the user as payload (just for testing)
	respondJSON(w, http.StatusCreated, user)
}
