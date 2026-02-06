package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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

	userResponse, err := h.userService.Login(r.Context(), userRequest.Username, userRequest.Password)
	if err != nil {
		log.Printf("user could not login: %v", err)
		http.Error(w, "Login failed", http.StatusBadRequest)
		return
	}

	log.Printf("User connected\n")
	respondJSON(w, http.StatusOK, userResponse)
}

// we can't revoke JWT because it's stateless
// but we can revoke the refresh token on the other hand
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// we should extract { "refresh_token": "<uuid>" } from r.body
	var refreshToken struct {
		Token string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		log.Printf("Decoding failed: %v", err)
		http.Error(w, "Decoding failed", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	token, err := uuid.Parse(refreshToken.Token)
	if err != nil {
		log.Printf("UUID parsing failed: %v", err)
		http.Error(w, "UUID parsing failed", http.StatusBadRequest)
		return
	}

	err = h.userService.Logout(r.Context(), token)
	if err != nil {
		log.Printf("user could not logout: %v", err)
		http.Error(w, "Logout failed", http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, nil)
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

	log.Printf("User registered: %v\n", user)
	// i'm assuming register is different in behavior than login
	// therefor i'm not assigning the user "AccessToken" and "RefreshToken"
	respondJSON(w, http.StatusCreated, nil)
}

// Task 2: Create POST /auth/refresh endpoint. Client sends { "refresh_token": "<uuid>" }. Look it up in
// Redis. If found → delete old one, generate new access token + new refresh token, store new refresh token
// in Redis, return both.
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// we should extract { "refresh_token": "<uuid>" } from r.body
	var refreshToken struct {
		Token string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		log.Printf("Decoding failed: %v", err)
		http.Error(w, "Decoding failed", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	refreshTokenUUID, err := uuid.Parse(refreshToken.Token)
	if err != nil {
		log.Printf("UUID parsing failed: %v", err)
		http.Error(w, "UUID parsing failed", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RefreshToken(r.Context(), refreshTokenUUID)
	if err != nil {
		log.Printf("refresh token not found: %v", err)
		http.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}
	// here we  should refreshToken via the service.
	respondJSON(w, http.StatusOK, user)
}
