package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
)

type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func GenerateRefreshToken() uuid.UUID {
	// for now, this is just a wrapper around
	// uuid.new() - better naming
	return uuid.New()
}

// should I use:
// - user domain.User as param
// - userId UUID, role string as params
// code design question? - it does not matter? - don't give a function more than it needs?
func GenerateAccessToken(userId uuid.UUID, role domain.Role) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(), // who the token is about?
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: string(role),
	}
	secretKey := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// This creates the signature using `HMAC(header + payload, secretKey)` and appends it:
	// eyJhbGc...header.eyJ1c2Vy...payload.SflKxw...signature
	// When the server validates, it recomputes the signature.
	// If an attacker modifies the payload, the signatures won't match → rejected.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
