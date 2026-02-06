package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/grainme/movie-api/internal/auth"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the token and validates it and then attach to context
		bearerToken := strings.Split(r.Header.Get("Authorization"), " ")
		if len(bearerToken) < 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access token missing"))
			return
		}

		token := bearerToken[1]
		claims, err := auth.ValidateAccessToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}
		newCtx := context.WithValue(r.Context(), "user", claims)
		requestWithModifiedCtx := r.WithContext(newCtx)

		next.ServeHTTP(w, requestWithModifiedCtx)
	})
}
