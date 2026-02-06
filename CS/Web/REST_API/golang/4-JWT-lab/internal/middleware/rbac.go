package middleware

import (
	"net/http"

	"github.com/grainme/movie-api/internal/auth"
	"github.com/grainme/movie-api/internal/domain"
)

// at this point, we already used the auth middleware
// which means that we have "user" as key in r.Context
// therefor we can extract it and check the role ("admin" or "regular")
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value("user").(*auth.Claims)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if userRole.Role != string(domain.Admin) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
