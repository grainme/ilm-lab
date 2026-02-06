package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	moviePrefix        = "movie:"
	viewsMoviePrefix   = "views:movie:"
	refreshTokenPrefix = "refresh_token:"
)

func MovieKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", moviePrefix, id.String())
}

func ViewsMovieKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", viewsMoviePrefix, id.String())
}

func RefreshTokenKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", refreshTokenPrefix, id.String())
}
